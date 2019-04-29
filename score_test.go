package score

import (
	"net/url"
	"testing"
	"fmt"

	"github.com/stretchr/testify/suite"
	"github.com/lectio/secret"
)

type ScoreSuite struct {
	suite.Suite
	keys Keys
	vault secret.Vault
}

func (suite *ScoreSuite) SetupSuite() {
	suite.keys = MakeDefaultKeys()

	var vaultErr error
	suite.vault, vaultErr = secret.Parse("env://LECTIO_VAULTPP_DEFAULT")
	if vaultErr != nil {
		panic(vaultErr)
	}
}

func (suite *ScoreSuite) TearDownSuite() {
}

func (suite *ScoreSuite) SharedCountAPIKey() (string, bool, Issue) {
	apiKey, err := suite.vault.DecryptText("0d4af7674abbfa18d01510fc107318ace74175c5cae32b1e3dfb1ec37ee5ceb1c8253d880ba027ed3c8280883cef0152d447f068a21f0a793f83c552fd89703aeecd53d5")
	fmt.Printf("SC API: %s (%+v)\n", apiKey, err)
	if err != nil {
		return "", false, newIssue("SharedCount.com", SecretManagementError, err.Error(), true)
	}
	return apiKey, true, nil
}

func (suite *ScoreSuite) TestScores() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")

	// [Get number of shares from social platforms](https://gist.github.com/ihorvorotnov/9132596)
	// https://rudrastyh.com/facebook/get-share-count-for-url.html#curl

	// https://graph.facebook.com/v2.0/id?https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html
	// https://api.facebook.com/v2.0/method/links.getStats?urls=https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html&format=json
	// http://buttons.reddit.com/button_info.json?url=https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html
	// http://www.linkedin.com/countserv/count/share?url=https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html&format=json

	fb, fbErr := GetFacebookLinkScoresForURL(scoreURL, suite.keys, UseFacebookAPI)
	suite.Nil(fbErr, "There shouldn't be a Facebook API error")
	suite.True(fb.IsValid(), "There shouldn't be a Facebook API error")
	suite.False(fb.SharesCount() == -1, "Facebook shares count shouldn't be the default")

	li, liErr := GetLinkedInLinkScoresForURL(scoreURL, suite.keys, UseLinkedInAPI)
	suite.Nil(liErr, "There shouldn't be a LinkedIn API error")
	suite.True(li.IsValid(), "There shouldn't be a LinkedIn API error")
	suite.False(li.SharesCount() == -1, "LinkedIn shares count shouldn't be the default")

	sharedCount, scErr := GetSharedCountLinkScoresForURL(suite, scoreURL, suite.keys, UseSharedCountAPI)
	suite.Nil(scErr, "There shouldn't be a SharedCount API error")
	suite.True(sharedCount.IsValid(), "SharedCount scores should be valid")
	suite.True(sharedCount.SharesCount() > 0, "SharedCount score should be greater than zero")

	aggregated := GetAggregatedLinkScores(scoreURL, suite.keys, -1, false)
	suite.True(aggregated.IsValid(), "All scores should be valid")
	suite.False(aggregated.SharesCount() == -1, "Aggregate count shouldn't be the default")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ScoreSuite))
}
