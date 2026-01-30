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

package atlasnetworkcontainer

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akomock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

var (
	// sample error test
	ErrTestFail = errors.New("failure")
)

const (
	testContainerID = "container-id"
)

func TestHandleCustomResource(t *testing.T) {
	type workflowResult struct {
		result ctrl.Result
		err    error
	}
	deletionTime := metav1.Now()
	tests := []struct {
		title            string
		networkContainer *akov2.AtlasNetworkContainer
		provider         atlas.Provider
		wantResult       workflowResult
		wantFinalizers   []string
		wantConditions   []api.Condition
	}{
		{
			title: "should skip reconciliation",
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			wantResult: workflowResult{
				result: ctrl.Result{},
				err:    nil,
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
		},
		{
			title: "should fail to validate resource",
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			wantResult: workflowResult{
				err: errors.New("wrong is not a valid semver version for label mongodb.com/atlas-resource-version"),
			},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.ResourceVersionStatus).
					WithReason(string(workflow.AtlasResourceVersionIsInvalid)).
					WithMessageRegexp("wrong is not a valid semver version for label mongodb.com/atlas-resource-version"),
			},
		},
		{
			title: "should fail when not supported",
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-container",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkContainerSpec{
					Provider: "AWS",
					AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
						Region:    "US_EAST_1",
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return false
				},
			},
			wantResult: workflowResult{
				result: ctrl.Result{},
				err:    nil,
			},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.AtlasGovUnsupported)).
					WithMessageRegexp("the AtlasNetworkContainer is not supported by Atlas for government"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title: "should fail to resolve credentials and remove finalizer",
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "network-container",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
			wantResult: workflowResult{
				err: errors.New("error resolving project reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
			},
			wantFinalizers: nil,
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("error resolving project reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title: "should fail to create sdk",
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					return nil, errors.New("failed to create sdk")
				},
			},
			wantResult: workflowResult{
				err: errors.New("failed to create sdk"),
			},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("failed to create sdk"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title: "should fail to resolve project and remove finalizers",
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "network-container",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					pAPI := mockadmin.NewProjectsApi(t)
					return &atlas.ClientSet{
						SdkClient20250312013: &admin.APIClient{ProjectsApi: pAPI},
					}, nil
				},
			},
			wantResult: workflowResult{
				err: errors.New("failed to get project via Kubernetes reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
			},
			wantFinalizers: nil,
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp("failed to get project via Kubernetes reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title: "should handle network container with unmanage",
			networkContainer: &akov2.AtlasNetworkContainer{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "network-container",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: &deletionTime,
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
						CIDRBlock: "11.10.0.0/21",
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					ncAPI := mockadmin.NewNetworkPeeringApi(t)
					ncAPI.EXPECT().ListGroupContainers(mock.Anything, mock.Anything).Return(
						admin.ListGroupContainersApiRequest{ApiService: ncAPI},
					)
					ncAPI.EXPECT().ListGroupContainersExecute(mock.Anything).Return(
						&admin.PaginatedCloudProviderContainer{
							Results: &[]admin.CloudProviderContainer{},
						}, nil, nil,
					)
					pAPI := mockadmin.NewProjectsApi(t)
					pAPI.EXPECT().GetGroupByName(mock.Anything, mock.Anything).Return(
						admin.GetGroupByNameApiRequest{ApiService: pAPI},
					)
					pAPI.EXPECT().GetGroupByNameExecute(mock.Anything).Return(
						&admin.Group{
							Id: pointer.MakePtr(testProjectID),
						}, nil, nil,
					)
					return &atlas.ClientSet{
						SdkClient20250312013: &admin.APIClient{
							NetworkPeeringApi: ncAPI,
							ProjectsApi:       pAPI,
						},
					}, nil
				},
			},
			wantResult: workflowResult{
				result: ctrl.Result{},
				err:    nil,
			},
			wantFinalizers: nil,
			wantConditions: nil,
		},
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
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, tc.networkContainer, &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-secret",
						Namespace: "default",
					},
					Data: map[string][]byte{
						"orgId":         []byte("orgId"),
						"publicApiKey":  []byte("publicApiKey"),
						"privateApiKey": []byte("privateApiKey"),
					},
				}).
				WithStatusSubresource(tc.networkContainer).
				Build()
			logger := zaptest.NewLogger(t)
			ctx := context.Background()
			r := testReconciler(k8sClient, tc.provider, logger)
			result, err := r.handleCustomResource(ctx, tc.networkContainer)
			nc := getNetworkContainer(t, ctx, k8sClient, client.ObjectKeyFromObject(tc.networkContainer))
			if tc.wantResult.err != nil {
				assert.Equal(t, tc.wantResult.err.Error(), err.Error())
				//assert.ErrorIs(t, tc.wantResult.err, err, "expected '%v',\ngot '%v'", tc.wantResult.err, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantResult.result, result)
			assert.Equal(t, tc.wantFinalizers, getFinalizers(nc))
			assert.Equal(t, cleanConditions(tc.wantConditions), cleanConditions(getConditions(nc)))
		})
	}
}

