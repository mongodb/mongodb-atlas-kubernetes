package model

import (
	// "os"
	"path/filepath"

	// . "github.com/onsi/gomega"
	// . "github.com/onsi/gomega/gstruct"
	// "go.mongodb.org/atlas/mongodbatlas"

	// kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kube"
	// mongocli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/mongocli"
	. "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

// TODO PRIVATE?
type UserInputs struct {
	ProjectName        string
	ProjectID          string
	KeyName            string
	Namespace          string
	K8sProjectName     string
	K8sFullProjectName string
	ProjectPath        string
	Clusters           []AC
	Users              []DBUser
}

// NewUsersInputs prepare users inputs
func NewUserInputs(keyName string, users []DBUser) UserInputs {
	uid := utils.GenUniqID()

	input := UserInputs{
		ProjectName:        uid,
		ProjectID:          "",
		KeyName:            keyName,
		Namespace:          "ns-" + uid,
		K8sProjectName:     "k-" + uid,
		K8sFullProjectName: "atlasproject.atlas.mongodb.com/k-" + uid,
		ProjectPath:        filepath.Join(DataFolder, uid, "resources", uid+".yaml"),
	}
	for _, user := range users {
		input.Users = append(input.Users, *user.WithProjectRef(input.K8sProjectName))
	}
	return input
}

func (u *UserInputs) GetAppFolder() string {
	return filepath.Join(DataFolder, u.ProjectName, "app")
}

func (u *UserInputs) GetOperatorFolder() string {
	return filepath.Join(DataFolder, u.ProjectName, "operator")
}

func (u *UserInputs) GetResourceFolder() string {
	return filepath.Dir(u.ProjectPath)
}

func (u *UserInputs) GetUsersFolder() string {
	return filepath.Join(u.GetResourceFolder(), "user")
}