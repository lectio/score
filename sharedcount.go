package score

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

// SharedCountAPIKeyEnvVarName is the environment variable which may be expected to contain the API key
const SharedCountAPIKeyEnvVarName = "LECTIO_SCORE_SHAREDCOUNT_API_KEY"

// SimulateSharedCountAPI is passed into GetSharedCountLinkScoresForURLText* if we want to simulate the API
const SimulateSharedCountAPI = true

// UseSharedCountAPI is passed into GetSharedCountLinkScoresForURLText* if we don't want to simulate the API, but actually run it
const UseSharedCountAPI = false

type SharedCountCredentials interface {
	SharedCountAPIKey() (string, bool, Issue)
}

// SharedCountLinkScores is the type-safe version of what SharedCount.com's API returns
type SharedCountLinkScores struct {
	MachineName         string                    `json:"scorer"`
	HumanName           string                    `json:"scorerName"`
	Simulated           bool                      `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL                 string                    `json:"url"`                   // part of lectio.score
	GloballyUniqueKey   string                    `json:"uniqueKey"`             // part of lectio.score
	APIEndpoint         string                    `json:"apiEndPoint"`           // part of lectio.score
	IssuesFound         []Issue                   `json:"issues"`                // part of lectio.score
	AggregatedScore     int                       `json:"aggregated_score"`      // part of lectio.score
	ErrorFromAPICall    string                    `json:"Error,omitempty"`       // direct mapping to SharedCount API result via Unmarshal httpRes.Body if there's an error
	ErrorType           string                    `json:"Type,omitempty"`        // direct mapping to SharedCount API result via Unmarshal httpRes.Body if there's an error
	ErrorHTTPStatusCode int                       `json:"HTTP_Code,omitempty"`   // direct mapping to SharedCount API result via Unmarshal httpRes.Body if there's an error
	StumbleUpon         int                       `json:"StumbleUpon"`           // direct mapping to SharedCount API result via Unmarshal httpRes.Body
	Pinterest           int                       `json:"Pinterest"`             // direct mapping to SharedCount API result via Unmarshal httpRes.Body
	LinkedIn            int                       `json:"LinkedIn"`              // direct mapping to SharedCount API result via Unmarshal httpRes.Body
	Facebook            SharedCountFacebookScores `json:"Facebook"`              // direct mapping to SharedCount API result via Unmarshal httpRes.Body
	GooglePlusOne       int                       `json:"GooglePlusOne"`         // direct mapping to SharedCount API result via Unmarshal httpRes.Body
}

// SharedCountFacebookScores returns the group of values returned by Facebook API
type SharedCountFacebookScores struct {
	TotalCount         int `json:"total_count"`
	CommentCount       int `json:"comment_count"`
	ReactionCount      int `json:"reaction_count"`
	ShareCount         int `json:"share_count"`
	CommentPluginCount int `json:"comment_plugin_count"`
}

// SourceID returns the name of the scoring engine
func (sc SharedCountLinkScores) SourceID() string {
	return sc.MachineName
}

// TargetURL is the URL that the scores were computed for
func (sc SharedCountLinkScores) TargetURL() string {
	return sc.URL
}

// IsValid returns true if the SharedCountLinkScores object is valid (did not return SharedCount error object)
func (sc SharedCountLinkScores) IsValid() bool {
	if sc.IssuesFound == nil || len(sc.IssuesFound) == 0 {
		return true
	}
	return false
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (sc SharedCountLinkScores) SharesCount() int {
	if sc.IsValid() {
		return sc.AggregatedScore
	}
	return -1
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (sc SharedCountLinkScores) CommentsCount() int {
	if sc.IsValid() {
		return sc.Facebook.CommentCount
	}
	return -1
}

// Issues contains all the problems detected in scoring
func (sc SharedCountLinkScores) Issues() Issues {
	return sc
}

// ErrorsAndWarnings contains the problems in this link plus satisfies the Link.Issues interface
func (sc SharedCountLinkScores) ErrorsAndWarnings() []Issue {
	return sc.IssuesFound
}

// IssueCounts returns the total, errors, and warnings counts
func (sc SharedCountLinkScores) IssueCounts() (uint, uint, uint) {
	if sc.IssuesFound == nil {
		return 0, 0, 0
	}
	var errors, warnings uint
	for _, i := range sc.IssuesFound {
		if i.IsError() {
			errors++
		} else {
			warnings++
		}
	}
	return uint(len(sc.IssuesFound)), errors, warnings
}

// HandleIssues loops through each issue and calls a particular handler
func (sc SharedCountLinkScores) HandleIssues(errorHandler func(Issue), warningHandler func(Issue)) {
	if sc.IssuesFound == nil {
		return
	}
	for _, i := range sc.IssuesFound {
		if i.IsError() && errorHandler != nil {
			errorHandler(i)
		}
		if i.IsWarning() && warningHandler != nil {
			warningHandler(i)
		}
	}
}

// GetSharedCountLinkScoresForURLText takes a text URL to score and returns the SharedCount share count
func GetSharedCountLinkScoresForURLText(creds SharedCountCredentials, url string, client *http.Client, keys Keys, simulateSharedCountAPI bool) *SharedCountLinkScores {
	result := new(SharedCountLinkScores)
	result.MachineName = "SharedCount.com"
	result.HumanName = "SharedCount.com"
	result.URL = url
	result.GloballyUniqueKey = keys.PrimaryKeyForURLText(url)
	if simulateSharedCountAPI {
		result.Simulated = true
		result.AggregatedScore = rand.Intn(50)
		return result
	}

	apiKey, apiKeyOK, issue := creds.SharedCountAPIKey()
	if !apiKeyOK && issue != nil {
		result.IssuesFound = append(result.IssuesFound, issue)
		return result
	}

	result.APIEndpoint = fmt.Sprintf("https://api.sharedcount.com/v1.0/?url=%s&apikey=%s", url, apiKey)
	httpRes, issue := getHTTPResult(result.APIEndpoint, client, HTTPUserAgent)
	if issue != nil {
		result.IssuesFound = append(result.IssuesFound, issue)
		return result
	}
	result.APIEndpoint = httpRes.apiEndpoint
	json.Unmarshal(*httpRes.body, result)

	if len(result.ErrorFromAPICall) > 0 {
		issue := NewIssue(url, APIErrorResponseFound, fmt.Sprintf("SharedCount API returned an error: %q, %q, %d", result.ErrorFromAPICall, result.ErrorType, result.ErrorHTTPStatusCode), true)
		result.IssuesFound = append(result.IssuesFound, issue)
		return result
	}

	result.AggregatedScore += result.Facebook.TotalCount
	result.AggregatedScore += result.LinkedIn
	result.AggregatedScore += result.StumbleUpon
	result.AggregatedScore += result.Pinterest
	result.AggregatedScore += result.GooglePlusOne

	return result
}

// GetSharedCountLinkScoresForURL takes a URL to score and returns the SharedCount share count
func GetSharedCountLinkScoresForURL(creds SharedCountCredentials, url *url.URL, client *http.Client, keys Keys, simulateSharedCountAPI bool) (*SharedCountLinkScores, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetSharedCountLinkScoresForURL")
	}
	return GetSharedCountLinkScoresForURLText(creds, url.String(), client, keys, simulateSharedCountAPI), nil
}
