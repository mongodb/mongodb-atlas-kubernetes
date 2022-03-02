package atlascluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestMergedAdvancedCluster(t *testing.T) {
	defaultAtlas := v1.DefaultAwsAdvancedCluster("default", "my-project")
	defaultAtlas.Spec.AdvancedClusterSpec.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName = "AWS"

	t.Run("Test merging clusters removes backing provider name if empty", func(t *testing.T) {
		advancedCluster := mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ID:        "123",
					ZoneName:  "Zone1",
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName:          "US_EAST_1",
							BackingProviderName: "",
						},
					},
				},
			},
		}

		merged, err := MergedAdvancedCluster(advancedCluster, defaultAtlas.Spec)
		assert.NoError(t, err)
		assert.Empty(t, merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})

	t.Run("Test merging clusters does not remove backing provider name if it is present in the atlas type", func(t *testing.T) {
		advancedCluster := mongodbatlas.AdvancedCluster{
			ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
				{
					NumShards: 1,
					ID:        "123",
					ZoneName:  "Zone1",
					RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
						{
							RegionName:          "US_EAST_1",
							BackingProviderName: "AWS",
						},
					},
				},
			},
		}

		merged, err := MergedAdvancedCluster(advancedCluster, defaultAtlas.Spec)
		assert.NoError(t, err)
		assert.Equal(t, "AWS", merged.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName)
	})
}
