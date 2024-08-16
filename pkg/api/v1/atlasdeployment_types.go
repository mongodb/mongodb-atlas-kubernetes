/*
Copyright 2020 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"errors"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
)

func init() {
	SchemeBuilder.Register(&AtlasDeployment{}, &AtlasDeploymentList{})
}

type DeploymentType string

const (
	TypeReplicaSet DeploymentType = "REPLICASET"
	TypeSharded    DeploymentType = "SHARDED"
	TypeGeoSharded DeploymentType = "GEOSHARDED"
)

// AtlasDeploymentSpec defines the desired state of AtlasDeployment
// Only one of DeploymentSpec, AdvancedDeploymentSpec and ServerlessSpec should be defined
type AtlasDeploymentSpec struct {
	// Project is a reference to AtlasProject resource the deployment belongs to
	Project common.ResourceRefNamespaced `json:"projectRef"`

	// Configuration for the advanced (v1.5) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters/
	// +optional
	DeploymentSpec *AdvancedDeploymentSpec `json:"deploymentSpec,omitempty"`

	// Backup schedule for the AtlasDeployment
	// +optional
	BackupScheduleRef common.ResourceRefNamespaced `json:"backupRef"`

	// Configuration for the serverless deployment API. https://www.mongodb.com/docs/atlas/reference/api/serverless-instances/
	// +optional
	ServerlessSpec *ServerlessSpec `json:"serverlessSpec,omitempty"`

	// ProcessArgs allows to modify Advanced Configuration Options
	// +optional
	ProcessArgs *ProcessArgs `json:"processArgs,omitempty"`
}

type SearchNode struct {
	// Hardware specification for the Search Node instance sizes.
	// +kubebuilder:validation:Enum:=S20_HIGHCPU_NVME;S30_HIGHCPU_NVME;S40_HIGHCPU_NVME;S50_HIGHCPU_NVME;S60_HIGHCPU_NVME;S70_HIGHCPU_NVME;S80_HIGHCPU_NVME;S30_LOWCPU_NVME;S40_LOWCPU_NVME;S50_LOWCPU_NVME;S60_LOWCPU_NVME;S80_LOWCPU_NVME;S90_LOWCPU_NVME;S100_LOWCPU_NVME;S110_LOWCPU_NVME
	InstanceSize string `json:"instanceSize,omitempty"`
	// Number of Search Nodes in the cluster.
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=32
	NodeCount uint8 `json:"nodeCount,omitempty"`
}

type AdvancedDeploymentSpec struct {
	// Applicable only for M10+ deployments.
	// Flag that indicates if the deployment uses Cloud Backups for backups.
	// +optional
	BackupEnabled *bool `json:"backupEnabled,omitempty"`
	// Configuration of BI Connector for Atlas on this deployment.
	// The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger deployments.
	// +optional
	BiConnector *BiConnectorSpec `json:"biConnector,omitempty"`
	// Type of the deployment that you want to create.
	// The parameter is required if replicationSpecs are set or if Global Deployments are deployed.
	// +kubebuilder:validation:Enum=REPLICASET;SHARDED;GEOSHARDED
	// +optional
	ClusterType string `json:"clusterType,omitempty"`
	// Capacity, in gigabytes, of the host's root volume.
	// Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB).
	// This value must be a positive integer.
	// The parameter is required if replicationSpecs are configured.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4096
	// +optional
	DiskSizeGB *int `json:"diskSizeGB,omitempty"`
	// Cloud service provider that offers Encryption at Rest.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;NONE
	// +optional
	EncryptionAtRestProvider string `json:"encryptionAtRestProvider,omitempty"`
	// Collection of key-value pairs that tag and categorize the deployment.
	// Each key and value has a maximum length of 255 characters.
	// +optional
	Labels []common.LabelSpec `json:"labels,omitempty"`
	// Version of the deployment to deploy.
	MongoDBMajorVersion string `json:"mongoDBMajorVersion,omitempty"`
	MongoDBVersion      string `json:"mongoDBVersion,omitempty"`
	// Name of the advanced deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name string `json:"name,omitempty"`
	// Flag that indicates whether the deployment should be paused.
	Paused *bool `json:"paused,omitempty"`
	// Flag that indicates the deployment uses continuous cloud backups.
	// +optional
	PitEnabled *bool `json:"pitEnabled,omitempty"`
	// Configuration for deployment regions.
	// +optional
	ReplicationSpecs []*AdvancedReplicationSpec `json:"replicationSpecs,omitempty"`
	RootCertType     string                     `json:"rootCertType,omitempty"`
	// Key-value pairs for resource tagging.
	// +kubebuilder:validation:MaxItems=50
	// +optional
	Tags                 []*TagSpec `json:"tags,omitempty"`
	VersionReleaseSystem string     `json:"versionReleaseSystem,omitempty"`
	// +optional
	CustomZoneMapping []CustomZoneMapping `json:"customZoneMapping,omitempty"`
	// +optional
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`
	// Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
	// +kubebuilder:default:=false
	TerminationProtectionEnabled bool `json:"terminationProtectionEnabled,omitempty"`
	// Settings for Search Nodes for the cluster. Currently, at most one search node configuration may be defined.
	// +kubebuilder:validation:MaxItems=1
	// +optional
	SearchNodes []SearchNode `json:"searchNodes,omitempty"`
	// A list of atlas search indexes configuration for the current deployment
	// +optional
	SearchIndexes []SearchIndex `json:"searchIndexes,omitempty"`
}

// ToAtlas converts the AdvancedDeploymentSpec to native Atlas client ToAtlas format.
func (s *AdvancedDeploymentSpec) ToAtlas() (*mongodbatlas.AdvancedCluster, error) {
	result := &mongodbatlas.AdvancedCluster{}
	err := compat.JSONCopy(result, s)
	return result, err
}

func (s *AdvancedDeploymentSpec) SearchNodesToAtlas() []admin.ApiSearchDeploymentSpec {
	if len(s.SearchNodes) == 0 {
		return nil
	}
	result := make([]admin.ApiSearchDeploymentSpec, len(s.SearchNodes))

	for i := 0; i < len(s.SearchNodes); i++ {
		result[i] = admin.ApiSearchDeploymentSpec{
			InstanceSize: s.SearchNodes[i].InstanceSize,
			NodeCount:    int(s.SearchNodes[i].NodeCount),
		}
	}
	return result
}

// ServerlessSpec defines the desired state of Atlas Serverless Instance
type ServerlessSpec struct {
	// Name of the serverless deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name string `json:"name"`
	// Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.
	ProviderSettings *ServerlessProviderSettingsSpec `json:"providerSettings"`
	PrivateEndpoints []ServerlessPrivateEndpoint     `json:"privateEndpoints,omitempty"`
	// Key-value pairs for resource tagging.
	// +kubebuilder:validation:MaxItems=50
	// +optional
	Tags []*TagSpec `json:"tags,omitempty"`

	// Serverless Backup Options
	BackupOptions ServerlessBackupOptions `json:"backupOptions,omitempty"`

	// Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
	// +kubebuilder:default:=false
	TerminationProtectionEnabled bool `json:"terminationProtectionEnabled,omitempty"`
}

// ToAtlas converts the ServerlessSpec to native Atlas client Cluster format.
func (s *ServerlessSpec) ToAtlas() (*mongodbatlas.Cluster, error) {
	result := &mongodbatlas.Cluster{}
	err := compat.JSONCopy(result, s)
	return result, err
}

// BiConnector specifies BI Connector for Atlas configuration on this deployment.
type BiConnector struct {
	Enabled        *bool  `json:"enabled,omitempty"`
	ReadPreference string `json:"readPreference,omitempty"`
}

// TagSpec holds a key-value pair for resource tagging on this deployment.
type TagSpec struct {
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9 @_.+`;`-]*$
	Key string `json:"key"`
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9@_.+`;`-]*$
	Value string `json:"value"`
}

type ServerlessBackupOptions struct {
	// ServerlessContinuousBackupEnabled
	// +kubebuilder:default:=true
	ServerlessContinuousBackupEnabled bool `json:"serverlessContinuousBackupEnabled,omitempty"`
}

// ConnectionStrings configuration for applications use to connect to this deployment.
type ConnectionStrings struct {
	Standard          string                `json:"standard,omitempty"`
	StandardSrv       string                `json:"standardSrv,omitempty"`
	PrivateEndpoint   []PrivateEndpointSpec `json:"privateEndpoint,omitempty"`
	AwsPrivateLink    map[string]string     `json:"awsPrivateLink,omitempty"`
	AwsPrivateLinkSrv map[string]string     `json:"awsPrivateLinkSrv,omitempty"`
	Private           string                `json:"private,omitempty"`
	PrivateSrv        string                `json:"privateSrv,omitempty"`
}

// PrivateEndpointSpec connection strings. Each object describes the connection strings
// you can use to connect to this deployment through a private endpoint.
// Atlas returns this parameter only if you deployed a private endpoint to all regions
// to which you deployed this deployment's nodes.
type PrivateEndpointSpec struct {
	ConnectionString    string         `json:"connectionString,omitempty"`
	Endpoints           []EndpointSpec `json:"endpoints,omitempty"`
	SRVConnectionString string         `json:"srvConnectionString,omitempty"`
	Type                string         `json:"type,omitempty"`
}

// EndpointSpec through which you connect to Atlas.
type EndpointSpec struct {
	EndpointID   string `json:"endpointId,omitempty"`
	ProviderName string `json:"providerName,omitempty"`
	Region       string `json:"region,omitempty"`
}

type AdvancedReplicationSpec struct {
	// Positive integer that specifies the number of shards to deploy in each specified zone.
	// If you set this value to 1 and clusterType is SHARDED, MongoDB Cloud deploys a single-shard sharded cluster.
	// Don't create a sharded cluster with a single shard for production environments.
	// Single-shard sharded clusters don't provide the same benefits as multi-shard configurations
	NumShards int `json:"numShards,omitempty"`
	// Human-readable label that identifies the zone in a Global Cluster.
	ZoneName string `json:"zoneName,omitempty"`
	// Hardware specifications for nodes set for a given region.
	// Each regionConfigs object describes the region's priority in elections and the number and type of MongoDB nodes that MongoDB Cloud deploys to the region.
	// Each regionConfigs object must have either an analyticsSpecs object, electableSpecs object, or readOnlySpecs object.
	// Tenant clusters only require electableSpecs. Dedicated clusters can specify any of these specifications, but must have at least one electableSpecs object within a replicationSpec.
	// Every hardware specification must use the same instanceSize.
	RegionConfigs []*AdvancedRegionConfig `json:"regionConfigs,omitempty"`
}

type AdvancedRegionConfig struct {
	AnalyticsSpecs *Specs                   `json:"analyticsSpecs,omitempty"`
	ElectableSpecs *Specs                   `json:"electableSpecs,omitempty"`
	ReadOnlySpecs  *Specs                   `json:"readOnlySpecs,omitempty"`
	AutoScaling    *AdvancedAutoScalingSpec `json:"autoScaling,omitempty"`
	// Cloud service provider on which the host for a multi-tenant deployment is provisioned.
	// This setting only works when "providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
	// Otherwise it should be equal to "providerName" value
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	BackingProviderName string `json:"backingProviderName,omitempty"`
	// Precedence is given to this region when a primary election occurs.
	// If your regionConfigs has only readOnlySpecs, analyticsSpecs, or both, set this value to 0.
	// If you have multiple regionConfigs objects (your cluster is multi-region or multi-cloud), they must have priorities in descending order.
	// The highest priority is 7
	Priority *int `json:"priority,omitempty"`
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT;SERVERLESS
	ProviderName string `json:"providerName,omitempty"`
	// Physical location of your MongoDB deployment.
	// The region you choose can affect network latency for clients accessing your databases.
	RegionName string `json:"regionName,omitempty"`
}

type Specs struct {
	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +optional
	DiskIOPS *int64 `json:"diskIOPS,omitempty"`
	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +kubebuilder:validation:Enum=STANDARD;PROVISIONED
	EbsVolumeType string `json:"ebsVolumeType,omitempty"`
	// Hardware specification for the instance sizes in this region.
	// Each instance size has a default storage and memory capacity.
	// The instance size you select applies to all the data-bearing hosts in your instance size
	InstanceSize string `json:"instanceSize,omitempty"`
	// Number of nodes of the given type for MongoDB Cloud to deploy to the region.
	NodeCount *int `json:"nodeCount,omitempty"`
}

// AutoScalingSpec configures your deployment to automatically scale its storage
type AutoScalingSpec struct {
	// Deprecated: This flag is not supported anymore.
	// Flag that indicates whether autopilot mode for Performance Advisor is enabled.
	// The default is false.
	AutoIndexingEnabled *bool `json:"autoIndexingEnabled,omitempty"`
	// Flag that indicates whether disk auto-scaling is enabled. The default is true.
	// +optional
	DiskGBEnabled *bool `json:"diskGBEnabled,omitempty"`

	// Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.
	// +optional
	Compute *ComputeSpec `json:"compute,omitempty"`
}

// AdvancedAutoScalingSpec configures your deployment to automatically scale its storage
type AdvancedAutoScalingSpec struct {
	// Flag that indicates whether disk auto-scaling is enabled. The default is true.
	// +optional
	DiskGB *DiskGB `json:"diskGB,omitempty"`

	// Collection of settings that configure how a deployment might scale its deployment tier and whether the deployment can scale down.
	// +optional
	Compute *ComputeSpec `json:"compute,omitempty"`
}

// DiskGB specifies whether disk auto-scaling is enabled. The default is true.
type DiskGB struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// ComputeSpec Specifies whether the deployment automatically scales its deployment tier and whether the deployment can scale down.
type ComputeSpec struct {
	// Flag that indicates whether deployment tier auto-scaling is enabled. The default is false.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Flag that indicates whether the deployment tier may scale down. Atlas requires this parameter if "autoScaling.compute.enabled" : true.
	// +optional
	ScaleDownEnabled *bool `json:"scaleDownEnabled,omitempty"`

	// Minimum instance size to which your deployment can automatically scale (such as M10). Atlas requires this parameter if "autoScaling.compute.scaleDownEnabled" : true.
	// +optional
	MinInstanceSize string `json:"minInstanceSize,omitempty"`

	// Maximum instance size to which your deployment can automatically scale (such as M40). Atlas requires this parameter if "autoScaling.compute.enabled" : true.
	// +optional
	MaxInstanceSize string `json:"maxInstanceSize,omitempty"`
}

type ProcessArgs struct {
	DefaultReadConcern               string `json:"defaultReadConcern,omitempty"`
	DefaultWriteConcern              string `json:"defaultWriteConcern,omitempty"`
	MinimumEnabledTLSProtocol        string `json:"minimumEnabledTlsProtocol,omitempty"`
	FailIndexKeyTooLong              *bool  `json:"failIndexKeyTooLong,omitempty"`
	JavascriptEnabled                *bool  `json:"javascriptEnabled,omitempty"`
	NoTableScan                      *bool  `json:"noTableScan,omitempty"`
	OplogSizeMB                      *int64 `json:"oplogSizeMB,omitempty"`
	SampleSizeBIConnector            *int64 `json:"sampleSizeBIConnector,omitempty"`
	SampleRefreshIntervalBIConnector *int64 `json:"sampleRefreshIntervalBIConnector,omitempty"`
	OplogMinRetentionHours           string `json:"oplogMinRetentionHours,omitempty"`
}

// Check compatibility with library type.
var _ = ComputeSpec(mongodbatlas.Compute{})

// BiConnectorSpec specifies BI Connector for Atlas configuration on this deployment
type BiConnectorSpec struct {
	// Flag that indicates whether or not BI Connector for Atlas is enabled on the deployment.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Source from which the BI Connector for Atlas reads data. Each BI Connector for Atlas read preference contains a distinct combination of readPreference and readPreferenceTags options.
	// +optional
	ReadPreference string `json:"readPreference,omitempty"`
}

// Check compatibility with library type.
var _ = BiConnectorSpec(mongodbatlas.BiConnector{})

// ServerlessProviderSettingsSpec configuration for the provisioned servers on which MongoDB runs. The available options are specific to the cloud service provider.
type ServerlessProviderSettingsSpec struct {
	// Cloud service provider on which the host for a multi-tenant deployment is provisioned.
	// This setting only works when "providerSetting.providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +optional
	BackingProviderName string `json:"backingProviderName,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +optional
	DiskIOPS *int64 `json:"diskIOPS,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Type of disk if you selected Azure as your cloud service provider.
	// +optional
	DiskTypeName string `json:"diskTypeName,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the deployment.
	// +optional
	EncryptEBSVolume *bool `json:"encryptEBSVolume,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Atlas provides different deployment tiers, each with a default storage capacity and RAM size. The deployment you select is used for all the data-bearing hosts in your deployment tier.
	// +optional
	InstanceSizeName string `json:"instanceSizeName,omitempty"`

	// Cloud service provider on which Atlas provisions the hosts.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT;SERVERLESS
	ProviderName provider.ProviderName `json:"providerName"`

	// Physical location of your MongoDB deployment.
	// The region you choose can affect network latency for clients accessing your databases.
	// +optional
	RegionName string `json:"regionName,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +kubebuilder:validation:Enum=STANDARD;PROVISIONED
	VolumeType string `json:"volumeType,omitempty"`

	// DEPRECATED FIELD. The value of this field doesn't take any effect. Range of instance sizes to which your deployment can scale.
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`
}

// Deployment converts the Spec to native Atlas client format.
func (spec *AtlasDeploymentSpec) Deployment() (*mongodbatlas.AdvancedCluster, error) {
	result := &mongodbatlas.AdvancedCluster{}
	if spec.DeploymentSpec == nil {
		return result, errors.New("AdvancedDeploymentSpec is empty")
	}
	err := compat.JSONCopy(result, *spec.DeploymentSpec)
	return result, err
}

var _ api.AtlasCustomResource = &AtlasDeployment{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Atlas State",type=string,JSONPath=`.status.stateName`
// +kubebuilder:printcolumn:name="MongoDB Version",type=string,JSONPath=`.status.mongoDBVersion`
// +kubebuilder:resource:categories=atlas,shortName=ad

// AtlasDeployment is the Schema for the atlasdeployments API
type AtlasDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasDeploymentSpec          `json:"spec,omitempty"`
	Status status.AtlasDeploymentStatus `json:"status,omitempty"`
}

func (c *AtlasDeployment) GetDeploymentName() string {
	if c.IsServerless() {
		return c.Spec.ServerlessSpec.Name
	}
	if c.IsAdvancedDeployment() {
		return c.Spec.DeploymentSpec.Name
	}

	return ""
}

// IsServerless returns true if the AtlasDeployment is configured to be a serverless instance
func (c *AtlasDeployment) IsServerless() bool {
	return c.Spec.ServerlessSpec != nil
}

// IsAdvancedDeployment returns true if the AtlasDeployment is configured to be an advanced deployment.
func (c *AtlasDeployment) IsAdvancedDeployment() bool {
	return c.Spec.DeploymentSpec != nil
}

func (c *AtlasDeployment) GetReplicationSetID() string {
	if len(c.Status.ReplicaSets) > 0 {
		return c.Status.ReplicaSets[0].ID
	}

	return ""
}

// +kubebuilder:object:root=true

// AtlasDeploymentList contains a list of AtlasDeployment
type AtlasDeploymentList struct {
	metav1.TypeMeta `                  json:",inline"`
	metav1.ListMeta `                  json:"metadata,omitempty"`
	Items           []AtlasDeployment `json:"items"`
}

func (c AtlasDeployment) AtlasProjectObjectKey() client.ObjectKey {
	ns := c.Namespace
	if c.Spec.Project.Namespace != "" {
		ns = c.Spec.Project.Namespace
	}
	return kube.ObjectKey(ns, c.Spec.Project.Name)
}

func (c *AtlasDeployment) GetStatus() api.Status {
	return c.Status
}

func (c *AtlasDeployment) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	c.Status.Conditions = conditions
	c.Status.ObservedGeneration = c.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.AtlasDeploymentStatusOption)
		v(&c.Status)
	}
}

// ************************************ Builder methods *************************************************

func NewDeployment(namespace, name, nameInAtlas string) *AtlasDeployment {
	return &AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasDeploymentSpec{
			DeploymentSpec: &AdvancedDeploymentSpec{
				ClusterType: "REPLICASET",
				Name:        nameInAtlas,
				ReplicationSpecs: []*AdvancedReplicationSpec{
					{
						ZoneName: "Zone 1",
						RegionConfigs: []*AdvancedRegionConfig{
							{
								Priority: pointer.MakePtr(7),
								ElectableSpecs: &Specs{
									InstanceSize: "M10",
									NodeCount:    pointer.MakePtr(3),
								},
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
}

func newServerlessInstance(namespace, name, nameInAtlas, backingProviderName, regionName string) *AtlasDeployment {
	return &AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasDeploymentSpec{
			ServerlessSpec: &ServerlessSpec{
				Name: nameInAtlas,
				ProviderSettings: &ServerlessProviderSettingsSpec{
					BackingProviderName: backingProviderName,
					ProviderName:        "SERVERLESS",
					RegionName:          regionName,
				},
			},
		},
	}
}

func addReplicaIfNotAdded(deployment *AtlasDeployment) {
	if deployment == nil {
		return
	}

	if deployment.Spec.DeploymentSpec == nil {
		return
	}

	if len(deployment.Spec.DeploymentSpec.ReplicationSpecs) == 0 {
		deployment.Spec.DeploymentSpec.ReplicationSpecs = append(deployment.Spec.DeploymentSpec.ReplicationSpecs, &AdvancedReplicationSpec{
			NumShards: 1,
			ZoneName:  "",
			RegionConfigs: []*AdvancedRegionConfig{
				{
					ElectableSpecs:      &Specs{},
					BackingProviderName: "",
					Priority:            pointer.MakePtr(7),
					ProviderName:        "",
				},
			},
		})
	}

	if len(deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs) == 0 {
		deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs = append(deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs, &AdvancedRegionConfig{
			ElectableSpecs:      &Specs{},
			BackingProviderName: "",
			Priority:            pointer.MakePtr(7),
			ProviderName:        "",
			RegionName:          "",
		})
	}
}

func (c *AtlasDeployment) WithName(name string) *AtlasDeployment {
	c.Name = name
	return c
}

func (c *AtlasDeployment) WithAtlasName(name string) *AtlasDeployment {
	c.Spec.DeploymentSpec.Name = name
	return c
}

func (c *AtlasDeployment) WithProjectName(projectName string) *AtlasDeployment {
	c.Spec.Project = common.ResourceRefNamespaced{Name: projectName}
	return c
}

func (c *AtlasDeployment) WithProviderName(name provider.ProviderName) *AtlasDeployment {
	addReplicaIfNotAdded(c)
	c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ProviderName = string(name)
	return c
}

func (c *AtlasDeployment) WithRegionName(name string) *AtlasDeployment {
	addReplicaIfNotAdded(c)
	c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].RegionName = name
	return c
}

func (c *AtlasDeployment) WithBackupScheduleRef(ref common.ResourceRefNamespaced) *AtlasDeployment {
	c.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
	c.Spec.BackupScheduleRef = ref
	return c
}

func (c *AtlasDeployment) WithDiskSizeGB(size int) *AtlasDeployment {
	c.Spec.DeploymentSpec.DiskSizeGB = &size
	return c
}

func (c *AtlasDeployment) WithAutoscalingDisabled() *AtlasDeployment {
	addReplicaIfNotAdded(c)
	c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].AutoScaling = nil
	return c
}

func (c *AtlasDeployment) WithInstanceSize(name string) *AtlasDeployment {
	addReplicaIfNotAdded(c)
	c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = name
	return c
}

func (c *AtlasDeployment) WithBackingProvider(name string) *AtlasDeployment {
	addReplicaIfNotAdded(c)
	c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].BackingProviderName = name
	return c
}

func (c *AtlasDeployment) WithSearchNodes(instanceSize string, count uint8) *AtlasDeployment {
	c.Spec.DeploymentSpec.SearchNodes = []SearchNode{
		{
			InstanceSize: instanceSize,
			NodeCount:    count,
		},
	}
	return c
}

// Lightweight makes the deployment work with small shared instance M2. This is useful for non-deployment tests (e.g.
// database users) and saves some money for the company.
func (c *AtlasDeployment) Lightweight() *AtlasDeployment {
	c.WithInstanceSize("M2")
	// M2 is restricted to some set of regions only - we need to ensure them
	switch provider.ProviderName(c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ProviderName) {
	case provider.ProviderAWS:
		{
			c.WithRegionName("US_EAST_1")
		}
	case provider.ProviderAzure:
		{
			c.WithRegionName("US_EAST_2")
		}
	case provider.ProviderGCP:
		{
			c.WithRegionName("CENTRAL_US")
		}
	}
	// Changing provider to tenant as this is shared now
	c.WithBackingProvider(c.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ProviderName)
	c.WithProviderName(provider.ProviderTenant)
	return c
}

func DefaultGCPDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-gcp-k8s", "test-deployment-gcp").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderGCP).
		WithBackingProvider(string(provider.ProviderGCP)).
		WithRegionName("EASTERN_US")
}

func DefaultAWSDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-aws-k8s", "test-deployment-aws").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderAWS).
		WithBackingProvider(string(provider.ProviderAWS)).
		WithRegionName("US_EAST_1")
}

func DefaultAzureDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-azure-k8s", "test-deployment-azure").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderAzure).
		WithBackingProvider(string(provider.ProviderAzure)).
		WithRegionName("EUROPE_NORTH")
}

func DefaultAwsAdvancedDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(
		namespace,
		"test-deployment-advanced-k8s",
		"test-deployment-advanced",
	).WithProjectName(projectName)
}

func NewDefaultAWSServerlessInstance(namespace, projectName string) *AtlasDeployment {
	return newServerlessInstance(
		namespace,
		"test-serverless-instance-k8s",
		"test-serverless-instance",
		"AWS",
		"US_EAST_1",
	).WithProjectName(projectName)
}

func (c *AtlasDeployment) AtlasName() string {
	if c.Spec.DeploymentSpec != nil {
		return c.Spec.DeploymentSpec.Name
	}
	if c.Spec.ServerlessSpec != nil {
		return c.Spec.ServerlessSpec.Name
	}
	return ""
}
