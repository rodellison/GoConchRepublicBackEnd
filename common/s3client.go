package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
)

var (
	S3SvcClient   s3iface.S3API
	uploader      *s3manager.Uploader
	PublishS3Func func(string) error
)

func init() {

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the eventbridge events service client, to be used for putting events
	S3SvcClient = s3.New(sess)
	uploader = s3manager.NewUploader(sess)
	PublishS3Func = PublishS3JSONFile

}

// func PublishSNSMessage uses an SDK service client to send an SNS Publish request
func PublishS3JSONFile(s3JsonEventData string) (err error) {

	//Create a reader that can be passed into the uploader process, with the Json Event data string as the source
	r := strings.NewReader(s3JsonEventData)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String("ConchRepublic/ConchRepublicEvents.json"),
		Body:   r,
	})
	if err != nil {
		return err
	}

	return nil
}
