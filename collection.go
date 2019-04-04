package score

import (
	"net/url"

	"gopkg.in/cheggaaa/pb.v1"
)

// TargetsIteratorFn is a function that computes the collection iteration start / end indices
type TargetsIteratorFn func() (startIndex int, endIndex int, retrievalFn TargetsIteratorRetrievalFn)

// TargetsIteratorRetrievalFn is a function that picks up a URL at a particular collection iterator index
type TargetsIteratorRetrievalFn func(index int) (ok bool, url *url.URL, globallyUniqueKey string)

// Collection is list of scored links
type Collection interface {
	ScoredLinks() []*AggregatedLinkScores                       // includes valid and invalid scores
	ValidScoredLinks() []*AggregatedLinkScores                  // only valid scores
	ScoredLink(targetURLUniqueKey string) *AggregatedLinkScores // specific link score
}

type defaultCollection struct {
	simulated        bool
	scoredLinksMap   map[string]*AggregatedLinkScores
	scoredLinks      []*AggregatedLinkScores
	validScoredLinks []*AggregatedLinkScores
}

// MakeCollection creates a new defaultCollection
func MakeCollection(getBoundaries TargetsIteratorFn, verbose bool, simulate bool) Collection {
	result := new(defaultCollection)
	result.simulated = simulate
	result.scoredLinksMap = make(map[string]*AggregatedLinkScores)

	startIndex, endIndex, getTarget := getBoundaries()
	var bar *pb.ProgressBar
	if verbose {
		bar = pb.StartNew(endIndex - startIndex + 1)
		bar.ShowCounters = true
	}
	ch := make(chan int)
	for i := startIndex; i <= endIndex; i++ {
		ok, url, key := getTarget(i)
		if ok {
			// because scores can take time, spin up a bunch concurrently
			go result.score(i, ch, url, key, simulate)
		}
	}

	for i := startIndex; i <= endIndex; i++ {
		_ = <-ch
		if verbose {
			bar.Increment()
		}
	}

	return result
}

func (c *defaultCollection) score(index int, ch chan<- int, url *url.URL, globallyUniqueKey string, simulate bool) {
	scores := GetAggregatedLinkScores(url, globallyUniqueKey, -1, simulate)
	c.scoredLinksMap[globallyUniqueKey] = scores
	c.scoredLinks = append(c.scoredLinks, scores)
	if scores.IsValid() {
		c.validScoredLinks = append(c.scoredLinks, scores)
	}
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