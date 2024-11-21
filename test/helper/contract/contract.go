package contract

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

type ContractHelper interface {
	AtlasClient() *admin.APIClient
	AddResources(ctx context.Context, timeout time.Duration, resources ...client.Object) error
	ProjectID(ctx context.Context, projectName string) (string, error)
}

type contractTest struct {
	credentials bool
	namespace   string
	resources   []client.Object
	k8sClient   client.Client
	atlasClient *admin.APIClient
}

func (ct *contractTest) cleanup(ctx context.Context) error {
	for i := len(ct.resources) - 1; i >= 0; i-- {
		resource := ct.resources[i]
		if err := ct.k8sClient.Delete(ctx, resource); err != nil {
			return fmt.Errorf("failed to delete contract test pre-set resource: %w", err)
		}
	}
	ct.resources = []client.Object{}
	if ct.namespace != "" {
		if err := ct.k8sClient.Delete(ctx, defaultNamespace(ct.namespace)); err != nil {
			return fmt.Errorf("failed to delete namespace %q: %w", ct.namespace, err)
		}
	}
	return nil
}

func RunGoContractTest(ctx context.Context, t *testing.T, name string, contractTest func(ch ContractHelper)) {
	if !control.Enabled("AKO_CONTRACT_TEST") {
		t.Skip("Skipping contract test as AKO_CONTRACT_TEST is unset")
		return
	}
	ct := newContractTest(ctx)
	defer func() {
		require.NoError(t, ct.cleanup(ctx))
	}()
	t.Run(name, func(t *testing.T) {
		contractTest(ct)
	})
}

func newContractTest(ctx context.Context) *contractTest {
	return &contractTest{
		k8sClient:   mustCreateK8sClient(),
		credentials: false,
		resources:   []client.Object{},
		atlasClient: mustCreateVersionedAtlasClient(ctx),
	}
}

func (ct *contractTest) AtlasClient() *admin.APIClient {
	return ct.atlasClient
}

func (ct *contractTest) AddResources(ctx context.Context, timeout time.Duration, resources ...client.Object) error {
	if !ct.credentials {
		akoTestNamespace := os.Getenv("HELM_AKO_NAMESPACE")
		if err := k8sRecreate(ctx, ct.k8sClient, globalSecret(akoTestNamespace)); err != nil {
			return fmt.Errorf("failed to set AKO namespace: %w", err)
		}
		ct.credentials = true
	}
	if ct.namespace == "" {
		ct.namespace = utils.RandomName("test-ns")
		if err := ct.k8sClient.Create(ctx, defaultNamespace(ct.namespace)); err != nil {
			return fmt.Errorf("failed to create test namespace: %w", err)
		}
	}
	for _, resource := range resources {
		resource.SetNamespace(ct.namespace)
		if err := ct.k8sClient.Create(ctx, resource); err != nil {
			return fmt.Errorf("failed to create resource: %w", err)
		}
	}
	ct.resources = append(ct.resources, resources...)
	if err := waitForReadyStatus(ct.k8sClient, resources, timeout); err != nil {
		return fmt.Errorf("failed to reach READY status: %w", err)
	}
	return nil
}

func (ct *contractTest) ProjectID(ctx context.Context, projectName string) (string, error) {
	project := akov2.AtlasProject{}
	key := types.NamespacedName{Namespace: ct.namespace, Name: projectName}
	if err := ct.k8sClient.Get(ctx, key, &project); err != nil {
		return "", fmt.Errorf("failed to get project ID: %w", err)
	}
	return project.Status.ID, nil
}
