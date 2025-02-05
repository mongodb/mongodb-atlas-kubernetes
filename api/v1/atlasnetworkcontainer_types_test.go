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

func TestCELChecks(t *testing.T) {
	for _, tc := range []struct {
		title          string
		old, obj       *AtlasNetworkContainer
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
		{
			title: "ID cannot be changed",
			old: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						ID: "old-id",
					},
				},
			},
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						ID: "new-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": id is immutable"},
		},
		{
			title: "ID can be unset",
			old: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
				},
			},
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
				},
			},
		},
		{
			title: "ID can be set",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						ID: "new-id",
					},
				},
			},
		},
		{
			title: "Region cannot be changed",
			old: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "old-region",
					},
				},
			},
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "new-region",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": region is immutable"},
		},
		{
			title: "Region can be unset (for GCP)",
			old: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
				},
			},
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderGCP),
				},
			},
		},
		{
			title: "Region can be set",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: AtlasNetworkContainerConfig{
						Region: "new-region",
					},
				},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.old)
			require.NoError(t, err)
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasnetworkcontainers.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, unstructuredOldObject)

			require.Equal(t, len(tc.expectedErrors), len(errs), fmt.Sprintf("errs: %v", errs))

			for i, err := range errs {
				assert.Equal(t, tc.expectedErrors[i], err.Error())
			}
		})
	}
}
