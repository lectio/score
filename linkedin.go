package score

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
)

// SimulateLinkedInAPI is passed into GetLinkedInShareCountForURL* if we want to simulate the API
const SimulateLinkedInAPI = true

// UseLinkedInAPI is passed into GetLinkedInShareCountForURL* if we don't want to simulate the API, but actually run it
const UseLinkedInAPI = false

// LinkedInLinkScoreResult is the type-safe version of what LinkedIn's share count API returns
type LinkedInLinkScoreResult struct {
	Simulated         bool   `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL               string `json:"url"`                   // part of lectio.score
	GloballyUniqueKey string `json:"uniqueKey"`             // part of lectio.score
	APIEndpoint       string `json:"apiEndPoint"`           // part of lectio.score
	HTTPError         error  `json:"httpError,omitempty"`   // part of lectio.score
	Count             int    `json:"count"`                 // direct mapping to LinkedIn API result via Unmarshal httpRes.Body
}

// Names returns the identities of the scorer
func (li LinkedInLinkScoreResult) Names() (string, string) {
	return "facebook", "Facebook"
}

// TargetURL is the URL that the scores were computed for
func (li LinkedInLinkScoreResult) TargetURL() string {
	return li.URL
}

// IsValid returns true if the LinkedInLinkScoreResult object is valid (did not return LinkedIn error object)
func (li LinkedInLinkScoreResult) IsValid() bool {
	if li.HTTPError == nil {
		return true
	}
	return false
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (li LinkedInLinkScoreResult) SharesCount() int {
	if li.IsValid() {
		return li.Count
	}
	return -1
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (li LinkedInLinkScoreResult) CommentsCount() int {
	return -1
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
		result.Count = rand.Intn(50)
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
