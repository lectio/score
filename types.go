package score

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

// LinkScorerIdentity uniquely identifies the link scoring engine (the "scorer")
type LinkScorerIdentity interface {
	MachineName() string                            // usually lowercase identifer useful for machine processing
	HumanName() string                              // can be any meaningful human identifer
	FileName(path string, scores LinkScores) string // create the name of this file for file storage
}

// LinkScores instances score a given link (by running an API or other computation)
type LinkScores interface {
	Identity() LinkScorerIdentity
	TargetURL() string
	TargetURLUniqueKey() string
	IsValid() bool
	SharesCount() int
	CommentsCount() int
}

// HTTPUserAgent may be passed into getHTTPResult as the default HTTP User-Agent header parameter
const HTTPUserAgent = "github.com/lectio/score"

// HTTPTimeout may be passed into getHTTPResult function as the default HTTP timeout parameter
const HTTPTimeout = time.Second * 90

type httpResult struct {
	apiEndpoint string
	body        *[]byte
}

// GetHTTPResult runs the apiEndpoint and returns the body of the HTTP result
func getHTTPResult(apiEndpoint string, userAgent string, timeout time.Duration) (*httpResult, error) {
	result := new(httpResult)
	result.apiEndpoint = apiEndpoint

	httpClient := http.Client{
		Timeout: timeout,
	}
	req, reqErr := http.NewRequest(http.MethodGet, apiEndpoint, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("Unable to create request %q: %v", apiEndpoint, reqErr)
	}
	req.Header.Set("User-Agent", userAgent)
	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return nil, fmt.Errorf("Unable to execute GET request %q: %v", apiEndpoint, getErr)
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Unable to read body from request %q: %v", apiEndpoint, readErr)
	}

	result.body = &body
	return result, nil
}

type defaultLinkScorerIdentity struct {
	machineName string
	humanName   string
}

func makeDefaultLinkScorerIdentity(machineName string, humanName string) *defaultLinkScorerIdentity {
	result := new(defaultLinkScorerIdentity)
	result.machineName = machineName
	result.humanName = humanName
	return result
}

// MachineName is usually lowercase identifer useful for machine processing
func (i defaultLinkScorerIdentity) MachineName() string {
	return i.machineName
}

// HumanName can be any meaningful human identifer
func (i defaultLinkScorerIdentity) HumanName() string {
	return i.humanName
}

// FileName creates the name of this file for file storage
func (i defaultLinkScorerIdentity) FileName(path string, scores LinkScores) string {
	suffix := i.MachineName()
	if !scores.IsValid() {
		suffix = suffix + "-error"
	}
	return fmt.Sprintf("%s.%s.json", filepath.Join(path, scores.TargetURLUniqueKey()), suffix)
}
