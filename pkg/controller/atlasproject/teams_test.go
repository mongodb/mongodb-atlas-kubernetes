package atlasproject

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestSyncAssignedTeams(t *testing.T) {
	tests := map[string]struct {
		teamsToAssign map[string]*akov2.Team
		expectedErr   error
	}{
		"should sync teams assigned": {
			teamsToAssign: map[string]*akov2.Team{
				"teamID_1": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_1",
					},
					Roles: []akov2.TeamRole{"GROUP_OWNER"},
				},
				"teamID_2": {
					TeamRef: common.ResourceRefNamespaced{
						Name: "teamName_2",
					},
					Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name: "projectName",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "projectName",
					Teams: []akov2.Team{
						{
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_1"},
							Roles:   []akov2.TeamRole{"GROUP_OWNER"},
						},
						{
							TeamRef: common.ResourceRefNamespaced{Name: "teamName_2"},
							Roles:   []akov2.TeamRole{"GROUP_READ_ONLY"},
						},
					},
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
					Teams: []status.ProjectTeamStatus{
						{
							ID: "teamID_1",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_1",
							},
						},
						{
							ID: "teamID_2",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_2",
							},
						},
						{
							ID: "teamID_3",
							TeamRef: common.ResourceRefNamespaced{
								Name: "teamName_3",
							},
						},
					},
				},
			}
			team1 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_1",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_1",
				},
				Status: status.TeamStatus{
					ID: "teamID_1",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}
			team2 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_2",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_2",
				},
				Status: status.TeamStatus{
					ID: "teamID_2",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}
			team3 := &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name: "teamName_3",
				},
				Spec: akov2.TeamSpec{
					Name: "teamName_3",
				},
				Status: status.TeamStatus{
					ID: "teamID_3",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(project, team1, team2, team3).
				WithStatusSubresource(project, team1, team2, team3).
				Build()

			atlasClient := &mongodbatlas.Client{
				Projects: &atlasmocks.ProjectsClientMock{
					GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
						return &mongodbatlas.TeamsAssigned{
							Links: nil,
							Results: []*mongodbatlas.Result{
								{
									Links:     nil,
									RoleNames: []string{"GROUP_OWNER"},
									TeamID:    "teamID_1",
								},
								{
									Links:     nil,
									RoleNames: []string{"GROUP_OWNER"},
									TeamID:    "teamID_2",
								},
								{
									Links:     nil,
									RoleNames: []string{"GROUP_READ_ONLY"},
									TeamID:    "teamID_3",
								},
							},
							TotalCount: 0,
						}, nil, nil
					},
					AddTeamsToProjectFunc: func(projectId string, teams []*mongodbatlas.ProjectTeam) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
						return &mongodbatlas.TeamsAssigned{}, nil, nil
					},
				},
				Teams: &atlasmocks.TeamsClientMock{
					ListFunc: func(orgID string) ([]mongodbatlas.Team, *mongodbatlas.Response, error) {
						return []mongodbatlas.Team{
							{
								ID:        "teamID_1",
								Name:      "teamName_1",
								Usernames: nil,
							},
							{
								ID:        "teamID_2",
								Name:      "teamName_2",
								Usernames: nil,
							},
							{
								ID:        "teamID_3",
								Name:      "teamName_3",
								Usernames: nil,
							},
						}, nil, nil
					},
					RemoveTeamFromProjectFunc: func(projectID string, teamID string) (*mongodbatlas.Response, error) {
						return nil, nil
					},
				},
			}

			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Log:    logger,
				Client: atlasClient,
			}
			r := &AtlasProjectReconciler{
				Client:        k8sClient,
				EventRecorder: record.NewFakeRecorder(10),
				Log:           logger,
			}

			err := r.syncAssignedTeams(ctx, "projectID", project, tt.teamsToAssign)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUpdateTeamState(t *testing.T) {
	tests := map[string]struct {
		team                     *akov2.AtlasTeam
		isRemoval                bool
		expectedAssignedProjects []status.TeamProject
	}{
		"should add project to team status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
			},
		},
		"should not duplicate projects already listed on status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{
							ID:   "projectID",
							Name: "projectName",
						},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
			},
		},
		"should remove project from team status": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
				},
			},
			isRemoval:                true,
			expectedAssignedProjects: nil,
		},
		"should not modify status of other assigned projects": {
			team: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testTeam",
					Namespace: "testNS",
				},
				Status: status.TeamStatus{
					ID: "testTeamStatus",
					Projects: []status.TeamProject{
						{
							ID:   "existingProjectID",
							Name: "existingProjectName",
						},
					},
				},
			},
			isRemoval: false,
			expectedAssignedProjects: []status.TeamProject{
				{
					ID:   "projectID",
					Name: "projectName",
				},
				{
					ID:   "existingProjectID",
					Name: "existingProjectName",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-secret",
				},
				Data: map[string][]byte{
					"orgId":         []byte("0987654321"),
					"publicApiKey":  []byte("api-pub-key"),
					"privateApiKey": []byte("api-priv-key"),
				},
				Type: "Opaque",
			}
			project := &akov2.AtlasProject{
				Spec: akov2.AtlasProjectSpec{
					Name: "projectName",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name: "my-secret",
					},
				},
				Status: status.AtlasProjectStatus{
					ID: "projectID",
				},
			}

			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(secret, project, tt.team).
				WithStatusSubresource(tt.team).
				Build()

			logger := zaptest.NewLogger(t).Sugar()
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}
			r := &AtlasProjectReconciler{
				Client:        k8sClient,
				EventRecorder: record.NewFakeRecorder(1),
				Log:           logger,
			}

			err := r.updateTeamState(ctx, project, &common.ResourceRefNamespaced{Name: tt.team.Name, Namespace: tt.team.Namespace}, tt.isRemoval)
			assert.NoError(t, err)

			actualTeam := &akov2.AtlasTeam{}
			assert.NoError(t, k8sClient.Get(context.Background(), client.ObjectKeyFromObject(tt.team), actualTeam))
			assert.Equal(t, tt.expectedAssignedProjects, actualTeam.Status.Projects)
		})
	}
}
