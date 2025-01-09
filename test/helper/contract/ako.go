package contract

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
)

func DefaultAtlasProject(name string) client.Object {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       akov2.AtlasProjectSpec{Name: name},
	}
}

func newVersionedClient(ctx context.Context) (*admin.APIClient, error) {
	domain := os.Getenv("MCLI_OPS_MANAGER_URL")
	pubKey := os.Getenv("MCLI_PUBLIC_API_KEY")
	prvKey := os.Getenv("MCLI_PRIVATE_API_KEY")
	client, err := atlas.NewClient(domain, pubKey, prvKey)
	if err != nil {
		return nil, fmt.Errorf("failed to setup Atlas Client: %w", err)
	}
	_, _, err = client.ProjectsApi.ListProjects(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("non working Atlas Client: %w", err)
	}
	return client, err
}

func mustCreateVersionedAtlasClient(ctx context.Context) *admin.APIClient {
	client, err := newVersionedClient(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to create an Atlas versioned client: %v", err))
	}
	return client
}

func mustCreateVersionedAtlasClientSet(ctx context.Context) *atlas.ClientSet {
	domain := os.Getenv("MCLI_OPS_MANAGER_URL")
	pubKey := os.Getenv("MCLI_PUBLIC_API_KEY")
	prvKey := os.Getenv("MCLI_PRIVATE_API_KEY")
	c2024, err := adminv20241113001.NewClient(
		adminv20241113001.UseBaseURL(domain),
		adminv20241113001.UseDigestAuth(pubKey, prvKey),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create an Atlas versioned client: %v", err))
	}
	_, _, err = c2024.ProjectsApi.ListProjects(ctx).Execute()
	if err != nil {
		panic(fmt.Sprintf("non working Atlas Client: %v", err))
	}

	return &atlas.ClientSet{
		SdkClient20241113001: c2024,
	}
}

func globalSecret(namespace string) client.Object {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mongodb-atlas-operator-api-key",
			Namespace: namespace,
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"orgId":         ([]byte)(os.Getenv("MCLI_ORG_ID")),
			"publicApiKey":  ([]byte)(os.Getenv("MCLI_PUBLIC_API_KEY")),
			"privateApiKey": ([]byte)(os.Getenv("MCLI_PRIVATE_API_KEY")),
		},
	}
}
