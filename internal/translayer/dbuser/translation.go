package dbuser

import (
	"fmt"
	"sort"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type User struct {
	akov2.AtlasDatabaseUserSpec
	Password  string
	ProjectID string
}

func toK8sDatabaseUser(dbUser *mongodbatlas.DatabaseUser) (*User, error) {
	deleteAfterDate, err := toK8sDateString(dbUser.DeleteAfterDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deleteAfterDate: %w", err)
	}
	scopes, err := toK8sScopes(dbUser.Scopes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scopes: %w", err)
	}
	return &User{
		ProjectID: dbUser.GroupID,
		Password:  dbUser.Password,
		AtlasDatabaseUserSpec: akov2.AtlasDatabaseUserSpec{
			DatabaseName:    dbUser.DatabaseName,
			DeleteAfterDate: deleteAfterDate,
			Roles:           toK8sRoles(dbUser.Roles),
			Scopes:          scopes,
			Username:        dbUser.Username,
			OIDCAuthType:    dbUser.OIDCAuthType,
			AWSIAMType:      dbUser.AWSIAMType,
			X509Type:        dbUser.X509Type,
		},
	}, nil
}

func toAtlas(au *User) *mongodbatlas.DatabaseUser {
	return &mongodbatlas.DatabaseUser{
		DatabaseName:    au.DatabaseName,
		DeleteAfterDate: au.DeleteAfterDate,
		X509Type:        au.X509Type,
		AWSIAMType:      au.AWSIAMType,
		GroupID:         au.ProjectID,
		Roles:           rolesToAtlas(au.Roles),
		Scopes:          scopesToAtlas(au.Scopes),
		Password:        au.Password,
		Username:        au.Username,
		OIDCAuthType:    au.OIDCAuthType,
	}
}

func rolesToAtlas(roles []akov2.RoleSpec) []mongodbatlas.Role {
	atlasRoles := []mongodbatlas.Role{}
	for _, role := range roles {
		atlasRoles = append(atlasRoles, mongodbatlas.Role{
			RoleName:       role.RoleName,
			DatabaseName:   role.DatabaseName,
			CollectionName: role.CollectionName,
		})
	}
	return atlasRoles
}

func scopesToAtlas(scopes []akov2.ScopeSpec) []mongodbatlas.Scope {
	atlasScopes := []mongodbatlas.Scope{}
	for _, scope := range scopes {
		atlasScopes = append(atlasScopes, mongodbatlas.Scope{
			Name: scope.Name,
			Type: string(scope.Type),
		})
	}
	return atlasScopes
}

func toK8sDateString(date string) (string, error) {
	if date != "" {
		d, err := timeutil.ParseISO8601(date)
		if err != nil {
			return "", err
		}
		return timeutil.FormatISO8601(d), nil
	}
	return "", nil
}

func toK8sScopes(scopes []mongodbatlas.Scope) ([]akov2.ScopeSpec, error) {
	specScopes := []akov2.ScopeSpec{}
	for _, scope := range scopes {
		scopeType, err := toK8sScopeType(scope.Type)
		if err != nil {
			return nil, err
		}
		specScopes = append(specScopes, akov2.ScopeSpec{
			Name: scope.Name,
			Type: scopeType,
		})
	}
	sort.Slice(specScopes, func(i, j int) bool {
		return specScopes[i].Name < specScopes[j].Name &&
			specScopes[i].Type < specScopes[j].Type
	})
	return specScopes, nil
}

func toK8sScopeType(scopeType string) (akov2.ScopeType, error) {
	switch akov2.ScopeType(scopeType) {
	case akov2.DeploymentScopeType:
		return akov2.DeploymentScopeType, nil
	case akov2.DataLakeScopeType:
		return akov2.DataLakeScopeType, nil
	default:
		return "", fmt.Errorf("unsupported scope type %s", scopeType)
	}
}

func toK8sRoles(roles []mongodbatlas.Role) []akov2.RoleSpec {
	specRoles := []akov2.RoleSpec{}
	for _, role := range roles {
		specRoles = append(specRoles, akov2.RoleSpec{
			RoleName:       role.RoleName,
			DatabaseName:   role.DatabaseName,
			CollectionName: role.CollectionName,
		})
	}
	sort.Slice(specRoles, func(i, j int) bool {
		return specRoles[i].RoleName < specRoles[j].RoleName &&
			specRoles[i].DatabaseName < specRoles[j].DatabaseName &&
			specRoles[i].CollectionName < specRoles[j].CollectionName
	})
	return specRoles
}
