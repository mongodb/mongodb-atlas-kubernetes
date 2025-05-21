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

package integrations

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

var (
	// sample error test
	ErrTestFail = errors.New("failure")
)

const (
// testProjectID = "project-id"

// testContainerName = "fake-container-name"

// testContainerID = "fake-container-id"

// testPeeringID = "peering-id"

// testVpcID = "vpc-id"
)

func TestHandleCustomResource(t *testing.T) {
	tests := []struct {
		title          string
		integration    *akov2next.AtlasThirdPartyIntegration
		provider       atlas.Provider
		wantResult     ctrl.Result
		wantFinalizers []string
		wantConditions []api.Condition
	}{
		{
			title: "should skip reconciliation",
			integration: &akov2next.AtlasThirdPartyIntegration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "integration",
					Namespace: "default",
					Annotations: map[string]string{
						customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
					},
					Finalizers: []string{customresource.FinalizerLabel},
				},
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
		},
		{
			title: "should fail to validate resource",
			integration: &akov2next.AtlasThirdPartyIntegration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "integration",
					Namespace: "default",
					Labels: map[string]string{
						customresource.ResourceVersion: "wrong",
					},
				},
			},
			wantResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.ResourceVersionStatus).
					WithReason(string(workflow.AtlasResourceVersionIsInvalid)).
					WithMessageRegexp("wrong is not a valid semver version for label mongodb.com/atlas-resource-version"),
			},
		},
		{
			title: "should fail when not supported",
			integration: &akov2next.AtlasThirdPartyIntegration{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "integration",
					Namespace: "default",
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
			wantResult: ctrl.Result{},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasGovUnsupported)).
					WithMessageRegexp("the AtlasThirdPartyIntegration is not supported by Atlas for government"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		// {
		// 	title: "should fail to resolve credentials and remove finalizer",
		// 	networkPeering: &akov2.AtlasNetworkPeering{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:      "network-peering",
		// 			Namespace: "default",
		// 		},
		// 		Spec: akov2.AtlasNetworkPeeringSpec{
		// 			ProjectDualReference: akov2.ProjectDualReference{
		// 				ProjectRef: &common.ResourceRefNamespaced{
		// 					Name: "my-no-existing-project",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	provider: &atlasmock.TestProvider{
		// 		IsSupportedFunc: func() bool {
		// 			return true
		// 		},
		// 	},
		// 	wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		// 	wantFinalizers: nil,
		// 	wantConditions: []api.Condition{
		// 		api.FalseCondition(api.ReadyType).
		// 			WithReason(string(workflow.NetworkPeeringNotConfigured)).
		// 			WithMessageRegexp("error resolving project reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
		// 		api.TrueCondition(api.ResourceVersionStatus),
		// 	},
		// },
		// {
		// 	title: "should fail to create sdk",
		// 	networkPeering: &akov2.AtlasNetworkPeering{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:      "network-peering",
		// 			Namespace: "default",
		// 		},
		// 		Spec: akov2.AtlasNetworkPeeringSpec{
		// 			ProjectDualReference: akov2.ProjectDualReference{
		// 				ConnectionSecret: &api.LocalObjectReference{
		// 					Name: "my-secret",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	provider: &atlasmock.TestProvider{
		// 		IsSupportedFunc: func() bool {
		// 			return true
		// 		},
		// 		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
		// 			return nil, errors.New("failed to create sdk")
		// 		},
		// 	},
		// 	wantResult: ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		// 	wantConditions: []api.Condition{
		// 		api.FalseCondition(api.ReadyType).
		// 			WithReason(string(workflow.NetworkPeeringNotConfigured)).
		// 			WithMessageRegexp("failed to create sdk"),
		// 		api.TrueCondition(api.ResourceVersionStatus),
		// 	},
		// },
		// {
		// 	title: "should fail to resolve project and remove finalizers",
		// 	networkPeering: &akov2.AtlasNetworkPeering{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:       "network-peering",
		// 			Namespace:  "default",
		// 			Finalizers: []string{customresource.FinalizerLabel},
		// 		},
		// 		Spec: akov2.AtlasNetworkPeeringSpec{
		// 			ProjectDualReference: akov2.ProjectDualReference{
		// 				ConnectionSecret: &api.LocalObjectReference{
		// 					Name: "my-secret",
		// 				},
		// 				ProjectRef: &common.ResourceRefNamespaced{
		// 					Name: "my-no-existing-project",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	provider: &atlasmock.TestProvider{
		// 		IsSupportedFunc: func() bool {
		// 			return true
		// 		},
		// 		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
		// 			pAPI := mockadmin.NewProjectsApi(t)
		// 			return &atlas.ClientSet{
		// 				SdkClient20231115008: &admin.APIClient{
		// 					ProjectsApi: pAPI,
		// 				},
		// 			}, nil
		// 		},
		// 	},
		// 	wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		// 	wantFinalizers: nil,
		// 	wantConditions: []api.Condition{
		// 		api.FalseCondition(api.ReadyType).
		// 			WithReason(string(workflow.NetworkPeeringNotConfigured)).
		// 			WithMessageRegexp("failed to get project via Kubernetes reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
		// 		api.TrueCondition(api.ResourceVersionStatus),
		// 	},
		// },
		// {
		// 	title: "should handle network peering but fail to find container id from kube",
		// 	networkPeering: &akov2.AtlasNetworkPeering{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:              "network-peering",
		// 			Namespace:         "default",
		// 			Finalizers:        []string{customresource.FinalizerLabel},
		// 			DeletionTimestamp: &deletionTime,
		// 		},
		// 		Spec: akov2.AtlasNetworkPeeringSpec{
		// 			ProjectDualReference: akov2.ProjectDualReference{
		// 				ConnectionSecret: &api.LocalObjectReference{
		// 					Name: "my-secret",
		// 				},
		// 				ProjectRef: &common.ResourceRefNamespaced{
		// 					Name: "my-project",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	provider: &atlasmock.TestProvider{
		// 		IsSupportedFunc: func() bool {
		// 			return true
		// 		},
		// 		SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
		// 			pAPI := mockadmin.NewProjectsApi(t)
		// 			pAPI.EXPECT().GetProjectByName(mock.Anything, mock.Anything).Return(
		// 				admin.GetProjectByNameApiRequest{ApiService: pAPI},
		// 			)
		// 			pAPI.EXPECT().GetProjectByNameExecute(mock.AnythingOfType("admin.GetProjectByNameApiRequest")).Return(
		// 				&admin.Group{
		// 					Id: pointer.MakePtr(testProjectID),
		// 				}, nil, nil,
		// 			)
		// 			npAPI := mockadmin.NewNetworkPeeringApi(t)
		// 			return &atlas.ClientSet{
		// 				SdkClient20231115008: &admin.APIClient{
		// 					ProjectsApi:       pAPI,
		// 					NetworkPeeringApi: npAPI,
		// 				},
		// 			}, nil
		// 		},
		// 	},
		// 	wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
		// 	wantFinalizers: []string{customresource.FinalizerLabel},
		// 	wantConditions: []api.Condition{
		// 		api.FalseCondition(api.ReadyType).WithReason(string(workflow.Internal)).WithMessageRegexp(
		// 			"failed to solve Network Container id from Kubernetes: failed to fetch the Kubernetes Network Container  info: atlasnetworkcontainers.atlas.mongodb.com \"\" not found",
		// 		),
		// 		api.TrueCondition(api.ResourceVersionStatus),
		// 	},
		// },
	}
	for _, tc := range tests {
		t.Run(tc.title, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			require.NoError(t, corev1.AddToScheme(testScheme))
			require.NoError(t, akov2next.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project,
					tc.integration,
					&corev1.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-secret",
							Namespace: "default",
						},
						Data: map[string][]byte{
							"orgId":         []byte("orgId"),
							"publicApiKey":  []byte("publicApiKey"),
							"privateApiKey": []byte("privateApiKey"),
						},
					},
				).
				WithStatusSubresource(tc.integration).
				Build()
			logger := zaptest.NewLogger(t)
			ctx := context.Background()
			r := testReconciler(k8sClient, tc.provider, logger)
			result, err := r.handleCustomResource(ctx, tc.integration)
			tpi := getIntegration(t, ctx, k8sClient, client.ObjectKeyFromObject(tc.integration))
			require.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			assert.Equal(t, tc.wantFinalizers, getFinalizers(tpi))
			assert.Equal(t, cleanConditions(tc.wantConditions), cleanConditions(getConditions(tpi)))
		})
	}
}

