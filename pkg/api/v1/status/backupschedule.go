package status

// +k8s:deepcopy-gen=false

// AtlasBackupScheduleStatusOption is the option that is applied to AtlasBackupSchedule Status
type AtlasBackupScheduleStatusOption func(s *BackupScheduleStatus)

func AtlasBackupScheduleSetDeploymentID(ID string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		IDs := copyListWithSkip(s.DeploymentIDs, ID)
		IDs = append(IDs, ID)

		s.DeploymentIDs = IDs
	}
}

func AtlasBackupScheduleUnsetDeploymentID(ID string) AtlasBackupScheduleStatusOption {
	return func(s *BackupScheduleStatus) {
		s.DeploymentIDs = copyListWithSkip(s.DeploymentIDs, ID)
	}
}

type BackupScheduleStatus struct {
	Common `json:",inline"`

	DeploymentIDs []string `json:"deploymentID,omitempty"`
}

func copyListWithSkip[T comparable](list []T, skip T) []T {
	newList := make([]T, 0, len(list))

	for _, item := range list {
		if item != skip {
			newList = append(newList, item)
		}
	}

	return newList
}
