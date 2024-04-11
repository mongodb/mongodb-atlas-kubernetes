package atlasdeployment

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestMergedAdvancedDeployment(t *testing.T) {
	defaultAtlas := makeDefaultAtlasSpec()
	atlasRegionConfig := defaultAtlas.ReplicationSpecs[0].RegionConfigs[0]
	fillInSpecs(atlasRegionConfig, "M10", "AWS")

	t.Run("Test merging clusters removes backing provider name if empty", func(t *testing.T) {
		advancedCluster := akov2.DefaultAwsAdvancedDeployment("default", "my-project")

		merged, _, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		assert.NoError(t, err)
		assert.Empty(t, merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})

	t.Run("Test merging clusters does not remove backing provider name if it is present in the atlas type", func(t *testing.T) {
		atlasRegionConfig.ElectableSpecs.InstanceSize = "M5"
		atlasRegionConfig.ProviderName = "TENANT"
		atlasRegionConfig.BackingProviderName = "AWS"

		advancedCluster := akov2.DefaultAwsAdvancedDeployment("default", "my-project")
		advancedCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "M5"
		advancedCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ProviderName = "TENANT"
		advancedCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName = "AWS"

		merged, _, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		assert.NoError(t, err)
		assert.Equal(t, atlasRegionConfig.BackingProviderName, merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})
}

func TestAdvancedDeploymentsEqual(t *testing.T) {
	defaultAtlas := makeDefaultAtlasSpec()
	regionConfig := defaultAtlas.ReplicationSpecs[0].RegionConfigs[0]
	fillInSpecs(regionConfig, "M10", "AWS")

	t.Run("Test equal advanced deployments", func(t *testing.T) {
		advancedCluster := akov2.DefaultAwsAdvancedDeployment("default", "my-project")

		merged, atlas, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		assert.NoError(t, err)
		beforeSpec := merged.DeepCopy()
		beforeAtlas := atlas.DeepCopy()

		logger, _ := zap.NewProduction()
		areEqual, _ := AdvancedDeploymentsEqual(logger.Sugar(), &merged, &atlas)
		assert.Equalf(t, merged, atlas, "Deployments should be equal")
		assert.True(t, areEqual, "Deployments should be equal")
		assert.Equal(t, beforeSpec, &merged, "Comparison should not change original spec values")
		assert.Equal(t, beforeAtlas, &atlas, "Comparison should not change original atlas values")
	})

	t.Run("Advanced deployments are equal when autoscaling is ON and only differ on instance sizes", func(t *testing.T) {
		advancedCluster := akov2.DefaultAwsAdvancedDeployment("default", "my-project")
		// set auto scaling ON
		advancedCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AutoScaling = &akov2.AdvancedAutoScalingSpec{
			DiskGB: &akov2.DiskGB{
				Enabled: pointer.MakePtr(false),
			},
			Compute: &akov2.ComputeSpec{
				Enabled:          pointer.MakePtr(true),
				ScaleDownEnabled: pointer.MakePtr(true),
				MinInstanceSize:  "M10",
				MaxInstanceSize:  "M30",
			},
		}

		merged, atlas, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		// copy autoscaling to atlas
		k8sRegion := advancedCluster.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0]
		atlas.ReplicationSpecs[0].RegionConfigs[0].AutoScaling = &akov2.AdvancedAutoScalingSpec{
			DiskGB: &akov2.DiskGB{
				Enabled: k8sRegion.AutoScaling.DiskGB.Enabled,
			},
			Compute: &akov2.ComputeSpec{
				Enabled:          k8sRegion.AutoScaling.Compute.Enabled,
				ScaleDownEnabled: k8sRegion.AutoScaling.Compute.ScaleDownEnabled,
				MinInstanceSize:  k8sRegion.AutoScaling.Compute.MinInstanceSize,
				MaxInstanceSize:  k8sRegion.AutoScaling.Compute.MaxInstanceSize,
			},
		}
		// inject difference
		atlas.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "something-else"
		assert.NoError(t, err)
		beforeSpec := merged.DeepCopy()
		beforeAtlas := atlas.DeepCopy()

		logger, _ := zap.NewProduction()
		areEqual, _ := AdvancedDeploymentsEqual(logger.Sugar(), &merged, &atlas)
		assert.True(t, areEqual, "Deployments should be equal")
		assert.Equal(t, beforeSpec, &merged, "Comparison should not change original spec values")
		assert.Equal(t, beforeAtlas, &atlas, "Comparison should not change original atlas values")
	})

	t.Run("Advanced deployments are different when autoscaling is OFF and only differ on instance sizes", func(t *testing.T) {
		advancedCluster := akov2.DefaultAwsAdvancedDeployment("default", "my-project")

		merged, atlas, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		// inject difference
		atlas.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "something-else"
		assert.NoError(t, err)
		beforeSpec := merged.DeepCopy()
		beforeAtlas := atlas.DeepCopy()

		logger, _ := zap.NewProduction()
		areEqual, _ := AdvancedDeploymentsEqual(logger.Sugar(), &merged, &atlas)
		assert.False(t, areEqual, "Deployments should be different")
		assert.Equal(t, beforeSpec, &merged, "Comparison should not change original spec values")
		assert.Equal(t, beforeAtlas, &atlas, "Comparison should not change original atlas values")
	})
}

func makeDefaultAtlasSpec() *mongodbatlas.AdvancedCluster {
	return &mongodbatlas.AdvancedCluster{
		ClusterType: "REPLICASET",
		Name:        "test-deployment-advanced",
		ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
			{
				NumShards: 1,
				ID:        "123",
				ZoneName:  "Zone 1",
				RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
					{
						ElectableSpecs: &mongodbatlas.Specs{
							InstanceSize: "M10",
							NodeCount:    pointer.MakePtr(3),
						},
						Priority:     pointer.MakePtr(7),
						ProviderName: "AWS",
						RegionName:   "US_EAST_1",
					},
				},
			},
		},
	}
}

