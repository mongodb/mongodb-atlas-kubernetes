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

package maintenancewindow

import (
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type MaintenanceWindow struct {
	*project.MaintenanceWindow
}

func NewMaintenanceWindow(spec *project.MaintenanceWindow) *MaintenanceWindow {
	if spec == nil {
		return nil
	}
	// No slices so no normalization needed
	return &MaintenanceWindow{MaintenanceWindow: spec}
}

func (mw *MaintenanceWindow) WithStartASAP(asap bool) *MaintenanceWindow {
	mw.StartASAP = asap

	return mw
}

func (mw *MaintenanceWindow) EqualTo(target *MaintenanceWindow) bool {
	return reflect.DeepEqual(mw, target)
}

func toAtlas(in *MaintenanceWindow) *admin.GroupMaintenanceWindow {
	if in == nil {
		return nil
	}
	return &admin.GroupMaintenanceWindow{
		AutoDeferOnceEnabled: pointer.MakePtrOrNil(in.AutoDefer),
		DayOfWeek:            in.DayOfWeek,
		HourOfDay:            pointer.MakePtrOrNil(in.HourOfDay),
		StartASAP:            pointer.MakePtrOrNil(in.StartASAP),
	}
}

func fromAtlas(in *admin.GroupMaintenanceWindow) *MaintenanceWindow {
	if in == nil {
		return nil
	}
	return &MaintenanceWindow{
		MaintenanceWindow: &project.MaintenanceWindow{
			DayOfWeek: in.GetDayOfWeek(),
			HourOfDay: in.GetHourOfDay(),
			AutoDefer: in.GetAutoDeferOnceEnabled(),
			StartASAP: in.GetStartASAP(),
		},
	}
}
