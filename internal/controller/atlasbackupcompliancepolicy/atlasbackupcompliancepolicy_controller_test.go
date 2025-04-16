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

package atlasbackupcompliancepolicy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
)

func TestReconcile(t *testing.T) {
	for _, tc := range []struct {
		name string

		objects     []client.Object
		isSupported bool

		wantErr              string
		wantResult           reconcile.Result
		wantStatusConditions []api.Condition
		wantFinalizers       []string
	}{
		{
			name:        "should terminate silently when resource is not found",
			isSupported: true,
		},
		{
			name: "should skip reconciliation when annotation is set",
			objects: []client.Object{
				&akov2.AtlasBackupCompliancePolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bcp",
						Namespace: "default",
						Annotations: map[string]string{
							customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
						},
					},
				},
			},
			isSupported: true,
		},
		{
			name: "should transition to error state when resource version is invalid",
			objects: []client.Object{
				&akov2.AtlasBackupCompliancePolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bcp",
						Namespace: "default",
						Labels: map[string]string{
							customresource.ResourceVersion: "blah",
						},
					},
				},
			},
			wantResult: ctrl.Result{
				RequeueAfter: workflow.DefaultRetry,
			},
			isSupported: true,
			wantStatusConditions: []api.Condition{
				{
					Type:   "Ready",
					Status: "False",
				},
				{
					Type:    "ResourceVersionIsValid",
					Status:  "False",
					Reason:  "AtlasResourceVersionIsInvalid",
					Message: "blah is not a valid semver version for label mongodb.com/atlas-resource-version",
				},
			},
		},
		{
			name: "should transition to error state when resource is unsupported",
			objects: []client.Object{
				&akov2.AtlasBackupCompliancePolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bcp",
						Namespace: "default",
					},
				},
			},
			isSupported: false,
			wantStatusConditions: []api.Condition{
				{
					Type:    "Ready",
					Status:  "False",
					Reason:  "AtlasGovUnsupported",
					Message: "the AtlasBackupCompliancePolicy is not supported by Atlas for government",
				},
				{
					Type:   "ResourceVersionIsValid",
					Status: "True",
				},
			},
		},
		{
			name: "should lock when there are references",
			objects: []client.Object{
				&akov2.AtlasBackupCompliancePolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bcp",
						Namespace: "default",
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
				},
				akov2.DefaultProject("default", "connection-secret").
					WithBackupCompliancePolicyNamespaced("bcp", "default"),
			},
			isSupported: true,
			wantStatusConditions: []api.Condition{
				{
					Type:   "Ready",
					Status: "True",
				},
				{
					Type:   "ResourceVersionIsValid",
					Status: "True",
				},
			},
			wantFinalizers: []string{"mongodbatlas/finalizer"},
		},
		{
			name: "should lock when there are references",
			objects: []client.Object{
				&akov2.AtlasBackupCompliancePolicy{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bcp",
						Namespace: "default",
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
				},
			},
			isSupported: true,
			wantStatusConditions: []api.Condition{
				{
					Type:   "Ready",
					Status: "True",
				},
				{
					Type:   "ResourceVersionIsValid",
					Status: "True",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			bcpIndexer := indexer.NewAtlasProjectByBackupCompliancePolicyIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.objects...).
				WithStatusSubresource(tc.objects...).
				WithIndex(
					bcpIndexer.Object(),
					bcpIndexer.Name(),
					bcpIndexer.Keys,
				).
				Build()

			reconciler := &AtlasBackupCompliancePolicyReconciler{
				Client:        k8sClient,
				Log:           zaptest.NewLogger(t).Sugar(),
				EventRecorder: record.NewFakeRecorder(1),
				AtlasProvider: &atlasmock.TestProvider{
					IsSupportedFunc: func() bool {
						return tc.isSupported
					},
				},
			}

			result, err := reconciler.Reconcile(
				context.Background(),
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Name:      "bcp",
						Namespace: "default",
					},
				},
			)

			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			assert.Equal(t, tc.wantErr, gotErr)
			assert.Equal(t, tc.wantResult, result)

			if len(tc.objects) == 0 {
				return
			}

			bcp := &akov2.AtlasBackupCompliancePolicy{}
			assert.NoError(t, k8sClient.Get(context.Background(), types.NamespacedName{Namespace: "default", Name: "bcp"}, bcp))

			for i := range bcp.Status.Conditions {
				bcp.Status.Conditions[i].LastTransitionTime = metav1.Time{}
			}

			assert.Equal(t, bcp.Status.Conditions, tc.wantStatusConditions)
			assert.Equal(t, bcp.Finalizers, tc.wantFinalizers)
		})
	}
}

func TestFindBCPForProjects(t *testing.T) {
	for _, tc := range []struct {
		name string

		project client.Object

		wantRequests []ctrl.Request
	}{
		{
			name:         "should return a slice of requests for BCP",
			project:      akov2.NewProject("test-project", "default", "connection-secret").WithBackupCompliancePolicyNamespaced("bcp1", "other-ns"),
			wantRequests: []ctrl.Request{{NamespacedName: types.NamespacedName{Name: "bcp1", Namespace: "other-ns"}}},
		},
		{
			name:         "should return nil when no BCP specified in project",
			project:      akov2.NewProject("test-project", "default", "connection-secret"),
			wantRequests: nil,
		},
		{
			name:         "should return nil when cannot cast object to project",
			project:      akov2.NewDeployment("default", "test-deployment", "test-deployment"),
			wantRequests: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.project).
				Build()

			reconciler := &AtlasBackupCompliancePolicyReconciler{
				Client: k8sClient,
				Log:    zaptest.NewLogger(t).Sugar(),
			}

			reqs := reconciler.findBCPForProjects(context.Background(), tc.project)
			assert.Equal(
				t,
				tc.wantRequests,
				reqs,
			)
		})
	}
}
