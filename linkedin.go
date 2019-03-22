package score

import (
	"encoding/json"
)

const simulateLinkedInAPI = false

// LinkedInCountServResult is the type-safe version of what LinkedIn's share count API returns
type LinkedInCountServResult struct {
	APIEndpoint string `json:"apiEndPoint"`
	HTTPError   error  `json:"httpError"`
	Count       int    `json:"count"`
}

// GetLinkedInShareCountForURL takes a URL to score and returns the LinkedIn share count
func GetLinkedInShareCountForURL(url string) (*LinkedInCountServResult, error) {
	result := new(LinkedInCountServResult)
	if simulateLinkedInAPI {
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
