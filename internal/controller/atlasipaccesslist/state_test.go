/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

package atlasipaccesslist

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func TestHandleCustomResource(t *testing.T) {
	tests := map[string]struct {
		ipAccessList       akov2.AtlasIPAccessList
		provider           atlas.Provider
		expectedResult     ctrl.Result
		expectedFinalizers []string
		expectedConditions []api.Condition
	}{
		"should skip reconciliation": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
					},
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			expectedResult:     ctrl.Result{},
			expectedFinalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to validate resource": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
					Labels: map[string]string{
						customresource.ResourceVersion: "wrong",
					},
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.ResourceVersionStatus).
					WithReason(string(workflow.AtlasResourceVersionIsInvalid)).
					WithMessageRegexp("wrong is not a valid semver version for label mongodb.com/atlas-resource-version"),
			},
		},
		"should fail when not supported": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasGovUnsupported)).
					WithMessageRegexp("the AtlasIPAccessList is not supported by Atlas for government"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to resolve credentials": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasAPIAccessNotConfigured)).
					WithMessageRegexp("missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to create sdk": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to create sdk")
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasAPIAccessNotConfigured)).
					WithMessageRegexp("failed to create sdk"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to resolve project": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{}, "", nil
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasAPIAccessNotConfigured)).
					WithMessageRegexp("failed to query Kubernetes: failed to get Project from Kubernetes: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should handle ip access list": {
			ipAccessList: akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			//nolint:dupl
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					ialAPI := mockadmin.NewProjectIPAccessListApi(t)
					ialAPI.EXPECT().ListProjectIpAccessLists(mock.Anything, "123").
						Return(admin.ListProjectIpAccessListsApiRequest{ApiService: ialAPI})
					ialAPI.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
						Return(
							&admin.PaginatedNetworkAccess{
								Results: &[]admin.NetworkPermissionEntry{
									{
										CidrBlock: pointer.MakePtr("192.168.0.0/24"),
									},
								},
							},
							nil,
							nil,
						)
					ialAPI.EXPECT().GetProjectIpAccessListStatus(mock.Anything, "123", "192.168.0.0/24").
						Return(admin.GetProjectIpAccessListStatusApiRequest{ApiService: ialAPI})
					ialAPI.EXPECT().GetProjectIpAccessListStatusExecute(mock.AnythingOfType("admin.GetProjectIpAccessListStatusApiRequest")).
						Return(
							&admin.NetworkPermissionEntryStatus{STATUS: "ACTIVE"},
							nil,
							nil,
						)

					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProjectByName(mock.Anything, "my-project").
						Return(admin.GetProjectByNameApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectByNameExecute(mock.Anything).
						Return(&admin.Group{Id: pointer.MakePtr("123")}, nil, nil)

					return &admin.APIClient{
						ProjectIPAccessListApi: ialAPI,
						ProjectsApi:            projectAPI,
					}, "", nil
				},
			},
			expectedResult:     ctrl.Result{},
			expectedFinalizers: []string{customresource.FinalizerLabel},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.IPAccessListReady),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, &tt.ipAccessList).
				WithStatusSubresource(&tt.ipAccessList).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasIPAccessListReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
				AtlasProvider: tt.provider,
				EventRecorder: record.NewFakeRecorder(10),
			}
			result := r.handleCustomResource(ctx.Context, &tt.ipAccessList)

			ipAccessList := &akov2.AtlasIPAccessList{}
			require.NoError(t, k8sClient.Get(ctx.Context, client.ObjectKeyFromObject(&tt.ipAccessList), ipAccessList))
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedFinalizers, ipAccessList.GetFinalizers())
			assert.True(t, cmp.Equal(tt.expectedConditions, ipAccessList.Status.GetConditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestHandleIPAccessList(t *testing.T) {
	deletionTime := metav1.Now()
	deleteAfterDate := metav1.NewTime(time.Now().Add(time.Minute * -5))
	tests := map[string]struct {
		akoIPAccessList     *akov2.AtlasIPAccessList
		partial             bool
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should fail to parse ip access list from crd": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("cidr 192.168.0/24 is invalid: invalid CIDR address: 192.168.0/24"),
			},
		},
		"should fail to list from atlas": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(nil, errors.New("failed to list ip access list"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to list ip access list"),
			},
		},
		"should release ip access list for deletion": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ip-access-list",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{}, nil)

				return s
			},
			expectedResult: ctrl.Result{},
		},
		"should delete ip access list from atlas": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "ip-access-list",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}, nil)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(nil)

				return s
			},
			expectedResult: ctrl.Result{},
		},
		"should add ip access list in atlas": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{}, nil)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}).
					Return(nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to add access list entries"),
			},
		},
		"should fail to add an expired entry": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
						{
							IPAddress:       "192.168.10.100",
							DeleteAfterDate: pointer.MakePtr(deleteAfterDate),
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}, nil)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{"192.168.10.100/32": {CIDR: "192.168.10.100/32", DeleteAfterDate: pointer.MakePtr(deleteAfterDate.Time)}}).
					Return(errors.New("fail to add, expired entry"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToCreate)).
					WithMessageRegexp("fail to add, expired entry"),
			},
		},
		"should delete ip access list entry in atlas": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}, "10.1.1.0/24": {CIDR: "10.1.1.0/24"}}, nil)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "10.1.1.0/24"}).
					Return(nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to delete access list entries"),
			},
		},
		"should fail to get status": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}, nil)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("", errors.New("failed to get status"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToGetState)).
					WithMessageRegexp("failed to get status"),
			},
		},
		"should be in pending state": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}, nil)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("PENDING", nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to add access list entries"),
			},
		},
		"should be in ready state": {
			akoIPAccessList: &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
				Spec: akov2.AtlasIPAccessListSpec{
					Entries: []akov2.IPAccessEntry{
						{
							CIDRBlock: "192.168.0.0/24",
						},
					},
				},
			},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().List(context.Background(), "").
					Return(ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}, nil)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("ACTIVE", nil)

				return s
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.IPAccessListReady),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.akoIPAccessList).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasIPAccessListReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
			}
			result := r.handleIPAccessList(ctx, tt.ipAccessListService(), "", tt.akoIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

/*
func TestManageIPAccessList(t *testing.T) {
	tests := map[string]struct {
		akoIPAccessList     ipaccesslist.IPAccessEntries
		atlasIPAccessList   ipaccesslist.IPAccessEntries
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should be no op task": {
			akoIPAccessList:   ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("ACTIVE", nil)

				return s
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.IPAccessListReady),
			},
		},
		"should add ip access list": {
			akoIPAccessList:   ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			atlasIPAccessList: ipaccesslist.IPAccessEntries{},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}).
					Return(nil)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("PENDING", nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to add access list entries"),
			},
		},
		"should remove ip access list": {
			akoIPAccessList:   ipaccesslist.IPAccessEntries{},
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to delete access list entries"),
			},
		},
	}
	//nolint:dupl
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ipAccessList := &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(ipAccessList).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasIPAccessListReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
			}
			result := r.manageIPAccessList(ctx, tt.ipAccessListService(), ipAccessList, "", tt.akoIPAccessList, tt.atlasIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestWatchState(t *testing.T) {
	tests := map[string]struct {
		ipAccessList        ipaccesslist.IPAccessEntries
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should fail to get status": {
			ipAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("", errors.New("failed to get status"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToGetState)).
					WithMessageRegexp("failed to get status"),
			},
		},
		"should be in progress when entries are pending": {
			ipAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("PENDING", nil)

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListPending)).
					WithMessageRegexp("Atlas has started to add access list entries"),
			},
		},
		"should be ready when entries are active": {
			ipAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Status(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return("ACTIVE", nil)

				return s
			},
			expectedResult: ctrl.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.IPAccessListReady),
			},
		},
	}
	//nolint:dupl
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ipAccessList := &akov2.AtlasIPAccessList{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ip-access-list",
					Namespace: "default",
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(ipAccessList).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasIPAccessListReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
			}
			result := r.watchState(ctx, tt.ipAccessListService(), ipAccessList, "", tt.ipAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
*/
