package score

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// LinkScores are the social or other score types associated with curated content
type LinkScores interface {
	TargetURL() *url.URL
	TargetGloballyUniqueKey() string
	TotalSharesCount() int
	FacebookGraph() *FacebookLinkScoreGraphResult
	LinkedInCount() *LinkedInLinkScoreResult
}

// DefaultInitialTotalSharesCount is the default value for shares count.
// This allows to distinguish whether something was computed or if there was an error.
const DefaultInitialTotalSharesCount = -1

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

type defaultCuratedLinkScores struct {
	url               *url.URL
	globallyUniqueKey string
	totalSharesCount  int
	facebookLinkScore *FacebookLinkScoreGraphResult
	linkedInLinkScore *LinkedInLinkScoreResult
}

func (d *defaultCuratedLinkScores) init(url *url.URL, globallyUniqueKey string, initialTotalCount int, simulateAPIs bool) {
	d.url = url
	d.globallyUniqueKey = globallyUniqueKey

	d.totalSharesCount = initialTotalCount

	var fbErr, liErr error
	d.facebookLinkScore, fbErr = GetFacebookGraphForURL(url, globallyUniqueKey, simulateAPIs)
	d.linkedInLinkScore, liErr = GetLinkedInShareCountForURL(url, globallyUniqueKey, simulateAPIs)

	if fbErr == nil && d.facebookLinkScore.IsValid() && d.facebookLinkScore != nil && d.facebookLinkScore.Shares != nil && d.facebookLinkScore.Shares.ShareCount > 0 {
		d.totalSharesCount = d.facebookLinkScore.Shares.ShareCount
	}
	if liErr == nil && d.linkedInLinkScore.IsValid() && d.linkedInLinkScore != nil && d.linkedInLinkScore.Count > 0 {
		if d.totalSharesCount == -1 {
			d.totalSharesCount = d.linkedInLinkScore.Count
		} else {
			d.totalSharesCount += d.linkedInLinkScore.Count
		}
	}
}

func (d defaultCuratedLinkScores) TargetURL() *url.URL {
	return d.url
}

func (d defaultCuratedLinkScores) TargetGloballyUniqueKey() string {
	return d.globallyUniqueKey
}

func (d defaultCuratedLinkScores) TotalSharesCount() int {
	return d.totalSharesCount
}

func (d defaultCuratedLinkScores) FacebookGraph() *FacebookLinkScoreGraphResult {
	return d.facebookLinkScore
}

func (d defaultCuratedLinkScores) LinkedInCount() *LinkedInLinkScoreResult {
	return d.linkedInLinkScore
}

// GetLinkScores computes and return social scores for a curated link
func GetLinkScores(url *url.URL, globallyUniqueKey string, initialTotalCount int, simulateAPIs bool) LinkScores {
	result := new(defaultCuratedLinkScores)
	result.init(url, globallyUniqueKey, initialTotalCount, simulateAPIs)
	return result
}
