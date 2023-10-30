package atlasdeployment

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func CreateBasicDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: "my-project",
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Name: "cluster-basics",
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 0,
						ZoneName:  name,
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								AnalyticsSpecs: &v1.Specs{},
								ElectableSpecs: &v1.Specs{
									InstanceSize: "M2",
									NodeCount:    toptr.MakePtr(3),
								},
								ReadOnlySpecs:       &v1.Specs{},
								AutoScaling:         &v1.AdvancedAutoScalingSpec{},
								BackingProviderName: "AWS",
								Priority:            toptr.MakePtr(7),
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
