package auditing_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	audit "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer/auditing"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/auditing"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	expectedProjectID = "12345"

	unexpectedProjectID = "54321"
)

var (
	// ErrorNotFound happens when a fake resource is not found
	ErrorNotFound = errors.New("not found")
)

type reconcileOutcome struct {
	result ctrl.Result
	err    error
}

type fakeAuditService struct {
	get func(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingSpec, error)
	set func(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingSpec) error
}

func (fas fakeAuditService) Get(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingSpec, error) {
	if fas.get == nil {
		panic("fake get is unset")
	}
	return fas.get(ctx, projectID)
}

func (fas fakeAuditService) Set(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingSpec) error {
	if fas.set == nil {
		panic("fake set is unset")
	}
	return fas.set(ctx, projectID, auditing)
}

func TestReconcile(t *testing.T) {
	testCases := []struct {
		title    string
		objects  []client.Object
		req      ctrl.Request
		expected reconcileOutcome
	}{
		{
			title: "Wrong object request gets ignored and skipped",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "this-does-not-exist"},
			},
			expected: reconcileOutcome{result: reconcile.Result{}, err: auditing.ErrorNotFound},
		},
		{
			title: "Right auditing object but with skip annotation also gets skipped",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta: metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-auditing",
						Annotations: map[string]string{
							customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
						},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			expected: reconcileOutcome{result: reconcile.Result{}, err: auditing.ErrorSkipped},
		},
		{
			title: "Right auditing object but empty gets rejected by type enum validation",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			expected: reconcileOutcome{result: reconcile.Result{RequeueAfter: workflow.DefaultRetry}, err: validate.ErrorBadEnum},
		},
		{
			title: "Right & proper auditing fails state evaluation if the project id is missing",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
					Spec: v1alpha1.AtlasAuditingSpec{
						Type:       v1alpha1.Standalone,
						ProjectIDs: []string{unexpectedProjectID},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			expected: reconcileOutcome{result: reconcile.Result{RequeueAfter: workflow.DefaultRetry}, err: ErrorNotFound},
		},
	}
	ctx := context.Background()
	as := fakeAuditService{
		get: func(_ context.Context, projectID string) (*v1alpha1.AtlasAuditingSpec, error) {
			if projectID == expectedProjectID {
				return &v1alpha1.AtlasAuditingSpec{}, nil
			}
			return nil, fmt.Errorf("%w project id %s does not exist", ErrorNotFound, projectID)
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			r := newTestReconciler(fakeK8sClient(tc.objects), as)
			result, err := r.Reconcile(ctx, tc.req)
			assert.Equal(t, tc.expected.result, result)
			assert.ErrorIs(t, err, tc.expected.err)
		})
	}
}

func newTestReconciler(k8sClient client.Client, as audit.Service) *auditing.AtlasAuditingReconciler {
	return &auditing.AtlasAuditingReconciler{
		Client:        k8sClient,
		Log:           zap.S(),
		Scheme:        &runtime.Scheme{},
		AuditService:  as,
		EventRecorder: nil,
	}
}

func fakeK8sClient(objects []client.Object) client.Client {
	sch := runtime.NewScheme()
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &v1alpha1.AtlasAuditing{})
	return fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(objects...).
		Build()
}
