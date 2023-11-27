package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
)

const (
	InstanceSizeM2  = "M2"
	InstanceSizeM10 = "M10"
	InstanceSizeM20 = "M20"
	InstanceSizeM30 = "M30"
	InstanceSizeM0  = "M0"
	AWSRegion       = "US_EAST_1"

	ServerlessProviderName = "SERVERLESS"
)

func CreateDeploymentWithKeepPolicy(name string) *v1.AtlasDeployment {
	deployment := CreateBasicDeployment(name)
	deployment.SetAnnotations(map[string]string{
		customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
	})
	return deployment
}

func CreateAdvancedGeoshardedDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				ClusterType: "GEOSHARDED",
				Name:        name,
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "US_EAST_1",
								Priority:     toptr.MakePtr(7),
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(3),
								},
							},
						},
					},
					{
						NumShards: 1,
						ZoneName:  "Zone 2",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ProviderName: "AZURE",
								RegionName:   "EUROPE_NORTH",
								Priority:     toptr.MakePtr(7),
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(3),
								},
							},
						},
					},
				},
			},
		},
	}
}

func CreateServerlessDeployment(name string, providerName string, regionName string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			ServerlessSpec: &v1.ServerlessSpec{
				Name: name,
				ProviderSettings: &v1.ProviderSettingsSpec{
					ProviderName:        ServerlessProviderName,
					BackingProviderName: providerName,
					RegionName:          regionName,
				},
			},
		},
	}
}

func CreateBasicDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				ClusterType: "REPLICASET",
				Name:        "cluster-basics",
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						ZoneName: "test zone 1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: "M2",
									NodeCount:    toptr.MakePtr(1),
								},
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

func CreateDeploymentWithBackup(name string) *v1.AtlasDeployment {
	deployment := &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				ClusterType:   "REPLICASET",
				Name:          "deployment-backup",
				BackupEnabled: toptr.MakePtr(true),
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						ZoneName: "Zone 1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(3),
								},
								Priority:            toptr.MakePtr(7),
								ProviderName:        "AWS",
								BackingProviderName: "AWS",
								RegionName:          "US_EAST_1",
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func NewDeploymentWithBackupSpec() v1.AtlasDeploymentSpec {
	return v1.AtlasDeploymentSpec{
		Project: common.ResourceRefNamespaced{
			Name: ProjectName,
		},
		DeploymentSpec: &v1.AdvancedDeploymentSpec{
			Name:          "deployment-backup",
			BackupEnabled: toptr.MakePtr(false),
			ReplicationSpecs: []*v1.AdvancedReplicationSpec{
				{
					ZoneName: "Zone 1",
					RegionConfigs: []*v1.AdvancedRegionConfig{
						{
							ElectableSpecs: &v1.Specs{
								InstanceSize: InstanceSizeM20,
								NodeCount:    toptr.MakePtr(3),
							},
							Priority:            toptr.MakePtr(7),
							ProviderName:        "AWS",
							BackingProviderName: "AWS",
							RegionName:          "US_EAST_1",
						},
					},
				},
			},
		},
	}
}

func CreateDeploymentWithMultiregionAWS(name string) *v1.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderAWS)
}

func CreateDeploymentWithMultiregionAzure(name string) *v1.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderAzure)
}

func CreateDeploymentWithMultiregionGCP(name string) *v1.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderGCP)
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
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Name:          "deployment-multiregion",
				BackupEnabled: toptr.MakePtr(true),
				ClusterType:   "REPLICASET",
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "US-Zone",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(2),
								},
								AutoScaling:  &v1.AdvancedAutoScalingSpec{},
								Priority:     toptr.MakePtr(7),
								ProviderName: string(providerName),
								RegionName:   regions[0],
							},
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(1),
								},
								AutoScaling:  &v1.AdvancedAutoScalingSpec{},
								Priority:     toptr.MakePtr(6),
								ProviderName: string(providerName),
								RegionName:   regions[1],
							},
						},
					},
				},
			},
		},
	}
}

func CreateFreeAdvancedDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Name:                 name,
				ClusterType:          string(v1.TypeReplicaSet),
				RootCertType:         "ISRGROOTX1",
				VersionReleaseSystem: "LTS",
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM0,
								},
								Priority:            toptr.MakePtr(7),
								ProviderName:        string(provider.ProviderTenant),
								BackingProviderName: string(provider.ProviderAWS),
								RegionName:          AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &v1.ProcessArgs{
				JavascriptEnabled:         toptr.MakePtr(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               toptr.MakePtr(false),
			},
		},
	}
}

func CreateAdvancedDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Name:          name,
				BackupEnabled: toptr.MakePtr(false),
				BiConnector: &v1.BiConnectorSpec{
					Enabled:        toptr.MakePtr(false),
					ReadPreference: "secondary",
				},
				ClusterType:              string(v1.TypeReplicaSet),
				EncryptionAtRestProvider: "NONE",
				PitEnabled:               toptr.MakePtr(false),
				Paused:                   toptr.MakePtr(false),
				RootCertType:             "ISRGROOTX1",
				VersionReleaseSystem:     "LTS",
				ReplicationSpecs: []*v1.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*v1.AdvancedRegionConfig{
							{
								ElectableSpecs: &v1.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    toptr.MakePtr(3),
								},
								Priority:     toptr.MakePtr(7),
								ProviderName: string(provider.ProviderAWS),
								RegionName:   AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &v1.ProcessArgs{
				JavascriptEnabled:         toptr.MakePtr(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               toptr.MakePtr(false),
			},
		},
	}
}
