package score

import (
	"net/url"
	"testing"
)

func TestLinkScore(t *testing.T) {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")
	res := GetLinkScores(scoreURL, DefaultInitialTotalSharesCount)
	if res.TotalSharesCount() == DefaultInitialTotalSharesCount {
		t.Errorf("Unable to score URL %q: %v.", scoreURL, res.TotalSharesCount())
	} else {
		t.Logf("Retrieved %q score: %+v", scoreURL, res)
	}
}
