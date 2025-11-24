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

package ipaccesslist

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.mongodb.org/atlas-sdk/v20250312009/mockadmin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestIPAccessList_List(t *testing.T) {
	projectID := "my-project"
	active := time.Now().UTC().Add(time.Minute * 1)
	apiErr := admin.GenericOpenAPIError{}
	apiErr.SetError("failed to list")
	tests := map[string]struct {
		service  func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		expected IPAccessEntries
		err      error
	}{
		"should return empty when atlas is also empty": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListProjectIpAccessLists(context.Background(), projectID).
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(&admin.PaginatedNetworkAccess{}, &http.Response{}, nil)

				return apiMock
			},
			expected: IPAccessEntries{},
		},
		"should return converted entries from atlas result": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListProjectIpAccessLists(context.Background(), projectID).
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(
						&admin.PaginatedNetworkAccess{
							Results: &[]admin.NetworkPermissionEntry{
								{
									IpAddress:       pointer.MakePtr("192.168.100.150"),
									CidrBlock:       pointer.MakePtr("192.168.100.150/32"),
									DeleteAfterDate: &active,
								},
								{
									CidrBlock: pointer.MakePtr("192.168.1.0/24"),
									Comment:   pointer.MakePtr("My Network"),
								},
								{
									AwsSecurityGroup: pointer.MakePtr("sg-12345"),
								},
							},
							TotalCount: pointer.MakePtr(3),
						},
						&http.Response{},
						nil,
					)

				return apiMock
			},
			expected: IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &active,
				},
				"192.168.1.0/24": {
					CIDR:    "192.168.1.0/24",
					Comment: "My Network",
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			},
		},
		"should return error when request fails": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().ListProjectIpAccessLists(context.Background(), projectID).
					Return(admin.ListProjectIpAccessListsApiRequest{ApiService: apiMock})
				apiMock.EXPECT().ListProjectIpAccessListsExecute(mock.AnythingOfType("admin.ListProjectIpAccessListsApiRequest")).
					Return(nil, &http.Response{}, apiErr)

				return apiMock
			},
			err: fmt.Errorf("failed to get ip access list from Atlas: %w", apiErr),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := &IPAccessList{
				ipAccessListAPI: tt.service(mockadmin.NewProjectIPAccessListApi(t)),
			}

			entries, err := i.List(context.Background(), projectID)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, entries)
		})
	}
}

