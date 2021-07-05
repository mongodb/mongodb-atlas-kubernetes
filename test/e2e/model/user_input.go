package model

import (
	"fmt"
	"path/filepath"

	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type UserInputs struct {
	AtlasKeyAccessType AtlasKeyType
	ProjectID          string
	KeyName            string
	Namespace          string
	K8sFullProjectName string
	ProjectPath        string
	Clusters           []AC
	Users              []DBUser
	Project            *AP
}

// NewUsersInputs prepare users inputs
func NewUserInputs(keyTestPrefix string, users []DBUser, r *AtlasKeyType) UserInputs {
	projectName := fmt.Sprintf("%s-%s", keyTestPrefix, utils.GenID())
	input := UserInputs{
		AtlasKeyAccessType: *r,
		ProjectID:          "",
		KeyName:            keyTestPrefix,
		Namespace:          "ns-" + projectName,
		K8sFullProjectName: "atlasproject.atlas.mongodb.com/k-" + projectName,
		ProjectPath:        filepath.Join(DataFolder, projectName, "resources", projectName+".yaml"),
	}
	input.Project = NewProject("k-"+projectName).ProjectName(projectName).WithIpAccess("0.0.0.0/0", "everyone")
	if !r.GlobalLevelKey {
		input.Project = input.Project.SecretRef(keyTestPrefix)
	}

	for _, user := range users {
		input.Users = append(input.Users, *user.WithProjectRef(input.Project.GetK8sMetaName()))
	}
	return input
}

func (u *UserInputs) GetAppFolder() string {
	return filepath.Join(DataFolder, u.Project.Spec.Name, "app")
}

func (u *UserInputs) GetOperatorFolder() string {
	return filepath.Join(DataFolder, u.Project.Spec.Name, "operator")
}

func (u *UserInputs) GetResourceFolder() string {
	return filepath.Dir(u.ProjectPath)
}

func (u *UserInputs) GetUsersFolder() string {
	return filepath.Join(u.GetResourceFolder(), "user")
}

func (u *UserInputs) GetServiceCatalogSourceFolder() string {
	return filepath.Join(DataFolder, u.Project.Spec.Name, "catalog")
}
