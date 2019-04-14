package container

import (
	"net/url"
	"testing"

	"github.com/lectio/score/cache"

	"github.com/lectio/score"
	"github.com/stretchr/testify/suite"
)

type ContainerSuite struct {
	suite.Suite
	keys  score.Keys
	cache cache.Cache
}

func (suite *ContainerSuite) SetupSuite() {
	suite.keys = score.MakeDefaultKeys()
	suite.cache = cache.MakeNullCache(suite.keys, -1, false)
}

func (suite *ContainerSuite) TearDownSuite() {
}

func (suite *ContainerSuite) TestCollection() {
	scoreURL, _ := url.Parse("https://www.cnbc.com/2019/03/18/bill-gates-says-he-talked-with-google-employees-about-ai-health-care.html")
	urls := [...]*url.URL{scoreURL}
	handler := func(index int) (*url.URL, error) {
		return urls[index], nil
	}
	iterator := func() (startIndex int, endIndex int, keys score.Keys, retrievalFn TargetsIteratorRetrievalFn) {
		return 0, len(urls) - 1, suite.keys, handler
	}
	sc := MakeCollection(suite.cache, iterator, nil)
	suite.NotNil(sc, "Scores collection should not be Nil")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ContainerSuite))
}
