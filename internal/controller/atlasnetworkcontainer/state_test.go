package atlasnetworkcontainer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/stretchr/testify/assert"
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
)

func TestHandleCustomResource(t *testing.T) {
	tests := map[string]struct {
		networkContainer   *akov2.AtlasNetworkContainer
		provider           atlas.Provider
		expectedResult     ctrl.Result
		expectedFinalizers []string
		expectedConditions []api.Condition
	}{
		"should skip reconciliation": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
					},
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
					},
				},
			},
			expectedResult:     ctrl.Result{},
			expectedFinalizers: []string{customresource.FinalizerLabel},
		},
		"should fail to validate resource": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
					Labels: map[string]string{
						customresource.ResourceVersion: "wrong",
					},
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
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
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
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
					WithMessageRegexp("the AtlasNetworkContainer is not supported by Atlas for government"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to resolve credentials": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
					},
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
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
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("can not fetch AtlasProject: " +
						"atlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to create sdk": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
					},
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return nil, "", errors.New("failed to create sdk")
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("failed to create sdk"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should fail to resolve project": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
					},
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{}, "", nil
				},
			},
			expectedResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("failed to query Kubernetes: failed to get Project from Kubernetes: can not fetch AtlasProject: atlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		"should handle network container": {
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/16",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					ialAPI := mockadmin.NewNetworkPeeringApi(t)
					return &atlas.ClientSet{
						SdkClient20231115008: &admin.APIClient{NetworkPeeringApi: ialAPI},
					}, "", nil
				},
			},
			expectedResult:     ctrl.Result{},
			expectedFinalizers: []string{customresource.FinalizerLabel},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.ReadyType),
				api.TrueCondition(api.ResourceVersionStatus),
				api.TrueCondition(api.NetworkContainerReady),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, tc.networkContainer).
				WithStatusSubresource(tc.networkContainer).
				Build()
			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
			}
			r := &AtlasNetworkContainerReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: k8sClient,
					Log:    logger,
				},
				AtlasProvider: tc.provider,
				EventRecorder: record.NewFakeRecorder(10),
			}
			result, err := r.handleCustomResource(ctx.Context, tc.networkContainer)

			require.NoError(t, err)
			networkContainer := &akov2.AtlasNetworkContainer{}
			require.NoError(t, k8sClient.Get(ctx.Context, client.ObjectKeyFromObject(tc.networkContainer), networkContainer))
			assert.Equal(t, tc.expectedResult, result)
			assert.Equal(t, tc.expectedFinalizers, networkContainer.GetFinalizers())
			assert.True(t, cmp.Equal(tc.expectedConditions, networkContainer.Status.GetConditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
