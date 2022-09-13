package model

import (
	"context"
	"fmt"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths                []string                  // init deployments configuration
	ConfUpdatePaths          []string                  // update configuration
	Resources                UserInputs                // struct of all user resoucers project,deployments,databaseusers
	Actions                  []func(*TestDataProvider) // additional actions for the current data set
	PortGroup                int                       // ports for the test application starts from _
	SkipAppConnectivityCheck bool
	Context                  context.Context
	K8SClient                client.Client
}

func NewTestDataProvider(keyTestPrefix string, project AProject, r *AtlasKeyType, initDeploymentConfigs []string, updateDeploymentConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
	var data TestDataProvider
	data.ConfPaths = initDeploymentConfigs
	data.ConfUpdatePaths = updateDeploymentConfig
	data.Resources = NewUserInputs(keyTestPrefix, project, users, r)
	data.Actions = actions
	data.PortGroup = portGroup
	for i := range initDeploymentConfigs {
		data.Resources.Deployments = append(data.Resources.Deployments, LoadUserDeploymentConfig(data.ConfPaths[i]))
	}
	data.Context = context.Background()
	k8sClient, err := kubecli.CreateNewClient()
	if err != nil {
		panic(fmt.Sprintf("failed to create k8s client: %v", err))
	}
	data.K8SClient = k8sClient
	return data
}
