// Copyright 2026 MongoDB Inc
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

package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312018/admin"
	"go.mongodb.org/atlas-sdk/v20250312018/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

func TestEnsureRegionalizedPrivateEndpointMode(t *testing.T) {
	for _, tc := range []struct {
		name               string
		spec               *project.RegionalizedPrivateEndpoint
		privateEndpointAPI *mockadmin.PrivateEndpointServicesApi

		isOK      bool
		isWarning bool

		wantRegionalizedMode bool
		wantReadyType        bool
		wantStatus           string
	}{
		{
			name: "nil spec should unset condition and return OK",
			spec: nil,
			privateEndpointAPI: func() *mockadmin.PrivateEndpointServicesApi {
				api := mockadmin.NewPrivateEndpointServicesApi(t)
				api.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
					Return(admin.GetRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
					Return(&admin.ProjectSettingItem{Enabled: false}, &http.Response{}, nil)
				return api
			}(),
			isOK:                 true,
			isWarning:            false,
			wantReadyType:        false,
			wantRegionalizedMode: false,
		},
		{
			name:                 "enabled in spec and disabled in atlas should toggle",
			spec:                 &project.RegionalizedPrivateEndpoint{Enabled: true},
			privateEndpointAPI:   mockAPIWithToggle(t, false, true),
			isOK:                 true,
			isWarning:            false,
			wantReadyType:        true,
			wantStatus:           "True",
			wantRegionalizedMode: true,
		},
		{
			name: "enabled in spec and enabled in atlas should not toggle",
			spec: &project.RegionalizedPrivateEndpoint{Enabled: true},
			privateEndpointAPI: func() *mockadmin.PrivateEndpointServicesApi {
				api := mockadmin.NewPrivateEndpointServicesApi(t)
				api.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
					Return(admin.GetRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
					Return(&admin.ProjectSettingItem{Enabled: true}, &http.Response{}, nil)
				return api
			}(),
			isOK:                 true,
			isWarning:            false,
			wantReadyType:        true,
			wantStatus:           "True",
			wantRegionalizedMode: true,
		},
		{
			name:                 "disabled in spec and enabled in atlas should toggle",
			spec:                 &project.RegionalizedPrivateEndpoint{Enabled: false},
			privateEndpointAPI:   mockAPIWithToggle(t, true, false),
			isOK:                 true,
			isWarning:            false,
			wantReadyType:        true,
			wantStatus:           "True",
			wantRegionalizedMode: false,
		},
		{
			name: "disabled in spec and disabled in atlas should not toggle",
			spec: &project.RegionalizedPrivateEndpoint{Enabled: false},
			privateEndpointAPI: func() *mockadmin.PrivateEndpointServicesApi {
				api := mockadmin.NewPrivateEndpointServicesApi(t)
				api.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
					Return(admin.GetRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
					Return(&admin.ProjectSettingItem{Enabled: false}, &http.Response{}, nil)
				return api
			}(),
			isOK:                 true,
			isWarning:            false,
			wantReadyType:        true,
			wantStatus:           "True",
			wantRegionalizedMode: false,
		},
		{
			name: "get current mode fails should terminate",
			spec: &project.RegionalizedPrivateEndpoint{Enabled: true},
			privateEndpointAPI: func() *mockadmin.PrivateEndpointServicesApi {
				api := mockadmin.NewPrivateEndpointServicesApi(t)
				api.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
					Return(admin.GetRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
					Return(nil, &http.Response{}, errors.New("failed to get mode"))
				return api
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
		},
		{
			name: "toggle fails should terminate",
			spec: &project.RegionalizedPrivateEndpoint{Enabled: true},
			privateEndpointAPI: func() *mockadmin.PrivateEndpointServicesApi {
				api := mockadmin.NewPrivateEndpointServicesApi(t)
				api.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
					Return(admin.GetRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
					Return(&admin.ProjectSettingItem{Enabled: false}, &http.Response{}, nil)
				api.EXPECT().ToggleRegionalEndpointMode(context.Background(), "testProjectID", mock.AnythingOfType("*admin.ProjectSettingItem")).
					Return(admin.ToggleRegionalEndpointModeApiRequest{ApiService: api})
				api.EXPECT().ToggleRegionalEndpointModeExecute(mock.Anything).
					Return(nil, &http.Response{}, errors.New("failed to toggle"))
				return api
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312: &admin.APIClient{
						PrivateEndpointServicesApi: tc.privateEndpointAPI,
					},
				},
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			atlasProject := &akov2.AtlasProject{
				Status: status.AtlasProjectStatus{
					ID: "testProjectID",
				},
				Spec: akov2.AtlasProjectSpec{
					RegionalizedPrivateEndpoint: tc.spec,
				},
			}

			r := AtlasProjectReconciler{}
			result := r.ensureRegionalizedPrivateEndpointMode(workflowCtx, atlasProject)

			assert.Equal(t, tc.isOK, result.IsOk())
			assert.Equal(t, tc.isWarning, result.IsWarning())

			con, ok := workflowCtx.GetCondition(api.RegionalizedPrivateEndpointReadyType)
			assert.Equal(t, tc.wantReadyType, ok)
			assert.Equal(t, tc.wantStatus, string(con.Status))
			if result.IsOk() {
				assert.Equal(t, tc.wantRegionalizedMode, atlasProject.Status.RegionalizedPrivateEndpoint.Enabled)
			}
		})
	}
}

func mockAPIWithToggle(t *testing.T, currentMode, afterToggle bool) *mockadmin.PrivateEndpointServicesApi {
	t.Helper()
	peAPI := mockadmin.NewPrivateEndpointServicesApi(t)
	peAPI.EXPECT().GetRegionalEndpointMode(context.Background(), "testProjectID").
		Return(admin.GetRegionalEndpointModeApiRequest{ApiService: peAPI})
	peAPI.EXPECT().GetRegionalEndpointModeExecute(mock.Anything).
		Return(&admin.ProjectSettingItem{Enabled: currentMode}, &http.Response{}, nil)
	peAPI.EXPECT().ToggleRegionalEndpointMode(context.Background(), "testProjectID", mock.AnythingOfType("*admin.ProjectSettingItem")).
		Return(admin.ToggleRegionalEndpointModeApiRequest{ApiService: peAPI})
	peAPI.EXPECT().ToggleRegionalEndpointModeExecute(mock.Anything).
		Return(&admin.ProjectSettingItem{Enabled: afterToggle}, &http.Response{}, nil)
	return peAPI
}
