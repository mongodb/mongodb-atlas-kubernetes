package dryrun

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type mockCluster struct {
	cluster.Cluster
	cache.Cache

	startErr               error
	waitForCacheSyncResult bool
	rec                    record.EventRecorder
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
			wantEvents:             []string{"Normal DryRun dry run finished"},
		},
		{
			name:                   "no errors",
			startErr:               nil,
			waitForCacheSyncResult: true,
			wantEvents:             []string{"Normal DryRun dry run finished"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := record.NewFakeRecorder(100)
			m := NewManager(&mockCluster{
				startErr:               tc.startErr,
				waitForCacheSyncResult: tc.waitForCacheSyncResult,
				rec:                    rec,
			}, zaptest.NewLogger(t))
			gotErr := ""
			if err := m.Start(context.Background()); err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tc.wantErr, gotErr)
			assertEqualEvents(t, tc.wantEvents, rec.Events)
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
			rec := record.NewFakeRecorder(100)
			m := NewManager(&mockCluster{rec: rec}, zaptest.NewLogger(t))
			m.reportError(obj, tc.err)
			assertEqualEvents(t, tc.wantEvents, rec.Events)
		})
	}
}

func assertEqualEvents(t *testing.T, expected []string, actual <-chan string) {
	c := time.After(time.Second)
	for _, e := range expected {
		select {
		case a := <-actual:
			if e != a {
				t.Errorf("Expected event %q, got %q", e, a)
				return
			}
		case <-c:
			t.Errorf("Expected event %q, got nothing", e)
			// continue iterating to print all expected events
		}
	}
	for {
		select {
		case a := <-actual:
			t.Errorf("Unexpected event: %q", a)
		default:
			return // No more events, as expected.
		}
	}
}
