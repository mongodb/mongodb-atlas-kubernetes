package model

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/test/helper"

	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/k8s"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

const (
	defaultE2EObjectProtectionDeletion    = false
	defaultE2ESubObjectProtectionDeletion = false
)

// Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths                   []string                  // init deployments configuration
	ConfUpdatePaths             []string                  // update configuration
	Resources                   UserInputs                // struct of all user resources project, deployments, database users
	Actions                     []func(*TestDataProvider) // additional actions for the current data set
	PortGroup                   int                       // ports for the test application starts from _
	SkipAppConnectivityCheck    bool
	Context                     context.Context
	K8SClient                   client.Client
	InitialDeployments          []*v1.AtlasDeployment
	Project                     *v1.AtlasProject
	Prefix                      string
	Users                       []*v1.AtlasDatabaseUser
	Teams                       []*v1.AtlasTeam
	ManagerContext              context.Context
	AWSResourcesGenerator       *helper.AwsResourcesGenerator
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

func DataProviderWithResources(keyTestPrefix string, project AProject, r *AtlasKeyType, initDeploymentConfigs []string, updateDeploymentConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
	var data TestDataProvider
	data.Prefix = keyTestPrefix
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
	Expect(err).NotTo(HaveOccurred(), "Failed to create k8s client")
	data.K8SClient = k8sClient

	data.AWSResourcesGenerator = helper.NewAwsResourcesGenerator(GinkgoT(), nil)

	data.ObjectDeletionProtection = defaultE2EObjectProtectionDeletion
	data.SubObjectDeletionProtection = defaultE2ESubObjectProtectionDeletion

	return data
}

func DataProvider(keyTestPrefix string, r *AtlasKeyType, portGroup int, actions []func(*TestDataProvider)) *TestDataProvider {
	var data TestDataProvider
	data.Prefix = keyTestPrefix
	data.Resources = NewSimpleUserInputs(keyTestPrefix, r)
	data.Actions = actions
	data.PortGroup = portGroup
	data.Context = context.Background()
	k8sClient, err := kubecli.CreateNewClient()
	Expect(err).NotTo(HaveOccurred(), "Failed to create k8s client")
	data.K8SClient = k8sClient

	data.AWSResourcesGenerator = helper.NewAwsResourcesGenerator(GinkgoT(), nil)

	data.ObjectDeletionProtection = defaultE2EObjectProtectionDeletion
	data.SubObjectDeletionProtection = defaultE2ESubObjectProtectionDeletion

	return &data
}

func (data TestDataProvider) WithInitialDeployments(deployments ...*v1.AtlasDeployment) *TestDataProvider {
	data.InitialDeployments = append(data.InitialDeployments, deployments...)
	return &data
}

func (data TestDataProvider) WithProject(project *v1.AtlasProject) *TestDataProvider {
	project.Spec.Name = fmt.Sprintf("%s-%s", data.Prefix, project.Spec.Name)
	data.Project = project
	return &data
}

func (data TestDataProvider) WithUsers(users ...*v1.AtlasDatabaseUser) *TestDataProvider {
	data.Users = append(data.Users, users...)
	return &data
}

func (data TestDataProvider) WithObjectDeletionProtection(protected bool) *TestDataProvider {
	data.ObjectDeletionProtection = protected
	return &data
}

func (data TestDataProvider) WithSubObjectDeletionProtection(protected bool) *TestDataProvider {
	data.SubObjectDeletionProtection = protected
	return &data
}
