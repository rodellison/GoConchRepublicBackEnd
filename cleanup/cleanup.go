package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"log"
	"os"
)

var (
	StrFormattedDateToday string
)

func init() {

	StrFormattedDateToday = common.GetFormattedDateToday()

}

// AWS Lambda Proxy Request functionality (default behavior)
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {

	fmt.Println("ConchRepublic Cleanup invoked")
	success := true

	if ctx != nil {
		//If context info is needed to be used, uncomment next line and do something in the function
		//contextHandler(&ctx)
	}

	fmt.Println("ConchRepublic Cleanup begin database purge")
	var successMessage string
	if count, err := common.DeleteDBEvents(StrFormattedDateToday); err != nil {
		success = false
		successMessage = "ConchRepublic Cleanup did NOT complete successfully."
	} else {
		successMessage = fmt.Sprintf("ConchRepublic Cleanup complete. Counts: Total purged: %d", count)
		fmt.Println(successMessage)

		//Send an SNS message reporting results
		if err := common.PublishSNSMessage(os.Getenv("SNS_TOPIC"), "Conch Republic Cleanup", successMessage); err != nil {
			fmt.Println("Error sending SNS message: ", err.Error())
		}
	}

	return responseHandler(success)

}

func contextHandler(ctx *context.Context) {
	lc, _ := lambdacontext.FromContext(*ctx)
	log.Print(lc)
}

func responseHandler(success bool) (Response, error) {

	var returnString string
	if success {
		returnString = "ConchRepublicBackend cleanup responding successful!"
	} else {
		returnString = "ConchRepublicBackend cleanup responding UNsuccessful!"
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
			"X-MyCompany-Func-Reply": "cleanup-handler",
		},
	}

	return resp, nil
}
func main() {
	lambda.Start(Handler)
}
