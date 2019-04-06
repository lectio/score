package score

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ScoreSuite struct {
	suite.Suite
}

func (suite *ScoreSuite) SetupSuite() {
}

func (suite *ScoreSuite) TearDownSuite() {
}

func (suite *ScoreSuite) TestScores() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")

	fb, fbErr := GetFacebookLinkScoresForURL(scoreURL, scoreURL.String(), UseFacebookAPI)
	suite.Nil(fbErr, "There shouldn't be a Facebook API error")
	suite.False(fb.SharesCount() == -1, "Facebook shares count shouldn't be the default")

	li, liErr := GetLinkedInLinkScoresForURL(scoreURL, scoreURL.String(), UseLinkedInAPI)
	suite.Nil(liErr, "There shouldn't be a LinkedIn API error")
	suite.False(li.SharesCount() == -1, "LinkedIn shares count shouldn't be the default")

	aggregated := GetAggregatedLinkScores(scoreURL, scoreURL.String(), -1, false)
	suite.False(aggregated.SharesCount() == -1, "Aggregate count shouldn't be the default")
}

func (suite *ScoreSuite) TestCollection() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")
	urls := [...]*url.URL{scoreURL}
	handler := func(index int) (*url.URL, string, error) {
		url := urls[index]
		return urls[index], url.String(), nil
	}
	iterator := func() (startIndex int, endIndex int, retrievalFn TargetsIteratorRetrievalFn) {
		return 0, len(urls) - 1, handler
	}
	sc := MakeCollection(iterator, false, true)
	suite.NotNil(sc, "Scores collection should not be Nil")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ScoreSuite))
}
