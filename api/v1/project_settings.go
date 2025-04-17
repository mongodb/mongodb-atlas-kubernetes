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

package v1

import (
	"go.mongodb.org/atlas-sdk/v20250312002/admin"
)

type ProjectSettings struct {
	IsCollectDatabaseSpecificsStatisticsEnabled *bool `json:"isCollectDatabaseSpecificsStatisticsEnabled,omitempty"`
	IsDataExplorerEnabled                       *bool `json:"isDataExplorerEnabled,omitempty"`
	IsExtendedStorageSizesEnabled               *bool `json:"isExtendedStorageSizesEnabled,omitempty"`
	IsPerformanceAdvisorEnabled                 *bool `json:"isPerformanceAdvisorEnabled,omitempty"`
	IsRealtimePerformancePanelEnabled           *bool `json:"isRealtimePerformancePanelEnabled,omitempty"`
	IsSchemaAdvisorEnabled                      *bool `json:"isSchemaAdvisorEnabled,omitempty"`
}

func (s ProjectSettings) ToAtlas() *admin.GroupSettings {
	atlas := &admin.GroupSettings{}

	atlas.IsCollectDatabaseSpecificsStatisticsEnabled = s.IsCollectDatabaseSpecificsStatisticsEnabled
	atlas.IsDataExplorerEnabled = s.IsDataExplorerEnabled
	atlas.IsExtendedStorageSizesEnabled = s.IsExtendedStorageSizesEnabled
	atlas.IsPerformanceAdvisorEnabled = s.IsPerformanceAdvisorEnabled
	atlas.IsRealtimePerformancePanelEnabled = s.IsRealtimePerformancePanelEnabled
	atlas.IsSchemaAdvisorEnabled = s.IsSchemaAdvisorEnabled

	return atlas
}

func ProjectSettingsFromAtlas(atlas *admin.GroupSettings) *ProjectSettings {
	ps := &ProjectSettings{}

	ps.IsCollectDatabaseSpecificsStatisticsEnabled = atlas.IsCollectDatabaseSpecificsStatisticsEnabled
	ps.IsDataExplorerEnabled = atlas.IsDataExplorerEnabled
	ps.IsExtendedStorageSizesEnabled = atlas.IsExtendedStorageSizesEnabled
	ps.IsPerformanceAdvisorEnabled = atlas.IsPerformanceAdvisorEnabled
	ps.IsRealtimePerformancePanelEnabled = atlas.IsRealtimePerformancePanelEnabled
	ps.IsSchemaAdvisorEnabled = atlas.IsSchemaAdvisorEnabled

	return ps
}
