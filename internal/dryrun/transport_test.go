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

package dryrun

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestDryRunTransport(t *testing.T) {
	for _, tc := range []struct {
		name    string
		req     *http.Request
		wantErr string
	}{
		{
			name: "GET request",
			req: &http.Request{
				Method: http.MethodGet,
			},
		},
		{
			name: "unknown verb",
			req: &http.Request{
				Method: "UNKNOWN",
				URL:    &url.URL{Path: "/test"},
			},
			wantErr: "DryRun event: Would execute UNKNOWN /test",
		},
		{
			name: "POST request",
			req: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/test"},
			},
			wantErr: "DryRun event: Would create (POST) /test",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			nopDelegate := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				return nil, nil
			})
			transport := NewDryRunTransport(nopDelegate)
			_, err := transport.RoundTrip(tc.req)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tc.wantErr, gotErr)
		})
	}
}
