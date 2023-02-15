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
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"go.mongodb.org/atlas/mongodbatlas"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
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

	// Configuration for the normal (v1) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters/
	// +optional
	DeploymentSpec *DeploymentSpec `json:"deploymentSpec,omitempty"`

	// Configuration for the advanced (v1.5) deployment API https://www.mongodb.com/docs/atlas/reference/api/clusters-advanced/
	// +optional
	AdvancedDeploymentSpec *AdvancedDeploymentSpec `json:"advancedDeploymentSpec,omitempty"`

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

type DeploymentSpec struct {
	// Collection of settings that configures auto-scaling information for the deployment.
	// If you specify the autoScaling object, you must also specify the providerSettings.autoScaling object.
	// +optional
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`

	// Configuration of BI Connector for Atlas on this deployment.
	// The MongoDB Connector for Business Intelligence for Atlas (BI Connector) is only available for M10 and larger deployments.
	// +optional
	BIConnector *BiConnectorSpec `json:"biConnector,omitempty"`

	// Type of the deployment that you want to create.
	// The parameter is required if replicationSpecs are set or if Global Deployments are deployed.
	// +kubebuilder:validation:Enum=REPLICASET;SHARDED;GEOSHARDED
	// +optional
	ClusterType DeploymentType `json:"clusterType,omitempty"`

	// Capacity, in gigabytes, of the host's root volume.
	// Increase this number to add capacity, up to a maximum possible value of 4096 (i.e., 4 TB).
	// This value must be a positive integer.
	// The parameter is required if replicationSpecs are configured.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4096
	// +optional
	DiskSizeGB *int `json:"diskSizeGB,omitempty"` // TODO: may cause issues due to mongodb/go-client-mongodb-atlas#140

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

	// Name of the deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name string `json:"name"`

	// Positive integer that specifies the number of shards to deploy for a sharded deployment.
	// The parameter is required if replicationSpecs are configured
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=50
	// +optional
	NumShards *int `json:"numShards,omitempty"`

	// Flag that indicates whether the deployment should be paused.
	Paused *bool `json:"paused,omitempty"`

	// Flag that indicates the deployment uses continuous cloud backups.
	// +optional
	PitEnabled *bool `json:"pitEnabled,omitempty"`

	// Applicable only for M10+ deployments.
	// Flag that indicates if the deployment uses Cloud Backups for backups.
	// +optional
	ProviderBackupEnabled *bool `json:"providerBackupEnabled,omitempty"`

	// Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.
	ProviderSettings *ProviderSettingsSpec `json:"providerSettings"`

	// Configuration for deployment regions.
	// +optional
	ReplicationSpecs []ReplicationSpec `json:"replicationSpecs,omitempty"`
	// +optional
	CustomZoneMapping []CustomZoneMapping `json:"customZoneMapping,omitempty"`
	// +optional
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`
}

