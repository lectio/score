package score

import "net/url"

// AggregatedLinkScores computes aggregate scores from multiple link scorers
type AggregatedLinkScores struct {
	Simulated              bool         `json:"isSimulated,omitempty"`
	URL                    string       `json:"url"`
	GloballyUniqueKey      string       `json:"uniqueKey"`
	AggregateSharesCount   int          `json:"aggregateSharesCount"`
	AggregateCommentsCount int          `json:"aggregateCommentsCount"`
	Scores                 []LinkScores `json:"scores"`
}

// GetAggregatedLinkScores returns a multiple scores structure
func GetAggregatedLinkScores(url *url.URL, globallyUniqueKey string, initialTotalCount int, simulate bool) *AggregatedLinkScores {
	result := new(AggregatedLinkScores)
	result.Simulated = simulate
	result.URL = url.String()

	if fb, fbErr := GetFacebookLinkScoresForURL(url, globallyUniqueKey, simulate); fbErr == nil {
		result.Scores = append(result.Scores, fb)
	}
	if li, liErr := GetLinkedInLinkScoresForURL(url, globallyUniqueKey, simulate); liErr == nil {
		result.Scores = append(result.Scores, li)
	}

	result.AggregateSharesCount = initialTotalCount   // this is often set to -1 to signify "uncalculated" or similar
	result.AggregateCommentsCount = initialTotalCount // this is often set to -1 to signify "uncalculated" or similar
	for _, scorer := range result.Scores {
		if scorer.IsValid() {
			shares := scorer.SharesCount()
			if shares > 0 {
				if result.AggregateSharesCount == initialTotalCount {
					result.AggregateSharesCount = shares
				} else {
					result.AggregateSharesCount += shares
				}
			}

			comments := scorer.CommentsCount()
			if comments > 0 {
				if result.AggregateCommentsCount == initialTotalCount {
					result.AggregateCommentsCount = comments
				} else {
					result.AggregateCommentsCount += comments
				}
			}
		}
	}

	return result
}

// Names returns the identities of the scorer
func (a AggregatedLinkScores) Names() (string, string) {
	return "aggregate", "Aggregate"
}

// TargetURL is the URL that the scores were computed for
func (a AggregatedLinkScores) TargetURL() string {
	return a.URL
}

// IsValid returns true if the FacebookLinkScoreGraphResult object is valid (did not return Facebook error object)
func (a AggregatedLinkScores) IsValid() bool {
	for _, scorer := range a.Scores {
		if !scorer.IsValid() {
			return false
		}
	}
	return true
}

// SharesCount is the count of how many times the given URL was shared by this scorer, -1 if invalid or not available
func (a AggregatedLinkScores) SharesCount() int {
	return a.AggregateSharesCount
}

// CommentsCount is the count of how many times the given URL was commented on, -1 if invalid or not available
func (a AggregatedLinkScores) CommentsCount() int {
	return a.AggregateCommentsCount
}
