//Copyright 2025 MongoDB Inc
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package atlasipaccesslist

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		akoIPAccessList     ipaccesslist.IPAccessEntries
		partial             bool
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should fail to add": {
			akoIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}}).
					Return(errors.New("failed to add ip access list"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToCreate)).
					WithMessageRegexp("failed to add ip access list"),
			},
		},
		"should add ip access list": {
			akoIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
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
			result := r.create(ctx, tt.ipAccessListService(), ipAccessList, "", tt.akoIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestDeleteAll(t *testing.T) {
	tests := map[string]struct {
		atlasIPAccessList   ipaccesslist.IPAccessEntries
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should fail to delete": {
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(errors.New("failed to delete ip access list"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToDelete)).
					WithMessageRegexp("failed to delete ip access list"),
			},
		},
		"should remove ip access list": {
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}, "10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(nil)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "10.0.0.0/24"}).
					Return(nil)

				return s
			},
			expectedResult: ctrl.Result{},
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
			result := r.deleteAll(ctx, tt.ipAccessListService(), ipAccessList, "", tt.atlasIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestDeletePartial(t *testing.T) {
	tests := map[string]struct {
		atlasIPAccessList   ipaccesslist.IPAccessEntries
		ipAccessListService func() ipaccesslist.IPAccessListService
		expectedResult      ctrl.Result
		expectedConditions  []api.Condition
	}{
		"should fail to delete": {
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(errors.New("failed to delete ip access list"))

				return s
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.IPAccessListReady).
					WithReason(string(workflow.IPAccessListFailedToDelete)).
					WithMessageRegexp("failed to delete ip access list"),
			},
		},
		"should remove ip access list": {
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}, "10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			ipAccessListService: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "192.168.0.0/24"}).
					Return(nil)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "10.0.0.0/24"}).
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
			result := r.deletePartial(ctx, tt.ipAccessListService(), ipAccessList, "", tt.atlasIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
