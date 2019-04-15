package score

// LinkScores instances score a given link (by running an API or other computation)
type LinkScores interface {
	SourceID() string
	TargetURL() string
	IsValid() bool
	SharesCount() int
	CommentsCount() int
}
