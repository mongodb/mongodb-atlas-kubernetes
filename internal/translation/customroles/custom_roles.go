package customroles

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

type CustomRoleService interface {
	Get(ctx context.Context, projectID string, roleName string) (CustomRole, error)
	List(ctx context.Context, projectID string) ([]CustomRole, error)
	Create(ctx context.Context, projectID string, role CustomRole) error
	Update(ctx context.Context, projectID string, roleName string, role CustomRole) error
	Delete(ctx context.Context, projectID string, roleName string) error
}

type CustomRoles struct {
	roleAPI admin.CustomDatabaseRolesApi
}

func NewCustomRoles(api admin.CustomDatabaseRolesApi) *CustomRoles {
	return &CustomRoles{roleAPI: api}
}

func (s *CustomRoles) Get(ctx context.Context, projectID string, roleName string) (CustomRole, error) {
	customRole, httpResp, err := s.roleAPI.GetCustomDatabaseRole(ctx, projectID, roleName).Execute()
	// handle RoleNotFound error
	if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
		return CustomRole{}, nil
	}
	if err != nil {
		return CustomRole{}, fmt.Errorf("failed to get custom roles from Atlas: %w", err)
	}

	return fromAtlas(customRole), err
}

func (s *CustomRoles) List(ctx context.Context, projectID string) ([]CustomRole, error) {
	atlasRoles, _, err := s.roleAPI.ListCustomDatabaseRoles(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list custom roles from Atlas: %w", err)
	}

	customRoles := make([]CustomRole, len(atlasRoles))

	for i := range atlasRoles {
		customRoles[i] = fromAtlas(&atlasRoles[i])
	}

	return customRoles, nil
}

func (s *CustomRoles) Create(ctx context.Context, projectID string, role CustomRole) error {
	_, _, err := s.roleAPI.CreateCustomDatabaseRole(ctx, projectID, toAtlas(&role)).Execute()
	return err
}

func (s *CustomRoles) Update(ctx context.Context, projectID string, roleName string, role CustomRole) error {
	_, _, err := s.roleAPI.UpdateCustomDatabaseRole(ctx, projectID, roleName, toAtlasUpdate(&role)).Execute()
	return err
}

func (s *CustomRoles) Delete(ctx context.Context, projectID string, roleName string) error {
	_, err := s.roleAPI.DeleteCustomDatabaseRole(ctx, projectID, roleName).Execute()
	return err
}
