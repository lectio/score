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
	TotalSharesCount() int
	FacebookGraph() *FacebookGraphResult
	LinkedInCount() *LinkedInCountServResult
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
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("Unable to read body from request %q: %v", apiEndpoint, readErr)
	}

	result.body = &body
	return result, nil
}

type defaultCuratedLinkScores struct {
	url              *url.URL
	totalSharesCount int
	facebookGraph    *FacebookGraphResult
	linkedInGraph    *LinkedInCountServResult
}

func (d *defaultCuratedLinkScores) init(url *url.URL, initialTotalCount int) {
	d.url = url
	d.totalSharesCount = initialTotalCount

	d.facebookGraph, _ = GetFacebookGraphForURL(url)
	d.linkedInGraph, _ = GetLinkedInShareCountForURL(url)

	if d.facebookGraph.IsValid() && d.facebookGraph != nil && d.facebookGraph.Shares != nil && d.facebookGraph.Shares.ShareCount > 0 {
		d.totalSharesCount = d.facebookGraph.Shares.ShareCount
	}
	if d.linkedInGraph.IsValid() && d.linkedInGraph != nil && d.linkedInGraph.Count > 0 {
		if d.totalSharesCount == -1 {
			d.totalSharesCount = d.linkedInGraph.Count
		} else {
			d.totalSharesCount += d.linkedInGraph.Count
		}
	}
}

func (d defaultCuratedLinkScores) TargetURL() *url.URL {
	return d.url
}

func (d defaultCuratedLinkScores) TotalSharesCount() int {
	return d.totalSharesCount
}

func (d defaultCuratedLinkScores) FacebookGraph() *FacebookGraphResult {
	return d.facebookGraph
}

func (d defaultCuratedLinkScores) LinkedInCount() *LinkedInCountServResult {
	return d.linkedInGraph
}

// GetLinkScores computes and return social scores for a curated link
func GetLinkScores(url *url.URL, initialTotalCount int) LinkScores {
	result := new(defaultCuratedLinkScores)
	result.init(url, initialTotalCount)
	return result
}
