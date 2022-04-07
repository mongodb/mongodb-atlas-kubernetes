package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: &mdbv1.AdvancedClusterSpec{}, ClusterSpec: &mdbv1.ClusterSpec{}}
			assert.Error(t, ClusterSpec(spec))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: nil, ClusterSpec: nil}
			assert.Error(t, ClusterSpec(spec))
		})
		t.Run("Instance size not empty when serverless", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: nil, ClusterSpec: &mdbv1.ClusterSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "M10",
					ProviderName:     "SERVERLESS",
				},
			}}
			assert.Error(t, ClusterSpec(spec))
		})
		t.Run("Instance size unset when not serverless", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: nil, ClusterSpec: &mdbv1.ClusterSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "",
					ProviderName:     "AWS",
				},
			}}
			assert.Error(t, ClusterSpec(spec))
		})
	})
	t.Run("Valid cluster specs", func(t *testing.T) {
		t.Run("Advanced cluster spec specified", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: &mdbv1.AdvancedClusterSpec{}, ClusterSpec: nil}
			assert.NoError(t, ClusterSpec(spec))
			assert.Nil(t, ClusterSpec(spec))
		})
		t.Run("Regular cluster specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: nil, ClusterSpec: &mdbv1.ClusterSpec{}}
			assert.NoError(t, ClusterSpec(spec))
			assert.Nil(t, ClusterSpec(spec))
		})

		t.Run("Serverless Cluster", func(t *testing.T) {
			spec := mdbv1.AtlasClusterSpec{AdvancedClusterSpec: nil, ClusterSpec: &mdbv1.ClusterSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					ProviderName: "SERVERLESS",
				},
			}}
			assert.NoError(t, ClusterSpec(spec))
			assert.Nil(t, ClusterSpec(spec))
		})
	})
}
