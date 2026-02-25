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
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312014/admin"
)

type MaintenanceWindowService interface {
	Get(context.Context, string) (*MaintenanceWindow, error)
	Update(context.Context, string, *MaintenanceWindow) error
	Reset(context.Context, string) error
	Defer(context.Context, string) error
	ToggleAutoDefer(context.Context, string) error
}

type MaintenanceWindowAPI struct {
	maintenanceAPI admin.MaintenanceWindowsApi
}

func NewMaintenanceWindowAPIService(api admin.MaintenanceWindowsApi) *MaintenanceWindowAPI {
	return &MaintenanceWindowAPI{maintenanceAPI: api}
}

func (mw *MaintenanceWindowAPI) Get(ctx context.Context, projectID string) (*MaintenanceWindow, error) {
	maintenanceWindow, _, err := mw.maintenanceAPI.GetMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get maintenance window from Atlas: %w", err)
	}
	return fromAtlas(maintenanceWindow), nil
}

func (mw *MaintenanceWindowAPI) Update(ctx context.Context, projectID string, maintenanceWindow *MaintenanceWindow) error {
	_, err := mw.maintenanceAPI.UpdateMaintenanceWindow(ctx, projectID, toAtlas(maintenanceWindow)).Execute()
	if err != nil {
		return fmt.Errorf("failed to update maintenance window in Atlas: %w", err)
	}
	return nil
}

func (mw *MaintenanceWindowAPI) Reset(ctx context.Context, projectID string) error {
	_, err := mw.maintenanceAPI.ResetMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to reset maintenance window in Atlas: %w", err)
	}
	return nil
}

func (mw *MaintenanceWindowAPI) Defer(ctx context.Context, projectID string) error {
	_, err := mw.maintenanceAPI.DeferMaintenanceWindow(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to defer maintenance window in Atlas: %w", err)
	}
	return nil
}

func (mw *MaintenanceWindowAPI) ToggleAutoDefer(ctx context.Context, projectID string) error {
	_, err := mw.maintenanceAPI.ToggleMaintenanceAutoDefer(ctx, projectID).Execute()
	if err != nil {
		return fmt.Errorf("failed to toggle auto defer for maintenance window in Atlas: %w", err)
	}
	return nil
}