type AdvancedDeploymentSpec struct {
	BackupEnabled            *bool              `json:"backupEnabled,omitempty"`
	BiConnector              *BiConnectorSpec   `json:"biConnector,omitempty"`
	ClusterType              string             `json:"clusterType,omitempty"`
	DiskSizeGB               *int               `json:"diskSizeGB,omitempty"`
	EncryptionAtRestProvider string             `json:"encryptionAtRestProvider,omitempty"`
	Labels                   []common.LabelSpec `json:"labels,omitempty"`
	MongoDBMajorVersion      string             `json:"mongoDBMajorVersion,omitempty"`
	MongoDBVersion           string             `json:"mongoDBVersion,omitempty"`
	// Name of the advanced deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name                 string                     `json:"name,omitempty"`
	Paused               *bool                      `json:"paused,omitempty"`
	PitEnabled           *bool                      `json:"pitEnabled,omitempty"`
	ReplicationSpecs     []*AdvancedReplicationSpec `json:"replicationSpecs,omitempty"`
	RootCertType         string                     `json:"rootCertType,omitempty"`
	VersionReleaseSystem string                     `json:"versionReleaseSystem,omitempty"`
	// +optional
	CustomZoneMapping []CustomZoneMapping `json:"customZoneMapping,omitempty"`
	// +optional
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`
}

// ToAtlas converts the AdvancedDeploymentSpec to native Atlas client ToAtlas format.
func (s *AdvancedDeploymentSpec) ToAtlas() (*mongodbatlas.AdvancedCluster, error) {
	result := &mongodbatlas.AdvancedCluster{}
	err := compat.JSONCopy(result, s)
	return result, err
}

// ServerlessSpec defines the desired state of Atlas Serverless Instance
type ServerlessSpec struct {
	// Name of the serverless deployment as it appears in Atlas.
	// After Atlas creates the deployment, you can't change its name.
	// Can only contain ASCII letters, numbers, and hyphens.
	// +kubebuilder:validation:Pattern:=^[a-zA-Z0-9][a-zA-Z0-9-]*$
	Name string `json:"name"`
	// Configuration for the provisioned hosts on which MongoDB runs. The available options are specific to the cloud service provider.
	ProviderSettings *ProviderSettingsSpec `json:"providerSettings"`

	PrivateEndpoints []ServerlessPrivateEndpoint `json:"privateEndpoints,omitempty"`
}

// BiConnector specifies BI Connector for Atlas configuration on this deployment.
type BiConnector struct {
	Enabled        *bool  `json:"enabled,omitempty"`
	ReadPreference string `json:"readPreference,omitempty"`
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
	NumShards     int                     `json:"numShards,omitempty"`
	ZoneName      string                  `json:"zoneName,omitempty"`
	RegionConfigs []*AdvancedRegionConfig `json:"regionConfigs,omitempty"`
}

type AdvancedRegionConfig struct {
	AnalyticsSpecs      *Specs                   `json:"analyticsSpecs,omitempty"`
	ElectableSpecs      *Specs                   `json:"electableSpecs,omitempty"`
	ReadOnlySpecs       *Specs                   `json:"readOnlySpecs,omitempty"`
	AutoScaling         *AdvancedAutoScalingSpec `json:"autoScaling,omitempty"`
	BackingProviderName string                   `json:"backingProviderName,omitempty"`
	Priority            *int                     `json:"priority,omitempty"`
	ProviderName        string                   `json:"providerName,omitempty"`
	RegionName          string                   `json:"regionName,omitempty"`
}

type Specs struct {
	DiskIOPS      *int64 `json:"diskIOPS,omitempty"`
	EbsVolumeType string `json:"ebsVolumeType,omitempty"`
	InstanceSize  string `json:"instanceSize,omitempty"`
	NodeCount     *int   `json:"nodeCount,omitempty"`
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

func (specArgs ProcessArgs) ToAtlas() (*mongodbatlas.ProcessArgs, error) {
	result := &mongodbatlas.ProcessArgs{}
	if err := convertOplogMinRetentionHours(&specArgs, result); err != nil {
		return nil, err
	}

	err := compat.JSONCopy(result, specArgs)
	return result, err
}

func convertOplogMinRetentionHours(specArgs *ProcessArgs, atlasArgs *mongodbatlas.ProcessArgs) error {
	if specArgs != nil && specArgs.OplogMinRetentionHours != "" {
		OplogMinRetentionHours, err := strconv.ParseFloat(specArgs.OplogMinRetentionHours, 64)
		if err != nil {
			return err
		}

		atlasArgs.OplogMinRetentionHours = &OplogMinRetentionHours
		specArgs.OplogMinRetentionHours = ""
	}

	return nil
}

func (specArgs ProcessArgs) IsEqual(newArgs interface{}) bool {
	specV := reflect.ValueOf(specArgs)
	newV := reflect.Indirect(reflect.ValueOf(newArgs))
	typeOfSpec := specV.Type()
	for i := 0; i < specV.NumField(); i++ {
		name := typeOfSpec.Field(i).Name
		specValue := specV.FieldByName(name)
		newValue := newV.FieldByName(name)

		if specValue.IsZero() {
			continue
		}
		if newValue.IsZero() {
			return false
		}

		if specValue.Kind() == reflect.Ptr {
			if specValue.IsNil() {
				continue
			}
			specValue = specValue.Elem()
		}

		if newValue.Kind() == reflect.Ptr {
			if newValue.IsNil() {
				return false
			}
			newValue = newValue.Elem()
		}

		if stringValue(specValue.Interface()) != stringValue(newValue.Interface()) {
			return false
		}
	}

	return true
}

var TrailingZerosRegex = regexp.MustCompile(`\.[0]*$`)

func stringValue(v interface{}) string {
	return TrailingZerosRegex.ReplaceAllString(fmt.Sprint(v), "")
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

// ProviderSettingsSpec configuration for the provisioned servers on which MongoDB runs. The available options are specific to the cloud service provider.
type ProviderSettingsSpec struct {
	// Cloud service provider on which the host for a multi-tenant deployment is provisioned.
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

	// Flag that indicates whether the Amazon EBS encryption feature encrypts the host's root volume for both data at rest within the volume and for data moving between the volume and the deployment.
	// +optional
	EncryptEBSVolume *bool `json:"encryptEBSVolume,omitempty"`

	// Atlas provides different deployment tiers, each with a default storage capacity and RAM size. The deployment you select is used for all the data-bearing hosts in your deployment tier.
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
	// +kubebuilder:validation:Enum=STANDARD;PROVISIONED
	VolumeType string `json:"volumeType,omitempty"`

	// Range of instance sizes to which your deployment can scale.
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`
}

