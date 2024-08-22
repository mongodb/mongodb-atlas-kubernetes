package customresource

import (
	"fmt"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestResourceShouldBeLeftInAtlas(t *testing.T) {
	t.Run("Empty annotations", func(t *testing.T) {
		assert.False(t, IsResourcePolicyKeep(&akov2.AtlasDatabaseUser{}))
	})

	t.Run("Other annotations", func(t *testing.T) {
		assert.False(t, IsResourcePolicyKeep(&akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{"foo": "bar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be removed", func(t *testing.T) {
		assert.False(t, IsResourcePolicyKeep(&akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				// Any other value except for "keep" is considered as "purge"
				Annotations: map[string]string{ResourcePolicyAnnotation: "foobar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be kept", func(t *testing.T) {
		assert.True(t, IsResourcePolicyKeep(&akov2.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{ResourcePolicyAnnotation: ResourcePolicyKeep},
			},
		}))
	})
}

func TestReconciliationShouldBeSkipped(t *testing.T) {
	newResourceTypes := func() []api.AtlasCustomResource {
		return []api.AtlasCustomResource{
			&akov2.AtlasDeployment{},
			&akov2.AtlasDatabaseUser{},
			&akov2.AtlasProject{},
		}
	}

	t.Run("Empty annotations", func(t *testing.T) {
		for _, resourceType := range newResourceTypes() {
			assert.False(t, ReconciliationShouldBeSkipped(resourceType))
		}
	})

	t.Run("Other resource types", func(t *testing.T) {
		for _, resourceType := range newResourceTypes() {
			resourceType.SetAnnotations(map[string]string{"foo": "bar"})
			assert.False(t, ReconciliationShouldBeSkipped(resourceType))
		}
	})

	t.Run("Annotation present, reconciliation should not be skipped", func(t *testing.T) {
		for _, resourceType := range newResourceTypes() {
			resourceType.SetAnnotations(map[string]string{ReconciliationPolicyAnnotation: "foobar"})
			assert.False(t, ReconciliationShouldBeSkipped(resourceType))
		}
	})

	t.Run("Annotation present, reconciliation should be skipped", func(t *testing.T) {
		for _, resourceType := range newResourceTypes() {
			resourceType.SetAnnotations(map[string]string{ReconciliationPolicyAnnotation: ReconciliationPolicySkip})
			assert.True(t, ReconciliationShouldBeSkipped(resourceType))
		}
	})
}

func TestResourceVersionIsValid(t *testing.T) {
	tests := []struct {
		name            string
		resource        api.AtlasCustomResource
		want            bool
		wantErr         assert.ErrorAssertionFunc
		operatorVersion string
	}{
		{
			name: "Resource version is LOWER than operator version",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.3.0",
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.4.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is EQUAL to the operator version",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.3.0",
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.5.0",
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            false,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version with ALLOWED OVERRIDE",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.5.0",
					},
					Annotations: map[string]string{
						ResourceVersionOverride: ResourceVersionAllow,
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version with DISALLOWED OVERRIDE",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.5.0",
					},
					Annotations: map[string]string{
						ResourceVersionOverride: "someValue",
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            false,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is INCORRECT, should return an error",
			resource: &akov2.AtlasProject{
				TypeMeta: metav1.TypeMeta{
					Kind:       "AtlasProject",
					APIVersion: "atlas.mongodb.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "TestProject",
					Labels: map[string]string{
						ResourceVersion: "1.incorrect.semantic.version",
					},
				},
				Spec:   akov2.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            false,
			operatorVersion: "1.3.0",
			wantErr:         assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version.Version = tt.operatorVersion
			got, err := ResourceVersionIsValid(tt.resource)
			if !tt.wantErr(t, err, fmt.Sprintf("ResourceVersionIsValid(%v)", tt.resource)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ResourceVersionIsValid(%v)", tt.resource)
		})
	}
}

func TestComputeSecretWithFallback(t *testing.T) {
	for _, tt := range []struct {
		name                 string
		fallback             bool
		project              *akov2.AtlasProject
		resource             api.ResourceWithCredentials
		expected             *types.NamespacedName
		expectedErrorMessage string
	}{
		{
			name:                 "nil inputs fails with project cannot be nil error without fallback",
			expectedErrorMessage: "resource cannot be nil",
		},

		{
			name: "nil project ignored if resource is set without fallback",
			resource: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Namespace: "local"},
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &corev1.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			expected: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},

		{
			name:                 "nil resource and empty project fails without fallback",
			project:              &akov2.AtlasProject{},
			expectedErrorMessage: "resource cannot be nil",
		},

		{
			name:                 "when both are set empty it fails without fallback",
			project:              &akov2.AtlasProject{},
			resource:             &akov2.AtlasDatabaseUser{},
			expectedErrorMessage: "failed to find credentials secret neither from resource",
		},

		{
			name: "empty resource and proper project get creds from project without fallback",
			project: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name:                    "",
					RegionUsageRestrictions: "",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name:      "project-secret",
						Namespace: "some-namespace",
					},
				},
			},
			resource: &akov2.AtlasDatabaseUser{},
			expected: &client.ObjectKey{
				Name:      "project-secret",
				Namespace: "some-namespace",
			},
		},

		{
			name: "when both are properly set the resource wins without fallback",
			project: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name:                    "",
					RegionUsageRestrictions: "",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name:      "project-secret",
						Namespace: "some-namespace",
					},
				},
			},
			resource: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Namespace: "local"},
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &corev1.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			expected: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},

		{
			name:     "nil inputs renders nil with fallback",
			fallback: true,
		},

		{
			name:     "nil project renders resource secret if set even with fallback",
			fallback: true,
			resource: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Namespace: "local"},
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &corev1.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			expected: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},

		{
			name:     "nil resource and empty project renders nil with fallback",
			fallback: true,
			project:  &akov2.AtlasProject{},
		},

		{
			name:     "when both are set empty it renders nil with fallback",
			fallback: true,
			project:  &akov2.AtlasProject{},
			resource: &akov2.AtlasDatabaseUser{},
		},

		{
			name:     "empty resource and proper project get creds from project even with fallback",
			fallback: true,
			project: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name:                    "",
					RegionUsageRestrictions: "",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name:      "project-secret",
						Namespace: "some-namespace",
					},
				},
			},
			resource: &akov2.AtlasDatabaseUser{},
			expected: &client.ObjectKey{
				Name:      "project-secret",
				Namespace: "some-namespace",
			},
		},

		{
			name:     "when both are properly set the resource wins even with fallback",
			fallback: true,
			project: &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name:                    "",
					RegionUsageRestrictions: "",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name:      "project-secret",
						Namespace: "some-namespace",
					},
				},
			},
			resource: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Namespace: "local"},
				Spec: akov2.AtlasDatabaseUserSpec{
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &corev1.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			expected: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ComputeSecretWithFallback(tt.fallback, tt.project, tt.resource)
			if tt.expectedErrorMessage != "" {
				assert.Nil(t, result, nil)
				assert.ErrorContains(t, err, tt.expectedErrorMessage)
			} else {
				assert.Equal(t, result, tt.expected)
				assert.NoError(t, err)
			}
		})
	}
}
