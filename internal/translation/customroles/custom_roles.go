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

package customroles

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
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
	customRole, httpResp, err := s.roleAPI.GetCustomDbRole(ctx, projectID, roleName).Execute()
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
	// custom database roles does not offer paginated resources.
	atlasRoles, _, err := s.roleAPI.ListCustomDbRoles(ctx, projectID).Execute()
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
	_, _, err := s.roleAPI.CreateCustomDbRole(ctx, projectID, toAtlas(&role)).Execute()
	return err
}

func (s *CustomRoles) Update(ctx context.Context, projectID string, roleName string, role CustomRole) error {
	_, _, err := s.roleAPI.UpdateCustomDbRole(ctx, projectID, roleName, toAtlasUpdate(&role)).Execute()
	return err
}

func (s *CustomRoles) Delete(ctx context.Context, projectID string, roleName string) error {
	_, err := s.roleAPI.DeleteCustomDbRole(ctx, projectID, roleName).Execute()
	return err
}
