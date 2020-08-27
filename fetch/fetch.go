package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"log"
	"os"
	"strconv"
	"sync/atomic"
)

type EventDetail struct {
	Month string
}

var (
	itemCount uint64
)

type Response events.APIGatewayProxyResponse

func Handler(ctx context.Context, theEvent *events.CloudWatchEvent) (Response, error) {

	fmt.Println("ConchRepublic Fetch invoked ")

	var thisEventsDetail EventDetail
	err := json.Unmarshal(theEvent.Detail, &thisEventsDetail)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Event detail:", thisEventsDetail.Month)

	chFinished := make(chan bool)
	defer func() {
		close(chFinished)
	}()

	go fetch(thisEventsDetail.Month, chFinished)

	// Subscribe to both channels
	select {
	case value := <-chFinished:
		//Its possible that the fetch of data or extraction of data fetched did not occur successfully
		if value != true {
			fmt.Println("Returned Error")
			return responseHandler(false)
		} else {
			fmt.Println("Returned Successfully. Items processed: " + strconv.FormatUint(itemCount, 10))
			return responseHandler(true)
		}
	}
}

func fetch(month string, chFinished chan bool) {
	itemCount = 0

	fullURL := os.Getenv("URLBASE") + os.Getenv("URLBASE2") + common.CalcSearchYYYYMMFromDate(month)
	fmt.Println(fullURL)
	//Use this to provide a return value to pass back through the channel once this routine is finished.
	//Default it to True, and return false only if any critical errors occur that should not allow us to proceed
	//to downstream processing
	returnValue := true

	//http get logic is in a common httpclient file
	if resp, err := common.GetURL(fullURL); err != nil || resp.StatusCode != 200 {
		//This is critical, return immediately and don't try to process anything further
		fmt.Println("ERROR: Failed to fetch:", fullURL)
		resp.Body.Close()
		chFinished <- false
	} else {
		defer func() {
			resp.Body.Close() // close Body when the function completes
			// Notify that we're done after this function
			chFinished <- returnValue
		}()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println("ERROR: Failed to get HTML from response body:", err.Error())
			returnValue = false
		}

		//	events := make(map[string]eventdata, 100)
		//Calendar entries will have one of these main listing block constants, depending on whether they
		//contain an image vs not contain image

		doc.Find(common.LISTING_BLOCK).Each(func(i int, s *goquery.Selection) {
			evdata := &common.Eventdata{" ", " ", " ", " ", " ", " ", " ", " ", " ", 0}
			err := evdata.ExtractEventData(i, s)
			if err != nil {
				//There can be many events, let the error go, print it, but move on to the next item
				fmt.Println("error caught extracting event detail: " + err.Error())
			} else {
				//Setup the event's expiry value - int64 epoch value that DynamoDB can use for automated record removal
				expYYYY, _ := strconv.Atoi(evdata.EndDate[0:4])
				expMM, _ := strconv.Atoi(evdata.EndDate[4:6])
				expDD, _ := strconv.Atoi(evdata.EndDate[6:8])
				evdata.EventExpiry = common.CalcLongEpochFromEndDate(expYYYY, expMM, expDD)

				detailStr, err := json.Marshal(evdata)
				if err != nil {
					//There can be many events, let the error go, print it, but move on to the next item
					fmt.Println("error caught extracting event detail: " + err.Error())
				} else {
					//Send an SQS Message with the Event Details so the Database module that will run later can poll/insert it.
					_, err := common.SQSSvcClient.SendMessage(&sqs.SendMessageInput{
						MessageBody: aws.String(string(detailStr)),
						QueueUrl:    aws.String(os.Getenv("SQS_TOPIC")),
					})
					if err != nil {
						fmt.Println("Error sending SQS message: ", err.Error())
					} else {
						atomic.AddUint64(&itemCount, 1)
					}

				}

			}
		})
	}
}

func responseHandler(success bool) (Response, error) {

	var returnString string
	if success {
		returnString = "ConchRepublicBackend fetch responding successful!"
	} else {
		returnString = "ConchRepublicBackend fetch responding UNsuccessful!"
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": returnString,
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "fetch-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
