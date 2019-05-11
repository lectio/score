package score

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTPUserAgent may be passed into getHTTPResult as the default HTTP User-Agent header parameter
const HTTPUserAgent = "github.com/lectio/score"

// HTTPTimeout may be passed into getHTTPResult function as the default HTTP timeout parameter
const HTTPTimeout = time.Second * 90

// HTTPResult encapsulates an API call
type httpResult struct {
	apiEndpoint string
	body        *[]byte
}

// GetHTTPResult runs the apiEndpoint and returns the body of the HTTP result
// TODO: Consider using [HTTP Cache](https://github.com/gregjones/httpcache)
func getHTTPResult(apiEndpoint string, client *http.Client, userAgent string) (*httpResult, Issue) {
	result := new(httpResult)
	result.apiEndpoint = apiEndpoint

	req, reqErr := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if reqErr != nil {
		return nil, NewIssue(apiEndpoint, UnableToCreateHTTPRequest, fmt.Sprintf("Unable to create HTTP request: %v", reqErr), true)
	}
	req.Header.Set("User-Agent", userAgent)
	resp, getErr := client.Do(req)
	if getErr != nil {
		return nil, NewIssue(apiEndpoint, UnableToExecuteHTTPGETRequest, fmt.Sprintf("Unable to execute HTTP GET request: %v", getErr), true)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NewHTTPResponseIssue(apiEndpoint, resp.StatusCode, fmt.Sprintf("HTTP response status is not 200: %v", resp.StatusCode), true)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, NewIssue(apiEndpoint, UnableToReadBodyFromHTTPResponse, fmt.Sprintf("Unable to read body from HTTP response: %v", readErr), true)
	}

	result.body = &body
	return result, nil
}
