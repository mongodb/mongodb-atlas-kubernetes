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

package atlasdeployment

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	"go.mongodb.org/atlas-sdk/v20250312010/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func TestHandleServerlessInstance(t *testing.T) {
	type workflowRes struct {
		res ctrl.Result
		err error
	}
	tests := map[string]struct {
		atlasDeployment    *akov2.AtlasDeployment
		deploymentInAtlas  *deployment.Serverless
		deploymentService  func() deployment.AtlasDeploymentsService
		sdkMock            func() *admin.APIClient
		expectedResult     workflowRes
		expectedConditions []api.Condition
	}{
		"fail to create a new serverless instance in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: true,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Serverless")).
					Return(nil, errors.New("failed to create serverless instance"))

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				err: errors.New("failed to create serverless instance"),
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotCreatedInAtlas)).
					WithMessageRegexp("failed to create serverless instance"),
			},
		},
		"create a new serverless instance in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: true,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Serverless")).
					Return(
						&deployment.Serverless{
							ProjectID: "project-id",
							State:     "CREATING",
							ServerlessSpec: &akov2.ServerlessSpec{
								Name: "instance0",
								ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
									ProviderName:        provider.ProviderServerless,
									BackingProviderName: "AWS",
									RegionName:          "us-east-1",
								},
								Tags: []*akov2.TagSpec{
									{
										Key:   "test",
										Value: "e2e",
									},
								},
								BackupOptions: akov2.ServerlessBackupOptions{
									ServerlessContinuousBackupEnabled: true,
								},
								TerminationProtectionEnabled: true,
							},
						},
						nil,
					)

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentCreating)).
					WithMessageRegexp("deployment is provisioning"),
			},
		},
		"fail to update a serverless instance in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: true,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Serverless")).
					Return(nil, errors.New("failed to update serverless instance"))

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				err: errors.New("failed to update serverless instance"),
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
					WithMessageRegexp("failed to update serverless instance"),
			},
		},
		"update a serverless instance in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: true,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Serverless")).
					Return(
						&deployment.Serverless{
							ProjectID: "project-id",
							State:     "UPDATING",
							ServerlessSpec: &akov2.ServerlessSpec{
								Name: "instance0",
								ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
									ProviderName:        provider.ProviderServerless,
									BackingProviderName: "AWS",
									RegionName:          "us-east-1",
								},
								Tags: []*akov2.TagSpec{
									{
										Key:   "test",
										Value: "e2e",
									},
								},
								BackupOptions: akov2.ServerlessBackupOptions{
									ServerlessContinuousBackupEnabled: false,
								},
								TerminationProtectionEnabled: true,
							},
						},
						nil,
					)

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"serverless instance is updating in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "UPDATING",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"update tags of a serverless instance in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
							{
								Key:   "newTag",
								Value: "newValue",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Serverless")).
					Return(
						&deployment.Serverless{
							ProjectID: "project-id",
							State:     "UPDATING",
							ServerlessSpec: &akov2.ServerlessSpec{
								Name: "instance0",
								ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
									ProviderName:        provider.ProviderServerless,
									BackingProviderName: "AWS",
									RegionName:          "us-east-1",
								},
								Tags: []*akov2.TagSpec{
									{
										Key:   "test",
										Value: "e2e",
									},
									{
										Key:   "newTag",
										Value: "newValue",
									},
								},
								BackupOptions: akov2.ServerlessBackupOptions{
									ServerlessContinuousBackupEnabled: false,
								},
								TerminationProtectionEnabled: true,
							},
						},
						nil,
					)

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"serverless instance fails when private endpoints fails": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
						PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
							{
								Name:                    "spe1",
								CloudProviderEndpointID: "arn-12345",
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				return service
			},
			sdkMock: func() *admin.APIClient {
				speClient := mockadmin.NewServerlessPrivateEndpointsApi(t)
				speClient.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance0").
					Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speClient})
				speClient.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
					Return(nil, &http.Response{}, errors.New("failed to list private endpoints"))

				return &admin.APIClient{ServerlessPrivateEndpointsApi: speClient}
			},
			expectedResult: workflowRes{
				err: errors.New("unable to retrieve list of serverless private endpoints from Atlas: failed to list private endpoints"),
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ServerlessPrivateEndpointReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("unable to retrieve list of serverless private endpoints from Atlas: failed to list private endpoints"),
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("unable to retrieve list of serverless private endpoints from Atlas: failed to list private endpoints"),
			},
		},
		"serverless flex instance fails when private endpoints are set": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
						PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
							{
								Name:                    "spe1",
								CloudProviderEndpointID: "arn-12345",
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				return service
			},
			sdkMock: func() *admin.APIClient {
				mockError := &admin.GenericOpenAPIError{}
				model := *admin.NewApiErrorWithDefaults()
				model.SetErrorCode("NOT_SERVERLESS_TENANT_CLUSTER")
				mockError.SetModel(model)

				speClient := mockadmin.NewServerlessPrivateEndpointsApi(t)
				speClient.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance0").
					Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speClient})
				speClient.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
					Return(nil, &http.Response{}, mockError)

				return &admin.APIClient{ServerlessPrivateEndpointsApi: speClient}
			},
			expectedResult: workflowRes{
				err: errors.New("serverless private endpoints are not supported: "),
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ServerlessPrivateEndpointReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("serverless private endpoints are not supported: "),
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("serverless private endpoints are not supported: "),
			},
		},
		"serverless instance is updating when private endpoints are in progress": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
						PrivateEndpoints: []akov2.ServerlessPrivateEndpoint{
							{
								Name:                    "spe1",
								CloudProviderEndpointID: "arn-12345",
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				return service
			},
			sdkMock: func() *admin.APIClient {
				speClient := mockadmin.NewServerlessPrivateEndpointsApi(t)
				speClient.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance0").
					Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speClient})
				speClient.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
					Return(nil, &http.Response{}, nil)

				speClient.EXPECT().CreateServerlessPrivateEndpoint(context.Background(), "project-id", "instance0", mock.AnythingOfType("*admin.ServerlessTenantCreateRequest")).
					Return(admin.CreateServerlessPrivateEndpointApiRequest{ApiService: speClient})
				speClient.EXPECT().CreateServerlessPrivateEndpointExecute(mock.AnythingOfType("admin.CreateServerlessPrivateEndpointApiRequest")).
					Return(&admin.ServerlessTenantEndpoint{}, &http.Response{}, nil)

				return &admin.APIClient{ServerlessPrivateEndpointsApi: speClient}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ServerlessPrivateEndpointReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointInProgress)).
					WithMessageRegexp("Waiting serverless private endpoint to be configured"),
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointInProgress)).
					WithMessageRegexp("Waiting serverless private endpoint to be configured"),
			},
		},
		"serverless instance is ready in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "IDLE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)

				return service
			},
			sdkMock: func() *admin.APIClient {
				speClient := mockadmin.NewServerlessPrivateEndpointsApi(t)
				speClient.EXPECT().ListServerlessPrivateEndpoints(context.Background(), "project-id", "instance0").
					Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speClient})
				speClient.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
					Return(nil, &http.Response{}, nil)

				return &admin.APIClient{ServerlessPrivateEndpointsApi: speClient}
			},
			expectedResult: workflowRes{
				res: ctrl.Result{},
				err: nil,
			},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.DeploymentReadyType),
				api.TrueCondition(api.ReadyType).
					WithMessageRegexp("WARNING: Serverless is deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details."),
			},
		},
		"serverless instance has an unknown state": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "instance0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "instance0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        provider.ProviderServerless,
							BackingProviderName: "AWS",
							RegionName:          "us-east-1",
						},
						Tags: []*akov2.TagSpec{
							{
								Key:   "test",
								Value: "e2e",
							},
						},
						BackupOptions: akov2.ServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: false,
						},
						TerminationProtectionEnabled: true,
					},
				},
			},
			deploymentInAtlas: &deployment.Serverless{
				ProjectID: "project-id",
				State:     "NEW_UNKNOWN_STATE",
				ServerlessSpec: &akov2.ServerlessSpec{
					Name: "instance0",
					ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
						ProviderName:        provider.ProviderServerless,
						BackingProviderName: "AWS",
						RegionName:          "us-east-1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "e2e",
						},
					},
					BackupOptions: akov2.ServerlessBackupOptions{
						ServerlessContinuousBackupEnabled: false,
					},
					TerminationProtectionEnabled: true,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: workflowRes{
				err: errors.New("unknown deployment state: NEW_UNKNOWN_STATE"),
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("unknown deployment state: NEW_UNKNOWN_STATE"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			logger := zaptest.NewLogger(t)
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			dbUserProjectIndexer := indexer.NewAtlasDatabaseUserByProjectIndexer(ctx, nil, logger)
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasDeployment).
				WithIndex(dbUserProjectIndexer.Object(), dbUserProjectIndexer.Name(), dbUserProjectIndexer.Keys).
				Build()
			reconciler := &AtlasDeploymentReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger.Sugar(),
				},
			}
			workflowCtx := &workflow.Context{
				Context: ctx,
				Log:     logger.Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312006: tt.sdkMock(),
				},
			}

			deploymentInAKO := deployment.NewDeployment("project-id", tt.atlasDeployment).(*deployment.Serverless)
			var projectService project.ProjectService
			result, err := reconciler.handleServerlessInstance(workflowCtx, projectService, tt.deploymentService(), deploymentInAKO, tt.deploymentInAtlas)
			//require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, workflowRes{
				res: result,
				err: err,
			})
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					workflowCtx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}
