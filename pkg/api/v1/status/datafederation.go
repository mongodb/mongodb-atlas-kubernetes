package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"

type DataFederationStatus struct {
	api.Common `json:",inline"`

	// MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`
}

// +k8s:deepcopy-gen=false

type DataFederationStatusOption func(s *DataFederationStatus)
