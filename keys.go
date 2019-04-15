package score

import (
	"crypto/sha1"
	"fmt"
	"net/url"
)

// Keys describes the different ways URL keys can be generated for URLs
type Keys interface {
	PrimaryKeyForURL(url *url.URL) string
	PrimaryKeyForURLText(urlText string) string
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
	return "url_is_nil_in_PrimaryKeyForURL"
}

func (k defaultKeys) PrimaryKeyForURLText(urlText string) string {
	// TODO: consider adding a key cache since sha1 is compute intensive
	h := sha1.New()
	h.Write([]byte(urlText))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
