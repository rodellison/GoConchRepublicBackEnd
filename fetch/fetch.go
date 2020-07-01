package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"log"
	"os"
)

const eventbusStr = "conchrepublic"
const sourceStr = "conchrepublicbackend.fetch"
const detailTypeStr = "conchrepublicbackend.database"

type EventDetail struct {
	Month string
}

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
			fmt.Println("Returned from Notification Channel with Error")
			return responseHandler(false)
		} else {
			fmt.Println("Returned from Notification Channel")
			return responseHandler(true)
		}
	}
}

func fetch(month string, chFinished chan bool) {

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
			evdata := &common.Eventdata{" ", " ", " ", " ", " ", " ", " ", " ", " "}
			err := evdata.ExtractEventData(i, s)
			if err != nil {
				//There can be many events, let the error go, print it, but move on to the next item
				fmt.Println("error caught extracting event detail: " + err.Error())
			} else {

				detailStr, err := json.Marshal(evdata)
				if err != nil {
					//There can be many events, let the error go, print it, but move on to the next item
					fmt.Println("error caught extracting event detail: " + err.Error())
				} else {
					//Send an Event with the Event Details so the Database module can insert it.
					//Call the sendEvents function and handle error if it occurs
					if err := common.SendEBEvent(eventbusStr, sourceStr, detailTypeStr, string(detailStr)); err != nil {
						//For this modules case, dont fatal out on error, just move along
						fmt.Println(err.Error())
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