func TestIPAccessList_Add(t *testing.T) {
	projectID := "my-project"
	active := time.Now().UTC().Add(time.Minute * 1)
	apiErr := admin.GenericOpenAPIError{}
	apiErr.SetError("failed to create")
	tests := map[string]struct {
		service func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		entries IPAccessEntries
		err     error
	}{
		"should add ip access list": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().CreateProjectIpAccessList(context.Background(), projectID, mock.AnythingOfType("*[]admin.NetworkPermissionEntry")).
					Return(admin.CreateProjectIpAccessListApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateProjectIpAccessListExecute(mock.AnythingOfType("admin.CreateProjectIpAccessListApiRequest")).
					Return(
						&admin.PaginatedNetworkAccess{
							Results: &[]admin.NetworkPermissionEntry{
								{
									IpAddress:       pointer.MakePtr("192.168.100.150"),
									CidrBlock:       pointer.MakePtr("192.168.100.150/32"),
									DeleteAfterDate: &active,
								},
								{
									CidrBlock: pointer.MakePtr("192.168.1.0/24"),
									Comment:   pointer.MakePtr("My Network"),
								},
								{
									AwsSecurityGroup: pointer.MakePtr("sg-12345"),
								},
							},
							TotalCount: pointer.MakePtr(3),
						},
						&http.Response{},
						nil,
					)

				return apiMock
			},
			entries: IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &active,
				},
				"192.168.1.0/24": {
					CIDR:    "192.168.1.0/24",
					Comment: "My Network",
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			},
		},
		"should return error when request fails": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().CreateProjectIpAccessList(context.Background(), projectID, mock.AnythingOfType("*[]admin.NetworkPermissionEntry")).
					Return(admin.CreateProjectIpAccessListApiRequest{ApiService: apiMock})
				apiMock.EXPECT().CreateProjectIpAccessListExecute(mock.AnythingOfType("admin.CreateProjectIpAccessListApiRequest")).
					Return(nil, &http.Response{}, apiErr)

				return apiMock
			},
			entries: IPAccessEntries{
				"192.168.100.150/32": {
					CIDR:            "192.168.100.150/32",
					DeleteAfterDate: &active,
				},
				"192.168.1.0/24": {
					CIDR:    "192.168.1.0/24",
					Comment: "My Network",
				},
				"sg-12345": {
					AWSSecurityGroup: "sg-12345",
				},
			},
			err: fmt.Errorf("failed to create ip access list from Atlas: %w", apiErr),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := &IPAccessList{
				ipAccessListAPI: tt.service(mockadmin.NewProjectIPAccessListApi(t)),
			}

			err := i.Add(context.Background(), projectID, tt.entries)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestIPAccessList_Delete(t *testing.T) {
	projectID := "my-project"
	apiErr := admin.GenericOpenAPIError{}
	apiErr.SetError("failed to delete")
	tests := map[string]struct {
		service func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		entry   *IPAccessEntry
		err     error
	}{
		"should delete ip access list": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().DeleteProjectIpAccessList(context.Background(), projectID, "192.168.100.150/32").
					Return(admin.DeleteProjectIpAccessListApiRequest{ApiService: apiMock})
				apiMock.EXPECT().DeleteProjectIpAccessListExecute(mock.AnythingOfType("admin.DeleteProjectIpAccessListApiRequest")).
					Return(
						&http.Response{},
						nil,
					)

				return apiMock
			},
			entry: &IPAccessEntry{
				CIDR: "192.168.100.150/32",
			},
		},
		"should return error when request fails": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().DeleteProjectIpAccessList(context.Background(), projectID, "192.168.100.150/32").
					Return(admin.DeleteProjectIpAccessListApiRequest{ApiService: apiMock})
				apiMock.EXPECT().DeleteProjectIpAccessListExecute(mock.AnythingOfType("admin.DeleteProjectIpAccessListApiRequest")).
					Return(&http.Response{}, apiErr)

				return apiMock
			},
			entry: &IPAccessEntry{
				CIDR: "192.168.100.150/32",
			},
			err: fmt.Errorf("failed to delete ip access list from Atlas: %w", apiErr),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := &IPAccessList{
				ipAccessListAPI: tt.service(mockadmin.NewProjectIPAccessListApi(t)),
			}

			err := i.Delete(context.Background(), projectID, tt.entry)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestIPAccessList_Status(t *testing.T) {
	projectID := "my-project"
	apiErr := admin.GenericOpenAPIError{}
	apiErr.SetError("failed to get status")
	tests := map[string]struct {
		service  func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi
		entry    *IPAccessEntry
		expected string
		err      error
	}{
		"should get status of ip access list": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().GetProjectIpAccessListStatus(context.Background(), projectID, "192.168.100.150/32").
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
			entry: &IPAccessEntry{
				CIDR: "192.168.100.150/32",
			},
			expected: "ACTIVE",
		},
		"should return error when request fails": {
			service: func(apiMock *mockadmin.ProjectIPAccessListApi) admin.ProjectIPAccessListApi {
				apiMock.EXPECT().GetProjectIpAccessListStatus(context.Background(), projectID, "192.168.100.150/32").
					Return(admin.GetProjectIpAccessListStatusApiRequest{ApiService: apiMock})
				apiMock.EXPECT().GetProjectIpAccessListStatusExecute(mock.AnythingOfType("admin.GetProjectIpAccessListStatusApiRequest")).
					Return(nil, &http.Response{}, apiErr)

				return apiMock
			},
			entry: &IPAccessEntry{
				CIDR: "192.168.100.150/32",
			},
			err: fmt.Errorf("failed to get status of ip access list from Atlas: %w", apiErr),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			i := &IPAccessList{
				ipAccessListAPI: tt.service(mockadmin.NewProjectIPAccessListApi(t)),
			}

			stat, err := i.Status(context.Background(), projectID, tt.entry)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, stat)
		})
	}
}
