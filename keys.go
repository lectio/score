package score

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

// KeysForURL describes the different ways URL keys can be generated for URLs, this interface is separated because
// it's useful outside of the package - for use internally in this package, use Keys interface below
type KeysForURL interface {
	PrimaryKeyForURL(url *url.URL) string
	PrimaryKeyForURLText(urlText string) string
}

// Keys describes the different ways URL keys can be generated
type Keys interface {
	KeysForURL
	PrimaryKeyForScores(scores LinkScores) string
}

// MakeDefaultKeys creates a default key generator for links
func MakeDefaultKeys() Keys {
	result := new(defaultKeys)
	return result
}

type defaultKeys struct {
}

func (k defaultKeys) PrimaryKeyForURL(url *url.URL) string {
	if url != nil {
		return k.PrimaryKeyForURLText(url.String())
	}
	return "url_is_nil_in_ScoreKeyForURL"
}

func (k defaultKeys) PrimaryKeyForURLText(urlText string) string {
	// TODO: consider adding a key cache since sha1 is compute intensive
	h := sha1.New()
	h.Write([]byte(urlText))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func (k defaultKeys) PrimaryKeyForScores(scores LinkScores) string {
	return k.PrimaryKeyForURLText(scores.TargetURL())
}
