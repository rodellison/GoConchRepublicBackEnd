package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	//IMPORTANT!! - for the test to use our mocked response below, we have to make sure to set the client to
	//be the mocked client, which will use the overridden versions of the function that makes calls
	//During testing, we'll override the endpoint to ensure testing against local DynamoDB Docker image
	cfg := aws.Config{
		Endpoint:   aws.String("http://localhost:8000"),
		Region:     aws.String("us-east-1"),
		MaxRetries: aws.Int(3),
	}

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create the eventbridge events service client, to be used for putting events
	common.DynamoDBSvcClient = dynamodb.New(sess, &cfg)

}

func TestHandlerCanInsertDynamoDBRequest(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend database responding successful!\"}"

	tests := []struct {
		request context.Context
		expect  string
		err     error
	}{
		{
			request: nil,
			expect:  expectedResult,
			err:     nil,
		},
	}

	evdata := []byte(`{ "EventID":"test-3", "StartDate":"20200606", "EndDate":"20200701", "EventName":"Test 3", "EventContact":"No Contact",
		"EventLocation":"Key West", "ImgURL":"http://someImgURL", "EventURL":"http://someEventURL", "EventDescription":"Test3-Description" }`)

	var testEvent = events.CloudWatchEvent{
		DetailType: "conchrepublicbackend.database",
		Source:     "goconchrepublicbackend.fetch",
		Time:       time.Now(),
		Detail:     evdata,
	}

	for _, test := range tests {
		response, err := Handler(test.request, &testEvent)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
