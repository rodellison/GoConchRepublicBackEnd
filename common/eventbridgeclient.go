package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/aws/aws-sdk-go/service/eventbridge/eventbridgeiface"
	"time"
)

var (
	EBSvcClient eventbridgeiface.EventBridgeAPI
)

func init() {

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the eventbridge events service client, to be used for putting events
	EBSvcClient = eventbridge.New(sess)

}

// func sendEvent uses an SDK service client to make a request to Amazon EventBridge.
func SendEBEvent(eventbusStr, sourceStr, detailTypeStr, detailStr string) (err error) {
	fmt.Println("Sending Event with detailStr:", detailStr)
	_, err = EBSvcClient.PutEvents(&eventbridge.PutEventsInput{
		Entries: []*eventbridge.PutEventsRequestEntry{
			{
				EventBusName: aws.String(eventbusStr),
				Source:       aws.String(sourceStr),
				DetailType:   aws.String(detailTypeStr),
				Detail:       aws.String(detailStr),
				Time:         aws.Time(time.Now()),
			},
		},
	})

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}

}
