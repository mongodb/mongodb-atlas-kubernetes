package atlasdeployment

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"

	"github.com/stretchr/testify/assert"
)

func TestConvertLegacyDeployment(t *testing.T) {
	deployment := CreateBasicDeployment("deplyment-name")

	t.Run("Legacy Deployment can be converted", func(t *testing.T) {
		err := ConvertLegacyDeployment(deployment)
		assert.NoError(t, err)

		result, err := json.MarshalIndent(deployment, "", "  ")
		assert.NoError(t, err)

		t.Log(string(result))
	})
}

func CreateBasicDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: "my-project",
			},
			DeploymentSpec: &v1.DeploymentSpec{
				Name: "cluster-basics",
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName:    "M2",
					ProviderName:        "TENANT",
					RegionName:          "US_EAST_1",
					BackingProviderName: "AWS",
				},
			},
		},
	}
}