// ReplicationSpec represents a configuration for deployment regions
type ReplicationSpec struct {
	// Number of shards to deploy in each specified zone.
	// The default value is 1.
	NumShards *int64 `json:"numShards,omitempty"`

	// Name for the zone in a Global Deployment.
	// Don't provide this value if deploymentType is not GEOSHARDED.
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

// Check compatibility with library type.
var _ = RegionsConfig(mongodbatlas.RegionsConfig{})

// Deployment converts the Spec to native Atlas client format.
func (spec *AtlasDeploymentSpec) Deployment() (*mongodbatlas.Cluster, error) {
	result := &mongodbatlas.Cluster{}
	err := compat.JSONCopy(result, *spec.DeploymentSpec)

	if result.AutoScaling != nil {
		result.AutoScaling.AutoIndexingEnabled = nil
	}

	if result.ProviderSettings != nil && result.ProviderSettings.AutoScaling != nil {
		result.ProviderSettings.AutoScaling.AutoIndexingEnabled = nil
	}

	return result, err
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// AtlasDeployment is the Schema for the atlasdeployments API
type AtlasDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasDeploymentSpec          `json:"spec,omitempty"`
	Status status.AtlasDeploymentStatus `json:"status,omitempty"`
}

func (c *AtlasDeployment) GetDeploymentName() string {
	if c.IsAdvancedDeployment() {
		return c.Spec.AdvancedDeploymentSpec.Name
	}
	if c.IsServerless() {
		return c.Spec.ServerlessSpec.Name
	}
	return c.Spec.DeploymentSpec.Name
}

// IsServerless returns true if the AtlasDeployment is configured to be a serverless instance
func (c *AtlasDeployment) IsServerless() bool {
	return c.Spec.ServerlessSpec != nil
}

// IsAdvancedDeployment returns true if the AtlasDeployment is configured to be an advanced deployment.
func (c *AtlasDeployment) IsAdvancedDeployment() bool {
	return c.Spec.AdvancedDeploymentSpec != nil
}

// +kubebuilder:object:root=true

// AtlasDeploymentList contains a list of AtlasDeployment
type AtlasDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasDeployment `json:"items"`
}

func (c AtlasDeployment) AtlasProjectObjectKey() client.ObjectKey {
	ns := c.Namespace
	if c.Spec.Project.Namespace != "" {
		ns = c.Spec.Project.Namespace
	}
	return kube.ObjectKey(ns, c.Spec.Project.Name)
}

func (c *AtlasDeployment) GetStatus() status.Status {
	return c.Status
}

func (c *AtlasDeployment) UpdateStatus(conditions []status.Condition, options ...status.Option) {
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
			DeploymentSpec: &DeploymentSpec{
				Name:             nameInAtlas,
				ProviderSettings: &ProviderSettingsSpec{InstanceSizeName: "M10"},
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
				ProviderSettings: &ProviderSettingsSpec{
					BackingProviderName: backingProviderName,
					ProviderName:        "SERVERLESS",
					RegionName:          regionName,
				},
			},
		},
	}
}

func NewAwsAdvancedDeployment(namespace, name, nameInAtlas string) *AtlasDeployment {
	return newAwsAdvancedDeployment(namespace, name, nameInAtlas, "M10", "AWS", "US_EAST_1", 3)
}

