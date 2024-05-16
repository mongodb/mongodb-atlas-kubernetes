package status

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
)

type AtlasAuditingStatusOption func(*AtlasAuditingStatus)

// AtlasAuditingStatus defines the observed state of Auditing for one of more projects
// +k8s:deepcopy-gen=true
type AtlasAuditingStatus struct {
	api.Common `json:",inline"`

	// Projects lists each associated project auditing status
	Projects []ProjectAuditingStatus `json:"projects"`
}

// ProjectAuditingStatus defined an individual Project Auditing status condition
// +k8s:deepcopy-gen=true
type ProjectAuditingStatus struct {
	api.Condition `json:",inline"`

	// ID represents the project ID
	ID string `json:"projects"`
}

func WithSuccess(projectID string) AtlasAuditingStatusOption {
	return func(s *AtlasAuditingStatus) {
		idx := projectIndex(s.Projects, projectID)
		successCondition := api.Condition{
			Type:               "Ready",
			Status:             "True",
			LastTransitionTime: metav1.Now(),
			Reason:             "",
		}
		if idx < 0 {
			s.Projects = append(s.Projects, ProjectAuditingStatus{
				ID:        projectID,
				Condition: successCondition,
			})
		} else {
			s.Projects[idx].Condition = successCondition
		}
	}
}

func WithProjectFailure(projectID string, err error) AtlasAuditingStatusOption {
	return func(s *AtlasAuditingStatus) {
		idx := projectIndex(s.Projects, projectID)
		failureCondition := api.Condition{
			Type:               "Ready",
			Status:             "False",
			LastTransitionTime: metav1.Now(),
			Reason:             err.Error(),
		}
		if idx < 0 {
			s.Projects = append(s.Projects, ProjectAuditingStatus{
				ID:        projectID,
				Condition: failureCondition,
			})
		} else {
			s.Projects[idx].Condition = failureCondition
		}
	}
}

func projectIndex(projects []ProjectAuditingStatus, projectID string) int {
	for i, project := range projects {
		if project.ID == projectID {
			return i
		}
	}
	return -1
}
