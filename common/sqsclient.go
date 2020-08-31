package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"
)

var (
	SQSIfaceClient sqsiface.SQSAPI
	SQSSvcClient   *sqs.SQS
)

func init() {

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the eventbridge events service client, to be used for putting events
	SQSSvcClient = sqs.New(sess)
	SQSIfaceClient = SQSSvcClient

	xray.AWS(SQSSvcClient.Client)

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func DeleteSQSMessage(thisContext aws.Context, queue *string, receiptHandle *string) (err error) {

	deleteInput := &sqs.DeleteMessageInput{
		QueueUrl:      queue,
		ReceiptHandle: receiptHandle,
	}

	_, err = SQSIfaceClient.DeleteMessageWithContext(thisContext, deleteInput)
	if err != nil {
		return err
	} else {
		return nil
	}

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func ReceiveSQSMessages(thisContext aws.Context, queueURL *string, maxMessages, maxWaitSeconds, visibilityTimeout *int64) (output *sqs.ReceiveMessageOutput, err error) {

	recieveMessageInput := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: maxMessages,
		QueueUrl:            queueURL,
		WaitTimeSeconds:     maxWaitSeconds,
		VisibilityTimeout:   visibilityTimeout,
	}

	return SQSIfaceClient.ReceiveMessageWithContext(thisContext, recieveMessageInput)

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func SendSQSMessage(thisContext aws.Context, detailString string) (output *sqs.SendMessageOutput, err error) {

	SendMessageInput := &sqs.SendMessageInput{
		MessageBody: aws.String(string(detailString)),
		QueueUrl:    aws.String(os.Getenv("SQS_TOPIC")),
	}
	return SQSIfaceClient.SendMessageWithContext(thisContext, SendMessageInput)

}
