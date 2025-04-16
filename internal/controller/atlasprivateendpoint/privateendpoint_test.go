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

package atlasprivateendpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

func TestHandlePrivateEndpointService(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	deletionTime := metav1.Now()

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		peClient             func() privateendpoint.PrivateEndpointService
		expectedResult       reconcile.Result
		expectedConditions   []api.Condition
	}{
		"failed to retrieve private endpoint": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(nil, errors.New("failed to get private endpoint"))

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get private endpoint"),
			},
		},
		"failed to create private endpoint service": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().ListPrivateEndpoints(ctx, projectID, "AWS").
					Return(nil, nil)
				c.EXPECT().CreatePrivateEndpointService(ctx, projectID, mock.AnythingOfType("*privateendpoint.AWSService")).
					Return(nil, errors.New("failed to create private endpoint"))

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointServiceFailedToCreate)).
					WithMessageRegexp("failed to create private endpoint"),
			},
		},
		"unmanage already deleted private endpoint": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pe1",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(nil, nil)

				return c
			},
			expectedResult: reconcile.Result{},
		},
		"unmanage protected private endpoint": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pe1",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
					Annotations: map[string]string{
						customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
					},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(&privateendpoint.AWSService{}, nil)

				return c
			},
			expectedResult: reconcile.Result{},
		},
		"failed to delete private endpoint": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pe1",
					Namespace:         "default",
					DeletionTimestamp: &deletionTime,
					Finalizers:        []string{customresource.FinalizerLabel},
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "AVAILABLE",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)
				c.EXPECT().DeleteEndpointService(ctx, projectID, "AWS", "pe-service-id").
					Return(errors.New("failed to delete private endpoint"))

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointFailedToDelete)).
					WithMessageRegexp("failed to delete private endpoint"),
			},
		},
		"private endpoint service is initiating": { //nolint:dupl
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "INITIATING",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "INITIATING",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointServiceInitializing)).
					WithMessageRegexp("Private Endpoint is being initialized"),
			},
		},
		"private endpoint service is pending": { //nolint:dupl
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "PENDING",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "PENDING",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointServiceInitializing)).
					WithMessageRegexp("Private Endpoint is waiting for human action"),
			},
		},
		"private endpoint service was rejected": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "REJECTED",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "REJECTED",
						Error:         "atlas could not connect the private endpoint",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointServiceFailedToConfigure)).
					WithMessageRegexp("atlas could not connect the private endpoint"),
			},
		},
		"private endpoint service is being deleted": { //nolint:dupl
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "DELETING",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "DELETING",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointServiceReady).
					WithReason(string(workflow.PrivateEndpointServiceDeleting)).
					WithMessageRegexp("Private Endpoint is being deleted"),
			},
		},
		"private endpoint service is available": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					ServiceID:     "pe-service-id",
					ServiceStatus: "AVAILABLE",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "US_EAST_1",
						ServiceStatus: "AVAILABLE",
						Interfaces:    privateendpoint.EndpointInterfaces{},
					},
				}

				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(awsService, nil)

				return c
			},
			expectedResult: reconcile.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.PrivateEndpointServiceReady),
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointConfigurationPending)).
					WithMessageRegexp("waiting for private endpoint configuration from customer side"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasPrivateEndpoint).
				WithStatusSubresource(tt.atlasPrivateEndpoint).
				Build()

			logger := zaptest.NewLogger(t).Sugar()
			r := &AtlasPrivateEndpointReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    logger,
				},
				EventRecorder: record.NewFakeRecorder(10),
			}
			workflowCtx := workflow.Context{
				Context: ctx,
				Log:     logger,
			}

			result, err := r.handlePrivateEndpointService(&workflowCtx, tt.peClient(), projectID, tt.atlasPrivateEndpoint)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					workflowCtx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestHandlePrivateEndpointInterfaces(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		atlasPEService       func() privateendpoint.EndpointService
		peClient             func() privateendpoint.PrivateEndpointService
		expectedResult       reconcile.Result
		expectedConditions   []api.Condition
	}{
		"failed to create private endpoint interface": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
					AWSConfiguration: []akov2.AWSPrivateEndpointConfiguration{
						{
							ID: "vpcpe-123456",
						},
					},
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionFalse,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionFalse,
								Reason:             string(workflow.PrivateEndpointConfigurationPending),
								Message:            "waiting for private endpoint configuration from customer side",
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ServiceName:   "aws/service/name",
					ServiceStatus: "AVAILABLE",
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						Interfaces: privateendpoint.EndpointInterfaces{},
					},
				}

				return awsService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().
					CreatePrivateEndpointInterface(ctx, projectID, "AWS", "pe-service-id", "", mock.AnythingOfType("*privateendpoint.AWSInterface")).
					Return(nil, errors.New("failed to create private endpoint interface"))

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointFailedToCreate)).
					WithMessageRegexp("failed to create private endpoint interface"),
			},
		},
		"create private endpoint interface": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "GCP",
					Region:   "EUROPE_WEST_3",
					GCPConfiguration: []akov2.GCPPrivateEndpointConfiguration{
						{
							ProjectID: "customer-project-id1",
							GroupName: "group-name1",
							Endpoints: []akov2.GCPPrivateEndpoint{
								{
									Name: "group-name1-pe1",
									IP:   "10.0.0.1",
								},
								{
									Name: "group-name1-pe2",
									IP:   "10.0.0.2",
								},
								{
									Name: "group-name1-pe3",
									IP:   "10.0.0.3",
								},
							},
						},
					},
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionFalse,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionFalse,
								Reason:             string(workflow.PrivateEndpointConfigurationPending),
								Message:            "waiting for private endpoint configuration from customer side",
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID: "pe-service-id",
					ServiceAttachmentNames: []string{
						"atls/service/attachment/name/1",
						"atls/service/attachment/name/2",
						"atls/service/attachment/name/3",
					},
					ServiceStatus: "AVAILABLE",
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						Interfaces: privateendpoint.EndpointInterfaces{},
					},
				}

				return awsService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().
					CreatePrivateEndpointInterface(
						ctx,
						projectID,
						"GCP",
						"pe-service-id",
						"customer-project-id1",
						mock.AnythingOfType("*privateendpoint.GCPInterface"),
					).
					Return(nil, nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointUpdating)).
					WithMessageRegexp("Private Endpoints are being updated"),
			},
		},
		"failed to configure private endpoint interface": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
					AWSConfiguration: []akov2.AWSPrivateEndpointConfiguration{
						{
							ID: "vpcpe-123456",
						},
					},
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionFalse,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionFalse,
								Reason:             string(workflow.PrivateEndpointConfigurationPending),
								Message:            "waiting for private endpoint configuration from customer side",
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ServiceName:   "aws/service/name",
					ServiceStatus: "AVAILABLE",
					Endpoints: []status.EndpointInterfaceStatus{
						{
							ID:     "vpcpe-123456",
							Status: "INITIATING",
						},
					},
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				awsInterface := &privateendpoint.AWSInterface{
					CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
						ID:              "vpcpe-123456",
						InterfaceStatus: "REJECTED",
						Error:           "failed to configure private endpoint interface",
					},
				}
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						Interfaces: privateendpoint.EndpointInterfaces{awsInterface},
					},
				}

				return awsService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointFailedToConfigure)).
					WithMessageRegexp("failed to configure private endpoint interface"),
			},
		},
		"failed to delete private endpoint interface": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionFalse,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionFalse,
								Reason:             string(workflow.PrivateEndpointConfigurationPending),
								Message:            "waiting for private endpoint configuration from customer side",
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ServiceName:   "aws/service/name",
					ServiceStatus: "AVAILABLE",
					Endpoints: []status.EndpointInterfaceStatus{
						{
							ID:     "vpcpe-123456",
							Status: "AVAILABLE",
						},
					},
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				awsInterface := &privateendpoint.AWSInterface{
					CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
						ID:              "vpcpe-123456",
						InterfaceStatus: "AVAILABLE",
					},
				}
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						Interfaces: privateendpoint.EndpointInterfaces{awsInterface},
					},
				}

				return awsService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().
					DeleteEndpointInterface(ctx, projectID, "AWS", "pe-service-id", "vpcpe-123456").
					Return(errors.New("failed to delete private endpoint interface"))

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointFailedToDelete)).
					WithMessageRegexp("failed to delete private endpoint interface"),
			},
		},
		"delete private endpoint interface": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ServiceName:   "aws/service/name",
					ServiceStatus: "AVAILABLE",
					Endpoints: []status.EndpointInterfaceStatus{
						{
							ID:     "vpcpe-123456",
							Status: "AVAILABLE",
						},
					},
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				awsInterface := &privateendpoint.AWSInterface{
					CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
						ID:              "vpcpe-123456",
						InterfaceStatus: "AVAILABLE",
					},
				}
				awsService := &privateendpoint.AWSService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						Interfaces: privateendpoint.EndpointInterfaces{awsInterface},
					},
				}

				return awsService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().
					DeleteEndpointInterface(ctx, projectID, "AWS", "pe-service-id", "vpcpe-123456").
					Return(nil)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointUpdating)).
					WithMessageRegexp("Private Endpoints are being updated"),
			},
		},
		"private endpoints are in progress": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{ //nolint:dupl
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AZURE",
					Region:   "GERMANY_NORTH",
					AzureConfiguration: []akov2.AzurePrivateEndpointConfiguration{
						{
							ID: "azure/resource/id",
							IP: "10.0.0.2",
						},
					},
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ResourceID:    "atlas/azure/resource/id",
					ServiceStatus: "AVAILABLE",
					Endpoints: []status.EndpointInterfaceStatus{
						{
							ID:             "azure/resource/id",
							ConnectionName: "atlas-connection-name",
							Status:         "INITIATING",
						},
					},
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				azureInterface := &privateendpoint.AzureInterface{
					CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
						ID:              "azure/resource/id",
						InterfaceStatus: "INITIATING",
					},
					IP:             "10.0.0.2",
					ConnectionName: "atlas-connection-name",
				}
				azureService := &privateendpoint.AzureService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "GERMANY_NORTH",
						ServiceStatus: "AVAILABLE",
						Interfaces:    privateendpoint.EndpointInterfaces{azureInterface},
					},
					ServiceName: "atlas/azure/service/name",
					ResourceID:  "atlas/azure/resource/id",
				}

				return azureService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)

				return c
			},
			expectedResult: reconcile.Result{RequeueAfter: workflow.DefaultRetry},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.ReadyType),
				api.FalseCondition(api.PrivateEndpointReady).
					WithReason(string(workflow.PrivateEndpointUpdating)).
					WithMessageRegexp("Private Endpoints are being updated"),
			},
		},
		"private endpoints are ready": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{ //nolint:dupl
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AZURE",
					Region:   "GERMANY_NORTH",
					AzureConfiguration: []akov2.AzurePrivateEndpointConfiguration{
						{
							ID: "azure/resource/id",
							IP: "10.0.0.2",
						},
					},
				},
				Status: status.AtlasPrivateEndpointStatus{
					Common: api.Common{
						Conditions: []api.Condition{
							{
								Type:               api.ReadyType,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointServiceReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
							{
								Type:               api.PrivateEndpointReady,
								Status:             corev1.ConditionTrue,
								LastTransitionTime: metav1.Now(),
							},
						},
					},
					ServiceID:     "pe-service-id",
					ResourceID:    "atlas/azure/resource/id",
					ServiceStatus: "AVAILABLE",
					Endpoints: []status.EndpointInterfaceStatus{
						{
							ID:             "azure/resource/id",
							ConnectionName: "atlas-connection-name",
							Status:         "AVAILABLE",
						},
					},
				},
			},
			atlasPEService: func() privateendpoint.EndpointService {
				azureInterface := &privateendpoint.AzureInterface{
					CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
						ID:              "azure/resource/id",
						InterfaceStatus: "AVAILABLE",
					},
					IP:             "10.0.0.2",
					ConnectionName: "atlas-connection-name",
				}
				azureService := &privateendpoint.AzureService{
					CommonEndpointService: privateendpoint.CommonEndpointService{
						ID:            "pe-service-id",
						CloudRegion:   "GERMANY_NORTH",
						ServiceStatus: "AVAILABLE",
						Interfaces:    privateendpoint.EndpointInterfaces{azureInterface},
					},
					ServiceName: "atlas/azure/service/name",
					ResourceID:  "atlas/azure/resource/id",
				}

				return azureService
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)

				return c
			},
			expectedResult: reconcile.Result{},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.PrivateEndpointServiceReady),
				api.TrueCondition(api.PrivateEndpointReady),
				api.TrueCondition(api.ReadyType),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasPrivateEndpoint).
				WithStatusSubresource(tt.atlasPrivateEndpoint).
				Build()

			logger := zaptest.NewLogger(t).Sugar()
			r := &AtlasPrivateEndpointReconciler{
				AtlasReconciler: reconciler.AtlasReconciler{
					Client: fakeClient,
					Log:    logger,
				},
				EventRecorder: record.NewFakeRecorder(10),
			}
			workflowCtx := workflow.Context{
				Context: ctx,
				Log:     logger,
			}

			akoPEService := privateendpoint.NewPrivateEndpoint(tt.atlasPrivateEndpoint)
			result, err := r.handlePrivateEndpointInterface(&workflowCtx, tt.peClient(), projectID, tt.atlasPrivateEndpoint, akoPEService, tt.atlasPEService())
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
			t.Log(cmp.Diff(
				tt.expectedConditions,
				workflowCtx.Conditions(),
				cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
			))
			assert.True(
				t,
				cmp.Equal(
					tt.expectedConditions,
					workflowCtx.Conditions(),
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}

func TestGetPrivateEndpointService(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		atlasPrivateEndpoint *akov2.AtlasPrivateEndpoint
		peClient             func() privateendpoint.PrivateEndpointService
		expectedResult       privateendpoint.EndpointService
		expectedErr          error
	}{
		"failed to list private endpoint to match": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().ListPrivateEndpoints(ctx, projectID, "AWS").
					Return(nil, errors.New("failed to list private endpoints"))

				return c
			},
			expectedErr: errors.New("failed to list private endpoints"),
		},
		"match private endpoint from the list": {
			atlasPrivateEndpoint: &akov2.AtlasPrivateEndpoint{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pe1",
					Namespace: "default",
				},
				Spec: akov2.AtlasPrivateEndpointSpec{
					ProjectDualReference: akov2.ProjectDualReference{
						ExternalProjectRef: &akov2.ExternalProjectReference{
							ID: projectID,
						},
						ConnectionSecret: &api.LocalObjectReference{},
					},
					Provider: "AWS",
					Region:   "US_EAST_1",
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().ListPrivateEndpoints(ctx, projectID, "AWS").
					Return(
						[]privateendpoint.EndpointService{
							&privateendpoint.AWSService{
								CommonEndpointService: privateendpoint.CommonEndpointService{
									ID:            "pe-service-id-1",
									CloudRegion:   "EU_CENTRAL_1",
									ServiceStatus: "AVAILABLE",
								},
							},
							&privateendpoint.AWSService{
								CommonEndpointService: privateendpoint.CommonEndpointService{
									ID:            "pe-service-id-2",
									CloudRegion:   "US_EAST_1",
									ServiceStatus: "AVAILABLE",
								},
							},
						},
						nil,
					)

				return c
			},
			expectedResult: &privateendpoint.AWSService{
				CommonEndpointService: privateendpoint.CommonEndpointService{
					ID:            "pe-service-id-2",
					CloudRegion:   "US_EAST_1",
					ServiceStatus: "AVAILABLE",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &AtlasPrivateEndpointReconciler{}

			result, err := r.getOrMatchPrivateEndpointService(ctx, tt.peClient(), projectID, tt.atlasPrivateEndpoint)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestDeletePrivateEndpoint(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"

	testScheme := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(testScheme))

	tests := map[string]struct {
		peService   privateendpoint.EndpointService
		peClient    func() privateendpoint.PrivateEndpointService
		expectedErr error
	}{
		"failed to delete private endpoint interface": {
			peService: &privateendpoint.AWSService{
				CommonEndpointService: privateendpoint.CommonEndpointService{
					ID:            "pe-service-id",
					CloudRegion:   "US_EAST_1",
					ServiceStatus: "AVAILABLE",
					Interfaces: privateendpoint.EndpointInterfaces{
						&privateendpoint.AWSInterface{
							CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
								ID:              "vpcpe-123456",
								InterfaceStatus: "AVAILABLE",
							},
						},
					},
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().DeleteEndpointInterface(ctx, projectID, "AWS", "pe-service-id", "vpcpe-123456").
					Return(errors.New("failed to delete private endpoint interface"))

				return c
			},
			expectedErr: errors.New("failed to delete private endpoint interface"),
		},
		"delete private endpoint interface": {
			peService: &privateendpoint.AWSService{
				CommonEndpointService: privateendpoint.CommonEndpointService{
					ID:            "pe-service-id",
					CloudRegion:   "US_EAST_1",
					ServiceStatus: "AVAILABLE",
					Interfaces: privateendpoint.EndpointInterfaces{
						&privateendpoint.AWSInterface{
							CommonEndpointInterface: privateendpoint.CommonEndpointInterface{
								ID:              "vpcpe-123456",
								InterfaceStatus: "AVAILABLE",
							},
						},
					},
				},
			},
			peClient: func() privateendpoint.PrivateEndpointService {
				c := translation.NewPrivateEndpointServiceMock(t)
				c.EXPECT().DeleteEndpointInterface(ctx, projectID, "AWS", "pe-service-id", "vpcpe-123456").
					Return(nil)
				c.EXPECT().GetPrivateEndpoint(ctx, projectID, "AWS", "pe-service-id").
					Return(&privateendpoint.AWSService{}, nil)

				return c
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &AtlasPrivateEndpointReconciler{}

			_, err := r.deletePrivateEndpoint(ctx, tt.peClient(), projectID, tt.peService)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
