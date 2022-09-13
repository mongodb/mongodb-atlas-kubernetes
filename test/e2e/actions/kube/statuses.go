package kube

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func GetReadyProjectStatus(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.ReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		return condition
	}
}

func GetProjectPEndpointStatus(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.PrivateEndpointReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		return condition
	}
}

func GetProjectNetworkPeerStatus(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.NetworkPeerReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		return condition
	}
}

func GetProjectPEndpointServiceStatus(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.PrivateEndpointServiceReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		return condition
	}
}

func GetProjectCloudAccessRolesStatus(data *model.TestDataProvider) func() string {
	return func() string {
		condition, _ := kubecli.GetProjectStatusCondition(data.Context, data.K8SClient, status.CloudProviderAccessReadyType, data.Resources.Namespace, data.Resources.Project.ObjectMeta.GetName())
		return condition
	}
}