func fillInSpecs(regionConfig *mongodbatlas.AdvancedRegionConfig, instanceSize string, provider string) {
	regionConfig.ProviderName = provider

	regionConfig.ElectableSpecs.InstanceSize = instanceSize
	regionConfig.AnalyticsSpecs = &mongodbatlas.Specs{
		InstanceSize: instanceSize,
		NodeCount:    pointer.MakePtr(0),
	}
	regionConfig.ReadOnlySpecs = &mongodbatlas.Specs{
		InstanceSize: instanceSize,
		NodeCount:    pointer.MakePtr(0),
	}
}

func TestDbUserBelongsToProjects(t *testing.T) {
	t.Run("Database User refer to a different project name", func(*testing.T) {
		dbUser := &akov2.AtlasDatabaseUser{
			Spec: akov2.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name: "project2",
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name: "project1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User is no", func(*testing.T) {
		dbUser := &akov2.AtlasDatabaseUser{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "ns-2",
			},
			Spec: akov2.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name: "project1",
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User refer to a project with same name but in another namespace", func(*testing.T) {
		dbUser := &akov2.AtlasDatabaseUser{
			Spec: akov2.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "project1",
					Namespace: "ns-2",
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User refer to a valid project and namespace", func(*testing.T) {
		dbUser := &akov2.AtlasDatabaseUser{
			Spec: akov2.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "project1",
					Namespace: "ns-1",
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.True(t, dbUserBelongsToProject(dbUser, project))
	})
}

func TestHandleSearchNodes(t *testing.T) {
	projectName := "test-project"
	projectID := "abc123"

	t.Run("get search nodes request errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)
		deployment.Spec.DeploymentSpec.SearchNodes = []akov2.SearchNode{
			{
				InstanceSize: "S80_LOWCPU_NVME",
				NodeCount:    3,
			},
		}

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				nil,
				&http.Response{},
				errors.New("get test error"),
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in AKO and atlas (update)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)
		deployment.Spec.DeploymentSpec.SearchNodes = []akov2.SearchNode{
			{
				InstanceSize: "S80_LOWCPU_NVME",
				NodeCount:    4,
			},
		}

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().UpdateAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name,
			&admin.ApiSearchDeploymentRequest{Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas()}).
			Return(admin.UpdateAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().UpdateAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S100_LOWCPU_NVME",
							NodeCount:    4,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
	})

	t.Run("search nodes are in AKO and atlas but update errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)
		deployment.Spec.DeploymentSpec.SearchNodes = []akov2.SearchNode{
			{
				InstanceSize: "S80_LOWCPU_NVME",
				NodeCount:    4,
			},
		}

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().UpdateAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name,
			&admin.ApiSearchDeploymentRequest{Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas()}).
			Return(admin.UpdateAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().UpdateAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				nil,
				&http.Response{},
				errors.New("update test error"),
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in AKO but not in Atlas (create)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)
		deployment.Spec.DeploymentSpec.SearchNodes = []akov2.SearchNode{
			{
				InstanceSize: "S80_LOWCPU_NVME",
				NodeCount:    3,
			},
		}

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				errors.New("no search nodes found"),
			)
		searchAPI.EXPECT().CreateAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name, mock.AnythingOfType("*admin.ApiSearchDeploymentRequest")).
			Return(admin.CreateAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().CreateAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S100_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
	})

	t.Run("search nodes are in AKO but not in Atlas but create errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)
		deployment.Spec.DeploymentSpec.SearchNodes = []akov2.SearchNode{
			{
				InstanceSize: "S80_LOWCPU_NVME",
				NodeCount:    3,
			},
		}

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				errors.New("no search nodes found"),
			)
		searchAPI.EXPECT().CreateAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name, mock.AnythingOfType("*admin.ApiSearchDeploymentRequest")).
			Return(admin.CreateAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().CreateAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{},
				errors.New("create test error"),
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in Atlas but not in AKO (delete)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().DeleteAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.DeleteAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().DeleteAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
	})

	t.Run("search nodes are in Atlas but not in AKO but delete errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().DeleteAtlasSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.DeleteAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().DeleteAtlasSearchDeploymentExecute(mock.Anything).
			Return(
				&http.Response{},
				errors.New("delete test error"),
			)

		ctx := &workflow.Context{
			SdkClient: &admin.APIClient{
				AtlasSearchApi: searchAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})
}
