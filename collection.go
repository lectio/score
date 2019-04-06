package score

import (
	"fmt"
	"net/url"
	"sync"

	"gopkg.in/cheggaaa/pb.v1"
)

// TargetsIteratorFn is a function that computes the collection iteration start / end indices
type TargetsIteratorFn func() (startIndex int, endIndex int, retrievalFn TargetsIteratorRetrievalFn)

// TargetsIteratorRetrievalFn is a function that picks up a URL at a particular collection iterator index
type TargetsIteratorRetrievalFn func(index int) (url *url.URL, globallyUniqueKey string, err error)

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
func MakeCollection(iterator TargetsIteratorFn, verbose bool, simulate bool) Collection {
	result := new(defaultCollection)
	result.simulated = simulate
	result.scoredLinksMap = make(map[string]*AggregatedLinkScores)

	startIndex, endIndex, getTarget := iterator()
	var bar *pb.ProgressBar
	if verbose {
		bar = pb.StartNew(endIndex - startIndex + 1)
		bar.ShowCounters = true
	}
	ch := make(chan int)
	for i := startIndex; i <= endIndex; i++ {
		url, key, err := getTarget(i)
		if err == nil {
			if url == nil || len(key) == 0 {
				result.errors = append(result.errors, fmt.Errorf("skipping scoring of item %d: url %q, key: %q", i, url, key))
				continue
			}
			// because scores can take time, spin up a bunch concurrently
			go result.score(i, ch, url, key, simulate)
		} else {
			result.errors = append(result.errors, fmt.Errorf("skipping scoring of item %d: %v", i, err))
		}
	}

	for i := startIndex; i <= endIndex; i++ {
		_ = <-ch
		if verbose {
			bar.Increment()
		}
	}

	if verbose {
		bar.FinishPrint(fmt.Sprintf("Completed scoring %d to %d in iterator: %d in map, %d in list, %d valid", startIndex, endIndex, len(result.scoredLinksMap), len(result.scoredLinks), len(result.validScoredLinks)))
	}

	return result
}

func (c *defaultCollection) score(index int, ch chan<- int, url *url.URL, globallyUniqueKey string, simulate bool) {
	scores := GetAggregatedLinkScores(url, globallyUniqueKey, -1, simulate)
	c.Lock()
	c.scoredLinksMap[globallyUniqueKey] = scores
	c.scoredLinks = append(c.scoredLinks, scores)
	if scores.IsValid() {
		c.validScoredLinks = append(c.scoredLinks, scores)
	}
	c.Unlock()
	ch <- index
}

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
