package v1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestGCPRegionCELCcheck(t *testing.T) {
	for _, tc := range []struct {
		title          string
		obj            *AtlasNetworkContainer
		expectedErrors []string
	}{
		{
			title: "GCP fails with a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "some-region",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must not set region for GCP containers"},
		},
		{
			title: "GCP succeeds without a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
				},
			},
		},
		{
			title: "AWS succeeds with a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "some-region",
					},
				},
			},
		},
		{
			title: "Azure succeeds with a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAzure),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "some-region",
					},
				},
			},
		},
		{
			title: "AWS fails without a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAWS),
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must set region for AWS and Azure containers"},
		},
		{
			title: "Azure fails without a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAzure),
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must set region for AWS and Azure containers"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasnetworkcontainers.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, nil)

			require.Equal(t, len(tc.expectedErrors), len(errs), fmt.Sprintf("errs: %v", errs))

			for i, err := range errs {
				assert.Equal(t, tc.expectedErrors[i], err.Error())
			}
		})
	}
}
