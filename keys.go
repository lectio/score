package score

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

// Keys describes the different ways link keys can be generated
type Keys interface {
	ScoreKeyForURL(url *url.URL) string
	ScoreKeyForURLText(urlText string) string
}

// MakeDefaultKeys creates a default key generator for links
func MakeDefaultKeys() Keys {
	result := new(defaultKeys)
	return result
}

type defaultKeys struct {
}

func (k defaultKeys) ScoreKeyForURL(url *url.URL) string {
	if url != nil {
		return k.ScoreKeyForURLText(url.String())
	}
	return "url_is_nil_in_ScoreKeyForURL"
}

func (k defaultKeys) ScoreKeyForURLText(urlText string) string {
	h := sha1.New()
	h.Write([]byte(urlText))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
