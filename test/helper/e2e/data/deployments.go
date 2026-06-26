// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

const (
	InstanceSizeM10 = "M10"
	InstanceSizeM20 = "M20"
	InstanceSizeM0  = "M0"
	InstanceSizeM30 = "M30"
	AWSRegion       = "US_EAST_1"
	AWSRegionWest   = "US_WEST_2"
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
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
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
								Priority:     new(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(3),
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
								Priority:     new(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(3),
								},
							},
						},
					},
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
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			FlexSpec: &akov2.FlexSpec{
				Name: "cluster-basics",
				ProviderSettings: &akov2.FlexProviderSettings{
					BackingProviderName: "AWS",
					RegionName:          "US_EAST_1",
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
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				ClusterType:   "REPLICASET",
				Name:          "deployment-backup",
				BackupEnabled: new(true),
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						ZoneName: "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(3),
								},
								Priority:            new(7),
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
		ProjectDualReference: akov2.ProjectDualReference{
			ProjectRef: &common.ResourceRefNamespaced{
				Name: ProjectName,
			},
		},
		DeploymentSpec: &akov2.AdvancedDeploymentSpec{
			Name:          "deployment-backup",
			BackupEnabled: new(false),
			ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
				{
					ZoneName: "Zone 1",
					RegionConfigs: []*akov2.AdvancedRegionConfig{
						{
							ElectableSpecs: &akov2.Specs{
								InstanceSize: InstanceSizeM20,
								NodeCount:    new(3),
							},
							Priority:            new(7),
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
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:          "deployment-multiregion",
				BackupEnabled: new(true),
				ClusterType:   "REPLICASET",
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "US-Zone",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(2),
								},
								AutoScaling:  &akov2.AdvancedAutoScalingSpec{},
								Priority:     new(7),
								ProviderName: string(providerName),
								RegionName:   regions[0],
							},
							{
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(1),
								},
								AutoScaling:  &akov2.AdvancedAutoScalingSpec{},
								Priority:     new(6),
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
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
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
								Priority:            new(7),
								ProviderName:        string(provider.ProviderTenant),
								BackingProviderName: string(provider.ProviderAWS),
								RegionName:          AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &akov2.ProcessArgs{
				JavascriptEnabled:         new(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               new(false),
			},
		},
	}
}

// autoscalingSpec returns an AdvancedAutoScalingSpec with compute and diskGB autoscaling enabled.
// minSize/maxSize are Atlas instance size names (e.g. "M10", "M30").
func autoscalingSpec(minSize, maxSize string) *akov2.AdvancedAutoScalingSpec {
	return &akov2.AdvancedAutoScalingSpec{
		Compute: &akov2.ComputeSpec{
			Enabled:          pointer.MakePtr(true),
			ScaleDownEnabled: pointer.MakePtr(true),
			MinInstanceSize:  minSize,
			MaxInstanceSize:  maxSize,
		},
		DiskGB: &akov2.DiskGB{
			Enabled: pointer.MakePtr(true),
		},
	}
}

// CreateDeploymentWithAutoscaling creates a REPLICASET in a single electable region with
// compute and diskGB autoscaling enabled. Use UpdateDeploymentAddReadOnlyRegion to add a
// second region later and exercise the AUTO_SCALINGS_MUST_MATCH code path.
func CreateDeploymentWithAutoscaling(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:          name,
				ClusterType:   string(akov2.TypeReplicaSet),
				BackupEnabled: new(false),
				PitEnabled:    new(false),
				ReplicationSpecs: []*akov2.AdvancedReplicationSpec{
					{
						NumShards: 1,
						ZoneName:  "Zone 1",
						RegionConfigs: []*akov2.AdvancedRegionConfig{
							{
								ProviderName: string(provider.ProviderAWS),
								RegionName:   AWSRegion,
								Priority:     new(7),
								ElectableSpecs: &akov2.Specs{
									InstanceSize: InstanceSizeM10,
									NodeCount:    new(3),
								},
								AutoScaling: autoscalingSpec(InstanceSizeM10, InstanceSizeM30),
							},
						},
					},
				},
			},
		},
	}
}

// UpdateDeploymentAddReadOnlyRegion appends a readOnly region with the same autoscaling
// settings to an existing deployment spec. This exercises the AUTO_SCALINGS_MUST_MATCH
// regression: the operator must send consistent analyticsAutoScaling/autoScaling values
// for all regions in the PATCH request.
func UpdateDeploymentAddReadOnlyRegion(d *akov2.AtlasDeployment) {
	d.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs = append(
		d.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs,
		&akov2.AdvancedRegionConfig{
			ProviderName: string(provider.ProviderAWS),
			RegionName:   AWSRegionWest,
			Priority:     new(0),
			ReadOnlySpecs: &akov2.Specs{
				InstanceSize: InstanceSizeM10,
				NodeCount:    new(1),
			},
			AutoScaling: autoscalingSpec(InstanceSizeM10, InstanceSizeM30),
		},
	)
}

func CreateAdvancedDeployment(name string) *akov2.AtlasDeployment {
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: akov2.AtlasDeploymentSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: ProjectName,
				},
			},
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name:          name,
				BackupEnabled: new(false),
				BiConnector: &akov2.BiConnectorSpec{
					Enabled:        new(false),
					ReadPreference: "secondary",
				},
				ClusterType:              string(akov2.TypeReplicaSet),
				EncryptionAtRestProvider: "NONE",
				PitEnabled:               new(false),
				Paused:                   new(false),
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
									NodeCount:    new(3),
								},
								Priority:     new(7),
								ProviderName: string(provider.ProviderAWS),
								RegionName:   AWSRegion,
							},
						},
					},
				},
			},
			ProcessArgs: &akov2.ProcessArgs{
				JavascriptEnabled:         new(true),
				MinimumEnabledTLSProtocol: "TLS1_2",
				NoTableScan:               new(false),
			},
		},
	}
}
