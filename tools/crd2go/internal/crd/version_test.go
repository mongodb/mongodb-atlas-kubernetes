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
