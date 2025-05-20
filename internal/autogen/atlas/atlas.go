package atlas

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/mongodb-forks/digest"
	admin20231115008 "go.mongodb.org/atlas-sdk/v20231115008/admin"
	admin20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	admin20250312002 "go.mongodb.org/atlas-sdk/v20250312002/admin"
)

type ctxKey int

const (
	ctxClientSet ctxKey = iota
)

type ClientSet struct {
	SdkClient20250312002 *admin20250312002.APIClient
	SdkClient20231115008 *admin20231115008.APIClient
	SdkClient20241113001 *admin20241113001.APIClient
}

func FromContext(ctx context.Context) *ClientSet {
	if v, ok := ctx.Value(ctxClientSet).(*ClientSet); ok {
		return v
	}
	return nil
}

func NewContext(ctx context.Context, clientSet *ClientSet) context.Context {
	return context.WithValue(ctx, ctxClientSet, clientSet)
}

func NewClientSet() (*ClientSet, error) {
	var transport http.RoundTripper = digest.NewTransport(os.Getenv("MCLI_PUBLIC_API_KEY"), os.Getenv("MCLI_PRIVATE_API_KEY"))
	httpClient := &http.Client{Transport: transport}

	atlas20231115008Client, err := admin20231115008.NewClient(
		admin20231115008.UseBaseURL(os.Getenv("MCLI_OPS_MANAGER_URL")),
		admin20231115008.UseHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create v20231115008 client: %w", err)
	}

	atlas20241113001Client, err := admin20241113001.NewClient(
		admin20241113001.UseBaseURL(os.Getenv("MCLI_OPS_MANAGER_URL")),
		admin20241113001.UseHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create v20241113001 client: %w", err)
	}

	atlas20250312002Client, err := admin20250312002.NewClient(
		admin20250312002.UseBaseURL(os.Getenv("MCLI_OPS_MANAGER_URL")),
		admin20250312002.UseHTTPClient(httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create v20250312002 client: %w", err)
	}

	return &ClientSet{
		SdkClient20231115008: atlas20231115008Client,
		SdkClient20241113001: atlas20241113001Client,
		SdkClient20250312002: atlas20250312002Client,
	}, nil
}
