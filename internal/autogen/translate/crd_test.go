// Copyright 2025 Google LLC
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

package translate

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestSelectVersion(t *testing.T) {
	// Define some sample versions to be reused in the tests
	v1 := apiextensionsv1.CustomResourceDefinitionVersion{Name: "v1", Served: true}
	v2 := apiextensionsv1.CustomResourceDefinitionVersion{Name: "v2", Served: true}
	v1beta1 := apiextensionsv1.CustomResourceDefinitionVersion{Name: "v1beta1", Served: true}

	testCases := []struct {
		name    string
		spec    *apiextensionsv1.CustomResourceDefinitionSpec
		version string
		want    *apiextensionsv1.CustomResourceDefinitionVersion
	}{
		{
			name: "should return nil when spec has no versions",
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{},
			},
			version: "v1",
			want:    nil,
		},
		{
			name: "should return the first version when version string is empty",
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{v1beta1, v1},
			},
			version: "",
			want:    &v1beta1,
		},
		{
			name: "should return the matching version when it exists",
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{v1, v2},
			},
			version: "v2",
			want:    &v2,
		},
		{
			name: "should return nil when the requested version does not exist",
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{v1, v2},
			},
			version: "v3",
			want:    nil,
		},
		{
			name: "should handle a single version in the spec",
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{v1},
			},
			version: "v1",
			want:    &v1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := selectVersion(tc.spec, tc.version)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestGetOpenAPIProperties(t *testing.T) {
	const kind = "MyTestCRD"

	happyPathVersion := &apiextensionsv1.CustomResourceDefinitionVersion{
		Name: "v1",
		Schema: &apiextensionsv1.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
				Properties: map[string]apiextensionsv1.JSONSchemaProps{
					"spec": {
						Type: "object",
						Properties: map[string]apiextensionsv1.JSONSchemaProps{
							"fieldA": {Type: "string"},
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		name           string
		version        *apiextensionsv1.CustomResourceDefinitionVersion
		wantProperties map[string]apiextensionsv1.JSONSchemaProps
		wantErrMsg     string
	}{
		{
			name:           "should return properties",
			version:        happyPathVersion,
			wantProperties: happyPathVersion.Schema.OpenAPIV3Schema.Properties,
			wantErrMsg:     "",
		},
		{
			name:           "should return error if version is nil",
			version:        nil,
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("missing version (nil) from %v spec", kind),
		},
		{
			name: "should return error if schema is nil",
			version: &apiextensionsv1.CustomResourceDefinitionVersion{
				Name:   "v1",
				Schema: nil, // The point of failure
			},
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("missing version schema from %v spec", kind),
		},
		{
			name: "error - should return error if OpenAPIV3Schema is nil",
			version: &apiextensionsv1.CustomResourceDefinitionVersion{
				Name: "v1",
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: nil, // The point of failure
				},
			},
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("missing version OpenAPI Schema from %v spec", kind),
		},
		{
			name: "should return error if Properties map is nil",
			version: &apiextensionsv1.CustomResourceDefinitionVersion{
				Name: "v1",
				Schema: &apiextensionsv1.CustomResourceValidation{
					OpenAPIV3Schema: &apiextensionsv1.JSONSchemaProps{
						Properties: nil, // The point of failure
					},
				},
			},
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("missing version OpenAPI Properties from %v spec", kind),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotProperties, err := getOpenAPIProperties(kind, tc.version)
			if tc.wantErrMsg != "" {
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErrMsg)
				require.Nil(t, gotProperties)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantProperties, gotProperties)
			}
		})
	}
}

func TestGetSpecPropertiesFor(t *testing.T) {
	const kind = "MyTestCRD"

	// Reusable properties for the test cases
	nestedProperties := map[string]apiextensionsv1.JSONSchemaProps{
		"replicas": {Type: "integer"},
		"image":    {Type: "string"},
	}

	testCases := []struct {
		name           string
		props          map[string]apiextensionsv1.JSONSchemaProps
		field          string
		wantProperties map[string]apiextensionsv1.JSONSchemaProps
		wantErrMsg     string
	}{
		{
			name: "should return nested properties",
			props: map[string]apiextensionsv1.JSONSchemaProps{
				"spec": {
					Type:       "object",
					Properties: nestedProperties,
				},
				"status": {
					Type: "object",
				},
			},
			field:          "spec",
			wantProperties: nestedProperties,
			wantErrMsg:     "", // Expect no error
		},
		{
			name: "field is missing from properties map",
			props: map[string]apiextensionsv1.JSONSchemaProps{
				"status": {Type: "object"},
			},
			field:          "spec", // This field does not exist
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("kind %q spec is missing field %q on", kind, "spec"),
		},
		{
			name: "field is not of type object",
			props: map[string]apiextensionsv1.JSONSchemaProps{
				"spec": {Type: "string"}, // The point of failure
			},
			field:          "spec",
			wantProperties: nil,
			wantErrMsg:     fmt.Sprintf("kind %q field %q expected to be object but is %v", kind, "spec", "string"),
		},
		{
			name: "field is an object but has no nested properties",
			props: map[string]apiextensionsv1.JSONSchemaProps{
				"spec": {
					Type:       "object",
					Properties: nil, // The nested map is nil
				},
			},
			field:          "spec",
			wantProperties: nil, // Expecting a nil map is correct here
			wantErrMsg:     "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			gotProperties, err := getSpecPropertiesFor(kind, tc.props, tc.field)

			// Assert
			if tc.wantErrMsg != "" {
				// We expect an error
				require.Error(t, err)
				require.EqualError(t, err, tc.wantErrMsg)
				require.Nil(t, gotProperties)
			} else {
				// We expect success
				require.NoError(t, err)
				require.Equal(t, tc.wantProperties, gotProperties)
			}
		})
	}
}
