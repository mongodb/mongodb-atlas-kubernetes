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
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

// Important:
// The procedure working with this file:
// 1. Edit the file
// 1. Run "make generate" to regenerate code
// 2. Run "make manifests" to regenerate the CRD

func init() {
	SchemeBuilder.Register(&AtlasDatabaseUser{}, &AtlasDatabaseUserList{})
}

type ScopeType string

const (
	DeploymentScopeType ScopeType = "CLUSTER" // todo: potentially rename to "DEPLOYMENT"
	DataLakeScopeType   ScopeType = "DATA_LAKE"
)

// AtlasDatabaseUserSpec defines the desired state of Database User in Atlas
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
type AtlasDatabaseUserSpec struct {
	api.LocalCredentialHolder `json:",inline"`

	// Project is a reference to AtlasProject resource the user belongs to
	Project *common.ResourceRefNamespaced `json:"projectRef,omitempty"`
	// ExternalProjectRef holds the Atlas project ID the user belongs to
	ExternalProjectRef *ExternalProjectReference `json:"externalProjectRef,omitempty"`

	// DatabaseName is a Database against which Atlas authenticates the user.
	// If the user authenticates with AWS IAM, x.509, LDAP, or OIDC Workload this value should be '$external'.
	// If the user authenticates with SCRAM-SHA or OIDC Workforce, this value should be 'admin'.
	// Default value is 'admin'.
	// +kubebuilder:default=admin
	DatabaseName string `json:"databaseName,omitempty"`

	// DeleteAfterDate is a timestamp in ISO 8601 date and time format in UTC after which Atlas deletes the user.
	// The specified date must be in the future and within one week.
	DeleteAfterDate string `json:"deleteAfterDate,omitempty"`

	// Labels is an array containing key-value pairs that tag and categorize the database user.
	// Each key and value has a maximum length of 255 characters.
	Labels []common.LabelSpec `json:"labels,omitempty"`

	// Roles is an array of this user's roles and the databases / collections on which the roles apply. A role allows
	// the user to perform particular actions on the specified database.
	// +kubebuilder:validation:MinItems=1
	Roles []RoleSpec `json:"roles"`

	// Scopes is an array of clusters and Atlas Data Lakes that this user has access to.
	Scopes []ScopeSpec `json:"scopes,omitempty"`

	// PasswordSecret is a reference to the Secret keeping the user password.
	PasswordSecret *common.ResourceRef `json:"passwordSecretRef,omitempty"`

	// Username is a username for authenticating to MongoDB
	// Human-readable label that represents the user that authenticates to MongoDB. The format of this label depends on the method of authentication:
	// In case of AWS IAM: the value should be AWS ARN for the IAM User/Role;
	// In case of OIDC Workload or Workforce: the value should be the Atlas OIDC IdP ID, followed by a '/', followed by the IdP group name;
	// In case of Plain text auth: the value can be anything
	// +kubebuilder:validation:MaxLength:=1024
	Username string `json:"username"`

	// Human-readable label that indicates whether the new database Username with OIDC federated authentication.
	// To create a federated authentication group (Workforce), specify the value of IDP_GROUP in this field.
	// To create a federated authentication user (Workload), specify the value of USER in this field.
	// +kubebuilder:default:=NONE
	// +kubebuilder:validation:Enum:=NONE;IDP_GROUP;USER
	// +optional
	OIDCAuthType string `json:"oidcAuthType,omitempty"`

	// Human-readable label that indicates whether the new database
	// user authenticates with the Amazon Web Services (AWS)
	// Identity and Access Management (IAM) credentials associated with
	// the user or the user's role
	// +kubebuilder:default:=NONE
	// +kubebuilder:validation:Enum:=NONE;USER;ROLE
	// +optional
	AWSIAMType string `json:"awsIamType,omitempty"`

	// X509Type is X.509 method by which the database authenticates the provided username
	// +kubebuilder:default:=NONE
	// +kubebuilder:validation:Enum:=NONE;MANAGED;CUSTOMER
	// +optional
	X509Type string `json:"x509Type,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Username",type=string,JSONPath=`.spec.username`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=adu

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
	CollectionName string `json:"collectionName,omitempty"`
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
	ns := p.Namespace
	if p.Spec.Project.Namespace != "" {
		ns = p.Spec.Project.Namespace
	}
	return kube.ObjectKey(ns, p.Spec.Project.Name)
}

func (p AtlasDatabaseUser) PasswordSecretObjectKey() *client.ObjectKey {
	if p.Spec.PasswordSecret != nil {
		key := kube.ObjectKey(p.Namespace, p.Spec.PasswordSecret.Name)
		return &key
	}
	return nil
}

func (p *AtlasDatabaseUser) GetStatus() api.Status {
	return p.Status
}

func (p *AtlasDatabaseUser) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	p.Status.Conditions = conditions
	p.Status.ObservedGeneration = p.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasDatabaseUserStatusOption)
		v(&p.Status)
	}
}

func (p *AtlasDatabaseUser) ReadPassword(ctx context.Context, kubeClient client.Client) (string, error) {
	if p.Spec.PasswordSecret != nil {
		secret := &corev1.Secret{}
		if err := kubeClient.Get(ctx, *p.PasswordSecretObjectKey(), secret); err != nil {
			return "", err
		}
		p, exist := secret.Data["password"]
		switch {
		case !exist:
			return "", fmt.Errorf("secret %s is invalid: it doesn't contain 'password' field", secret.Name)
		case len(p) == 0:
			return "", fmt.Errorf("secret %s is invalid: the 'password' field is empty", secret.Name)
		default:
			return string(p), nil
		}
	}
	return "", nil
}

// ToAtlasSDK is clone of ToAtlas used temporarily for test migrations
func (p AtlasDatabaseUser) ToAtlasSDK(ctx context.Context, kubeClient client.Client) (*admin.CloudDatabaseUser, error) {
	password, err := p.ReadPassword(ctx, kubeClient)
	if err != nil {
		return nil, err
	}

	result := admin.NewCloudDatabaseUserWithDefaults()
	err = compat.JSONCopy(result, p.Spec)
	if err != nil {
		return nil, err
	}

	result.Password = &password

	return result, nil
}

func (p AtlasDatabaseUser) GetScopes(scopeType ScopeType) []string {
	var scopeClusters []string
	for _, scope := range p.Spec.Scopes {
		if scope.Type == scopeType {
			scopeClusters = append(scopeClusters, scope.Name)
		}
	}
	return scopeClusters
}

// ************************************ Builder methods *************************************************

func NewDBUser(namespace, name, dbUserName, projectName string) *AtlasDatabaseUser {
	return &AtlasDatabaseUser{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasDatabaseUserSpec{
			Username:       dbUserName,
			Project:        &common.ResourceRefNamespaced{Name: projectName},
			PasswordSecret: &common.ResourceRef{},
			Roles:          []RoleSpec{},
			Scopes:         []ScopeSpec{},
		},
	}
}

func (p *AtlasDatabaseUser) WithName(name string) *AtlasDatabaseUser {
	p.Name = name
	return p
}

func (p *AtlasDatabaseUser) WithAtlasUserName(name string) *AtlasDatabaseUser {
	p.Spec.Username = name
	return p
}

func (p *AtlasDatabaseUser) WithPasswordSecret(name string) *AtlasDatabaseUser {
	p.Spec.PasswordSecret.Name = name
	return p
}

func (p *AtlasDatabaseUser) WithRole(roleName, databaseName, collectionName string) *AtlasDatabaseUser {
	p.Spec.Roles = append(p.Spec.Roles, RoleSpec{RoleName: roleName, DatabaseName: databaseName, CollectionName: collectionName})
	return p
}

func (p *AtlasDatabaseUser) WithScope(scopeType ScopeType, name string) *AtlasDatabaseUser {
	p.Spec.Scopes = append(p.Spec.Scopes, ScopeSpec{Name: name, Type: scopeType})
	return p
}

func (p *AtlasDatabaseUser) ClearScopes() *AtlasDatabaseUser {
	p.Spec.Scopes = make([]ScopeSpec, 0)
	return p
}

func (p *AtlasDatabaseUser) WithDeleteAfterDate(date string) *AtlasDatabaseUser {
	p.Spec.DeleteAfterDate = date
	return p
}

func (p *AtlasDatabaseUser) WithExternaLProject(projectID, credentialsName string) *AtlasDatabaseUser {
	p.Spec.Project = nil
	p.Spec.ExternalProjectRef = &ExternalProjectReference{
		ID: projectID,
	}
	p.Spec.LocalCredentialHolder = api.LocalCredentialHolder{
		ConnectionSecret: &api.LocalObjectReference{
			Name: credentialsName,
		},
	}

	return p
}

func (p AtlasDatabaseUser) Credentials() *api.LocalObjectReference {
	return p.Spec.Credentials()
}

func DefaultDBUser(namespace, username, projectName string) *AtlasDatabaseUser {
	return NewDBUser(namespace, username, username, projectName).WithRole("clusterMonitor", "admin", "")
}
