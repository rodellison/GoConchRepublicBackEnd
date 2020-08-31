package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSQSSvcClient struct {
	sqsiface.SQSAPI
}

var (
	MockDoDeleteSQSMessage             func(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	MockDoReceiveSQSMessage            func(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	MockDoReceiveSQSMessageWithContext func(aws.Context, *sqs.ReceiveMessageInput, ...request.Option) (*sqs.ReceiveMessageOutput, error)
	MockDoSendSQSMessageWithContext    func(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error)
	MockDoSendSQSMessage               func(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
)

//This is the mocked version of the real function
func (s *MockSQSSvcClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return MockDoDeleteSQSMessage(input)
}

func (s *MockSQSSvcClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return MockDoReceiveSQSMessage(input)
}
func (s *MockSQSSvcClient) ReceiveMessageWithContext(ctx aws.Context, input *sqs.ReceiveMessageInput, options ...request.Option) (*sqs.ReceiveMessageOutput, error) {
	return MockDoReceiveSQSMessageWithContext(ctx, input)
}
func (s *MockSQSSvcClient) SendMessageWithContext(ctx aws.Context, input *sqs.SendMessageInput, options ...request.Option) (*sqs.SendMessageOutput, error) {
	return MockDoSendSQSMessageWithContext(ctx, input)
}
