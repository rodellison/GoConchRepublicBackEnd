package main

import (

	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type ResponseOutput struct {
	Message string   `json:"message"`
}

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

func (c *sqsConsumer) consumeAndProcess(ctx aws.Context) error {
	itemCount = 0
	var wg sync.WaitGroup

	for {
		output, err := common.ReceiveSQSMessages(ctx, &c.QueueURL, &c.maxMessages, &c.maxWaitSeconds, &c.visibilityTimeout)
		if err != nil {
			return err
		}
		if output != nil && len(output.Messages) > 0 {

			fmt.Println("This loop is processing " + strconv.Itoa(len(output.Messages)) + " messages.")
			wg.Add(len(output.Messages))
			for _, message := range output.Messages {
				go func(c *sqsConsumer, m *sqs.Message) {
					defer wg.Done()
					var theEvent common.Eventdata
					messagebodyBytes := []byte(*m.Body)

					if err := json.Unmarshal(messagebodyBytes, &theEvent); err != nil {
						panic(err)
					}

					fmt.Println("SQS Message Item: " + *m.MessageId + ", with EventID: " + theEvent.EventID)
					dberr := common.InsertDBEvent(ctx, theEvent)
					if dberr != nil {
						fmt.Println("Error occurred inserting Data via InsertDBEvent")
					} else {
						atomic.AddUint64(&itemCount, 1)
						//If we inserted the Event, then Delete the SQS message
						err := common.DeleteSQSMessage(ctx, &c.QueueURL, m.ReceiptHandle) //MESSAGE CONSUMED
						if err != nil {
							fmt.Println("Error deleting SQS message")
						}
					}
				}(c, message)
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

func Handler(ctx aws.Context) (ResponseOutput, error) {
	xray.Configure(xray.Config{LogLevel: "trace"})
	fmt.Println("ConchRepublic Database starting...")

	//This calls the main process to process SQS messages and perform a DB insert for each message/item received
	err := mySQSConsumer.consumeAndProcess(ctx)
	if err != nil {
		fmt.Println("Error receiving messagage from SQS: " + err.Error())
		return responseHandler(false, "ConchRepublicBackend database processed UNsuccessful!")
	} else {
		if itemCount > 0 {
			return responseHandler(true, "ConchRepublicBackend database processed successful! Count of items processed: " + strconv.FormatUint(itemCount, 10))
		} else {
			return responseHandler(true, "ConchRepublicBackend database processed successful! There were no items to process.")
		}
	}
}

func responseHandler(success bool, message string) (ResponseOutput, error) {

	fmt.Println("ConchRepublic database processing completed.")
	return ResponseOutput{
		Message: message,
	}, nil


}

func main() {
	lambda.Start(Handler)
}
