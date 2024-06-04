package contract

import (
	"context"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

func NewVersionedClient(ctx context.Context) (*admin.APIClient, error) {
	domain := os.Getenv("MCLI_OPS_MANAGER_URL")
	pubKey := os.Getenv("MCLI_PUBLIC_API_KEY")
	prvKey := os.Getenv("MCLI_PRIVATE_API_KEY")
	return atlas.NewClient(domain, pubKey, prvKey)
}

func MustVersionedClient(t *testing.T, ctx context.Context) *admin.APIClient {
	client, err := NewVersionedClient(ctx)
	if err != nil {
		t.Fatalf("Failed to get Atlas versioned client: %v", err)
	}
	return client
}
