package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
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

	common.DynamoDBSvcClient = dynamodb.New(sess, &cfg) //use this one for actual dB interaction - test or prod
	//	common.DynamoDBSvcClient = &mocks.MockDynamoDBSvcClient{}
	//common.SNSIfaceClient = &mocks.MockSNSSvcClient{}
	common.SQSIfaceClient = &mocks.MockSQSSvcClient{}

}

func TestHandlerCanInsertDynamoDBRequest(t *testing.T) {

	expectedResult := "ConchRepublicBackend database processed successful!"

	tests := []struct {
		request aws.Context
		expect  string
		err     error
	}{
		{
			request: nil,
			expect:  expectedResult,
			err:     nil,
		},
	}

	mocks.MockDoDeleteSQSMessage = func(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
		//Creating a mock for Delete as we don't really want to delete any items
		return &sqs.DeleteMessageOutput{}, nil
	}

	// build mock DynamoDB put
	//mocks.MockDynamoPutItem = func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	//	//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
	//	fmt.Println("Mock DynamoDB Put Item called")
	//	return &dynamodb.PutItemOutput{
	//		Attributes:            nil,
	//		ConsumedCapacity:      nil,
	//		ItemCollectionMetrics: nil,
	//	}, nil
	//}

	//mocks.MockDoSNSPublish = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS Publish called with info: " + *input.Message)
	//	return &sns.PublishOutput{}, nil
	//}
	//mocks.MockDoSNSPublishWithContext = func(ctx aws.Context, input *sns.PublishInput, options ...request.Option) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS PublishWithContext called with info: " + *input.Message)
	//	return &sns.PublishOutput{}, nil
	//}
	mocks.MockDoReceiveSQSMessageWithContext = func(ctx aws.Context, input *sqs.ReceiveMessageInput, options ...request.Option) (*sqs.ReceiveMessageOutput, error) {
		fmt.Println("Mock SQS ReceiveMessageWithontext called with info: ")
		return &sqs.ReceiveMessageOutput{Messages: nil}, nil
	}
	mocks.MockDoReceiveSQSMessage = func(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
		//Creating a mock for Delete as we don't really want to delete any items
		return &sqs.ReceiveMessageOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Contains(t, response.Message, test.expect )
	}

}
