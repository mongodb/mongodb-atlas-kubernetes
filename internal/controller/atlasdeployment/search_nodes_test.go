// Copyright 2024 MongoDB Inc
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

package atlasdeployment

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"go.mongodb.org/atlas-sdk/v20250312013/mockadmin"
	"go.uber.org/zap/zaptest"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

func TestHandleSearchNodes(t *testing.T) {
	projectName := "test-project"
	projectID := "abc123"

	t.Run("get search nodes request errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 3)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				nil,
				&http.Response{},
				errors.New("get test error"),
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in AKO and atlas (update)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 4)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().UpdateClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name,
			&admin.ApiSearchDeploymentRequest{Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas()}).
			Return(admin.UpdateClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().UpdateClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S100_LOWCPU_NVME",
							NodeCount:    4,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
		assert.True(t, ctx.HasReason(workflow.SearchNodesUpdating))
	})

	t.Run("search nodes are in AKO and atlas, but are not IDLE", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 4)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("UPDATING"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    4,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
		assert.True(t, ctx.HasReason(workflow.SearchNodesUpdating))
	})

	t.Run("search nodes are in AKO and atlas but update errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 4)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().UpdateClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name,
			&admin.ApiSearchDeploymentRequest{Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas()}).
			Return(admin.UpdateClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().UpdateClusterSearchDeploymentExecute(mock.Anything).
			Return(
				nil,
				&http.Response{},
				errors.New("update test error"),
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in AKO but not in Atlas (create)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 3)

		mockError := makeMockError()

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				mockError,
			)
		searchAPI.EXPECT().CreateClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name, mock.AnythingOfType("*admin.ApiSearchDeploymentRequest")).
			Return(admin.CreateClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().CreateClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S100_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
		assert.True(t, ctx.HasReason(workflow.SearchNodesCreating))
	})

	t.Run("search nodes are in AKO but not in Atlas but create errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 3)

		mockError := makeMockError()

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				mockError,
			)
		searchAPI.EXPECT().CreateClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name, mock.AnythingOfType("*admin.ApiSearchDeploymentRequest")).
			Return(admin.CreateClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().CreateClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{},
				errors.New("create test error"),
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("search nodes are in Atlas but not in AKO (delete)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().DeleteClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.DeleteClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().DeleteClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
		assert.True(t, ctx.HasReason(workflow.SearchNodesDeleting))
	})

	t.Run("search nodes are in Atlas but not in AKO but delete errors", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)
		searchAPI.EXPECT().DeleteClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.DeleteClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().DeleteClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&http.Response{},
				errors.New("delete test error"),
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.False(t, result.IsOk())
		assert.True(t, result.IsWarning())
	})

	t.Run("no search nodes in Atlas nor in AKO (unmanaged)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		mockError := makeMockError()

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				mockError,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
		assert.Empty(t, ctx.Conditions())
	})

	t.Run("search nodes are still creating in Atlas, AKO is waiting (handle upserting)", func(t *testing.T) { //nolint:dupl
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 3)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("UPDATING"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		ctx.SetConditionFromResult(api.SearchNodesReadyType, workflow.InProgress(workflow.SearchNodesCreating, "search nodes creating"))

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
	})

	t.Run("search nodes are created in Atlas, AKO is waiting (handle upserting)", func(t *testing.T) { //nolint:dupl
		deployment := akov2.DefaultAWSDeployment("default", projectName).WithSearchNodes("S80_LOWCPU_NVME", 3)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("IDLE"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		ctx.SetConditionFromResult(api.SearchNodesReadyType, workflow.InProgress(workflow.SearchNodesCreating, "search nodes creating"))

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
	})

	t.Run("search nodes are deleting in Atlas, AKO is waiting (handle deleting)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{
					GroupId:   pointer.MakePtr(projectID),
					StateName: pointer.MakePtr("UPDATING"),
					Specs: &[]admin.ApiSearchDeploymentSpec{
						{
							InstanceSize: "S80_LOWCPU_NVME",
							NodeCount:    3,
						},
					},
				},
				&http.Response{},
				nil,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		ctx.SetConditionFromResult(api.SearchNodesReadyType, workflow.InProgress(workflow.SearchNodesDeleting, "search nodes creating"))

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsInProgress())
	})

	t.Run("search nodes are deleted in Atlas, AKO is waiting (handle deleting)", func(t *testing.T) {
		deployment := akov2.DefaultAWSDeployment("default", projectName)

		mockError := makeMockError()

		searchAPI := mockadmin.NewAtlasSearchApi(t)
		searchAPI.EXPECT().GetClusterSearchDeployment(context.Background(), projectID, deployment.Spec.DeploymentSpec.Name).
			Return(admin.GetClusterSearchDeploymentApiRequest{ApiService: searchAPI})
		searchAPI.EXPECT().GetClusterSearchDeploymentExecute(mock.Anything).
			Return(
				&admin.ApiSearchDeploymentResponse{},
				&http.Response{
					Status:     http.StatusText(http.StatusBadRequest),
					StatusCode: http.StatusBadRequest,
				},
				mockError,
			)

		ctx := &workflow.Context{
			SdkClientSet: &atlas.ClientSet{
				SdkClient20250312012: &admin.APIClient{
					AtlasSearchApi: searchAPI,
				},
			},
			Context: context.Background(),
			Log:     zaptest.NewLogger(t).Sugar(),
		}

		ctx.SetConditionFromResult(api.SearchNodesReadyType, workflow.InProgress(workflow.SearchNodesDeleting, "search nodes deleting"))

		result := handleSearchNodes(ctx, deployment, projectID)

		assert.True(t, result.IsOk())
	})
}

func makeMockError() *admin.GenericOpenAPIError {
	mockError := &admin.GenericOpenAPIError{}
	model := *admin.NewApiErrorWithDefaults()
	model.SetError(http.StatusBadRequest)
	mockError.SetModel(model)
	return mockError
}
