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

package operator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
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

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	generatedv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
)

type managerMock struct {
	ManagerProvider
	manager.Manager
	client.FieldIndexer

	client client.Client
	scheme *runtime.Scheme

	opts ctrl.Options
}

func (m *managerMock) GetCache() cache.Cache {
	return &informertest.FakeInformers{}
}

func (m *managerMock) Add(_ manager.Runnable) error {
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

func (m *managerMock) GetEventRecorderFor(_ string) record.EventRecorder {
	return record.NewFakeRecorder(100)
}

func (m *managerMock) GetClient() client.Client {
	return m.client
}

func (m *managerMock) GetFieldIndexer() client.FieldIndexer {
	return &informertest.FakeInformers{}
}

func (m *managerMock) New(_ *rest.Config, options manager.Options) (manager.Manager, error) {
	m.opts = options
	m.scheme = options.Scheme
	m.client = fake.NewClientBuilder().
		WithScheme(options.Scheme).
		Build()

	return m, nil
}

func (m *managerMock) AddHealthzCheck(_ string, _ healthz.Checker) error {
	return nil
}

func (m *managerMock) AddReadyzCheck(_ string, _ healthz.Checker) error {
	return nil
}

func TestBuildManager(t *testing.T) {
	tests := map[string]struct {
		configure                func(b *Builder)
		expectedSyncPeriod       time.Duration
		expectedClusterWideCache bool
		expectedNamespacedCache  bool
		expectedError            error
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
					WithIndependentSyncPeriod(15 * time.Minute).
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
		"should error when independentSyncPeriod is misconfigured": {
			configure: func(b *Builder) {
				b.WithIndependentSyncPeriod(4 * time.Minute)
			},
			expectedSyncPeriod:       DefaultSyncPeriod,
			expectedClusterWideCache: false,
			expectedNamespacedCache:  true,
			expectedError:            errors.New("wrong value for independentSyncPeriod. Value should be greater or equal to 5"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			akoScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(akoScheme))
			require.NoError(t, corev1.AddToScheme(akoScheme))
			require.NoError(t, generatedv1.AddToScheme(akoScheme))

			mgrMock := &managerMock{}
			builder := NewBuilder(mgrMock, akoScheme, 5*time.Minute)
			tt.configure(builder)
			// this is necessary for tests
			builder.WithSkipNameValidation(true)
			_, err := builder.Build(context.Background())
			require.Equal(t, tt.expectedError, err)

			if err == nil {
				assert.Equal(t, tt.expectedSyncPeriod, *mgrMock.opts.Cache.SyncPeriod)
				assert.Equal(t, tt.expectedClusterWideCache, len(mgrMock.opts.Cache.ByObject) > 0)
				assert.Equal(t, tt.expectedNamespacedCache, len(mgrMock.opts.Cache.DefaultNamespaces) > 0)
			}
		})
	}
}
