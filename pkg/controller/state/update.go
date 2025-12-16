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
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
)

// ShouldUpdate returns true if the object should be updated based on generation change, reapply period, or error status.
// Note: a generation change or error status will only be detected if the object implements StatusObject interface.
//
// Returns an error if there is an issue checking the reapply period.
func ShouldUpdate(obj metav1.Object) (bool, error) {
	generationChanged, hasErrorState := false, false

	if statusObj, ok := obj.(StatusObject); ok {
		if stateCondition := meta.FindStatusCondition(statusObj.GetConditions(), state.StateCondition); stateCondition != nil {
			generationChanged = stateCondition.ObservedGeneration != obj.GetGeneration()
		}
		if errorCondition := meta.FindStatusCondition(statusObj.GetConditions(), state.ReadyCondition); errorCondition != nil {
			hasErrorState = errorCondition.Reason == ReadyReasonError
		}
	}

	shouldReapply, err := ShouldReapply(obj)
	if err != nil {
		return false, fmt.Errorf("failed to check reapply period: %w", err)
	}

	return generationChanged || shouldReapply || hasErrorState, nil
}
