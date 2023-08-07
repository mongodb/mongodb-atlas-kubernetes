//go:build e2e

package kube

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func ProjectReadyCondition(data *model.TestDataProvider) string {
	condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.ReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
	return condition
}

func DeploymentReadyCondition(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetDeploymentStatusCondition(data.Context, data.K8SClient, status.ReadyType, data.Resources.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
		return condition
	}
}

func GetDeploymentStatus(data *model.TestDataProvider) status.AtlasDeploymentStatus {
	err := data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.InitialDeployments[0].ObjectMeta.GetName(),
		Namespace: data.Resources.Namespace}, data.InitialDeployments[0])
	if err != nil {
		return status.AtlasDeploymentStatus{}
	}
	return data.InitialDeployments[0].Status
}

func GetProjectStatus(data *model.TestDataProvider) status.AtlasProjectStatus {
	err := data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.ObjectMeta.GetName(),
		Namespace: data.Resources.Namespace}, data.Project)
	if err != nil {
		return status.AtlasProjectStatus{}
	}
	return data.Project.Status
}

func GetProjectStatusCondition(data *model.TestDataProvider, statusType status.ConditionType) (string, error) {
	conditions, err := GetAllProjectConditions(data)
	if err != nil {
		return "", err
	}
	for _, condition := range conditions {
		if condition.Type == statusType {
			return string(condition.Status), err
		}
	}
	return "", fmt.Errorf("condition %s not found, conditions: %v", statusType, conditions)
}

func GetAllProjectConditions(data *model.TestDataProvider) (result []status.Condition, err error) {
	err = data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Project.Namespace}, data.Project)
	if err != nil {
		return result, err
	}

	return data.Project.Status.Conditions, nil
}
