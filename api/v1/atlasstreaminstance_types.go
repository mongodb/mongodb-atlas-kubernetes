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
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

type AtlasStreamInstanceSpec struct {
	// Human-readable label that identifies the stream connection.
	Name string `json:"name"`
	// The configuration to be used to connect to an Atlas Cluster.
	Config Config `json:"clusterConfig"`
	// Project which the instance belongs to.
	Project common.ResourceRefNamespaced `json:"projectRef"`
	// List of connections of the stream instance for the specified project.
	ConnectionRegistry []common.ResourceRefNamespaced `json:"connectionRegistry,omitempty"`
}

type Config struct {
	// Name of the cluster configured for this connection.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT;SERVERLESS
	// +kubebuilder:default=AWS
	Provider string `json:"provider"`
	// Name of the cloud provider region hosting Atlas Stream Processing.
	Region string `json:"region"`
	// Selected tier for the Stream Instance. Configures Memory / VCPU allowances.
	// +kubebuilder:validation:Enum=SP10;SP30;SP50
	// +kubebuilder:default=SP10
	Tier string `json:"tier"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Atlas ID",type=string,JSONPath=`.status.id`
// +kubebuilder:resource:categories=atlas,shortName=asi

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

func (f *AtlasStreamInstance) GetStatus() api.Status {
	return f.Status
}

func (f *AtlasStreamInstance) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	f.Status.Conditions = conditions
	f.Status.ObservedGeneration = f.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasStreamInstanceStatusOption)
		v(&f.Status)
	}
}

func (f *AtlasStreamInstance) AtlasProjectObjectKey() client.ObjectKey {
	ns := f.Namespace
	if f.Spec.Project.Namespace != "" {
		ns = f.Spec.Project.Namespace
	}

	return kube.ObjectKey(ns, f.Spec.Project.Name)
}
