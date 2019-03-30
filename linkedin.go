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
func GetLinkedInShareCountForURLText(url string, simulateLinkedInAPI bool) (*LinkedInCountServResult, error) {
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
func GetLinkedInShareCountForURL(url *url.URL, simulateLinkedInAPI bool) (*LinkedInCountServResult, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetFacebookGraphForURL")
	}
	return GetLinkedInShareCountForURLText(url.String(), simulateLinkedInAPI)
}
