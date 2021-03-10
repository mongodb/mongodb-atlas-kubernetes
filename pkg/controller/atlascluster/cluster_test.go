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

		merged, err := mergedCluster(atlasCluster, operatorCluster)
		assert.NoError(t, err)

		equal := clustersEqual(zap.S(), atlasCluster, merged)
		assert.True(t, equal)
	})

	t.Run("Clusters don't match (enums)", func(t *testing.T) {
		atlasClusterEnum := mongodbatlas.Cluster{ClusterType: "GEOSHARDED"}
		operatorClusterEnum := mdbv1.AtlasClusterSpec{ClusterType: mdbv1.TypeReplicaSet}

		merged, err := mergedCluster(atlasClusterEnum, operatorClusterEnum)
		assert.NoError(t, err)

		equal := clustersEqual(zap.S(), atlasClusterEnum, merged)
		assert.False(t, equal)
	})
}
