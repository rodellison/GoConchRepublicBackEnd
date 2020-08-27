package mocks

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSQSSvcClient struct {
	sqsiface.SQSAPI
}

var (
	MockDoDeleteSQSMessage func(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
)

//This is the mocked version of the real function
func (s *MockSQSSvcClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return MockDoDeleteSQSMessage(input)
}
