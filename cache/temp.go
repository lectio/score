package cache

import (
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/lectio/score"
)

type tempCache struct {
	fileCache        *fileCache
	removeAllOnClose bool
}

// MakeTempCache creates an instance of a cache, which stores links on disk, in a temp directory
func MakeTempCache(keys score.Keys, initialTotalCount int, simulate bool, removeAllOnClose bool) (Cache, error) {
	validScoresPath, err := ioutil.TempDir("", "valid-scores")
	if err != nil {
		return nil, err
	}
	invalidScoresPath, err := ioutil.TempDir("", "invalid-scores")
	if err != nil {
		return nil, err
	}
	result := new(tempCache)
	fc, fcErr := MakeFileCache(validScoresPath, invalidScoresPath, false, keys, initialTotalCount, simulate)
	if fcErr != nil {
		return nil, fcErr
	}
	result.fileCache = fc.(*fileCache)
	result.removeAllOnClose = removeAllOnClose
	return result, nil
}

func (c tempCache) Score(url *url.URL) (score.LinkScores, error) {
	return c.fileCache.Score(url)
}

func (c tempCache) Get(url *url.URL) (score.LinkScores, error) {
	return c.fileCache.Get(url)
}

func (c tempCache) Find(url *url.URL) (scores score.LinkScores, found bool, expired bool, err error) {
	return c.Find(url)
}

func (c tempCache) Save(scores score.LinkScores, autoExpire time.Duration) error {
	return c.Save(scores, autoExpire)
}

func (c tempCache) Close() error {
	if c.removeAllOnClose {
		os.RemoveAll(c.fileCache.validScoresPath)
		os.RemoveAll(c.fileCache.invalidScoresPath)
	}
	return nil
}
