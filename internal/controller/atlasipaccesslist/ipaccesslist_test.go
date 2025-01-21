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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func TestHandleIPAccessList(t *testing.T) {
	deletionTime := metav1.Now()
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
		"should create ip access list in atlas": {
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
		"should manage ip access list in atlas": {
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

func TestCreateIPAccessList(t *testing.T) {
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
			result := r.createIPAccessList(ctx, tt.ipAccessListService(), ipAccessList, "", tt.akoIPAccessList)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestDeleteIPAccessList(t *testing.T) {
	tests := map[string]struct {
		atlasIPAccessList   ipaccesslist.IPAccessEntries
		partial             bool
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
		"should remove ip access list partially": {
			atlasIPAccessList: ipaccesslist.IPAccessEntries{"192.168.0.0/24": {CIDR: "192.168.0.0/24"}, "10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			partial:           true,
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
			result := r.deleteIPAccessList(ctx, tt.ipAccessListService(), ipAccessList, "", tt.atlasIPAccessList, tt.partial)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

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
