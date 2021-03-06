package common

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"

	//	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"
)

var (
	DynamoDBIfaceClient dynamodbiface.DynamoDBAPI
	DynamoDBSvcClient   *dynamodb.DynamoDB
	TableName           string
)

func init() {

	//During testing, we'll override the endpoint to ensure testing against local DynamoDB Docker image
	cfg := aws.Config{
		//		Endpoint: aws.String("http://localhost:8000"),
		Region:     aws.String("us-east-1"),
		MaxRetries: aws.Int(3),
	}

	//Get Session, credentials
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create the DynamoDB service client, to be used for inserting DB entries
	DynamoDBSvcClient = dynamodb.New(sess, &cfg)
	DynamoDBIfaceClient = DynamoDBSvcClient

	//Making the Tablename an environmental variable so that it can be changed outside of program
	TableName = os.Getenv("DYNAMO_DB_TABLENAME")

	//Note: XRay is unnecessary - but using it to try out tracing services..
	xray.AWS(DynamoDBSvcClient.Client)

}

// func InsertDBEvent converts Eventdata into appropriate DynamoDB table attributes, and puts the item into the DB.
func InsertDBEvent(ctx aws.Context, data Eventdata) (err error) {

	//First, Marshal the incoming EventItem JSON string data into a DynamoDB attribute map
	evItem, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		fmt.Println("Error occurred during marshalling new Eventdata item: ", err.Error())
		return err
	}

	//	fmt.Println(evItem)

	_, err = DynamoDBIfaceClient.PutItemWithContext(ctx, &dynamodb.PutItemInput{
		//	_, err = DynamoDBIfaceClient.PutItem(&dynamodb.PutItemInput{
		Item:      evItem,
		TableName: &TableName,
	})

	if err != nil {
		fmt.Println(err)
		return err
	} else {
		return nil
	}

}
