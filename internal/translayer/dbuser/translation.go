package dbuser

import (
	"fmt"
	"sort"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type User struct {
	akov2.AtlasDatabaseUserSpec
	Password  string
	ProjectID string
}

func toK8s(dbUser *admin.CloudDatabaseUser) (*User, error) {
	deleteAfterDate, err := dateStringToK8s(dbUser.DeleteAfterDate.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse deleteAfterDate: %w", err)
	}
	scopes, err := scopesToK8s(dbUser.GetScopes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse scopes: %w", err)
	}
	return &User{
		ProjectID: dbUser.GroupId,
		Password:  dbUser.GetPassword(),
		AtlasDatabaseUserSpec: akov2.AtlasDatabaseUserSpec{
			DatabaseName:    dbUser.DatabaseName,
			DeleteAfterDate: deleteAfterDate,
			Roles:           rolesToK8s(dbUser.GetRoles()),
			Scopes:          scopes,
			Username:        dbUser.Username,
			OIDCAuthType:    dbUser.GetOidcAuthType(),
			AWSIAMType:      dbUser.GetAwsIAMType(),
			X509Type:        dbUser.GetX509Type(),
		},
	}, nil
}

func toAtlas(au *User) *admin.CloudDatabaseUser {
	return &admin.CloudDatabaseUser{
		DatabaseName: au.DatabaseName,
		//DeleteAfterDate: au.DeleteAfterDate,
		X509Type:     pointer.MakePtr(au.X509Type),
		AwsIAMType:   pointer.MakePtr(au.AWSIAMType),
		GroupId:      au.ProjectID,
		Roles:        rolesToAtlas(au.Roles),
		Scopes:       scopesToAtlas(au.Scopes),
		Password:     pointer.MakePtr(au.Password),
		Username:     au.Username,
		OidcAuthType: pointer.MakePtr(au.OIDCAuthType),
	}
}

func rolesToAtlas(roles []akov2.RoleSpec) *[]admin.DatabaseUserRole {
	atlasRoles := []admin.DatabaseUserRole{}
	for _, role := range roles {
		atlasRoles = append(atlasRoles, admin.DatabaseUserRole{
			RoleName:       role.RoleName,
			DatabaseName:   role.DatabaseName,
			CollectionName: pointer.MakePtr(role.CollectionName),
		})
	}
	return &atlasRoles
}

func scopesToAtlas(scopes []akov2.ScopeSpec) *[]admin.UserScope {
	atlasScopes := []admin.UserScope{}
	for _, scope := range scopes {
		atlasScopes = append(atlasScopes, admin.UserScope{
			Name: scope.Name,
			Type: string(scope.Type),
		})
	}
	return &atlasScopes
}

func dateStringToK8s(date string) (string, error) {
	if date != "" {
		d, err := timeutil.ParseISO8601(date)
		if err != nil {
			return "", err
		}
		return timeutil.FormatISO8601(d), nil
	}
	return "", nil
}

func scopesToK8s(scopes []admin.UserScope) ([]akov2.ScopeSpec, error) {
	specScopes := []akov2.ScopeSpec{}
	for _, scope := range scopes {
		scopeType, err := scopeTypeToK8s(scope.Type)
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

func scopeTypeToK8s(scopeType string) (akov2.ScopeType, error) {
	switch akov2.ScopeType(scopeType) {
	case akov2.DeploymentScopeType:
		return akov2.DeploymentScopeType, nil
	case akov2.DataLakeScopeType:
		return akov2.DataLakeScopeType, nil
	default:
		return "", fmt.Errorf("unsupported scope type %s", scopeType)
	}
}

func rolesToK8s(roles []admin.DatabaseUserRole) []akov2.RoleSpec {
	specRoles := []akov2.RoleSpec{}
	for _, role := range roles {
		specRoles = append(specRoles, akov2.RoleSpec{
			RoleName:       role.RoleName,
			DatabaseName:   role.DatabaseName,
			CollectionName: role.GetCollectionName(),
		})
	}
	sort.Slice(specRoles, func(i, j int) bool {
		return specRoles[i].RoleName < specRoles[j].RoleName &&
			specRoles[i].DatabaseName < specRoles[j].DatabaseName &&
			specRoles[i].CollectionName < specRoles[j].CollectionName
	})
	return specRoles
}
