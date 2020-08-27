package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"github.com/rodellison/GoConchRepublicBackEnd/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	//IMPORTANT!! - for the test to use our mocked response below, we have to make sure to set the client to
	//be the mocked client(s), which will use the overridden versions of the function that makes calls
	common.EBSvcClient = &mocks.MockEBSvcClient{}
	common.SNSSvcClient = &mocks.MockSNSSvcClient{}
}

func TestInitiateHandlerSuccess(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend initiate responding successful!\"}"

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
	mocks.MockDoPutEvent = func(input *eventbridge.PutEventsInput) (*eventbridge.PutEventsOutput, error) {
		fmt.Println("Mock PutEvents called")
		return &eventbridge.PutEventsOutput{}, nil
	}

	mocks.MockDoSNSPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func TestInitiateHandlerEventPubFail(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend initiate responding UNsuccessful!\"}"

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
	mocks.MockDoPutEvent = func(input *eventbridge.PutEventsInput) (*eventbridge.PutEventsOutput, error) {
		fmt.Println("Mock PutEvents called")
		return &eventbridge.PutEventsOutput{}, errors.New("Could not send Event")
	}

	mocks.MockDoSNSPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}

func TestInitiateHandlerSNSPubFail(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend initiate responding UNsuccessful!\"}"

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
	mocks.MockDoPutEvent = func(input *eventbridge.PutEventsInput) (*eventbridge.PutEventsOutput, error) {
		fmt.Println("Mock PutEvents called")
		return &eventbridge.PutEventsOutput{}, nil
	}

	mocks.MockDoSNSPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, errors.New("Could not send SMS")
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
