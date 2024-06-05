package dbuser

import (
	"fmt"
	"reflect"
	"sort"
	"time"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type User struct {
	*akov2.AtlasDatabaseUserSpec
	Password  string
	ProjectID string
}

// NewUser wraps a Kubernetes Atlas User Spec pointer augmenting it with projectID and password.
func NewUser(spec *akov2.AtlasDatabaseUserSpec, projectID, password string) (*User, error) {
	return normalize(&User{AtlasDatabaseUserSpec: spec, ProjectID: projectID, Password: password})
}

// DiffSpecs returns all differences found in the user Spec fields or a spec user and an atlas user.
// Non Spec fields are not compared. Inputs are dbuser.User so they are normalized
func DiffSpecs(specUser, atlasUser *User) []string {
	diffs := []string{}
	if atlasUser.Username != specUser.Username {
		diffs = append(diffs, fmt.Sprintf("Usernames differs from spec: %q <> %q\n",
			atlasUser.Username, specUser.Username))
	}
	if atlasUser.DatabaseName != specUser.DatabaseName {
		diffs = append(diffs, fmt.Sprintf("DatabaseName differs from spec: %q <> %q\n",
			atlasUser.DatabaseName, specUser.DatabaseName))
	}
	if atlasUser.DeleteAfterDate != specUser.DeleteAfterDate {
		diffs = append(diffs, fmt.Sprintf("DeleteAfterDate differs from spec: %q <> %q\n",
			atlasUser.DeleteAfterDate, specUser.DeleteAfterDate))
	}
	if atlasUser.OIDCAuthType != specUser.OIDCAuthType {
		diffs = append(diffs, fmt.Sprintf("OIDCAuthType differs from spec: %q <> %q\n",
			atlasUser.OIDCAuthType, specUser.OIDCAuthType))
	}
	if atlasUser.AWSIAMType != specUser.AWSIAMType {
		diffs = append(diffs, fmt.Sprintf("AWSIAMType differs from spec: %q <> %q\n",
			atlasUser.AWSIAMType, specUser.AWSIAMType))
	}
	if atlasUser.X509Type != specUser.X509Type {
		diffs = append(diffs, fmt.Sprintf("X509Type differs from spec: %q <> %q\n",
			atlasUser.X509Type, specUser.X509Type))
	}
	if !reflect.DeepEqual(atlasUser.Roles, specUser.Roles) {
		diffs = append(diffs, fmt.Sprintf("Roles differs from spec: %v <> %v\n",
			atlasUser.Roles, specUser.Roles))
	}
	if !reflect.DeepEqual(atlasUser.Scopes, specUser.Scopes) {
		diffs = append(diffs, fmt.Sprintf("Scopes differs from spec: %v <> %v END\n",
			atlasUser.Scopes, specUser.Scopes))
	}
	return diffs
}

func normalize(user *User) (*User, error) {
	if user.Roles == nil {
		user.Roles = []akov2.RoleSpec{}
	} else {
		orderRoles(user.Roles)
	}
	if user.Scopes == nil {
		user.Scopes = []akov2.ScopeSpec{}
	} else {
		orderScopes(user.Scopes)
	}
	if user.DeleteAfterDate != "" { // enforce date format
		operatorDeleteDate, err := timeutil.ParseISO8601(user.DeleteAfterDate)
		if err != nil {
			return nil, err
		}
		user.DeleteAfterDate = timeutil.FormatISO8601(operatorDeleteDate)
	}
	// Ensure comparisons succeed on default value
	// TODO: Check if this is required after updating the Helm Chart CRDs
	if user.X509Type == "" {
		user.X509Type = "NONE"
	}
	if user.OIDCAuthType == "" {
		user.OIDCAuthType = "NONE"
	}
	if user.AWSIAMType == "" {
		user.AWSIAMType = "NONE"
	}
	return user, nil
}

func fromAtlas(dbUser *admin.CloudDatabaseUser) (*User, error) {
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
			Roles:           rolesFromAtlas(dbUser.GetRoles()),
			Scopes:          scopes,
			Username:        dbUser.Username,
			OIDCAuthType:    dbUser.GetOidcAuthType(),
			AWSIAMType:      dbUser.GetAwsIAMType(),
			X509Type:        dbUser.GetX509Type(),
		},
	}
	return normalize(u)
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
	return specScopes, nil
}

func orderScopes(scopes []akov2.ScopeSpec) {
	sort.Slice(scopes, func(i, j int) bool {
		return scopes[i].Name < scopes[j].Name &&
			scopes[i].Type < scopes[j].Type
	})
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
	return specRoles
}

func orderRoles(roles []akov2.RoleSpec) {
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].RoleName < roles[j].RoleName &&
			roles[i].DatabaseName < roles[j].DatabaseName &&
			roles[i].CollectionName < roles[j].CollectionName
	})
}