func getIntegration(t *testing.T, ctx context.Context, k8sClient client.Client, key client.ObjectKey) *akov2next.AtlasThirdPartyIntegration {
	integration := &akov2next.AtlasThirdPartyIntegration{}
	if err := k8sClient.Get(ctx, key, integration); err != nil && !k8serrors.IsNotFound(err) {
		require.NoError(t, err)
	}
	return integration
}

func getFinalizers(integration *akov2next.AtlasThirdPartyIntegration) []string {
	if integration == nil {
		return nil
	}
	return integration.GetFinalizers()
}

func getConditions(integration *akov2next.AtlasThirdPartyIntegration) []api.Condition {
	if integration == nil {
		return nil
	}
	return integration.Status.GetConditions()
}

func testReconciler(k8sClient client.Client, provider atlas.Provider, logger *zap.Logger) *AtlasThirdPartyIntegrationsReconciler {
	return &AtlasThirdPartyIntegrationsReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: k8sClient,
			Log:    logger.Sugar(),
			GlobalSecretRef: client.ObjectKey{
				Namespace: "default",
				Name:      "secret",
			},
			AtlasProvider: provider,
		},
		EventRecorder: record.NewFakeRecorder(10),
	}
}

func cleanConditions(inputs []api.Condition) []api.Condition {
	outputs := make([]api.Condition, 0, len(inputs))
	for _, condition := range inputs {
		clean := condition
		clean.LastTransitionTime = metav1.Time{}
		outputs = append(outputs, clean)
	}
	return outputs
}
