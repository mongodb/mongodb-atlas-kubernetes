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
	"github.com/jinzhu/copier"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	SchemeBuilder.Register(&AtlasCluster{}, &AtlasClusterList{})
}

// AtlasClusterSpec defines the desired state of AtlasCluster
type AtlasClusterSpec struct {
	// ConnectionSecret is the name of the Kubernetes Secret which contains the information about the way to connect to
	// Atlas (organization ID, API keys). The default Operator connection configuration will be used if not provided.
	// +optional
	ConnectionSecret *SecretRef `json:"connectionSecretRef,omitempty"`

	// Collection of settings that configures auto-scaling information for the cluster.
	// If you specify the autoScaling object, you must also specify the providerSettings.autoScaling object.
	// +optional
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`

	// Deprecated: do not use.
	// Flag that indicates whether legacy backups have been enabled.
	// Applicable only for M10+ clusters.
	// +optional
	BackupEnabled *bool `json:"backupEnabled,omitempty"`

	// Configuration of BI Connector for Atlas on this cluster.
	// The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger clusters.
	// +optional
	BIConnector *BiConnectorSpec `json:"biConnector,omitempty"`

	// Type of the cluster that you want to create.
	// +kubebuilder:validation:Enum=REPLICASET;SHARDED;GEOSHARDED
	// +optional
	ClusterType string `json:"clusterType,omitempty"`

	// Capacity, in gigabytes, of the host's root volume.
	// Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB).
	// This value must be a positive integer.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4096
	// +optional
	DiskSizeGB *int `json:"diskSizeGB,omitempty"` // TODO: why on earth is this *float64 in mongodbatlas?

	// Cloud service provider that offers Encryption at Rest.
	// +optional
	EncryptionAtRestProvider string `json:"encryptionAtRestProvider,omitempty"`

	// Collection of key-value pairs that tag and categorize the cluster.
	// Each key and value has a maximum length of 255 characters.
	// +optional
	Labels []LabelSpec `json:"labels,omitempty"`

	// Version of the cluster to deploy.
	// +kubebuilder:validation:Enum="3.6";"4.0";"4.2";"4.4"
	// +optional
	MongoDBMajorVersion string `json:"mongoDBMajorVersion,omitempty"`

	// Name of the cluster as it appears in Atlas. After Atlas creates the cluster, you can't change its name.
	Name string `json:"name,omitempty"`

	// Positive integer that specifies the number of shards to deploy for a sharded cluster.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=50
	// +optional
	NumShards *int `json:"numShards,omitempty"`

	// Flag that indicates the cluster uses continuous cloud backups.
	// +optional
	PitEnabled *bool `json:"pitEnabled,omitempty"`

	// Applicable only for M10+ clusters.
	// Flag that indicates if the cluster uses Cloud Backups for backups.
	// +optional
	ProviderBackupEnabled *bool `json:"providerBackupEnabled,omitempty"`

	// Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.
	ProviderSettings *ProviderSettingsSpec `json:"providerSettings,omitempty"`

	// Configuration for cluster regions.
	// +optional
	ReplicationSpecs []ReplicationSpec `json:"replicationSpecs,omitempty"`
}

// AutoScalingSpec configures your cluster to automatically scale its storage
type AutoScalingSpec struct {
	// Flag that indicates whether disk auto-scaling is enabled. The default is true.
	// +optional
	DiskGBEnabled *bool `json:"diskGBEnabled,omitempty"`

	// Collection of settings that configure how a cluster might scale its cluster tier and whether the cluster can scale down.
	// +optional
	Compute *ComputeSpec `json:"compute,omitempty"`
}

// ComputeSpec Specifies whether the cluster automatically scales its cluster tier and whether the cluster can scale down.
type ComputeSpec struct {
	// Flag that indicates whether cluster tier auto-scaling is enabled. The default is false.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Flag that indicates whether the cluster tier may scale down. Atlas requires this parameter if "autoScaling.compute.enabled" : true.
	// +optional
	ScaleDownEnabled *bool `json:"scaleDownEnabled,omitempty"`

	// Minimum instance size to which your cluster can automatically scale (such as M10). Atlas requires this parameter if "autoScaling.compute.scaleDownEnabled" : true.
	// +optional
	MinInstanceSize string `json:"minInstanceSize,omitempty"`

	// Maximum instance size to which your cluster can automatically scale (such as M40). Atlas requires this parameter if "autoScaling.compute.enabled" : true.
	// +optional
	MaxInstanceSize string `json:"maxInstanceSize,omitempty"`
}

// BiConnectorSpec specifies BI Connector for Atlas configuration on this cluster
type BiConnectorSpec struct {
	// Flag that indicates whether or not BI Connector for Atlas is enabled on the cluster.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Source from which the BI Connector for Atlas reads data. Each BI Connector for Atlas read preference contains a distinct combination of readPreference and readPreferenceTags options.
	// +optional
	ReadPreference string `json:"readPreference,omitempty"`
}

// ProviderSettingsSpec configuration for the provisioned servers on which MongoDB runs. The available options are specific to the cloud service provider.
type ProviderSettingsSpec struct {
	// Cloud service provider on which the host for a multi-tenant cluster is provisioned.
	// This setting only works when "providerSetting.providerName" : "TENANT" and "providerSetting.instanceSizeName" : M2 or M5.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE
	// +optional
	BackingProviderName string `json:"backingProviderName,omitempty"`

	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +optional
	DiskIOPS *int64 `json:"diskIOPS,omitempty"`

	// Type of disk if you selected Azure as your cloud service provider.
	// +optional
	DiskTypeName string `json:"diskTypeName,omitempty"`

	// Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the cluster.
	// +optional
	EncryptEBSVolume *bool `json:"encryptEBSVolume,omitempty"`

	// Atlas provides different cluster tiers, each with a default storage capacity and RAM size. The cluster you select is used for all the data-bearing hosts in your cluster tier.
	InstanceSizeName string `json:"instanceSizeName,omitempty"`

	// Cloud service provider on which Atlas provisions the hosts.
	// +kubebuilder:validation:Enum=AWS;GCP;AZURE;TENANT
	ProviderName string `json:"providerName,omitempty"`

	// Physical location of your MongoDB cluster.
	// The region you choose can affect network latency for clients accessing your databases.
	// +optional
	RegionName string `json:"regionName,omitempty"`

	// Disk IOPS setting for AWS storage.
	// Set only if you selected AWS as your cloud service provider.
	// +kubebuilder:validation:Enum=STANDARD;PROVISIONED
	VolumeType string `json:"volumeType,omitempty"`

	// Range of instance sizes to which your cluster can scale.
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`
}

