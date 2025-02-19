package dryrun

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	crFake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

type mockCluster struct {
	cluster.Cluster
	cache.Cache

	startErr               error
	waitForCacheSyncResult bool
	rec                    record.EventRecorder
	client                 client.Client
}

func (c *mockCluster) GetScheme() *runtime.Scheme {
	akoScheme := runtime.NewScheme()
	utilruntime.Must(scheme.AddToScheme(akoScheme))
	utilruntime.Must(akov2.AddToScheme(akoScheme))
	return akoScheme
}

func (c *mockCluster) GetHTTPClient() *http.Client {
	return http.DefaultClient
}
func (c *mockCluster) GetConfig() *rest.Config {
	return &rest.Config{}
}
func (c *mockCluster) GetEventRecorderFor(name string) record.EventRecorder {
	return c.rec
}

func (c *mockCluster) GetCache() cache.Cache {
	return c
}

func (c *mockCluster) WaitForCacheSync(context.Context) bool {
	return c.waitForCacheSyncResult
}

func (c *mockCluster) Start(ctx context.Context) error {
	<-ctx.Done() // block until context is canceled
	return c.startErr
}

func (c *mockCluster) GetClient() client.Client {
	return c.client
}

func TestManagerStart(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		startErr               error
		waitForCacheSyncResult bool
		wantErr                string
		wantEvents             []string
	}{
		{
			name:                   "no start error but cache sync failed",
			startErr:               nil,
			waitForCacheSyncResult: false,
			wantErr:                "cluster cache sync failed",
		},
		{
			name:                   "cache sync error is preferred over start error",
			startErr:               errors.New("start error"),
			waitForCacheSyncResult: false,
			wantErr:                "cluster cache sync failed",
		},
		{
			name:                   "start error",
			startErr:               errors.New("start error"),
			waitForCacheSyncResult: true,
			wantErr:                "cluster start failed: start error",
			wantEvents:             []string{"finished"},
		},
		{
			name:                   "no errors",
			startErr:               nil,
			waitForCacheSyncResult: true,
			wantEvents:             []string{"finished"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mckCluster := mockCluster{
				startErr:               tc.startErr,
				waitForCacheSyncResult: tc.waitForCacheSyncResult,
			}
			eventsGetter := fake.NewClientset().CoreV1()
			m, err := NewManager(&mckCluster, eventsGetter, zaptest.NewLogger(t))
			require.NoError(t, err)

			gotErr := ""
			if err := m.Start(context.Background()); err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tc.wantErr, gotErr)

			if len(tc.wantEvents) > 0 {
				ev, err := eventsGetter.Events("default").List(context.Background(), metav1.ListOptions{})
				require.NoError(t, err)
				require.NotEmpty(t, ev.Items)

				var gotEvMsgs []string
				for i := range ev.Items {
					gotEvMsgs = append(gotEvMsgs, ev.Items[i].Message)
				}
				require.Equal(t, tc.wantEvents, gotEvMsgs)
			}
		})
	}
}

