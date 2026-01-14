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

package atlasnetworkpeering

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312012/admin"
	"go.mongodb.org/atlas-sdk/v20250312012/mockadmin"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akomock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

var (
	// sample error test
	ErrTestFail = errors.New("failure")
)

const (
	testProjectID = "project-id"

	testContainerName = "fake-container-name"

	testContainerID = "fake-container-id"

	testPeeringID = "peering-id"

	// testVpcID = "vpc-id"
)

func TestHandleCustomResource(t *testing.T) {
	deletionTime := metav1.Now()
	tests := []struct {
		title          string
		networkPeering *akov2.AtlasNetworkPeering
		provider       atlas.Provider
		wantResult     ctrl.Result
		wantError      bool
		wantFinalizers []string
		wantConditions []api.Condition
	}{
		{
			title: "should skip reconciliation",
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-peering",
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
			title:     "should fail to validate resource",
			wantError: true,
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-peering",
					Namespace: "default",
					Labels: map[string]string{
						customresource.ResourceVersion: "wrong",
					},
				},
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
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-peering",
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
					WithMessageRegexp("the AtlasNetworkPeering is not supported by Atlas for government"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title:     "should fail to resolve credentials and remove finalizer",
			wantError: true,
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-peering",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkPeeringSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
			wantFinalizers: nil,
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkPeeringNotConfigured)).
					WithMessageRegexp("error resolving project reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title:     "should fail to create sdk",
			wantError: true,
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "network-peering",
					Namespace: "default",
				},
				Spec: akov2.AtlasNetworkPeeringSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
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
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkPeeringNotConfigured)).
					WithMessageRegexp("failed to create sdk"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title:     "should fail to resolve project and remove finalizers",
			wantError: true,
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "network-peering",
					Namespace:  "default",
					Finalizers: []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasNetworkPeeringSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-no-existing-project",
						},
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
						SdkClient20250312011: &admin.APIClient{
							ProjectsApi: pAPI,
						},
					}, nil
				},
			},
			wantFinalizers: nil,
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.NetworkPeeringNotConfigured)).
					WithMessageRegexp("failed to get project via Kubernetes reference: missing Kubernetes Atlas Project\natlasprojects.atlas.mongodb.com \"my-no-existing-project\" not found"),
				api.TrueCondition(api.ResourceVersionStatus),
			},
		},
		{
			title:     "should handle network peering but fail to find container id from kube",
			wantError: true,
			networkPeering: &akov2.AtlasNetworkPeering{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "network-peering",
					Namespace:         "default",
					Finalizers:        []string{customresource.FinalizerLabel},
					DeletionTimestamp: &deletionTime,
				},
				Spec: akov2.AtlasNetworkPeeringSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-secret",
						},
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
				},
			},
			provider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientSetFunc: func(ctx context.Context, creds *atlas.Credentials, log *zap.SugaredLogger) (*atlas.ClientSet, error) {
					pAPI := mockadmin.NewProjectsApi(t)
					pAPI.EXPECT().GetGroupByName(mock.Anything, mock.Anything).Return(
						admin.GetGroupByNameApiRequest{ApiService: pAPI},
					)
					pAPI.EXPECT().GetGroupByNameExecute(mock.AnythingOfType("admin.GetGroupByNameApiRequest")).Return(
						&admin.Group{
							Id: pointer.MakePtr(testProjectID),
						}, nil, nil,
					)
					npAPI := mockadmin.NewNetworkPeeringApi(t)
					return &atlas.ClientSet{
						SdkClient20250312011: &admin.APIClient{
							ProjectsApi:       pAPI,
							NetworkPeeringApi: npAPI,
						},
					}, nil
				},
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.Internal)).WithMessageRegexp(
					"failed to solve Network Container id from Kubernetes: failed to fetch the Kubernetes Network Container  info: atlasnetworkcontainers.atlas.mongodb.com \"\" not found",
				),
				api.TrueCondition(api.ResourceVersionStatus),
			},
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
				WithObjects(project,
					tc.networkPeering,
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
				WithStatusSubresource(tc.networkPeering).
				Build()
			logger := zaptest.NewLogger(t)
			ctx := context.Background()
			r := testReconciler(k8sClient, tc.provider, logger)
			result, err := r.handleCustomResource(ctx, tc.networkPeering)
			np := getNetworkPeering(t, ctx, k8sClient, client.ObjectKeyFromObject(tc.networkPeering))
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantResult, result)
			assert.Equal(t, tc.wantFinalizers, getFinalizers(np))
			assert.Equal(t, cleanConditions(tc.wantConditions), cleanConditions(getConditions(np)))
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
			title: "create succeeds and goes in progress",
			req: &reconcileRequest{
				projectID:      testProjectID,
				networkPeering: testNetworkPeering(),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Create(mock.Anything, testProjectID, testContainerID, mock.Anything).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "10.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "CREATING",
						},
						nil,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.FalseCondition(api.NetworkPeerReadyType).
					WithMessageRegexp(fmt.Sprintf("Network Peering Connection %s is CREATING", testPeeringID)),
				api.FalseCondition(api.ReadyType),
			},
		},

		{
			title: "create fails",
			req: &reconcileRequest{
				projectID:      testProjectID,
				networkPeering: testNetworkPeering(),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Create(mock.Anything, testProjectID, testContainerID, mock.Anything).Return(
						nil, ErrTestFail,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantFinalizers: nil,
			wantErr:        fmt.Errorf("failed to create peering connection: %w", ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.NetworkPeeringNotConfigured)).
					WithMessageRegexp(fmt.Sprintf("failed to create peering connection: %v", ErrTestFail)),
			},
		},

		{
			title: "peering in sync",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(testNetworkPeering(), status.AtlasNetworkPeeringStatus{
					ID: testPeeringID,
				}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "10.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.TrueCondition(api.NetworkPeerReadyType),
				api.TrueCondition(api.ReadyType),
			},
		},

		{
			title: "peering connecting",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "10.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "NOT YET AVAILABLE",
						},
						nil,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.FalseCondition(api.NetworkPeerReadyType).WithMessageRegexp(
					"Network Peering Connection peering-id is NOT YET AVAILABLE",
				),
				api.FalseCondition(api.ReadyType),
			},
		},

		{
			title: "peering creation failed",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "10.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID:  testContainerID,
							Status:       "OOPs!",
							ErrorMessage: ErrTestFail.Error(),
						},
						nil,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr:        fmt.Errorf("peering connection failed: %s", ErrTestFail.Error()),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp(
						fmt.Sprintf("peering connection failed: %s", ErrTestFail.Error()),
					),
			},
		},

		{
			title: "update succeeds",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					)
					nps.EXPECT().Update(mock.Anything, testProjectID, testPeeringID, testContainerID, mock.Anything).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "10.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "UPDATING",
						},
						nil,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.FalseCondition(api.NetworkPeerReadyType).WithMessageRegexp(
					"Network Peering Connection peering-id is UPDATING",
				),
				api.FalseCondition(api.ReadyType),
			},
		},

		{
			title: "update fails",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					)
					nps.EXPECT().Update(mock.Anything, testProjectID, testPeeringID, testContainerID, mock.Anything).Return(
						nil, ErrTestFail,
					)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr:        fmt.Errorf("failed to update peering connection: %w", ErrTestFail),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.Internal)).
					WithMessageRegexp(fmt.Sprintf("failed to update peering connection: %v", ErrTestFail)),
			},
		},

		{
			title: "delete succeeds",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					withDeletionTimestamp(
						WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
						&deletionTime,
					),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					).Once()
					nps.EXPECT().Delete(mock.Anything, testProjectID, testPeeringID).Return(nil)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "DELETING",
						},
						nil,
					).Once()
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{RequeueAfter: workflow.DefaultRetry},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantConditions: []api.Condition{
				api.FalseCondition(api.NetworkPeerReadyType).WithMessageRegexp(
					"Network Peering Connection peering-id is DELETING",
				),
				api.FalseCondition(api.ReadyType),
			},
		},

		{
			title: "delete fails",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					withDeletionTimestamp(
						WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
						&deletionTime,
					),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					)
					nps.EXPECT().Delete(mock.Anything, testProjectID, testPeeringID).Return(ErrTestFail)
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr: fmt.Errorf("failed to delete peer connection %s: %s",
				testPeeringID, ErrTestFail.Error()),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.Internal)).
					WithMessageRegexp(fmt.Sprintf("failed to delete peer connection %s: %s",
						testPeeringID, ErrTestFail.Error())),
			},
		},

		{
			title: "delete fails getting closing peering",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					withDeletionTimestamp(
						WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
						&deletionTime,
					),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					).Once()
					nps.EXPECT().Delete(mock.Anything, testProjectID, testPeeringID).Return(nil)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(nil, ErrTestFail).Once()
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantFinalizers: []string{customresource.FinalizerLabel},
			wantErr: fmt.Errorf("failed to get closing peer connection %s: %s",
				testPeeringID, ErrTestFail.Error()),
			wantConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).WithReason(string(workflow.Internal)).
					WithMessageRegexp(fmt.Sprintf("failed to get closing peer connection %s: %s",
						testPeeringID, ErrTestFail.Error())),
			},
		},

		{
			title: "delete immediately success with not found",
			req: &reconcileRequest{
				projectID: testProjectID,
				networkPeering: withStatus(
					withDeletionTimestamp(
						WithFinalizers(testNetworkPeering(), []string{customresource.FinalizerLabel}),
						&deletionTime,
					),
					status.AtlasNetworkPeeringStatus{
						ID: testPeeringID,
					}),
				service: func() networkpeering.NetworkPeeringService {
					nps := akomock.NewNetworkPeeringServiceMock(t)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						&networkpeering.NetworkPeer{
							AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
								ID:       testPeeringID,
								Provider: "AWS",
								AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
									AccepterRegionName:  "us-east-1",
									AWSAccountID:        "some-aws-id",
									RouteTableCIDRBlock: "11.0.0.0/8",
									VpcID:               "vpc-id-test",
								},
							},
							ContainerID: testContainerID,
							Status:      "AVAILABLE",
						},
						nil,
					).Once()
					nps.EXPECT().Delete(mock.Anything, testProjectID, testPeeringID).Return(nil)
					nps.EXPECT().Get(mock.Anything, testProjectID, testPeeringID).Return(
						nil, networkpeering.ErrNotFound,
					).Once()
					return nps
				}(),
				containerService: func() networkcontainer.NetworkContainerService {
					ncs := akomock.NewNetworkContainerServiceMock(t)
					ncs.EXPECT().Get(mock.Anything, testProjectID, testContainerID).Return(
						testAtlasContainer(), nil,
					)
					return ncs
				}(),
			},
			wantResult:     ctrl.Result{},
			wantConditions: []api.Condition{},
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
				WithObjects(tc.req.networkPeering, testContainer()).
				Build()
			r := testReconciler(k8sClient, emptyProvider, logger)
			result, err := r.handle(workflowCtx, tc.req)
			if tc.wantErr != nil {
				assert.Equal(t, err.Error(), tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.wantResult, result)
			nc := getNetworkPeering(t, workflowCtx.Context, k8sClient, client.ObjectKeyFromObject(tc.req.networkPeering))
			assert.Equal(t, tc.wantFinalizers, getFinalizers(nc))
			assert.Equal(t, cleanConditions(tc.wantConditions), cleanConditions(workflowCtx.Conditions()))
		})
	}
}

