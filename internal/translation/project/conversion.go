package project

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/tag"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type Project struct {
	OrgID                     string
	ID                        string
	Name                      string
	RegionUsageRestrictions   string
	WithDefaultAlertsSettings bool
	Tags                      []tag.Tag
}

func NewProject(project *akov2.AtlasProject, orgID string) *Project {
	return &Project{
		OrgID:                     orgID,
		ID:                        project.ID(),
		Name:                      project.Spec.Name,
		RegionUsageRestrictions:   project.Spec.RegionUsageRestrictions,
		WithDefaultAlertsSettings: project.Spec.WithDefaultAlertsSettings,
	}
}

func fromAtlas(group *admin.Group) *Project {
	return &Project{
		OrgID:                     group.GetOrgId(),
		ID:                        group.GetId(),
		Name:                      group.GetName(),
		RegionUsageRestrictions:   group.GetRegionUsageRestrictions(),
		WithDefaultAlertsSettings: group.GetWithDefaultAlertsSettings(),
		Tags:                      tag.FromAtlas(group.GetTags()),
	}
}

func toAtlas(project *Project) *admin.Group {
	return &admin.Group{
		OrgId:                     project.OrgID,
		Name:                      project.Name,
		RegionUsageRestrictions:   &project.RegionUsageRestrictions,
		Tags:                      tag.ToAtlas(project.Tags),
		WithDefaultAlertsSettings: &project.WithDefaultAlertsSettings,
	}
}
