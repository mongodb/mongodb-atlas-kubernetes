package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func TestClusterValidation(t *testing.T) {
	t.Run("Invalid cluster specs", func(t *testing.T) {
		t.Run("Multiple specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("No specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: nil}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("Instance size not empty when serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "M10",
					ProviderName:     "SERVERLESS",
				},
			}}
			assert.Error(t, DeploymentSpec(spec))
		})
		t.Run("Instance size unset when not serverless", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					InstanceSizeName: "",
					ProviderName:     "AWS",
				},
			}}
			assert.Error(t, DeploymentSpec(spec))
		})
	})
	t.Run("Valid cluster specs", func(t *testing.T) {
		t.Run("Advanced cluster spec specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: &mdbv1.AdvancedDeploymentSpec{}, DeploymentSpec: nil}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})
		t.Run("Regular cluster specs specified", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{}}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})

		t.Run("Serverless Cluster", func(t *testing.T) {
			spec := mdbv1.AtlasDeploymentSpec{AdvancedDeploymentSpec: nil, DeploymentSpec: &mdbv1.DeploymentSpec{
				ProviderSettings: &mdbv1.ProviderSettingsSpec{
					ProviderName: "SERVERLESS",
				},
			}}
			assert.NoError(t, DeploymentSpec(spec))
			assert.Nil(t, DeploymentSpec(spec))
		})
	})
}
