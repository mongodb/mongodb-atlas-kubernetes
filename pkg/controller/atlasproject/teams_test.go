package atlasproject

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/internal/mocks/atlas"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func TestCanAssignedTeamsReconcile(t *testing.T) {
	team1 := &mdbv1.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "team1",
			Namespace: "default",
		},
		Status: status.TeamStatus{
			ID: "team1",
		},
	}
	team2 := &mdbv1.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "team2",
			Namespace: "default",
		},
		Status: status.TeamStatus{
			ID: "team2",
		},
	}

	testScheme := runtime.NewScheme()
	testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasProject{})
	testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasTeam{})
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(team1, team2).
		Build()

	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		result, err := canAssignedTeamsReconcile(context.TODO(), mongodbatlas.Client{}, k8sClient, false, &mdbv1.AtlasProject{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		result, err := canAssignedTeamsReconcile(context.TODO(), mongodbatlas.Client{}, k8sClient, true, akoProject)
		assert.EqualError(t, err, "invalid character 'w' looking for beginning of object key string")
		assert.False(t, result)
	})

	t.Run("should return error when unable to fetch data from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.EqualError(t, err, "failed to retrieve data")
		assert.False(t, result)
	})

	t.Run("should return true when return nil from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return nil, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when return empty list from Atlas", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{TotalCount: 0}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are no difference between current Atlas and previous applied configuration", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{
						Results: []*mongodbatlas.Result{
							{
								TeamID:    "team1",
								RoleNames: []string{"GROUP_OWNER"},
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Teams: []mdbv1.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team2",
							Namespace: "default",
						},
						Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return true when there are differences but new configuration synchronize operator", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{
						Results: []*mongodbatlas.Result{
							{
								TeamID:    "team2",
								RoleNames: []string{"GROUP_READ_ONLY"},
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Teams: []mdbv1.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team2",
							Namespace: "default",
						},
						Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false when unable to reconcile assigned teams", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{
						Results: []*mongodbatlas.Result{
							{
								TeamID:    "team2",
								RoleNames: []string{"GROUP_READ_ONLY"},
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Teams: []mdbv1.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team3",
							Namespace: "default",
						},
						Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		result, err := canAssignedTeamsReconcile(context.TODO(), atlasClient, k8sClient, true, akoProject)

		assert.NoError(t, err)
		assert.False(t, result)
	})
}

func TestEnsureAssignedTeams(t *testing.T) {
	t.Run("should failed to reconcile when unable to decide resource ownership", func(t *testing.T) {
		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return nil, nil, errors.New("failed to retrieve data")
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    logger,
		}
		reconciler := &AtlasProjectReconciler{
			Log: logger,
		}
		result := reconciler.ensureAssignedTeams(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		team1 := &mdbv1.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "team1",
				Namespace: "default",
			},
			Status: status.TeamStatus{
				ID: "team1",
			},
		}
		team2 := &mdbv1.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "team2",
				Namespace: "default",
			},
			Status: status.TeamStatus{
				ID: "team2",
			},
		}

		testScheme := runtime.NewScheme()
		testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasProject{})
		testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasTeam{})
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(team1, team2).
			Build()

		atlasClient := mongodbatlas.Client{
			Projects: &atlas.ProjectsClientMock{
				GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
					return &mongodbatlas.TeamsAssigned{
						Results: []*mongodbatlas.Result{
							{
								TeamID:    "team2",
								RoleNames: []string{"GROUP_READ_ONLY"},
							},
						},
						TotalCount: 1,
					}, nil, nil
				},
			},
		}
		akoProject := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Teams: []mdbv1.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team3",
							Namespace: "default",
						},
						Roles: []mdbv1.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := &workflow.Context{
			Client: atlasClient,
			Log:    logger,
		}
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    logger,
		}
		result := reconciler.ensureAssignedTeams(context.TODO(), workflowCtx, akoProject, true)

		require.Equal(
			t,
			workflow.Terminate(
				workflow.AtlasDeletionProtection,
				"unable to reconcile Assigned Teams due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
			),
			result,
		)
	})
}

func TestUpdateTeamState(t *testing.T) {
	t.Run("should not duplicate projects listed", func(t *testing.T) {
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := &workflow.Context{
			Log: logger,
		}
		testScheme := runtime.NewScheme()
		testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasProject{})
		testScheme.AddKnownTypes(mdbv1.GroupVersion, &mdbv1.AtlasTeam{})
		project := &mdbv1.AtlasProject{
			Spec: mdbv1.AtlasProjectSpec{
				Name: "projectName",
			},
			Status: status.AtlasProjectStatus{
				ID: "projectID",
			},
		}
		team := &mdbv1.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testTeam",
				Namespace: "testNS",
			},
			Status: status.TeamStatus{
				ID: "testTeamStaatus",
				Projects: []status.TeamProject{
					{
						ID:   project.Status.ID,
						Name: project.Spec.Name,
					},
				},
			},
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(team).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    logger,
		}
		teamRef := &common.ResourceRefNamespaced{
			Name:      team.Name,
			Namespace: "testNS",
		}
		// check we have exactly 1 project in status
		assert.Equal(t, 1, len(team.Status.Projects))

		// "reconcile" the team state and check we still have 1 project in status
		err := reconciler.updateTeamState(workflowCtx, project, teamRef, false)
		assert.NoError(t, err)
		k8sClient.Get(context.TODO(), types.NamespacedName{Name: team.ObjectMeta.Name, Namespace: team.ObjectMeta.Namespace}, team)
		assert.Equal(t, 1, len(team.Status.Projects))
	})
}
