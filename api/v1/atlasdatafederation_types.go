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
	// Project is a reference to AtlasProject resource the deployment belongs to
	Project common.ResourceRefNamespaced `json:"projectRef"`
	Name    string                       `json:"name"`

	// +optional
	CloudProviderConfig *CloudProviderConfig `json:"cloudProviderConfig,omitempty"`

	// +optional
	DataProcessRegion *DataProcessRegion `json:"dataProcessRegion,omitempty"`

	// +optional
	Storage *Storage `json:"storage,omitempty"`

	// +optional
	PrivateEndpoints []DataFederationPE `json:"privateEndpoints,omitempty"`
}

type CloudProviderConfig struct {
	AWS *AWSProviderConfig `json:"aws,omitempty"`
}

type AWSProviderConfig struct {
	RoleID       string `json:"roleId,omitempty"`
	TestS3Bucket string `json:"testS3Bucket,omitempty"`
}

type DataProcessRegion struct {
	// +kubebuilder:validation:Enum:=AWS
	CloudProvider string `json:"cloudProvider,omitempty"`
	// +kubebuilder:validation:Enum:=SYDNEY_AUS;MUMBAI_IND;FRANKFURT_DEU;DUBLIN_IRL;LONDON_GBR;VIRGINIA_USA;OREGON_USA;SAOPAULO_BRA;SINGAPORE_SGP
	Region string `json:"region,omitempty"`
}

type Storage struct {
	Databases []Database `json:"databases,omitempty"`
	Stores    []Store    `json:"stores,omitempty"`
}

type Database struct {
	Collections            []Collection `json:"collections,omitempty"`
	MaxWildcardCollections int          `json:"maxWildcardCollections,omitempty"`
	Name                   string       `json:"name,omitempty"`
	Views                  []View       `json:"views,omitempty"`
}

type Collection struct {
	DataSources []DataSource `json:"dataSources,omitempty"`
	Name        string       `json:"name,omitempty"`
}

type View struct {
	Name     string `json:"name,omitempty"`
	Pipeline string `json:"pipeline,omitempty"`
	Source   string `json:"source,omitempty"`
}

type DataSource struct {
	AllowInsecure   bool   `json:"allowInsecure,omitempty"`
	Collection      string `json:"collection,omitempty"`
	CollectionRegex string `json:"collectionRegex,omitempty"`
	Database        string `json:"database,omitempty"`
	DatabaseRegex   string `json:"databaseRegex,omitempty"`
	// +kubebuilder:validation:Enum:=.avro;.avro.bz2;.avro.gz;.bson;.bson.bz2;.bson.gz;.bsonx;.csv;.csv.bz2;.csv.gz;.json;.json.bz2;.json.gz;.orc;.parquet;.tsv;.tsv.bz2;.tsv.gz
	DefaultFormat       string   `json:"defaultFormat,omitempty"`
	Path                string   `json:"path,omitempty"`
	ProvenanceFieldName string   `json:"provenanceFieldName,omitempty"`
	StoreName           string   `json:"storeName,omitempty"`
	Urls                []string `json:"urls,omitempty"`
}

type Store struct {
	Name                     string   `json:"name,omitempty"`
	Provider                 string   `json:"provider,omitempty"`
	AdditionalStorageClasses []string `json:"additionalStorageClasses,omitempty"`
	Bucket                   string   `json:"bucket,omitempty"`
	Delimiter                string   `json:"delimiter,omitempty"`
	IncludeTags              bool     `json:"includeTags,omitempty"`
	Prefix                   string   `json:"prefix,omitempty"`
	Public                   bool     `json:"public,omitempty"`
	Region                   string   `json:"region,omitempty"`
}

type DataFederationPE struct {
	EndpointID string `json:"endpointId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Type       string `json:"type,omitempty"`
}

var _ api.AtlasCustomResource = &AtlasDataFederation{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com
// +kubebuilder:resource:categories=atlas,shortName=adf

// AtlasDataFederation is the Schema for the Atlas Data Federation API
type AtlasDataFederation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DataFederationSpec          `json:"spec,omitempty"`
	Status status.DataFederationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AtlasDataFederationList contains a list of AtlasDataFederationList
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
