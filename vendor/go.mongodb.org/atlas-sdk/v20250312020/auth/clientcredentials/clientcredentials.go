package clientcredentials

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"

	"go.mongodb.org/atlas-sdk/v20250312020/auth"
	"go.mongodb.org/atlas-sdk/v20250312020/internal/core"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	// TokenAPIPath for getting OAuth Access Token from server
	TokenAPIPath = "/api/oauth/token" //nolint:gosec //url only
	// serverTokenURL for Token Atlas API
	serverTokenURL = core.DefaultCloudURL + TokenAPIPath
	// RevokeAPIPath for revoking OAuth Access Token from server
	RevokeAPIPath = "/api/oauth/revoke"

	// serverRevokeURL for Revoke Atlas API
	serverRevokeURL = core.DefaultCloudURL + RevokeAPIPath
	userAgent       = "User-Agent"
)

func NewConfig(clientID, clientSecret string) *Config {
	c := &Config{}
	c.ClientID = clientID
	c.ClientSecret = clientSecret
	c.RevokeURL = serverRevokeURL
	c.TokenURL = serverTokenURL
	c.AuthStyle = oauth2.AuthStyleInHeader
	c.userAgent = core.DefaultUserAgent

	return c
}

// Config describes a 2-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
//
// NOTE: Config values are used only internally
// and should not be overridden by clients
type Config struct {
	clientcredentials.Config
	RevokeURL string
	userAgent string
}

func (c *Config) Client(ctx context.Context) *http.Client {
	client := c.Config.Client(ctx)
	client.Transport = &Transport{
		Base:      client.Transport,
		UserAgent: core.DefaultUserAgent,
	}
	return client
}

// RevokeToken revokes OAuth Token
func (c *Config) RevokeToken(ctx context.Context, t *auth.Token) error {
	if c.RevokeURL == "" {
		return errors.New("endpoint missing RevokeURL")
	}
	if !t.Valid() {
		return nil // nothing to do
	}
	v := url.Values{
		"token":           {t.AccessToken},
		"token_type_hint": {"access_token"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.RevokeURL, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(url.QueryEscape(c.ClientID), url.QueryEscape(c.ClientSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set(userAgent, c.userAgent)

	client := auth.NewClient(ctx, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusTooManyRequests {
			msg, _ := io.ReadAll(resp.Body)
			formattedMessage := fmt.Sprintf("%s %s: HTTP %v Detail: %v Reason: %v",
				http.MethodPost, c.RevokeURL, resp.StatusCode,
				"Token Revocation request was rate limited", string(msg))
			return errors.New(formattedMessage)
		}
		formattedMessage := fmt.Sprintf("%s %s: HTTP %v Detail: %v Reason: %v",
			http.MethodPost, c.RevokeURL, resp.StatusCode,
			"Failed to revoke Access Token when fetching new OAuth Token from remote server",
			resp.Header.Get("www-authenticate"))
		return errors.New(formattedMessage)
	}
	return nil
}
