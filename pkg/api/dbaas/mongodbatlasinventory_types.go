/*
Copyright 2021.
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
	"sigs.k8s.io/controller-runtime/pkg/client"

	dbaasv1beta1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1beta1"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +groupName:=dbaas.redhat.com
// +versionName:=v1beta1

// MongoDBAtlasInventory is the Schema for the MongoDBAtlasInventory API
type MongoDBAtlasInventory struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   dbaasv1beta1.DBaaSInventorySpec   `json:"spec,omitempty"`
	Status dbaasv1beta1.DBaaSInventoryStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBAtlasInventoryList contains a list of DBaaSInventories
type MongoDBAtlasInventoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDBAtlasInventory `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDBAtlasInventory{}, &MongoDBAtlasInventoryList{})
}

func (p *MongoDBAtlasInventory) ConnectionSecretObjectKey() *client.ObjectKey {
	if p.Spec.CredentialsRef != nil {
		key := kube.ObjectKey(p.Namespace, p.Spec.CredentialsRef.Name)
		return &key
	}
	return nil
}

// +k8s:deepcopy-gen=false

// SetInventoryCondition sets a condition on the status object. If the condition already
// exists, it will be replaced. SetCondition does not update the resource in
// the cluster.
func SetInventoryCondition(inv *MongoDBAtlasInventory, condType string, status metav1.ConditionStatus, reason, msg string) {
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

// GetInventoryCondition return the condition with the passed condition type from
// the status object. If the condition is not already present, return nil
func GetInventoryCondition(inv *MongoDBAtlasInventory, condType string) *metav1.Condition {
	for i := range inv.Status.Conditions {
		if inv.Status.Conditions[i].Type == condType {
			return &inv.Status.Conditions[i]
		}
	}
	return nil
}
