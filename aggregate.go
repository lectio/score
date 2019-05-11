package score

import (
	"net/http"
	"net/url"
)

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

	issues []Issue
}

// GetAggregatedLinkScores returns a multiple scores structure
func GetAggregatedLinkScores(url *url.URL, client *http.Client, keys Keys, initialTotalCount int, simulate bool) *AggregatedLinkScores {
	result := new(AggregatedLinkScores)
	result.MachineName = "aggregate"
	result.HumanName = "Aggregate"
	result.Simulated = simulate
	result.URL = url.String()
	result.GloballyUniqueKey = keys.PrimaryKeyForURL(url)

	if fb, fbErr := GetFacebookLinkScoresForURL(url, client, keys, simulate); fbErr == nil {
		result.Scores = append(result.Scores, fb)
		if fb.Issues() != nil {
			for _, issue := range fb.Issues().ErrorsAndWarnings() {
				result.issues = append(result.issues, issue)
			}
		}
	}
	if li, liErr := GetLinkedInLinkScoresForURL(url, client, keys, simulate); liErr == nil {
		result.Scores = append(result.Scores, li)
		if li.Issues() != nil {
			for _, issue := range li.Issues().ErrorsAndWarnings() {
				result.issues = append(result.issues, issue)
			}
		}
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

// SourceID returns the name of the scoring engine
func (a AggregatedLinkScores) SourceID() string {
	return a.MachineName
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

// Issues contains all the problems detected in scoring
func (a AggregatedLinkScores) Issues() Issues {
	return a
}

// ErrorsAndWarnings contains the problems in this link plus satisfies the Link.Issues interface
func (a AggregatedLinkScores) ErrorsAndWarnings() []Issue {
	return a.issues
}

// IssueCounts returns the total, errors, and warnings counts
func (a AggregatedLinkScores) IssueCounts() (uint, uint, uint) {
	if a.issues == nil {
		return 0, 0, 0
	}
	var errors, warnings uint
	for _, i := range a.issues {
		if i.IsError() {
			errors++
		} else {
			warnings++
		}
	}
	return uint(len(a.issues)), errors, warnings
}

// HandleIssues loops through each issue and calls a particular handler
func (a AggregatedLinkScores) HandleIssues(errorHandler func(Issue), warningHandler func(Issue)) {
	if a.issues == nil {
		return
	}
	for _, i := range a.issues {
		if i.IsError() && errorHandler != nil {
			errorHandler(i)
		}
		if i.IsWarning() && warningHandler != nil {
			warningHandler(i)
		}
	}
}
