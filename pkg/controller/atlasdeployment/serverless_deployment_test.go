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
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestHandleServerlessInstance(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment    *akov2.AtlasDeployment
		deploymentInAtlas  *deployment.Serverless
		deploymentService  func() deployment.AtlasDeploymentsService
		sdkMock            func() *admin.APIClient
		expectedResult     ctrl.Result
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ServerlessPrivateEndpointReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("unable to retrieve list of serverless private endpoints from Atlas: failed to list private endpoints"),
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.ServerlessPrivateEndpointFailed)).
					WithMessageRegexp("unable to retrieve list of serverless private endpoints from Atlas: failed to list private endpoints"),
			},
		},
		"serverless instance is updating when private endpoints are in progress": {
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
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
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.DeploymentReadyType),
				api.TrueCondition(api.ReadyType),
			},
		},
		"serverless instance has an unknown state": {
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
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("unknown deployment state: NEW_UNKNOWN_STATE"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasDeployment).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			reconciler := &AtlasDeploymentReconciler{
				Client:            k8sClient,
				Log:               logger,
				deploymentService: tt.deploymentService(),
			}
			ctx := &workflow.Context{
				Context:   context.Background(),
				Log:       logger,
				SdkClient: tt.sdkMock(),
			}

			deploymentInAKO := deployment.NewDeployment("project-id", tt.atlasDeployment).(*deployment.Serverless)
			result, err := reconciler.handleServerlessInstance(ctx, deploymentInAKO, tt.deploymentInAtlas)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					ctx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}
