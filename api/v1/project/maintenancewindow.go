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

package project

type MaintenanceWindow struct {
	// Day of the week when you would like the maintenance window to start as a 1-based integer.
	// Sunday 1, Monday 2, Tuesday 3, Wednesday 4, Thursday 5, Friday 6, Saturday 7.
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
	// Cannot be specified if defer is true
	// +optional
	StartASAP bool `json:"startASAP,omitempty"`
	// Flag indicating whether the next scheduled project maintenance should be deferred for one week.
	// Cannot be specified if startASAP is true
	// +optional
	Defer bool `json:"defer,omitempty"`
}
