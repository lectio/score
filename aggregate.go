package score

import "net/url"

// AggregatedLinkScores computes aggregate scores from multiple link scorers
type AggregatedLinkScores struct {
	MachineName            string       `json:"scorer"`
	HumanName              string       `json:"scorerName"`
	Simulated              bool         `json:"isSimulated,omitempty"`
	URL                    string       `json:"url"`
	GloballyUniqueKey      string       `json:"uniqueKey"`
	AggregateSharesCount   int          `json:"aggregateSharesCount"`
	AggregateCommentsCount int          `json:"aggregateCommentsCount"`
	Scores                 []LinkScores `json:"scores"`
}

// GetAggregatedLinkScores returns a multiple scores structure
func GetAggregatedLinkScores(url *url.URL, keys Keys, initialTotalCount int, simulate bool) *AggregatedLinkScores {
	result := new(AggregatedLinkScores)
	result.MachineName = "aggregate"
	result.HumanName = "Aggregate"
	result.Simulated = simulate
	result.URL = url.String()
	result.GloballyUniqueKey = keys.PrimaryKeyForURL(url)

	if fb, fbErr := GetFacebookLinkScoresForURL(url, keys, simulate); fbErr == nil {
		result.Scores = append(result.Scores, fb)
	}
	if li, liErr := GetLinkedInLinkScoresForURL(url, keys, simulate); liErr == nil {
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

// ScorerMachineName returns the name of the scoring engine suitable for machine processing
func (a AggregatedLinkScores) ScorerMachineName() string {
	return a.MachineName
}

// ScorerHumanName returns the name of the scoring engine suitable for humans
func (a AggregatedLinkScores) ScorerHumanName() string {
	return a.HumanName
}

// Scorer returns the scoring engine information
func (a AggregatedLinkScores) Scorer() LinkScorer {
	return a
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
