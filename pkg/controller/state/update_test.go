//Copyright 2025 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package state

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// statusObj embeds an Unstructured to satisfy metav1.Object and
// implements GetConditions() to satisfy StatusObject.
type statusObj struct {
	*unstructured.Unstructured
	conditions []metav1.Condition
}

func (s *statusObj) GetConditions() []metav1.Condition { return s.conditions }

func newStatusObj(gen int64, conditions []metav1.Condition) *statusObj {
	u := &unstructured.Unstructured{}
	u.SetGeneration(gen)
	return &statusObj{Unstructured: u, conditions: conditions}
}

func TestShouldUpdate(t *testing.T) {
	now := time.Now()
	pastMillis := strconv.FormatInt(now.Add(-2*time.Hour).UnixMilli(), 10)

	tests := []struct {
		name         string
		obj          metav1.Object
		shouldUpdate bool
		wantErr      string
	}{
		{
			name: "generation changed",
			obj: newStatusObj(2, []metav1.Condition{
				{
					Type:               state.StateCondition,
					ObservedGeneration: 1,
				},
			}),
			shouldUpdate: true,
		},
		{
			name: "generation did not change",
			obj: newStatusObj(1, []metav1.Condition{
				{
					Type:               state.StateCondition,
					ObservedGeneration: 1,
				},
			}),
			shouldUpdate: false,
		},
		{
			name: "error status (ready reason error)",
			obj: newStatusObj(1, []metav1.Condition{
				{
					Type:   state.ReadyCondition,
					Reason: ReadyReasonError,
				},
			}),
			shouldUpdate: true,
		},
		{
			name: "reapply due (old timestamp + period)",
			obj: newUnstructuredObj(map[string]string{
				AnnotationReapplyTimestamp:   pastMillis,
				"mongodb.com/reapply-period": "1h",
			}),
			shouldUpdate: true,
		},
		{
			name:         "no update needed",
			obj:          newUnstructuredObj(map[string]string{}),
			shouldUpdate: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ShouldUpdate(tc.obj)
			if tc.wantErr != "" {
				assert.ErrorContains(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.shouldUpdate, got)
		})
	}
}
