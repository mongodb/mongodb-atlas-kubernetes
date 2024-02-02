package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestMergedAdvancedDeployment(t *testing.T) {
	defaultAtlas := makeDefaultAtlasSpec()
	atlasRegionConfig := defaultAtlas.ReplicationSpecs[0].RegionConfigs[0]
	fillInSpecs(atlasRegionConfig, "M10", "AWS")

	t.Run("Test merging clusters removes backing provider name if empty", func(t *testing.T) {
		advancedCluster := mdbv1.DefaultAwsAdvancedDeployment("default", "my-project")

		merged, _, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		assert.NoError(t, err)
		assert.Empty(t, merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})

	t.Run("Test merging clusters does not remove backing provider name if it is present in the atlas type", func(t *testing.T) {
		atlasRegionConfig.ElectableSpecs.InstanceSize = "M5"
		atlasRegionConfig.ProviderName = "TENANT"
		atlasRegionConfig.BackingProviderName = "AWS"

		advancedCluster := mdbv1.DefaultAwsAdvancedDeployment("default", "my-project")
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

	t.Run("Test ", func(t *testing.T) {
		advancedCluster := mdbv1.DefaultAwsAdvancedDeployment("default", "my-project")

		merged, atlas, err := MergedAdvancedDeployment(*defaultAtlas, *advancedCluster.Spec.DeploymentSpec)
		assert.NoError(t, err)

		logger, _ := zap.NewProduction()
		areEqual, _ := AdvancedDeploymentsEqual(logger.Sugar(), merged, atlas)
		assert.Equalf(t, merged, atlas, "Deployments should be equal")
		assert.True(t, areEqual, "Deployments should be equal")
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
		dbUser := &mdbv1.AtlasDatabaseUser{
			Spec: mdbv1.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name: "project2",
				},
			},
		}
		project := &mdbv1.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name: "project1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User is no", func(*testing.T) {
		dbUser := &mdbv1.AtlasDatabaseUser{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "ns-2",
			},
			Spec: mdbv1.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name: "project1",
				},
			},
		}
		project := &mdbv1.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User refer to a project with same name but in another namespace", func(*testing.T) {
		dbUser := &mdbv1.AtlasDatabaseUser{
			Spec: mdbv1.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "project1",
					Namespace: "ns-2",
				},
			},
		}
		project := &mdbv1.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.False(t, dbUserBelongsToProject(dbUser, project))
	})

	t.Run("Database User refer to a valid project and namespace", func(*testing.T) {
		dbUser := &mdbv1.AtlasDatabaseUser{
			Spec: mdbv1.AtlasDatabaseUserSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "project1",
					Namespace: "ns-1",
				},
			},
		}
		project := &mdbv1.AtlasProject{
			ObjectMeta: v1.ObjectMeta{
				Name:      "project1",
				Namespace: "ns-1",
			},
		}

		assert.True(t, dbUserBelongsToProject(dbUser, project))
	})
}
