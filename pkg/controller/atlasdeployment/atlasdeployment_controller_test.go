package atlasdeployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestUniqueKey(t *testing.T) {
	t.Run("Test duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &mdbv1.AtlasDeploymentSpec{
			AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
				Tags: []*mdbv1.TagSpec{{Key: "foo", Value: "true"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &mdbv1.AtlasDeploymentSpec{
			AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{
				Tags: []*mdbv1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foobar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
	t.Run("Test duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &mdbv1.AtlasDeploymentSpec{
			ServerlessSpec: &mdbv1.ServerlessSpec{
				Tags: []*mdbv1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &mdbv1.AtlasDeploymentSpec{
			ServerlessSpec: &mdbv1.ServerlessSpec{
				Tags: []*mdbv1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
}