func TestHandle(t *testing.T) {
	deletionTime := metav1.Now()
	emptyProvider := &atlasmock.TestProvider{}
	logger := zaptest.NewLogger(t)
	for _, tc := range []struct {
		title          string
		req            *reconcileRequest
		wantResult     ctrl.Result
		wantErr        error
		wantFinalizers []string
		wantConditions []api.Condition
	}{
		{
			title: "create succeeds",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-container",
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						nil, networkcontainer.ErrNotFound,
					)
					ncs.EXPECT().Create(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: false,
						},
						nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.TrueCondition(api.NetworkContainerReady).
					WithMessageRegexp(fmt.Sprintf("Network Container %s is ready", testContainerID)),
				api.TrueCondition(api.ReadyType),
			},
		},

		{
			title: "create fails",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-container",
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						nil, networkcontainer.ErrNotFound,
					)
					ncs.EXPECT().Create(mock.Anything, testProjectID, mock.Anything).Return(
						nil,
						ErrTestFail,
					)
					return ncs
				}(),
			},
			wantErr:        fmt.Errorf("failed to create container: %w", ErrTestFail),
			wantFinalizers: nil,
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to create container: %v", ErrTestFail)),
			},
		},

		{
			title: "in sync",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-container",
						Finalizers: []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: true,
						}, nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.TrueCondition(api.NetworkContainerReady).
					WithMessageRegexp(fmt.Sprintf("Network Container %s is ready", testContainerID)),
				api.TrueCondition(api.ReadyType),
			},
		},

		{
			title: "existent container in sync",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-container",
						Finalizers: []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        testContainerID,
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: true,
						}, nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.TrueCondition(api.NetworkContainerReady).
					WithMessageRegexp(fmt.Sprintf("Network Container %s is ready", testContainerID)),
				api.TrueCondition(api.ReadyType),
			},
		},

		{
			title: "update succeeds",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-container",
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.12.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: true,
						}, nil,
					)
					ncs.EXPECT().Update(mock.Anything, testProjectID, testContainerID, &networkcontainer.NetworkContainerConfig{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.12.0.0/21",
						},
					}).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.12.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: false,
						}, nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.TrueCondition(api.NetworkContainerReady).
					WithMessageRegexp(fmt.Sprintf("Network Container %s is ready", testContainerID)),
				api.TrueCondition(api.ReadyType),
			},
		},

		{
			title: "update fails",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-container",
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.12.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region:    "US_EAST_1",
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: true,
						}, nil,
					)
					ncs.EXPECT().Update(mock.Anything, testProjectID, testContainerID, mock.Anything).Return(
						nil, ErrTestFail,
					)
					return ncs
				}(),
			},
			wantFinalizers: nil,
			wantErr:        fmt.Errorf("failed to update container: %w", ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to update container: %v", ErrTestFail)),
			},
		},

		{
			title: "delete succeeds",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "test-container",
						Finalizers:        []string{customresource.FinalizerLabel},
						DeletionTimestamp: &deletionTime,
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.12.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "AWS",
								AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
									Region: "US_EAST_1",
									// different CIDR, but it should not matter as we are removing
									CIDRBlock: "10.11.0.0/21",
								},
							},
							ID:          testContainerID,
							Provisioned: true,
						}, nil,
					)
					ncs.EXPECT().Delete(mock.Anything, testProjectID, testContainerID).Return(
						nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: nil,
			wantConditions: []api.Condition{},
		},

		{
			title: "delete fails",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "test-another-container",
						Finalizers:        []string{customresource.FinalizerLabel},
						DeletionTimestamp: &deletionTime,
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "Azure",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_2",
							CIDRBlock: "10.14.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(
						&networkcontainer.NetworkContainer{
							NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
								Provider: "Azure", // almost empty, but we are removing anyways
							},
							ID: testContainerID,
						}, nil,
					)
					ncs.EXPECT().Delete(mock.Anything, testProjectID, testContainerID).Return(
						ErrTestFail,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr:        fmt.Errorf("failed to delete container: %w", ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotDeleted)).
					WithMessageRegexp(fmt.Sprintf("failed to delete container: %v", ErrTestFail)),
			},
		},

		{
			title: "discover find fails abnormally",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-container",
						Finalizers: []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Find(mock.Anything, testProjectID, mock.Anything).Return(nil, ErrTestFail)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr: fmt.Errorf("failed to find container from project %s: %w",
				testProjectID, ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to find container from project %s: %v",
						testProjectID, ErrTestFail)),
			},
		},

		{
			title: "discover get fails abnormally",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-container",
						Finalizers: []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        testContainerID,
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(nil, ErrTestFail)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr: fmt.Errorf("failed to get container %s from project %s: %w",
				testContainerID, testProjectID, ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to get container %s from project %s: %v",
						testContainerID, testProjectID, ErrTestFail)),
			},
		},

		{
			title: "discover get fails with not found",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkContainer: &akov2.AtlasNetworkContainer{
					ObjectMeta: metav1.ObjectMeta{
						Name:       "test-container",
						Finalizers: []string{customresource.FinalizerLabel},
					},
					Spec: akov2.AtlasNetworkContainerSpec{
						Provider: "AWS",
						AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
							ID:        testContainerID,
							Region:    "US_EAST_1",
							CIDRBlock: "10.11.0.0/21",
						},
					},
				},
				service: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						nil,
						networkcontainer.ErrNotFound,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr: fmt.Errorf("failed to get container %s from project %s: %w",
				testContainerID, testProjectID, networkcontainer.ErrNotFound),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkContainerNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to get container %s from project %s: %v",
						testContainerID, testProjectID, networkcontainer.ErrNotFound)),
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			workflowCtx := &workflow.Context{
				Context: context.Background(),
			}
			testScheme := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.req.networkContainer).
				Build()
			r := testReconciler(k8sClient, emptyProvider, logger)
			result, err := r.handle(workflowCtx, tc.req)
			if tc.wantErr != nil {
				assert.Equal(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			//assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantResult, result)
			nc := getNetworkContainer(t, workflowCtx.Context, k8sClient, client.ObjectKeyFromObject(tc.req.networkContainer))
			assert.Equal(t, tc.wantFinalizers, getFinalizers(nc))
			assert.Equal(t, cleanConditions(tc.wantConditions), cleanConditions(workflowCtx.Conditions()))
		})
	}
}

func getNetworkContainer(t *testing.T, ctx context.Context, k8sClient client.Client, key client.ObjectKey) *akov2.AtlasNetworkContainer {
	networkContainer := &akov2.AtlasNetworkContainer{}
	if err := k8sClient.Get(ctx, key, networkContainer); err != nil && !k8serrors.IsNotFound(err) {
		require.NoError(t, err)
	}
	return networkContainer
}

func getFinalizers(networkContainer *akov2.AtlasNetworkContainer) []string {
	if networkContainer == nil {
		return nil
	}
	return networkContainer.GetFinalizers()
}

func getConditions(networkContainer *akov2.AtlasNetworkContainer) []api.Condition {
	if networkContainer == nil {
		return nil
	}
	return networkContainer.Status.GetConditions()
}

func testReconciler(k8sClient client.Client, provider atlas.Provider, logger *zap.Logger) *AtlasNetworkContainerReconciler {
	return &AtlasNetworkContainerReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:        k8sClient,
			Log:           logger.Sugar(),
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
