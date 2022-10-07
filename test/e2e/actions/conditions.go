package actions

import (
	"fmt"
	"time"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	v1 "k8s.io/api/core/v1"

	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/types"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func WaitForConditionsToBecomeTrue(userData *model.TestDataProvider, conditonTypes ...status.ConditionType) {
	Eventually(allConditionsAreTrueFunc(userData, conditonTypes...)).
		WithTimeout(15*time.Minute).WithPolling(20*time.Second).
		Should(BeTrue(), fmt.Sprintf("Status conditions %v are not all 'True'", conditonTypes))
}

func allConditionsAreTrueFunc(userData *model.TestDataProvider, conditonTypes ...status.ConditionType) func(g types.Gomega) bool {
	return func(g Gomega) bool {
		conditions, err := kube.GetAllProjectConditions(userData)
		g.Expect(err).ShouldNot(HaveOccurred())

		for _, conditionType := range conditonTypes {
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
