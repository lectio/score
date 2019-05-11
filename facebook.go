package score

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
)

// TODO: If Facebook rate limiting gets in the way, try https://github.com/ustayready/fireprox

// SimulateFacebookAPI is passed into GetFacebookGraphForURL* if we want to simulate the API
const SimulateFacebookAPI = true

// UseFacebookAPI is passed into GetFacebookGraphForURL* if we don't want to simulate the API, but actually run it
const UseFacebookAPI = false

// FacebookLinkScores is the type-safe version of what Facebook API Graph returns
type FacebookLinkScores struct {
	MachineName       string                 `json:"scorer"`
	HumanName         string                 `json:"scorerName"`
	Simulated         bool                   `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL               string                 `json:"url"`                   // part of lectio.score
	GloballyUniqueKey string                 `json:"uniqueKey"`             // part of lectio.score
	APIEndpoint       string                 `json:"apiEndPoint"`           // part of lectio.score
	IssuesFound       []Issue                `json:"issues"`                // part of lectio.score
	APIError          *FacebookGraphAPIError `json:"error,omitempty"`       // direct mapping to Facebook API result via Unmarshal httpRes.Body
	ID                string                 `json:"id"`                    // direct mapping to Facebook API result via Unmarshal httpRes.Body
	Shares            *FacebookGraphShares   `json:"share"`                 // direct mapping to Facebook API result via Unmarshal httpRes.Body
	OpenGraph         *FacebookGraphOGObject `json:"og_object"`             // direct mapping to Facebook API result via Unmarshal httpRes.Body
}

// SourceID returns the name of the scoring engine
func (fb FacebookLinkScores) SourceID() string {
	return fb.MachineName
}

// TargetURL is the URL that the scores were computed for
func (fb FacebookLinkScores) TargetURL() string {
	return fb.URL
}

// IsValid returns true if the FacebookLinkScores object is valid (did not return Facebook error object)
func (fb FacebookLinkScores) IsValid() bool {
	if fb.IssuesFound == nil || len(fb.IssuesFound) == 0 {
		return true
	}
	return false
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (fb FacebookLinkScores) SharesCount() int {
	if fb.IsValid() && fb.Shares != nil {
		return fb.Shares.ShareCount
	}
	return -1
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (fb FacebookLinkScores) CommentsCount() int {
	if fb.IsValid() && fb.Shares != nil {
		return fb.Shares.CommentCount
	}
	return -1
}

// Issues contains all the problems detected in scoring
func (fb FacebookLinkScores) Issues() Issues {
	return fb
}

// ErrorsAndWarnings contains the problems in this link plus satisfies the Link.Issues interface
func (fb FacebookLinkScores) ErrorsAndWarnings() []Issue {
	return fb.IssuesFound
}

// IssueCounts returns the total, errors, and warnings counts
func (fb FacebookLinkScores) IssueCounts() (uint, uint, uint) {
	if fb.IssuesFound == nil {
		return 0, 0, 0
	}
	var errors, warnings uint
	for _, i := range fb.IssuesFound {
		if i.IsError() {
			errors++
		} else {
			warnings++
		}
	}
	return uint(len(fb.IssuesFound)), errors, warnings
}

// HandleIssues loops through each issue and calls a particular handler
func (fb FacebookLinkScores) HandleIssues(errorHandler func(Issue), warningHandler func(Issue)) {
	if fb.IssuesFound == nil {
		return
	}
	for _, i := range fb.IssuesFound {
		if i.IsError() && errorHandler != nil {
			errorHandler(i)
		}
		if i.IsWarning() && warningHandler != nil {
			warningHandler(i)
		}
	}
}

// FacebookGraphAPIError is the type-safe version of a Facebook API Graph error (e.g. rate limiting)
type FacebookGraphAPIError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Transient bool   `json:"is_transient"`
	Code      string `json:"code"`
	TraceID   string `json:"fbtrace_id"`
}

// FacebookGraphShares is the type-safe version of a Facebook API Graph shares object
type FacebookGraphShares struct {
	ShareCount   int `json:"share_count"`
	CommentCount int `json:"comment_count"`
}

// FacebookGraphOGObject is the type-safe version of a Facebook API OpenGraph object
type FacebookGraphOGObject struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// GetFacebookLinkScoresForURLText takes a text URL to score and returns the Facebook graph (and share counts)
func GetFacebookLinkScoresForURLText(url string, client *http.Client, keys Keys, simulateFacebookAPI bool) *FacebookLinkScores {
	apiEndpoint := "https://graph.facebook.com/?id=" + url
	result := new(FacebookLinkScores)
	result.MachineName = "facebook"
	result.HumanName = "Facebook"
	result.URL = url
	result.APIEndpoint = apiEndpoint
	result.GloballyUniqueKey = keys.PrimaryKeyForURLText(url)
	if simulateFacebookAPI {
		result.Simulated = simulateFacebookAPI
		result.Shares = new(FacebookGraphShares)
		result.Shares.ShareCount = rand.Intn(750)
		result.Shares.CommentCount = rand.Intn(2500)
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

// GetFacebookLinkScoresForURL takes a URL to score and returns the Facebook graph (and share counts)
func GetFacebookLinkScoresForURL(url *url.URL, client *http.Client, keys Keys, simulateFacebookAPI bool) (*FacebookLinkScores, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetFacebookLinkScoresForURL")
	}
	return GetFacebookLinkScoresForURLText(url.String(), client, keys, simulateFacebookAPI), nil
}
