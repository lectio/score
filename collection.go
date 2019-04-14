package score

import (
	"fmt"
	"io"
	"net/url"
	"sync"
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
type TargetsIteratorFn func() (startIndex int, endIndex int, keys Keys, retrievalFn TargetsIteratorRetrievalFn)

// TargetsIteratorRetrievalFn is a function that picks up a URL at a particular collection iterator index
type TargetsIteratorRetrievalFn func(index int) (url *url.URL, err error)

// Collection is list of scored links
type Collection interface {
	ScoredLinks() []*AggregatedLinkScores                       // includes valid and invalid scores
	ValidScoredLinks() []*AggregatedLinkScores                  // only valid scores
	ScoredLink(targetURLUniqueKey string) *AggregatedLinkScores // specific link score
	Errors() []error
}

type defaultCollection struct {
	sync.RWMutex
	simulated        bool
	scoredLinksMap   map[string]*AggregatedLinkScores
	scoredLinks      []*AggregatedLinkScores
	validScoredLinks []*AggregatedLinkScores
	errors           []error
}

// MakeCollection creates a new defaultCollection
func MakeCollection(iterator TargetsIteratorFn, pr ProgressReporter, simulate bool) Collection {
	result := new(defaultCollection)
	result.simulated = simulate
	result.scoredLinksMap = make(map[string]*AggregatedLinkScores)

	startIndex, endIndex, keys, getTarget := iterator()
	ch := make(chan int)
	for i := startIndex; i <= endIndex; i++ {
		url, err := getTarget(i)
		go result.score(i, ch, url, keys, err, simulate)
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

func (c *defaultCollection) score(index int, ch chan<- int, url *url.URL, keys Keys, getTargetErr error, simulate bool) {
	c.Lock()
	key := keys.ScoreKeyForURL(url)
	if getTargetErr != nil {
		c.errors = append(c.errors, fmt.Errorf("skipping scoring of item %d: %v", index, getTargetErr))
	} else if url == nil || len(key) == 0 {
		c.errors = append(c.errors, fmt.Errorf("skipping scoring of item %d: url %q, key: %q", index, url, key))
	} else {
		scores := GetAggregatedLinkScores(url, keys, -1, simulate)
		c.scoredLinksMap[key] = scores
		c.scoredLinks = append(c.scoredLinks, scores)
		if scores.IsValid() {
			c.validScoredLinks = append(c.validScoredLinks, scores)
		}
	}
	c.Unlock()
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

func (c defaultCollection) ScoredLinks() []*AggregatedLinkScores {
	return c.scoredLinks
}

func (c defaultCollection) ValidScoredLinks() []*AggregatedLinkScores {
	return c.validScoredLinks
}

func (c defaultCollection) ScoredLink(targetURLUniqueKey string) *AggregatedLinkScores {
	return c.scoredLinksMap[targetURLUniqueKey]
}

func (c defaultCollection) Errors() []error {
	return c.errors
}
