package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rodellison/GoConchRepublicBackEnd/common"
	"github.com/rodellison/GoConchRepublicBackEnd/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var (
	testGoodHTML = "<!DOCTYPE html><html lang=\"en\">" +
		"<head><title>Official Florida Keys Tourism Council Calendar of Events</title></head>" +
		"<body id=\"top\" class=\"tdcus section section-florida-keys page-calendar unscrolled\" style=\"position: relative; min-height: 100%; top: 0px;\">" +
		"<div class=\"listing-block listing-calendar listing-calendar-img\" id=\"calendar-5117\">" +
		"<div class=\"listing-img\">" +
		"<div class=\"calendar-photos\" id=\"calendar-photos-5117\"><a class=\"swipebox expand-img\" href=\"/calendarofevents/img/5117-ev.jpg\"><span class=\"fa fa-expand\"></span><img src=\"/calendarofevents/img/5117-ev.jpg\" alt=\"Image for Key West Art in the Garden\"></a></div></div>" +
		"<ul class=\"ui\">" +
		"<li class=\"listing-info\"><span class=\"listing-date\">" +
		"<span class=\"fa fa-fw fa-calendar\"></span>Jun 1, 2020 - Jul 31, 2020</span>" +
		"<span class=\"listing-location\"><a href=\"/calendar/key-west/\">" +
		"<span class=\"fa fa-fw fa-map-marker\"></span>Location: Key West</a></span>" +
		"<span class=\"listing-type\"><a href=\"/calendar/arts-culture/\">" +
		"<span class=\"fa fa-fw fa-tags\"></span>Category: Arts &amp; Culture</a></span></li>" +
		"<li class=\"listing-name\"><a rel=\"nofollow\" href=\"https://www.keywest.garden/\" target=\"_blank\" title=\"View website\">Key West Art in the Garden</a></li>" +
		"<li class=\"listing-contact\">" +
		"<span class=\"listing-website\"><a href=\"https://www.keywest.garden/\" target=\"_blank\">" +
		"<span class=\"fa fa-fw fa-external-link\"></span>Website</a></span></li>" +
		"<li class=\"listing-desc\">The 10th annual Key West ART in the GARDEN opens June 1 at the Key West Tropical Forest &amp; Botanical Garden with artistic expressions emphasizing harmony with nature. These works are earth-friendly in the selection of materials (including recycled and natural material) and themes. Sculptures by local artists accent the natural beauty of the conservation site on Stock Island through July 31, 2020, 10am-4pm seven days a week.</li>" +
		"</ul>" +
		"</div> " +
		"<div class=\"listing-block listing-calendar\" id=\"calendar-4955\">" +
		"<ul class=\"ui\">" +
		"<li class=\"listing-info\">" +
		"<span class=\"listing-date\"><span class=\"fa fa-fw fa-calendar\"></span>Jun 2, 2020 - Jun 6, 2020</span>" +
		"<span class=\"listing-location\"><a href=\"/calendar/islamorada/\">" +
		"<span class=\"fa fa-fw fa-map-marker\"></span>Location: The Lower Keys</a></span>" +
		"<span class=\"listing-type\"><a href=\"/calendar/fishing/\">" +
		"<span class=\"fa fa-fw fa-tags\"></span>Category: Fishing</a></span></li>" +
		"<li class=\"listing-name\"><a rel=\"nofollow\" href=\"https://guidestrustfoundation.org/\" target=\"_blank\" title=\"View website\">46th Annual Don Hawley Tarpon Fly Tournament</a></li>" +
		"<li class=\"listing-contact\">" +
		"<span class=\"listing-website\"><a href=\"https://guidestrustfoundation.org/\" target=\"_blank\">" +
		"<span class=\"fa fa-fw fa-external-link\"></span>Website</a></span></li>" +
		"<li class=\"listing-desc\">Up to 25 of the world’s top fly-rod anglers endure a five-day test of patience and finesse, fishing Keys waters using fly tackle and 12-pound tippet. Named for the late fly fisherman and conservationist Don Hawley, the tournament benefits the Guides Trust Foundation, assisting professional fishing guides and supporting backcountry fishery conservation programs. </li>" +
		"</ul>" +
		"</div> " +
		"</body>" +
		"</html>"
)

