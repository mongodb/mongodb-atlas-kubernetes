package atlasdeployment

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func CreateBasicDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name: "cluster-basics",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 0,
						ZoneName:  name,
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								AnalyticsSpecs: &akov2.Specs{},
								ElectableSpecs: &akov2.Specs{
									InstanceSize: "M2",
									NodeCount:    pointer.MakePtr(3),
								},
								ReadOnlySpecs:       &akov2.Specs{},
								AutoScaling:         &akov2.AdvancedAutoScalingSpec{},
								BackingProviderName: "AWS",
								Priority:            pointer.MakePtr(7),
								ProviderName:        "TENANT",
								RegionName:          "US_EAST_1",
							},
						},
					},
				},
			},
		},
	}
}
