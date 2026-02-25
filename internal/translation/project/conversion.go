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

package project

import (
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/tag"
)

type Project struct {
	OrgID                     string
	ID                        string
	Name                      string
	RegionUsageRestrictions   string
	WithDefaultAlertsSettings bool
	Tags                      []*akov2.TagSpec
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
