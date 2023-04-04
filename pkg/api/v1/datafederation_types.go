package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	RoleId       string `json:"roleId,omitempty"`
	TestS3Bucket string `json:"testS3Bucket,omitempty"`
}

type DataProcessRegion struct {
	CloudProvider string `json:"cloudProvider,omitempty"` // "AWS" always
	Region        string `json:"region,omitempty"`
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
	Allowinsecure       bool     `json:"allowInsecure,omitempty"`
	Collection          string   `json:"collection,omitempty"`
	Collectionregex     string   `json:"collectionRegex,omitempty"`
	Database            string   `json:"database,omitempty"`
	Databaseregex       string   `json:"databaseRegex,omitempty"`
	Defaultformat       string   `json:"defaultFormat,omitempty"`
	Path                string   `json:"path,omitempty"`
	Provenancefieldname string   `json:"provenanceFieldName,omitempty"`
	Storename           string   `json:"storeName,omitempty"`
	Urls                []string `json:"urls,omitempty"`
}

type Store struct {
	Name                     string   `json:"name,omitempty"`
	Provider                 string   `json:"provider,omitempty"`
	Additionalstorageclasses []string `json:"additionalStorageClasses,omitempty"`
	Bucket                   string   `json:"bucket,omitempty"`
	Delimiter                string   `json:"delimiter,omitempty"`
	Includetags              bool     `json:"includeTags,omitempty"`
	Prefix                   string   `json:"prefix,omitempty"`
	Public                   bool     `json:"public,omitempty"`
	Region                   string   `json:"region,omitempty"`
}

type DataFederationPE struct {
	EndpointID string `json:"endpointId,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Type       string `json:"type,omitempty"`
}

func (pe *DataFederationPE) Identifier() interface{} {
	return pe.EndpointID
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:subresource:status
// +groupName:=atlas.mongodb.com

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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

func (c *AtlasDataFederation) GetStatus() status.Status {
	return c.Status
}

func (c *AtlasDataFederation) UpdateStatus(conditions []status.Condition, options ...status.Option) {
	c.Status.Conditions = conditions
	c.Status.ObservedGeneration = c.ObjectMeta.Generation

	for _, o := range options {
		// This will fail if the Option passed is incorrect - which is expected
		v := o.(status.DataFederationStatusOption)
		v(&c.Status)
	}
}
