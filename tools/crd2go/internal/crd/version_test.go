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

package crd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestSelectVersion(t *testing.T) {
	tests := map[string]struct {
		spec            *apiextensionsv1.CustomResourceDefinitionSpec
		version         string
		expectedVersion *VersionedCRD
	}{
		"no versions": {
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind: "TestKind",
				},
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{},
			},
			version:         "v1",
			expectedVersion: nil,
		},
		"empty version string, returns first version": {
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind: "TestKind",
				},
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{Name: "v1"},
					{Name: "v2"},
				},
			},
			version: "",
			expectedVersion: &VersionedCRD{
				Kind:    "TestKind",
				Version: &apiextensionsv1.CustomResourceDefinitionVersion{Name: "v1"},
			},
		},
		"matching version found": {
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind: "TestKind",
				},
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{Name: "v1"},
					{Name: "v2"},
				},
			},
			version: "v2",
			expectedVersion: &VersionedCRD{
				Kind:    "TestKind",
				Version: &apiextensionsv1.CustomResourceDefinitionVersion{Name: "v2"},
			},
		},
		"no matching version found": {
			spec: &apiextensionsv1.CustomResourceDefinitionSpec{
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind: "TestKind",
				},
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{Name: "v1"},
					{Name: "v2"},
				},
			},
			version:         "v3",
			expectedVersion: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			v := SelectVersion(tt.spec, tt.version)
			assert.Equal(t, tt.expectedVersion, v)
		})
	}
}
