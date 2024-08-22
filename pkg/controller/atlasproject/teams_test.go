package atlasproject

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func TestUpdateTeamState(t *testing.T) {
	t.Run("should not duplicate projects listed", func(t *testing.T) {
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := defaultTestWorkflow(logger)
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
			HasGlobalFallbackSecretFunc: func() bool { return true },
		}
		k8sClient := buildFakeKubernetesClient(secret, project, team)
		reconciler := &AtlasProjectReconciler{
			Client:        k8sClient,
			Log:           logger,
			AtlasProvider: atlasProvider,
		}
		// check we have exactly 1 project in status
		assert.Equal(t, 1, len(team.Status.Projects))

		// "reconcile" the team state and check we still have 1 project in status
		err := reconciler.updateTeamState(workflowCtx, project, reference(team), false)
		assert.NoError(t, err)
		k8sClient.Get(context.Background(), types.NamespacedName{Name: team.ObjectMeta.Name, Namespace: team.ObjectMeta.Namespace}, team)
		assert.Equal(t, 1, len(team.Status.Projects))
	})

	t.Run("must remove a team from Atlas is a team is unassigned", func(t *testing.T) {
		logger := zaptest.NewLogger(t).Sugar()
		workflowCtx := defaultTestWorkflow(logger)
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
			HasGlobalFallbackSecretFunc: func() bool { return true },
		}
		k8sClient := buildFakeKubernetesClient(secret, project, team)
		reconciler := &AtlasProjectReconciler{
			Client:        k8sClient,
			Log:           logger,
			AtlasProvider: atlasProvider,
		}

		err := reconciler.updateTeamState(workflowCtx, project, reference(team), true)
		assert.NoError(t, err)
		k8sClient.Get(context.Background(), types.NamespacedName{Name: team.ObjectMeta.Name, Namespace: team.ObjectMeta.Namespace}, team)
		assert.Len(t, teamsMock.RemoveTeamFromOrganizationRequests, 1)
	})

	t.Run("must honor deletion protection flag for Teams", func(t *testing.T) {
		for _, tc := range []struct {
			title              string
			deletionProtection bool
			keepFlag           bool
			expectRemoval      bool
		}{
			{
				title:              "with deletion protection unassigned teams are not removed",
				deletionProtection: true,
				keepFlag:           false,
				expectRemoval:      false,
			},
			{
				title:              "without deletion protection unassigned teams are removed",
				deletionProtection: false,
				keepFlag:           false,
				expectRemoval:      true,
			},
			{
				title:              "with deletion protection & keep flag teams are not removed",
				deletionProtection: false,
				keepFlag:           true,
				expectRemoval:      false,
			},
			{
				title:              "without deletion protection but keep flag teams are not removed",
				deletionProtection: true,
				keepFlag:           true,
				expectRemoval:      false,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				logger := zaptest.NewLogger(t).Sugar()
				workflowCtx := defaultTestWorkflow(logger)
				project := &akov2.AtlasProject{
					Spec: akov2.AtlasProjectSpec{
						Name: "projectName",
					},
				}
				team := &akov2.AtlasTeam{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "testTeam",
						Namespace: "testNS",
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
					HasGlobalFallbackSecretFunc: func() bool { return true },
				}
				reconciler := &AtlasProjectReconciler{
					Client:                   buildFakeKubernetesClient(project, team),
					Log:                      logger,
					AtlasProvider:            atlasProvider,
					ObjectDeletionProtection: tc.deletionProtection,
				}
				if tc.keepFlag {
					customresource.SetAnnotation(project,
						customresource.ResourcePolicyAnnotation, customresource.ResourcePolicyKeep)
				}
				err := reconciler.updateTeamState(workflowCtx, project, reference(team), true)
				assert.NoError(t, err)
				expectedRemovals := 0
				if tc.expectRemoval {
					expectedRemovals = 1
				}
				assert.Len(t, teamsMock.RemoveTeamFromOrganizationRequests, expectedRemovals)
			})
		}
	})
}

func defaultTestWorkflow(logger *zap.SugaredLogger) *workflow.Context {
	return &workflow.Context{
		Context: context.Background(),
		Log:     logger,
	}
}

func defaultTestScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	akov2.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	return scheme
}

func buildFakeKubernetesClient(objects ...client.Object) client.WithWatch {
	return fake.NewClientBuilder().
		WithScheme(defaultTestScheme()).
		WithObjects(objects...).
		Build()
}

func reference(obj client.Object) *common.ResourceRefNamespaced {
	return &common.ResourceRefNamespaced{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	}
}
