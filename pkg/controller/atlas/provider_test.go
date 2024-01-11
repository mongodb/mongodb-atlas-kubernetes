package atlas

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"
)

func TestProvider_Client(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-secret",
			Namespace: "default",
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"orgId":         []byte("1234567890"),
			"publicApiKey":  []byte("a1b2c3"),
			"privateApiKey": []byte("abcdef123456"),
		},
		Type: "Opaque",
	}

	sch := runtime.NewScheme()
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{})
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secret).
		Build()

	t.Run("should return Atlas API client and organization id using global secret", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodb.com/", client.ObjectKey{Name: "api-secret", Namespace: "default"}, k8sClient)

		c, id, err := p.Client(context.Background(), nil, zaptest.NewLogger(t).Sugar())
		assert.NoError(t, err)
		assert.Equal(t, "1234567890", id)
		assert.NotNil(t, c)
	})

	t.Run("should return Atlas API client and organization id using connection secret", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodb.com/", client.ObjectKey{Name: "global-secret", Namespace: "default"}, k8sClient)

		c, id, err := p.Client(context.Background(), &client.ObjectKey{Name: "api-secret", Namespace: "default"}, zaptest.NewLogger(t).Sugar())
		assert.NoError(t, err)
		assert.Equal(t, "1234567890", id)
		assert.NotNil(t, c)
	})
}

func TestProvider_IsCloudGov(t *testing.T) {
	t.Run("should return false for invalid domain", func(t *testing.T) {
		p := NewProductionProvider("http://x:namedport", client.ObjectKey{}, nil)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return false for commercial Atlas domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodb.com/", client.ObjectKey{}, nil)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return true for Atlas for government domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodbgov.com/", client.ObjectKey{}, nil)
		assert.True(t, p.IsCloudGov())
	})
}

func TestProvider_IsResourceSupported(t *testing.T) {
	dataProvider := map[string]struct {
		domain      string
		resource    akov2.AtlasCustomResource
		expectation bool
	}{
		"should return true when it's commercial Atlas": {
			domain:      "https://cloud.mongodb.com",
			resource:    &akov2.AtlasDataFederation{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Project": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasProject{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Team": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasTeam{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupSchedule": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupSchedule{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupPolicy": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is DatabaseUser": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is regular Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{},
			},
			expectation: true,
		},
		"should return false when it's Atlas Gov and resource is DataFederation": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasDataFederation{},
			expectation: false,
		},
		"should return false when it's Atlas Gov and resource is Serverless Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{},
				},
			},
			expectation: false,
		},
	}

	for desc, data := range dataProvider {
		t.Run(desc, func(t *testing.T) {
			p := NewProductionProvider(data.domain, client.ObjectKey{}, nil)
			assert.Equal(t, data.expectation, p.IsResourceSupported(data.resource))
		})
	}
}

func TestValidateSecretData(t *testing.T) {
	t.Run("should be invalid and all missing data", func(t *testing.T) {
		missing, ok := validateSecretData(&credentialsSecret{})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"orgId", "publicApiKey", "privateApiKey"})
	})

	t.Run("should be invalid and organization id is missing", func(t *testing.T) {
		missing, ok := validateSecretData(&credentialsSecret{PublicKey: "abcdef", PrivateKey: "123456"})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"orgId"})
	})

	t.Run("should be invalid and public key id is missing", func(t *testing.T) {
		missing, ok := validateSecretData(&credentialsSecret{OrgID: "abcdef", PrivateKey: "123456"})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"publicApiKey"})
	})

	t.Run("should be invalid and private key id is missing", func(t *testing.T) {
		missing, ok := validateSecretData(&credentialsSecret{PublicKey: "abcdef", OrgID: "123456"})
		assert.False(t, ok)
		assert.Equal(t, missing, []string{"privateApiKey"})
	})

	t.Run("should be valid", func(t *testing.T) {
		missing, ok := validateSecretData(&credentialsSecret{OrgID: "my-org", PublicKey: "abcdef", PrivateKey: "123456"})
		assert.True(t, ok)
		assert.Empty(t, missing)
	})
}

func TestOperatorUserAgent(t *testing.T) {
	userAgent := operatorUserAgent()

	require.Contains(t, userAgent, "MongoDBAtlasKubernetesOperator")
	require.Contains(t, userAgent, version.Version)
}
