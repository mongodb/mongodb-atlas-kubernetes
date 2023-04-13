package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
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
			AdvancedDeploymentSpec: &v1.AdvancedDeploymentSpec{
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

func CreateRegularGeoshardedDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.DeploymentSpec{
				Name: name,
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName: InstanceSizeM30,
					ProviderName:     "AWS",
				},
				ClusterType: "GEOSHARDED",
				ReplicationSpecs: []v1.ReplicationSpec{
					{
						NumShards: toptr.MakePtr(int64(1)),
						ZoneName:  "Zone 1",
						RegionsConfig: map[string]v1.RegionsConfig{
							"US_EAST_1": {
								AnalyticsNodes: toptr.MakePtr(int64(0)),
								ElectableNodes: toptr.MakePtr(int64(3)),
								Priority:       toptr.MakePtr(int64(7)),
								ReadOnlyNodes:  toptr.MakePtr(int64(0)),
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
			DeploymentSpec: &v1.DeploymentSpec{
				Name: "cluster-basics",
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName:    InstanceSizeM2,
					ProviderName:        "TENANT",
					RegionName:          "US_EAST_1",
					BackingProviderName: "AWS",
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
			DeploymentSpec: &v1.DeploymentSpec{
				Name:                  "deployment-backup",
				ProviderBackupEnabled: toptr.MakePtr(true),
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName: InstanceSizeM10,
					ProviderName:     "AWS",
					RegionName:       "US_EAST_1",
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
		DeploymentSpec: &v1.DeploymentSpec{
			Name:                  "deployment-backup",
			ProviderBackupEnabled: toptr.MakePtr(false),
			ProviderSettings: &v1.ProviderSettingsSpec{
				InstanceSizeName: InstanceSizeM20,
				ProviderName:     "AWS",
				RegionName:       "US_EAST_1",
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
			DeploymentSpec: &v1.DeploymentSpec{
				Name:                  "deployment-multiregion",
				ProviderBackupEnabled: toptr.MakePtr(true),
				ClusterType:           "REPLICASET",
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName: InstanceSizeM10,
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

func CreateBasicFreeDeployment(name string) *v1.AtlasDeployment {
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1.AtlasDeploymentSpec{
			Project: common.ResourceRefNamespaced{
				Name: ProjectName,
			},
			DeploymentSpec: &v1.DeploymentSpec{
				Name: name,
				ProviderSettings: &v1.ProviderSettingsSpec{
					InstanceSizeName:    InstanceSizeM0,
					ProviderName:        "TENANT",
					RegionName:          "US_EAST_1",
					BackingProviderName: "AWS",
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
			AdvancedDeploymentSpec: &v1.AdvancedDeploymentSpec{
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
			AdvancedDeploymentSpec: &v1.AdvancedDeploymentSpec{
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
