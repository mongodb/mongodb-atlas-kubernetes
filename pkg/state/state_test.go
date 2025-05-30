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

package state_test

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

func TestGetState(t *testing.T) {
	tests := []struct {
		name      string
		conds     []metav1.Condition
		wantState state.ResourceState
	}{
		{
			name:      "no conditions returns initial",
			conds:     nil,
			wantState: state.StateInitial,
		},
		{
			name:      "empty conditions returns initial",
			conds:     []metav1.Condition{},
			wantState: state.StateInitial,
		},
		{
			name: "unrelated condition returns initial",
			conds: []metav1.Condition{
				{Type: "Other", Reason: "Ignored"},
			},
			wantState: state.StateInitial,
		},
		{
			name: "state condition returns correct state",
			conds: []metav1.Condition{
				{Type: state.StateCondition, Reason: string(state.StateCreated)},
			},
			wantState: state.StateCreated,
		},
		{
			name: "multiple conditions picks state condition",
			conds: []metav1.Condition{
				{Type: "Other", Reason: "Whatever"},
				{Type: state.StateCondition, Reason: string(state.StateDeleted)},
			},
			wantState: state.StateDeleted,
		},
		{
			name: "state condition with empty reason returns empty string state",
			conds: []metav1.Condition{
				{Type: state.StateCondition, Reason: ""},
			},
			wantState: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := state.GetState(tc.conds)
			if got != tc.wantState {
				t.Errorf("GetState() = %v, want %v", got, tc.wantState)
			}
		})
	}
}

func TestEnsureState(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name              string
		status            bool
		expectedCondition metav1.ConditionStatus
		state             state.ResourceState
		msg               string
	}{
		{
			name:              "sets ConditionTrue",
			status:            true,
			expectedCondition: metav1.ConditionTrue,
			state:             state.StateImported,
			msg:               "Import successful",
		},
		{
			name:              "sets ConditionFalse",
			status:            false,
			expectedCondition: metav1.ConditionFalse,
			state:             state.StateDeleting,
			msg:               "Deletion in progress",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var conds []metav1.Condition
			observedGen := int64(123)
			state.EnsureState(&conds, observedGen, tc.state, tc.msg, tc.status)
			if len(conds) != 1 {
				t.Fatalf("expected 1 condition, got %d", len(conds))
			}
			got := conds[0]
			if got.Type != state.StateCondition {
				t.Errorf("Condition Type = %v, want %v", got.Type, state.StateCondition)
			}
			if got.Status != tc.expectedCondition {
				t.Errorf("Condition Status = %v, want %v", got.Status, tc.expectedCondition)
			}
			if got.Reason != string(tc.state) {
				t.Errorf("Condition Reason = %v, want %v", got.Reason, tc.state)
			}
			if got.Message != tc.msg {
				t.Errorf("Condition Message = %v, want %v", got.Message, tc.msg)
			}
			if got.ObservedGeneration != observedGen {
				t.Errorf("ObservedGeneration = %v, want %v", got.ObservedGeneration, observedGen)
			}
			if got.LastTransitionTime.Time.Before(now) {
				t.Errorf("LastTransitionTime = %v, should be >= test start time", got.LastTransitionTime)
			}
		})
	}
}
