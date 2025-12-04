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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasOrgSettings{})
	SchemeBuilder.Register(&AtlasOrgSettingsList{})
}

// AtlasOrgSettingsSpec defines the desired state of AtlasOrgSettings.
type AtlasOrgSettingsSpec struct {
	// OrgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects.
	// +required
	OrgID string `json:"orgID"`

	// ConnectionSecretRef is the name of the Kubernetes Secret which contains the information about the way to connect to Atlas (Public & Private API keys).
	ConnectionSecretRef *api.LocalObjectReference `json:"connectionSecretRef,omitempty"`

	// ApiAccessListRequired Flag that indicates whether to require API operations to originate from an IP Address added to the API access list for the specified organization.
	// +optional
	ApiAccessListRequired *bool `json:"apiAccessListRequired,omitempty"`

	// GenAIFeaturesEnabled Flag that indicates whether this organization has access to generative AI features. This setting only applies to Atlas Commercial and is enabled by default.
	// Once this setting is turned on, Project Owners may be able to enable or disable individual AI features at the project level.
	// +optional
	GenAIFeaturesEnabled *bool `json:"genAIFeaturesEnabled,omitempty"`

	// MaxServiceAccountSecretValidityInHours Number that represents the maximum period before expiry in hours for new Atlas Admin API Service Account secrets within the specified organization.
	// +optional
	MaxServiceAccountSecretValidityInHours *int `json:"maxServiceAccountSecretValidityInHours,omitempty"`

	// MultiFactorAuthRequired Flag that indicates whether to require users to set up Multi-Factor Authentication (MFA) before accessing the specified organization.
	// To learn more, see: https://www.mongodb.com/docs/atlas/security-multi-factor-authentication/.
	// +optional
	MultiFactorAuthRequired *bool `json:"multiFactorAuthRequired,omitempty"`

	// RestrictEmployeeAccess Flag that indicates whether to block MongoDB Support from accessing Atlas infrastructure and cluster logs for any deployment in the specified organization without explicit permission.
	// Once this setting is turned on, you can grant MongoDB Support a 24-hour bypass access to the Atlas deployment to resolve support issues.
	// To learn more, see: https://www.mongodb.com/docs/atlas/security-restrict-support-access/.
	// +optional
	RestrictEmployeeAccess *bool `json:"restrictEmployeeAccess,omitempty"`

	// SecurityContact String that specifies a single email address for the specified organization to receive security-related notifications.
	// Specifying a security contact does not grant them authorization or access to Atlas for security decisions or approvals.
	// An empty string is valid and clears the existing security contact (if any).
	// +optional
	SecurityContact *string `json:"securityContact,omitempty"`

	// StreamsCrossGroupEnabled Flag that indicates whether a group's Atlas Stream Processing instances in this organization can create connections to other group's clusters in the same organization.
	// +optional
	StreamsCrossGroupEnabled *bool `json:"streamsCrossGroupEnabled,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:name:plural=AtlasOrgSettings, singular=AtlasOrgSettings
// +kubebuilder:resource:categories=atlas,shortName=aos
// +kubebuilder:subresource:status
type AtlasOrgSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasOrgSettingsSpec          `json:"spec,omitempty"`
	Status status.AtlasOrgSettingsStatus `json:"status,omitempty"`
}

func (aos *AtlasOrgSettings) Credentials() *api.LocalObjectReference {
	return aos.Spec.ConnectionSecretRef
}

func (aos *AtlasOrgSettings) GetConditions() []metav1.Condition {
	if aos.Status.Conditions == nil {
		return []metav1.Condition{}
	}
	return aos.Status.Conditions
}

// +kubebuilder:object:root=true
type AtlasOrgSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasOrgSettings `json:"items"`
}
