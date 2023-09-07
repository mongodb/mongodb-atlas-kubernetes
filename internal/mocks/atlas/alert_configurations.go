package atlas

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
)

type MockAlertConfigurationsMock struct {
	CreateFunc     func(projectID string, alertConfig *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	CreateRequests map[string]*mongodbatlas.AlertConfiguration

	EnableAnAlertConfigFunc     func(projectID string, alertConfigID string, enabled *bool) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	EnableAnAlertConfigRequests map[string]*bool

	GetAnAlertConfigFunc     func(projectID string, alertConfigID string) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	GetAnAlertConfigRequests map[string]struct{}

	GetOpenAlertsConfigFunc     func(projectID string, alertConfigID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	GetOpenAlertsConfigRequests map[string]struct{}

	ListFunc     func(projectID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	ListRequests map[string]struct{}

	ListMatcherFieldsFunc  func() ([]string, *mongodbatlas.Response, error)
	ListMatcherFieldsCalls int

	UpdateFunc     func(projectID string, alertConfigID string, alertConfig *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error)
	UpdateRequests map[string]*mongodbatlas.AlertConfiguration

	DeleteFunc     func(projectID string, alertConfigID string) (*mongodbatlas.Response, error)
	DeleteRequests map[string]struct{}
}

func (c *MockAlertConfigurationsMock) Create(_ context.Context, projectID string, alertConfig *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.CreateRequests == nil {
		c.CreateRequests = map[string]*mongodbatlas.AlertConfiguration{}
	}

	c.CreateRequests[projectID] = alertConfig

	return c.CreateFunc(projectID, alertConfig)
}
func (c *MockAlertConfigurationsMock) EnableAnAlertConfig(_ context.Context, projectID string, alertConfigID string, enabled *bool) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.EnableAnAlertConfigRequests == nil {
		c.EnableAnAlertConfigRequests = map[string]*bool{}
	}

	c.EnableAnAlertConfigRequests[fmt.Sprintf("%s.%s", projectID, alertConfigID)] = enabled

	return c.EnableAnAlertConfigFunc(projectID, alertConfigID, enabled)
}
func (c *MockAlertConfigurationsMock) GetAnAlertConfig(_ context.Context, projectID string, alertConfigID string) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.GetAnAlertConfigRequests == nil {
		c.GetAnAlertConfigRequests = map[string]struct{}{}
	}

	c.GetAnAlertConfigRequests[fmt.Sprintf("%s.%s", projectID, alertConfigID)] = struct{}{}

	return c.GetAnAlertConfigFunc(projectID, alertConfigID)
}
func (c *MockAlertConfigurationsMock) GetOpenAlertsConfig(_ context.Context, projectID string, alertConfigID string) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.GetOpenAlertsConfigRequests == nil {
		c.GetOpenAlertsConfigRequests = map[string]struct{}{}
	}

	c.GetOpenAlertsConfigRequests[fmt.Sprintf("%s.%s", projectID, alertConfigID)] = struct{}{}

	return c.GetOpenAlertsConfigFunc(projectID, alertConfigID)
}
func (c *MockAlertConfigurationsMock) List(_ context.Context, projectID string, _ *mongodbatlas.ListOptions) ([]mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.ListRequests == nil {
		c.ListRequests = map[string]struct{}{}
	}

	c.ListRequests[projectID] = struct{}{}

	return c.ListFunc(projectID)
}
func (c *MockAlertConfigurationsMock) ListMatcherFields(_ context.Context) ([]string, *mongodbatlas.Response, error) {
	c.ListMatcherFieldsCalls++

	return c.ListMatcherFieldsFunc()
}
func (c *MockAlertConfigurationsMock) Update(_ context.Context, projectID string, alertConfigID string, alertConfig *mongodbatlas.AlertConfiguration) (*mongodbatlas.AlertConfiguration, *mongodbatlas.Response, error) {
	if c.UpdateRequests == nil {
		c.UpdateRequests = map[string]*mongodbatlas.AlertConfiguration{}
	}

	c.UpdateRequests[fmt.Sprintf("%s.%s", projectID, alertConfigID)] = alertConfig

	return c.UpdateFunc(projectID, alertConfigID, alertConfig)
}
func (c *MockAlertConfigurationsMock) Delete(_ context.Context, projectID string, alertConfigID string) (*mongodbatlas.Response, error) {
	if c.DeleteRequests == nil {
		c.DeleteRequests = map[string]struct{}{}
	}

	c.DeleteRequests[fmt.Sprintf("%s.%s", projectID, alertConfigID)] = struct{}{}

	return c.DeleteFunc(projectID, alertConfigID)
}
