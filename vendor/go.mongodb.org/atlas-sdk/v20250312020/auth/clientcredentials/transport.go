package clientcredentials

import (
	"net/http"
)

// Transport supplies custom user agent to token requests
type Transport struct {
	Base      http.RoundTripper
	UserAgent string
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.UserAgent != "" {
		req.Header.Set("User-Agent", t.UserAgent)
	}
	return t.base().RoundTrip(req)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}
