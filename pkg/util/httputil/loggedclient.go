package httputil

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingTransport is the option adding logging capability to an http Client
func LoggingTransport(log *zap.SugaredLogger) ClientOpt {
	return func(c *http.Client) error {
		c.Transport = &loggedRoundTripper{rt: c.Transport, log: log}
		return nil
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
