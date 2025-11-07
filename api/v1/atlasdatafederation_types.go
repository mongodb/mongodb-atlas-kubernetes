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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
)

func init() {
	SchemeBuilder.Register(&AtlasDataFederation{}, &AtlasDataFederationList{})
}

type DataFederationSpec struct {
	// Project is a reference to AtlasProject resource the deployment belongs to.
	Project common.ResourceRefNamespaced `json:"projectRef"`
	// Human-readable label that identifies the Federated Database Instance.
	Name string `json:"name"`
	// Configuration for the cloud provider where this Federated Database Instance is hosted.
	// +optional
	CloudProviderConfig *CloudProviderConfig `json:"cloudProviderConfig,omitempty"`
	// Information about the cloud provider region to which the Federated Database Instance routes client connections.
	// +optional
	DataProcessRegion *DataProcessRegion `json:"dataProcessRegion,omitempty"`
	// Configuration information for each data store and its mapping to MongoDB Atlas databases.
	// +optional
	Storage *Storage `json:"storage,omitempty"`
	// Private endpoint for Federated Database Instances and Online Archives to add to the specified project.
	// +optional
	PrivateEndpoints []DataFederationPE `json:"privateEndpoints,omitempty"`
}

type CloudProviderConfig struct {
	// Configuration for running Data Federation in AWS.
	AWS *AWSProviderConfig `json:"aws,omitempty"`
}

type AWSProviderConfig struct {
	// Unique identifier of the role that the data lake can use to access the data stores.Required if specifying cloudProviderConfig.
	RoleID string `json:"roleId,omitempty"`
	// Name of the S3 data bucket that the provided role ID is authorized to access.Required if specifying cloudProviderConfig.
	TestS3Bucket string `json:"testS3Bucket,omitempty"`
}

type DataProcessRegion struct {
	// Name of the cloud service that hosts the Federated Database Instance's infrastructure.
	// +kubebuilder:validation:Enum:=AWS
	CloudProvider string `json:"cloudProvider,omitempty"`
	// Name of the region to which the data lake routes client connections.
	// +kubebuilder:validation:Enum:=SYDNEY_AUS;MUMBAI_IND;FRANKFURT_DEU;DUBLIN_IRL;LONDON_GBR;VIRGINIA_USA;OREGON_USA;SAOPAULO_BRA;SINGAPORE_SGP
	Region string `json:"region,omitempty"`
}

type Storage struct {
	// Array that contains the queryable databases and collections for this data lake.
	Databases []Database `json:"databases,omitempty"`
	// Array that contains the data stores for the data lake.
	Stores []Store `json:"stores,omitempty"`
}

// Database associated with this data lake. Databases contain collections and views.
type Database struct {
	// Array of collections and data sources that map to a stores data store.
	Collections []Collection `json:"collections,omitempty"`
	// Maximum number of wildcard collections in the database. This only applies to S3 data sources.
	// Minimum value is 1, maximum value is 1000. Default value is 100.
	MaxWildcardCollections int `json:"maxWildcardCollections,omitempty"`
	// Human-readable label that identifies the database to which the data lake maps data.
	Name string `json:"name,omitempty"`
	// Array of aggregation pipelines that apply to the collection. This only applies to S3 data sources.
	Views []View `json:"views,omitempty"`
}

// Collection maps to a stores data store.
type Collection struct {
	// Array that contains the data stores that map to a collection for this data lake.
	DataSources []DataSource `json:"dataSources,omitempty"`
	// Human-readable label that identifies the collection to which MongoDB Atlas maps the data in the data stores.
	Name string `json:"name,omitempty"`
}

type View struct {
	// Human-readable label that identifies the view, which corresponds to an aggregation pipeline on a collection.
	Name string `json:"name,omitempty"`
	// Aggregation pipeline stages to apply to the source collection.
	Pipeline string `json:"pipeline,omitempty"`
	// Human-readable label that identifies the source collection for the view.
	Source string `json:"source,omitempty"`
}

