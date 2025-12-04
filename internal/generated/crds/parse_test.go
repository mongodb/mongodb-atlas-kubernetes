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

package crds

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestParseCRD(t *testing.T) {
	tests := map[string]struct {
		scanner     *bufio.Scanner
		expectedCrd *apiextensionsv1.CustomResourceDefinition
		expectedErr string
	}{
		"valid CRD": {
			scanner:     bufio.NewScanner(strings.NewReader(validCRDManifest(t))),
			expectedCrd: validCRDObject(t),
		},
		"not a CRD": {
			scanner:     bufio.NewScanner(strings.NewReader("apiVersion: autoscaling/__internal\nkind: Scale\nmetadata:\n  name: test-scale\n")),
			expectedErr: "failed to decode YAML: no kind \"Scale\" is registered for the internal version of group \"autoscaling\" in scheme \"pkg/runtime/scheme.go:110\"",
		},
		"empty input": {
			scanner:     bufio.NewScanner(strings.NewReader("")),
			expectedErr: "EOF",
		},
		"only comments": {
			scanner:     bufio.NewScanner(strings.NewReader("# This is a comment\n# Another comment line\n")),
			expectedErr: "EOF",
		},
		"multiple CRDs, returns first": {
			scanner:     bufio.NewScanner(strings.NewReader(validCRDManifest(t) + "---\n" + validCRDManifest(t))),
			expectedCrd: validCRDObject(t),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseCRDs(tt.scanner)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, tt.expectedErr, gotErr)
			assert.Equal(t, tt.expectedCrd, got)
		})
	}
}

func validCRDManifest(t *testing.T) string {
	t.Helper()

	return `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: tests.example.com
spec:
  group: example.com
  names:
    kind: Test
    plural: tests
    singular: test
  scope: Namespaced
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                field1:
                  type: string
  storedVersions: ["v1"]
`
}

func validCRDObject(t *testing.T) *apiextensionsv1.CustomResourceDefinition {
	t.Helper()

	return &apiextensionsv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tests.example.com",
		},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group: "example.com",
			Names: apiextensionsv1.CustomResourceDefinitionNames{
				Kind:     "Test",
				Plural:   "tests",
				Singular: "test",
			},
			Scope: "Namespaced",
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
				{
					Name:    "v1",
					Served:  true,
					Storage: true,
					Schema: &apiextensionsv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
							Type: "object",
							Properties: map[string]apiextensionsv1.JSONSchemaProps{
								"spec": {
									Type: "object",
									Properties: map[string]apiextensionsv1.JSONSchemaProps{
										"field1": {
											Type: "string",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
