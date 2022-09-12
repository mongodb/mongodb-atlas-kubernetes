package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestDeploymentMatchesSpec(t *testing.T) {
	t.Run("Deployments match (enums)", func(t *testing.T) {
		atlasDeployment := mongodbatlas.Cluster{
			ProviderSettings: &mongodbatlas.ProviderSettings{
				ProviderName: "AWS",
			},
			ClusterType: "GEOSHARDED",
		}
		operatorDeployment := mdbv1.AtlasDeploymentSpec{
			DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					ProviderName: provider.ProviderAWS,
				},
				ClusterType: mdbv1.TypeGeoSharded,
			},
		}

		merged, err := MergedDeployment(atlasDeployment, operatorDeployment)
		assert.NoError(t, err)

		equal := DeploymentsEqual(zap.S(), atlasDeployment, merged)
		assert.True(t, equal)
	})
	t.Run("Deployments don't match (enums)", func(t *testing.T) {
		atlasDeploymentEnum := mongodbatlas.Cluster{ClusterType: "GEOSHARDED"}
		operatorDeploymentEnum := mdbv1.AtlasDeploymentSpec{DeploymentSpec: &mdbv1.DeploymentSpec{ClusterType: mdbv1.TypeReplicaSet}}

		merged, err := MergedDeployment(atlasDeploymentEnum, operatorDeploymentEnum)
		assert.NoError(t, err)

		equal := DeploymentsEqual(zap.S(), atlasDeploymentEnum, merged)
		assert.False(t, equal)
	})
	t.Run("Deployments match (ProviderSettings.RegionName ignored)", func(t *testing.T) {
		common := mdbv1.DefaultAWSDeployment("test-ns", "project-name")
		// Note, that in reality it seems that Atlas nullifies ProviderSettings.RegionName only if RegionsConfig are specified
		// but it's ok not to overcomplicate
		common.Spec.DeploymentSpec.ReplicationSpecs = append(common.Spec.DeploymentSpec.ReplicationSpecs, mdbv1.ReplicationSpec{
			NumShards: int64ptr(2),
		})
		// Emulating Atlas behavior when it nullifies the ProviderSettings.RegionName
		atlasDeployment, err := common.DeepCopy().WithRegionName("").Spec.Deployment()
		assert.NoError(t, err)
		operatorDeployment := common.DeepCopy()

		merged, err := MergedDeployment(*atlasDeployment, operatorDeployment.Spec)
		assert.NoError(t, err)

		equal := DeploymentsEqual(zap.S(), *atlasDeployment, merged)
		assert.True(t, equal)
	})
	t.Run("Deployments don't match (ProviderSettings.RegionName was changed)", func(t *testing.T) {
		atlasDeployment, err := mdbv1.DefaultAWSDeployment("test-ns", "project-name").WithRegionName("US_WEST_1").Spec.Deployment()
		assert.NoError(t, err)
		// RegionName has changed and no ReplicationSpecs are specified (meaning ProviderSettings.RegionName is mandatory)
		operatorDeployment := mdbv1.DefaultAWSDeployment("test-ns", "project-name").WithRegionName("EU_EAST_1")

		merged, err := MergedDeployment(*atlasDeployment, operatorDeployment.Spec)
		assert.NoError(t, err)

		equal := DeploymentsEqual(zap.S(), *atlasDeployment, merged)
		assert.False(t, equal)
	})
	t.Run("Deployments match when Atlas adds default ReplicationSpecs", func(t *testing.T) {
		atlasDeployment, err := mdbv1.DefaultAWSDeployment("test-ns", "project-name").Spec.Deployment()
		assert.NoError(t, err)
		atlasDeployment.ReplicationSpecs = []mongodbatlas.ReplicationSpec{
			{
				ID:        "id",
				NumShards: int64ptr(1),
				ZoneName:  "zone1",
				RegionsConfig: map[string]mongodbatlas.RegionsConfig{
					"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			},
		}
		operatorDeployment := mdbv1.DefaultAWSDeployment("test-ns", "project-name")
		operatorDeployment.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
		}}

		merged, err := MergedDeployment(*atlasDeployment, operatorDeployment.Spec)
		assert.NoError(t, err)

		equal := DeploymentsEqual(zap.S(), *atlasDeployment, merged)
		assert.True(t, equal)
	})
	t.Run("Deployments don't match when Atlas adds default ReplicationSpecs and Operator overrides something", func(t *testing.T) {
		atlasDeployment, err := mdbv1.DefaultAWSDeployment("test-ns", "project-name").Spec.Deployment()
		assert.NoError(t, err)
		atlasDeployment.ReplicationSpecs = []mongodbatlas.ReplicationSpec{
			{
				ID:        "id",
				NumShards: int64ptr(1),
				ZoneName:  "zone1",
				RegionsConfig: map[string]mongodbatlas.RegionsConfig{
					"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			},
		}
		operatorDeployment := mdbv1.DefaultAWSDeployment("test-ns", "project-name")
		operatorDeployment.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(2),
			ZoneName:  "zone5",
		}}

		merged, err := MergedDeployment(*atlasDeployment, operatorDeployment.Spec)
		assert.NoError(t, err)

		expectedReplicationSpecs := []mongodbatlas.ReplicationSpec{
			{
				ID:        "id",
				NumShards: int64ptr(2),
				ZoneName:  "zone5",
				RegionsConfig: map[string]mongodbatlas.RegionsConfig{
					"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			},
		}
		assert.Equal(t, expectedReplicationSpecs, merged.ReplicationSpecs)

		equal := DeploymentsEqual(zap.S(), *atlasDeployment, merged)
		assert.False(t, equal)
	})

	t.Run("Deployments don't match - Operator removed the region", func(t *testing.T) {
		atlasDeployment, err := mdbv1.DefaultAWSDeployment("test-ns", "project-name").Spec.Deployment()
		assert.NoError(t, err)
		atlasDeployment.ReplicationSpecs = []mongodbatlas.ReplicationSpec{{
			ID:        "id",
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mongodbatlas.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				"US_WEST": {AnalyticsNodes: int64ptr(2), ElectableNodes: int64ptr(5), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)},
			}},
		}
		operatorDeployment := mdbv1.DefaultAWSDeployment("test-ns", "project-name")
		operatorDeployment.Spec.DeploymentSpec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mdbv1.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
			}},
		}

		merged, err := MergedDeployment(*atlasDeployment, operatorDeployment.Spec)
		assert.NoError(t, err)

		expectedReplicationSpecs := []mongodbatlas.ReplicationSpec{{
			ID:        "id",
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mongodbatlas.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)}}},
		}
		assert.Equal(t, expectedReplicationSpecs, merged.ReplicationSpecs)

		equal := DeploymentsEqual(zap.S(), *atlasDeployment, merged)
		assert.False(t, equal)
	})
}

func int64ptr(i int64) *int64 {
	return &i
}
