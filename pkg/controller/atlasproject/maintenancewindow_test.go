package atlasproject

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/project"
)

func TestValidateMaintenanceWindow(t *testing.T) {
	testCases := []struct {
		in    project.MaintenanceWindow
		valid bool
	}{
		// All fields empty, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: false,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: true,
		},

		// Modify just the window (hour and day), valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            1, // Sunday
				HourOfDay:            2, // 2am
				AutoDeferOnceEnabled: false,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: true,
		},

		// Modify window and autoDefer, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            3,  // Tuesday
				HourOfDay:            14, // 2pm
				AutoDeferOnceEnabled: true,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: true,
		},

		// Modify window and startASAP, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            3,  // Tuesday
				HourOfDay:            14, // 2pm
				AutoDeferOnceEnabled: false,
				StartASAP:            true,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: true,
		},

		// startASAP only, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: false,
				StartASAP:            true,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: true,
		},

		// Defer only, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: false,
				StartASAP:            false,
				Defer:                true,
				AutoDefer:            false,
			},
			valid: true,
		},

		// Auto-defer only, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: false,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            true,
			},
			valid: true,
		},

		// Modify window, both startASAP and autoDeferOnceEnabled activated, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            1, // Sunday
				HourOfDay:            2, // 2am
				AutoDeferOnceEnabled: true,
				StartASAP:            true,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: false,
		},

		// Defer and other fields enabled, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            1,
				HourOfDay:            0,
				AutoDeferOnceEnabled: false,
				StartASAP:            true,
				Defer:                true,
				AutoDefer:            false,
			},
			valid: false,
		},

		// Auto-defer and another field enabled, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: true,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            true,
			},
			valid: false,
		},

		// AutoDeferOnceEnabled only, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            0,
				HourOfDay:            0,
				AutoDeferOnceEnabled: true,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            false,
			},
			valid: false,
		},

		// One out of two fields of the window specified, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:            2, // Monday
				HourOfDay:            0, // Empty
				AutoDeferOnceEnabled: true,
				StartASAP:            false,
				Defer:                false,
				AutoDefer:            false,
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
