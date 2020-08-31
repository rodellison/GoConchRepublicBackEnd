package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type MockDynamoDBSvcClient struct {
	dynamodbiface.DynamoDBAPI
}

var (
	MockDynamoScan               func(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
	MockDynamoPutItem            func(*dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	MockDynamoPutItemWithContext func(aws.Context, *dynamodb.PutItemInput, ...request.Option) (*dynamodb.PutItemOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test routines
func (m *MockDynamoDBSvcClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	return MockDynamoScan(input)
}
func (m *MockDynamoDBSvcClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return MockDynamoPutItem(input)
}
func (m *MockDynamoDBSvcClient) PutItemWithContext(context aws.Context, input *dynamodb.PutItemInput, options ...request.Option) (*dynamodb.PutItemOutput, error) {
	return MockDynamoPutItemWithContext(context, input)
}
