package atlasdeployment

import (
	"context"
	"errors"
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleFlexInstance(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment    *akov2.AtlasDeployment
		deploymentInAtlas  *deployment.Flex
		deploymentService  func() deployment.AtlasDeploymentsService
		expectedResult     ctrl.Result
		expectedConditions []api.Condition
	}{
		"fail to create flex cluster in atlas": {
			atlasDeployment:   basicFlexCluster(),
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Flex")).
					Return(nil, errors.New("failed to create flex cluster"))

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotCreatedInAtlas)).
					WithMessageRegexp("failed to create flex cluster"),
			},
		},
		"create a flex cluster in atlas": {
			atlasDeployment:   basicFlexCluster(),
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Flex")).
					Return(
						&deployment.Flex{
							ProjectID: "project-id",
							State:     "CREATING",
							FlexSpec: &akov2.FlexSpec{
								Name: "cluster0",
								ProviderSettings: &akov2.FlexProviderSettings{
									BackingProviderName: "AWS",
									RegionName:          "US_EAST_1",
								},
								TerminationProtectionEnabled: false,
							},
						},
						nil,
					)

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentCreating)).
					WithMessageRegexp("deployment is provisioning"),
			},
		},
		"fail to update a flex cluster in atlas": {
			atlasDeployment: basicFlexCluster(),
			deploymentInAtlas: &deployment.Flex{
				ProjectID: "project-id",
				State:     "IDLE",
				FlexSpec: &akov2.FlexSpec{
					Name: "cluster0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_2",
					},
					TerminationProtectionEnabled: false,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Flex")).
					Return(nil, errors.New("failed to update flex cluster"))

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
					WithMessageRegexp("failed to update flex cluster"),
			},
		},
		"update a flex cluster in atlas": {
			atlasDeployment: basicFlexCluster(),
			deploymentInAtlas: &deployment.Flex{
				ProjectID: "project-id",
				State:     "IDLE",
				FlexSpec: &akov2.FlexSpec{
					Name: "cluster0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					Tags: []*akov2.TagSpec{
						{
							Key:   "test",
							Value: "value",
						},
					},
					TerminationProtectionEnabled: false,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Flex")).
					Return(
						&deployment.Flex{
							ProjectID: "project-id",
							State:     "UPDATING",
							FlexSpec: &akov2.FlexSpec{
								Name: "cluster0",
								ProviderSettings: &akov2.FlexProviderSettings{
									BackingProviderName: "AWS",
									RegionName:          "US_EAST_1",
								},
								Tags: []*akov2.TagSpec{
									{
										Key:   "test",
										Value: "value",
									},
								},
								TerminationProtectionEnabled: false,
							},
						},
						nil,
					)

				return service
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"flex cluster is updating": {
			atlasDeployment: basicFlexCluster(),
			deploymentInAtlas: &deployment.Flex{
				ProjectID: "project-id",
				State:     "UPDATING",
				FlexSpec: &akov2.FlexSpec{
					Name: "cluster0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					TerminationProtectionEnabled: false,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"flex cluster is deleting": {
			atlasDeployment: basicFlexCluster(),
			deploymentInAtlas: &deployment.Flex{
				ProjectID: "project-id",
				State:     "DELETING",
				FlexSpec: &akov2.FlexSpec{
					Name: "cluster0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					TerminationProtectionEnabled: false,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{},
		},
		"flex cluster has unknown state in atlas": {
			atlasDeployment: basicFlexCluster(),
			deploymentInAtlas: &deployment.Flex{
				ProjectID: "project-id",
				State:     "NONSENSE",
				FlexSpec: &akov2.FlexSpec{
					Name: "cluster0",
					ProviderSettings: &akov2.FlexProviderSettings{
						BackingProviderName: "AWS",
						RegionName:          "US_EAST_1",
					},
					TerminationProtectionEnabled: false,
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("unknown deployment state: NONSENSE"),
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
			}

			deploymentInAKO := deployment.NewDeployment("project-id", tt.atlasDeployment).(*deployment.Flex)
			var projectService project.ProjectService
			result, err := reconciler.handleFlexInstance(workflowCtx, projectService, tt.deploymentService(), deploymentInAKO, tt.deploymentInAtlas)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
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

func basicFlexCluster() *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster0",
			Namespace: "default",
		},
		Spec: akov2.AtlasDeploymentSpec{
			FlexSpec: &akov2.FlexSpec{
				Name: "cluster0",
				ProviderSettings: &akov2.FlexProviderSettings{
					BackingProviderName: "AWS",
					RegionName:          "US_EAST_1",
				},
				TerminationProtectionEnabled: false,
			},
		},
	}
}
