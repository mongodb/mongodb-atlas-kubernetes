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

package model

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	helper "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/api/aws"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/k8s"
)

const (
	defaultE2EObjectProtectionDeletion    = false
	defaultE2ESubObjectProtectionDeletion = false
)

// TestDataProvider is the Full Data set for the current test case
type TestDataProvider struct {
	ConfPaths                   []string                  // init deployments configuration
	ConfUpdatePaths             []string                  // update configuration
	Resources                   UserInputs                // struct of all user resources project, deployments, database users
	Actions                     []func(*TestDataProvider) // additional actions for the current data set
	PortGroup                   int                       // ports for the test application starts from _
	SkipAppConnectivityCheck    bool
	Context                     context.Context
	K8SClient                   client.Client
	InitialDeployments          []*akov2.AtlasDeployment
	Project                     *akov2.AtlasProject
	Prefix                      string
	Users                       []*akov2.AtlasDatabaseUser
	Teams                       []*akov2.AtlasTeam
	ManagerContext              context.Context
	AWSResourcesGenerator       *helper.AwsResourcesGenerator
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

func DataProviderWithResources(ctx context.Context, keyTestPrefix string, project AProject, r *AtlasKeyType, initDeploymentConfigs []string, updateDeploymentConfig []string, users []DBUser, portGroup int, actions []func(*TestDataProvider)) TestDataProvider {
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

	data.AWSResourcesGenerator = helper.NewAwsResourcesGenerator(ctx, GinkgoT(), nil)

	data.ObjectDeletionProtection = defaultE2EObjectProtectionDeletion
	data.SubObjectDeletionProtection = defaultE2ESubObjectProtectionDeletion

	return data
}

func DataProvider(ctx context.Context, keyTestPrefix string, r *AtlasKeyType, portGroup int, actions []func(*TestDataProvider)) *TestDataProvider {
	var data TestDataProvider
	data.Prefix = keyTestPrefix
	data.Resources = NewSimpleUserInputs(keyTestPrefix, r)
	data.Actions = actions
	data.PortGroup = portGroup
	data.Context = context.Background()
	k8sClient, err := kubecli.CreateNewClient()
	Expect(err).NotTo(HaveOccurred(), "Failed to create k8s client")
	data.K8SClient = k8sClient

	data.AWSResourcesGenerator = helper.NewAwsResourcesGenerator(ctx, GinkgoT(), nil)

	data.ObjectDeletionProtection = defaultE2EObjectProtectionDeletion
	data.SubObjectDeletionProtection = defaultE2ESubObjectProtectionDeletion

	return &data
}

func (data TestDataProvider) WithInitialDeployments(deployments ...*akov2.AtlasDeployment) *TestDataProvider {
	data.InitialDeployments = append(data.InitialDeployments, deployments...)
	return &data
}

func (data TestDataProvider) WithProject(project *akov2.AtlasProject) *TestDataProvider {
	project.Spec.Name = fmt.Sprintf("%s-%s", data.Prefix, project.Spec.Name)
	data.Project = project
	return &data
}

func (data TestDataProvider) WithUsers(users ...*akov2.AtlasDatabaseUser) *TestDataProvider {
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
