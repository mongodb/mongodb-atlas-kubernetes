package atlasproject

import (
	"context"
	"errors"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestCanAssignedTeamsReconcile(t *testing.T) {
	team1 := &akov2.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "team1",
			Namespace: "default",
		},
		Status: status.TeamStatus{
			ID: "team1",
		},
	}
	team2 := &akov2.AtlasTeam{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "team2",
			Namespace: "default",
		},
		Status: status.TeamStatus{
			ID: "team2",
		},
	}

	testScheme := runtime.NewScheme()
	testScheme.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasProject{})
	testScheme.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasTeam{})
	k8sClient := fake.NewClientBuilder().
		WithScheme(testScheme).
		WithObjects(team1, team2).
		Build()

	t.Run("should return true when subResourceDeletionProtection is disabled", func(t *testing.T) {
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, false, &akov2.AtlasProject{})
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return error when unable to deserialize last applied configuration", func(t *testing.T) {
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{wrong}"})
		workflowCtx := &workflow.Context{
			Client:  &mongodbatlas.Client{},
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)
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
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				Teams: []akov2.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team2",
							Namespace: "default",
						},
						Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				Teams: []akov2.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team2",
							Namespace: "default",
						},
						Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				Teams: []akov2.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team3",
							Namespace: "default",
						},
						Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Context: context.Background(),
		}
		result, err := canAssignedTeamsReconcile(workflowCtx, k8sClient, true, akoProject)

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
		akoProject := &akov2.AtlasProject{}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: "{}"})
		logger := zaptest.NewLogger(t).Sugar()

		testScheme := runtime.NewScheme()
		akov2.AddToScheme(testScheme)
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			Build()

		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Log:     logger,
			Context: context.Background(),
		}
		reconciler := &AtlasProjectReconciler{
			Log:    logger,
			Client: k8sClient,
		}
		result := reconciler.ensureAssignedTeams(workflowCtx, akoProject, true)

		require.Equal(t, workflow.Terminate(workflow.Internal, "unable to resolve ownership for deletion protection: failed to retrieve data"), result)
	})

	t.Run("should failed to reconcile when unable to synchronize with Atlas", func(t *testing.T) {
		team1 := &akov2.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "team1",
				Namespace: "default",
			},
			Status: status.TeamStatus{
				ID: "team1",
			},
		}
		team2 := &akov2.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "team2",
				Namespace: "default",
			},
			Status: status.TeamStatus{
				ID: "team2",
			},
		}

		testScheme := runtime.NewScheme()
		akov2.AddToScheme(testScheme)
		corev1.AddToScheme(testScheme)
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
		akoProject := &akov2.AtlasProject{
			Spec: akov2.AtlasProjectSpec{
				Teams: []akov2.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "team3",
							Namespace: "default",
						},
						Roles: []akov2.TeamRole{"GROUP_READ_ONLY"},
					},
				},
			},
		}
		akoProject.WithAnnotations(map[string]string{customresource.AnnotationLastAppliedConfiguration: `{"teams":[{"teamRef":{"name":"team1","namespace":"default"},"roles":["GROUP_OWNER"]}]}`})
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := &workflow.Context{
			Client:  &atlasClient,
			Log:     logger,
			Context: context.Background(),
		}
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    logger,
		}
		result := reconciler.ensureAssignedTeams(workflowCtx, akoProject, true)

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
			Context: context.Background(),
			Log:     logger,
		}
		testScheme := runtime.NewScheme()
		akov2.AddToScheme(testScheme)
		corev1.AddToScheme(testScheme)
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
		team := &akov2.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testTeam",
				Namespace: "testNS",
			},
			Status: status.TeamStatus{
				ID: "testTeamStatus",
				Projects: []status.TeamProject{
					{
						ID:   project.Status.ID,
						Name: project.Spec.Name,
					},
				},
			},
		}
		atlasProvider := &atlas.TestProvider{
			ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
				return &mongodbatlas.Client{}, "0987654321", nil
			},
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(secret, project, team).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client:        k8sClient,
			Log:           logger,
			AtlasProvider: atlasProvider,
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
		k8sClient.Get(context.Background(), types.NamespacedName{Name: team.ObjectMeta.Name, Namespace: team.ObjectMeta.Namespace}, team)
		assert.Equal(t, 1, len(team.Status.Projects))
	})

	t.Run("must remove a team from Atlas is a team is unassigned", func(t *testing.T) {
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := &workflow.Context{
			Context: context.Background(),
			Log:     logger,
		}
		testScheme := runtime.NewScheme()
		akov2.AddToScheme(testScheme)
		corev1.AddToScheme(testScheme)
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
		team := &akov2.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "testTeam",
				Namespace: "testNS",
			},
			Status: status.TeamStatus{
				ID:       "testTeamStatus",
				Projects: []status.TeamProject{},
			},
		}
		teamsMock := &atlas.TeamsClientMock{
			RemoveTeamFromOrganizationFunc: func(orgID string, teamID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
			RemoveTeamFromOrganizationRequests: map[string]struct{}{},
		}
		atlasProvider := &atlas.TestProvider{
			ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
				return &mongodbatlas.Client{
					Teams: teamsMock,
				}, "0987654321", nil
			},
		}
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(secret, project, team).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client:        k8sClient,
			Log:           logger,
			AtlasProvider: atlasProvider,
		}
		teamRef := &common.ResourceRefNamespaced{
			Name:      team.Name,
			Namespace: "testNS",
		}

		err := reconciler.updateTeamState(workflowCtx, project, teamRef, true)
		assert.NoError(t, err)
		k8sClient.Get(context.Background(), types.NamespacedName{Name: team.ObjectMeta.Name, Namespace: team.ObjectMeta.Namespace}, team)
		assert.Len(t, teamsMock.RemoveTeamFromOrganizationRequests, 1)
	})
}