type DataSource struct {
	// Flag that validates the scheme in the specified URLs.
	// If true, allows insecure HTTP scheme, doesn't verify the server's certificate chain and hostname, and accepts any certificate with any hostname presented by the server.
	// If false, allows secure HTTPS scheme only.
	AllowInsecure bool `json:"allowInsecure,omitempty"`
	// Human-readable label that identifies the collection in the database. For creating a wildcard (*) collection, you must omit this parameter.
	Collection string `json:"collection,omitempty"`
	// Regex pattern to use for creating the wildcard (*) collection.
	CollectionRegex string `json:"collectionRegex,omitempty"`
	// Human-readable label that identifies the database, which contains the collection in the cluster. You must omit this parameter to generate wildcard (*) collections for dynamically generated databases.
	Database string `json:"database,omitempty"`
	// Regex pattern to use for creating the wildcard (*) database.
	DatabaseRegex string `json:"databaseRegex,omitempty"`
	// File format that MongoDB Cloud uses if it encounters a file without a file extension while searching storeName.
	// +kubebuilder:validation:Enum:=.avro;.avro.bz2;.avro.gz;.bson;.bson.bz2;.bson.gz;.bsonx;.csv;.csv.bz2;.csv.gz;.json;.json.bz2;.json.gz;.orc;.parquet;.tsv;.tsv.bz2;.tsv.gz
	DefaultFormat string `json:"defaultFormat,omitempty"`
	// File path that controls how MongoDB Cloud searches for and parses files in the storeName before mapping them to a collection.
	// Specify / to capture all files and folders from the prefix path.
	Path string `json:"path,omitempty"`
	// Name for the field that includes the provenance of the documents in the results. MongoDB Atlas returns different fields in the results for each supported provider.
	ProvenanceFieldName string `json:"provenanceFieldName,omitempty"`
	// Human-readable label that identifies the data store that MongoDB Cloud maps to the collection.
	StoreName string `json:"storeName,omitempty"`
	// URLs of the publicly accessible data files. You can't specify URLs that require authentication.
	// Atlas Data Lake creates a partition for each URL. If empty or omitted, Data Lake uses the URLs from the store specified in the storeName parameter.
	Urls []string `json:"urls,omitempty"`
}

// Store is a group of settings that define where the data is stored.
type Store struct {
	// Human-readable label that identifies the data store. The storeName field references this values as part of the mapping configuration.
	// To use MongoDB Atlas as a data store, the data lake requires a serverless instance or an M10 or higher cluster.
	Name string `json:"name,omitempty"`
	// The provider used for data stores.
	Provider string `json:"provider,omitempty"`
	// Collection of AWS S3 storage classes. Atlas Data Lake includes the files in these storage classes in the query results.
	AdditionalStorageClasses []string `json:"additionalStorageClasses,omitempty"`
	// Human-readable label that identifies the AWS S3 bucket.
	// This label must exactly match the name of an S3 bucket that the data lake can access with the configured AWS Identity and Access Management (IAM) credentials.
	Bucket string `json:"bucket,omitempty"`
	// The delimiter that separates path segments in the data store.
	// MongoDB Atlas uses the delimiter to efficiently traverse S3 buckets with a hierarchical directory structure. You can specify any character supported by the S3 object keys as the delimiter.
	Delimiter string `json:"delimiter,omitempty"`
	// Flag that indicates whether to use S3 tags on the files in the given path as additional partition attributes.
	// If set to true, data lake adds the S3 tags as additional partition attributes and adds new top-level BSON elements associating each tag to each document.
	IncludeTags bool `json:"includeTags,omitempty"`
	// Prefix that MongoDB Cloud applies when searching for files in the S3 bucket.
	// The data store prepends the value of prefix to the path to create the full path for files to ingest.
	// If omitted, MongoDB Cloud searches all files from the root of the S3 bucket.
	Prefix string `json:"prefix,omitempty"`
	// Flag that indicates whether the bucket is public.
	// If set to true, MongoDB Cloud doesn't use the configured AWS Identity and Access Management (IAM) role to access the S3 bucket.
	// If set to false, the configured AWS IAM role must include permissions to access the S3 bucket.
	Public bool `json:"public,omitempty"`
	// Physical location where MongoDB Cloud deploys your AWS-hosted MongoDB cluster nodes. The region you choose can affect network latency for clients accessing your databases.
	// When MongoDB Atlas deploys a dedicated cluster, it checks if a VPC or VPC connection exists for that provider and region. If not, MongoDB Atlas creates them as part of the deployment.
	// To limit a new VPC peering connection to one CIDR block and region, create the connection first. Deploy the cluster after the connection starts.
	Region string `json:"region,omitempty"`
}

