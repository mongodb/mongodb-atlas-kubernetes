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

// +k8s:deepcopy-gen=false

type Reader interface {
	// GetStatus returns the status of the object.
	GetStatus() Status
}

// +k8s:deepcopy-gen=false

type Writer interface {
	// UpdateStatus allows to do the update of the status of an Atlas Custom resource.
	UpdateStatus(conditions []Condition, option ...Option)
}

// +k8s:deepcopy-gen=false

// Status is a generic status for any Custom Resource managed by Atlas Operator
type Status interface {
	GetConditions() []Condition

	GetObservedGeneration() int64
}

var _ Status = &Common{}

// Common is the struct shared by all statuses in existing Custom Resources.
type Common struct {
	// Conditions is the list of statuses showing the current state of the Atlas Custom Resource
	Conditions []Condition `json:"conditions"`

	// ObservedGeneration indicates the generation of the resource specification of which the Atlas Operator is aware.
	// The Atlas Operator updates this field to the value of 'metadata.generation' as soon as it starts reconciliation of the resource.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func (c Common) GetConditions() []Condition {
	return c.Conditions
}

func (c Common) GetObservedGeneration() int64 {
	return c.ObservedGeneration
}
