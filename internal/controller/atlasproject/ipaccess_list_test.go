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

package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"
	"go.mongodb.org/atlas-sdk/v20250312014/mockadmin"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
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
		expectedResult     workflow.DeprecatedResult
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
		expectedResult     workflow.DeprecatedResult
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
		expectedResult     workflow.DeprecatedResult
		expectedConditions []api.Condition
	}{
		"should fail getting last applied configuration": {
			annotations: map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"},
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				return apiMock
			},
			expectedResult: workflow.Terminate(workflow.Internal, errors.New("failed to get last applied configuration: error reading AtlasProject Spec from annotation [mongodb.com/last-applied-configuration]: invalid character 'w' looking for beginning of object key string")),
		},
		"should successfully handle ip access list reconciliation": {
			expectedCalls: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListAccessListEntries(context.Background(), "project-id").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListAccessListEntriesExecute(mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).
					Return(
						&admin.PaginatedNetworkAccess{
							Results:    &[]admin.NetworkPermissionEntry{},
							TotalCount: pointer.MakePtr(0),
						},
						&http.Response{},
						nil,
					)
				apiMock.EXPECT().CreateAccessListEntry(context.Background(), "project-id", mock.AnythingOfType("*[]admin.NetworkPermissionEntry")).
					Return(admin.CreateAccessListEntryApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateAccessListEntryExecute(mock.AnythingOfType("admin.CreateAccessListEntryApiRequest")).
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
				apiMock.EXPECT().GetAccessListStatus(context.Background(), "project-id", "192.168.100.150/32").
					Return(admin.GetAccessListStatusApiRequest{ApiService: apiMock})
				apiMock.EXPECT().GetAccessListStatusExecute(mock.AnythingOfType("admin.GetAccessListStatusApiRequest")).
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
				apiMock.EXPECT().ListAccessListEntries(context.Background(), "project-id").
					Return(admin.ListAccessListEntriesApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListAccessListEntriesExecute(mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).
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
				Context: context.Background(),
				Log:     zaptest.NewLogger(t).Sugar(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						ProjectIPAccessListApi: tt.expectedCalls(mockadmin.NewProjectIPAccessListApi(t)),
					},
				},
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
			if tt.expectedResult.GetError() != nil {
				assert.ErrorContains(t, result.GetError(), tt.expectedResult.GetError().Error())
			} else {
				assert.NoError(t, result.GetError())
			}
			assert.Equal(t, tt.expectedResult.CloneWithoutError(), result.CloneWithoutError())
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
			annotations:   map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"},
			expectedError: "invalid character 'w' looking for beginning of object key string",
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

func TestIPAccessListNonGreedyBehaviour(t *testing.T) {
	for _, tc := range []struct {
		title                   string
		lastAppliedIPAccessList []string
		specIPAccessList        []string
		atlasIPAccessList       []string
		wantRemoved             []string
	}{
		{
			title:                   "no last applied no removal in Atlas",
			lastAppliedIPAccessList: []string{},
			specIPAccessList:        []string{},
			atlasIPAccessList:       []string{"100.90.0.0/24", "101.99.0.0/24"},
			wantRemoved:             []string{},
		},
		{
			title:                   "removed from last applied removes from Atlas",
			lastAppliedIPAccessList: []string{"100.90.0.0/24", "101.99.0.0/24"},
			specIPAccessList:        []string{"100.90.0.0/24"},
			atlasIPAccessList:       []string{"100.90.0.0/24", "101.99.0.0/24"},
			wantRemoved:             []string{"101.99.0.0/24"},
		},
		{
			title:                   "removed all from last applied removes all from Atlas",
			lastAppliedIPAccessList: []string{"100.90.0.0/24", "101.99.0.0/24"},
			specIPAccessList:        []string{},
			atlasIPAccessList:       []string{"100.90.0.0/24", "101.99.0.0/24"},
			wantRemoved:             []string{"100.90.0.0/24", "101.99.0.0/24"},
		},
		{
			title:                   "not in last applied still removed from Atlas",
			lastAppliedIPAccessList: []string{"100.90.0.0/24"},
			specIPAccessList:        []string{"100.90.0.0/24"},
			atlasIPAccessList:       []string{"100.90.0.0/24", "101.99.0.0/24"},
			wantRemoved:             []string{"101.99.0.0/24"},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			prj := newIPAccessListTestProject(tc.specIPAccessList)
			lastPrj := newIPAccessListTestProject(tc.lastAppliedIPAccessList)
			prj.Annotations[customresource.AnnotationLastAppliedConfiguration] = jsonize(t, lastPrj.Spec)

			ipAccessAPI := mockadmin.NewProjectIPAccessListApi(t)
			ipAccessAPI.EXPECT().ListAccessListEntries(mock.Anything, mock.Anything).
				Return(admin.ListAccessListEntriesApiRequest{ApiService: ipAccessAPI}).Once()
			ipAccessAPI.EXPECT().ListAccessListEntriesExecute(
				mock.AnythingOfType("admin.ListAccessListEntriesApiRequest")).Return(
				synthesizeAtlasIPAccessList(tc.atlasIPAccessList), nil, nil,
			).Once()
			// CreateAccessListEntry is a non destrutive operation, it does not remove entries
			ipAccessAPI.EXPECT().CreateAccessListEntry(mock.Anything, mock.Anything, mock.Anything).
				Return(admin.CreateAccessListEntryApiRequest{ApiService: ipAccessAPI}).Once()
			ipAccessAPI.EXPECT().CreateAccessListEntryExecute(
				mock.AnythingOfType("admin.CreateAccessListEntryApiRequest")).Return(
				nil, nil, nil,
			).Once()

			removals := len(tc.wantRemoved)
			if removals > 0 {
				ipAccessAPI.EXPECT().DeleteAccessListEntry(
					mock.Anything, mock.Anything, mock.Anything,
				).Return(admin.DeleteAccessListEntryApiRequest{ApiService: ipAccessAPI}).Times(removals)
				ipAccessAPI.EXPECT().DeleteAccessListEntryExecute(
					mock.AnythingOfType("admin.DeleteAccessListEntryApiRequest")).Return(
					nil, nil,
				).Times(removals)
			}

			unset := len(tc.specIPAccessList) == 0
			if !unset {
				ipAccessAPI.EXPECT().GetAccessListStatus(
					mock.Anything, mock.Anything, mock.Anything,
				).Return(admin.GetAccessListStatusApiRequest{ApiService: ipAccessAPI}).Times(removals)
				ipAccessAPI.EXPECT().GetAccessListStatusExecute(
					mock.AnythingOfType("admin.GetAccessListStatusApiRequest")).Return(
					nil, nil, nil,
				).Times(removals)
			}

			workflowCtx := workflow.Context{
				Log:     zaptest.NewLogger(t).Sugar(),
				Context: context.Background(),
				SdkClientSet: &atlas.ClientSet{
					SdkClient20250312013: &admin.APIClient{
						ProjectIPAccessListApi: ipAccessAPI,
					},
				},
			}

			result := handleIPAccessList(&workflowCtx, prj)
			require.Equal(t, workflow.OK(), result)
		})
	}
}

func newIPAccessListTestProject(ipAccessList []string) *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
		},
		Spec: akov2.AtlasProjectSpec{
			Name:                "test-project",
			ProjectIPAccessList: synthesizeIPAccessList(ipAccessList),
		},
	}
}

func synthesizeIPAccessList(ipAccessList []string) []project.IPAccessList {
	peers := make([]project.IPAccessList, 0, len(ipAccessList))
	for _, cidr := range ipAccessList {
		peers = append(peers, project.IPAccessList{
			CIDRBlock: cidr,
			Comment:   fmt.Sprintf("fake CIDR block %s", cidr),
		})
	}
	return peers
}

func synthesizeAtlasIPAccessList(peeringIDs []string) *admin.PaginatedNetworkAccess {
	atlasIPAccessList := make([]admin.NetworkPermissionEntry, 0, len(peeringIDs))
	for _, cidr := range peeringIDs {
		atlasIPAccessList = append(atlasIPAccessList, admin.NetworkPermissionEntry{
			CidrBlock: pointer.MakePtr(cidr),
			Comment:   pointer.MakePtr(fmt.Sprintf("fake CIDR block %s", cidr)),
		})
	}
	return &admin.PaginatedNetworkAccess{
		Results: &atlasIPAccessList,
	}
}
