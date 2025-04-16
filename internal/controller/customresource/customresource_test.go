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

package customresource

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
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

func TestComputeSecret(t *testing.T) {
	for _, tt := range []struct {
		name         string
		project      *akov2.AtlasProject
		resource     api.ObjectWithCredentials
		wantRef      *types.NamespacedName
		wantErrorMsg string
	}{
		{
			name:         "nil inputs fails with resource cannot be nil",
			wantErrorMsg: "resource cannot be nil",
		},

		{
			name: "nil project ignored if resource is set",
			resource: &akov2.AtlasDatabaseUser{
				ObjectMeta: metav1.ObjectMeta{Namespace: "local"},
				Spec: akov2.AtlasDatabaseUserSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			wantRef: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},

		{
			name:         "nil resource and empty project fails",
			project:      &akov2.AtlasProject{},
			wantErrorMsg: "resource cannot be nil",
		},

		{
			name:     "when both are set empty it renders nil",
			project:  &akov2.AtlasProject{},
			resource: &akov2.AtlasDatabaseUser{},
		},

		{
			name: "empty resource and proper project get creds from project",
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
			wantRef: &client.ObjectKey{
				Name:      "project-secret",
				Namespace: "some-namespace",
			},
		},

		{
			name: "when both are properly set the resource wins",
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
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{Name: "local-secret"},
					},
				},
			},
			wantRef: &client.ObjectKey{
				Name:      "local-secret",
				Namespace: "local",
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ComputeSecret(tt.project, tt.resource)
			if tt.wantErrorMsg != "" {
				assert.Nil(t, result, nil)
				assert.ErrorContains(t, err, tt.wantErrorMsg)
			} else {
				assert.Equal(t, result, tt.wantRef)
				assert.NoError(t, err)
			}
		})
	}
}
