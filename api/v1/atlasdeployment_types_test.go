package v1

import (
	"testing"
)

func TestDedploymentProjectRefCELValidations(t *testing.T) {
	launchProjectRefCELTests(
		t,
		func(pdr *ProjectDualReference) AtlasCustomResource {
			d := AtlasDeployment{}
			if pdr != nil {
				setDualRef(d.ProjectDualRef(), pdr)
			}
			return &d
		},
		"../../config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml",
	)
}
