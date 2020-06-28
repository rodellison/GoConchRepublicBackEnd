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
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, theEvent *events.CloudWatchEvent) (Response, error) {

	fmt.Println("ConchRepublic Database invoked ")
	success := true
	if ctx != nil {
		//If context info is needed to be used, uncomment next line and do something in the function
		//contextHandler(&ctx)
	}

	//Unmarshal the incoming JSON Event detail attribute contents into an Eventdata struct
	var inInterface common.Eventdata
	err := json.Unmarshal([]byte(theEvent.Detail), &inInterface)
	if err != nil {
		fmt.Println("Error during Unmarshal of incoming Event Detail data")
		success = false
	} else {
		dberr := common.InsertDBEvent(inInterface)
		if dberr != nil {
			fmt.Println("Error occurred inserting InsertDBEvent")
			success = false
		}
	}

	fmt.Println("ConchRepublic database processing completed.")
	return responseHandler(success)

}

func responseHandler(success bool) (Response, error) {

	var returnString string
	if success {
		returnString = "ConchRepublicBackend database responding successful!"
	} else {
		returnString = "ConchRepublicBackend database responding UNsuccessful!"
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

func contextHandler(ctx *context.Context) {
	lc, _ := lambdacontext.FromContext(*ctx)
	log.Print(lc)
}

func main() {
	lambda.Start(Handler)
}
