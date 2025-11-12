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
//

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
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
