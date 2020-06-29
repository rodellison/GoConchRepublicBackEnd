package main

import (
	"context"
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

func TestInitiateHandler(t *testing.T) {

	expectedResult := "{\"message\":\"ConchRepublicBackend initiate responding successfully!\"}"

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
	// build response from mocked EventBridge PutEvents call
	mocks.MockDoPublishEvent = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
		fmt.Println("Mock SNS Publish called")
		return &sns.PublishOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}
}
