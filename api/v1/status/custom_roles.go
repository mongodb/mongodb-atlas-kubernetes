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

// +k8s:deepcopy-gen=false

type CustomRoleStatus string

const (
	CustomRoleStatusOK     CustomRoleStatus = "OK"
	CustomRoleStatusFailed CustomRoleStatus = "FAILED"
)

type CustomRole struct {
	// Role name which is unique
	Name string `json:"name"`
	// The status of the given custom role (OK or FAILED)
	Status CustomRoleStatus `json:"status"`
	// The message when the custom role is in the FAILED status
	Error string `json:"error,omitempty"`
}

// AtlasCustomRoleStatus is a status for the AtlasCustomRole Custom resource.
// Not the one included in the AtlasProject
type AtlasCustomRoleStatus struct {
	api.Common `json:",inline"`
}
