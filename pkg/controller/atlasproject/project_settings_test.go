package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func TestAreSettingsInSync(t *testing.T) {
	atlas := &v1.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
		IsDataExplorerEnabled:                       toptr.MakePtr(true),
		IsPerformanceAdvisorEnabled:                 toptr.MakePtr(true),
		IsRealtimePerformancePanelEnabled:           toptr.MakePtr(true),
		IsSchemaAdvisorEnabled:                      toptr.MakePtr(true),
	}
	spec := &v1.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: toptr.MakePtr(true),
		IsDataExplorerEnabled:                       toptr.MakePtr(true),
	}

	areEqual := areSettingsInSync(atlas, spec)
	assert.True(t, areEqual, "Only fields which are set should be compared")

	spec.IsPerformanceAdvisorEnabled = toptr.MakePtr(false)
	areEqual = areSettingsInSync(atlas, spec)
	assert.False(t, areEqual, "Field values should be the same ")
}
