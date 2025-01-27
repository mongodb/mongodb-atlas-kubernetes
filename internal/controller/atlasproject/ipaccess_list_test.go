package atlasproject

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func TestIpAccessListController_reconcile(t *testing.T) {
	tests := map[string]struct {
		service            func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService
		ipAccessList       []project.IPAccessList
		expectedResult     workflow.Result
		expectedConditions []api.Condition
	}{
		"should fail to convert wrongly defined ip access list": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{}, nil)

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "wrong-ip",
				},
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("ip wrong-ip is invalid")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("ip wrong-ip is invalid"),
			},
		},
		"should unmanage ip access list when unset on both Atlas and AKO": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{}, nil)

				return serviceMock
			},
			ipAccessList:       nil,
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
		"should fail to retrieve ip access list config from Atlas": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(nil, errors.New("failed to list"))

				return serviceMock
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to list")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to list"),
			},
		},
		"should fail to add ip access list in Atlas": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{}, nil)
				serviceMock.EXPECT().Add(context.Background(), "my-project", mock.AnythingOfType("ipaccesslist.IPAccessEntries")).
					Return(errors.New("failed to add ip access list"))

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "192.168.100.200",
				},
			},
			expectedResult: workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, errors.New("failed to add ip access list")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
					WithMessageRegexp("failed to add ip access list"),
			},
		},
		"should fail to remove ip access list in Atlas": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(
						ipaccesslist.IPAccessEntries{
							"sg-12345": {
								AWSSecurityGroup: "sg-12345",
							},
						},
						nil,
					)
				serviceMock.EXPECT().Add(context.Background(), "my-project", mock.AnythingOfType("ipaccesslist.IPAccessEntries")).
					Return(nil)
				serviceMock.EXPECT().Delete(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return(errors.New("failed to delete"))

				return serviceMock
			},
			ipAccessList:   nil,
			expectedResult: workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, errors.New("failed to delete")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
					WithMessageRegexp("failed to delete"),
			},
		},
		"should remove ip access list in Atlas and unmanage resource": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(
						ipaccesslist.IPAccessEntries{
							"sg-12345": {
								AWSSecurityGroup: "sg-12345",
							},
						},
						nil,
					)
				serviceMock.EXPECT().Add(context.Background(), "my-project", mock.AnythingOfType("ipaccesslist.IPAccessEntries")).
					Return(nil)
				serviceMock.EXPECT().Delete(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return(nil)

				return serviceMock
			},
			ipAccessList:       nil,
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
		"should add ip access list in Atlas and fail(request) to get status": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{}, nil)
				serviceMock.EXPECT().Add(context.Background(), "my-project", mock.AnythingOfType("ipaccesslist.IPAccessEntries")).
					Return(nil)
				serviceMock.EXPECT().Status(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return("", errors.New("failed to get status"))

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "192.168.100.200",
				},
			},
			expectedResult: workflow.Terminate(
				workflow.Internal,
				errors.New("failed to get status")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get status"),
			},
		},
		"should add ip access list in Atlas and wait for status to progress": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{}, nil)
				serviceMock.EXPECT().Add(context.Background(), "my-project", mock.AnythingOfType("ipaccesslist.IPAccessEntries")).
					Return(nil)
				serviceMock.EXPECT().Status(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return("PENDING", nil)

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "192.168.100.200",
				},
			},
			expectedResult: workflow.InProgress(
				workflow.ProjectIPAccessListNotActive,
				"atlas is adding access. this entry may not apply to all cloud providers at the time of this request"),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPAccessListNotActive)).
					WithMessageRegexp("atlas is adding access. this entry may not apply to all cloud providers at the time of this request"),
			},
		},
		"should terminate when status is failed": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{
						"sg-12345": {
							AWSSecurityGroup: "sg-12345",
						},
					}, nil)
				serviceMock.EXPECT().Status(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return("FAILED", nil)

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					AwsSecurityGroup: "sg-12345",
				},
			},
			expectedResult: workflow.Terminate(
				workflow.ProjectIPNotCreatedInAtlas,
				errors.New("atlas didn't succeed in adding this access entry")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
					WithMessageRegexp("atlas didn't succeed in adding this access entry"),
			},
		},
		"should be ready when status is active": {
			service: func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService {
				serviceMock.EXPECT().List(context.Background(), "my-project").
					Return(ipaccesslist.IPAccessEntries{
						"sg-12345": {
							AWSSecurityGroup: "sg-12345",
						},
					}, nil)
				serviceMock.EXPECT().Status(context.Background(), "my-project", mock.AnythingOfType("*ipaccesslist.IPAccessEntry")).
					Return("ACTIVE", nil)

				return serviceMock
			},
			ipAccessList: []project.IPAccessList{
				{
					AwsSecurityGroup: "sg-12345",
				},
			},
			expectedResult: workflow.OK(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.IPAccessListReadyType),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := &ipAccessListController{
				ctx: &workflow.Context{
					Context: context.Background(),
					Log:     zaptest.NewLogger(t).Sugar(),
				},
				project: &akov2.AtlasProject{
					Spec: akov2.AtlasProjectSpec{
						ProjectIPAccessList: tt.ipAccessList,
					},
					Status: status.AtlasProjectStatus{
						ID: "my-project",
					},
				},
				service: tt.service(translation.NewIPAccessListServiceMock(t)),
			}

			result := c.reconcile()
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, c.ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestHandleIPAccessList(t *testing.T) {
	tests := map[string]struct {
		ipAccessList       []project.IPAccessList
		expectedCalls      func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		expectedResult     workflow.Result
		expectedConditions []api.Condition
	}{
		"should successfully handle ip access list reconciliation": {
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListProjectIpAccessLists(context.Background(), "project-id").
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(
						&admin.PaginatedNetworkAccess{
							Results:    &[]admin.NetworkPermissionEntry{},
							TotalCount: pointer.MakePtr(0),
						},
						&http.Response{},
						nil,
					)
				apiMock.EXPECT().CreateProjectIpAccessList(context.Background(), "project-id", mock.AnythingOfType("*[]admin.NetworkPermissionEntry")).
					Return(admin.CreateProjectIpAccessListApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateProjectIpAccessListExecute(mock.AnythingOfType("admin.CreateProjectIpAccessListApiRequest")).
					Return(
						&admin.PaginatedNetworkAccess{
							Results: &[]admin.NetworkPermissionEntry{
								{
									IpAddress: pointer.MakePtr("192.168.100.150"),
									CidrBlock: pointer.MakePtr("192.168.100.150/32"),
								},
							},
							TotalCount: pointer.MakePtr(1),
						},
						&http.Response{},
						nil,
					)
				apiMock.EXPECT().GetProjectIpAccessListStatus(context.Background(), "project-id", "192.168.100.150/32").
					Return(admin.GetProjectIpAccessListStatusApiRequest{ApiService: apiMock})
				apiMock.EXPECT().GetProjectIpAccessListStatusExecute(mock.AnythingOfType("admin.GetProjectIpAccessListStatusApiRequest")).
					Return(
						&admin.NetworkPermissionEntryStatus{
							STATUS: "ACTIVE",
						},
						&http.Response{},
						nil,
					)

				return apiMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "192.168.100.150",
				},
			},
			expectedResult: workflow.OK(),
			expectedConditions: []api.Condition{
				api.TrueCondition(api.IPAccessListReadyType),
			},
		},
		"should fail to handle ip access list reconciliation": {
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListProjectIpAccessLists(context.Background(), "project-id").
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(
						nil,
						&http.Response{},
						errors.New("failed to list"),
					)

				return apiMock
			},
			ipAccessList: []project.IPAccessList{
				{
					IPAddress: "192.168.100.150",
				},
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to get ip access list from Atlas: failed to list")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.Internal)).
					WithMessageRegexp("failed to get ip access list from Atlas: failed to list"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := &workflow.Context{
				Context:   context.Background(),
				Log:       zaptest.NewLogger(t).Sugar(),
				SdkClient: &admin.APIClient{ProjectIPAccessListApi: tt.expectedCalls(mockadmin.NewProjectIPAccessListApi(t))},
			}
			p := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					ProjectIPAccessList: tt.ipAccessList,
				},
				Status: status.AtlasProjectStatus{
					ID: "project-id",
				},
			}

			result := handleIPAccessList(ctx, p)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}
