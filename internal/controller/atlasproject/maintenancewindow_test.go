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

package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/maintenancewindow"
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

		// Wave assignment only, valid
		{
			in: project.MaintenanceWindow{
				WaveAssignment: pointer.MakePtr(2),
			},
			valid: true,
		},

		// Wave assignment with a schedule, valid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:      3, // Tuesday
				HourOfDay:      14,
				WaveAssignment: pointer.MakePtr(1),
			},
			valid: true,
		},

		// Wave assignment together with autoDefer, invalid
		{
			in: project.MaintenanceWindow{
				DayOfWeek:      3, // Tuesday
				HourOfDay:      14,
				AutoDefer:      true,
				WaveAssignment: pointer.MakePtr(1),
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

func TestEnsureMaintenanceWindow(t *testing.T) {
	for _, tc := range []struct {
		name string

		maintenanceWindow  project.MaintenanceWindow
		maintenanceService maintenancewindow.MaintenanceWindowService

		isOK             bool
		wantStatus       string
		conditionMissing bool
	}{
		{
			name: "maintenance window is the same in atlas and kube",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "validation fails",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				StartASAP: true,
				Defer:     true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				return nil
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name: "get request errors",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{},
					errors.New("TEST GET ERROR"),
				)
				return service
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name: "maintenance window is different (update)",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 4,
							HourOfDay: 18,
						},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "start maintenance ASAP",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				StartASAP: true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "update request fails",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				StartASAP: true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 4,
							HourOfDay: 18,
						},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).
					Return(errors.New("TEST UPDATE ERROR"))
				return service
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name:              "maintenance window not in AKO (delete/reset)",
			maintenanceWindow: project.MaintenanceWindow{},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Reset(context.Background(), "testid123").Return(nil)
				return service
			}(),
			isOK:             true,
			wantStatus:       "False",
			conditionMissing: true,
		},
		{
			name:              "reset request errors",
			maintenanceWindow: project.MaintenanceWindow{},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Reset(context.Background(), "testid123").Return(errors.New("TEST RESET ERROR"))
				return service
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name: "auto defer toggled",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				AutoDefer: true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				service.EXPECT().ToggleAutoDefer(context.Background(), "testid123").Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},

		{
			name: "auto defer toggle errors",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				AutoDefer: true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				service.EXPECT().ToggleAutoDefer(context.Background(), "testid123").Return(errors.New("TEST AUTO DEFER ERROR"))
				return service
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name: "defer maintenance",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				Defer:     true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				service.EXPECT().Defer(context.Background(), "testid123").Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "defer maintenance and update",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				Defer:     true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 4,
							HourOfDay: 18,
						},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).Return(nil)
				service.EXPECT().Defer(context.Background(), "testid123").Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "defer maintenance errors",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
				Defer:     true,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek: 2,
							HourOfDay: 14,
						},
					},
					nil,
				)
				service.EXPECT().Defer(context.Background(), "testid123").Return(errors.New("TEST DEFER ERROR"))
				return service
			}(),
			isOK:       false,
			wantStatus: "False",
		},
		{
			name: "wave-only window is created without a schedule",
			maintenanceWindow: project.MaintenanceWindow{
				WaveAssignment: pointer.MakePtr(2),
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "wave assignment drift triggers update",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek:      2,
				HourOfDay:      14,
				WaveAssignment: pointer.MakePtr(3),
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek:      2,
							HourOfDay:      14,
							WaveAssignment: pointer.MakePtr(1),
						},
					},
					nil,
				)
				service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).Return(nil)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			// The spec manages only the wave, and it already matches. Atlas keeps
			// its own schedule, which AKO must leave alone instead of reporting
			// drift it cannot correct.
			name: "wave-only window leaves an existing atlas schedule alone",
			maintenanceWindow: project.MaintenanceWindow{
				WaveAssignment: pointer.MakePtr(2),
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek:      3,
							HourOfDay:      5,
							WaveAssignment: pointer.MakePtr(2),
						},
					},
					nil,
				)
				// No Update expectation: any call fails the test.
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
		{
			name: "unmanaged wave assignment does not trigger update",
			maintenanceWindow: project.MaintenanceWindow{
				DayOfWeek: 2,
				HourOfDay: 14,
			},
			maintenanceService: func() maintenancewindow.MaintenanceWindowService {
				service := translation.NewMaintenanceWindowServiceMock(t)
				service.EXPECT().Get(context.Background(), "testid123").Return(
					&maintenancewindow.MaintenanceWindow{
						MaintenanceWindow: &project.MaintenanceWindow{
							DayOfWeek:      2,
							HourOfDay:      14,
							WaveAssignment: pointer.MakePtr(1),
						},
					},
					nil,
				)
				return service
			}(),
			isOK:       true,
			wantStatus: "True",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			proj := &akov2.AtlasProject{
				Status: status.AtlasProjectStatus{
					ID: "testid123",
				},
				Spec: akov2.AtlasProjectSpec{
					MaintenanceWindow: tc.maintenanceWindow,
				},
			}

			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     zap.S(),
			}
			ctx.SetConditionTrue(api.MaintenanceWindowReadyType)

			r := AtlasProjectReconciler{}

			result := r.ensureMaintenanceWindow(ctx, proj, tc.maintenanceService)

			assert.Equal(t, tc.isOK, result.IsOk())

			con, ok := ctx.GetCondition(api.MaintenanceWindowReadyType)
			if tc.conditionMissing {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.Equal(t, tc.wantStatus, string(con.Status))
			}
		})
	}
}

