package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type MockSNSSvcClient struct {
	snsiface.SNSAPI
}

var (
	MockDoSNSPublish            func(input *sns.PublishInput) (*sns.PublishOutput, error)
	MockDoSNSPublishWithContext func(ctx aws.Context, input *sns.PublishInput, options ...request.Option) (*sns.PublishOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test routines
func (s *MockSNSSvcClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return MockDoSNSPublish(input)
}
func (s *MockSNSSvcClient) PublishWithContext(ctx aws.Context, input *sns.PublishInput, options ...request.Option) (*sns.PublishOutput, error) {
	return MockDoSNSPublishWithContext(ctx, input)
}
