package atlasdatafederation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestReconcilePrivateEndpoints(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	tests := map[string]struct {
		privateEndpoints []akov2.DataFederationPE
		service          func() datafederation.DataFederationPrivateEndpointService
		expectedErr      error
	}{
		"no changes in private endpoint": {
			privateEndpoints: []akov2.DataFederationPE{
				{
					Provider:   "AWS",
					Type:       "DATA_LAKE",
					EndpointID: "vpcpe-123456",
				},
				{
					Provider:   "AZURE",
					Type:       "DATA_LAKE",
					EndpointID: "azure/resource/id",
				},
			},
			service: func() datafederation.DataFederationPrivateEndpointService {
				serviceMock := translation.NewDataFederationPrivateEndpointServiceMock(t)
				serviceMock.EXPECT().List(ctx, projectID).
					Return(
						[]datafederation.PrivateEndpoint{
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AZURE",
									Type:       "DATA_LAKE",
									EndpointID: "azure/resource/id",
								},
								ProjectID: projectID,
							},
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AWS",
									Type:       "DATA_LAKE",
									EndpointID: "vpcpe-123456",
								},
								ProjectID: projectID,
							},
						},
						nil,
					)

				return serviceMock
			},
			expectedErr: nil,
		},
		"failed to create a private endpoint": {
			privateEndpoints: []akov2.DataFederationPE{
				{
					Provider:   "AWS",
					Type:       "DATA_LAKE",
					EndpointID: "vpcpe-123456",
				},
				{
					Provider:   "AZURE",
					Type:       "DATA_LAKE",
					EndpointID: "azure/resource/id",
				},
			},
			service: func() datafederation.DataFederationPrivateEndpointService {
				serviceMock := translation.NewDataFederationPrivateEndpointServiceMock(t)
				serviceMock.EXPECT().List(ctx, projectID).
					Return(
						[]datafederation.PrivateEndpoint{
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AZURE",
									Type:       "DATA_LAKE",
									EndpointID: "azure/resource/id",
								},
								ProjectID: projectID,
							},
						},
						nil,
					)
				serviceMock.EXPECT().Create(ctx, mock.AnythingOfType("*datafederation.PrivateEndpoint")).
					Return(errors.New("failed to create private endpoint"))

				return serviceMock
			},
			expectedErr: errors.New("failed to create private endpoint"),
		},
		"create a private endpoint": {
			privateEndpoints: []akov2.DataFederationPE{
				{
					Provider:   "AWS",
					Type:       "DATA_LAKE",
					EndpointID: "vpcpe-123456",
				},
				{
					Provider:   "AZURE",
					Type:       "DATA_LAKE",
					EndpointID: "azure/resource/id",
				},
			},
			service: func() datafederation.DataFederationPrivateEndpointService {
				serviceMock := translation.NewDataFederationPrivateEndpointServiceMock(t)
				serviceMock.EXPECT().List(ctx, projectID).
					Return(
						[]datafederation.PrivateEndpoint{
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AZURE",
									Type:       "DATA_LAKE",
									EndpointID: "azure/resource/id",
								},
								ProjectID: projectID,
							},
						},
						nil,
					)
				serviceMock.EXPECT().Create(ctx, mock.AnythingOfType("*datafederation.PrivateEndpoint")).
					Return(nil)

				return serviceMock
			},
			expectedErr: nil,
		},
		"failed to delete a private endpoint": {
			privateEndpoints: []akov2.DataFederationPE{
				{
					Provider:   "AWS",
					Type:       "DATA_LAKE",
					EndpointID: "vpcpe-123456",
				},
			},
			service: func() datafederation.DataFederationPrivateEndpointService {
				serviceMock := translation.NewDataFederationPrivateEndpointServiceMock(t)
				serviceMock.EXPECT().List(ctx, projectID).
					Return(
						[]datafederation.PrivateEndpoint{
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AZURE",
									Type:       "DATA_LAKE",
									EndpointID: "azure/resource/id",
								},
								ProjectID: projectID,
							},
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AWS",
									Type:       "DATA_LAKE",
									EndpointID: "vpcpe-123456",
								},
								ProjectID: projectID,
							},
						},
						nil,
					)
				serviceMock.EXPECT().Delete(ctx, mock.AnythingOfType("*datafederation.PrivateEndpoint")).
					Return(errors.New("failed to delete private endpoint"))

				return serviceMock
			},
			expectedErr: errors.New("failed to delete private endpoint"),
		},
		"delete a private endpoint": {
			privateEndpoints: []akov2.DataFederationPE{
				{
					Provider:   "AWS",
					Type:       "DATA_LAKE",
					EndpointID: "vpcpe-123456",
				},
			},
			service: func() datafederation.DataFederationPrivateEndpointService {
				serviceMock := translation.NewDataFederationPrivateEndpointServiceMock(t)
				serviceMock.EXPECT().List(ctx, projectID).
					Return(
						[]datafederation.PrivateEndpoint{
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AZURE",
									Type:       "DATA_LAKE",
									EndpointID: "azure/resource/id",
								},
								ProjectID: projectID,
							},
							{
								DataFederationPE: &akov2.DataFederationPE{
									Provider:   "AWS",
									Type:       "DATA_LAKE",
									EndpointID: "vpcpe-123456",
								},
								ProjectID: projectID,
							},
						},
						nil,
					)
				serviceMock.EXPECT().Delete(ctx, mock.AnythingOfType("*datafederation.PrivateEndpoint")).
					Return(nil)

				return serviceMock
			},
			expectedErr: nil,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &AtlasDataFederationReconciler{
				privateEndpointService: tt.service(),
			}
			workflowCtx := &workflow.Context{
				Context: ctx,
				Log:     zaptest.NewLogger(t).Sugar(),
			}

			err := r.reconcilePrivateEndpoints(workflowCtx, projectID, tt.privateEndpoints)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestEnsurePrivateEndpoints(t *testing.T) {
	ctx := context.Background()
	projectID := "project-id"
	tests := map[string]struct {
		dataFederation     *akov2.AtlasDataFederation
		sdk                func() *admin.APIClient
		expectedConditions []api.Condition
		expectedResult     workflow.Result
	}{
		"failed to reconcile private endpoints": {
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
					},
				},
			},
			sdk: func() *admin.APIClient {
				dataFederationAPI := mockadmin.NewDataFederationApi(t)
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpoints(ctx, projectID).
					Return(admin.ListDataFederationPrivateEndpointsApiRequest{ApiService: dataFederationAPI})
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpointsExecute(mock.AnythingOfType("admin.ListDataFederationPrivateEndpointsApiRequest")).
					Return(nil, nil, errors.New("request failed"))

				return &admin.APIClient{DataFederationApi: dataFederationAPI}
			},
			expectedConditions: []api.Condition{
				api.FalseCondition(api.DataFederationPEReadyType).
					WithMessageRegexp("failed to list data federation private endpoints from Atlas: request failed"),
			},
			expectedResult: workflow.Terminate(workflow.Internal, "failed to list data federation private endpoints from Atlas: request failed"),
		},
		"no private endpoints to reconcile": {
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{},
				},
			},
			sdk: func() *admin.APIClient {
				dataFederationAPI := mockadmin.NewDataFederationApi(t)
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpoints(ctx, projectID).
					Return(admin.ListDataFederationPrivateEndpointsApiRequest{ApiService: dataFederationAPI})
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpointsExecute(mock.AnythingOfType("admin.ListDataFederationPrivateEndpointsApiRequest")).
					Return(nil, nil, nil)

				return &admin.APIClient{DataFederationApi: dataFederationAPI}
			},
			expectedConditions: []api.Condition{},
			expectedResult:     workflow.OK(),
		},
		"private endpoints reconcile": {
			dataFederation: &akov2.AtlasDataFederation{
				Spec: akov2.DataFederationSpec{
					PrivateEndpoints: []akov2.DataFederationPE{
						{
							Provider:   "AWS",
							Type:       "DATA_LAKE",
							EndpointID: "vpcpe-123456",
						},
					},
				},
			},
			sdk: func() *admin.APIClient {
				dataFederationAPI := mockadmin.NewDataFederationApi(t)
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpoints(ctx, projectID).
					Return(admin.ListDataFederationPrivateEndpointsApiRequest{ApiService: dataFederationAPI})
				dataFederationAPI.EXPECT().ListDataFederationPrivateEndpointsExecute(mock.AnythingOfType("admin.ListDataFederationPrivateEndpointsApiRequest")).
					Return(
						&admin.PaginatedPrivateNetworkEndpointIdEntry{
							Results: &[]admin.PrivateNetworkEndpointIdEntry{
								{
									Provider:   pointer.MakePtr("AWS"),
									Type:       pointer.MakePtr("DATA_LAKE"),
									EndpointId: "vpcpe-123456",
								},
							},
							TotalCount: pointer.MakePtr(1),
						},
						nil,
						nil,
					)

				return &admin.APIClient{DataFederationApi: dataFederationAPI}
			},
			expectedConditions: []api.Condition{
				api.TrueCondition(api.DataFederationPEReadyType),
			},
			expectedResult: workflow.OK(),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := &AtlasDataFederationReconciler{}
			workflowCtx := &workflow.Context{
				Context:   ctx,
				SdkClient: tt.sdk(),
				Log:       zaptest.NewLogger(t).Sugar(),
			}

			result := r.ensurePrivateEndpoints(workflowCtx, projectID, tt.dataFederation)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, workflowCtx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
