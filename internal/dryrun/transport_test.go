package dryrun

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/record"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestDryRunTransport(t *testing.T) {
	for _, tc := range []struct {
		name       string
		req        *http.Request
		obj        runtime.Object
		wantEvents []string
		wantErr    string
	}{
		{
			name: "GET request",
			req: &http.Request{
				Method: http.MethodGet,
			},
		},
		{
			name: "no object",
			req: &http.Request{
				Method: http.MethodPost,
			},
			wantErr: "no object present in context: cannot record dry run without an object",
		},
		{
			name: "unknown verb",
			req: &http.Request{
				Method: "UNKNOWN",
				URL:    &url.URL{Path: "/test"},
			},
			obj: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
			},
			wantErr: "DryRun event GVK=atlas.mongodb.com/v1, Kind=AtlasProject, Namespace=bar, Name=foo, EventType=Normal, Reason=DryRun, Message=Would execute UNKNOWN /test",
			wantEvents: []string{
				"Normal DryRun Would execute UNKNOWN /test",
			},
		},
		{
			name: "POST request",
			req: &http.Request{
				Method: http.MethodPost,
				URL:    &url.URL{Path: "/test"},
			},
			obj:     &akov2.AtlasProject{},
			wantErr: "DryRun event GVK=atlas.mongodb.com/v1, Kind=AtlasProject, Namespace=, Name=, EventType=Normal, Reason=DryRun, Message=Would create (POST) /test",
			wantEvents: []string{
				"Normal DryRun Would create (POST) /test",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := record.NewFakeRecorder(100)
			nopDelegate := RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				return nil, nil
			})
			transport := NewDryRunTransport(rec, nopDelegate)

			if tc.obj != nil {
				scheme := runtime.NewScheme()
				assert.NoError(t, akov2.AddToScheme(scheme))
				gvks, _, err := scheme.ObjectKinds(tc.obj)
				require.NoError(t, err)
				meta := tc.obj.(schema.ObjectKind)
				meta.SetGroupVersionKind(gvks[0])
			}

			req := tc.req.WithContext(WithRuntimeObject(context.Background(), tc.obj))
			_, err := transport.RoundTrip(req)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			require.Equal(t, tc.wantErr, gotErr)
			assertEqualEvents(t, tc.wantEvents, rec.Events)
		})
	}
}

func assertEqualEvents(t *testing.T, expected []string, actual <-chan string) {
	c := time.After(wait.ForeverTestTimeout)
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
