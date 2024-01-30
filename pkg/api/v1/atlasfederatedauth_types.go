package v1

import (
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasFederatedAuth{}, &AtlasFederatedAuthList{})
}

type AtlasFederatedAuthSpec struct {
	// +kubebuilder:default:=false
	Enabled bool `json:"enabled,omitempty"`
	// Connection secret with API credentials for configuring the federation.
	// These credentials must have OrganizationOwner permissions.
	ConnectionSecretRef common.ResourceRefNamespaced `json:"connectionSecretRef,omitempty"`
	// Approved domains that restrict users who can join the organization based on their email address.
	// +optional
	DomainAllowList []string `json:"domainAllowList,omitempty"`
	// Prevent users in the federation from accessing organizations outside of the federation, and creating new organizations.
	// This option applies to the entire federation.
	// See more information at https://www.mongodb.com/docs/atlas/security/federation-advanced-options/#restrict-user-membership-to-the-federation
	// +kubebuilder:default:=false
	DomainRestrictionEnabled *bool `json:"domainRestrictionEnabled,omitempty"`
	// +kubebuilder:default:=false
	// +optional
	SSODebugEnabled *bool `json:"ssoDebugEnabled,omitempty"`
	// Atlas roles that are granted to a user in this organization after authenticating.
	// +optional
	PostAuthRoleGrants []string `json:"postAuthRoleGrants,omitempty"`
	// Map IDP groups to Atlas roles.
	// +optional
	RoleMappings []RoleMapping `json:"roleMappings,omitempty"`
}

func (f *AtlasFederatedAuthSpec) ToAtlas(orgID, idpID string, projectNameToID map[string]string) (*admin.ConnectedOrgConfig, error) {
	var errs []error
	atlasRoleMappings := make([]admin.AuthFederationRoleMapping, 0, len(f.RoleMappings))

	for i := range f.RoleMappings {
		roleMapping := &f.RoleMappings[i]
		atlasRoleAssignments := make([]admin.RoleAssignment, 0, len(roleMapping.RoleAssignments))
		for j := range roleMapping.RoleAssignments {
			atlasRoleAssignment := admin.RoleAssignment{}
			roleAssignment := &roleMapping.RoleAssignments[j]
			if roleAssignment.ProjectName != "" {
				id, ok := projectNameToID[roleAssignment.ProjectName]
				if !ok {
					errs = append(errs, fmt.Errorf("project name '%s' doesn't exists in the organization", roleAssignment.ProjectName))
					continue
				}
				atlasRoleAssignment.SetGroupId(id)
			} else {
				atlasRoleAssignment.SetOrgId(orgID)
			}
			atlasRoleAssignment.SetRole(roleAssignment.Role)
			atlasRoleAssignments = append(atlasRoleAssignments, atlasRoleAssignment)
		}
		atlasRoleMappings = append(atlasRoleMappings, admin.AuthFederationRoleMapping{
			ExternalGroupName: roleMapping.ExternalGroupName,
			Id:                &idpID,
			RoleAssignments:   &atlasRoleAssignments,
		})
	}

	result := &admin.ConnectedOrgConfig{
		DomainAllowList:          &f.DomainAllowList,
		DomainRestrictionEnabled: *f.DomainRestrictionEnabled,
		IdentityProviderId:       idpID,
		OrgId:                    orgID,
		PostAuthRoleGrants:       &f.PostAuthRoleGrants,
	}

	if len(atlasRoleMappings) > 0 {
		result.SetRoleMappings(atlasRoleMappings)
	}

	return result, errors.Join(errs...)
}

// RoleMapping maps an external group from an identity provider to roles within Atlas.
type RoleMapping struct {
	// ExternalGroupName is the name of the IDP group to which this mapping applies.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=200
	ExternalGroupName string `json:"externalGroupName,omitempty"`
	// RoleAssignments define the roles within projects that should be given to members of the group.
	RoleAssignments []RoleAssignment `json:"roleAssignments,omitempty"`
}

type RoleAssignment struct {
	// The Atlas project in the same org in which the role should be given.
	ProjectName string `json:"projectName,omitempty"`
	// The role in Atlas that should be given to group members.
	// +kubebuilder:validation:Enum=ORG_MEMBER;ORG_READ_ONLY;ORG_BILLING_ADMIN;ORG_GROUP_CREATOR;ORG_OWNER;ORG_BILLING_READ_ONLY;ORG_TEAM_MEMBERS_ADMIN;GROUP_AUTOMATION_ADMIN;GROUP_BACKUP_ADMIN;GROUP_MONITORING_ADMIN;GROUP_OWNER;GROUP_READ_ONLY;GROUP_USER_ADMIN;GROUP_BILLING_ADMIN;GROUP_DATA_ACCESS_ADMIN;GROUP_DATA_ACCESS_READ_ONLY;GROUP_DATA_ACCESS_READ_WRITE;GROUP_CHARTS_ADMIN;GROUP_CLUSTER_MANAGER;GROUP_SEARCH_INDEX_EDITOR
	Role string `json:"role,omitempty"`
}

// AtlasFederatedAuth is the Schema for the Atlasfederatedauth API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type AtlasFederatedAuth struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasFederatedAuthSpec          `json:"spec,omitempty"`
	Status status.AtlasFederatedAuthStatus `json:"status,omitempty"`
}

func (f *AtlasFederatedAuth) ConnectionSecretObjectKey() *client.ObjectKey {
	var key client.ObjectKey
	if f.Spec.ConnectionSecretRef.Namespace != "" {
		key = kube.ObjectKey(f.Spec.ConnectionSecretRef.Namespace, f.Spec.ConnectionSecretRef.Name)
	} else {
		key = kube.ObjectKey(f.Namespace, f.Spec.ConnectionSecretRef.Name)
	}
	return &key
}

func (f *AtlasFederatedAuth) GetStatus() status.Status {
	return f.Status
}

func (f *AtlasFederatedAuth) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	f.Status.Conditions = conditions
	f.Status.ObservedGeneration = f.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasFederatedAuthStatusOption)
		v(&f.Status)
	}
}

// AtlasFederatedAuthList contains a list of AtlasFederatedAuth
// +kubebuilder:object:root=true
type AtlasFederatedAuthList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasFederatedAuth `json:"items"`
}
