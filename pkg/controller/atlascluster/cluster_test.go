package atlascluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

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
				ProviderName: mdbv1.ProviderAWS,
			},
			ClusterType: mdbv1.TypeGeoSharded,
		}

		merged, err := MergedCluster(atlasCluster, operatorCluster)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), atlasCluster, merged)
		assert.True(t, equal)
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

	t.Run("Clusters don't match (enums)", func(t *testing.T) {
		atlasClusterEnum := mongodbatlas.Cluster{ClusterType: "GEOSHARDED"}
		operatorClusterEnum := mdbv1.AtlasClusterSpec{ClusterType: mdbv1.TypeReplicaSet}

		merged, err := MergedCluster(atlasClusterEnum, operatorClusterEnum)
		assert.NoError(t, err)

		equal := ClustersEqual(zap.S(), atlasClusterEnum, merged)
		assert.False(t, equal)
	})
}
func int64ptr(i int64) *int64 {
	return &i
}
