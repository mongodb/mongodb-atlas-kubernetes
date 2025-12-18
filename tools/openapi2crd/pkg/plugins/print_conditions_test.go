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
			title:      "columns added to spec",
			apiVersion: "v1",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1beta1", Served: true, Storage: false},
				{Name: "v1", Served: true, Storage: true},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				// Verify columns were added at the spec level
				require.Len(t, crd.Spec.AdditionalPrinterColumns, 3)

				assert.Equal(t, "Ready", crd.Spec.AdditionalPrinterColumns[0].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].status`, crd.Spec.AdditionalPrinterColumns[0].JSONPath)

				assert.Equal(t, "Reason", crd.Spec.AdditionalPrinterColumns[1].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].reason`, crd.Spec.AdditionalPrinterColumns[1].JSONPath)

				assert.Equal(t, "State", crd.Spec.AdditionalPrinterColumns[2].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="State")].reason`, crd.Spec.AdditionalPrinterColumns[2].JSONPath)
			},
		},
		{
			title:      "replace existing columns",
			apiVersion: "v1",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1", Served: true, Storage: true},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				require.Len(t, crd.Spec.AdditionalPrinterColumns, 3)
				assert.Equal(t, "Ready", crd.Spec.AdditionalPrinterColumns[0].Name)
				assert.Equal(t, "Reason", crd.Spec.AdditionalPrinterColumns[1].Name)
				assert.Equal(t, "State", crd.Spec.AdditionalPrinterColumns[2].Name)
			},
		},
		{
			title:      "columns added with single version",
			apiVersion: "",
			versions: []apiextensions.CustomResourceDefinitionVersion{
				{Name: "v1", Served: true, Storage: true},
			},
			validate: func(t *testing.T, crd *apiextensions.CustomResourceDefinition) {
				// Verify columns were added at the spec level
				require.Len(t, crd.Spec.AdditionalPrinterColumns, 3)

				assert.Equal(t, "Ready", crd.Spec.AdditionalPrinterColumns[0].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].status`, crd.Spec.AdditionalPrinterColumns[0].JSONPath)

				assert.Equal(t, "Reason", crd.Spec.AdditionalPrinterColumns[1].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="Ready")].reason`, crd.Spec.AdditionalPrinterColumns[1].JSONPath)

				assert.Equal(t, "State", crd.Spec.AdditionalPrinterColumns[2].Name)
				assert.Equal(t, `.status.conditions[?(@.type=="State")].reason`, crd.Spec.AdditionalPrinterColumns[2].JSONPath)
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
