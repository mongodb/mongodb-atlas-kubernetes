package model

import (
	"path/filepath"

	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO PRIVATE?
type UserInputs struct {
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
func NewUserInputs(keyName string, users []DBUser) UserInputs {
	uid := utils.GenUniqID()
	input := UserInputs{
		ProjectID:          "",
		KeyName:            keyName,
		Namespace:          "ns-" + uid,
		K8sFullProjectName: "atlasproject.atlas.mongodb.com/k-" + uid,
		ProjectPath:        filepath.Join(DataFolder, uid, "resources", uid+".yaml"),
	}
	input.Project = NewProject("k-"+uid).ProjectName(uid).SecretRef(keyName).WithIpAccess("0.0.0.0/0", "everyone")
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
