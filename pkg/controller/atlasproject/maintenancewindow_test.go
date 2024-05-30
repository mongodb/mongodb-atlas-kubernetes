package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

func TestValidateMaintenanceWindow(t *testing.T) {
	testCases := []struct {
		in    project.MaintenanceWindow
		valid bool
	}{
		// All fields empty, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: false,
				Defer:     false,
				AutoDefer: false,
			},
			valid: false,
		},

		// Only dayOfWeek specified, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 1, // Sunday
				HourOfDay: 0, // Will default to midnight
				StartASAP: false,
				Defer:     false,
				AutoDefer: false,
			},
			valid: true,
		},

		// Specify window and autoDefer, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: false,
				Defer:     false,
				AutoDefer: true,
			},
			valid: true,
		},

		// Specify window, autoDefer, and startASAP, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: true,
				Defer:     false,
				AutoDefer: true,
			},
			valid: true,
		},

		// Specify window, autoDefer, and defer, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: false,
				Defer:     true,
				AutoDefer: true,
			},
			valid: true,
		},

		// StartASAP only, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: true,
				Defer:     false,
				AutoDefer: false,
			},
			valid: false,
		},

		// AutoDefer only, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 0,
				HourOfDay: 0,
				StartASAP: false,
				Defer:     false,
				AutoDefer: true,
			},
			valid: false,
		},

		// Both startASAP and Defer specified, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek: 3,  // Tuesday
				HourOfDay: 14, // 2pm
				StartASAP: true,
				Defer:     true,
				AutoDefer: true,
			},
			valid: false,
		},
	}

	for _, testCase := range testCases {
		t.Run("", func(t *testing.T) {
			err := validateMaintenanceWindow(testCase.in)
			if !testCase.valid {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
