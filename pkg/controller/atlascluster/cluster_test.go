package atlascluster

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

func TestClusterMatchesSpec(t *testing.T) {
	t.Run("Clusters match (enums)", func(t *testing.T) {
		atlasCluster := mongodbatlas.Cluster{
			ProviderSettings: &mongodbatlas.ProviderSettings{
				ProviderName: "AWS",
			},
			ClusterType: "GEOSHARDED",
		}
		operatorCluster := mdbv1.AtlasClusterSpec{
			ProviderSettings: &mdbv1.ProviderSettingsSpec{
				ProviderName: provider.ProviderAWS,
			},
			ClusterType: mdbv1.TypeGeoSharded,
		}

		merged, err := MergedCluster(atlasCluster, operatorCluster)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), atlasCluster, merged)
		assert.True(t, equal)
	})
	t.Run("Clusters don't match (enums)", func(t *testing.T) {
		atlasClusterEnum := mongodbatlas.Cluster{ClusterType: "GEOSHARDED"}
		operatorClusterEnum := mdbv1.AtlasClusterSpec{ClusterType: mdbv1.TypeReplicaSet}

		merged, err := MergedCluster(atlasClusterEnum, operatorClusterEnum)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), atlasClusterEnum, merged)
		assert.False(t, equal)
	})
	t.Run("Clusters match (ProviderSettings.RegionName ignored)", func(t *testing.T) {
		common := mdbv1.DefaultAWSCluster("test-ns", "project-name")
		// Note, that in reality it seems that Atlas nullifies ProviderSettings.RegionName only if RegionsConfig are specified
		// but it's ok not to overcomplicate
		common.Spec.ReplicationSpecs = append(common.Spec.ReplicationSpecs, mdbv1.ReplicationSpec{
			NumShards: int64ptr(2),
		})
		// Emulating Atlas behavior when it nullifies the ProviderSettings.RegionName
		atlasCluster, err := common.DeepCopy().WithRegionName("").Spec.Cluster()
		assert.NoError(t, err)
		operatorCluster := common.DeepCopy()

		merged, err := MergedCluster(*atlasCluster, operatorCluster.Spec)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), *atlasCluster, merged)
		assert.True(t, equal)
	})
	t.Run("Clusters don't match (ProviderSettings.RegionName was changed)", func(t *testing.T) {
		atlasCluster, err := mdbv1.DefaultAWSCluster("test-ns", "project-name").WithRegionName("US_WEST_1").Spec.Cluster()
		assert.NoError(t, err)
		// RegionName has changed and no ReplicationSpecs are specified (meaning ProviderSettings.RegionName is mandatory)
		operatorCluster := mdbv1.DefaultAWSCluster("test-ns", "project-name").WithRegionName("EU_EAST_1")

		merged, err := MergedCluster(*atlasCluster, operatorCluster.Spec)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), *atlasCluster, merged)
		assert.False(t, equal)
	})
	t.Run("Clusters match when Atlas adds default ReplicationSpecs", func(t *testing.T) {
		atlasCluster, err := mdbv1.DefaultAWSCluster("test-ns", "project-name").Spec.Cluster()
		assert.NoError(t, err)
		atlasCluster.ReplicationSpecs = []mongodbatlas.ReplicationSpec{
			{
				ID:        "id",
				NumShards: int64ptr(1),
				ZoneName:  "zone1",
				RegionsConfig: map[string]mongodbatlas.RegionsConfig{
					"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			},
		}
		operatorCluster := mdbv1.DefaultAWSCluster("test-ns", "project-name")
		operatorCluster.Spec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
		}}

		merged, err := MergedCluster(*atlasCluster, operatorCluster.Spec)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), *atlasCluster, merged)
		assert.True(t, equal)
	})
	t.Run("Clusters don't match when Atlas adds default ReplicationSpecs and Operator overrides something", func(t *testing.T) {
		atlasCluster, err := mdbv1.DefaultAWSCluster("test-ns", "project-name").Spec.Cluster()
		assert.NoError(t, err)
		atlasCluster.ReplicationSpecs = []mongodbatlas.ReplicationSpec{
			{
				ID:        "id",
				NumShards: int64ptr(1),
				ZoneName:  "zone1",
				RegionsConfig: map[string]mongodbatlas.RegionsConfig{
					"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				},
			},
		}
		operatorCluster := mdbv1.DefaultAWSCluster("test-ns", "project-name")
		operatorCluster.Spec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(2),
			ZoneName:  "zone5",
		}}

		merged, err := MergedCluster(*atlasCluster, operatorCluster.Spec)
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

		equal := ClustersEqual(zap.S(), *atlasCluster, merged)
		assert.False(t, equal)
	})

	t.Run("Clusters don't match - Operator removed the region", func(t *testing.T) {
		atlasCluster, err := mdbv1.DefaultAWSCluster("test-ns", "project-name").Spec.Cluster()
		assert.NoError(t, err)
		atlasCluster.ReplicationSpecs = []mongodbatlas.ReplicationSpec{{
			ID:        "id",
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mongodbatlas.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
				"US_WEST": {AnalyticsNodes: int64ptr(2), ElectableNodes: int64ptr(5), Priority: int64ptr(6), ReadOnlyNodes: int64ptr(0)},
			}},
		}
		operatorCluster := mdbv1.DefaultAWSCluster("test-ns", "project-name")
		operatorCluster.Spec.ReplicationSpecs = []mdbv1.ReplicationSpec{{
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mdbv1.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)},
			}},
		}

		merged, err := MergedCluster(*atlasCluster, operatorCluster.Spec)
		assert.NoError(t, err)

		expectedReplicationSpecs := []mongodbatlas.ReplicationSpec{{
			ID:        "id",
			NumShards: int64ptr(1),
			ZoneName:  "zone1",
			RegionsConfig: map[string]mongodbatlas.RegionsConfig{
				"US_EAST": {AnalyticsNodes: int64ptr(0), ElectableNodes: int64ptr(3), Priority: int64ptr(7), ReadOnlyNodes: int64ptr(0)}}},
		}
		assert.Equal(t, expectedReplicationSpecs, merged.ReplicationSpecs)

		equal := ClustersEqual(zap.S(), *atlasCluster, merged)
		assert.False(t, equal)
	})
}

func int64ptr(i int64) *int64 {
	return &i
}
