package contract

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/control"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"
)

const (
	operatorInstallName = "test-atlas-operator"
	akoTestNamespace    = "ako-test"
)

type ContractTest struct {
	t            *testing.T
	akoInstalled bool
	namespace    string
	resources    []client.Object

	Ctx         context.Context
	AtlasClient *admin.APIClient
	K8sClient   client.Client
}

func (ct *ContractTest) cleanup() {
	for i := len(ct.resources) - 1; i >= 0; i-- {
		resource := ct.resources[i]
		require.NoError(ct.t, ct.K8sClient.Delete(ct.Ctx, resource))
	}
	ct.resources = []client.Object{}
	if ct.namespace != "" {
		require.NoError(ct.t, ct.K8sClient.Delete(ct.Ctx, defaultNamespace(ct.namespace)))
	}
}

func RunContractTest(t *testing.T, name string, ctFn func(ct *ContractTest)) {
	t.Helper()
	if !control.Enabled("AKO_CONTRACT_TEST") {
		t.Skip("Skipping contract test as AKO_CONTRACT_TEST is unset")
		return
	}
	ctx := context.Background()
	ct := &ContractTest{
		t:            t,
		Ctx:          ctx,
		K8sClient:    mustCreateK8sClient(),
		akoInstalled: false,
		resources:    []client.Object{},
		AtlasClient:  mustCreateVersionedAtlasClient(ctx),
	}
	defer ct.cleanup()
	t.Run(name, func(t *testing.T) {
		ctFn(ct)
	})
}

func (ct *ContractTest) AddResources(timeout time.Duration, resources ...client.Object) {
	if !ct.akoInstalled {
		require.NoError(ct.t, ensureTestAtlasOperator(akoTestNamespace))
		require.NoError(ct.t, k8sRecreate(ct.Ctx, ct.K8sClient, globalSecret(akoTestNamespace)))
		ct.akoInstalled = true
	}
	if ct.namespace == "" {
		ct.namespace = utils.RandomName("test-ns")
		require.NoError(ct.t, ct.K8sClient.Create(ct.Ctx, defaultNamespace(ct.namespace)))
	}
	for _, resource := range resources {
		resource.SetNamespace(ct.namespace)
		require.NoError(ct.t, ct.K8sClient.Create(ct.Ctx, resource))
	}
	ct.resources = append(ct.resources, resources...)
	require.NoError(ct.t, waitForReadyStatus(ct.K8sClient, resources, timeout))
}

func (ct *ContractTest) Namespace() string {
	return ct.namespace
}
