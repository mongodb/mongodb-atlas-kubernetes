package model

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/utils"
)

type UserSpec v1.AtlasDatabaseUserSpec

type DBUser struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

type UserCustomRoleType string

const (
	// build-in dbroles
	RoleBuildInAdmin        string = "atlasAdmin"
	RoleBuildInReadWriteAny string = "readWriteAnyDatabase"
	RoleBuildInReadAny      string = "readAnyDatabase"

	RoleCustomAdmin     UserCustomRoleType = "dbAdmin"
	RoleCustomReadWrite UserCustomRoleType = "readWrite"
	RoleCustomRead      UserCustomRoleType = "read"
)

func NewDBUser(userName string) *DBUser {
	return &DBUser{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "atlas.mongodb.com/v1",
			Kind:       "AtlasDatabaseUser",
		},
		ObjectMeta: &metav1.ObjectMeta{
			Name: "k-" + userName,
		},
		Spec: UserSpec{
			Username: userName,
			Project: v1.ResourceRefNamespaced{
				Name: "my-project",
			},
		},
	}
}

func (u *DBUser) WithAuthDatabase(name string) *DBUser {
	u.Spec.DatabaseName = name
	return u
}

func (s *DBUser) WithProjectRef(name string) *DBUser {
	s.Spec.Project.Name = name
	return s
}

func (s *DBUser) WithSecretRef(name string) *DBUser {
	s.Spec.PasswordSecret = &v1.ResourceRef{Name: name}
	return s
}

func (s *DBUser) AddBuildInAdminRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, v1.RoleSpec{
		RoleName:       RoleBuildInAdmin,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddBuildInReadAnyRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, v1.RoleSpec{
		RoleName:       RoleBuildInReadAny,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddBuildInReadWriteRole() *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, v1.RoleSpec{
		RoleName:       RoleBuildInReadWriteAny,
		DatabaseName:   "admin",
		CollectionName: "",
	})
	return s
}

func (s *DBUser) AddCustomRole(role UserCustomRoleType, db string, collection string) *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, v1.RoleSpec{
		RoleName:       string(role),
		DatabaseName:   db,
		CollectionName: collection,
	})
	return s
}

func (s *DBUser) DeleteAllRoles() *DBUser {
	s.Spec.Roles = []v1.RoleSpec{}
	return s
}

func (s *DBUser) GetFilePath(projectName string) string {
	return filepath.Join(projectName, "user", "user-"+s.ObjectMeta.Name+".yaml")
}

func (s *DBUser) SaveConfigurationTo(folder string) {
	folder = filepath.Dir(folder)
	yamlConf := utils.JSONToYAMLConvert(s)
	utils.SaveToFile(s.GetFilePath(folder), yamlConf)
}