// ReplicationSpec represents a configuration for cluster regions
type ReplicationSpec struct {
	// Number of shards to deploy in each specified zone.
	// The default value is 1.
	NumShards *int64 `json:"numShards,omitempty"`

	// Name for the zone in a Global Cluster.
	// Don't provide this value if clusterType is not GEOSHARDED.
	// +optional
	ZoneName string `json:"zoneName,omitempty"`

	// Configuration for a region.
	// Each regionsConfig object describes the region's priority in elections and the number and type of MongoDB nodes that Atlas deploys to the region.
	// +optional
	RegionsConfig map[string]RegionsConfig `json:"regionsConfig,omitempty"`
}

// RegionsConfig describes the region’s priority in elections and the number and type of MongoDB nodes Atlas deploys to the region.
type RegionsConfig struct {
	// The number of analytics nodes for Atlas to deploy to the region.
	// Analytics nodes are useful for handling analytic data such as reporting queries from BI Connector for Atlas.
	// Analytics nodes are read-only, and can never become the primary.
	// If you do not specify this option, no analytics nodes are deployed to the region.
	// +optional
	AnalyticsNodes *int64 `json:"analyticsNodes,omitempty"`

	// Number of electable nodes for Atlas to deploy to the region.
	// Electable nodes can become the primary and can facilitate local reads.
	// +optional
	ElectableNodes *int64 `json:"electableNodes,omitempty"`

	// Election priority of the region.
	// For regions with only replicationSpecs[n].regionsConfig.<region>.readOnlyNodes, set this value to 0.
	// +optional
	Priority *int64 `json:"priority,omitempty"`

	// Number of read-only nodes for Atlas to deploy to the region.
	// Read-only nodes can never become the primary, but can facilitate local-reads.
	// +optional
	ReadOnlyNodes *int64 `json:"readOnlyNodes,omitempty"`
}

// Cluster converts the Spec to native Atlas client format.
func (spec *AtlasClusterSpec) Cluster() *mongodbatlas.Cluster {
	result := &mongodbatlas.Cluster{}
	err := copier.Copy(result, spec)
	if err != nil {
		panic(err)
	}

	return result
}

// AtlasClusterStatus defines the observed state of AtlasCluster.
type AtlasClusterStatus struct {
	// TODO: this is a stub; will implement the Conditions proposal here

	GroupID   string `json:"groupId,omitempty"`
	ID        string `json:"id,omitempty"`
	StateName string `json:"stateName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AtlasCluster is the Schema for the atlasclusters API
type AtlasCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasClusterSpec   `json:"spec,omitempty"`
	Status AtlasClusterStatus `json:"status,omitempty"`
}

func (c *AtlasCluster) ConnectionSecretObjectKey() *client.ObjectKey {
	if c.Spec.ConnectionSecret != nil {
		key := kube.ObjectKey(c.Namespace, c.Spec.ConnectionSecret.Name)
		return &key
	}
	return nil
}

// +kubebuilder:object:root=true

// AtlasClusterList contains a list of AtlasCluster
type AtlasClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasCluster `json:"items"`
}
