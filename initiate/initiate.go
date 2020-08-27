package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"os"
	"time"
)

const eventbusStr = "conchrepublic"
const sourceStr = "conchrepublicbackend.initiate"
const detailTypeStr = "conchrepublicbackend.fetch"

// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {

	fmt.Println("ConchRepublic Initiate invoked")

	fmt.Println("ConchRepublic Initiate begin sending events..")
	//This creates 12 events, one for each month
	//If there are ANY failures sending events, then return false from this function and get out, otherwise return true
	for i := 1; i <= 12; i++ {
		detailStr := fmt.Sprintf("{ \"month\": \"%d\" }", i)

		//Call the sendEvents function and handle error if it occurs
		if err := common.SendEBEvent(eventbusStr, sourceStr, detailTypeStr, detailStr); err != nil {
			return responseHandler(false)
		}

		//NOTE: Sleep is being called here only to act as a throttler, to limit the number of lambda instances
		//that will started up to receive and process the events being put from this module.
		//Just doing this since we don't really need 12 parallel lambda instances running (each having its own cloudwatch file).
		time.Sleep(100 * time.Millisecond)

	}
	fmt.Println("ConchRepublic initiate send events completed.")

	if err := common.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Conch Republic Initiate", "Conch Republic Backend process initiated."); err != nil {
		fmt.Println("Error sending SNS message: ", err.Error())
		return responseHandler(false)
	}

	return responseHandler(true)
}

func responseHandler(success bool) (Response, error) {

	var returnString string
	if success {
		returnString = "ConchRepublicBackend initiate responding successful!"
	} else {
		returnString = "ConchRepublicBackend initiate responding UNsuccessful!"
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
			"X-MyCompany-Func-Reply": "initiate-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
