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

package dbuser

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nsf/jsondiff"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

type User struct {
	*akov2.AtlasDatabaseUserSpec
	Password  string
	ProjectID string
}

// NewUser wraps and normalizes a Kubernetes Atlas User Spec pointer augmenting it with projectID and password.
func NewUser(spec *akov2.AtlasDatabaseUserSpec, projectID, password string) (*User, error) {
	if spec == nil {
		return nil, nil
	}
	user := &User{AtlasDatabaseUserSpec: spec, ProjectID: projectID, Password: password}
	if err := normalize(user.AtlasDatabaseUserSpec); err != nil {
		return nil, fmt.Errorf("failed to create internal user type: %w", err)
	}
	return user, nil
}

// EqualSpecs returns true if the given users have the same specs
func EqualSpecs(spec, atlas *User) bool {
	if !spec.hasSpec() && !atlas.hasSpec() { // both missing spec are same
		return true
	}
	if !spec.hasSpec() || !atlas.hasSpec() { // only one missing spec are different
		return false
	}
	// note users are normalized at construction time
	return reflect.DeepEqual(spec.clearedSpecClone(), atlas.clearedSpecClone())
}

func DiffSpecs(a, b *User) string {
	opts := jsondiff.DefaultJSONOptions()
	_, result := jsondiff.Compare(
		cmp.JSON(a.clearedSpecClone()),
		cmp.JSON(b.clearedSpecClone()),
		&opts)
	return result
}

func (u *User) hasSpec() bool {
	return u != nil && u.AtlasDatabaseUserSpec != nil
}

func (u *User) clearedSpecClone() *akov2.AtlasDatabaseUserSpec {
	if u == nil || u.AtlasDatabaseUserSpec == nil {
		return nil
	}
	clone := *u.AtlasDatabaseUserSpec
	clone.ProjectRef = nil
	clone.PasswordSecret = nil
	clone.ExternalProjectRef = nil
	clone.ConnectionSecret = nil
	return &clone
}

func normalize(spec *akov2.AtlasDatabaseUserSpec) error {
	cmp.NormalizeSlice(spec.Labels, func(a, b common.LabelSpec) int {
		return strings.Compare(a.Key+a.Value, b.Key+b.Value)
	})
	cmp.NormalizeSlice(spec.Roles, func(a, b akov2.RoleSpec) int {
		return strings.Compare(
			a.RoleName+a.DatabaseName+a.CollectionName,
			b.RoleName+b.DatabaseName+b.CollectionName)
	})
	cmp.NormalizeSlice(spec.Scopes, func(a, b akov2.ScopeSpec) int {
		return strings.Compare(
			a.Name+string(a.Type),
			b.Name+string(b.Type))
	})
	if spec.Scopes == nil {
		spec.Scopes = []akov2.ScopeSpec{}
	}
	if spec.DeleteAfterDate != "" { // enforce date format
		operatorDeleteDate, err := timeutil.ParseISO8601(spec.DeleteAfterDate)
		if err != nil {
			return fmt.Errorf("failed to parse %q to an ISO date: %w", spec.DeleteAfterDate, err)
		}
		spec.DeleteAfterDate = timeutil.FormatISO8601(operatorDeleteDate)
	}
	return nil
}

func fromAtlas(dbUser *admin.CloudDatabaseUser) (*User, error) {
	if dbUser == nil {
		return nil, nil
	}
	scopes, err := scopesFromAtlas(dbUser.GetScopes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse scopes: %w", err)
	}
	u := &User{
		ProjectID: dbUser.GroupId,
		Password:  dbUser.GetPassword(),
		AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
			DatabaseName:    dbUser.DatabaseName,
			DeleteAfterDate: dateFromAtlas(dbUser.DeleteAfterDate),
			Description:     dbUser.GetDescription(),
			Labels:          labelsFromAtlas(dbUser.Labels),
			Roles:           rolesFromAtlas(dbUser.GetRoles()),
			Scopes:          scopes,
			Username:        dbUser.Username,
			OIDCAuthType:    dbUser.GetOidcAuthType(),
			AWSIAMType:      dbUser.GetAwsIAMType(),
			X509Type:        dbUser.GetX509Type(),
		},
	}
	if err := normalize(u.AtlasDatabaseUserSpec); err != nil {
		return nil, fmt.Errorf("failed to normalize spec from Atlas: %w", err)
	}
	return u, nil
}

