package score

import (
	"encoding/json"
	"net/url"
)

var simulateLinkedInAPI = false

// GetSimulateLinkedInAPI returns the value of simulateLinkedInAPI flag
func GetSimulateLinkedInAPI() bool {
	return simulateLinkedInAPI
}

// SetSimulateLinkedInAPI sets the simulateLinkedInAPI flag to indicate whether to run the HTTP API or
// just simulate an execution
func SetSimulateLinkedInAPI(flag bool) {
	simulateLinkedInAPI = flag
}

// LinkedInCountServResult is the type-safe version of what LinkedIn's share count API returns
type LinkedInCountServResult struct {
	Simulated   bool   `json:"isSimulated"`
	APIEndpoint string `json:"apiEndPoint"`
	HTTPError   error  `json:"httpError"`
	Count       int    `json:"count"`
}

// IsValid returns true if the LinkedInCountServResult object is valid (did not return LinkedIn error object)
func (licsr LinkedInCountServResult) IsValid() bool {
	if licsr.HTTPError == nil {
		return true
	}
	return false
}

// GetLinkedInShareCountForURLText takes a text URL to score and returns the LinkedIn share count
func GetLinkedInShareCountForURLText(url string) (*LinkedInCountServResult, error) {
	result := new(LinkedInCountServResult)
	if simulateLinkedInAPI {
		result.Simulated = true
		return result, nil
	}
	httpRes, httpErr := getHTTPResult("https://www.linkedin.com/countserv/count/share?format=json&url="+url, HTTPUserAgent, HTTPTimeout)
	result.APIEndpoint = httpRes.apiEndpoint
	result.HTTPError = httpErr
	if httpErr != nil {
		return result, httpErr
	}
	json.Unmarshal(*httpRes.body, result)
	return result, nil
}

// GetLinkedInShareCountForURL takes a URL to score and returns the LinkedIn share count
func GetLinkedInShareCountForURL(url *url.URL) (*LinkedInCountServResult, error) {
	return GetLinkedInShareCountForURLText(url.String())
}
