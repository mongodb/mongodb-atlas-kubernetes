// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func NewLoggingTransport(log *zap.SugaredLogger, logBody bool, delegate http.RoundTripper) http.RoundTripper {
	return &loggedRoundTripper{
		rt:      delegate,
		log:     log,
		logBody: logBody,
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
		statusCode := StatusCode(res)
		l.log.Debugf("HTTP Request (%s) %s [time (ms): %d, status: %d]", req.Method, req.URL, duration, statusCode)
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
