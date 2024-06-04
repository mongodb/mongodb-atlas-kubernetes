package operator

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type managerMock struct {
	ManagerProvider
	manager.Manager
	client.FieldIndexer

	client client.Client
	scheme *runtime.Scheme

	gotHealthzCheck string
	gotReadyzCheck  string

	opts ctrl.Options
}

func (m *managerMock) GetCache() cache.Cache {
	return &informertest.FakeInformers{}
}

func (m *managerMock) Add(runnable manager.Runnable) error {
	return nil
}

func (m *managerMock) GetLogger() logr.Logger {
	return logr.Logger{}
}

func (m *managerMock) GetControllerOptions() config.Controller {
	return config.Controller{}
}

func (m *managerMock) GetScheme() *runtime.Scheme {
	return m.scheme
}

func (m *managerMock) GetEventRecorderFor(name string) record.EventRecorder {
	return record.NewFakeRecorder(100)
}

func (m *managerMock) GetClient() client.Client {
	return m.client
}

func (m *managerMock) GetFieldIndexer() client.FieldIndexer {
	return &informertest.FakeInformers{}
}

func (m *managerMock) New(config *rest.Config, options manager.Options) (manager.Manager, error) {
	m.opts = options
	m.scheme = options.Scheme
	m.client = fake.NewClientBuilder().
		WithScheme(options.Scheme).
		Build()

	return m, nil
}

func (m *managerMock) AddHealthzCheck(name string, check healthz.Checker) error {
	m.gotHealthzCheck = name
	return nil
}

func (m *managerMock) AddReadyzCheck(name string, check healthz.Checker) error {
	m.gotReadyzCheck = name
	return nil
}

func TestBuildManager(t *testing.T) {
	tests := map[string]struct {
		configure                func(b *Builder)
		expectedSyncPeriod       time.Duration
		expectedClusterWideCache bool
		expectedNamespacedCache  bool
	}{
		"should build the manager with default values": {
			configure:                func(b *Builder) {},
			expectedSyncPeriod:       DefaultSyncPeriod,
			expectedClusterWideCache: true,
			expectedNamespacedCache:  false,
		},
		"should build the manager with namespace config": {
			configure: func(b *Builder) {
				b.WithNamespaces("ns1", "ns2")
			},
			expectedSyncPeriod:       DefaultSyncPeriod,
			expectedClusterWideCache: false,
			expectedNamespacedCache:  true,
		},
		"should build the manager with custom config": {
			configure: func(b *Builder) {
				b.WithConfig(&rest.Config{}).
					WithNamespaces("ns1").
					WithLogger(zaptest.NewLogger(t)).
					WithSyncPeriod(time.Hour).
					WithMetricAddress(":9090").
					WithProbeAddress(":9091").
					WithLeaderElection(true).
					WithAtlasDomain("https://cloud-qa.mongodb.com").
					WithPredicates([]predicate.Predicate{predicate.GenerationChangedPredicate{}}).
					WithAPISecret(client.ObjectKey{Namespace: "ns1", Name: "creds"}).
					WithAtlasProvider(&atlas.TestProvider{}).
					WithFeatureFlags(featureflags.NewFeatureFlags(func() []string { return nil })).
					WithDeletionProtection(true)
			},
			expectedSyncPeriod:       time.Hour,
			expectedClusterWideCache: false,
			expectedNamespacedCache:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			akoScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(akoScheme))

			mgrMock := &managerMock{}
			builder := NewBuilder(mgrMock, akoScheme)
			tt.configure(builder)
			_, err := builder.Build(context.Background())
			require.NoError(t, err)

			assert.Equal(t, tt.expectedSyncPeriod, *mgrMock.opts.Cache.SyncPeriod)
			assert.Equal(t, tt.expectedClusterWideCache, len(mgrMock.opts.Cache.ByObject) > 0)
			assert.Equal(t, tt.expectedNamespacedCache, len(mgrMock.opts.Cache.DefaultNamespaces) > 0)
		})
	}
}
