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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	v1status "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasProject{}, &AtlasProjectList{})
}

var _ api.AtlasCustomResource = &AtlasProject{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.experimental.mongodb.com

// AtlasProject is the Schema for the atlasprojects API
type AtlasProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasProjectSpec            `json:"spec,omitempty"`
	Status v1status.AtlasProjectStatus `json:"status,omitempty"`
}

type AtlasProjectSpec struct {
	v1.AtlasProjectSpec `json:",inline"`

	// AuditRef references an existing audit configuration
	AuditRef common.ResourceRefNamespaced `json:"auditRef,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasProjectList contains a list of AtlasProject
type AtlasProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasProject `json:"items"`
}

func (p *AtlasProject) GetStatus() api.Status {
	return p.Status
}

func (p *AtlasProject) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	p.Status.Conditions = conditions
	p.Status.ObservedGeneration = p.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(v1status.AtlasProjectStatusOption)
		v(&p.Status)
	}
}

func (p *AtlasProject) ConnectionSecretObjectKey() *client.ObjectKey {
	if p.Spec.ConnectionSecret != nil {
		var key client.ObjectKey
		if p.Spec.ConnectionSecret.Namespace != "" {
			key = kube.ObjectKey(p.Spec.ConnectionSecret.Namespace, p.Spec.ConnectionSecret.Name)
		} else {
			key = kube.ObjectKey(p.Namespace, p.Spec.ConnectionSecret.Name)
		}
		return &key
	}
	return nil
}