func getNetworkPeering(t *testing.T, ctx context.Context, k8sClient client.Client, key client.ObjectKey) *akov2.AtlasNetworkPeering {
	networkPeering := &akov2.AtlasNetworkPeering{}
	if err := k8sClient.Get(ctx, key, networkPeering); err != nil && !k8serrors.IsNotFound(err) {
		require.NoError(t, err)
	}
	return networkPeering
}

func getFinalizers(networkContainer *akov2.AtlasNetworkPeering) []string {
	if networkContainer == nil {
		return nil
	}
	return networkContainer.GetFinalizers()
}

func getConditions(networkContainer *akov2.AtlasNetworkPeering) []api.Condition {
	if networkContainer == nil {
		return nil
	}
	return networkContainer.Status.GetConditions()
}

func testReconciler(k8sClient client.Client, provider atlas.Provider, logger *zap.Logger) *AtlasNetworkPeeringReconciler {
	return &AtlasNetworkPeeringReconciler{
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

func testContainer() *akov2.AtlasNetworkContainer {
	return &akov2.AtlasNetworkContainer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testContainerName,
			Namespace: "default",
			Annotations: map[string]string{
				customresource.ReconciliationPolicyAnnotation: customresource.ReconciliationPolicySkip,
			},
		},
		Spec: akov2.AtlasNetworkContainerSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: testProjectID,
				},
				ConnectionSecret: &api.LocalObjectReference{},
			},
			Provider: "AWS",
			AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
				Region:    "US_EAST_1",
				CIDRBlock: "10.0.0.0/18",
			},
		},
		Status: status.AtlasNetworkContainerStatus{
			ID:          testContainerID,
			Provisioned: true,
		},
	}
}

