package actions

import (
	"fmt"
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
			for _, unsetConditionType := range unsetConditionTypes {
				if condition.Type == unsetConditionType {
					return false
				}
			}
		}

		return true
	}
}
