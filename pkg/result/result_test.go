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

package result

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

var defaultRequeueResult = reconcile.Result{RequeueAfter: 15 * time.Second}

func TestNextState(t *testing.T) {
	tests := []struct {
		name        string
		state       state.ResourceState
		msg         string
		expected    ctrlstate.Result
		expectedErr error
	}{
		{
			name:  "StateInitial",
			state: state.StateInitial,
			msg:   "",
			expected: ctrlstate.Result{
				NextState: state.StateInitial,
				StateMsg:  "",
			},
		},
		{
			name:  "StateCreated",
			state: state.StateCreated,
			msg:   "Resource created",
			expected: ctrlstate.Result{
				NextState: state.StateCreated,
				StateMsg:  "Resource created.",
			},
		},
		{
			name:  "StateCreating with requeue",
			state: state.StateCreating,
			msg:   "Creating resource",
			expected: ctrlstate.Result{
				Result:    defaultRequeueResult,
				NextState: state.StateCreating,
				StateMsg:  "Creating resource.",
			},
		},
		{
			name:  "StateUpdated",
			state: state.StateUpdated,
			msg:   "Resource updated",
			expected: ctrlstate.Result{
				NextState: state.StateUpdated,
				StateMsg:  "Resource updated.",
			},
		},
		{
			name:  "StateUpdating with requeue",
			state: state.StateUpdating,
			msg:   "Updating resource",
			expected: ctrlstate.Result{
				Result:    defaultRequeueResult,
				NextState: state.StateUpdating,
				StateMsg:  "Updating resource.",
			},
		},
		{
			name:  "StateDeleted",
			state: state.StateDeleted,
			msg:   "Resource deleted",
			expected: ctrlstate.Result{
				NextState: state.StateDeleted,
				StateMsg:  "Resource deleted.",
			},
		},
		{
			name:  "StateDeletionRequested",
			state: state.StateDeletionRequested,
			msg:   "Resource delete",
			expected: ctrlstate.Result{
				Result:    defaultRequeueResult,
				NextState: state.StateDeletionRequested,
				StateMsg:  "Resource delete.",
			},
		},
		{
			name:  "StateDeleting",
			state: state.StateDeleting,
			msg:   "Deleting resource",
			expected: ctrlstate.Result{
				Result:    defaultRequeueResult,
				NextState: state.StateDeleting,
				StateMsg:  "Deleting resource.",
			},
		},
		{
			name:  "StateImported",
			state: state.StateImported,
			msg:   "Resource imported",
			expected: ctrlstate.Result{
				NextState: state.StateImported,
				StateMsg:  "Resource imported.",
			},
		},
		{
			name:  "StateImportRequested",
			state: state.StateImportRequested,
			msg:   "Resource import",
			expected: ctrlstate.Result{
				NextState: state.StateImportRequested,
				StateMsg:  "Resource import.",
			},
		},
		{
			name:        "Unknown state",
			state:       state.ResourceState("Unknown"),
			msg:         "Unknown state",
			expectedErr: fmt.Errorf("unknown state %v", state.ResourceState("Unknown")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NextState(tt.state, tt.msg)
			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestError(t *testing.T) {
	err := fmt.Errorf("an error occurred")
	st := state.StateCreating

	s, returnedErr := Error(st, err)

	require.Equal(t, ctrlstate.Result{
		Result: reconcile.Result{
			Requeue:      false,
			RequeueAfter: 0,
		},
		NextState: st,
		StateMsg:  "",
	}, s)
	require.EqualError(t, returnedErr, err.Error())
}
