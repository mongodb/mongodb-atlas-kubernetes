package project

import (
	"go.mongodb.org/atlas/mongodbatlas"
)

type MaintenanceWindow struct {
	// Day of the week when you would like the maintenance window to start as a 1-based integer.
	// Sunday 1, Monday 2, Tuesday 3, Wednesday 4, Thursday 5, Friday 6, Saturday 7
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	DayOfWeek int `json:"dayOfWeek,omitempty"`
	// Hour of the day when you would like the maintenance window to start.
	// This parameter uses the 24-hour clock, where midnight is 0, noon is 12.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	HourOfDay int `json:"hourOfDay,omitempty"`
	// Flag indicating whether any scheduled project maintenance should be deferred automatically for one week.
	// +optional
	AutoDefer bool `json:"autoDefer,omitempty"`
	// Flag indicating whether project maintenance has been directed to start immediately.
	// +optional
	StartASAP bool `json:"startASAP,omitempty"`
	// Flag indicating whether the next scheduled project maintenance should be deferred for one week.
	// +optional
	Defer bool `json:"defer,omitempty"`
}

// ToAtlas converts the ProjectMaintenanceWindow to native Atlas client format.
func (m MaintenanceWindow) ToAtlas() (*mongodbatlas.MaintenanceWindow, error) {
	return &mongodbatlas.MaintenanceWindow{
		DayOfWeek:            m.DayOfWeek,
		HourOfDay:            &m.HourOfDay,
		StartASAP:            &m.StartASAP,
		NumberOfDeferrals:    0,
		AutoDeferOnceEnabled: &m.AutoDefer,
	}, nil
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
	m.AutoDefer = autoDefer
	return m
}

func (m MaintenanceWindow) WithStartASAP(startASAP bool) MaintenanceWindow {
	m.StartASAP = startASAP
	return m
}

func (m MaintenanceWindow) WithDefer(isDefer bool) MaintenanceWindow {
	m.Defer = isDefer
	return m
}
