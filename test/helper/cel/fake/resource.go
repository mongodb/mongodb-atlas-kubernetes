// Package fake contains API Schema definitions for the resource v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=test.mongodb.com
package fake

// Regenerate with:
// controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./test/helper/cel/fake/..."
// controller-gen "crd:crdVersions=v1,ignoreUnexportedFields=true" rbac:roleName=manager-role webhook paths="./test/helper/cel/fake/..." output:crd:artifacts:config=test/helper/cel/fake
// make fmt

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "test.mongodb.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

// ResourceSpec defines the desired state of Resource
// +kubebuilder:validation:XValidation:rule=!has(self.deprecatedSet) || has(oldSelf.deprecatedSet), message="setting new deprecated set values is invalid: use the NewThing CRD instead."
type ResourceSpec struct {
	// DeprecatedField for the resource
	DeprecatedSet []string `json:"deprecatedSet,omitempty"`
}

// ResourceStatus defines the observed state of FakeRemote
type ResourceStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Resource is the Schema for the resource API
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec,omitempty"`
	Status ResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceList contains a list of Resource
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resource{}, &ResourceList{})
}
