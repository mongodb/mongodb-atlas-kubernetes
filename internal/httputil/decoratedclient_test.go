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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_DecorateClient(t *testing.T) {
	httpClient := &http.Client{Transport: http.DefaultTransport}
	withDigest := Digest("publiApi", "privateApi")
	withLogging := LoggingTransport(zap.S())

	decorated, err := DecorateClient(&http.Client{Transport: http.DefaultTransport}, withDigest, withLogging)
	a := assert.New(t)
	a.NoError(err)
	a.Equal(httpClient.Timeout, decorated.Timeout)
	a.Equal(httpClient.Jar, decorated.Jar)
	a.NotNil(decorated.Transport)

	// not going deeper here, just need to confirm that transport was changed
	a.NotEqual(t, httpClient.Transport, decorated.Transport)
}

type dummyTripper struct{}

// RoundTrip implements http.RoundTripper.
func (*dummyTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func Test_DecorateClientCustomTransport(t *testing.T) {
	dt := &dummyTripper{}
	withTransport := CustomTransport(dt)

	decorated, err := DecorateClient(&http.Client{Transport: http.DefaultTransport}, withTransport)
	a := assert.New(t)
	a.NoError(err)
	a.Equal(decorated.Transport, dt)
}
