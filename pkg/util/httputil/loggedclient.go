package httputil

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewLoggingClient returns the http.Client with wrapped Transport that is capable of logging request
func NewLoggingClient(client http.Client, log *zap.SugaredLogger) http.Client {
	return http.Client{
		Transport:     &loggedRoundTripper{rt: client.Transport, log: log},
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}
}

type loggedRoundTripper struct {
	rt  http.RoundTripper
	log *zap.SugaredLogger
}

func (l *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	startTime := time.Now()
	response, err := l.rt.RoundTrip(request)
	duration := time.Since(startTime)
	l.logResponse(request, response, err, duration)
	return response, err
}

// LogResponse logs path, host, status code and duration in milliseconds
// This can be extended further by providing the custom logger pattern but not necessary so far.
func (l loggedRoundTripper) logResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	duration /= time.Millisecond
	if err != nil {
		l.log.Debugf("HTTP Request (%s) %s [time (ms): %d, error=%q]", req.Method, req.URL, duration, err.Error())
	} else {
		l.log.Debugf("HTTP Request (%s) %s [time (ms): %d, status: %d]", req.Method, req.URL, duration, res.StatusCode)
	}
}
