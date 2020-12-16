package status

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*

conditions:
  - lastTransitionTime: "2020-12-15T20:46:55Z"
    status: "True"
    message: the following fields are missing in the Secret secret: %v
    reason: AtlasCredentialsNotProvided
    type: ProjectReady
  - lastTransitionTime: "2020-12-15T20:46:55Z"
    message: NOT_ALLOWED. You don't have enough permissions to perform the operation
    reason: AtlasApiError
    status: "False"
    type: IPAccessListReady
  - privateLink
  - lastTransitionTime: "2020-12-15T20:46:55Z"
    status: "True"
    type: Ready

*/
const (
	ReadyType ConditionType = "Ready"

	// AtlasProject condition types
	ProjectReadyType      ConditionType = "ProjectReady"
	IPAccessListReadyType ConditionType = "IPAccessListReady"
)

type ConditionType string

// Condition describes the state of an Atlas Custom Resource at a certain point.
type Condition struct {
	// Type of Atlas Custom Resource condition.
	Type ConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}
