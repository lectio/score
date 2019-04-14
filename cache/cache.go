package cache

import (
	"net/url"
	"time"

	"github.com/lectio/score"
)

// Cache allows storing and retrieving scores from disk, RAM, etc.
type Cache interface {
	Score(url *url.URL) (score.LinkScores, error)
	Get(url *url.URL) (score.LinkScores, error)
	Find(url *url.URL) (scores score.LinkScores, found bool, expired bool, err error)
	Save(scores score.LinkScores, autoExpire time.Duration) error
	Close() error
}
