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

package api

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

type ConditionType string

const (
	ReadyType           ConditionType = "Ready"
	ValidationSucceeded ConditionType = "ValidationSucceeded"
)

// AtlasProject condition types
const (
	ProjectReadyType                  ConditionType = "ProjectReady"
	IPAccessListReadyType             ConditionType = "IPAccessListReady"
	MaintenanceWindowReadyType        ConditionType = "MaintenanceWindowReady"
	PrivateEndpointServiceReadyType   ConditionType = "PrivateEndpointServiceReady"
	PrivateEndpointReadyType          ConditionType = "PrivateEndpointReady"
	NetworkPeerReadyType              ConditionType = "NetworkPeerReady"
	CloudProviderIntegrationReadyType ConditionType = "CloudProviderIntegrationReady"
	IntegrationReadyType              ConditionType = "ThirdPartyIntegrationReady"
	AlertConfigurationReadyType       ConditionType = "AlertConfigurationReady"
	EncryptionAtRestReadyType         ConditionType = "EncryptionAtRestReady"
	AuditingReadyType                 ConditionType = "AuditingReady"
	ProjectSettingsReadyType          ConditionType = "ProjectSettingsReady"
	ProjectCustomRolesReadyType       ConditionType = "ProjectCustomRolesReady"
	ProjectTeamsReadyType             ConditionType = "ProjectTeamsReady"
	SearchIndexesReadyType            ConditionType = "AtlasSearchIndexesReady"
	BackupComplianceReadyType         ConditionType = "BackupCompliancePolicyReady"
	X509AuthReadyType                 ConditionType = "X509AuthReady"
)

// AtlasDeployment condition types
const (
	DeploymentReadyType                ConditionType = "DeploymentReady"
	ServerlessPrivateEndpointReadyType ConditionType = "ServerlessPrivateEndpointReady"
	ManagedNamespacesReadyType         ConditionType = "ManagedNamespacesReady"
	CustomZoneMappingReadyType         ConditionType = "CustomZoneMappingReady"
	SearchNodesReadyType               ConditionType = "SearchNodesReady"
)

// AtlasDatabaseUser condition types
const (
	DatabaseUserReadyType ConditionType = "DatabaseUserReady"
)

// Atlas Data Federation condition types
const (
	DataFederationReadyType   ConditionType = "DataFederationReady"
	DataFederationPEReadyType ConditionType = "DataFederationPrivateEndpointsReady"
)

// Atlas Federated Auth condition types
const (
	FederatedAuthReadyType      ConditionType = "FederatedAuthReady"
	FederatedAuthRolesReadyType ConditionType = "RolesReady"
)

// Atlas Streams condition types
const (
	StreamInstanceReadyType   ConditionType = "StreamInstanceReady"
	StreamConnectionReadyType ConditionType = "StreamConnectionReady"
)

const (
	SearchIndexesNamesAreNotUnique = "SearchIndexesNamesAreNotUnique"
	SearchIndexesNotReady          = "SearchIndexesNotReady"
)

// Atlas Teams condition types
const (
	TeamUnmanaged ConditionType = "TeamUnmanaged"
)

// Atlas Private Endpoint condition types
const (
	PrivateEndpointServiceReady ConditionType = "PrivateEndpointServiceReady"
	PrivateEndpointReady        ConditionType = "PrivateEndpointReady"
)

// Atlas IP Access List condition types
const (
	IPAccessListReady ConditionType = "IPAccessListReady"
)

// Atlas Network Container condition types
const (
	NetworkContainerReady ConditionType = "NetworkContainerReady"
)

// Generic condition type
const (
	ResourceVersionStatus ConditionType = "ResourceVersionIsValid"
)

// Condition describes the state of an Atlas Custom Resource at a certain point.
type Condition struct {
	// Type of Atlas Custom Resource condition.
	Type ConditionType `json:"type"`
	// Status of the condition; one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// Last time the condition transitioned from one status to another.
	// Represented in ISO 8601 format.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// A message providing details about the transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// TrueCondition returns the Condition that has the 'Status' set to 'true' and 'Type' to 'conditionType'.
// It explicitly omits the 'Reason' and 'Message' fields.
func TrueCondition(conditionType ConditionType) Condition {
	return Condition{
		Type:               conditionType,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
	}
}

// FalseCondition returns the Condition that has the 'Status' set to 'false' and 'Type' to 'conditionType'.
// The reason and message can be provided optionally
func FalseCondition(conditionType ConditionType) Condition {
	condition := Condition{
		Type:               conditionType,
		Status:             corev1.ConditionFalse,
		LastTransitionTime: metav1.Now(),
	}
	return condition
}

func HasConditionType(typ ConditionType, source []Condition) bool {
	for i := range source {
		if source[i].Type == typ {
			return true
		}
	}
	return false
}

func HasReadyCondition(conditions []Condition) bool {
	for _, c := range conditions {
		if c.Type == ReadyType && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// EnsureConditionExists adds or updates the condition in the copy of a 'source' slice
func EnsureConditionExists(condition Condition, source []Condition) []Condition {
	condition.LastTransitionTime = metav1.Now()
	target := make([]Condition, len(source))
	copy(target, source)
	for i, c := range source {
		if c.Type == condition.Type {
			// We don't update the last transition time in case status hasn't changed.
			if c.Status == condition.Status {
				condition.LastTransitionTime = c.LastTransitionTime
			}
			//goland:noinspection GoNilness
			target[i] = condition
			return target
		}
	}
	// Condition not found - appending
	target = append(target, condition)
	return target
}

func RemoveConditionIfExists(conditionType ConditionType, source []Condition) []Condition {
	updatedConditions := []Condition{}
	for _, cond := range source {
		if cond.Type != conditionType {
			updatedConditions = append(updatedConditions, cond)
		}
	}
	return updatedConditions
}

func (c Condition) WithReason(reason string) Condition {
	c.Reason = reason
	return c
}

func (c Condition) WithMessageRegexp(msg string) Condition {
	c.Message = msg
	return c
}
