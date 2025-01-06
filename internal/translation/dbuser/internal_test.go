package dbuser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
)

const (
	testProjectID = "project-id"

	testPassword = "something-secret-here"

	testDB = "test-db"

	testUsername = "testUsername"
)

var (
	testDate = timeutil.FormatISO8601(time.Now())
)

func TestToAndFromAtlas(t *testing.T) {
	for _, tc := range []struct {
		title string
		user  *User
	}{
		{
			title: "Nil user renders nil atlas user",
		},

		{
			title: "Nil user spec renders nil atlas user",
		},

		{
			title: "Default user spec converts back and forth correctly",
			user:  defaultTestUser(),
		},

		{
			title: "Default user spec with correct date succeeds",
			user: func() *User {
				u := defaultTestUser()
				u.DeleteAfterDate = testDate
				return u
			}(),
		},

		{
			title: "Default user spec with ordered roles succeeds",
			user: func() *User {
				u := defaultTestUser()
				u.Roles = []akov2.RoleSpec{
					{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
					{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
					{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
				}
				return u
			}(),
		},

		{
			title: "Default user spec with ordered scopes succeeds",
			user: func() *User {
				u := defaultTestUser()
				u.Scopes = []akov2.ScopeSpec{
					{Name: "cluster1", Type: "CLUSTER"},
					{Name: "cluster2", Type: "CLUSTER"},
					{Name: "lake1", Type: "DATA_LAKE"},
					{Name: "lake2", Type: "DATA_LAKE"},
				}
				return u
			}(),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			atlasUser, err := toAtlas(tc.user)
			require.NoError(t, err)
			converted, err := fromAtlas(atlasUser)
			require.NoError(t, err)
			assert.Equal(t, tc.user, converted)
		})
	}
}

func TestToAtlasDateFailure(t *testing.T) {
	user := defaultTestUser()
	user.DeleteAfterDate = "bad-date-string"
	expectedErrMsg := "failed to parse deleteAfterDate value"

	_, err := toAtlas(user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrMsg)
}

func TestFromAtlasScopeFailure(t *testing.T) {
	atlasUser := defaultTestAtlasUser()
	atlasUser.Scopes = &[]admin.UserScope{{Name: "some-name", Type: "not-a-proper-cluster"}}
	expectedErrMsg := "unsupported scope type not-a-proper-cluster"

	_, err := fromAtlas(atlasUser)
	require.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrMsg)
}

func defaultTestUser() *User {
	return &User{
		AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{
			DatabaseName: testDB,
			Username:     testUsername,
		},
		Password:  testPassword,
		ProjectID: testProjectID,
	}
}

func defaultTestAtlasUser() *admin.CloudDatabaseUser {
	return &admin.CloudDatabaseUser{
		DatabaseName: testDB,
		GroupId:      testProjectID,
		Password:     pointer.MakePtr(testPassword),
		Username:     testUsername,
	}
}
