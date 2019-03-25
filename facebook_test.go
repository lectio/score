package score

import (
	"net/url"
	"testing"
)

func TestFacebookScore(t *testing.T) {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")
	res, err := GetFacebookGraphForURL(scoreURL, UseFacebookAPI)
	if err != nil {
		t.Errorf("Unable to score URL %q: %v.", scoreURL, err)
	} else {
		t.Logf("Retrieved %q score: %+v", scoreURL, res)
	}
}
