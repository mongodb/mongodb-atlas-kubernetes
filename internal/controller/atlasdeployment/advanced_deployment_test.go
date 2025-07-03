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
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
)

func TestHandleAdvancedDeployment(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment    *akov2.AtlasDeployment
		deploymentInAtlas  *deployment.Cluster
		deploymentService  func() deployment.AtlasDeploymentsService
		sdkMock            func() *admin.APIClient
		expectedResult     ctrl.Result
		expectedConditions []api.Condition
	}{
		"fail to create a new cluster in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Cluster")).
					Return(nil, errors.New("failed to create cluster"))

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotCreatedInAtlas)).
					WithMessageRegexp("failed to create cluster"),
			},
		},
		"create a new cluster in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: nil,
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().CreateDeployment(context.Background(), mock.AnythingOfType("*deployment.Cluster")).
					Return(
						&deployment.Cluster{
							ProjectID: "project-id",
							State:     "CREATING",
							AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
								Name:        "cluster0",
								ClusterType: "REPLICASET",
								ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
									{
										RegionConfigs: []*akov2.AdvancedRegionConfig{
											{
												ProviderName: "AWS",
												RegionName:   "US_WEST_1",
												Priority:     pointer.MakePtr(7),
												ElectableSpecs: &akov2.Specs{
													InstanceSize: "M10",
													NodeCount:    pointer.MakePtr(3),
												},
											},
										},
									},
								},
								BackupEnabled:            pointer.MakePtr(false),
								EncryptionAtRestProvider: "NONE",
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
		"fail to update a cluster in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "IDLE",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Cluster")).
					Return(nil, errors.New("failed to update cluster"))

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentNotUpdatedInAtlas)).
					WithMessageRegexp("failed to update cluster"),
			},
		},
		"update a cluster in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "IDLE",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().UpdateDeployment(context.Background(), mock.AnythingOfType("*deployment.Cluster")).
					Return(
						&deployment.Cluster{
							ProjectID: "project-id",
							State:     "UPDATING",
							AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
								Name:        "cluster0",
								ClusterType: "REPLICASET",
								ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
									{
										RegionConfigs: []*akov2.AdvancedRegionConfig{
											{
												ProviderName: "AWS",
												RegionName:   "US_WEST_1",
												Priority:     pointer.MakePtr(7),
												ElectableSpecs: &akov2.Specs{
													InstanceSize: "M20",
													NodeCount:    pointer.MakePtr(3),
												},
											},
										},
									},
								},
								BackupEnabled:            pointer.MakePtr(false),
								EncryptionAtRestProvider: "NONE",
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
		"cluster is updating in atlas": { //nolint:dupl
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "UPDATING",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M20",
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
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
					WithReason(string(workflow.DeploymentUpdating)).
					WithMessageRegexp("deployment is updating"),
			},
		},
		"cluster was deleted in atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "DELETING",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M20",
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				return translation.NewAtlasDeploymentsServiceMock(t)
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult:     ctrl.Result{},
			expectedConditions: nil,
		},
		"cluster has an unknown state in atlas": { //nolint:dupl
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M20",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
				},
			},
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "LOST",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:        "cluster0",
					ClusterType: "REPLICASET",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M20",
										NodeCount:    pointer.MakePtr(3),
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
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
					WithMessageRegexp("unknown deployment state: LOST"),
			},
		},
		"fail to update a cluster process args": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
					ProcessArgs: &akov2.ProcessArgs{
						JavascriptEnabled:         pointer.MakePtr(true),
						MinimumEnabledTLSProtocol: "TLS1_2",
						DefaultReadConcern:        "available",
					},
				},
			},
			//nolint:dupl
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "IDLE",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					EncryptionAtRestProvider: "NONE",
					MongoDBMajorVersion:      "7.0",
					VersionReleaseSystem:     "LTS",
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
										},
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ClusterWithProcessArgs(context.Background(), mock.AnythingOfType("*deployment.Cluster")).
					Return(errors.New("failed to get process args"))

				return service
			},
			sdkMock: func() *admin.APIClient {
				return &admin.APIClient{}
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DeploymentReadyType).
					WithReason(string(workflow.DeploymentAdvancedOptionsReady)).
					WithMessageRegexp("failed to get process args"),
			},
		},
		"update a cluster process args": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name:        "cluster0",
						ClusterType: "REPLICASET",
						ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
							{
								RegionConfigs: []*akov2.AdvancedRegionConfig{
									{
										ProviderName: "AWS",
										RegionName:   "US_WEST_1",
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &akov2.Specs{
											InstanceSize: "M10",
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
					ProcessArgs: &akov2.ProcessArgs{
						JavascriptEnabled:         pointer.MakePtr(true),
						MinimumEnabledTLSProtocol: "TLS1_2",
						DefaultReadConcern:        "available",
					},
				},
			},
			//nolint:dupl
			deploymentInAtlas: &deployment.Cluster{
				ProjectID: "project-id",
				State:     "IDLE",
				AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
					Name:                     "cluster0",
					ClusterType:              "REPLICASET",
					BackupEnabled:            pointer.MakePtr(false),
					EncryptionAtRestProvider: "NONE",
					MongoDBMajorVersion:      "7.0",
					VersionReleaseSystem:     "LTS",
					Paused:                   pointer.MakePtr(false),
					PitEnabled:               pointer.MakePtr(false),
					RootCertType:             "ISRGROOTX1",
					ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
						{
							ZoneName:  "Zone 1",
							NumShards: 1,
							RegionConfigs: []*akov2.AdvancedRegionConfig{
								{
									ProviderName: "AWS",
									RegionName:   "US_WEST_1",
									Priority:     pointer.MakePtr(7),
									ElectableSpecs: &akov2.Specs{
										InstanceSize: "M10",
										NodeCount:    pointer.MakePtr(3),
									},
									AutoScaling: &akov2.AdvancedAutoScalingSpec{
										Compute: &akov2.ComputeSpec{
											Enabled: pointer.MakePtr(false),
										},
										DiskGB: &akov2.DiskGB{
											Enabled: pointer.MakePtr(false),
										},
									},
								},
							},
						},
					},
					Tags: []*akov2.TagSpec{},
				},
			},
			deploymentService: func() deployment.AtlasDeploymentsService {
				service := translation.NewAtlasDeploymentsServiceMock(t)
				service.EXPECT().ClusterWithProcessArgs(context.Background(), mock.Anything).
					RunAndReturn(func(_ context.Context, cluster *deployment.Cluster) error {
						cluster.ProcessArgs = &akov2.ProcessArgs{
							JavascriptEnabled:         pointer.MakePtr(true),
							MinimumEnabledTLSProtocol: "LTS1_2",
							NoTableScan:               pointer.MakePtr(false),
							DefaultReadConcern:        "available",
						}
						return nil
					})
				service.EXPECT().UpdateProcessArgs(context.Background(), mock.Anything).
					RunAndReturn(func(_ context.Context, cluster *deployment.Cluster) error {
						cluster.ProcessArgs = &akov2.ProcessArgs{
							JavascriptEnabled:         pointer.MakePtr(true),
							MinimumEnabledTLSProtocol: "LTS1_2",
							NoTableScan:               pointer.MakePtr(false),
						}
						return nil
					})
				service.EXPECT().GetDeployment(context.Background(), mock.Anything, mock.Anything).
					Return(
						&deployment.Cluster{
							ProjectID: "project-id",
							State:     "UPDATING",
							AdvancedDeploymentSpec: &akov2.AdvancedDeploymentSpec{
								Name:        "cluster0",
								ClusterType: "REPLICASET",
								ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
									{
										RegionConfigs: []*akov2.AdvancedRegionConfig{
											{
												ProviderName: "AWS",
												RegionName:   "US_WEST_1",
												Priority:     pointer.MakePtr(7),
												ElectableSpecs: &akov2.Specs{
													InstanceSize: "M10",
													NodeCount:    pointer.MakePtr(3),
												},
											},
										},
									},
								},
								BackupEnabled:            pointer.MakePtr(false),
								EncryptionAtRestProvider: "NONE",
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
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312002: tt.sdkMock(),
				},
			}

			deploymentInAKO := deployment.NewDeployment("project-id", tt.atlasDeployment).(*deployment.Cluster)
			var projectService project.ProjectService // nil projetc service
			result, err := reconciler.handleAdvancedDeployment(ctx, projectService, tt.deploymentService(), deploymentInAKO, tt.deploymentInAtlas)
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
