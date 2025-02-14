package dryrun

import (
	"net/http"
)

var verbMap = map[string]string{
	http.MethodPost:   "create (" + http.MethodPost + ")",
	http.MethodPut:    "update (" + http.MethodPut + ")",
	http.MethodPatch:  "update (" + http.MethodPatch + ")",
	http.MethodDelete: "delete (" + http.MethodDelete + ")",
}

type DryRunTransport struct {
	Delegate http.RoundTripper
}

func NewDryRunTransport(delegate http.RoundTripper) *DryRunTransport {
	return &DryRunTransport{
		Delegate: delegate,
	}
}

func (t *DryRunTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	switch req.Method {
	case http.MethodGet:
	case http.MethodConnect:
	case http.MethodTrace:
	case http.MethodHead:
	default:
		verb, ok := verbMap[req.Method]
		if !ok {
			verb = "execute " + req.Method
		}
		msg := "Would %v %v"

		return nil, NewDryRunError(msg, verb, req.URL.Path)
	}

	return t.Delegate.RoundTrip(req)
}