func toAtlas(au *User) (*admin.CloudDatabaseUser, error) {
	if au == nil || au.AtlasDatabaseUserSpec == nil {
		return nil, nil
	}
	date, err := dateToAtlas(au.DeleteAfterDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deleteAfterDate value: %w", err)
	}
	return &admin.CloudDatabaseUser{
		DatabaseName:    au.DatabaseName,
		DeleteAfterDate: date,
		X509Type:        pointer.MakePtrOrNil(au.X509Type),
		AwsIAMType:      pointer.MakePtrOrNil(au.AWSIAMType),
		GroupId:         au.ProjectID,
		Description:     pointer.MakePtr(au.Description),
		Labels:          labelsToAtlas(au.Labels),
		Roles:           rolesToAtlas(au.Roles),
		Scopes:          scopesToAtlas(au.Scopes),
		Username:        au.Username,
		Password:        pointer.MakePtrOrNil(au.Password),
		OidcAuthType:    pointer.MakePtrOrNil(au.OIDCAuthType),
	}, nil
}

func dateToAtlas(d string) (*time.Time, error) {
	if d == "" {
		return nil, nil
	}
	date, err := timeutil.ParseISO8601(d)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q to an ISO date: %w", d, err)
	}
	return pointer.MakePtr(date), nil
}

func rolesToAtlas(roles []akov2.RoleSpec) *[]admin.DatabaseUserRole {
	if len(roles) == 0 {
		return nil
	}
	atlasRoles := []admin.DatabaseUserRole{}
	for _, role := range roles {
		ar := admin.DatabaseUserRole{
			RoleName:     role.RoleName,
			DatabaseName: role.DatabaseName,
		}
		if role.CollectionName != "" {
			ar.CollectionName = pointer.MakePtr(role.CollectionName)
		}
		atlasRoles = append(atlasRoles, ar)
	}
	return &atlasRoles
}

func scopesToAtlas(scopes []akov2.ScopeSpec) *[]admin.UserScope {
	if len(scopes) == 0 {
		return &[]admin.UserScope{}
	}
	atlasScopes := []admin.UserScope{}
	for _, scope := range scopes {
		atlasScopes = append(atlasScopes, admin.UserScope{
			Name: scope.Name,
			Type: string(scope.Type),
		})
	}
	return &atlasScopes
}

func dateFromAtlas(date *time.Time) string {
	if date == nil {
		return ""
	}
	return timeutil.FormatISO8601(*date)
}

func scopesFromAtlas(scopes []admin.UserScope) ([]akov2.ScopeSpec, error) {
	if len(scopes) == 0 {
		return []akov2.ScopeSpec{}, nil
	}
	specScopes := make([]akov2.ScopeSpec, 0, len(scopes))
	for _, scope := range scopes {
		scopeType, err := scopeTypeFromAtlas(scope.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to parse atlas scopes: %w", err)
		}
		specScopes = append(specScopes, akov2.ScopeSpec{
			Name: scope.Name,
			Type: scopeType,
		})
	}
	return specScopes, nil
}

func scopeTypeFromAtlas(scopeType string) (akov2.ScopeType, error) {
	switch akov2.ScopeType(scopeType) {
	case akov2.DeploymentScopeType:
		return akov2.DeploymentScopeType, nil
	case akov2.DataLakeScopeType:
		return akov2.DataLakeScopeType, nil
	default:
		return "", fmt.Errorf("unsupported scope type %s", scopeType)
	}
}

func rolesFromAtlas(roles []admin.DatabaseUserRole) []akov2.RoleSpec {
	if len(roles) == 0 {
		return nil
	}
	specRoles := []akov2.RoleSpec{}
	for _, role := range roles {
		specRoles = append(specRoles, akov2.RoleSpec{
			RoleName:       role.RoleName,
			DatabaseName:   role.DatabaseName,
			CollectionName: role.GetCollectionName(),
		})
	}
	return specRoles
}

func labelsToAtlas(labels []common.LabelSpec) *[]admin.ComponentLabel {
	if len(labels) == 0 {
		return nil
	}
	atlasLabels := make([]admin.ComponentLabel, 0, len(labels))
	for _, label := range labels {
		atlasLabels = append(atlasLabels, admin.ComponentLabel{
			Key:   pointer.MakePtr(label.Key),
			Value: pointer.MakePtr(label.Value),
		})
	}
	return &atlasLabels
}

func labelsFromAtlas(atlasLabels *[]admin.ComponentLabel) []common.LabelSpec {
	if atlasLabels == nil || len(*atlasLabels) == 0 {
		return nil
	}
	labels := make([]common.LabelSpec, 0, len(*atlasLabels))
	for _, atlasLabel := range *atlasLabels {
		labels = append(labels, common.LabelSpec{
			Key:   pointer.GetOrDefault(atlasLabel.Key, ""),
			Value: pointer.GetOrDefault(atlasLabel.Value, ""),
		})
	}
	return labels
}
