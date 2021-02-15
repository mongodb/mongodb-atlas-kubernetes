/*
Copyright 2020 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// Important:
// The procedure working with this file:
// 1. Edit the file
// 1. Run "make generate" to regenerate code
// 2. Run "make manifests" to regenerate the CRD

func init() {
	SchemeBuilder.Register(&AtlasProject{}, &AtlasProjectList{})
}

type ScopeType string

const (
	ClusterScopeType  ScopeType = "CLUSTER"
	DataLakeScopeType ScopeType = "DATA_LAKE"
)

// AtlasDatabaseUserSpec defines the desired state of Database User in Atlas
type AtlasDatabaseUserSpec struct {
	// Project is a reference to AtlasProject resource the user belongs to
	Project ResourceRef `json:"projectRef"`

	// DatabaseName is a Database against which Atlas authenticates the user. Default value is 'admin'.
	// +kubebuilder:default=admin
	DatabaseName string `json:"databaseName,omitempty"`

	// DeleteAfterDate is a timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the user.
	// The specified date must be in the future and within one week.
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`

	// Labels is an array containing key-value pairs that tag and categorize the database user.
	// Each key and value has a maximum length of 255 characters.
	Labels []LabelSpec `json:"labels,omitempty"`

	// Roles is an array of this user's roles and the databases / collections on which the roles apply. A role allows
	// the user to perform particular actions on the specified database.
	// +kubebuilder:validation:MinItems=1
	Roles []RoleSpec `json:"roles"`

	// Scopes is an array of clusters and Atlas Data Lakes that this user has access to.
	Scopes []ScopeSpec `json:"scopes,omitempty"`

	// PasswordSecret is a reference to the Secret keeping the user password.
	PasswordSecret *ResourceRef `json:"passwordSecretRef"`

	// Username is a username for authenticating to MongoDB.
	Username string `json:"username"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com

// AtlasDatabaseUser is the Schema for the Atlas Database User API
type AtlasDatabaseUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasDatabaseUserSpec          `json:"spec,omitempty"`
	Status status.AtlasDatabaseUserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasDatabaseUserList contains a list of AtlasDatabaseUser
type AtlasDatabaseUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasDatabaseUser `json:"items"`
}

// RoleSpec allows the user to perform particular actions on the specified database.
// A role on the admin database can include privileges that apply to the other databases as well.
type RoleSpec struct {
	// RoleName is a name of the role. This value can either be a built-in role or a custom role.
	RoleName string `json:"roleName"`

	// DatabaseName is a database on which the user has the specified role. A role on the admin database can include
	// privileges that apply to the other databases.
	DatabaseName string `json:"databaseName"`

	// CollectionName is a collection for which the role applies.
	CollectionName string `json:"collectionName"`
}

// ScopeSpec if present a database user only have access to the indicated resource (Cluster or Atlas Data Lake)
// if none is given then it has access to all.
// It's highly recommended to restrict the access of the database users only to a limited set of resources.
type ScopeSpec struct {
	// Name is a name of the cluster or Atlas Data Lake that the user has access to.
	Name string `json:"name"`
	// Type is a type of resource that the user has access to.
	// +kubebuilder:validation:Enum=CLUSTER;DATA_LAKE
	Type ScopeType `json:"type"`
}

func (p AtlasDatabaseUser) AtlasProjectObjectKey() client.ObjectKey {
	return kube.ObjectKey(p.Namespace, p.Spec.Project.Name)
}

func (p *AtlasDatabaseUser) GetStatus() status.Status {
	return p.Status
}

func (p *AtlasDatabaseUser) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	p.Status.Conditions = conditions
	p.Status.ObservedGeneration = p.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasDatabaseUserStatusOption)
		v(&p.Status)
	}
}
