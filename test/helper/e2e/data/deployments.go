package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
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

func CreateDeploymentWithKeepPolicy(name string) *akov2.AtlasDeployment {
	deployment := CreateBasicDeployment(name)
	deployment.SetAnnotations(map[string]string{
		customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
	})
	return deployment
}

func CreateAdvancedGeoshardedDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ClusterType: "GEOSHARDED",
				Name:        name,
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AWS",
								RegionName:   "US_EAST_1",
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(3),
								},
							},
						},
					},
					{
						NumShards: 1,
						ZoneName:  "Zone 2",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: "AZURE",
								RegionName:   "EUROPE_NORTH",
								Priority:     pointer.MakePtr(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(3),
								},
							},
						},
					},
				},
			},
		},
	}
}

func CreateServerlessDeployment(name string, providerName string, regionName string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			ServerlessSpec: &akov2.ServerlessSpec{
				Name: name,
				ProviderSettings: &akov2.ProviderSettingsSpec{
					ProviderName:        ServerlessProviderName,
					BackingProviderName: providerName,
					RegionName:          regionName,
				},
			},
		},
	}
}

func CreateBasicDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ClusterType: "REPLICASET",
				Name:        "cluster-basics",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						ZoneName: "test zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: "M2",
									NodeCount:    pointer.MakePtr(1),
								},
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

func CreateDeploymentWithBackup(name string) *akov2.AtlasDeployment {
	deployment := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ClusterType:   "REPLICASET",
				Name:          "deployment-backup",
				BackupEnabled: pointer.MakePtr(true),
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						ZoneName: "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(3),
								},
								Priority:            pointer.MakePtr(7),
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

func NewDeploymentWithBackupSpec() akov2.AtlasDeploymentSpec {
	return akov2.AtlasDeploymentSpec{
		Project: common.ResourceRefNamespaced{
			Name: ProjectName,
		},
		DeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:          "deployment-backup",
			BackupEnabled: pointer.MakePtr(false),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					ZoneName: "Zone 1",
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ElectableSpecs: &akov2.Specs{
								InstanceSize: InstanceSizeM20,
								NodeCount:    pointer.MakePtr(3),
							},
							Priority:            pointer.MakePtr(7),
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

func CreateDeploymentWithMultiregionAWS(name string) *akov2.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderAWS)
}

func CreateDeploymentWithMultiregionAzure(name string) *akov2.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderAzure)
}

func CreateDeploymentWithMultiregionGCP(name string) *akov2.AtlasDeployment {
	return CreateDeploymentWithMultiregion(name, provider.ProviderGCP)
}

func CreateDeploymentWithMultiregion(name string, providerName provider.ProviderName) *akov2.AtlasDeployment {
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

	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:          "deployment-multiregion",
				BackupEnabled: pointer.MakePtr(true),
				ClusterType:   "REPLICASET",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "US-Zone",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(2),
								},
								AutoScaling:  &akov2.AdvancedAutoScalingSpec{},
								Priority:     pointer.MakePtr(7),
								ProviderName: string(providerName),
								RegionName:   regions[0],
							},
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(1),
								},
								AutoScaling:  &akov2.AdvancedAutoScalingSpec{},
								Priority:     pointer.MakePtr(6),
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

func CreateFreeAdvancedDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:                 name,
				ClusterType:          string(akov2.TypeReplicaSet),
				RootCertType:         "ISRGROOTX1",
				VersionReleaseSystem: "LTS",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM0,
								},
								Priority:            pointer.MakePtr(7),
								ProviderName:        string(provider.ProviderTenant),
								BackingProviderName: string(provider.ProviderAWS),
								RegionName:          AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &akov2.ProcessArgs{
				JavascriptEnabled:         pointer.MakePtr(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               pointer.MakePtr(false),
			},
		},
	}
}

func CreateAdvancedDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:          name,
				BackupEnabled: pointer.MakePtr(false),
				BiConnector: &akov2.BiConnectorSpec{
					Enabled:        pointer.MakePtr(false),
					ReadPreference: "secondary",
				},
				ClusterType:              string(akov2.TypeReplicaSet),
				EncryptionAtRestProvider: "NONE",
				PitEnabled:               pointer.MakePtr(false),
				Paused:                   pointer.MakePtr(false),
				RootCertType:             "ISRGROOTX1",
				VersionReleaseSystem:     "LTS",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    pointer.MakePtr(3),
								},
								Priority:     pointer.MakePtr(7),
								ProviderName: string(provider.ProviderAWS),
								RegionName:   AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &akov2.ProcessArgs{
				JavascriptEnabled:         pointer.MakePtr(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               pointer.MakePtr(false),
			},
		},
	}
}
