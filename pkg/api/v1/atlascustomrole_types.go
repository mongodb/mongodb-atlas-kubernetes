package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

func init() {
	SchemeBuilder.Register(&AtlasCustomRole{}, &AtlasCustomRoleList{})
}

// AtlasCustomRole is the Schema for the AtlasCustomRole API
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.role.name`
// +kubebuilder:printcolumn:name="Project ID",type=string,JSONPath=`.spec.projectIDRef.id`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=acr
type AtlasCustomRole struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasCustomRoleSpec          `json:"spec,omitempty"`
	Status status.AtlasCustomRoleStatus `json:"status,omitempty"`
}

func (in *AtlasCustomRole) Credentials() *api.LocalObjectReference {
	return in.Spec.ConnectionSecret
}

func (in *AtlasCustomRole) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	in.Status.Conditions = conditions
	in.Status.ObservedGeneration = in.ObjectMeta.Generation
}

func (in *AtlasCustomRole) GetStatus() api.Status {
	return in.Status
}

// AtlasCustomRoleList contains a list of AtlasCustomRole
// +kubebuilder:object:root=true
type AtlasCustomRoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasCustomRole `json:"items"`
}

// AtlasCustomRoleSpec defines the desired state of CustomRole in Atlas
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
type AtlasCustomRoleSpec struct {
	// ProjectReference is the dual external or kubernetes reference with access credentials
	ProjectDualReference `json:",inline"`
	Role                 CustomRole `json:"role"`
}
