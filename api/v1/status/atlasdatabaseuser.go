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

// AtlasDatabaseUserStatusOption is the option that is applied to Atlas Project Status
type AtlasDatabaseUserStatusOption func(s *AtlasDatabaseUserStatus)

func AtlasDatabaseUserPasswordVersion(passwordVersion string) AtlasDatabaseUserStatusOption {
	return func(s *AtlasDatabaseUserStatus) {
		s.PasswordVersion = passwordVersion
	}
}

func AtlasDatabaseUserNameOption(name string) AtlasDatabaseUserStatusOption {
	return func(s *AtlasDatabaseUserStatus) {
		s.UserName = name
	}
}

// AtlasDatabaseUserStatus defines the observed state of AtlasProject
type AtlasDatabaseUserStatus struct {
	api.Common `json:",inline"`

	// PasswordVersion is the 'ResourceVersion' of the password Secret that the Atlas Operator is aware of
	PasswordVersion string `json:"passwordVersion,omitempty"`

	// UserName is the current name of database user.
	UserName string `json:"name,omitempty"`
}
