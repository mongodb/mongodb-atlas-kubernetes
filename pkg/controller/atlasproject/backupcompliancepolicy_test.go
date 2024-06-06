/*
Copyright 2024 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestEnsureBackupCompliance(t *testing.T) {
	t.Run("get BCP request errors", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
			Return(
				nil,
				&http.Response{},
				errors.New("get test error"),
			)

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		reconciler := &AtlasProjectReconciler{
			AtlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.False(t, result.IsOk())

		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
		assert.True(t, result.IsWarning())
	})

	t.Run("BCP is in AKO, but not Atlas (create)", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				nil,
			)
		backupAPI.EXPECT().UpdateDataProtectionSettings(context.Background(), project.ID(), mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
			Return(admin.UpdateDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().UpdateDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{
					State: pointer.MakePtr("ACTIVE"),
				},
				&http.Response{},
				nil,
			)

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bcp",
				Namespace: project.Namespace,
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

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.True(t, result.IsOk())
		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
	})

	t.Run("BCP is in AKO, but not Atlas, create errors", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				nil,
			)
		backupAPI.EXPECT().UpdateDataProtectionSettings(context.Background(), project.ID(), mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
			Return(admin.UpdateDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().UpdateDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				errors.New("create test error"),
			)

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bcp",
				Namespace: project.Namespace,
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

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.False(t, result.IsOk())

		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
		assert.True(t, result.IsWarning())
	})

	t.Run("BCP are still creating in Atlas, AKO is waiting", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
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

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bcp",
				Namespace: project.Namespace,
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

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.False(t, result.IsOk())
		assert.False(t, result.IsWarning())
		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
		assert.True(t, workflowCtx.HasReason(workflow.ProjectBackupCompliancePolicyUpdating))
	})

	t.Run("BCP is no longer UPDATING in Atlas, AKO checks", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
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

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bcp",
				Namespace: project.Namespace,
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

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			Build()

		workflowCtx.SetConditionFromResult(api.BackupComplianceReadyType, workflow.InProgress(workflow.ProjectBackupCompliancePolicyUpdating, "test state"))

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.True(t, result.IsOk())
		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
	})

	t.Run("BCP is in AKO, but not Atlas, and BCP is not met", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection").WithBackupCompliancePolicy("test-bcp")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				nil,
			)

		mockError := &admin.GenericOpenAPIError{}
		model := *admin.NewApiError()
		model.SetErrorCode(atlas.BackupComplianceNotMet)
		mockError.SetModel(model)

		backupAPI.EXPECT().UpdateDataProtectionSettings(context.Background(), project.ID(), mock.AnythingOfType("*admin.DataProtectionSettings20231001")).
			Return(admin.UpdateDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().UpdateDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				mockError,
			)

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		bcp := &akov2.AtlasBackupCompliancePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-bcp",
				Namespace: project.Namespace,
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

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(bcp).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.False(t, result.IsOk())

		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
		assert.True(t, workflowCtx.HasReason(workflow.ProjectBackupCompliancePolicyNotMet))
	})

	t.Run("BCP is in Atlas but not AKO (attempted delete)", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
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

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.False(t, result.IsOk())

		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.True(t, ok)
		assert.True(t, workflowCtx.HasReason(workflow.ProjectBackupCompliancePolicyCannotDelete))
	})

	t.Run("BCP is not in AKO nor Atlas (unmanage)", func(t *testing.T) {
		project := akov2.DefaultProject("test-namespace", "test-connection")
		backupAPI := mockadmin.NewCloudBackupsApi(t)

		backupAPI.EXPECT().GetDataProtectionSettings(context.Background(), project.ID()).
			Return(admin.GetDataProtectionSettingsApiRequest{ApiService: backupAPI})
		backupAPI.EXPECT().GetDataProtectionSettingsExecute(mock.Anything).
			Return(
				&admin.DataProtectionSettings20231001{},
				&http.Response{},
				nil,
			)

		workflowCtx := &workflow.Context{
			SdkClient: &admin.APIClient{
				CloudBackupsApi: backupAPI,
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}

		result := reconciler.ensureBackupCompliance(workflowCtx, project)

		assert.True(t, result.IsOk())

		_, ok := workflowCtx.GetCondition(api.BackupComplianceReadyType)
		assert.False(t, ok)
	})
}
