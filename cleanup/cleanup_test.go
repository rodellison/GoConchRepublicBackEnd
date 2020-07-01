package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"github.com/rodellison/GoConchRepublicBackEnd/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
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

	// Create the DynamoDB service client, to be used for deleting events
	common.DynamoDBSvcClient = dynamodb.New(sess, &cfg)

	// Create the sns publish service client, to be used for publishing SNS messages
	common.SNSSvcClient = &mocks.MockSNSSvcClient{}

}

func TestCleanupHandler(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend cleanup responding successful!\"}"

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

	// build response from mocked EventBridge PutEvents call
	mocks.MockDoPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, nil
	}

	//Override the Cleanup Handlers variable depending on needed result..
	//a far future date will ensure ALL of the DB items are selected..
	StrFormattedDateToday = "20300101"

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func TestCleanupHandlerSMSError(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend cleanup responding UNsuccessful!\"}"

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

	// build response from mocked EventBridge PutEvents call
	mocks.MockDoPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, errors.New("Error to test SMS error logic")
	}

	//A bad date used to test DB Delete error logic
	//Note, this text itself will not trip up DynamoDB Scan logic, which would just return 0 records..
	StrFormattedDateToday = "20300101"

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func TestCleanupHandlerDBDeleteError(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend cleanup responding UNsuccessful!\"}"

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

	// build response from mocked EventBridge PutEvents call
	mocks.MockDoPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, nil
	}

	//For this test override the DynamoDBSvcClient to be our Mocked version, so as to test handling Scan error logic
	common.DynamoDBSvcClient = &mocks.MockDynamoDBSvcClient{}
	mocks.MockDynamoScan = func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
		fmt.Println("Mock DynamoDB Scan error invoked")
		return &dynamodb.ScanOutput{}, errors.New("A DynamoDB Scan error occurred")
	}
	//A bad date used to test DB Delete error logic
	//Note, this text itself will not trip up DynamoDB Scan logic, which would just return 0 records..
	StrFormattedDateToday = "?*^%$^%$"

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

