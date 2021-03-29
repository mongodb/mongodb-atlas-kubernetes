package utils

import (
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

type UserSpec v1.AtlasDatabaseUserSpec

type DBUser struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      *metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

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
			Project: v1.ResourceRef{
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

func (s *DBUser) AddRole(role string, db string, collection string) *DBUser {
	s.Spec.Roles = append(s.Spec.Roles, v1.RoleSpec{
		RoleName:       role,
		DatabaseName:   db,
		CollectionName: collection,
	})
	return s
}

func (s *DBUser) GetFilePath(ns string) string {
	return filepath.Join("data", ns, "user", "user-"+s.ObjectMeta.Name+".yaml")
}

func (s *DBUser) SaveConfigurationTo(ns string) {
	yamlConf := JSONToYAMLConvert(s)
	SaveToFile(s.GetFilePath(ns), yamlConf)
}
