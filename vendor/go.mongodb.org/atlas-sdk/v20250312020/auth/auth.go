package auth

import "golang.org/x/oauth2"

// oauth2 alias
var (
	// HTTPClient is the context key to use with golang.org/x/net/context's
	// WithValue function to associate an *http.Client value with a context.
	HTTPClient = oauth2.HTTPClient
	// ReuseTokenSource returns a TokenSource which repeatedly returns the
	// same token as long as it's valid, starting with t.
	// When its cached token is invalid, a new token is obtained from src.
	//
	// ReuseTokenSource is typically used to reuse tokens from a cache
	// (such as a file on disk) between runs of a program, rather than
	// obtaining new tokens unnecessarily.
	//
	// The initial token t may be nil, in which case the TokenSource is
	// wrapped in a caching version if it isn't one already. This also
	// means it's always safe to wrap ReuseTokenSource around any other
	// TokenSource without adverse effects.
	ReuseTokenSource = oauth2.ReuseTokenSource
	// NewClient creates an *http.Client from a Context and TokenSource.
	// The returned client is not valid beyond the lifetime of the context.
	//
	// Note that if a custom *http.Client is provided via the Context it
	// is used only for token acquisition and is not used to configure the
	// *http.Client returned from NewClient.
	//
	// As a special case, if src is nil, a non-OAuth2 client is returned
	// using the provided context. This exists to support related OAuth2
	// packages.
	NewClient = oauth2.NewClient
)
