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

package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestContainerCELChecks(t *testing.T) {
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
			expectedErrors: []string{"spec: Invalid value: must not set region for GCP containers"},
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
			expectedErrors: []string{"spec: Invalid value: must set region for AWS and Azure containers"},
		},
		{
			title: "Azure fails without a region",
			obj: &AtlasNetworkContainer{
				Spec: AtlasNetworkContainerSpec{
					Provider: string(provider.ProviderAzure),
				},
			},
			expectedErrors: []string{"spec: Invalid value: must set region for AWS and Azure containers"},
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
			expectedErrors: []string{"spec: Invalid value: id is immutable"},
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
			expectedErrors: []string{"spec: Invalid value: region is immutable"},
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

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
