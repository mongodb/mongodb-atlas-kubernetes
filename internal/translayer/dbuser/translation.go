package dbuser

import (
	"fmt"
	"sort"
	"time"

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

func NewUser(spec akov2.AtlasDatabaseUserSpec, projectID, password string) *User {
	return &User{AtlasDatabaseUserSpec: spec, ProjectID: projectID, Password: password}
}

func Normalize(spec *akov2.AtlasDatabaseUserSpec) *akov2.AtlasDatabaseUserSpec {
	if spec.Roles == nil {
		spec.Roles = []akov2.RoleSpec{}
	}
	if spec.Scopes == nil {
		spec.Scopes = []akov2.ScopeSpec{}
	}
	return spec
}

func fromAtlas(dbUser *admin.CloudDatabaseUser) (*User, error) {
	scopes, err := scopesFromAtlas(dbUser.GetScopes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse scopes: %w", err)
	}
	u := &User{
		ProjectID: dbUser.GroupId,
		Password:  dbUser.GetPassword(),
		AtlasDatabaseUserSpec: akov2.AtlasDatabaseUserSpec{
			DatabaseName:    dbUser.DatabaseName,
			DeleteAfterDate: dateFromAtlas(dbUser.DeleteAfterDate),
			Roles:           rolesFromAtlas(dbUser.GetRoles()),
			Scopes:          scopes,
			Username:        dbUser.Username,
			OIDCAuthType:    dbUser.GetOidcAuthType(),
			AWSIAMType:      dbUser.GetAwsIAMType(),
			X509Type:        dbUser.GetX509Type(),
		},
	}
	return u, nil
}

func toAtlas(au *User) (*admin.CloudDatabaseUser, error) {
	date, err := dateToAtlas(au.DeleteAfterDate)
	if err != nil {
		return nil, err
	}
	return &admin.CloudDatabaseUser{
		DatabaseName:    au.DatabaseName,
		DeleteAfterDate: date,
		X509Type:        pointer.MakePtrOrNil(au.X509Type),
		AwsIAMType:      pointer.MakePtrOrNil(au.AWSIAMType),
		GroupId:         au.ProjectID,
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
		return nil, err
	}
	return pointer.MakePtr(date), nil
}

func rolesToAtlas(roles []akov2.RoleSpec) *[]admin.DatabaseUserRole {
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
	specScopes := []akov2.ScopeSpec{}
	for _, scope := range scopes {
		scopeType, err := scopeTypeFromAtlas(scope.Type)
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
