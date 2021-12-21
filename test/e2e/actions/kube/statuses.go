package kube

import (
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func GetReadyProjectStatus(data *model.TestDataProvider) func() string {
	return func() string {
		return kubecli.GetStatusCondition("Ready", data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
	}
}

func GetProjectPEndpointStatus(data *model.TestDataProvider) func() string {
	return func() string {
		return kubecli.GetStatusCondition("PrivateEndpointReady", data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
	}
}

func GetProjectPEndpointServiceStatus(data *model.TestDataProvider) func() string {
	return func() string {
		return kubecli.GetStatusCondition("PrivateEndpointServiceReady", data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
	}
}
