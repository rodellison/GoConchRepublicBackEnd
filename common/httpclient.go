package common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//An HTTPClient interface declared to allow for easier Mock testing
//Basically, ensure the custom interface has definitions for the functions that need to be mocked (so
//as to not make 'real' requests
//Establish a Variable that is of the interface's type that can be used to hold the 'real' client (when
//not running tests, as well as be a variable we can set during 'test' time
//And setup an init() function that sets the variable to the 'real' interface as it's default when not testing

type HTTPClient interface {
	Do (req *http.Request) (*http.Response, error)
}
var (
	TheHTTPClient HTTPClient
)
func init() {
	//IF we're running a test, we'll swap this variable's value to use a mock instead, but when not
	//testing, the value will be preset to ensure that it uses the 'real' httpClient interface
	TheHTTPClient = &http.Client{}
}

//Post sends a post request to the URL with the body
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
	//Client := &http.Client{}  using the variable/interface above to facilitate easier mock testing later
	return TheHTTPClient.Do(request)
}

//// Post sends a post request to the URL with the body
//func PostURL(url string, body interface{}, headers http.Header) (*http.Response, error) {
//	jsonBytes, err := json.Marshal(body)
//	if err != nil {
//		return nil, err
//	}
//	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
//	if err != nil {
//		return nil, err
//	}
//	request.Header = headers
////	client := &http.Client{}
//	return TheHTTPClient.Do(request)
//}