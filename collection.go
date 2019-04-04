package score

import "net/url"

// GatherURLs is a function that returns a set of URLs that should be added to a collection
type GatherURLs func() []*url.URL

// ComputeGloballyUniqueKeyForURL is a function that can compute a unique key for a given URL
type ComputeGloballyUniqueKeyForURL func(url *url.URL) string

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
	uniqueKeyFn      ComputeGloballyUniqueKeyForURL
}

// MakeMutableCollection creates a new defaultCollection
func MakeMutableCollection(urlsFn GatherURLs, uniqueKeyFn ComputeGloballyUniqueKeyForURL, verbose bool, simulate bool) Collection {
	result := new(defaultCollection)
	result.simulated = simulate
	result.scoredLinksMap = make(map[string]*AggregatedLinkScores)
	result.uniqueKeyFn = uniqueKeyFn

	urls := urlsFn()
	for _, url := range urls {
		key := uniqueKeyFn(url)
		scores := GetAggregatedLinkScores(url, key, -1, simulate)
		result.scoredLinksMap[key] = scores
		result.scoredLinks = append(result.scoredLinks, scores)
		if scores.IsValid() {
			result.validScoredLinks = append(result.scoredLinks, scores)
		}
	}

	return result
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
