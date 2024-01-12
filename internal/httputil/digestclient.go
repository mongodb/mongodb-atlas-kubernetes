package httputil

import (
	"net/http"

	"github.com/mongodb-forks/digest"
)

// Digest is the option adding digest authentication capability to an http client
func Digest(publicKey, privateKey string) ClientOpt {
	return func(c *http.Client) error {
		t := &digest.Transport{
			Username:  publicKey,
			Password:  privateKey,
			Transport: c.Transport,
		}
		c.Transport = t
		return nil
	}
}
