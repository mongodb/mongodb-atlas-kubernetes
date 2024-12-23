package v1

import (
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
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
