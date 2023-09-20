package customresource

import (
	"fmt"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestResourceShouldBeLeftInAtlas(t *testing.T) {
	t.Run("Empty annotations", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{}))
	})

	t.Run("Other annotations", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{"foo": "bar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be removed", func(t *testing.T) {
		assert.False(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				// Any other value except for "keep" is considered as "purge"
				Annotations: map[string]string{ResourcePolicyAnnotation: "foobar"},
			},
		}))
	})

	t.Run("Annotation present, resources should be kept", func(t *testing.T) {
		assert.True(t, ResourceShouldBeLeftInAtlas(&v1.AtlasDatabaseUser{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{ResourcePolicyAnnotation: ResourcePolicyKeep},
			},
		}))
	})
}

func TestReconciliationShouldBeSkipped(t *testing.T) {
	newResourceTypes := func() []v1.AtlasCustomResource {
		return []v1.AtlasCustomResource{
			&v1.AtlasDeployment{},
			&v1.AtlasDatabaseUser{},
			&v1.AtlasProject{},
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
		resource        v1.AtlasCustomResource
		want            bool
		wantErr         assert.ErrorAssertionFunc
		operatorVersion string
	}{
		{
			name: "Resource version is LOWER than operator version",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.4.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is EQUAL to the operator version",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            false,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version with ALLOWED OVERRIDE",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            true,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is GREATER than the operator version with DISALLOWED OVERRIDE",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
				Status: status.AtlasProjectStatus{},
			},
			want:            false,
			operatorVersion: "1.3.0",
			wantErr:         assert.NoError,
		},
		{
			name: "Resource version is INCORRECT, should return an error",
			resource: &v1.AtlasProject{
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
				Spec:   v1.AtlasProjectSpec{},
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

func TestIsGov(t *testing.T) {
	t.Run("should return false for invalid domain", func(t *testing.T) {
		assert.False(t, IsGov("http://x:namedport"))
	})

	t.Run("should return false for commercial Atlas domain", func(t *testing.T) {
		assert.False(t, IsGov("https://cloud.mongodb.com/"))
	})

	t.Run("should return true for Atlas for government domain", func(t *testing.T) {
		assert.True(t, IsGov("https://cloud.mongodbgov.com/"))
	})
}

func TestIsSupportedByCloudGov(t *testing.T) {
	dataProvider := map[string]struct {
		domain      string
		resource    v1.AtlasCustomResource
		expectation bool
	}{
		"should return true when it's commercial Atlas": {
			domain:      "https://cloud.mongodb.com",
			resource:    &v1.AtlasDataFederation{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Project": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasProject{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Team": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasTeam{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupSchedule": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasBackupSchedule{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupPolicy": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is DatabaseUser": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is regular Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &v1.AtlasDeployment{
				Spec: v1.AtlasDeploymentSpec{},
			},
			expectation: true,
		},
		"should return false when it's Atlas Gov and resource is DataFederation": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &v1.AtlasDataFederation{},
			expectation: false,
		},
		"should return false when it's Atlas Gov and resource is Serverless Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &v1.AtlasDeployment{
				Spec: v1.AtlasDeploymentSpec{
					ServerlessSpec: &v1.ServerlessSpec{},
				},
			},
			expectation: false,
		},
	}

	for desc, data := range dataProvider {
		t.Run(desc, func(t *testing.T) {
			assert.Equal(t, data.expectation, IsResourceSupportedInDomain(data.resource, data.domain))
		})
	}
}
