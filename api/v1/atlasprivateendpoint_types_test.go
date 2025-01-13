package v1

import (
	"testing"
)

func TestPrivateEndpointProjectRefCELValidations(t *testing.T) {
	launchProjectRefCELTests(
		t,
		func(pdr *ProjectDualReference) AtlasCustomResource {
			pe := AtlasPrivateEndpoint{}
			if pdr != nil {
				setDualRef(pe.ProjectDualRef(), pdr)
			}
			return &pe
		},
		"../../config/crd/bases/atlas.mongodb.com_atlasprivateendpoints.yaml",
	)
}
