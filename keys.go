package score

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

// Keys describes the different ways link keys can be generated
type Keys interface {
	ScoresKeyForURL(url *url.URL) string
	ScoresKeyForURLText(urlText string) string
	ScoresKey(scores LinkScores) string
}

// MakeDefaultKeys creates a default key generator for links
func MakeDefaultKeys() Keys {
	result := new(defaultKeys)
	return result
}

type defaultKeys struct {
}

func (k defaultKeys) ScoresKeyForURL(url *url.URL) string {
	if url != nil {
		return k.ScoresKeyForURLText(url.String())
	}
	return "url_is_nil_in_ScoreKeyForURL"
}

func (k defaultKeys) ScoresKeyForURLText(urlText string) string {
	// TODO: consider adding a key cache since sha1 is compute intensive
	h := sha1.New()
	h.Write([]byte(urlText))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (k defaultKeys) ScoresKey(scores LinkScores) string {
	return k.ScoresKeyForURLText(scores.TargetURL())
}
