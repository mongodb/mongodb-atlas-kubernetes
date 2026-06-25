// Copyright 2026 MongoDB Inc
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

package serviceaccounttoken

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312021/auth/clientcredentials"
	"golang.org/x/oauth2"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newOAuthHTTPClient(t *testing.T, fn roundTripFunc) *http.Client {
	t.Helper()
	return &http.Client{Transport: fn}
}

func TestAtlasTokenProvider_FetchToken_Success(t *testing.T) {
	testCases := []struct {
		name        string
		atlasDomain string
	}{
		{name: "domain without trailing slash", atlasDomain: "https://atlas.example.com"},
		{name: "domain with trailing slash", atlasDomain: "https://atlas.example.com/"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expiry := time.Now().UTC().Add(1 * time.Hour).Truncate(time.Second)
			provider := NewAtlasTokenProvider(tc.atlasDomain)

			client := newOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
				require.Equal(t, http.MethodPost, req.Method)
				require.Equal(t, clientcredentials.TokenAPIPath, req.URL.Path)

				err := req.ParseForm()
				require.NoError(t, err)
				assert.Equal(t, "client_credentials", req.Form.Get("grant_type"))

				user, pass, ok := req.BasicAuth()
				require.True(t, ok, "oauth2 client credentials request should use basic auth")
				assert.Equal(t, "my-client-id", user)
				assert.Equal(t, "my-client-secret", pass)

				body := fmt.Sprintf(`{"access_token":"token-123","token_type":"Bearer","expires_in":%d,"expiry":"%s"}`,
					time.Until(expiry)/time.Second,
					expiry.Format(time.RFC3339),
				)

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			})

			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
			token, tokenExpiry, err := provider.FetchToken(ctx, "my-client-id", "my-client-secret")

			require.NoError(t, err)
			assert.Equal(t, "token-123", token)
			assert.True(t, tokenExpiry.After(time.Now().UTC()), "returned expiry should be in the future")
		})
	}
}

func TestAtlasTokenProvider_FetchToken_HTTPTransportError(t *testing.T) {
	provider := NewAtlasTokenProvider("https://atlas.example.com")

	client := newOAuthHTTPClient(t, func(_ *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("boom")
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	token, tokenExpiry, err := provider.FetchToken(ctx, "my-client-id", "my-client-secret")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to acquire OAuth token")
	assert.Contains(t, err.Error(), "boom")
	assert.Empty(t, token)
	assert.True(t, tokenExpiry.IsZero())
}

func TestAtlasTokenProvider_FetchToken_OAuthErrorResponse(t *testing.T) {
	provider := NewAtlasTokenProvider("https://atlas.example.com")

	client := newOAuthHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		require.Equal(t, clientcredentials.TokenAPIPath, req.URL.Path)

		values := url.Values{}
		values.Set("error", "invalid_client")
		values.Set("error_description", "client authentication failed")

		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			Body:       io.NopCloser(strings.NewReader(values.Encode())),
		}, nil
	})

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	token, tokenExpiry, err := provider.FetchToken(ctx, "bad-id", "bad-secret")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to acquire OAuth token")
	assert.Contains(t, err.Error(), "invalid_client")
	assert.Empty(t, token)
	assert.True(t, tokenExpiry.IsZero())
}
