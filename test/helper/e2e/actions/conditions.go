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

package actions

import (
	"fmt"
	"slices"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	v1 "k8s.io/api/core/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

func WaitForConditionsToBecomeTrue(userData *model.TestDataProvider, conditionTypes ...api.ConditionType) {
	Eventually(allConditionsAreTrueFunc(userData, conditionTypes...)).
		WithTimeout(15*time.Minute).WithPolling(20*time.Second).
		Should(BeTrue(), fmt.Sprintf("Status conditions %v are not all 'True'", conditionTypes))
}

// CheckProjectConditionsNotSet wait for Ready condition to become true and checks that input conditions are unset
func CheckProjectConditionsNotSet(userData *model.TestDataProvider, conditionTypes ...api.ConditionType) {
	Eventually(conditionsAreUnset(userData, conditionTypes...)).
		WithTimeout(15*time.Minute).WithPolling(20*time.Second).
		Should(BeTrue(), fmt.Sprintf("Status conditions %v should be unset. project status: %v",
			conditionTypes, userData.Project.Status.Conditions))
}

func allConditionsAreTrueFunc(userData *model.TestDataProvider, conditionTypes ...api.ConditionType) func(g types.Gomega) bool {
	return func(g Gomega) bool {
		conditions, err := kube.GetAllProjectConditions(userData)
		g.Expect(err).ShouldNot(HaveOccurred())

		for _, conditionType := range conditionTypes {
			foundTrue := false
			for _, condition := range conditions {
				if condition.Type == conditionType && condition.Status == v1.ConditionTrue {
					foundTrue = true
					break
				}
			}

			if !foundTrue {
				return false
			}
		}

		return true
	}
}

func conditionsAreUnset(userData *model.TestDataProvider, unsetConditionTypes ...api.ConditionType) func(g types.Gomega) bool {
	return func(g Gomega) bool {
		conditions, err := kube.GetAllProjectConditions(userData)
		g.Expect(err).ShouldNot(HaveOccurred())

		isReady := false
		for _, condition := range conditions {
			if condition.Type == api.ReadyType && condition.Status == v1.ConditionTrue {
				isReady = true
				break
			}
		}

		if !isReady {
			return false
		}

		for _, condition := range conditions {
			if slices.Contains(unsetConditionTypes, condition.Type) {
				return false
			}
		}

		return true
	}
}
