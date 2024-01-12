package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// LoggingTransport is the option adding logging capability to an http Client
func LoggingTransport(log *zap.SugaredLogger) ClientOpt {
	return func(c *http.Client) error {
		c.Transport = &loggedRoundTripper{rt: c.Transport, log: log, logBody: false}
		return nil
	}
}

type loggedRoundTripper struct {
	rt      http.RoundTripper
	log     *zap.SugaredLogger
	logBody bool
}

func (l *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	startTime := time.Now()
	if l.logBody && request.Body != nil {
		l.writeBodyToLog(request.GetBody)
	}
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

func (l loggedRoundTripper) writeBodyToLog(body func() (io.ReadCloser, error)) {
	bodyCopy, err := body()
	if err != nil {
		l.log.Error(err)
		return
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(bodyCopy)
	if err != nil {
		l.log.Errorf("Failed to read to buffer: %s", err)
		return
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, buf.Bytes(), "", "\t")
	if err != nil {
		l.log.Errorf("JSON parse error: %s", err)
		return
	}

	l.log.Debugf(">> \n%s \n", prettyJSON.String())
}
