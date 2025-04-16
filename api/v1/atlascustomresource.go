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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

//+k8s:deepcopy-gen=false

// AtlasCustomResource is the interface common for all Atlas entities
type AtlasCustomResource interface {
	metav1.Object
	runtime.Object
	api.Reader
	api.Writer
}

var _ AtlasCustomResource = &AtlasProject{}
var _ AtlasCustomResource = &AtlasTeam{}
var _ AtlasCustomResource = &AtlasDeployment{}
var _ AtlasCustomResource = &AtlasDatabaseUser{}
var _ AtlasCustomResource = &AtlasDataFederation{}
var _ AtlasCustomResource = &AtlasBackupSchedule{}
var _ AtlasCustomResource = &AtlasBackupPolicy{}
var _ AtlasCustomResource = &AtlasFederatedAuth{}
var _ AtlasCustomResource = &AtlasStreamInstance{}
var _ AtlasCustomResource = &AtlasStreamConnection{}
var _ AtlasCustomResource = &AtlasSearchIndexConfig{}
var _ AtlasCustomResource = &AtlasBackupCompliancePolicy{}
var _ AtlasCustomResource = &AtlasPrivateEndpoint{}
var _ AtlasCustomResource = &AtlasCustomRole{}

// InitCondition initializes the underlying type of the given condition to the given default value.
func InitCondition(resource AtlasCustomResource, defaultCondition api.Condition) []api.Condition {
	return api.EnsureConditionExists(defaultCondition, resource.GetStatus().GetConditions())
}
