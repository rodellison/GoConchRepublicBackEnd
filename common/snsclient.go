package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	SNSIfaceClient snsiface.SNSAPI
	SNSSvcClient   *sns.SNS
)

func init() {

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the SNS service client, to be used for putting events
	SNSSvcClient = sns.New(sess)
	SNSIfaceClient = SNSSvcClient

	xray.AWS(SNSSvcClient.Client)

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func PublishSNSMessage(thisContext aws.Context, snsTopic, snsSubject, snsMessage string) (err error) {

	pubInput := &sns.PublishInput{
		Message:  aws.String(snsMessage),
		Subject:  aws.String(snsSubject),
		TopicArn: aws.String(snsTopic),
	}

	_, err = SNSIfaceClient.PublishWithContext(thisContext, pubInput)
	return err

}