func newAwsAdvancedDeployment(namespace, name, nameInAtlas, instanceSize, providerName, regionName string, nodeCount int) *AtlasDeployment {
	priority := 7
	return &AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: AtlasDeploymentSpec{
			AdvancedDeploymentSpec: &AdvancedDeploymentSpec{
				Name:        nameInAtlas,
				ClusterType: string(TypeReplicaSet),
				ReplicationSpecs: []*AdvancedReplicationSpec{
					{
						RegionConfigs: []*AdvancedRegionConfig{
							{
								Priority: &priority,
								ElectableSpecs: &Specs{
									InstanceSize: instanceSize,
									NodeCount:    &nodeCount,
								},
								ProviderName: providerName,
								RegionName:   regionName,
							},
						},
					}},
			},
		},
	}
}

func (c *AtlasDeployment) WithName(name string) *AtlasDeployment {
	c.Spec.DeploymentSpec.Name = name
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
	c.Spec.DeploymentSpec.ProviderSettings.ProviderName = name
	return c
}

func (c *AtlasDeployment) WithRegionName(name string) *AtlasDeployment {
	c.Spec.DeploymentSpec.ProviderSettings.RegionName = name
	return c
}

func (c *AtlasDeployment) WithBackupScheduleRef(ref common.ResourceRefNamespaced) *AtlasDeployment {
	t := true
	c.Spec.DeploymentSpec.ProviderBackupEnabled = &t
	c.Spec.BackupScheduleRef = ref
	return c
}

func (c *AtlasDeployment) WithDiskSizeGB(size int) *AtlasDeployment {
	c.Spec.DeploymentSpec.DiskSizeGB = &size
	return c
}

func (c *AtlasDeployment) WithAutoscalingDisabled() *AtlasDeployment {
	f := false
	c.Spec.DeploymentSpec.AutoScaling = &AutoScalingSpec{
		DiskGBEnabled: &f,
		Compute: &ComputeSpec{
			Enabled:          &f,
			ScaleDownEnabled: &f,
			MinInstanceSize:  "",
			MaxInstanceSize:  "",
		},
	}
	return c
}

func (c *AtlasDeployment) WithInstanceSize(name string) *AtlasDeployment {
	c.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = name
	return c
}
func (c *AtlasDeployment) WithBackingProvider(name string) *AtlasDeployment {
	c.Spec.DeploymentSpec.ProviderSettings.BackingProviderName = name
	return c
}

// Lightweight makes the deployment work with small shared instance M2. This is useful for non-deployment tests (e.g.
// database users) and saves some money for the company.
func (c *AtlasDeployment) Lightweight() *AtlasDeployment {
	c.WithInstanceSize("M2")
	// M2 is restricted to some set of regions only - we need to ensure them
	switch c.Spec.DeploymentSpec.ProviderSettings.ProviderName {
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
	c.WithBackingProvider(string(c.Spec.DeploymentSpec.ProviderSettings.ProviderName))
	c.WithProviderName(provider.ProviderTenant)
	return c
}

func DefaultGCPDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-gcp-k8s", "test-deployment-gcp").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderGCP).
		WithRegionName("EASTERN_US")
}

func DefaultAWSDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-aws-k8s", "test-deployment-aws").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderAWS).
		WithRegionName("US_WEST_2")
}

func DefaultAzureDeployment(namespace, projectName string) *AtlasDeployment {
	return NewDeployment(namespace, "test-deployment-azure-k8s", "test-deployment-azure").
		WithProjectName(projectName).
		WithProviderName(provider.ProviderAzure).
		WithRegionName("EUROPE_NORTH")
}

func DefaultAwsAdvancedDeployment(namespace, projectName string) *AtlasDeployment {
	return NewAwsAdvancedDeployment(namespace, "test-deployment-advanced-k8s", "test-deployment-advanced").WithProjectName(projectName)
}

func NewDefaultAWSServerlessInstance(namespace, projectName string) *AtlasDeployment {
	return newServerlessInstance(namespace, "test-serverless-instance-k8s", "test-serverless-instance", "AWS", "US_EAST_1").WithProjectName(projectName)
}

func (c *AtlasDeployment) AtlasName() string {
	if c.Spec.DeploymentSpec != nil {
		return c.Spec.DeploymentSpec.Name
	}
	if c.Spec.AdvancedDeploymentSpec != nil {
		return c.Spec.AdvancedDeploymentSpec.Name
	}
	if c.Spec.ServerlessSpec != nil {
		return c.Spec.ServerlessSpec.Name
	}
	return ""
}
