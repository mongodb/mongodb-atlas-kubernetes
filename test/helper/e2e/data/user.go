package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

const (
	// build-in dbroles
	RoleBuildInAdmin        string = "atlasAdmin"
	RoleBuildInReadWriteAny string = "readWriteAnyDatabase"
	RoleBuildInReadAny      string = "readAnyDatabase"

	DefaultDatabaseName = "admin"
)

func BasicUser(crName, atlasUserName string, add ...func(user *v1.AtlasDatabaseUser)) *v1.AtlasDatabaseUser {
	user := &v1.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name: crName,
		},
		Spec: v1.AtlasDatabaseUserSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			Username: atlasUserName,
		},
	}
	for _, f := range add {
		f(user)
	}
	return user
}

func WithSecretRef(name string) func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.PasswordSecret = &common.ResourceRef{Name: name}
	}
}

func WithAdminRole() func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, v1.RoleSpec{
			RoleName:       RoleBuildInAdmin,
			DatabaseName:   DefaultDatabaseName,
			CollectionName: "",
		})
	}
}

func WithReadWriteRole() func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, v1.RoleSpec{
			RoleName:       RoleBuildInReadWriteAny,
			DatabaseName:   DefaultDatabaseName,
			CollectionName: "",
		})
	}
}

func WithX509(newUserName string) func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.Username = newUserName
		user.Spec.DatabaseName = "$external"
		user.Spec.X509Type = "CUSTOMER"
	}
}

func WithCustomRole(role, db, collection string) func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, v1.RoleSpec{
			RoleName:       role,
			DatabaseName:   db,
			CollectionName: collection,
		})
	}
}

func WithNamespace(namespace string) func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Namespace = namespace
	}
}

func WithOIDCEnabled() func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.OIDCAuthType = "IDP_GROUP"
	}
}

func WithProject(project *v1.AtlasProject) func(user *v1.AtlasDatabaseUser) {
	return func(user *v1.AtlasDatabaseUser) {
		user.Spec.Project = common.ResourceRefNamespaced{
			Name:      project.Name,
			Namespace: project.Namespace,
		}
	}
}
