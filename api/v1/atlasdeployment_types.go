// Copyright 2020 MongoDB Inc
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

package v1

import (
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
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

// AtlasDeploymentSpec defines the target state of AtlasDeployment.
// Only one of DeploymentSpec, AdvancedDeploymentSpec and ServerlessSpec should be defined.
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && !has(self.projectRef)) || (!has(self.externalProjectRef) && has(self.projectRef))",message="must define only one project reference through externalProjectRef or projectRef"
// +kubebuilder:validation:XValidation:rule="(has(self.externalProjectRef) && has(self.connectionSecret)) || !has(self.externalProjectRef)",message="must define a local connection secret when referencing an external project"
// +kubebuilder:validation:XValidation:rule="!has(self.serverlessSpec) || (oldSelf.hasValue() && oldSelf.value().serverlessSpec != null)",optionalOldSelf=true,message="serverlessSpec cannot be added - serverless instances are deprecated",fieldPath=.serverlessSpec
type AtlasDeploymentSpec struct {
	// ProjectReference is the dual external or kubernetes reference with access credentials
	ProjectDualReference `json:",inline"`

	//  upgradeToDedicated, when set to true, triggers the migration from a Flex to a
	//  Dedicated cluster. The user MUST provide the new dedicated cluster configuration.
	//  This flag is ignored if the cluster is already dedicated.
	// +optional
	UpgradeToDedicated bool `json:"upgradeToDedicated,omitempty"`

	// Configuration for the advanced (v1.5) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters/
	// +optional
	DeploymentSpec *AdvancedDeploymentSpec `json:"deploymentSpec,omitempty"`

	// Reference to the backup schedule for the AtlasDeployment.
	// +optional
	BackupScheduleRef common.ResourceRefNamespaced `json:"backupRef"`

	// Configuration for the serverless deployment API. https://www.mongodb.com/docs/atlas/reference/api/serverless-instances/
	// DEPRECATED: Serverless instances are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.
	// +optional
	ServerlessSpec *ServerlessSpec `json:"serverlessSpec,omitempty"`

	// ProcessArgs allows modification of Advanced Configuration Options.
	// +optional
	ProcessArgs *ProcessArgs `json:"processArgs,omitempty"`

	// Configuration for the Flex cluster API. https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Flex-Clusters
	// +optional
	FlexSpec *FlexSpec `json:"flexSpec,omitempty"`
}

type SearchNode struct {
	// Hardware specification for the Search Node instance sizes.
	// See https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-creategroupclustersearchdeployment#operation-creategroupclustersearchdeployment-body-application-vnd-atlas-2024-05-30-json-specs-instancesize for available values
	InstanceSize string `json:"instanceSize,omitempty"`
	// Number of Search Nodes in the cluster.
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=32
	NodeCount uint8 `json:"nodeCount,omitempty"`
}

