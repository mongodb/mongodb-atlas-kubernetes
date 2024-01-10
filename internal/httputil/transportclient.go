package httputil

import "net/http"

// CustomTransport is the option adding a custom transport on a http Client
func CustomTransport(t http.RoundTripper) ClientOpt {
	return func(c *http.Client) error {
		c.Transport = t
		return nil
	}
}
