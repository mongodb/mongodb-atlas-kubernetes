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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"
)

// +k8s:deepcopy-gen=false

// AtlasBackupPolicyStatusOption is the option that is applied to AtlasBackupPolicy Status
type AtlasBackupPolicyStatusOption func(s *BackupPolicyStatus)

func AtlasBackupPolicySetScheduleID(ID string) AtlasBackupPolicyStatusOption {
	return func(s *BackupPolicyStatus) {
		IDs := collection.CopyWithSkip(s.BackupScheduleIDs, ID)
		IDs = append(IDs, ID)

		s.BackupScheduleIDs = IDs
	}
}

func AtlasBackupPolicyUnsetScheduleID(ID string) AtlasBackupPolicyStatusOption {
	return func(s *BackupPolicyStatus) {
		s.BackupScheduleIDs = collection.CopyWithSkip(s.BackupScheduleIDs, ID)
	}
}

// BackupPolicyStatus defines the observed state of AtlasBackupPolicy.
type BackupPolicyStatus struct {
	api.Common `json:",inline"`

	// DeploymentID of the deployment using the backup policy
	BackupScheduleIDs []string `json:"backupScheduleIDs,omitempty"`
}
