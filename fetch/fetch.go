package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	itemCount uint64
)

type ResponseOutput struct {
	Message string   `json:"message"`
}

func Handler(ctx aws.Context) (ResponseOutput, error) {
	xray.Configure(xray.Config{LogLevel: "trace"})

	fmt.Println("ConchRepublic Initiate Fetch invoked ")

	chFinished := make(chan bool)
	defer func() {
		close(chFinished)
	}()

	if err := common.PublishSNSMessage(ctx, os.Getenv("SNS_TOPIC"), "Conch Republic Initiate Fetch", "Conch Republic Backend process initiated."); err != nil {
		fmt.Println("Error sending SNS message: ", err.Error())
		return responseHandler(false)
	}

	go fetch(ctx, chFinished)

	// Subscribe to both channels
	select {
	case value := <-chFinished:
		//Its possible that the fetch of data or extraction of data fetched did not occur successfully
		if value != true {
			fmt.Println("Did NOT return successfully. Check errors above.")
			return responseHandler(false)
		} else {
			fmt.Println("Returned Successfully. Items processed: " + strconv.FormatUint(itemCount, 10))
			return responseHandler(true)
		}
	}
}

func fetch(ctx aws.Context, chFinished chan bool) {
	itemCount = 0

	var wg sync.WaitGroup
	returnVal := true

	for monthVal := 1; monthVal <= 12; monthVal++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, month int) {
			fullURL := os.Getenv("URLBASE") + os.Getenv("URLBASE2") + common.CalcSearchYYYYMMFromDate(month)
			fmt.Println("Attempting to Fetch URL: " + fullURL)
			resp, err := common.GetURLWithContext(ctx, fullURL)
			if err != nil || resp.StatusCode != 200 {
				//This is critical, return immediately and don't try to process anything further
				if err == nil {
					fmt.Println("ERROR: Failed to fetch URL:" + fullURL + ", Response code: " + strconv.Itoa(resp.StatusCode))
				} else {
					fmt.Println("ERROR: Failed to fetch URL:" + fullURL + ", Response code: " + strconv.Itoa(resp.StatusCode) + ", Error: " + err.Error())
				}

				resp.Body.Close()
				returnVal = false
				wg.Done()
			} else {
				defer func() {
					resp.Body.Close() // close Body when the function completes
					// Notify that we're done after this function
					wg.Done()
				}()

				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					fmt.Println("ERROR: Failed to get HTML from response body:", err.Error())
					returnVal = false
				} else {
					doc.Find(common.LISTING_BLOCK).Each(func(i int, s *goquery.Selection) {
						evdata := &common.Eventdata{" ", " ", " ", " ", " ", " ", " ", " ", " ", 0}
						err := evdata.ExtractEventData(s)
						if err != nil {
							//There can be many events, let the error go, print it, but move on to the next item
							fmt.Println("Error caught extracting event detail: " + err.Error())
							//chFinished <- false
						} else {
							//Setup the event's expiry value - int64 epoch value that DynamoDB can use for automated record removal
							expYYYY, _ := strconv.Atoi(evdata.EndDate[0:4])
							expMM, _ := strconv.Atoi(evdata.EndDate[4:6])
							expDD, _ := strconv.Atoi(evdata.EndDate[6:8])
							evdata.EventExpiry = common.CalcLongEpochFromEndDate(expYYYY, expMM, expDD)

							detailStr, err := json.Marshal(evdata)
							if err != nil {
								//There can be many events, let the error go, print it, but move on to the next item
								fmt.Println("Error caught extracting event detail: " + err.Error())
								//chFinished <- false
							} else {
								//Send an SQS Message with the Event Details so the Database module that will run later can poll/insert it.
								//End if an error occurs as an AWS service issue is something we want to know
								_, err := common.SendSQSMessage(ctx, string(detailStr))
								if err != nil {
									fmt.Println("Error sending SQS message: ", err.Error())
									returnVal = false
								} else {
									atomic.AddUint64(&itemCount, 1)
								}
							}

						}
					})
				}
			}

		} (&wg, monthVal)
	}
	wg.Wait()
	chFinished <- returnVal

}

func responseHandler(success bool) (ResponseOutput, error) {

	if !success {
		return ResponseOutput{
			Message: "ConchRepublicBackend fetch responding UNsuccessful!",
		}, nil
	} else
	{
		return ResponseOutput{
			Message: "ConchRepublicBackend fetch responding successful!",
		}, nil
	}

}

func main() {
	lambda.Start(Handler)
}