type AdvancedDeploymentSpec struct {
	// Flag that indicates if the deployment uses Cloud Backups for backups.
	// Applicable only for M10+ deployments.
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
	// DEPRECATED: Cluster labels are deprecated and will be removed in a future release. We strongly recommend that you use Resource Tags instead.
	// +optional
	Labels []common.LabelSpec `json:"labels,omitempty"`
	// MongoDB major version of the cluster. Set to the binary major version.
	// +optional
	MongoDBMajorVersion string `json:"mongoDBMajorVersion,omitempty"`
	// Version of MongoDB that the cluster runs.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`
	// Name of the advanced deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Name cannot be modified after deployment creation"
	Name string `json:"name,omitempty"`
	// Flag that indicates whether the deployment should be paused.
	Paused *bool `json:"paused,omitempty"`
	// Flag that indicates the deployment uses continuous cloud backups.
	// +optional
	PitEnabled *bool `json:"pitEnabled,omitempty"`
	// Configuration for deployment regions.
	// +optional
	ReplicationSpecs []*AdvancedReplicationSpec `json:"replicationSpecs,omitempty"`
	// Root Certificate Authority that MongoDB Atlas cluster uses.
	// +optional
	RootCertType string `json:"rootCertType,omitempty"`
	// Key-value pairs for resource tagging.
	// +kubebuilder:validation:MaxItems=50
	// +optional
	Tags []*TagSpec `json:"tags,omitempty"`
	// Method by which the cluster maintains the MongoDB versions.
	// If value is CONTINUOUS, you must not specify mongoDBMajorVersion.
	// +optional
	VersionReleaseSystem string `json:"versionReleaseSystem,omitempty"`
	// List that contains Global Cluster parameters that map zones to geographic regions.
	// +optional
	CustomZoneMapping []CustomZoneMapping `json:"customZoneMapping,omitempty"`
	// List that contains information to create a managed namespace in a specified Global Cluster to create.
	// +optional
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`
	// Flag that indicates whether termination protection is enabled on the cluster. If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
	// +kubebuilder:default:=false
	TerminationProtectionEnabled bool `json:"terminationProtectionEnabled,omitempty"`
	// Settings for Search Nodes for the cluster. Currently, at most one search node configuration may be defined.
	// +kubebuilder:validation:MaxItems=1
	// +optional
	SearchNodes []SearchNode `json:"searchNodes,omitempty"`
	// An array of SearchIndex objects with fields that describe the search index.
	// +optional
	SearchIndexes []SearchIndex `json:"searchIndexes,omitempty"`
	// Config Server Management Mode for creating or updating a sharded cluster.
	// +kubebuilder:validation:Enum=ATLAS_MANAGED;FIXED_TO_DEDICATED
	// +optional
	ConfigServerManagementMode string `json:"configServerManagementMode,omitempty"`
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

// ServerlessSpec defines the target state of Atlas Serverless Instance.
// DEPRECATED: Serverless instances are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.
type ServerlessSpec struct {
	// Name of the serverless deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name string `json:"name"`
	// Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.
	ProviderSettings *ServerlessProviderSettingsSpec `json:"providerSettings"`
	// List that contains the private endpoint configurations for the Serverless instance.
	// DEPRECATED: Serverless private endpoints are deprecated. See https://dochub.mongodb.org/core/atlas-flex-migration for details.
	PrivateEndpoints []ServerlessPrivateEndpoint `json:"privateEndpoints,omitempty"`
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

// BiConnector specifies Business Intelligence Connector for Atlas configuration on this deployment.
type BiConnector struct {
	// Flag that indicates whether MongoDB Connector for Business Intelligence is enabled on the specified cluster.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// Data source node designated for the MongoDB Connector for Business Intelligence on MongoDB Cloud.
	// The MongoDB Connector for Business Intelligence on MongoDB Cloud reads data from the primary, secondary, or analytics node based on your read preferences.
	// +optional
	ReadPreference string `json:"readPreference,omitempty"`
}

// TagSpec holds a key-value pair for resource tagging on this deployment.
type TagSpec struct {
	// Constant that defines the set of the tag.
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9 @_.+`;`-]*$
	Key string `json:"key"`
	// Variable that belongs to the set of the tag.
	// +kubebuilder:validation:MaxLength:=255
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9 @_.+`;`-]*$
	Value string `json:"value"`
}

type ServerlessBackupOptions struct {
	// ServerlessContinuousBackupEnabled indicates whether the cluster uses continuous cloud backups.
	// DEPRECATED: Serverless instances are deprecated, and no longer support continuous backup. See https://dochub.mongodb.org/core/atlas-flex-migration for details.
	// +kubebuilder:default:=true
	ServerlessContinuousBackupEnabled bool `json:"serverlessContinuousBackupEnabled,omitempty"`
}

// ConnectionStrings is a collection of Uniform Resource Locators that point to the MongoDB database.
type ConnectionStrings struct {
	// Public connection string that you can use to connect to this cluster. This connection string uses the mongodb:// protocol.
	Standard string `json:"standard,omitempty"`
	// Public connection string that you can use to connect to this cluster. This connection string uses the mongodb+srv:// protocol.
	StandardSrv string `json:"standardSrv,omitempty"`
	// List of private endpoint-aware connection strings that you can use to connect to this cluster through a private endpoint.
	// This parameter returns only if you deployed a private endpoint to all regions to which you deployed this clusters' nodes.
	PrivateEndpoint []PrivateEndpointSpec `json:"privateEndpoint,omitempty"`
	// Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink.
	// Each key identifies an Amazon Web Services (AWS) interface endpoint.
	// Each value identifies the related mongodb:// connection string that you use to connect to MongoDB Cloud through the interface endpoint that the key names.
	AwsPrivateLink map[string]string `json:"awsPrivateLink,omitempty"`
	// Private endpoint-aware connection strings that use AWS-hosted clusters with Amazon Web Services (AWS) PrivateLink.
	// Each key identifies an Amazon Web Services (AWS) interface endpoint.
	// Each value identifies the related mongodb:// connection string that you use to connect to Atlas through the interface endpoint that the key names.
	// If the cluster uses an optimized connection string, awsPrivateLinkSrv contains the optimized connection string.
	// If the cluster has the non-optimized (legacy) connection string, awsPrivateLinkSrv contains the non-optimized connection string even if an optimized connection string is also present.
	AwsPrivateLinkSrv map[string]string `json:"awsPrivateLinkSrv,omitempty"`
	// Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster.
	// This connection string uses the mongodb+srv:// protocol. The resource returns this parameter once someone creates a network peering connection to this cluster.
	// This protocol tells the application to look up the host seed list in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster.
	// If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the URI if the nodes change.
	// Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private.
	// For Amazon Web Services (AWS) clusters, this resource returns this parameter only if you enable custom DNS.
	Private string `json:"private,omitempty"`
	// Network peering connection strings for each interface Virtual Private Cloud (VPC) endpoint that you configured to connect to this cluster.
	// This connection string uses the mongodb+srv:// protocol. The resource returns this parameter when someone creates a network peering connection to this cluster.
	// This protocol tells the application to look up the host seed list in the Domain Name System (DNS).
	// This list synchronizes with the nodes in a cluster. If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change.
	// Use this Uniform Resource Identifier (URI) format if your driver supports it. If it doesn't, use connectionStrings.private.
	// For Amazon Web Services (AWS) clusters, this parameter returns only if you enable custom DNS.
	PrivateSrv string `json:"privateSrv,omitempty"`
}

// PrivateEndpointSpec connection strings. Each object describes the connection strings
// you can use to connect to this deployment through a private endpoint.
// Atlas returns this parameter only if you deployed a private endpoint to all regions
// to which you deployed this deployment's nodes.
type PrivateEndpointSpec struct {
	// Private endpoint-aware connection string that uses the mongodb:// protocol to connect to MongoDB Cloud through a private endpoint.
	ConnectionString string `json:"connectionString,omitempty"`
	// List that contains the private endpoints through which you connect to MongoDB Atlas.
	Endpoints []EndpointSpec `json:"endpoints,omitempty"`
	// Private endpoint-aware connection string that uses the mongodb+srv:// protocol to connect to MongoDB Cloud through a private endpoint.
	// The mongodb+srv protocol tells the driver to look up the seed list of hosts in the Domain Name System (DNS). This list synchronizes with the nodes in a cluster.
	// If the connection string uses this Uniform Resource Identifier (URI) format, you don't need to append the seed list or change the Uniform Resource Identifier (URI) if the nodes change.
	// Use this Uniform Resource Identifier (URI) format if your application supports it.
	SRVConnectionString string `json:"srvConnectionString,omitempty"`
	// MongoDB process type to which your application connects.
	// Use MONGOD for replica sets and MONGOS for sharded clusters.
	Type string `json:"type,omitempty"`
}

// EndpointSpec through which you connect to Atlas.
type EndpointSpec struct {
	// Unique string that the cloud provider uses to identify the private endpoint.
	EndpointID string `json:"endpointId,omitempty"`
	// Cloud provider in which MongoDB Cloud deploys the private endpoint.
	ProviderName string `json:"providerName,omitempty"`
	// Region where the private endpoint is deployed.
	Region string `json:"region,omitempty"`
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
	// Hardware specifications for analytics nodes deployed in the region.
	AnalyticsSpecs *Specs `json:"analyticsSpecs,omitempty"`
	// Hardware specifications for nodes deployed in the region.
	ElectableSpecs *Specs `json:"electableSpecs,omitempty"`
	// Hardware specifications for read only nodes deployed in the region.
	ReadOnlySpecs *Specs `json:"readOnlySpecs,omitempty"`
	// Options that determine how this cluster handles resource scaling.
	AutoScaling *AdvancedAutoScalingSpec `json:"autoScaling,omitempty"`
	// Cloud service provider on which the host for a multi-tenant deployment is provisioned.
	// This setting only works when "providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
	// Otherwise, it should be equal to the "providerName" value.
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
	// The instance size you select applies to all the data-bearing hosts in your instance size.
	InstanceSize string `json:"instanceSize,omitempty"`
	// Number of nodes of the given type for MongoDB Cloud to deploy to the region.
	NodeCount *int `json:"nodeCount,omitempty"`
}

// AutoScalingSpec configures your deployment to automatically scale its storage
type AutoScalingSpec struct {
	// Flag that indicates whether autopilot mode for Performance Advisor is enabled.
	// The default is false.
	// DEPRECATED: This flag is no longer supported.
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
	// Flag that indicates whether this cluster enables disk auto-scaling.
	// The maximum memory allowed for the selected cluster tier and the oplog size can limit storage auto-scaling.
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
	// String that indicates the default level of acknowledgment requested from MongoDB for read operations set for this cluster.
	DefaultReadConcern string `json:"defaultReadConcern,omitempty"`
	// String that indicates the default level of acknowledgment requested from MongoDB for write operations set for this cluster.
	DefaultWriteConcern string `json:"defaultWriteConcern,omitempty"`
	// String that indicates the minimum TLS version that the cluster accepts for incoming connections.
	// Clusters using TLS 1.0 or 1.1 should consider setting TLS 1.2 as the minimum TLS protocol version.
	MinimumEnabledTLSProtocol string `json:"minimumEnabledTlsProtocol,omitempty"`
	// Flag that indicates whether to fail the operation and return an error when you insert or update documents where all indexed entries exceed 1024 bytes.
	// If you set this to false, mongod writes documents that exceed this limit, but doesn't index them.
	FailIndexKeyTooLong *bool `json:"failIndexKeyTooLong,omitempty"`
	// Flag that indicates whether the cluster allows execution of operations that perform server-side executions of JavaScript.
	JavascriptEnabled *bool `json:"javascriptEnabled,omitempty"`
	// Flag that indicates whether the cluster disables executing any query that requires a collection scan to return results.
	NoTableScan *bool `json:"noTableScan,omitempty"`
	// Number that indicates the storage limit of a cluster's oplog expressed in megabytes.
	// A value of null indicates that the cluster uses the default oplog size that Atlas calculates.
	OplogSizeMB *int64 `json:"oplogSizeMB,omitempty"`
	// Number that indicates the interval in seconds at which the mongosqld process re-samples data to create its relational schema.
	SampleSizeBIConnector *int64 `json:"sampleSizeBIConnector,omitempty"`
	// Number that indicates the documents per database to sample when gathering schema information.
	SampleRefreshIntervalBIConnector *int64 `json:"sampleRefreshIntervalBIConnector,omitempty"`
	// Minimum retention window for cluster's oplog expressed in hours. A value of null indicates that the cluster uses the default minimum oplog window that MongoDB Cloud calculates.
	OplogMinRetentionHours string `json:"oplogMinRetentionHours,omitempty"`
}

// BiConnectorSpec specifies BI Connector for Atlas configuration on this deployment
type BiConnectorSpec struct {
	// Flag that indicates whether the Business Intelligence Connector for Atlas is enabled on the deployment.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Source from which the BI Connector for Atlas reads data. Each BI Connector for Atlas read preference contains a distinct combination of readPreference and readPreferenceTags options.
	// +optional
	ReadPreference string `json:"readPreference,omitempty"`
}

// ServerlessProviderSettingsSpec configuration for the provisioned servers on which MongoDB runs. The available options are specific to the cloud service provider.
type ServerlessProviderSettingsSpec struct {
	// Cloud service provider on which the host for a multi-tenant deployment is provisioned.
	// This setting only works when "providerSetting.providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +optional
	BackingProviderName string `json:"backingProviderName,omitempty"`

	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// DEPRECATED: The value of this field doesn't take any effect.
	// +optional
	DiskIOPS *int64 `json:"diskIOPS,omitempty"`

	// Type of disk if you selected Azure as your cloud service provider.
	// DEPRECATED: The value of this field doesn't take any effect.
	// +optional
	DiskTypeName string `json:"diskTypeName,omitempty"`

	// Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the deployment.
	// DEPRECATED: The value of this field doesn't take any effect.
	// +optional
	EncryptEBSVolume *bool `json:"encryptEBSVolume,omitempty"`

	// Atlas provides different deployment tiers, each with a default storage capacity and RAM size. The deployment you select is used for all the data-bearing hosts in your deployment tier.
	// DEPRECATED: The value of this field doesn't take any effect.
	// +optional
	InstanceSizeName string `json:"instanceSizeName,omitempty"`

	// Cloud service provider on which Atlas provisions the hosts.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT;SERVERLESS
	ProviderName provider.ProviderName `json:"providerName"`

	// Physical location of your MongoDB deployment.
	// The region you choose can affect network latency for clients accessing your databases.
	// +optional
	RegionName string `json:"regionName,omitempty"`

	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// DEPRECATED: The value of this field doesn't take any effect.
	// +kubebuilder:validation:Enum=STANDARD;PROVISIONED
	VolumeType string `json:"volumeType,omitempty"`

	// Range of instance sizes to which your deployment can scale.
	// DEPRECATED: The value of this field doesn't take any effect.
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`
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
	if c.IsFlex() {
		return c.Spec.FlexSpec.Name
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

func (c *AtlasDeployment) IsFlex() bool {
	return c.Spec.FlexSpec != nil
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
	if c.Spec.ProjectRef.Namespace != "" {
		ns = c.Spec.ProjectRef.Namespace
	}
	return kube.ObjectKey(ns, c.Spec.ProjectRef.Name)
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

func (c *AtlasDeployment) Credentials() *api.LocalObjectReference {
	return c.Spec.ConnectionSecret
}

func (c *AtlasDeployment) ProjectDualRef() *ProjectDualReference {
	return &c.Spec.ProjectDualReference
}

type FlexSpec struct {
	// Human-readable label that identifies the instance.
	// +required
	Name string `json:"name"`

	// List that contains key-value pairs between 1 and 255 characters in length for tagging and categorizing the instance.
	// +kubebuilder:validation:MaxItems=50
	// +optional
	Tags []*TagSpec `json:"tags,omitempty"`

	// Flag that indicates whether termination protection is enabled on the cluster.
	// If set to true, MongoDB Cloud won't delete the cluster. If set to false, MongoDB Cloud will delete the cluster.
	// +kubebuilder:default:=false
	// +optional
	TerminationProtectionEnabled bool `json:"terminationProtectionEnabled,omitempty"`

	// Group of cloud provider settings that configure the provisioned MongoDB flex cluster.
	// +required
	ProviderSettings *FlexProviderSettings `json:"providerSettings"`
}

type FlexProviderSettings struct {
	// Cloud service provider on which MongoDB Atlas provisions the flex cluster.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Backing Provider cannot be modified after cluster creation"
	// +required
	BackingProviderName string `json:"backingProviderName,omitempty"`

	// Human-readable label that identifies the geographic location of your MongoDB flex cluster.
	// The region you choose can affect network latency for clients accessing your databases.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Region Name cannot be modified after cluster creation"
	// +required
	RegionName string `json:"regionName,omitempty"`
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

func newFlexInstance(namespace, name, nameInAtlas, backingProviderName, regionName string) *AtlasDeployment {
	return &AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasDeploymentSpec{
			FlexSpec: &FlexSpec{
				Name: nameInAtlas,
				ProviderSettings: &FlexProviderSettings{
					BackingProviderName: backingProviderName,
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
	if c.Spec.DeploymentSpec != nil {
		c.Spec.DeploymentSpec.Name = name
	} else if c.Spec.ServerlessSpec != nil {
		c.Spec.ServerlessSpec.Name = name
	} else if c.Spec.FlexSpec != nil {
		c.Spec.FlexSpec.Name = name
	}
	return c
}

func (c *AtlasDeployment) WithProjectName(projectName string) *AtlasDeployment {
	c.Spec.ProjectRef = &common.ResourceRefNamespaced{Name: projectName}
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

func (c *AtlasDeployment) WithExternaLProject(projectID, credentialsName string) *AtlasDeployment {
	c.Spec.ProjectRef = nil
	c.Spec.ExternalProjectRef = &ExternalProjectReference{
		ID: projectID,
	}
	c.Spec.ConnectionSecret = &api.LocalObjectReference{
		Name: credentialsName,
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

func NewDefaultAWSFlexInstance(namespace, projectName string) *AtlasDeployment {
	return newFlexInstance(
		namespace,
		"test-flex-instance-aws-k8s",
		"test-flex-instance-aws",
		"AWS",
		"US_EAST_1",
	).WithProjectName(projectName)
}

func NewDefaultAzureFlexInstance(namespace, projectName string) *AtlasDeployment {
	return newFlexInstance(
		namespace,
		"test-flex-instance-az-k8s",
		"test-flex-instance-az",
		"AZURE",
		"US_EAST_2",
	).WithProjectName(projectName)
}

func (c *AtlasDeployment) AtlasName() string {
	if c.Spec.DeploymentSpec != nil {
		return c.Spec.DeploymentSpec.Name
	}
	if c.Spec.ServerlessSpec != nil {
		return c.Spec.ServerlessSpec.Name
	}
	if c.Spec.FlexSpec != nil {
		return c.Spec.FlexSpec.Name
	}
	return ""
}
