// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

const (
	// build-in dbroles
	RoleBuildInAdmin        string = "atlasAdmin"
	RoleBuildInReadWriteAny string = "readWriteAnyDatabase"
	RoleBuildInReadAny      string = "readAnyDatabase"

	DefaultDatabaseName = "admin"
)

func BasicUser(crName, atlasUserName string, add ...func(user *akov2.AtlasDatabaseUser)) *akov2.AtlasDatabaseUser {
	user := &akov2.AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name: crName,
		},
		Spec: akov2.AtlasDatabaseUserSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			Username: atlasUserName,
		},
	}
	for _, f := range add {
		f(user)
	}
	return user
}

func WithSecretRef(name string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.PasswordSecret = &common.ResourceRef{Name: name}
	}
}

func WithAdminRole() func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, akov2.RoleSpec{
			RoleName:       RoleBuildInAdmin,
			DatabaseName:   DefaultDatabaseName,
			CollectionName: "",
		})
	}
}

func WithReadWriteRole() func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, akov2.RoleSpec{
			RoleName:       RoleBuildInReadWriteAny,
			DatabaseName:   DefaultDatabaseName,
			CollectionName: "",
		})
	}
}

func WithX509(newUserName string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.Username = newUserName
		user.Spec.DatabaseName = "$external"
		user.Spec.X509Type = "CUSTOMER"
	}
}

func WithCustomRole(role, db, collection string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.Roles = append(user.Spec.Roles, akov2.RoleSpec{
			RoleName:       role,
			DatabaseName:   db,
			CollectionName: collection,
		})
	}
}

func WithNamespace(namespace string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Namespace = namespace
	}
}

func WithOIDCEnabled() func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.OIDCAuthType = "IDP_GROUP"
	}
}

func WithProject(project *akov2.AtlasProject) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.ExternalProjectRef = nil
		user.Spec.ProjectRef = &common.ResourceRefNamespaced{
			Name:      project.Name,
			Namespace: project.Namespace,
		}
	}
}

func WithLabels(labels []common.LabelSpec) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.Labels = labels
	}
}

func WithCredentials(secretName string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.ConnectionSecret = &api.LocalObjectReference{Name: secretName}
	}
}

func WithExternalProjectRef(projectID, credentialsName string) func(user *akov2.AtlasDatabaseUser) {
	return func(user *akov2.AtlasDatabaseUser) {
		user.Spec.ProjectRef = nil
		user.Spec.ExternalProjectRef = &akov2.ExternalProjectReference{
			ID: projectID,
		}
		user.Spec.ConnectionSecret = &api.LocalObjectReference{
			Name: credentialsName,
		}
	}
}
