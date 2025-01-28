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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/ipaccesslist"
)

func TestIpAccessListController_reconcile(t *testing.T) {
	tests := map[string]struct {
		service            func(serviceMock *translation.IPAccessListServiceMock) ipaccesslist.IPAccessListService
		ipAccessList       []project.IPAccessList
		lastApplied        ipaccesslist.IPAccessEntries
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
			ipAccessList: nil,
			lastApplied: ipaccesslist.IPAccessEntries{
				"sg-12345": {AWSSecurityGroup: "sg-12345"},
			},
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
			ipAccessList: nil,
			lastApplied: ipaccesslist.IPAccessEntries{
				"sg-12345": {AWSSecurityGroup: "sg-12345"},
			},
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
				service:     tt.service(translation.NewIPAccessListServiceMock(t)),
				lastApplied: tt.lastApplied,
			}

			result := c.reconcile()
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, c.ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestConfigure(t *testing.T) {
	tests := map[string]struct {
		current            ipaccesslist.IPAccessEntries
		desired            ipaccesslist.IPAccessEntries
		lastApplied        ipaccesslist.IPAccessEntries
		expectedCalls      func() ipaccesslist.IPAccessListService
		expectedResult     workflow.Result
		expectedConditions []api.Condition
	}{
		"should fail to add ip access list": {
			current:     ipaccesslist.IPAccessEntries{},
			desired:     ipaccesslist.IPAccessEntries{"10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			lastApplied: ipaccesslist.IPAccessEntries{},
			expectedCalls: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{"10.0.0.0/24": {CIDR: "10.0.0.0/24"}}).
					Return(errors.New("failed to add ip access list"))

				return s
			},
			expectedResult: workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, errors.New("failed to add ip access list")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
					WithMessageRegexp("failed to add ip access list"),
			},
		},
		"should fail to delete ip access list": {
			current:     ipaccesslist.IPAccessEntries{"10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			desired:     ipaccesslist.IPAccessEntries{},
			lastApplied: ipaccesslist.IPAccessEntries{"10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			expectedCalls: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{}).
					Return(nil)
				s.EXPECT().Delete(context.Background(), "", &ipaccesslist.IPAccessEntry{CIDR: "10.0.0.0/24"}).
					Return(errors.New("failed to delete ip access list entry"))

				return s
			},
			expectedResult: workflow.Terminate(workflow.ProjectIPNotCreatedInAtlas, errors.New("failed to delete ip access list entry")),
			expectedConditions: []api.Condition{
				api.FalseCondition(api.IPAccessListReadyType).
					WithReason(string(workflow.ProjectIPNotCreatedInAtlas)).
					WithMessageRegexp("failed to delete ip access list entry"),
			},
		},
		"should no delete ip access list which were not previously managed": {
			current:     ipaccesslist.IPAccessEntries{"10.0.0.0/24": {CIDR: "10.0.0.0/24"}},
			desired:     ipaccesslist.IPAccessEntries{},
			lastApplied: ipaccesslist.IPAccessEntries{},
			expectedCalls: func() ipaccesslist.IPAccessListService {
				s := translation.NewIPAccessListServiceMock(t)
				s.EXPECT().Add(context.Background(), "", ipaccesslist.IPAccessEntries{}).
					Return(nil)

				return s
			},
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{api.TrueCondition(api.IPAccessListReadyType)},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
			}
			reconciler := ipAccessListController{
				ctx:         ctx,
				project:     &akov2.AtlasProject{},
				service:     tt.expectedCalls(),
				lastApplied: tt.lastApplied,
			}

			result := reconciler.configure(tt.current, tt.desired, false)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestHandleIPAccessList(t *testing.T) {
	tests := map[string]struct {
		ipAccessList       []project.IPAccessList
		annotations        map[string]string
		expectedCalls      func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		expectedResult     workflow.Result
		expectedConditions []api.Condition
	}{
		"should fail resolving last skipped flag": {
			annotations: map[string]string{customresource.AnnotationLastSkippedConfiguration: "{wrong}"},
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				return apiMock
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to parse last skipped configuration: invalid character 'w' looking for beginning of object key string")),
		},
		"should skip reconciliation": {
			annotations: map[string]string{customresource.AnnotationLastSkippedConfiguration: "{\"projectIpAccessList\": []}"},
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				return apiMock
			},
			expectedResult:     workflow.OK(),
			expectedConditions: []api.Condition{},
		},
		"should fail getting last applied configuration": {
			annotations: map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"},
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				return apiMock
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to get last applied configuration: error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]: invalid character 'w' looking for beginning of object key string")),
		},
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
			p.WithAnnotations(tt.annotations)

			result := handleIPAccessList(ctx, p)
			assert.Equal(t, tt.expectedResult, result)
			assert.True(t, cmp.Equal(tt.expectedConditions, ctx.Conditions(), cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime")))
		})
	}
}

func TestMapLastAppliedIPAccessList(t *testing.T) {
	tests := map[string]struct {
		annotations         map[string]string
		expectedIPAccessLst ipaccesslist.IPAccessEntries
		expectedError       string
	}{
		"should return error when last spec annotation is wrong": {
			annotations: map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"},
			expectedError: "error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]:" +
				" invalid character 'w' looking for beginning of object key string",
		},
		"should return nil when there is no last spec": {},
		"should return map of last ip access list": {
			annotations: map[string]string{
				customresource.AnnotationLastAppliedConfiguration: "{\"projectIpAccessList\": [{\"ipAddress\":\"192.168.0.100\"},{\"awsSecurityGroup\":\"sg-123456\",\"comment\":\"My AWS SG\"}]}"},
			expectedIPAccessLst: ipaccesslist.IPAccessEntries{
				"192.168.0.100/32": {
					CIDR: "192.168.0.100/32",
				},
				"sg-123456": {
					AWSSecurityGroup: "sg-123456",
					Comment:          "My AWS SG",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &akov2.AtlasProject{}
			p.WithAnnotations(tt.annotations)

			result, err := mapLastAppliedIPAccessList(p)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expectedIPAccessLst, result)
		})
	}
}

func TestHasSkippedIPAccessListConfiguration(t *testing.T) {
	tests := map[string]struct {
		annotations   map[string]string
		expected      bool
		expectedError string
	}{
		"should return error when last spec annotation is wrong": {
			annotations: map[string]string{customresource.AnnotationLastSkippedConfiguration: "{wrong}"},
			expectedError: "failed to parse last skipped configuration:" +
				" invalid character 'w' looking for beginning of object key string",
		},
		"should return false when there is no annotation": {},
		"should return false where there are last ip access list": {
			annotations: map[string]string{
				customresource.AnnotationLastSkippedConfiguration: "{\"projectIpAccessList\": [{\"ipAddress\":\"192.168.0.100\"},{\"awsSecurityGroup\":\"sg-123456\",\"comment\":\"My AWS SG\"}]}"},
		},
		"should return true where last ip access list is empty": {
			annotations: map[string]string{
				customresource.AnnotationLastSkippedConfiguration: "{\"projectIpAccessList\": []}"},
			expected: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			p := &akov2.AtlasProject{}
			p.WithAnnotations(tt.annotations)

			result, err := shouldIPAccessListSkipReconciliation(p)
			if err != nil {
				assert.ErrorContains(t, err, tt.expectedError)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
