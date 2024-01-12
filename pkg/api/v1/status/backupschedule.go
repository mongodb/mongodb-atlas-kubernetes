package status

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/collection"

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

type BackupScheduleStatus struct {
	Common `json:",inline"`

	DeploymentIDs []string `json:"deploymentID,omitempty"`
}