type DataFederationPE struct {
	// Unique 22-character alphanumeric string that identifies the private endpoint.
	EndpointID string `json:"endpointId,omitempty"`
	// Human-readable label that identifies the cloud service provider. Atlas Data Lake supports Amazon Web Services only.
	Provider string `json:"provider,omitempty"`
	// Human-readable label that identifies the resource type associated with this private endpoint.
	Type string `json:"type,omitempty"`
}

var _ api.AtlasCustomResource = &AtlasDataFederation{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=adf

// AtlasDataFederation is the Schema for the Atlas Data Federation API.
type AtlasDataFederation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataFederationSpec          `json:"spec,omitempty"`
	Status status.DataFederationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasDataFederationList contains a list of AtlasDataFederationList.
type AtlasDataFederationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasDataFederation `json:"items"`
}

func (c AtlasDataFederation) AtlasProjectObjectKey() client.ObjectKey {
	ns := c.Namespace
	if c.Spec.Project.Namespace != "" {
		ns = c.Spec.Project.Namespace
	}
	return kube.ObjectKey(ns, c.Spec.Project.Name)
}

func (c *AtlasDataFederation) GetStatus() api.Status {
	return c.Status
}

func (c *AtlasDataFederation) UpdateStatus(conditions []api.Condition, options ...api.Option) {
	c.Status.Conditions = conditions
	c.Status.ObservedGeneration = c.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.DataFederationStatusOption)
		v(&c.Status)
	}
}

func NewDataFederationInstance(projectName, instanceName, namespace string) *AtlasDataFederation {
	return &AtlasDataFederation{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: namespace,
		},
		Spec: DataFederationSpec{
			Project: common.ResourceRefNamespaced{
				Name:      projectName,
				Namespace: namespace,
			},
			Name:                instanceName,
			CloudProviderConfig: nil,
			DataProcessRegion:   nil,
			Storage:             nil,
			PrivateEndpoints:    nil,
		},
	}
}

func (c *AtlasDataFederation) WithAWSCloudProviderConfig(AWSRoleID, S3Bucket string) *AtlasDataFederation {
	c.Spec.CloudProviderConfig = &CloudProviderConfig{
		AWS: &AWSProviderConfig{
			RoleID:       AWSRoleID,
			TestS3Bucket: S3Bucket,
		}}
	return c
}

func (c *AtlasDataFederation) WithDataProcessingRegion(AWSRegion string) *AtlasDataFederation {
	c.Spec.DataProcessRegion = &DataProcessRegion{
		CloudProvider: "AWS",
		Region:        AWSRegion,
	}
	return c
}

func (c *AtlasDataFederation) WithStorage(storage *Storage) *AtlasDataFederation {
	c.Spec.Storage = storage
	return c
}

func (c *AtlasDataFederation) WithPrivateEndpoint(endpointID, provider, endpointType string) *AtlasDataFederation {
	c.Spec.PrivateEndpoints = append(c.Spec.PrivateEndpoints, DataFederationPE{
		EndpointID: endpointID,
		Provider:   provider,
		Type:       endpointType,
	})
	return c
}

func (c *AtlasDataFederation) WithAnnotations(annotations map[string]string) *AtlasDataFederation {
	c.Annotations = annotations
	return c
}
