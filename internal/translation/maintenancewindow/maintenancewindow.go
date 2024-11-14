package maintenancewindow

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
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
	_, _, err := mw.maintenanceAPI.UpdateMaintenanceWindow(ctx, projectID, toAtlas(maintenanceWindow)).Execute()
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
