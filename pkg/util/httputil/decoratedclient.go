package httputil

import "net/http"

type ClientOpt func(*http.Client) error

// DecorateClient performs some custom modifications to an http Client
func DecorateClient(c *http.Client, opts ...ClientOpt) (*http.Client, error) {
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}
