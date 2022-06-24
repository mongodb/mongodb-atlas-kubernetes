package testutil

import (
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

// MatchMaintenanceWindow returns the GomegaMatcher that checks if the 'actual' mongodbatlas.MaintenanceWindow matches
// the 'expected' mdbv1.MaintenanceWindow  one.
// Note, that we cannot compare them by all the fields as Atlas tends to set default fields after MaintenanceWindow
// requests execution so we need to compare only the fields that remain in the same state
func MatchMaintenanceWindow(expected project.MaintenanceWindow) types.GomegaMatcher {
	return &maintenanceWindowMatcher{ExpectedMaintenanceWindow: expected}
}

type maintenanceWindowMatcher struct {
	ExpectedMaintenanceWindow project.MaintenanceWindow
}

func (m *maintenanceWindowMatcher) Match(actual interface{}) (success bool, err error) {
	var c mongodbatlas.MaintenanceWindow
	var ok bool
	if c, ok = actual.(mongodbatlas.MaintenanceWindow); !ok {
		panic("Expected mongodbatlas.ProjectIPAccessList")
	}
	if m.ExpectedMaintenanceWindow.DayOfWeek != 0 && c.DayOfWeek != m.ExpectedMaintenanceWindow.DayOfWeek {
		return false, nil
	}
	if *c.HourOfDay != m.ExpectedMaintenanceWindow.HourOfDay {
		return false, nil
	}
	// TODO : check assumption : an autoDefer POST request enable the field AutoDeferOnceEnabled of the maintenance object
	if *c.AutoDeferOnceEnabled != m.ExpectedMaintenanceWindow.AutoDefer {
		return false, nil
	}
	return true, nil
}

func (m *maintenanceWindowMatcher) FailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "to match", m.ExpectedMaintenanceWindow)
}

func (m *maintenanceWindowMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return format.Message(actual, "not to match", m.ExpectedMaintenanceWindow)
}
