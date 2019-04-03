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

// LinkedInLinkScores is the type-safe version of what LinkedIn's share count API returns
type LinkedInLinkScores struct {
	Simulated         bool   `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL               string `json:"url"`                   // part of lectio.score
	GloballyUniqueKey string `json:"uniqueKey"`             // part of lectio.score
	APIEndpoint       string `json:"apiEndPoint"`           // part of lectio.score
	HTTPError         error  `json:"httpError,omitempty"`   // part of lectio.score
	Count             int    `json:"count"`                 // direct mapping to LinkedIn API result via Unmarshal httpRes.Body
}

// Names returns the identities of the scorer
func (li LinkedInLinkScores) Names() (string, string) {
	return "linkedin", "LinkedIn"
}

// TargetURL is the URL that the scores were computed for
func (li LinkedInLinkScores) TargetURL() string {
	return li.URL
}

// IsValid returns true if the LinkedInLinkScores object is valid (did not return LinkedIn error object)
func (li LinkedInLinkScores) IsValid() bool {
	if li.HTTPError == nil {
		return true
	}
	return false
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (li LinkedInLinkScores) SharesCount() int {
	if li.IsValid() {
		return li.Count
	}
	return -1
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (li LinkedInLinkScores) CommentsCount() int {
	return -1
}

// GetLinkedInLinkScoresForURLText takes a text URL to score and returns the LinkedIn share count
func GetLinkedInLinkScoresForURLText(url string, globallyUniqueKey string, simulateLinkedInAPI bool) (*LinkedInLinkScores, error) {
	apiEndpoint := "https://www.linkedin.com/countserv/count/share?format=json&url=" + url
	result := new(LinkedInLinkScores)
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

// GetLinkedInLinkScoresForURL takes a URL to score and returns the LinkedIn share count
func GetLinkedInLinkScoresForURL(url *url.URL, globallyUniqueKey string, simulateLinkedInAPI bool) (*LinkedInLinkScores, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetLinkedInLinkScoresForURL")
	}
	return GetLinkedInLinkScoresForURLText(url.String(), globallyUniqueKey, simulateLinkedInAPI)
}
