package v1

import (
	"testing"
)

func TestDatabaseUsersProjectRefCELValidations(t *testing.T) {
	launchProjectRefCELTests(
		t,
		func(pdr *ProjectDualReference) AtlasCustomResource {
			dbu := AtlasDatabaseUser{}
			if pdr != nil {
				setDualRef(dbu.ProjectDualRef(), pdr)
			}
			return &dbu
		},
		"../../config/crd/bases/atlas.mongodb.com_atlasdatabaseusers.yaml",
	)
}