// TestWaveOnlyWindowConvergesAcrossReconciles guards against the PATCH loop the
// wave-only window used to cause: toAtlas drops the zero dayOfWeek/hourOfDay, so
// Atlas never adopts them and a naive comparison reports drift forever. Asserting
// on a single reconcile would not catch it — the first Update looks legitimate.
func TestWaveOnlyWindowConvergesAcrossReconciles(t *testing.T) {
	const reconciles = 3

	updates := 0
	service := translation.NewMaintenanceWindowServiceMock(t)
	// Atlas keeps its schedule, because the request never carries one.
	service.EXPECT().Get(context.Background(), "testid123").Return(
		&maintenancewindow.MaintenanceWindow{
			MaintenanceWindow: &project.MaintenanceWindow{
				DayOfWeek:      3,
				HourOfDay:      5,
				WaveAssignment: pointer.MakePtr(2),
			},
		},
		nil,
	).Times(reconciles)
	service.EXPECT().Update(context.Background(), "testid123", mock.AnythingOfType("*maintenancewindow.MaintenanceWindow")).
		Run(func(context.Context, string, *maintenancewindow.MaintenanceWindow) { updates++ }).
		Return(nil).
		Maybe()

	r := &AtlasProjectReconciler{}
	windowSpec := project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)}

	for range reconciles {
		ctx := &workflow.Context{Context: context.Background(), Log: zap.S()}
		result := r.syncAtlasWithSpec(ctx, "testid123", windowSpec, service)
		assert.True(t, result.IsOk())
	}

	assert.Zero(t, updates, "a converged wave-only window must not keep patching Atlas")
}

func TestWaveIsDifferent(t *testing.T) {
	testCases := []struct {
		name        string
		inAtlas     project.MaintenanceWindow
		inAKO       project.MaintenanceWindow
		isDifferent bool
	}{
		{
			name:        "both unmanaged",
			inAtlas:     project.MaintenanceWindow{},
			inAKO:       project.MaintenanceWindow{},
			isDifferent: false,
		},
		{
			name:        "AKO unmanaged, Atlas has a wave",
			inAtlas:     project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)},
			inAKO:       project.MaintenanceWindow{},
			isDifferent: false,
		},
		{
			name:        "AKO manages a wave, Atlas has none",
			inAtlas:     project.MaintenanceWindow{},
			inAKO:       project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)},
			isDifferent: true,
		},
		{
			name:        "both manage the same wave",
			inAtlas:     project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)},
			inAKO:       project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)},
			isDifferent: false,
		},
		{
			name:        "both manage different waves",
			inAtlas:     project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(1)},
			inAKO:       project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(2)},
			isDifferent: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inAtlas := maintenancewindow.MaintenanceWindow{MaintenanceWindow: &tc.inAtlas}
			inAKO := maintenancewindow.MaintenanceWindow{MaintenanceWindow: &tc.inAKO}
			assert.Equal(t, tc.isDifferent, waveIsDifferent(inAtlas, inAKO))
		})
	}
}

func TestIsEmptyWindow(t *testing.T) {
	assert.True(t, isEmptyWindow(project.MaintenanceWindow{}))
	assert.False(t, isEmptyWindow(project.MaintenanceWindow{WaveAssignment: pointer.MakePtr(1)}))
}
