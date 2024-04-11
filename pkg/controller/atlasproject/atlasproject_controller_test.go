package atlasproject

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestAtlasProjectReconciler_handleDeletion(t *testing.T) {
	t.Run("Should delete team from Atlas when AtlasProject with finalizer is deleted", func(t *testing.T) {
		sch := runtime.NewScheme()
		akov2.AddToScheme(sch)
		corev1.AddToScheme(sch)
		deletionTS := metav1.Now()

		testTeam := &akov2.AtlasTeam{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: akov2.TeamSpec{
				Name: "teamName",
			},
			Status: status.TeamStatus{
				ID: "teamID",
				Projects: []status.TeamProject{
					{
						ID:   "projectID",
						Name: "testProject",
					},
				},
			},
		}
		testProject := &akov2.AtlasProject{
			TypeMeta: metav1.TypeMeta{
				Kind:       "AtlasProject",
				APIVersion: "atlas.mongodb.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              "testProject",
				Namespace:         "default",
				DeletionTimestamp: &deletionTS,
				Finalizers:        []string{customresource.FinalizerLabel},
			},
			Spec: akov2.AtlasProjectSpec{},
			Status: status.AtlasProjectStatus{
				ID: "projectID",
				Teams: []status.ProjectTeamStatus{
					{
						ID: testTeam.Status.ID,
					},
				},
			},
		}

		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(testProject, testTeam).Build()

		teamsMock := &atlasmocks.TeamsClientMock{
			RemoveTeamFromOrganizationFunc: func(orgID string, teamID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
			RemoveTeamFromOrganizationRequests: map[string]struct{}{},
			ListFunc: func(orgID string) ([]mongodbatlas.Team, *mongodbatlas.Response, error) {
				return []mongodbatlas.Team{
					{
						ID:        testTeam.Status.ID,
						Name:      testTeam.Name,
						Usernames: nil,
					},
				}, nil, nil
			},
			RemoveTeamFromProjectFunc: func(projectID string, teamID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
		}

		atlasClient := mongodbatlas.Client{
			Projects: &atlasmocks.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{
						Links: nil,
						Results: []*mongodbatlas.Result{
							{
								Links:     nil,
								RoleNames: nil,
								TeamID:    testTeam.Status.ID,
							},
						},
						TotalCount: 0,
					}, nil, nil
				},
				DeleteFunc: func(projectID string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			},
			Teams: teamsMock,
		}
		reconciler := AtlasProjectReconciler{
			Client:                      k8sClient,
			ResourceWatcher:             watch.ResourceWatcher{},
			Log:                         zap.S(),
			Scheme:                      sch,
			ObjectDeletionProtection:    false,
			SubObjectDeletionProtection: false,
			AtlasProvider: &atlasmocks.TestProvider{
				ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
					return &atlasClient, "123", nil
				},
			},
		}

		mockPrivateEndpointAPI := mockadmin.NewPrivateEndpointServicesApi(t)
		mockPrivateEndpointAPI.EXPECT().
			ListPrivateEndpointServices(mock.Anything, mock.Anything, mock.Anything).
			Return(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI})
		mockPrivateEndpointAPI.EXPECT().
			ListPrivateEndpointServicesExecute(admin.ListPrivateEndpointServicesApiRequest{ApiService: mockPrivateEndpointAPI}).
			Return([]admin.EndpointService{}, nil, nil)

		mockPeeringEndpointAPI := mockadmin.NewNetworkPeeringApi(t)
		mockPeeringEndpointAPI.EXPECT().ListPeeringConnectionsWithParams(mock.Anything, mock.Anything).
			Return(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI})
		mockPeeringEndpointAPI.EXPECT().
			ListPeeringConnectionsExecute(admin.ListPeeringConnectionsApiRequest{ApiService: mockPeeringEndpointAPI}).
			Return(&admin.PaginatedContainerPeer{}, nil, nil)

		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
			SdkClient: &admin.APIClient{
				PrivateEndpointServicesApi: mockPrivateEndpointAPI,
				NetworkPeeringApi:          mockPeeringEndpointAPI,
			},
			Log: zap.S(),
		}
		_ = reconciler.handleDeletion(workflowCtx, workflowCtx.Client, testProject)
		//assert.True(t, result.IsOk())
		fmt.Println("DEBUG", teamsMock.RemoveTeamFromOrganizationRequests)
		assert.Len(t, teamsMock.RemoveTeamFromOrganizationRequests, 1)
	})
}
