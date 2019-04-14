package container

import (
	"fmt"
	"io"
	"net/url"
	"sync"

	"github.com/lectio/score"
	"github.com/lectio/score/cache"
)

// ProgressReporter is sent to this package's methods if activity progress reporting is expected
type ProgressReporter interface {
	IsProgressReportingRequested() bool
	StartReportableActivity(expectedItems int)
	StartReportableReaderActivityInBytes(exepectedBytes int64, inputReader io.Reader) io.Reader
	IncrementReportableActivityProgress()
	IncrementReportableActivityProgressBy(incrementBy int)
	CompleteReportableActivityProgress(summary string)
}

// TargetsIteratorFn is a function that computes the collection iteration start / end indices
type TargetsIteratorFn func() (startIndex int, endIndex int, keys score.Keys, retrievalFn TargetsIteratorRetrievalFn)

// TargetsIteratorRetrievalFn is a function that picks up a URL at a particular collection iterator index
type TargetsIteratorRetrievalFn func(index int) (url *url.URL, err error)

// Collection is list of scored links
type Collection interface {
	ScoredLinks() []score.LinkScores                       // includes valid and invalid scores
	ValidScoredLinks() []score.LinkScores                  // only valid scores
	ScoredLink(targetURLUniqueKey string) score.LinkScores // specific link score
	Errors() []error
}

type defaultCollection struct {
	mutex            *sync.RWMutex
	cache            cache.Cache
	scoredLinksMap   map[string]score.LinkScores
	scoredLinks      []score.LinkScores
	validScoredLinks []score.LinkScores
	errors           []error
}

// MakeCollection creates a new defaultCollection
func MakeCollection(cache cache.Cache, iterator TargetsIteratorFn, pr ProgressReporter) Collection {
	result := new(defaultCollection)
	result.mutex = new(sync.RWMutex)
	result.cache = cache
	result.scoredLinksMap = make(map[string]score.LinkScores)

	startIndex, endIndex, keys, getTarget := iterator()
	ch := make(chan int)
	for i := startIndex; i <= endIndex; i++ {
		url, err := getTarget(i)
		go result.score(i, ch, url, keys, err)
	}

	if pr != nil && pr.IsProgressReportingRequested() {
		pr.StartReportableActivity(endIndex - startIndex + 1)
	}

	for i := startIndex; i <= endIndex; i++ {
		_ = <-ch
		if pr != nil && pr.IsProgressReportingRequested() {
			pr.IncrementReportableActivityProgress()
		}
	}

	if pr != nil && pr.IsProgressReportingRequested() {
		pr.CompleteReportableActivityProgress(fmt.Sprintf("Completed scoring %d items in iterator: %d in map, %d in list, %d valid", endIndex-startIndex+1, len(result.scoredLinksMap), len(result.scoredLinks), len(result.validScoredLinks)))
	}

	return result
}

func (c *defaultCollection) score(index int, ch chan<- int, url *url.URL, keys score.Keys, getTargetErr error) {
	c.mutex.Lock()
	key := keys.PrimaryKeyForURL(url)
	if getTargetErr != nil {
		c.errors = append(c.errors, fmt.Errorf("skipping scoring of item %d: %v", index, getTargetErr))
	} else if url == nil || len(key) == 0 {
		c.errors = append(c.errors, fmt.Errorf("skipping scoring of item %d: url %q, key: %q", index, url, key))
	} else {
		scores, getErr := c.cache.Get(url)
		if getErr == nil {
			c.scoredLinksMap[key] = scores
			c.scoredLinks = append(c.scoredLinks, scores)
			if scores.IsValid() {
				c.validScoredLinks = append(c.validScoredLinks, scores)
			}
		} else {
			c.errors = append(c.errors, fmt.Errorf("skipping scoring of item %d: %v", index, getErr))
		}
	}
	c.mutex.Unlock()
	ch <- index
}

// if url == nil || len(key) == 0 {
// 	result.errors = append(result.errors, fmt.Errorf("skipping scoring of item %d: url %q, key: %q", i, url, key))
// 	continue
// }
// // because scores can take time, spin up a bunch concurrently
// go result.score(i, ch, url, key, err, simulate)
// scoreConcurrentCount++
// } else {

func (c defaultCollection) ScoredLinks() []score.LinkScores {
	return c.scoredLinks
}

func (c defaultCollection) ValidScoredLinks() []score.LinkScores {
	return c.validScoredLinks
}

func (c defaultCollection) ScoredLink(targetURLUniqueKey string) score.LinkScores {
	return c.scoredLinksMap[targetURLUniqueKey]
}

func (c defaultCollection) Errors() []error {
	return c.errors
}
