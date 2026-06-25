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

package serviceaccounttoken

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go.mongodb.org/atlas-sdk/v20250312021/auth/clientcredentials"
)

// TokenProvider abstracts the OAuth token acquisition so it can be mocked in tests.
type TokenProvider interface {
	FetchToken(ctx context.Context, clientID, clientSecret string) (token string, expiry time.Time, err error)
}

// AtlasTokenProvider implements TokenProvider using the Atlas SDK's client credentials flow.
type AtlasTokenProvider struct {
	atlasDomain string
}

func NewAtlasTokenProvider(atlasDomain string) *AtlasTokenProvider {
	return &AtlasTokenProvider{atlasDomain: atlasDomain}
}

func (p *AtlasTokenProvider) FetchToken(ctx context.Context, clientID, clientSecret string) (string, time.Time, error) {
	cfg := clientcredentials.NewConfig(clientID, clientSecret)

	baseURL, err := url.Parse(p.atlasDomain)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid Atlas domain %q: %w", p.atlasDomain, err)
	}

	tokenURL := baseURL.ResolveReference(&url.URL{Path: clientcredentials.TokenAPIPath})
	cfg.TokenURL = tokenURL.String()

	token, err := cfg.Token(ctx)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to acquire OAuth token: %w", err)
	}

	return token.AccessToken, token.Expiry, nil
}
