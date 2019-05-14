package score

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
)

// SimulateLinkedInAPI is passed into GetLinkedInShareCountForURL* if we want to simulate the API
const SimulateLinkedInAPI = true

// UseLinkedInAPI is passed into GetLinkedInShareCountForURL* if we don't want to simulate the API, but actually run it
const UseLinkedInAPI = false

// LinkedInLinkScores is the type-safe version of what LinkedIn's share count API returns
type LinkedInLinkScores struct {
	MachineName string  `json:"scorer"`
	HumanName   string  `json:"scorerName"`
	Simulated   bool    `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL         string  `json:"url"`                   // part of lectio.score
	APIEndpoint string  `json:"apiEndPoint"`           // part of lectio.score
	IssuesFound []Issue `json:"issues"`                // part of lectio.score
	Count       int     `json:"count"`                 // direct mapping to LinkedIn API result via Unmarshal httpRes.Body
}

// SourceID returns the name of the scoring engine
func (li LinkedInLinkScores) SourceID() string {
	return li.MachineName
}

// TargetURL is the URL that the scores were computed for
func (li LinkedInLinkScores) TargetURL() string {
	return li.URL
}

// IsValid returns true if the LinkedInLinkScores object is valid (did not return LinkedIn error object)
func (li LinkedInLinkScores) IsValid() bool {
	if li.IssuesFound == nil || len(li.IssuesFound) == 0 {
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

// Issues contains all the problems detected in scoring
func (li LinkedInLinkScores) Issues() Issues {
	return li
}

// ErrorsAndWarnings contains the problems in this link plus satisfies the Link.Issues interface
func (li LinkedInLinkScores) ErrorsAndWarnings() []Issue {
	return li.IssuesFound
}

// IssueCounts returns the total, errors, and warnings counts
func (li LinkedInLinkScores) IssueCounts() (uint, uint, uint) {
	if li.IssuesFound == nil {
		return 0, 0, 0
	}
	var errors, warnings uint
	for _, i := range li.IssuesFound {
		if i.IsError() {
			errors++
		} else {
			warnings++
		}
	}
	return uint(len(li.IssuesFound)), errors, warnings
}

// HandleIssues loops through each issue and calls a particular handler
func (li LinkedInLinkScores) HandleIssues(errorHandler func(Issue), warningHandler func(Issue)) {
	if li.IssuesFound == nil {
		return
	}
	for _, i := range li.IssuesFound {
		if i.IsError() && errorHandler != nil {
			errorHandler(i)
		}
		if i.IsWarning() && warningHandler != nil {
			warningHandler(i)
		}
	}
}

// GetLinkedInLinkScoresForURLText takes a text URL to score and returns the LinkedIn share count
func GetLinkedInLinkScoresForURLText(url string, client *http.Client, simulateLinkedInAPI bool) *LinkedInLinkScores {
	apiEndpoint := "https://www.linkedin.com/countserv/count/share?format=json&url=" + url
	result := new(LinkedInLinkScores)
	result.MachineName = "linkedin"
	result.HumanName = "LinkedIn"
	result.URL = url
	result.APIEndpoint = apiEndpoint
	if simulateLinkedInAPI {
		result.Simulated = true
		result.Count = rand.Intn(50)
		return result
	}
	httpRes, issue := getHTTPResult(apiEndpoint, client, HTTPUserAgent)
	result.APIEndpoint = httpRes.apiEndpoint
	if issue != nil {
		result.IssuesFound = append(result.IssuesFound, issue)
		return result
	}
	json.Unmarshal(*httpRes.body, result)
	return result
}

// GetLinkedInLinkScoresForURL takes a URL to score and returns the LinkedIn share count
func GetLinkedInLinkScoresForURL(url *url.URL, client *http.Client, simulateLinkedInAPI bool) (*LinkedInLinkScores, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetLinkedInLinkScoresForURL")
	}
	return GetLinkedInLinkScoresForURLText(url.String(), client, simulateLinkedInAPI), nil
}
