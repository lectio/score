package cache

import (
	"net/url"
	"testing"

	"github.com/lectio/score"
	"github.com/stretchr/testify/suite"
)

type CacheSuite struct {
	suite.Suite
	cache Cache
	keys  score.Keys
}

func (suite *CacheSuite) SetupSuite() {
	keys := score.MakeDefaultKeys()
	cache, err := MakeFileCache("valid-scores", "invalid-scores", true, keys, -1, false)
	if err != nil {
		panic(err)
	}
	suite.cache = cache
}

func (suite *CacheSuite) TearDownSuite() {
	suite.cache.Close()
}

func (suite *CacheSuite) TestSingle() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")
	scores, err := suite.cache.Get(scoreURL)
	suite.Nil(err, "There should be no cache error")
	suite.NotNil(scores, "Score should not be nil")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}
