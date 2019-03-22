package score

import (
	"testing"
)

const scoreURL = "https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html"

func TestFacebookScore(t *testing.T) {
	res, err := GetFacebookGraphForURL(scoreURL)
	if err != nil {
		t.Errorf("Unable to score URL %q: %v.", scoreURL, err)
	} else {
		t.Logf("Retrieved %q score: %+v", scoreURL, res)
	}
}
