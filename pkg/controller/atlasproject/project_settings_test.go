package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestAreSettingsInSync(t *testing.T) {
	atlasDef := &akov2.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
		IsDataExplorerEnabled:                       pointer.MakePtr(true),
		IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
		IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
		IsSchemaAdvisorEnabled:                      pointer.MakePtr(true),
	}
	specDef := &akov2.ProjectSettings{
		IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
		IsDataExplorerEnabled:                       pointer.MakePtr(true),
	}

	areEqual := areSettingsInSync(atlasDef, specDef)
	assert.True(t, areEqual, "Only fields which are set should be compared")

	specDef.IsPerformanceAdvisorEnabled = pointer.MakePtr(false)
	areEqual = areSettingsInSync(atlasDef, specDef)
	assert.False(t, areEqual, "Field values should be the same ")
}

func TestEnsureProjectSettings(t *testing.T) {
	for _, tc := range []struct {
		name       string
		settings   *akov2.ProjectSettings
		projectAPI *mockadmin.ProjectsApi

		isOK      bool
		isWarning bool

		wantReadyType bool   // whether the ProjectSettingsReadyType is expected
		wantStatus    string // the expected status of ProjectSettingsReadyType ("True", "False")
	}{
		{
			name:     "Project Settings unset in AKO & Atlas",
			settings: nil,
			projectAPI: func() *mockadmin.ProjectsApi {
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.Anything).
					Return(
						&admin.GroupSettings{ // These are the default settings on a fresh project
							IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(true),
							IsDataExplorerEnabled:                       pointer.MakePtr(true),
							IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
							IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
							IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
							IsSchemaAdvisorEnabled:                      pointer.MakePtr(true),
						},
						&http.Response{},
						nil,
					)
				return projectAPI
			}(),
			isOK:      true,
			isWarning: false,

			wantReadyType: false,
		},
		{
			name:     "GET Atlas Project Settings errors",
			settings: nil,
			projectAPI: func() *mockadmin.ProjectsApi {
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.Anything).
					Return(
						&admin.GroupSettings{},
						&http.Response{},
						errors.New("TEST GET ERROR"),
					)
				return projectAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
		},
		{
			name: "Project Settings are equal in AKO & Atlas",
			settings: &akov2.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
				IsDataExplorerEnabled:                       pointer.MakePtr(false),
				IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
				IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
				IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
			},
			projectAPI: func() *mockadmin.ProjectsApi {
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.Anything).
					Return(
						&admin.GroupSettings{
							IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
							IsDataExplorerEnabled:                       pointer.MakePtr(false),
							IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
							IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
							IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
							IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
						},
						&http.Response{},
						nil,
					)
				return projectAPI
			}(),
			isOK:      true,
			isWarning: false,

			wantReadyType: true,
			wantStatus:    "True",
		},
		{
			name: "Project Settings are different in AKO & Atlas",
			settings: &akov2.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
				IsDataExplorerEnabled:                       pointer.MakePtr(false),
				IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
				IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
				IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
			},
			projectAPI: func() *mockadmin.ProjectsApi {
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.Anything).
					Return(
						&admin.GroupSettings{
							IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
							IsDataExplorerEnabled:                       pointer.MakePtr(true),
							IsExtendedStorageSizesEnabled:               pointer.MakePtr(true),
							IsPerformanceAdvisorEnabled:                 pointer.MakePtr(false),
							IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
							IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
						},
						&http.Response{},
						nil,
					)
				projectAPI.EXPECT().UpdateProjectSettings(context.Background(), "", mock.AnythingOfType("*admin.GroupSettings")).
					Return(admin.UpdateProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().UpdateProjectSettingsExecute(mock.Anything).
					Return(&admin.GroupSettings{}, &http.Response{}, nil)

				return projectAPI
			}(),
			isOK:      true,
			isWarning: false,

			wantReadyType: true,
			wantStatus:    "True",
		},
		{
			name: "PATCH Atlas Project Settings errors",
			settings: &akov2.ProjectSettings{
				IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
				IsDataExplorerEnabled:                       pointer.MakePtr(false),
				IsExtendedStorageSizesEnabled:               pointer.MakePtr(false),
				IsPerformanceAdvisorEnabled:                 pointer.MakePtr(true),
				IsRealtimePerformancePanelEnabled:           pointer.MakePtr(true),
				IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
			},
			projectAPI: func() *mockadmin.ProjectsApi {
				projectAPI := mockadmin.NewProjectsApi(t)
				projectAPI.EXPECT().GetProjectSettings(context.Background(), "").
					Return(admin.GetProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().GetProjectSettingsExecute(mock.Anything).
					Return(
						&admin.GroupSettings{
							IsCollectDatabaseSpecificsStatisticsEnabled: pointer.MakePtr(false),
							IsDataExplorerEnabled:                       pointer.MakePtr(true),
							IsExtendedStorageSizesEnabled:               pointer.MakePtr(true),
							IsPerformanceAdvisorEnabled:                 pointer.MakePtr(false),
							IsRealtimePerformancePanelEnabled:           pointer.MakePtr(false),
							IsSchemaAdvisorEnabled:                      pointer.MakePtr(false),
						},
						&http.Response{},
						nil,
					)
				projectAPI.EXPECT().UpdateProjectSettings(context.Background(), "", mock.AnythingOfType("*admin.GroupSettings")).
					Return(admin.UpdateProjectSettingsApiRequest{ApiService: projectAPI})
				projectAPI.EXPECT().UpdateProjectSettingsExecute(mock.Anything).
					Return(&admin.GroupSettings{}, &http.Response{}, errors.New("TEST PATCH ERROR"))

				return projectAPI
			}(),
			isOK:      false,
			isWarning: true,

			wantReadyType: true,
			wantStatus:    "False",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := &workflow.Context{
				SdkClient: &admin.APIClient{
					ProjectsApi: tc.projectAPI,
				},
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			project := akov2.DefaultProject("test-ns", "test-conn")
			project.Spec.Settings = tc.settings

			result := ensureProjectSettings(ctx, project)
			assert.Equal(t, tc.isOK, result.IsOk())
			assert.Equal(t, tc.isWarning, result.IsWarning())

			con, ok := ctx.GetCondition(api.ProjectSettingsReadyType)
			assert.Equal(t, tc.wantReadyType, ok)

			assert.Equal(t, tc.wantStatus, string(con.Status))
		})
	}
}
