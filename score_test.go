package score

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ScoreSuite struct {
	suite.Suite
	keys Keys
}

func (suite *ScoreSuite) SetupSuite() {
	suite.keys = MakeDefaultKeys()
}

func (suite *ScoreSuite) TearDownSuite() {
}

func (suite *ScoreSuite) TestScores() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")

	fb, fbErr := GetFacebookLinkScoresForURL(scoreURL, suite.keys, UseFacebookAPI)
	suite.Nil(fbErr, "There shouldn't be a Facebook API error")
	suite.False(fb.SharesCount() == -1, "Facebook shares count shouldn't be the default")

	li, liErr := GetLinkedInLinkScoresForURL(scoreURL, suite.keys, UseLinkedInAPI)
	suite.Nil(liErr, "There shouldn't be a LinkedIn API error")
	suite.False(li.SharesCount() == -1, "LinkedIn shares count shouldn't be the default")

	aggregated := GetAggregatedLinkScores(scoreURL, suite.keys, -1, false)
	suite.False(aggregated.SharesCount() == -1, "Aggregate count shouldn't be the default")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ScoreSuite))
}
