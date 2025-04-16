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

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

type DataFederationStatus struct {
	api.Common `json:",inline"`

	// MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`
}

// +k8s:deepcopy-gen=false

type DataFederationStatusOption func(s *DataFederationStatus)
