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

// AtlasBackupScheduleStatusOption is the option that is applied to AtlasBackupSchedule Status
type AtlasBackupScheduleStatusOption func(s *BackupScheduleStatus)

func AtlasBackupScheduleSetDeploymentID(ID string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		IDs := collection.CopyWithSkip(s.DeploymentIDs, ID)
		IDs = append(IDs, ID)

		s.DeploymentIDs = IDs
	}
}

func AtlasBackupScheduleUnsetDeploymentID(ID string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		s.DeploymentIDs = collection.CopyWithSkip(s.DeploymentIDs, ID)
	}
}

// BackupScheduleStatus defines the observed state of AtlasBackupSchedule.
type BackupScheduleStatus struct {
	api.Common `json:",inline"`

	// List of the human-readable names of all deployments utilizing this backup schedule.
	DeploymentIDs []string `json:"deploymentID,omitempty"`
}
