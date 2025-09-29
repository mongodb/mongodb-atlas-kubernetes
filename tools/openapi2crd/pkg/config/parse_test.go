package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		raw            []byte
		expectedConfig *v1alpha1.Config
		expectedError  error
	}{
		"valid config": {
			raw: []byte(`
kind: Config
apiVersion: atlas2crd.mongodb.com/v1alpha1
spec:
  crd: []
  openapi: []
`),
			expectedConfig: &v1alpha1.Config{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Config",
					APIVersion: "atlas2crd.mongodb.com/v1alpha1",
				},
				Spec: v1alpha1.Spec{
					CRDConfig:          []v1alpha1.CRDConfig{},
					OpenAPIDefinitions: []v1alpha1.OpenAPIDefinition{},
				},
			},
		},
		"invalid config": {
			raw:            []byte("invalid yaml"),
			expectedConfig: nil,
			expectedError: fmt.Errorf(
				"error unmarshalling config type: %w",
				fmt.Errorf("error unmarshaling JSON: %w", fmt.Errorf("while decoding JSON: json: cannot unmarshal string into Go value of type v1alpha1.Config")),
			),
		},
		"invalid kind": {
			raw: []byte(`
kind: Other
apiVersion: atlas2crd.mongodb.com/v1alpha1
spec:
  crd: []
  openapi: []
`),
			expectedConfig: nil,
			expectedError:  fmt.Errorf("invalid config type: %s", "Other"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := Parse(tt.raw)
			require.Equal(t, tt.expectedError, err)
			require.Equal(t, tt.expectedConfig, cfg)
		})
	}
}
