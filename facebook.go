package score

import (
	"encoding/json"
	"errors"
	"net/url"
)

// SimulateFacebookAPI is passed into GetFacebookGraphForURL* if we want to simulate the API
const SimulateFacebookAPI = true

// UseFacebookAPI is passed into GetFacebookGraphForURL* if we don't want to simulate the API, but actually run it
const UseFacebookAPI = false

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
func GetFacebookGraphForURLText(url string, simulateFacebookAPI bool) (*FacebookGraphResult, error) {
	result := new(FacebookGraphResult)
	if simulateFacebookAPI {
		result.Simulated = simulateFacebookAPI
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
func GetFacebookGraphForURL(url *url.URL, simulateFacebookAPI bool) (*FacebookGraphResult, error) {
	if url == nil {
		return nil, errors.New("Null URL passed to GetFacebookGraphForURL")
	}
	return GetFacebookGraphForURLText(url.String(), simulateFacebookAPI)
}
