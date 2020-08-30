package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type Response events.APIGatewayProxyResponse

type sqsConsumer struct {
	QueueURL          string
	maxMessages       int64
	maxWaitSeconds    int64
	visibilityTimeout int64
}

var (
	mySQSConsumer sqsConsumer
	itemCount     uint64
)

func init() {

	mySQSConsumer = sqsConsumer{
		QueueURL:          os.Getenv("SQS_TOPIC"),
		maxMessages:       10,
		maxWaitSeconds:    10,
		visibilityTimeout: 30,
	}

}

func (c *sqsConsumer) consumeAndProcess() error {
	itemCount = 0
	var wg sync.WaitGroup

	for {
		output, err := common.SQSSvcClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			MaxNumberOfMessages: &c.maxMessages,
			QueueUrl:            &c.QueueURL,
			WaitTimeSeconds:     &c.maxWaitSeconds,
			VisibilityTimeout:   &c.visibilityTimeout,
		})
		if err != nil {
			return err
		}
		if len(output.Messages) > 0 {

			fmt.Println("This loop is processing " + strconv.Itoa(len(output.Messages)) + " messages.")
			wg.Add(len(output.Messages))
			for _, message := range output.Messages {
				go func(m *sqs.Message) {
					defer wg.Done()
					var theEvent common.Eventdata
					messagebodyBytes := []byte(*m.Body)

					if err := json.Unmarshal(messagebodyBytes, &theEvent); err != nil {
						panic(err)
					}

					fmt.Println("SQS Message Item: " + *m.MessageId + ", with EventID: " + theEvent.EventID)
					dberr := common.InsertDBEvent(theEvent)
					if dberr != nil {
						fmt.Println("Error occurred inserting Data via InsertDBEvent")
					} else {
						atomic.AddUint64(&itemCount, 1)
						//If we inserted the Event, then Delete the SQS message
						_, err := common.SQSSvcClient.DeleteMessage(&sqs.DeleteMessageInput{
							QueueUrl:      &c.QueueURL,
							ReceiptHandle: message.ReceiptHandle,
						}) //MESSAGE CONSUMED
						if err != nil {
							fmt.Println("Error deleting SQS message")
						}
					}
				}(message)
			}
			wg.Wait()
		} else {
			//There are no more items for this worker so get out
			fmt.Println("No additional messages to process")
			break
		}
	}

	return nil
}

func Handler(ctx context.Context) (Response, error) {

	fmt.Println("ConchRepublic Database starting...")
	success := true

	//This calls the main process to process SQS messages and perform a DB insert for each message/item received
	err := mySQSConsumer.consumeAndProcess()
	if err != nil {
		fmt.Println("Error receiving messagage from SQS: " + err.Error())
		success = false
	} else {
		if itemCount > 0 {
			snsBody := "Conch Republic Backend process completed. Count of items processed: " + strconv.FormatUint(itemCount, 10)
			fmt.Println(snsBody)
			err := common.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Conch Republic Database", snsBody)
			if err != nil {
				fmt.Println("Error sending SNS message: ", err.Error())
			}
		}
	}

	fmt.Println("ConchRepublic database processing completed.")
	return responseHandler(success)

}

func responseHandler(success bool) (Response, error) {

	var returnString string
	if success {
		returnString = "ConchRepublicBackend database processed successful!"
	} else {
		returnString = "ConchRepublicBackend database processed UNsuccessful!"
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": returnString,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "database-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
