package actions

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	v1 "k8s.io/api/core/v1"

	. "github.com/onsi/gomega"

	kube "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/actions/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func WaitForConditionsToBecomeTrue(userData *model.TestDataProvider, conditonTypes ...status.ConditionType) {
	allConditionsAreTrueFunc := func(g Gomega) bool {
		conditions, err := kube.GetAllProjectConditions(userData)
		g.Expect(err).ShouldNot(HaveOccurred())

		for _, conditionType := range conditonTypes {
			found := false
			for _, condition := range conditions {
				if condition.Type == conditionType && condition.Status == v1.ConditionTrue {
					found = true
					break
				}
			}

			if !found {
				return false
			}
		}

		return true
	}

	Eventually(allConditionsAreTrueFunc).Should(BeTrue(), fmt.Sprintf("Status conditions %v are not 'True'", conditonTypes))
}
