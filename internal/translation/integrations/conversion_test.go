package integrations

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRoundTrip_Integrations(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &Integration{}
		f.Fuzz(fuzzed)
		fuzzed, err := NewIntegration(&fuzzed.Integration)
		require.NoError(t, err)

		// Don't fuzz secrets that we haven't created in the fake client
		fuzzed.LicenseKeyRef = common.ResourceRefNamespaced{}
		fuzzed.WriteTokenRef = common.ResourceRefNamespaced{}
		fuzzed.ReadTokenRef = common.ResourceRefNamespaced{}
		fuzzed.APIKeyRef = common.ResourceRefNamespaced{}
		fuzzed.ServiceKeyRef = common.ResourceRefNamespaced{}
		fuzzed.APITokenRef = common.ResourceRefNamespaced{}
		fuzzed.RoutingKeyRef = common.ResourceRefNamespaced{}
		fuzzed.SecretRef = common.ResourceRefNamespaced{}
		fuzzed.PasswordRef = common.ResourceRefNamespaced{}

		// Don't expect the 'dud' fields to be converted
		fuzzed.FlowName = ""
		fuzzed.OrgName = ""
		fuzzed.Name = ""
		fuzzed.Scheme = ""

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		assert.NoError(t, corev1.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		toAtlasResult, err := toAtlas(*fuzzed, context.Background(), k8sClient, "testNS")
		require.NoError(t, err)

		fromAtlasResult, err := fromAtlas(toAtlasResult)
		require.NoError(t, err)

		equals := cmp.Diff(fuzzed, fromAtlasResult) == ""
		require.True(t, equals)
	}
}

func TestIntegrationsReadPassword(t *testing.T) {
	in := &Integration{}

	in.LicenseKeyRef = common.ResourceRefNamespaced{
		Name:      "secret-name",
		Namespace: "secret-namespace",
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret-name",
			Namespace: "secret-namespace",
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"password": []byte("Passw0rd!"),
		},
	}

	testScheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(testScheme))
	assert.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(secret).
		Build()

	toAtlasResult, err := toAtlas(*in, context.Background(), k8sClient, "test-namespace")
	require.NoError(t, err)

	require.Equal(t, string(secret.Data["password"]), toAtlasResult.GetLicenseKey())
}

func TestIntegrationsFailReadPassword(t *testing.T) {
	in := &Integration{}

	in.LicenseKeyRef = common.ResourceRefNamespaced{
		Name:      "secret-name",
		Namespace: "secret-namespace",
	}

	testScheme := runtime.NewScheme()
	assert.NoError(t, akov2.AddToScheme(testScheme))
	assert.NoError(t, corev1.AddToScheme(testScheme))
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		Build()

	_, err := toAtlas(*in, context.Background(), k8sClient, "test-namespace")
	require.Error(t, err)
	assert.Contains(t, err.Error(), `secrets "secret-name" not found`)
}
