package score

import (
	"net/url"
)

// LinkScores instances score a given link (by running an API or other computation)
type LinkScores interface {
	SourceID() string
	TargetURL() string
	IsValid() bool
	SharesCount() int
	CommentsCount() int
}

// Lifecycle defines common creation / destruction methods
type Lifecycle interface {
	ScoreLink(*url.URL) (LinkScores, error)
}

// Reader defines common reader methods
type Reader interface {
	ReadLinkScores(*url.URL) (LinkScores, error)
	FindLinkScores(*url.URL) (LinkScores, bool, error)
}

// Writer defines common writer methods
type Writer interface {
	WriteLinkScores(LinkScores) error
}

// Store pulls together all the lifecyle, reader, and writer methods
type Store interface {
	Lifecycle
	Reader
	Writer
}
