package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestIsTagsEqual(t *testing.T) {
	t.Run("Test tags are equal and in same order", func(t *testing.T) {
		k8sCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		atlasCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		areEqual := isTagsEqual(*(atlasCluster.Tags), *(k8sCluster.Tags))
		assert.True(t, areEqual, "Deployments should be equal")
	})
	t.Run("Test tags are different lengths", func(t *testing.T) {
		k8sCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foobar", Value: "true"}}}
		atlasCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		areEqual := isTagsEqual(*(atlasCluster.Tags), *(k8sCluster.Tags))
		assert.False(t, areEqual, "Deployments should not be equal")
	})
}
