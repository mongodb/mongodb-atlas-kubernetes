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

package dbuser_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/timeutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
)

const (
	testProjectID = "project-id"

	testPassword = "something-secret-here"

	testDB = "test-db"

	testUsername = "testUsername"
)

var (
	testDate = timeutil.FormatISO8601(time.Now())

	testOtherDate = timeutil.FormatISO8601(time.Now().Add(time.Hour))
)

func TestNewUser(t *testing.T) {
	for _, tc := range []struct {
		title            string
		spec             *akov2.AtlasDatabaseUserSpec
		projectID        string
		password         string
		expectedUser     *dbuser.User
		expectedErrorMsg string
	}{
		{
			title: "Nil spec returns nil user",
		},

		{
			title:        "Empty spec returns empty user",
			spec:         &akov2.AtlasDatabaseUserSpec{},
			expectedUser: &dbuser.User{AtlasDatabaseUserSpec: &akov2.AtlasDatabaseUserSpec{Scopes: []akov2.ScopeSpec{}}},
		},

		{
			title:     "Basic user is properly created",
			spec:      defaultTestSpec(),
			projectID: testProjectID,
			password:  testPassword,
			expectedUser: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
				ProjectID:             testProjectID,
				Password:              testPassword,
			},
		},

		{
			title: "Spec with bad date gets rejected",
			spec: func() *akov2.AtlasDatabaseUserSpec {
				spec := defaultTestSpec()
				spec.DeleteAfterDate = "bad-date"
				return spec
			}(),
			projectID:        testProjectID,
			password:         testPassword,
			expectedUser:     nil,
			expectedErrorMsg: "failed to parse \"bad-date\" to an ISO date",
		},

		{
			title: "Spec with proper date gets created",
			spec: func() *akov2.AtlasDatabaseUserSpec {
				spec := defaultTestSpec()
				spec.DeleteAfterDate = testDate
				return spec
			}(),
			projectID: testProjectID,
			password:  testPassword,
			expectedUser: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.DeleteAfterDate = testDate
					return spec
				}(),
				ProjectID: testProjectID,
				Password:  testPassword,
			},
		},

		{
			title: "Spec with unordered roles renders a normalized user with ordered entries",
			spec: func() *akov2.AtlasDatabaseUserSpec {
				spec := defaultTestSpec()
				spec.Roles = []akov2.RoleSpec{
					{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
					{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
					{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
				}
				return spec
			}(),
			projectID: testProjectID,
			password:  testPassword,
			expectedUser: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Roles = []akov2.RoleSpec{
						{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
						{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
						{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
					}
					return spec
				}(),
				ProjectID: testProjectID,
				Password:  testPassword,
			},
		},

		{
			title: "Spec with unordered scopes renders a normalized user with ordered entries",
			spec: func() *akov2.AtlasDatabaseUserSpec {
				spec := defaultTestSpec()
				spec.Scopes = []akov2.ScopeSpec{
					{Name: "cluster2", Type: "CLUSTER"},
					{Name: "lake2", Type: "DATA_LAKE"},
					{Name: "lake1", Type: "DATA_LAKE"},
					{Name: "cluster1", Type: "CLUSTER"},
				}
				return spec
			}(),
			projectID: testProjectID,
			password:  testPassword,
			expectedUser: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					return spec
				}(),
				ProjectID: testProjectID,
				Password:  testPassword,
			},
		},

		{
			title: "Spec with unordered labels renders a normalized user with ordered entries",
			spec: func() *akov2.AtlasDatabaseUserSpec {
				spec := defaultTestSpec()
				spec.Labels = []common.LabelSpec{
					{Key: "label3", Value: "value3"},
					{Key: "label2", Value: "value2"},
					{Key: "label1", Value: "value1"},
				}
				return spec
			}(),
			projectID: testProjectID,
			password:  testPassword,
			expectedUser: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Labels = []common.LabelSpec{
						{Key: "label1", Value: "value1"},
						{Key: "label2", Value: "value2"},
						{Key: "label3", Value: "value3"},
					}
					return spec
				}(),
				ProjectID: testProjectID,
				Password:  testPassword,
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			user, err := dbuser.NewUser(tc.spec, tc.projectID, tc.password)
			if tc.expectedErrorMsg != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.expectedUser, user)
		})
	}
}

