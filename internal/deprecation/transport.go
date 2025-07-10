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

package deprecation

import (
	"net/http"

	"go.uber.org/zap"
)

type Transport struct {
	delegate http.RoundTripper
	logger   *zap.Logger
}

func NewLoggingTransport(delegate http.RoundTripper, logger *zap.Logger) *Transport {
	return &Transport{
		delegate: delegate,
		logger:   logger.Named("deprecated"),
	}
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.delegate.RoundTrip(req)
	if resp != nil {
		javaMethod := resp.Header.Get("X-Java-Method")
		deprecation := resp.Header.Get("Deprecation")
		sunset := resp.Header.Get("Sunset")

		if deprecation != "" {
			t.logger.Warn("deprecation", zap.String("type", "deprecation"), zap.String("date", deprecation), zap.String("javaMethod", javaMethod), zap.String("path", req.URL.Path), zap.String("method", req.Method))
		}

		if sunset != "" {
			t.logger.Warn("sunset", zap.String("type", "sunset"), zap.String("date", sunset), zap.String("javaMethod", javaMethod), zap.String("path", req.URL.Path), zap.String("method", req.Method))
		}
	}
	return resp, err
}
