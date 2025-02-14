package reconciler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

// projectReferrer allows dual ref generic tests without supporting CRDs
type projectReferrer struct {
	akov2.AtlasProject
	pdr akov2.ProjectDualReference
}

func (pr *projectReferrer) ProjectDualRef() *akov2.ProjectDualReference {
	return &pr.pdr
}

func TestResolveCredentials(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		title         string
		objects       []client.Object
		input         project.ProjectReferrerObject
		expected      *client.ObjectKey
		expectedError error
	}{
		{
			title: "all empty returns nil",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{},
				pdr:          akov2.ProjectDualReference{},
			},
		},
		{
			title: "credential set with rest empty returns credential",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{},
				pdr: akov2.ProjectDualReference{
					ConnectionSecret: &api.LocalObjectReference{Name: "some-secret"},
				},
			},
			expected: &client.ObjectKey{Name: "some-secret"},
		},
		{
			title: "credential set with rest empty but for the namespace returns credential",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{ObjectMeta: v1.ObjectMeta{Namespace: "ns"}},
				pdr: akov2.ProjectDualReference{
					ConnectionSecret: &api.LocalObjectReference{Name: "local-secret"},
				},
			},
			expected: &client.ObjectKey{Name: "local-secret", Namespace: "ns"},
		},
		{
			title: "credential unset gets credential from project",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{Name: "project", Namespace: "project-ns"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "project-secret",
							Namespace: "project-ns",
						},
					},
				},
				pdr: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "project", Namespace: "project-ns"},
				},
			},
			expected: &client.ObjectKey{Name: "project-secret", Namespace: "project-ns"},
		},
		{
			title: "when project namespace is unset, credential unset remnder credential project from same namespace",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{Name: "project", Namespace: "project-ns"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "project-secret",
							Namespace: "project-ns",
						},
					},
				},
				pdr: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "project"},
				},
			},
			expected: &client.ObjectKey{Name: "project-secret", Namespace: "project-ns"},
		},
		{
			title: "credential unset gets credential from project",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{Name: "project", Namespace: "project-ns"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "project-secret",
							Namespace: "project-ns",
						},
					},
				},
				pdr: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "project", Namespace: "project-ns"},
				},
			},
			expected: &client.ObjectKey{Name: "project-secret", Namespace: "project-ns"},
		},
		{
			title: "credential unset with non matching project fails",
			input: &projectReferrer{
				AtlasProject: akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{Name: "other-project", Namespace: "project-ns"},
				},
				pdr: akov2.ProjectDualReference{
					ProjectRef: &common.ResourceRefNamespaced{Name: "project", Namespace: "project-ns"},
				},
			},
			expectedError: reconciler.ErrMissingKubeProject,
		},
		{
			title: "should select CustomRoleCredentials from CustomRole when externalProjectID is set",
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "testProjectID"},
						ConnectionSecret:   &api.LocalObjectReference{Name: "credentials"},
					},
				},
			},
			expected: &client.ObjectKey{Name: "credentials", Namespace: "testNamespace"},
		},
		{
			title: "should NOT select CustomRoleCredentials from CustomRole when externalProjectID is set with empty local credentials",
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{ID: "testProjectID"},
						ConnectionSecret:   nil,
					},
				},
			},
			expected: nil,
		},
		{
			title: "should select CustomRoleCredentials when projectRef is set with LocalCredentials",
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "testCredentials",
							Namespace: "testNamespace",
						},
					},
				},
			},
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef:       &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
						ConnectionSecret: &api.LocalObjectReference{Name: "credentials"},
					},
				},
			},
			expected: &client.ObjectKey{Name: "credentials", Namespace: "testNamespace"},
		},
		{
			title: "should select AtlasProject credentials when projectRef is set but LocalCredentials are empty",
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{
							Name:      "testCredentials",
							Namespace: "testNamespace",
						},
					},
				},
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "testCredentials",
						Namespace: "testNamespace",
					},
				},
			},
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
					},
				},
			},
			expected: &client.ObjectKey{Name: "testCredentials", Namespace: "testNamespace"},
		},
		{
			title: "should NOT select AtlasProject credentials when projectRef is set but AtlasProject credentials are not configured",
			objects: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: v1.ObjectMeta{
						Name:      "testProject",
						Namespace: "testNamespace",
					},
					Spec:   akov2.AtlasProjectSpec{},
					Status: status.AtlasProjectStatus{},
				},
				&corev1.Secret{
					ObjectMeta: v1.ObjectMeta{
						Name:      "testCredentials",
						Namespace: "testNamespace",
					},
				},
			},
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
					},
				},
			},
		}, {
			title:   "should NOT select AtlasProject credentials when projectRef is set and AtlasProject doesn't exist",
			objects: []client.Object{},
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec: akov2.AtlasCustomRoleSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{Name: "testProject", Namespace: "testNamespace"},
					},
				},
			},
			expectedError: reconciler.ErrMissingKubeProject,
		},
		{
			title:   "should NOT select credentials when both projectRef and externalProjectId are empty",
			objects: []client.Object{},
			input: &akov2.AtlasCustomRole{
				ObjectMeta: v1.ObjectMeta{
					Name:      "testRole",
					Namespace: "testNamespace",
				},
				Spec:   akov2.AtlasCustomRoleSpec{},
				Status: status.AtlasCustomRoleStatus{},
			},
			expected: nil,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fakeClient := newFakeKubeClient(t, mergeSchemaObject(tc.objects, tc.input)...)
			r := reconciler.AtlasReconciler{
				Client: fakeClient,
			}
			credential, err := r.ResolveCredentials(ctx, tc.input)
			if tc.expectedError != nil {
				require.Nil(t, credential)
				assert.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, credential)
			}
		})
	}
}

func newFakeKubeClient(t *testing.T, objs ...client.Object) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, akov2.AddToScheme(scheme))
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func mergeSchemaObject(objs []client.Object, pro project.ProjectReferrerObject) []client.Object {
	pr, ok := pro.(*projectReferrer)
	if ok {
		return append(objs, &pr.AtlasProject)
	}
	return append(objs, pro)
}
