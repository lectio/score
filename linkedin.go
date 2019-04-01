package score

import (
	"encoding/json"
	"errors"
	"net/url"
)

// SimulateLinkedInAPI is passed into GetLinkedInShareCountForURL* if we want to simulate the API
const SimulateLinkedInAPI = true

// UseLinkedInAPI is passed into GetLinkedInShareCountForURL* if we don't want to simulate the API, but actually run it
const UseLinkedInAPI = false

// LinkedInLinkScoreResult is the type-safe version of what LinkedIn's share count API returns
type LinkedInLinkScoreResult struct {
	Simulated         bool   `json:"isSimulated"`
	URL               string `json:"url"`
	GloballyUniqueKey string `json:"uniqueKey"`
	APIEndpoint       string `json:"apiEndPoint"`
	HTTPError         error  `json:"httpError"`
	Count             int    `json:"count"`
}

// IsValid returns true if the LinkedInLinkScoreResult object is valid (did not return LinkedIn error object)
func (licsr LinkedInLinkScoreResult) IsValid() bool {
	if licsr.HTTPError == nil {
		return true
	}
	return false
}

// GetLinkedInShareCountForURLText takes a text URL to score and returns the LinkedIn share count
func GetLinkedInShareCountForURLText(url string, globallyUniqueKey string, simulateLinkedInAPI bool) (*LinkedInLinkScoreResult, error) {
	apiEndpoint := "https://www.linkedin.com/countserv/count/share?format=json&url=" + url
	result := new(LinkedInLinkScoreResult)
	result.URL = url
	result.APIEndpoint = apiEndpoint
	result.GloballyUniqueKey = globallyUniqueKey
	if simulateLinkedInAPI {
		result.Simulated = true
		return result, nil
	}
	httpRes, httpErr := getHTTPResult(apiEndpoint, HTTPUserAgent, HTTPTimeout)
	result.APIEndpoint = httpRes.apiEndpoint
	result.HTTPError = httpErr
	if httpErr != nil {
		return result, httpErr
	}
	json.Unmarshal(*httpRes.body, result)
	return result, nil
}

// GetLinkedInShareCountForURL takes a URL to score and returns the LinkedIn share count
func GetLinkedInShareCountForURL(url *url.URL, globallyUniqueKey string, simulateLinkedInAPI bool) (*LinkedInLinkScoreResult, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetFacebookGraphForURL")
	}
	return GetLinkedInShareCountForURLText(url.String(), globallyUniqueKey, simulateLinkedInAPI)
}
