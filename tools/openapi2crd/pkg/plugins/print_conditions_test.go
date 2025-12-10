package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/plugins"
)

func TestPrintConditions_Name(t *testing.T) {
	p := &plugins.PrintConditions{}
	assert.Equal(t, "print_conditions", p.Name())
}

func TestPrintConditions_Process(t *testing.T) {
	tests := []struct {
		title      string
		apiVersion string
		versions   []apiextensions.CustomResourceDefinitionVersion
		wantErr    string
		// validate allows us to inspect the CRD after processing to ensure side effects occurred
		validate func(t *testing.T, crd *apiextensions.CustomResourceDefinition)
	}{
		{
			title:      "version updated",
			apiVersion: "v1",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1beta1", Served: true, Storage: false},
				{Name: "v1", Served: true, Storage: true},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				// Find the v1 version
				var targetVersion apiextensions.CustomResourceDefinitionVersion
				for _, v := range crd.Spec.Versions {
					if v.Name == "v1" {
						targetVersion = v
						break
					}
				}

				// Verify columns were added
				require.NotEmpty(t, targetVersion.Name, "Should have found the target version")
				require.Len(t, targetVersion.AdditionalPrinterColumns, 2)

				assert.Equal(t, "Ready", targetVersion.AdditionalPrinterColumns[0].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].status`, targetVersion.AdditionalPrinterColumns[0].JSONPath)

				assert.Equal(t, "State", targetVersion.AdditionalPrinterColumns[1].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="State")].reason`, targetVersion.AdditionalPrinterColumns[1].JSONPath)
			},
		},
		{
			title:      "missing version",
			apiVersion: "v2",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1", Served: true, Storage: true},
			},
			wantErr: "apiVersion \"v2\" not listed in spec",
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				// Ensure no weird side effects happened to existing versions
				assert.Empty(t, crd.Spec.Versions[0].AdditionalPrinterColumns)
			},
		},
		{
			title:      "replace columns",
			apiVersion: "v1",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name: "v1",
					AdditionalPrinterColumns: []apiextensions.CustomResourceColumnDefinition{
						{Name: "OldColumn", JSONPath: ".old"},
					},
				},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				targetVersion := crd.Spec.Versions[0]
				require.Len(t, targetVersion.AdditionalPrinterColumns, 2)
				assert.Equal(t, "Ready", targetVersion.AdditionalPrinterColumns[0].Name)
			},
		},
		{
			title:      "unique version updated",
			apiVersion: "",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1", Served: true, Storage: true},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				// Find the v1 version
				var targetVersion apiextensions.CustomResourceDefinitionVersion
				for _, v := range crd.Spec.Versions {
					if v.Name == "v1" {
						targetVersion = v
						break
					}
				}

				// Verify columns were added
				require.NotEmpty(t, targetVersion.Name, "Should have found the target version")
				require.Len(t, targetVersion.AdditionalPrinterColumns, 2)

				assert.Equal(t, "Ready", targetVersion.AdditionalPrinterColumns[0].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].status`, targetVersion.AdditionalPrinterColumns[0].JSONPath)

				assert.Equal(t, "State", targetVersion.AdditionalPrinterColumns[1].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="State")].reason`, targetVersion.AdditionalPrinterColumns[1].JSONPath)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			p := &plugins.PrintConditions{}
			crd := &apiextensions.CustomResourceDefinition{
				Spec: apiextensions.CustomResourceDefinitionSpec{
					Versions: tc.versions,
				},
			}
			crd.APIVersion = tc.apiVersion
			req := &plugins.MappingProcessorRequest{
				CRD: crd,
			}

			err := p.Process(req)

			if tc.wantErr != "" {
				assert.EqualError(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			if tc.validate != nil {
				tc.validate(t, crd)
			}
		})
	}
}
