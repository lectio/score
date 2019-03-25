package score

import (
	"encoding/json"
	"net/url"
)

var simulateFacebookAPI = false

// GetSimulateFacebookAPI returns the value of simulateFacebookAPI flag
func GetSimulateFacebookAPI() bool {
	return simulateFacebookAPI
}

// SetSimulateFacebookAPI sets the simulateFacebookAPI flag to indicate whether to run the HTTP API or
// just simulate an execution
func SetSimulateFacebookAPI(flag bool) {
	simulateFacebookAPI = flag
}

// FacebookGraphResult is the type-safe version of what Facebook API Graph returns
type FacebookGraphResult struct {
	Simulated   bool                   `json:"isSimulated"`
	APIEndpoint string                 `json:"apiEndPoint"`
	HTTPError   error                  `json:"httpError"`
	APIError    *FacebookGraphAPIError `json:"error"`
	ID          string                 `json:"id"`
	Shares      *FacebookGraphShares   `json:"share"`
	OpenGraph   *FacebookGraphOGObject `json:"og_object"`
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

// IsValid returns true if the FacebookGraphResult object is valid (did not return Facebook error object)
func (fbgr FacebookGraphResult) IsValid() bool {
	if fbgr.HTTPError == nil && fbgr.APIError == nil {
		return true
	}
	return false
}

// GetFacebookGraphForURLText takes a text URL to score and returns the Facebook graph (and share counts)
func GetFacebookGraphForURLText(url string) (*FacebookGraphResult, error) {
	result := new(FacebookGraphResult)
	if simulateFacebookAPI {
		result.Simulated = true
		return result, nil
	}
	httpRes, httpErr := getHTTPResult("https://graph.facebook.com/?id="+url, HTTPUserAgent, HTTPTimeout)
	result.APIEndpoint = httpRes.apiEndpoint
	result.HTTPError = httpErr
	if httpErr != nil {
		return result, httpErr
	}
	json.Unmarshal(*httpRes.body, result)
	return result, nil
}

// GetFacebookGraphForURL takes a URL to score and returns the Facebook graph (and share counts)
func GetFacebookGraphForURL(url *url.URL) (*FacebookGraphResult, error) {
	return GetFacebookGraphForURLText(url.String())
}
