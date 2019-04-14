package cache

import (
	"net/url"
	"time"

	"github.com/lectio/score"
)

type nullCache struct {
	keys              score.Keys
	initialTotalCount int
	simulate          bool
}

// MakeNullCache creates an instance of a "noop" cache, which always runs Score (no storage)
func MakeNullCache(keys score.Keys, initialTotalCount int, simulate bool) Cache {
	cache := new(nullCache)
	cache.keys = keys
	cache.initialTotalCount = initialTotalCount
	cache.simulate = simulate
	return cache
}

func (c nullCache) Score(url *url.URL) (score.LinkScores, error) {
	als := score.GetAggregatedLinkScores(url, c.keys, c.initialTotalCount, c.simulate)
	return als, nil
}

func (c nullCache) Get(url *url.URL) (score.LinkScores, error) {
	return c.Score(url)
}

func (c nullCache) Find(url *url.URL) (scores score.LinkScores, found bool, expired bool, err error) {
	return nil, false, true, nil
}

func (c nullCache) Save(scores score.LinkScores, autoExpire time.Duration) error {
	return nil
}

func (c nullCache) Close() error {
	return nil
}
