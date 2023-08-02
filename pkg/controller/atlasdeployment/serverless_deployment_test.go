package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
)

func TestIsEqual(t *testing.T) {
	//var k8sCluster *mongodbatlas.Cluster
	//var atlasCluster *mongodbatlas.Cluster
	t.Run("Test tags are equal and in same order", func(t *testing.T) {
		k8sCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		atlasCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		areEqual := isEqual(atlasCluster, k8sCluster)
		assert.True(t, areEqual, "Deployments should be equal")
	})
	t.Run("Test tags are different lengths", func(t *testing.T) {
		k8sCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foobar", Value: "true"}}}
		atlasCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		areEqual := isEqual(atlasCluster, k8sCluster)
		assert.False(t, areEqual, "Deployments should not be equal")
	})
	t.Run("Test tags are equal and in wrong order", func(t *testing.T) {
		k8sCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "bar", Value: "false"}, {Key: "foo", Value: "true"}}}
		atlasCluster := &mongodbatlas.Cluster{Tags: &[]*mongodbatlas.Tag{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}}}
		areEqual := isEqual(atlasCluster, k8sCluster)
		assert.False(t, areEqual, "Deployments should not be equal")
	})
}