func testAtlasContainer() *networkcontainer.NetworkContainer {
	return &networkcontainer.NetworkContainer{
		NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
			Provider: "AWS",
			AtlasNetworkContainerConfig: akov2.AtlasNetworkContainerConfig{
				Region:    "US_EAST_1",
				CIDRBlock: "10.0.0.0/18",
			},
		},
		ID: testContainerID,
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

func testNetworkPeering() *akov2.AtlasNetworkPeering {
	return &akov2.AtlasNetworkPeering{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-peering",
			Namespace: "default",
		},
		Spec: akov2.AtlasNetworkPeeringSpec{
			ProjectDualReference: akov2.ProjectDualReference{
				ExternalProjectRef: &akov2.ExternalProjectReference{
					ID: testProjectID,
				},
				ConnectionSecret: &api.LocalObjectReference{
					Name: "fake-secret",
				},
			},
			ContainerRef: akov2.ContainerDualReference{
				Name: testContainerName,
			},
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				Provider: "AWS",
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  "us-east-1",
					AWSAccountID:        "some-aws-id",
					RouteTableCIDRBlock: "10.0.0.0/8",
					VpcID:               "vpc-id-test",
				},
			},
		},
	}
}

func withStatus(networkPeering *akov2.AtlasNetworkPeering, status status.AtlasNetworkPeeringStatus) *akov2.AtlasNetworkPeering {
	networkPeering.Status = status
	return networkPeering
}

func WithFinalizers(networkPeering *akov2.AtlasNetworkPeering, finalizers []string) *akov2.AtlasNetworkPeering {
	networkPeering.Finalizers = finalizers
	return networkPeering
}

func withDeletionTimestamp(networkPeering *akov2.AtlasNetworkPeering, deletionTimestamp *metav1.Time) *akov2.AtlasNetworkPeering {
	networkPeering.DeletionTimestamp = deletionTimestamp
	return networkPeering
}
