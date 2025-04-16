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

package kube

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

func ProjectReadyCondition(data *model.TestDataProvider) string {
	condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, api.ReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
	return condition
}

func DeploymentReadyCondition(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetDeploymentStatusCondition(data.Context, data.K8SClient, api.ReadyType, data.Resources.Namespace, data.InitialDeployments[0].ObjectMeta.GetName())
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

func GetProjectStatusCondition(data *model.TestDataProvider, statusType api.ConditionType) (string, error) {
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

func GetAllProjectConditions(data *model.TestDataProvider) (result []api.Condition, err error) {
	err = data.K8SClient.Get(data.Context, types.NamespacedName{Name: data.Project.Name, Namespace: data.Project.Namespace}, data.Project)
	if err != nil {
		return result, err
	}

	return data.Project.Status.Conditions, nil
}
