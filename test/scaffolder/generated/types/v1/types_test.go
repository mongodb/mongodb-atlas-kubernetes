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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestGetConditions verifies the generated GetConditions helper methods.
func TestGetConditions(t *testing.T) {
	t.Run("Parent with conditions", func(t *testing.T) {
		conditions := []metav1.Condition{
			{Type: "Ready", Status: metav1.ConditionTrue, Reason: "Settled"},
		}
		parent := &Parent{
			Status: ParentStatus{
				Conditions: &conditions,
			},
		}

		got := parent.GetConditions()
		require.Len(t, got, 1)
		assert.Equal(t, "Ready", got[0].Type)
	})

	t.Run("Parent without conditions", func(t *testing.T) {
		parent := &Parent{}
		got := parent.GetConditions()
		assert.Nil(t, got)
	})

	t.Run("Child with conditions", func(t *testing.T) {
		conditions := []metav1.Condition{
			{Type: "Ready", Status: metav1.ConditionFalse, Reason: "Pending"},
		}
		child := &Child{
			Status: ChildStatus{
				Conditions: &conditions,
			},
		}

		got := child.GetConditions()
		require.Len(t, got, 1)
		assert.Equal(t, "Pending", got[0].Reason)
	})

	t.Run("Child without conditions", func(t *testing.T) {
		child := &Child{}
		got := child.GetConditions()
		assert.Nil(t, got)
	})
}
