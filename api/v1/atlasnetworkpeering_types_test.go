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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestPeeringCELChecks(t *testing.T) {
	for _, tc := range []struct {
		title          string
		old            *AtlasNetworkPeering
		obj            *AtlasNetworkPeering
		expectedErrors []string
	}{
		{
			title: "Missing container ref in peering fails",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{},
			},
			expectedErrors: []string{"spec: Invalid value: must either have a container Atlas id or Kubernetes name, but not both (or neither)"},
		},

		{
			title: "Named container ref works",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "Some-name",
					},
				},
			},
		},

		{
			title: "Container id ref works",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
				},
			},
		},

		{
			title: "Both container id and name ref fails",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "Some-name",
						ID:   "some-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: must either have a container Atlas id or Kubernetes name, but not both (or neither)"},
		},

		{
			title: "ID set works",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
		},

		{
			title: "ID added in fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": no such key: id evaluating rule: id is immutable"},
		},

		{
			title: "ID remove fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": no such key: id evaluating rule: id is immutable"},
		},

		{
			title: "ID changed fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "another-peering-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: id is immutable"},
		},

		{
			title: "ID unchanged works",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
		},

		{
			title: "container ID changed fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-other-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: container ref id is immutable"},
		},

		{
			title: "container name changed fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-other-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			expectedErrors: []string{"spec: Invalid value: container ref name is immutable"},
		},

		{
			title: "change container name to id fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": no such key: name evaluating rule: container ref name is immutable",
				"spec: Invalid value: \"object\": no such key: id evaluating rule: container ref id is immutable",
			},
		},

		{
			title: "change container id to name fails",
			old: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						ID: "some-id",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{
					ContainerRef: ContainerDualReference{
						Name: "some-name",
					},
					AtlasNetworkPeeringConfig: AtlasNetworkPeeringConfig{
						ID: "some-peering-id",
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": no such key: name evaluating rule: container ref name is immutable",
				"spec: Invalid value: \"object\": no such key: id evaluating rule: container ref id is immutable",
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			unstructuredOld := map[string]any{}
			if tc.old != nil {
				var err error
				tc.old.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
				unstructuredOld, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.old)
				require.NoError(t, err)
			}

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasnetworkpeerings.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, unstructuredOld)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
