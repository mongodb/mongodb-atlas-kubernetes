package auditctl_test

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

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer/audit"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/auditctl"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
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
	get func(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingConfig, error)
	set func(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingConfig) error
}

func (fas fakeAuditService) Get(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingConfig, error) {
	if fas.get == nil {
		panic("fake get is unset")
	}
	return fas.get(ctx, projectID)
}

func (fas fakeAuditService) Set(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingConfig) error {
	if fas.set == nil {
		panic("fake set is unset")
	}
	return fas.set(ctx, projectID, auditing)
}

func TestReconcile(t *testing.T) {
	cfgInAtlas := v1alpha1.AtlasAuditingConfig{
		Enabled:                   true,
		AuditAuthorizationSuccess: true,
		AuditFilter:               &apiextensionsv1.JSON{Raw: ([]byte)("{}")},
	}
	atlasUpdated := false
	as := fakeAuditService{
		get: func(_ context.Context, projectID string) (*v1alpha1.AtlasAuditingConfig, error) {
			if projectID == expectedProjectID {
				return &cfgInAtlas, nil
			}
			return nil, fmt.Errorf("%w project id %s does not exist", ErrorNotFound, projectID)
		},
		set: func(_ context.Context, projectID string, _ *v1alpha1.AtlasAuditingConfig) error {
			if projectID == expectedProjectID {
				atlasUpdated = true
				return nil
			}
			return fmt.Errorf("setting project id %s is not supported by this fake", projectID)
		},
	}

	testCases := []struct {
		title        string
		objects      []client.Object
		req          ctrl.Request
		atlasUpdated bool
		expected     reconcileOutcome
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
			expected: reconcileOutcome{result: reconcile.Result{}, err: auditctl.ErrorNotFound},
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
			expected: reconcileOutcome{result: ctrl.Result{}, err: auditctl.ErrorSkipped},
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
			expected: reconcileOutcome{result: ctrl.Result{}, err: validate.ErrorBadEnum},
		},
		{
			title: "Right & proper auditing fails reconciling when project is missing",
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
			expected: reconcileOutcome{result: ctrl.Result{}, err: ErrorNotFound},
		},
		{
			title: "Right & proper auditing reconciles ok when project id is correct",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
					Spec: v1alpha1.AtlasAuditingSpec{
						Type:       v1alpha1.Standalone,
						ProjectIDs: []string{expectedProjectID},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			atlasUpdated: true,
			expected:     reconcileOutcome{result: ctrl.Result{}, err: nil},
		},
		{
			title: "Right & proper auditing reconciles ok when project id is correct",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
					Spec: v1alpha1.AtlasAuditingSpec{
						Type:       v1alpha1.Standalone,
						ProjectIDs: []string{expectedProjectID},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			atlasUpdated: true,
			expected:     reconcileOutcome{result: ctrl.Result{}, err: nil},
		},
		{
			title: "Right & proper auditing reconciles ok with correct project (idle case)",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
					Spec: v1alpha1.AtlasAuditingSpec{
						Type:       v1alpha1.Standalone,
						ProjectIDs: []string{expectedProjectID},
						AtlasAuditingConfig: v1alpha1.AtlasAuditingConfig{
							Enabled:                   cfgInAtlas.Enabled,
							AuditAuthorizationSuccess: cfgInAtlas.AuditAuthorizationSuccess,
							AuditFilter:               cfgInAtlas.AuditFilter,
						},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			atlasUpdated: false,
			expected:     reconcileOutcome{result: ctrl.Result{}, err: nil},
		},
		{
			title: "Right & proper auditing gets reconciling without projects",
			objects: []client.Object{
				&v1alpha1.AtlasAuditing{
					TypeMeta:   metav1.TypeMeta{Kind: "AtlasAuditing", APIVersion: "v1alpha1"},
					ObjectMeta: metav1.ObjectMeta{Name: "test-auditing"},
					Spec: v1alpha1.AtlasAuditingSpec{
						Type:       v1alpha1.Standalone,
						ProjectIDs: []string{},
					},
				},
			},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Name: "test-auditing"},
			},
			atlasUpdated: false,
			expected:     reconcileOutcome{result: ctrl.Result{}, err: nil},
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			atlasUpdated = false
			r := newTestReconciler(fakeK8sClient(tc.objects), as)
			result, err := r.Reconcile(ctx, tc.req)
			assert.Equal(t, tc.expected.result, result)
			assert.ErrorIs(t, err, tc.expected.err)
			assert.Equal(t, tc.atlasUpdated, atlasUpdated)
		})
	}
}

func newTestReconciler(k8sClient client.Client, as audit.Service) *auditctl.AtlasAuditingReconciler {
	return &auditctl.AtlasAuditingReconciler{
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
