package maintenancewindow

import (
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
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
