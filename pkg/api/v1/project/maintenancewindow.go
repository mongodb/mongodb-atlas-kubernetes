package project

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

type MaintenanceWindow struct {
	// Day of the week when you would like the maintenance window to start as a 1-based integer.
	// Sunday 1, Monday 2, Tuesday 3, Wednesday 4, Thursday 5, Friday 6, Saturday 7
	// +optional
	DayOfWeek int `json:"dayOfWeek,omitempty"`
	// Hour of the day when you would like the maintenance window to start.
	// This parameter uses the 24-hour clock, where midnight is 0, noon is 12.
	// +optional
	HourOfDay int `json:"hourOfDay,omitempty"`
	// Flag that indicates whether you want to defer all maintenance windows one week they would be triggered.
	// +optional
	AutoDeferOnceEnabled bool `json:"autoDeferOnceEnabled,omitempty"`
	// Flag indicating whether project maintenance has been directed to start immediately.
	// +optional
	StartASAP bool `json:"startASAP,omitempty"`
}

// ToAtlas converts the ProjectMaintenanceWindow to native Atlas client format.
func (m MaintenanceWindow) ToAtlas() (*mongodbatlas.MaintenanceWindow, error) {
	result := &mongodbatlas.MaintenanceWindow{}
	err := compat.JSONCopy(result, m)
	return result, err
}

// ************************************ Builder methods *************************************************
// Note, that we don't use pointers here as the AtlasProject uses this without pointers

func NewMaintenanceWindow() MaintenanceWindow {
	return MaintenanceWindow{}
}

func (m MaintenanceWindow) WithDay(day int) MaintenanceWindow {
	m.DayOfWeek = day
	return m
}

func (m MaintenanceWindow) WithHour(hour int) MaintenanceWindow {
	m.HourOfDay = hour
	return m
}

func (m MaintenanceWindow) WithAutoDefer(autoDefer bool) MaintenanceWindow {
	m.AutoDeferOnceEnabled = autoDefer
	return m
}

func (m MaintenanceWindow) WithStartASAP(startASAP bool) MaintenanceWindow {
	m.StartASAP = startASAP
	return m
}
