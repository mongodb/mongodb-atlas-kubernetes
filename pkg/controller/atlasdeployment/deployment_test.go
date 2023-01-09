package atlasdeployment

import (
	"encoding/json"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	"github.com/stretchr/testify/assert"
)

func TestConvertLegacyDeployment(t *testing.T) {
	deployment := CreateBasicDeployment("deplyment-name")

	t.Run("Legacy Deployment can be converted", func(t *testing.T) {
		err := ConvertLegacyDeployment(&deployment.Spec)
		assert.NoError(t, err)

		result, err := json.MarshalIndent(deployment, "", "  ")
		assert.NoError(t, err)

		t.Log(string(result))
	})

	deploymentMultiregion := CreateDeploymentWithMultiregion("deplyment-multiregion-name", provider.ProviderAWS)

	t.Run("Legacy Multiregion Deployment can be converted", func(t *testing.T) {
		err := ConvertLegacyDeployment(&deploymentMultiregion.Spec)
		assert.NoError(t, err)

		result, err := json.MarshalIndent(deploymentMultiregion, "", "  ")
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

func CreateDeploymentWithMultiregion(name string, providerName provider.ProviderName) *v1.AtlasDeployment {
	var regions []string
	switch providerName {
	case provider.ProviderAWS:
		regions = []string{"US_EAST_1", "US_WEST_2"}
	case provider.ProviderAzure:
		regions = []string{"NORWAY_EAST", "GERMANY_NORTH"}
	case provider.ProviderGCP:
		regions = []string{"CENTRAL_US", "EASTERN_US"}
	}

	if len(regions) == 0 {
		panic("unknown provider")
	}

	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: "my-project",
			},
			DeploymentSpec: &v1.DeploymentSpec{
				Name:                  "deployment-multiregion",
				ProviderBackupEnabled: toptr.MakePtr(true),
				ClusterType:           "REPLICASET",
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName: "M10",
					ProviderName:     providerName,
				},
				ReplicationSpecs: []v1.ReplicationSpec{
					{
						NumShards: toptr.MakePtr(int64(1)),
						ZoneName:  "US-Zone",
						RegionsConfig: map[string]v1.RegionsConfig{
							regions[0]: {
								AnalyticsNodes: toptr.MakePtr(int64(0)),
								ElectableNodes: toptr.MakePtr(int64(1)),
								Priority:       toptr.MakePtr(int64(6)),
								ReadOnlyNodes:  toptr.MakePtr(int64(0)),
							},
							regions[1]: {
								AnalyticsNodes: toptr.MakePtr(int64(0)),
								ElectableNodes: toptr.MakePtr(int64(2)),
								Priority:       toptr.MakePtr(int64(7)),
								ReadOnlyNodes:  toptr.MakePtr(int64(0)),
							},
						},
					},
				},
			},
		}}
}
