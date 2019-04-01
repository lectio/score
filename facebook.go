package score

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
)

// SimulateFacebookAPI is passed into GetFacebookGraphForURL* if we want to simulate the API
const SimulateFacebookAPI = true

// UseFacebookAPI is passed into GetFacebookGraphForURL* if we don't want to simulate the API, but actually run it
const UseFacebookAPI = false

// FacebookLinkScoreGraphResult is the type-safe version of what Facebook API Graph returns
type FacebookLinkScoreGraphResult struct {
	Simulated         bool                   `json:"isSimulated,omitempty"` // part of lectio.score, omitted if it's false
	URL               string                 `json:"url"`                   // part of lectio.score
	GloballyUniqueKey string                 `json:"uniqueKey"`             // part of lectio.score
	APIEndpoint       string                 `json:"apiEndPoint"`           // part of lectio.score
	HTTPError         error                  `json:"httpError,omitempty"`   // part of lectio.score
	APIError          *FacebookGraphAPIError `json:"error,omitempty"`       // direct mapping to Facebook API result via Unmarshal httpRes.Body
	ID                string                 `json:"id"`                    // direct mapping to Facebook API result via Unmarshal httpRes.Body
	Shares            *FacebookGraphShares   `json:"share"`                 // direct mapping to Facebook API result via Unmarshal httpRes.Body
	OpenGraph         *FacebookGraphOGObject `json:"og_object"`             // direct mapping to Facebook API result via Unmarshal httpRes.Body
}

// Names returns the identities of the scorer
func (fb FacebookLinkScoreGraphResult) Names() (string, string) {
	return "facebook", "Facebook"
}

// TargetURL is the URL that the scores were computed for
func (fb FacebookLinkScoreGraphResult) TargetURL() string {
	return fb.URL
}

// IsValid returns true if the FacebookLinkScoreGraphResult object is valid (did not return Facebook error object)
func (fb FacebookLinkScoreGraphResult) IsValid() bool {
	if fb.HTTPError == nil && fb.APIError == nil {
		return true
	}
	return false
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (fb FacebookLinkScoreGraphResult) SharesCount() int {
	if fb.IsValid() && fb.Shares != nil {
		return fb.Shares.ShareCount
	}
	return -1
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (fb FacebookLinkScoreGraphResult) CommentsCount() int {
	if fb.IsValid() && fb.Shares != nil {
		return fb.Shares.CommentCount
	}
	return -1
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

// GetFacebookGraphForURLText takes a text URL to score and returns the Facebook graph (and share counts)
func GetFacebookGraphForURLText(url string, globallyUniqueKey string, simulateFacebookAPI bool) (*FacebookLinkScoreGraphResult, error) {
	apiEndpoint := "https://graph.facebook.com/?id=" + url
	result := new(FacebookLinkScoreGraphResult)
	result.URL = url
	result.APIEndpoint = apiEndpoint
	result.GloballyUniqueKey = globallyUniqueKey
	if simulateFacebookAPI {
		result.Simulated = simulateFacebookAPI
		result.Shares = new(FacebookGraphShares)
		result.Shares.ShareCount = rand.Intn(750)
		result.Shares.CommentCount = rand.Intn(2500)
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

// GetFacebookGraphForURL takes a URL to score and returns the Facebook graph (and share counts)
func GetFacebookGraphForURL(url *url.URL, globallyUniqueKey string, simulateFacebookAPI bool) (*FacebookLinkScoreGraphResult, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetFacebookGraphForURL")
	}
	return GetFacebookGraphForURLText(url.String(), globallyUniqueKey, simulateFacebookAPI)
}
