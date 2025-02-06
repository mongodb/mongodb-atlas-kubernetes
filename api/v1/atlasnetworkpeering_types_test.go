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
		obj            *AtlasNetworkPeering
		expectedErrors []string
	}{
		{
			title: "Missing container ref in peering fails",
			obj: &AtlasNetworkPeering{
				Spec: AtlasNetworkPeeringSpec{},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must either have a container Atlas id or Kubernetes name, but not both (or neither)"},
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
			expectedErrors: []string{"spec: Invalid value: \"object\": must either have a container Atlas id or Kubernetes name, but not both (or neither)"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			// inject a project to avoid other CEL validations being hit
			tc.obj.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: "some-project"}
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			crdPath := "../../config/crd/bases/atlas.mongodb.com_atlasnetworkpeerings.yaml"
			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, nil)

			require.Equal(t, tc.expectedErrors, cel.ErrorListAsStrings(errs))
		})
	}
}
