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

package maintenance

import (
	"errors"
	"reflect"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
)

// MatchMaintenanceWindow returns the GomegaMatcher that checks if the 'actual' mongodbatlas.MaintenanceWindow matches
// the 'expected' akov2.MaintenanceWindow one.
// Note, that we cannot compare them by all the fields as Atlas tends to set default fields after MaintenanceWindow
// requests execution so we need to compare only the fields that remain in the same state
func MatchMaintenanceWindow(expected project.MaintenanceWindow) types.GomegaMatcher {
	return &maintenanceWindowMatcher{ExpectedMaintenanceWindow: expected}
}

type maintenanceWindowMatcher struct {
	ExpectedMaintenanceWindow project.MaintenanceWindow
}

func (m *maintenanceWindowMatcher) Match(actual any) (success bool, err error) {
	var c *admin.GroupMaintenanceWindow
	var ok bool
	if c, ok = actual.(*admin.GroupMaintenanceWindow); !ok {
		actualType := reflect.TypeOf(actual)
		return false, errors.New("Expected *mongodbatlas.MaintenanceWindow but received type " + actualType.String())
	}
	if c.GetDayOfWeek() != m.ExpectedMaintenanceWindow.DayOfWeek {
		return false, nil
	}
	if c.GetHourOfDay() != m.ExpectedMaintenanceWindow.HourOfDay {
		return false, nil
	}
	if c.GetAutoDeferOnceEnabled() != m.ExpectedMaintenanceWindow.AutoDefer {
		return false, nil
	}
	return true, nil
}

func (m *maintenanceWindowMatcher) FailureMessage(actual any) (message string) {
	return format.Message(actual, "to match", m.ExpectedMaintenanceWindow)
}

func (m *maintenanceWindowMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "not to match", m.ExpectedMaintenanceWindow)
}
