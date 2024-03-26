package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

type AtlasStreamInstanceSpec struct {
	// Human-readable label that identifies the stream connection
	Name string `json:"name"`
	// The configuration to be used to connect to a Atlas Cluster
	Config Config `json:"clusterConfig"`
	// Project which the instance belongs to
	Project common.ResourceRefNamespaced `json:"projectRef"`
	// List of connections of the stream instance for the specified project
	ConnectionRegistry []common.ResourceRefNamespaced `json:"connectionRegistry,omitempty"`
}

type Config struct {
	// Name of the cluster configured for this connection
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT;SERVERLESS
	// +kubebuilder:default=AWS
	Provider string `json:"provider"`
	// The name of a Built in or Custom DB Role to connect to an Atlas Cluster
	Region string `json:"region"`
	// Selected tier for the Stream Instance. Configures Memory / VCPU allowances.
	// +kubebuilder:validation:Enum=SP10;SP30;SP50
	// +kubebuilder:default=SP10
	Tier string `json:"tier"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`

// AtlasStreamInstance is the Schema for the atlasstreaminstances API
type AtlasStreamInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasStreamInstanceSpec          `json:"spec,omitempty"`
	Status status.AtlasStreamInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasStreamInstanceList contains a list of AtlasStreamInstance
type AtlasStreamInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasStreamInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasStreamInstance{}, &AtlasStreamInstanceList{})
}

func (f *AtlasStreamInstance) GetStatus() status.Status {
	return f.Status
}

func (f *AtlasStreamInstance) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	f.Status.Conditions = conditions
	f.Status.ObservedGeneration = f.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasStreamInstanceStatusOption)
		v(&f.Status)
	}
}
