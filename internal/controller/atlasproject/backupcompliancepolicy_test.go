// Copyright 2024 MongoDB Inc
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
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestEnsureBackupCompliance(t *testing.T) {
	testBCP := &akov2.AtlasBackupCompliancePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-bcp",
			Namespace: "test-namespace",
		},
		Spec: akov2.AtlasBackupCompliancePolicySpec{
			AuthorizedEmail:         "test@example.com",
			AuthorizedUserFirstName: "John",
			AuthorizedUserLastName:  "Doe",
			CopyProtectionEnabled:   false,
			EncryptionAtRestEnabled: false,
			PITEnabled:              false,
			RestoreWindowDays:       42,
			ScheduledPolicyItems: []akov2.AtlasBackupPolicyItem{
				{
					FrequencyType:     "monthly",
					FrequencyInterval: 4,
					RetentionUnit:     "months",
					RetentionValue:    1,
				},
			},
			OnDemandPolicy: akov2.AtlasOnDemandPolicy{
				RetentionUnit:  "weeks",
				RetentionValue: 3,
			},
		},
	}

	for _, tc := range []struct {
		name string

		project         *akov2.AtlasProject
		conditionStatus workflow.ConditionReason
		bcp             *akov2.AtlasBackupCompliancePolicy

		backupAPI *mockadmin.CloudBackupsApi

		isOK      bool
		isWarning bool

		wantReadyType bool                     // whether the BackupComplianceReadyType is expected
		wantStatus    string                   // the expected status of BackupComplianceReadyType ("True", "False")
		wantReason    workflow.ConditionReason // the expected reason of BackupComplianceReadyType
	}{
		{
			name:    "BCP GET request errors",
			project: akov2.DefaultProject("test-namespace", "test-connection"),
			bcp:     &akov2.AtlasBackupCompliancePolicy{},
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						nil,
						&http.Response{},
						errors.New("bcp get test error"),
					)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
		},
		{
			name:    "BCP is in AKO, but not Atlas (create)",
			project: akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:     testBCP,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{}, &http.Response{}, nil)
				backupAPI.EXPECT().UpdateCompliancePolicy(context.Background(), "", mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
					Return(admin.UpdateCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().UpdateCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{State: pointer.MakePtr("ACTIVE")}, &http.Response{}, nil)
				return backupAPI
			}(),
			isOK:          true,
			wantReadyType: true,
			wantStatus:    "True",
		},
		{
			name:    "BCP is in AKO, but not Atlas, create errors",
			project: akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:     testBCP,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{}, &http.Response{}, nil)
				backupAPI.EXPECT().UpdateCompliancePolicy(context.Background(), "", mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
					Return(admin.UpdateCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().UpdateCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{}, &http.Response{}, errors.New("create test error"))
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
		},
		{
			name:    "BCP are still creating in Atlas, AKO is waiting",
			project: akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:     testBCP,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{
							AuthorizedEmail:         "test@example.com",
							AuthorizedUserFirstName: "John",
							AuthorizedUserLastName:  "Doe",
							CopyProtectionEnabled:   pointer.MakePtr(false),
							EncryptionAtRestEnabled: pointer.MakePtr(false),
							PitEnabled:              pointer.MakePtr(false),
							RestoreWindowDays:       pointer.MakePtr(42),
							ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
								{
									FrequencyType:     "monthly",
									FrequencyInterval: 4,
									RetentionUnit:     "months",
									RetentionValue:    1,
								},
							},
							OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
								RetentionUnit:  "weeks",
								RetentionValue: 3,
							},
							State: pointer.MakePtr("UPDATING"),
						},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     false,
			wantReadyType: true,
			wantStatus:    "False",
			wantReason:    workflow.ProjectBackupCompliancePolicyUpdating,
		},
		{
			name:            "BCP is still updating in Atlas, AKO status persists",
			project:         akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:             testBCP,
			conditionStatus: workflow.ProjectBackupCompliancePolicyUpdating,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{
							AuthorizedEmail:         "test@example.com",
							AuthorizedUserFirstName: "John",
							AuthorizedUserLastName:  "Doe",
							CopyProtectionEnabled:   pointer.MakePtr(false),
							EncryptionAtRestEnabled: pointer.MakePtr(false),
							PitEnabled:              pointer.MakePtr(false),
							RestoreWindowDays:       pointer.MakePtr(42),
							ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
								{
									FrequencyType:     "monthly",
									FrequencyInterval: 4,
									RetentionUnit:     "months",
									RetentionValue:    1,
								},
							},
							OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
								RetentionUnit:  "weeks",
								RetentionValue: 3,
							},
							State: pointer.MakePtr("UPDATING"),
						},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     false,
			wantReadyType: true,
			wantStatus:    "False",
			wantReason:    workflow.ProjectBackupCompliancePolicyUpdating,
		},
		{
			name:            "BCP removed from project in AKO with UPDATING status",
			project:         akov2.DefaultProject("test-namespace", "test-connection"),
			bcp:             testBCP,
			conditionStatus: workflow.ProjectBackupCompliancePolicyUpdating,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{State: pointer.MakePtr("UPDATING")},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
			wantReason:    workflow.Internal,
		},
		{
			name:            "BCP is no longer UPDATING in Atlas, AKO checks",
			project:         akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:             testBCP,
			conditionStatus: workflow.ProjectBackupCompliancePolicyUpdating,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{
							AuthorizedEmail:         "test@example.com",
							AuthorizedUserFirstName: "John",
							AuthorizedUserLastName:  "Doe",
							CopyProtectionEnabled:   pointer.MakePtr(false),
							EncryptionAtRestEnabled: pointer.MakePtr(false),
							PitEnabled:              pointer.MakePtr(false),
							RestoreWindowDays:       pointer.MakePtr(42),
							ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
								{
									FrequencyType:     "monthly",
									FrequencyInterval: 4,
									RetentionUnit:     "months",
									RetentionValue:    1,
								},
							},
							OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
								RetentionUnit:  "weeks",
								RetentionValue: 3,
							},
							State: pointer.MakePtr("ACTIVE"),
						},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          true,
			isWarning:     false,
			wantReadyType: true,
			wantStatus:    "True",
		},
		{
			name:    "BCP is in AKO, but not Atlas, and BCP is not met",
			project: akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp"),
			bcp:     testBCP,
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{}, &http.Response{}, nil)

				mockError := &admin.GenericOpenAPIError{}
				model := *admin.NewApiErrorWithDefaults()
				model.SetErrorCode(atlas.BackupComplianceNotMet)
				mockError.SetModel(model)

				backupAPI.EXPECT().UpdateCompliancePolicy(context.Background(), "", mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
					Return(admin.UpdateCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().UpdateCompliancePolicyExecute(mock.Anything).
					Return(&admin.DataProtectionSettings20231001{}, &http.Response{}, mockError)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
			wantReason:    workflow.ProjectBackupCompliancePolicyNotMet,
		},
		{
			name:    "BCP is in Atlas but not AKO (attempted delete)",
			project: akov2.DefaultProject("test-namespace", "test-connection"),
			bcp:     &akov2.AtlasBackupCompliancePolicy{},
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{
							AuthorizedEmail:         "test@example.com",
							AuthorizedUserFirstName: "John",
							AuthorizedUserLastName:  "Doe",
							CopyProtectionEnabled:   pointer.MakePtr(false),
							EncryptionAtRestEnabled: pointer.MakePtr(false),
							PitEnabled:              pointer.MakePtr(false),
							RestoreWindowDays:       pointer.MakePtr(42),
							ScheduledPolicyItems: &[]admin.BackupComplianceScheduledPolicyItem{
								{
									FrequencyType:     "monthly",
									FrequencyInterval: 4,
									RetentionUnit:     "months",
									RetentionValue:    1,
								},
							},
							OnDemandPolicyItem: &admin.BackupComplianceOnDemandPolicyItem{
								RetentionUnit:  "weeks",
								RetentionValue: 3,
							},
							State: pointer.MakePtr("ACTIVE"),
						},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          false,
			isWarning:     true,
			wantReadyType: true,
			wantStatus:    "False",
			wantReason:    workflow.ProjectBackupCompliancePolicyCannotDelete,
		},
		{
			name:    "BCP is not in AKO nor Atlas (unmanage)",
			project: akov2.DefaultProject("test-namespace", "test-connection"),
			bcp:     &akov2.AtlasBackupCompliancePolicy{},
			backupAPI: func() *mockadmin.CloudBackupsApi {
				backupAPI := mockadmin.NewCloudBackupsApi(t)
				backupAPI.EXPECT().GetCompliancePolicy(context.Background(), "").
					Return(admin.GetCompliancePolicyApiRequest{ApiService: backupAPI})
				backupAPI.EXPECT().GetCompliancePolicyExecute(mock.Anything).
					Return(
						&admin.DataProtectionSettings20231001{},
						&http.Response{},
						nil,
					)
				return backupAPI
			}(),
			isOK:          true,
			isWarning:     false,
			wantReadyType: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						CloudBackupsApi: tc.backupAPI,
					},
				},
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithStatusSubresource(tc.bcp).
				WithObjects(tc.bcp).
				Build()

			reconciler := &AtlasProjectReconciler{
				AtlasProvider: &atlasmock.TestProvider{},
				Client:        k8sClient,
			}

			workflowCtx.SetConditionFromResult(api.BackupComplianceReadyType, workflow.InProgress(tc.conditionStatus, "test state"))

			result := reconciler.ensureBackupCompliance(workflowCtx, tc.project)

			assert.Equal(t, tc.isOK, result.IsOk())
			assert.Equal(t, tc.isWarning, result.IsWarning())

			con, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
			assert.Equal(t, tc.wantReadyType, ok)

			assert.Equal(t, tc.wantStatus, string(con.Status))
			if tc.wantReason != "" {
				assert.True(t, workflowCtx.HasReason(tc.wantReason))
			}
		})
	}
}
