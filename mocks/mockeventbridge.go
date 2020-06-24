package mocks

import (
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
)

type MockEBSvcClient struct {
	eventbridgeiface.EventBridgeAPI
}

var (
	MockDoPutEvent func(input *eventbridge.PutEventsInput) (*eventbridge.PutEventsOutput, error)
)

//This is the mocked version of the real function
//It returns the variable above, which is a function that can be overloaded in our test functions
func (m *MockEBSvcClient) PutEvents(input *eventbridge.PutEventsInput) (*eventbridge.PutEventsOutput, error){
	return MockDoPutEvent(input)
}

