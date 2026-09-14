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

// AtlasBackupScheduleSetDeploymentID records a deployment as a user of this
// backup schedule.
//
// staleKeys are keys the same deployment may already be recorded under, written
// by an operator version that used a different key format. They are dropped
// here so an existing status converges on the current format the first time the
// deployment reconciles; without that, garbage collection would never match the
// old entries and the deletion finalizer would stay on forever.
func AtlasBackupScheduleSetDeploymentID(key string, staleKeys ...string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		keys := collection.CopyWithSkip(s.DeploymentIDs, key)
		for _, staleKey := range staleKeys {
			keys = collection.CopyWithSkip(keys, staleKey)
		}

		s.DeploymentIDs = append(keys, key)
	}
}

func AtlasBackupScheduleUnsetDeploymentID(keys ...string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		remaining := s.DeploymentIDs
		for _, key := range keys {
			remaining = collection.CopyWithSkip(remaining, key)
		}

		s.DeploymentIDs = remaining
	}
}

// BackupScheduleStatus defines the observed state of AtlasBackupSchedule.
type BackupScheduleStatus struct {
	api.Common `json:",inline"`

	// List of keys identifying the deployments that use this backup schedule, in
	// "namespace/name" form. A schedule can be referenced from another namespace,
	// so the namespace is part of the key.
	// The json tag is kept for compatibility with statuses written earlier.
	DeploymentIDs []string `json:"deploymentID,omitempty"`
}