func TestDryRunReportError(t *testing.T) {
	obj := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
		},
	}

	for _, tc := range []struct {
		name       string
		err        error
		wantEvents []string
	}{
		{
			name: "no error",
			err:  nil,
		},
		{
			name:       "error",
			err:        errors.New("some random error"),
			wantEvents: []string{"Warning DryRun some random error"},
		},
		{
			name: "nested errors",
			err: fmt.Errorf("%w, %w",
				errors.New("some random error"),
				errors.Join(
					errors.New("another random error"),
					errors.New("yet another random error"),
				),
			),
			wantEvents: []string{"Warning DryRun some random error, another random error\nyet another random error"},
		},
		{
			name:       "dry run error",
			err:        &DryRunError{Msg: "dry run error"},
			wantEvents: []string{"Normal DryRun dry run error"},
		},
		{
			name: "multiple nested dry run errors",
			err: fmt.Errorf("%w", fmt.Errorf("%w, %w",
				&DryRunError{Msg: "dry run error 1"},
				&DryRunError{Msg: "dry run error 2"},
			)),
			wantEvents: []string{
				"Normal DryRun dry run error 1",
				"Normal DryRun dry run error 2",
			},
		},
		{
			name: "multiple nested dry run errors in errors.Join",
			err: fmt.Errorf("%w, %w", nil, errors.Join(
				&DryRunError{Msg: "dry run error 1"},
				nil,
				&DryRunError{Msg: "dry run error 2"},
				fmt.Errorf("%w, %w, %w, %w",
					&DryRunError{Msg: "dry run error 3"},
					nil,
					errors.Join(nil, fmt.Errorf("%w", &DryRunError{Msg: "dry run error 4"}), nil),
					&DryRunError{Msg: "dry run error 5"},
				),
			)),
			wantEvents: []string{
				"Normal DryRun dry run error 1",
				"Normal DryRun dry run error 2",
				"Normal DryRun dry run error 3",
				"Normal DryRun dry run error 4",
				"Normal DryRun dry run error 5",
			},
		},
		{
			name: "forgot to wrap dry run error",
			err: fmt.Errorf("%w",
				fmt.Errorf("%w %w",
					//nolint:errorlint
					fmt.Errorf("errors occurred: %v, %v, %v", &DryRunError{Msg: "dry run error 1"}, &DryRunError{Msg: "dry run error 2"}, errors.Join(&DryRunError{Msg: "dry run error 3"}, &DryRunError{Msg: "dry run error 4"})),
					fmt.Errorf("errors occurred: %w, %w", &DryRunError{Msg: "dry run error 5"}, &DryRunError{Msg: "dry run error 6"}),
				),
			),
			wantEvents: []string{
				"Normal DryRun errors occurred: DryRun event: dry run error 1, DryRun event: dry run error 2, DryRun event: dry run error 3\nDryRun event: dry run error 4",
				"Normal DryRun dry run error 5",
				"Normal DryRun dry run error 6",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			eventsGetter := fake.NewClientset().CoreV1()
			m, err := NewManager(&mockCluster{}, eventsGetter, zaptest.NewLogger(t))
			require.NoError(t, err)
			m.reportError(context.Background(), obj, tc.err)

			if len(tc.wantEvents) > 0 {
				ev, err := eventsGetter.Events(obj.Namespace).List(context.Background(), metav1.ListOptions{})
				require.NoError(t, err)
				require.NotEmpty(t, ev.Items)

				var gotEvMsgs []string
				for i := range ev.Items {
					gotEvMsgs = append(gotEvMsgs, fmt.Sprintf("%s %s %s", ev.Items[i].Type, ev.Items[i].Reason, ev.Items[i].Message))
				}
				require.Equal(t, tc.wantEvents, gotEvMsgs)
			}
		})
	}
}

type mockReconciler struct {
	reconcile.Reconciler
	ErrFail  error
	Resource akov2.AtlasCustomResource
}

func (m *mockReconciler) Reconcile(_ context.Context, _ ctrl.Request) (ctrl.Result, error) {
	return ctrl.Result{}, m.ErrFail
}

func (m *mockReconciler) For() (client.Object, builder.Predicates) {
	return m.Resource.(client.Object), builder.Predicates{}
}

func TestManager_dryRunReconcilers(t *testing.T) {
	tests := []struct {
		name        string
		reconcilers []reconciler
		logger      *zap.Logger
		instanceUID string
		ctx         context.Context
		wantErr     bool
	}{
		{
			name: "Should run dry run without errors for AtlasProject resource",
			reconcilers: []reconciler{
				&mockReconciler{
					Resource: &akov2.AtlasProject{},
					ErrFail:  nil,
				},
			},
			wantErr: false,
			ctx:     context.Background(),
			logger:  zap.L(),
		},
		{
			name: "Should not fail when a reconciler fails",
			reconcilers: []reconciler{
				&mockReconciler{
					Resource: &akov2.AtlasProject{},
					ErrFail:  fmt.Errorf("failed to reconcile"),
				},
			},
			wantErr: false,
			ctx:     context.Background(),
			logger:  zap.L(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventsGetter := fake.NewClientset().CoreV1()
			schm := scheme.Scheme
			require.NoError(t, akov2.AddToScheme(schm))

			clstr := &mockCluster{
				startErr:               nil,
				waitForCacheSyncResult: true,
				client: crFake.NewClientBuilder().WithScheme(schm).WithObjects(
					&akov2.AtlasProject{
						TypeMeta: metav1.TypeMeta{
							Kind:       "AtlasProject",
							APIVersion: "atlas.mongodb.com/v1",
						},
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test",
							Namespace: "test",
						},
						Spec:   akov2.AtlasProjectSpec{},
						Status: status.AtlasProjectStatus{},
					}).Build(),
			}
			m := &Manager{
				Cluster:      clstr,
				reconcilers:  tt.reconcilers,
				logger:       tt.logger,
				instanceUID:  tt.instanceUID,
				eventsClient: eventsGetter,
			}
			if err := m.dryRunReconcilers(tt.ctx); (err != nil) != tt.wantErr {
				t.Errorf("dryRunReconcilers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
