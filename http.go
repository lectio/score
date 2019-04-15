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
func getHTTPResult(apiEndpoint string, userAgent string, timeout time.Duration) (*httpResult, error) {
	result := new(httpResult)
	result.apiEndpoint = apiEndpoint

	httpClient := http.Client{
		Timeout: timeout,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("Unable to create request %q: %v", apiEndpoint, reqErr)
	}
	req.Header.Set("User-Agent", userAgent)
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return nil, fmt.Errorf("Unable to execute GET request %q: %v", apiEndpoint, getErr)
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Unable to read body from request %q: %v", apiEndpoint, readErr)
	}

	result.body = &body
	return result, nil
}
