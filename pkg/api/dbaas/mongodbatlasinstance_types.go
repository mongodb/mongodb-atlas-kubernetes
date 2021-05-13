/*
Copyright 2022.
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

package dbaas

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	dbaasv1beta1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1beta1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=dbaas.redhat.com
// +versionName:=v1beta1

// MongoDBAtlasInstance is the Schema for the MongoDBAtlasInstance API
type MongoDBAtlasInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbaasv1beta1.DBaaSInstanceSpec   `json:"spec,omitempty"`
	Status dbaasv1beta1.DBaaSInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBAtlasInstanceList contains a list of DBaaSInstances
type MongoDBAtlasInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBAtlasInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBAtlasInstance{}, &MongoDBAtlasInstanceList{})
}

// +k8s:deepcopy-gen=false

// SetInstanceCondition sets a condition on the status object. If the condition already
// exists, it will be replaced. SetCondition does not update the resource in
// the cluster.
func SetInstanceCondition(inv *MongoDBAtlasInstance, condType string, status metav1.ConditionStatus, reason, msg string) {
	now := metav1.Now()
	for i := range inv.Status.Conditions {
		if inv.Status.Conditions[i].Type == condType {
			var lastTransitionTime metav1.Time
			if inv.Status.Conditions[i].Status != status {
				lastTransitionTime = now
			} else {
				lastTransitionTime = inv.Status.Conditions[i].LastTransitionTime
			}
			inv.Status.Conditions[i] = metav1.Condition{
				LastTransitionTime: lastTransitionTime,
				Status:             status,
				Type:               condType,
				Reason:             reason,
				Message:            msg,
			}
			return
		}
	}

	// If the condition does not exist,
	// initialize the lastTransitionTime
	inv.Status.Conditions = append(inv.Status.Conditions, metav1.Condition{
		LastTransitionTime: now,
		Type:               condType,
		Status:             status,
		Reason:             reason,
		Message:            msg,
	})
}

// GetInstanceCondition return the condition with the passed condition type from
// the status object. If the condition is not already present, return nil
func GetInstanceCondition(inv *MongoDBAtlasInstance, condType string) *metav1.Condition {
	for i := range inv.Status.Conditions {
		if inv.Status.Conditions[i].Type == condType {
			return &inv.Status.Conditions[i]
		}
	}
	return nil
}
