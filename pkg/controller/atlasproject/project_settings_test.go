package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

func TestAreSettingsInSync(t *testing.T) {
	atlasDef := &akov2.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
		IsDataExplorerEnabled:                       pointer.MakePtr(true),
		IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
		IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
		IsSchemaAdvisorEnabled:                      pointer.MakePtr(true),
	}
	specDef := &akov2.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
		IsDataExplorerEnabled:                       pointer.MakePtr(true),
	}

	areEqual := areSettingsInSync(atlasDef, specDef)
	assert.True(t, areEqual, "Only fields which are set should be compared")

	specDef.IsPerformanceAdvisorEnabled = pointer.MakePtr(false)
	areEqual = areSettingsInSync(atlasDef, specDef)
	assert.False(t, areEqual, "Field values should be the same ")
}
