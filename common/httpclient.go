package common

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-xray-sdk-go/xray"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
)

//An HTTPClient interface declared to allow for easier Mock testing
//Basically, ensure the custom interface has definitions for the functions that need to be mocked (so
//as to not make 'real' requests
//Establish a Variable that is of the interface's type that can be used to hold the 'real' client (when
//not running tests, as well as be a variable we can set during 'test' time
//And setup an init() function that sets the variable to the 'real' interface as it's default when not testing

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)

}

var (
	TheHTTPClient HTTPClient
	ctxClient *http.Client
	DoHTTPWithCTX func(ctx aws.Context, client *http.Client, req *http.Request) (*http.Response, error)
)

func init() {
	//IF we're running a test, we'll swap this variable's value to use a mock instead, but when not
	//testing, the value will be preset to ensure that it uses the 'real' httpClient interface
    ctxClient =  xray.Client(&http.Client{})
	TheHTTPClient = ctxClient
	DoHTTPWithCTX = ctxhttp.Do

}

//func GetURL fetches raw HTML data from the input url.. essentially a screen-scrape
func GetURL(url string) (*http.Response, error) {
	//Empty body for now
	jsonBytes, err := json.Marshal("")
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return TheHTTPClient.Do(request)
}
//func GetURL fetches raw HTML data from the input url.. essentially a screen-scrape
func GetURLWithContext(ctx aws.Context, url string) (*http.Response, error) {
	//Empty body for now
	jsonBytes, err := json.Marshal("")
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	return DoHTTPWithCTX(ctx, ctxClient, request)
}

