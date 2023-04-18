package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type ProjectSettings struct {
	IsCollectDatabaseSpecificsStatisticsEnabled *bool `json:"isCollectDatabaseSpecificsStatisticsEnabled,omitempty"`
	IsDataExplorerEnabled                       *bool `json:"isDataExplorerEnabled,omitempty"`
	IsExtendedStorageSizesEnabled               *bool `json:"isExtendedStorageSizesEnabled,omitempty"`
	IsPerformanceAdvisorEnabled                 *bool `json:"isPerformanceAdvisorEnabled,omitempty"`
	IsRealtimePerformancePanelEnabled           *bool `json:"isRealtimePerformancePanelEnabled,omitempty"`
	IsSchemaAdvisorEnabled                      *bool `json:"isSchemaAdvisorEnabled,omitempty"`
}

func (s ProjectSettings) ToAtlas() (*mongodbatlas.ProjectSettings, error) {
	result := &mongodbatlas.ProjectSettings{}
	err := compat.JSONCopy(result, s)
	return result, err
}
