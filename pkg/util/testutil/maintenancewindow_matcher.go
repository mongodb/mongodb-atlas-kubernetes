// TODO: move away from pkg, this code is only usable from tests
package testutil

import (
	"errors"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

// MatchMaintenanceWindow returns the GomegaMatcher that checks if the 'actual' mongodbatlas.MaintenanceWindow matches
// the 'expected' mdbv1.MaintenanceWindow one.
// Note, that we cannot compare them by all the fields as Atlas tends to set default fields after MaintenanceWindow
// requests execution so we need to compare only the fields that remain in the same state
func MatchMaintenanceWindow(expected project.MaintenanceWindow) types.GomegaMatcher {
	return &maintenanceWindowMatcher{ExpectedMaintenanceWindow: expected}
}

type maintenanceWindowMatcher struct {
	ExpectedMaintenanceWindow project.MaintenanceWindow
}

func (m *maintenanceWindowMatcher) Match(actual interface{}) (success bool, err error) {
	var c *mongodbatlas.MaintenanceWindow
	var ok bool
	if c, ok = actual.(*mongodbatlas.MaintenanceWindow); !ok {
		actualType := reflect.TypeOf(actual)
		return false, errors.New("Expected *mongodbatlas.MaintenanceWindow but received type " + actualType.String())
	}
	if c.DayOfWeek != m.ExpectedMaintenanceWindow.DayOfWeek {
		return false, nil
	}
	if *c.HourOfDay != m.ExpectedMaintenanceWindow.HourOfDay {
		return false, nil
	}
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