func init() {
	//IMPORTANT!! - for the test to use our mocked response below, we have to make sure to set the client to
	//be the mocked client, which will use the overridden versions of the function that makes calls
//	common.TheHTTPClient = &mocks.MockHTTPClient{}
	common.TheHTTPClient = &http.Client{}   //use this if testing for real data fetch, not mocked.
	//common.EBIfaceClient = &mocks.MockEBSvcClient{}
	//common.SNSIfaceClient = &mocks.MockSNSSvcClient{}
	common.SQSIfaceClient = &mocks.MockSQSSvcClient{}

}


func TestHandlerCanProcessSingleGoodRequest(t *testing.T) {


	//For use with non-mocked HTTP handler
	expectedResult := "ConchRepublicBackend fetch responding successful!"

	tests := []struct {
		request aws.Context
		expect  string
		err     error
	}{
		{
			request: nil,
			expect:  expectedResult,
			err:     nil,
		},
	}


	mocks.MockDoSendSQSMessageWithContext = func(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		return &sqs.SendMessageOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Message)
	}
}


 /*

func TestHandlerCanProcessGoodRequest(t *testing.T) {

	expectedResult := "ConchRepublicBackend fetch responding successful!"

	tests := []struct {
		request context.Context
		expect  string
		err     error
	}{
		{
			request: context.Background(),
			expect:  expectedResult,
			err:     nil,
		},
	}

	// build response html
	// create a new reader with that html
	mocks.GetDoHTTPFunc = func(*http.Request) (*http.Response, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		r := ioutil.NopCloser(bytes.NewReader([]byte(testGoodHTML)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	//mocks.MockDoSNSPublishWithContext = func(ctx aws.Context, input *sns.PublishInput, options ...request.Option) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS Publish called")
	//	return &sns.PublishOutput{}, nil
	//}

	//Mock out the ctxhttp context sensitive http do function
	common.DoHTTPWithCTX = func(ctx aws.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(testGoodHTML)))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	mocks.MockDoSendSQSMessageWithContext = func(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		return &sqs.SendMessageOutput{}, nil
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Message)
	}
}
func TestHandlerCanProcessBadRequest(t *testing.T) {

	expectedResult := "ConchRepublicBackend fetch responding UNsuccessful!"

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

	// build response html
	// create a new reader with that html
	mocks.GetDoHTTPFunc = func(*http.Request) (*http.Response, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		r := ioutil.NopCloser(bytes.NewReader([]byte("")))
		return &http.Response{
			StatusCode: 500, //for this test, just using a bad return code to signify http get error
			Body:       r,
		}, nil
	}

	//mocks.MockDoSNSPublish = func(input *sns.PublishInput) (*sns.PublishOutput, error) {
	//	fmt.Println("Mock SNS Publish called")
	//	return &sns.PublishOutput{}, nil
	//}

	//Mock out the ctxhttp context sensitive http do function
	common.DoHTTPWithCTX = func(ctx aws.Context, client *http.Client, req *http.Request) (*http.Response, error) {
		r := ioutil.NopCloser(bytes.NewReader([]byte(testGoodHTML)))
		return &http.Response{
			StatusCode: 500, //for this test, just using a bad return code to signify http get error
			Body:       r,
		}, nil
	}

	mocks.MockDoSendSQSMessageWithContext = func(aws.Context, *sqs.SendMessageInput, ...request.Option) (*sqs.SendMessageOutput, error) {
		//Placing the NopCloser inside as EACH time the GetDoFunc is called the reader will be 'drained'
		return &sqs.SendMessageOutput{}, nil
	}


	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Message)
	}
}


 */


