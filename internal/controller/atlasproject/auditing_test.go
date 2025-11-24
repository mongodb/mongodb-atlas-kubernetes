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

package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
)

func TestAuditController_reconcile(t *testing.T) {
	tests := map[string]struct {
		service            audit.AuditLogService
		audit              *akov2.Auditing
		expectedResult     workflow.DeprecatedResult
		expectedConditions []api.Condition
	}{
		"should unmanage audit config when unset on both Atlas and AKO": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{
							AuditFilter: "{}",
						},
					}, nil
				},
			},
			audit:              nil,
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
		"should fail to retrieve audit config from Atlas": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return nil, errors.New("failed to get audit log config")
				},
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to get audit log config")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.AuditingReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get audit log config"),
			},
		},
		"should fail to configure audit config in Atlas": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{},
					}, nil
				},
				UpdateFunc: func(projectID string, auditing *audit.AuditConfig) error {
					return errors.New("failed to set audit log config")
				},
			},
			audit: &akov2.Auditing{
				Enabled: true,
			},
			expectedResult: workflow.Terminate(workflow.ProjectAuditingReady, errors.New("failed to set audit log config")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.AuditingReadyType).
					WithReason(string(workflow.ProjectAuditingReady)).
					WithMessageRegexp("failed to set audit log config"),
			},
		},
		"should successfully configure audit config in Atlas": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{},
					}, nil
				},
				UpdateFunc: func(projectID string, auditing *audit.AuditConfig) error {
					return nil
				},
			},
			audit: &akov2.Auditing{
				Enabled: true,
			},
			expectedResult: workflow.OK(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.AuditingReadyType),
			},
		},
		"should be ready when not change is applied": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{
							Enabled:     true,
							AuditFilter: "{}",
						},
					}, nil
				},
			},
			audit: &akov2.Auditing{
				Enabled: true,
			},
			expectedResult: workflow.OK(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.AuditingReadyType),
			},
		},
		"should unmanage when unset in AKO and disable in Atlas": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{},
					}, nil
				},
				UpdateFunc: func(projectID string, auditing *audit.AuditConfig) error {
					return nil
				},
			},
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
		"should disable audit config in Atlas when unset in AKO": {
			service: &translation.AuditLogMock{
				GetFunc: func(projectID string) (*audit.AuditConfig, error) {
					return &audit.AuditConfig{
						Auditing: &akov2.Auditing{
							Enabled: true,
						},
					}, nil
				},
				UpdateFunc: func(projectID string, auditing *audit.AuditConfig) error {
					return nil
				},
			},
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a := &auditController{
				ctx: &workflow.Context{
					Context: context.Background(),
					Log:     zaptest.NewLogger(t).Sugar(),
				},
				project: &akov2.AtlasProject{
					Spec: akov2.AtlasProjectSpec{
						Auditing: tt.audit,
					},
				},
				service: tt.service,
			}

			result := a.reconcile()
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, a.ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestHandleAudit(t *testing.T) {
	tests := map[string]struct {
		audit              *akov2.Auditing
		expectedCalls      func(api *mockadmin.AuditingApi) admin.AuditingApi
		expectedResult     workflow.DeprecatedResult
		expectedConditions []api.Condition
	}{
		"should successfully handle audit reconciliation": {
			audit: &akov2.Auditing{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
			},
			expectedCalls: func(api *mockadmin.AuditingApi) admin.AuditingApi {
				api.EXPECT().GetAuditingConfiguration(context.Background(), "project-id").
					Return(admin.GetAuditingConfigurationApiRequest{ApiService: api})
				api.EXPECT().GetAuditingConfigurationExecute(mock.AnythingOfType("admin.GetAuditingConfigurationApiRequest")).
					Return(
						&admin.AuditLog{
							AuditAuthorizationSuccess: pointer.MakePtr(false),
							ConfigurationType:         pointer.MakePtr("NONE"),
							Enabled:                   pointer.MakePtr(false),
						},
						&http.Response{},
						nil,
					)
				api.EXPECT().UpdateAuditingConfiguration(context.Background(), "project-id", mock.AnythingOfType("*admin.AuditLog")).
					Return(admin.UpdateAuditingConfigurationApiRequest{ApiService: api})
				api.EXPECT().UpdateAuditingConfigurationExecute(mock.AnythingOfType("admin.UpdateAuditingConfigurationApiRequest")).
					Return(
						&admin.AuditLog{
							AuditAuthorizationSuccess: pointer.MakePtr(true),
							ConfigurationType:         pointer.MakePtr("FILTER_JSON"),
							Enabled:                   pointer.MakePtr(true),
						},
						&http.Response{},
						nil,
					)

				return api
			},
			expectedResult: workflow.OK(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.AuditingReadyType),
			},
		},
		"should fail to handle audit reconciliation": {
			audit: &akov2.Auditing{
				Enabled:                   true,
				AuditAuthorizationSuccess: true,
			},
			expectedCalls: func(api *mockadmin.AuditingApi) admin.AuditingApi {
				api.EXPECT().GetAuditingConfiguration(context.Background(), "project-id").
					Return(admin.GetAuditingConfigurationApiRequest{ApiService: api})
				api.EXPECT().GetAuditingConfigurationExecute(mock.AnythingOfType("admin.GetAuditingConfigurationApiRequest")).
					Return(
						&admin.AuditLog{
							AuditAuthorizationSuccess: pointer.MakePtr(false),
							ConfigurationType:         pointer.MakePtr("NONE"),
							Enabled:                   pointer.MakePtr(false),
						},
						&http.Response{},
						nil,
					)
				api.EXPECT().UpdateAuditingConfiguration(context.Background(), "project-id", mock.AnythingOfType("*admin.AuditLog")).
					Return(admin.UpdateAuditingConfigurationApiRequest{ApiService: api})
				api.EXPECT().UpdateAuditingConfigurationExecute(mock.AnythingOfType("admin.UpdateAuditingConfigurationApiRequest")).
					Return(
						nil,
						&http.Response{},
						errors.New("failed to configure audit log"),
					)

				return api
			},
			expectedResult: workflow.Terminate(workflow.ProjectAuditingReady, errors.New("failed to configure audit log")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.AuditingReadyType).
					WithReason(string(workflow.ProjectAuditingReady)).
					WithMessageRegexp("failed to configure audit log"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312009: &admin.APIClient{
						AuditingApi: tt.expectedCalls(mockadmin.NewAuditingApi(t)),
					},
				},
			}
			project := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Auditing: tt.audit,
				},
				Status: status.AtlasProjectStatus{
					ID: "project-id",
				},
			}

			result := handleAudit(ctx, project)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