func TestDiffSpecs(t *testing.T) {
	for _, tc := range []struct {
		title         string
		spec          *dbuser.User
		atlas         *dbuser.User
		expectedDiffs []string
	}{
		{
			title: "Nil users are equal",
		},

		{
			title: "Nil spec side is flagged",
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
			},
			expectedDiffs: []string{"\"changed\":[null, {}]"},
		},

		{
			title: "Nil atlas side is flagged",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
			},
			expectedDiffs: []string{"\"changed\":[{}, null]"},
		},

		{
			title: "Sample users have no diffs",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
			},
		},

		{
			title: "Only the spec part of each user is compared",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
				ProjectID:             "",
				Password:              testPassword,
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
				ProjectID:             "",
				Password:              testPassword,
			},
		},

		{
			title: "All simple field diffs are flagged",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: defaultTestSpec(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Username = spec.Username + "1"
					spec.DatabaseName = spec.DatabaseName + "2"
					spec.DeleteAfterDate = testOtherDate
					spec.OIDCAuthType = "IDP_GROUP"
					spec.AWSIAMType = "USER"
					spec.X509Type = "MANAGED"
					return spec
				}(),
			},
			expectedDiffs: []string{
				"username",
				"databaseName",
				"deleteAfterDate",
				"oidcAuthType",
				"awsIamType",
				"x509Type",
			},
		},

		{
			title: "Role diffs are flagged",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Roles = []akov2.RoleSpec{
						{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
						{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
						{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Roles = []akov2.RoleSpec{
						{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
						{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{"roles", `"prop-removed":{"roleName": "role1"}`},
		},

		{
			title: "Equal role lists show no diffs",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Roles = []akov2.RoleSpec{
						{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
						{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
						{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Roles = []akov2.RoleSpec{
						{RoleName: "role2", DatabaseName: "database2", CollectionName: "collection2"},
						{RoleName: "role2", DatabaseName: "database1", CollectionName: "collection1"},
						{RoleName: "role1", DatabaseName: "database1", CollectionName: "collection1"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{},
		},

		{
			title: "Scope diffs are flagged",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{"scopes", `prop-added":{"name": "lake1",}`},
		},

		{
			title: "Equal scopes show no diffs",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{},
		},

		{
			title: "Scopes with different references show no diffs",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					spec.ProjectRef = &common.ResourceRefNamespaced{
						Name:      "some-project",
						Namespace: "some-namespace",
					}
					spec.PasswordSecret = &common.ResourceRef{Name: "some-secret-ref"}
					spec.ConnectionSecret = &api.LocalObjectReference{Name: "some-local-secret-ref"}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Scopes = []akov2.ScopeSpec{
						{Name: "cluster1", Type: "CLUSTER"},
						{Name: "cluster2", Type: "CLUSTER"},
						{Name: "lake1", Type: "DATA_LAKE"},
						{Name: "lake2", Type: "DATA_LAKE"},
					}
					spec.ProjectRef = &common.ResourceRefNamespaced{
						Name:      "another-project",
						Namespace: "another-namespace",
					}
					spec.PasswordSecret = &common.ResourceRef{Name: "another-secret-ref"}
					spec.ConnectionSecret = &api.LocalObjectReference{Name: "another-local-secret-ref"}
					return spec
				}(),
			},
			expectedDiffs: []string{},
		},

		{
			title: "Different Labels fail comparison",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Labels = []common.LabelSpec{
						{Key: "label1", Value: "value1"},
						{Key: "label2", Value: "value2"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Labels = []common.LabelSpec{
						{Key: "label1", Value: "value1"},
						{Key: "label2", Value: "value2"},
						{Key: "label3", Value: "value3"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{"labels", `prop-added":{"value": "value3"}`},
		},

		{
			title: "Same Labels show no diffs",
			spec: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Labels = []common.LabelSpec{
						{Key: "label1", Value: "value1"},
						{Key: "label2", Value: "value2"},
					}
					return spec
				}(),
			},
			atlas: &dbuser.User{
				AtlasDatabaseUserSpec: func() *akov2.AtlasDatabaseUserSpec {
					spec := defaultTestSpec()
					spec.Labels = []common.LabelSpec{
						{Key: "label1", Value: "value1"},
						{Key: "label2", Value: "value2"},
					}
					return spec
				}(),
			},
			expectedDiffs: []string{},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			equal := dbuser.EqualSpecs(tc.spec, tc.atlas)
			if len(tc.expectedDiffs) == 0 {
				assert.Equal(t, true, equal)
			} else {
				diff := dbuser.DiffSpecs(tc.spec, tc.atlas)
				for _, expected := range tc.expectedDiffs {
					assert.Contains(t, diff, expected)
				}
			}
		})
	}
}

func defaultTestSpec() *akov2.AtlasDatabaseUserSpec {
	return &akov2.AtlasDatabaseUserSpec{
		DatabaseName: testDB,
		Username:     testUsername,
		Scopes:       []akov2.ScopeSpec{},
	}
}
