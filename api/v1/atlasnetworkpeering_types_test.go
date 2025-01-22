package v1 // nolint: dupl

import (
	"testing"
)

func TestNetworkPeeringProjectRefCELValidations(t *testing.T) {
	launchProjectRefCELTests(
		t,
		func(pdr *ProjectDualReference) AtlasCustomResource {
			np := AtlasNetworkPeering{}
			if pdr != nil {
				setDualRef(np.ProjectDualRef(), pdr)
			}
			return &np
		},
		"../../config/crd/bases/atlas.mongodb.com_atlasnetworkpeerings.yaml",
	)
}
