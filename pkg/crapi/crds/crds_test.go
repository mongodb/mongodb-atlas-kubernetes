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

package crds_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/crds"
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
			got := crds.SelectVersion(tc.spec, tc.version)
			require.Equal(t, tc.want, got)
		})
	}
}
