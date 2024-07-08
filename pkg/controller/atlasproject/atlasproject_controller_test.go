package atlasproject

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	atlasmocks "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestHandleDeletion(t *testing.T) {
	t.Run("should fail when unable to set finalizer", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			WithInterceptorFuncs(interceptor.Funcs{Patch: func(ctx context.Context, client client.WithWatch, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
				return errors.New("failed to set finalizer to the project")
			}}).
			Build()
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    zap.New(core).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.Terminate(workflow.AtlasFinalizerNotSet, "failed to set finalizer to the project"),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Empty(t, project.Finalizers)
		assert.Equal(t, "Add deletion finalizer", logs.All()[0].Message)
	})

	t.Run("should ensure finalizer is set while project exists", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			Build()
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    zap.New(core).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.OK(),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Equal(t, []string{customresource.FinalizerLabel}, project.Finalizers)
		assert.Equal(t, "Add deletion finalizer", logs.All()[0].Message)
	})

	t.Run("should fail when unable to check project dependencies", func(t *testing.T) {
		deletionTime := metav1.Now()
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-project",
				Namespace:         "default",
				Finalizers:        []string{customresource.FinalizerLabel},
				DeletionTimestamp: &deletionTime,
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			WithIndex(
				instancesIndexer.Object(),
				instancesIndexer.Name(),
				instancesIndexer.Keys,
			).
			WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list streams instances")
			}}).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.Terminate(workflow.Internal, "failed to determine if project has dependencies: failed to list streams instances"),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Equal(t, []string{customresource.FinalizerLabel}, project.Finalizers)
	})

	t.Run("should fail when project was deleted but it has dependencies", func(t *testing.T) {
		deletionTime := metav1.Now()
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-project",
				Namespace:         "default",
				Finalizers:        []string{customresource.FinalizerLabel},
				DeletionTimestamp: &deletionTime,
			},
		}
		streamsInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance-0",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamsInstance).
			WithIndex(
				instancesIndexer.Object(),
				instancesIndexer.Name(),
				instancesIndexer.Keys,
			).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.Terminate(workflow.Internal, "the project cannot be deleted until dependencies were removed"),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Equal(t, []string{customresource.FinalizerLabel}, project.Finalizers)
	})

	t.Run("should do soft deletion when deletion protection is enabled", func(t *testing.T) {
		deletionTime := metav1.Now()
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-project",
				Namespace:         "default",
				DeletionTimestamp: &deletionTime,
				Finalizers:        []string{customresource.FinalizerLabel},
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			WithIndex(
				instancesIndexer.Object(),
				instancesIndexer.Name(),
				instancesIndexer.Keys,
			).
			Build()
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasProjectReconciler{
			Client:                   k8sClient,
			Log:                      zap.New(core).Sugar(),
			ObjectDeletionProtection: true,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.OK(),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Equal(
			t,
			"Not removing Project from Atlas as per configuration",
			logs.All()[0].Message,
		)
		assert.Empty(t, project.Finalizers)
	})

	t.Run("should do soft deletion when resource policy is set to keep", func(t *testing.T) {
		deletionTime := metav1.Now()
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-project",
				Namespace:         "default",
				DeletionTimestamp: &deletionTime,
				Finalizers:        []string{customresource.FinalizerLabel},
				Annotations: map[string]string{
					customresource.ResourcePolicyAnnotation: customresource.ResourcePolicyKeep,
				},
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithIndex(
				instancesIndexer.Object(),
				instancesIndexer.Name(),
				instancesIndexer.Keys,
			).
			WithObjects(project).
			Build()
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
			Log:    zap.New(core).Sugar(),
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		assert.Equal(
			t,
			workflow.OK(),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Equal(
			t,
			"Not removing Project from Atlas as per configuration",
			logs.All()[0].Message,
		)
		assert.Empty(t, project.Finalizers)
	})

	t.Run("Should delete team from Atlas when AtlasProject with finalizer is deleted", func(t *testing.T) {
		deletionTS := metav1.Now()
		team := &akov2.AtlasTeam{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-team",
				Namespace: "default",
			},
			Spec: akov2.TeamSpec{
				Name: "teamName",
			},
			Status: status.TeamStatus{
				ID: "teamID",
				Projects: []status.TeamProject{
					{
						ID:   "projectID",
						Name: "project",
					},
				},
			},
		}
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "my-project",
				Namespace:         "default",
				DeletionTimestamp: &deletionTS,
				Finalizers:        []string{customresource.FinalizerLabel},
			},
			Spec: akov2.AtlasProjectSpec{
				Teams: []akov2.Team{
					{
						TeamRef: common.ResourceRefNamespaced{
							Name:      "my-team",
							Namespace: "default",
						},
						Roles: []akov2.TeamRole{
							"PROJECT_OWNER",
						},
					},
				},
			},
			Status: status.AtlasProjectStatus{
				ID: "projectID",
				Teams: []status.ProjectTeamStatus{
					{
						ID: team.Status.ID,
						TeamRef: common.ResourceRefNamespaced{
							Name:      "my-team",
							Namespace: "default",
						},
					},
				},
			},
		}

		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		require.NoError(t, corev1.AddToScheme(testScheme))

		instancesIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))

		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, team).
			WithIndex(
				instancesIndexer.Object(),
				instancesIndexer.Name(),
				instancesIndexer.Keys,
			).
			Build()

		projectsMock := &atlasmocks.ProjectsClientMock{
			GetProjectTeamsAssignedFunc: func(projectID string) (*mongodbatlas.TeamsAssigned, *mongodbatlas.Response, error) {
				return &mongodbatlas.TeamsAssigned{
					Links: nil,
					Results: []*mongodbatlas.Result{
						{
							Links:     nil,
							RoleNames: nil,
							TeamID:    team.Status.ID,
						},
					},
					TotalCount: 0,
				}, nil, nil
			},
			DeleteFunc: func(projectID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
		}
		teamsMock := &atlasmocks.TeamsClientMock{
			RemoveTeamFromOrganizationFunc: func(orgID string, teamID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
			RemoveTeamFromOrganizationRequests: map[string]struct{}{},
			ListFunc: func(orgID string) ([]mongodbatlas.Team, *mongodbatlas.Response, error) {
				return []mongodbatlas.Team{
					{
						ID:        team.Status.ID,
						Name:      team.Name,
						Usernames: nil,
					},
				}, nil, nil
			},
			RemoveTeamFromProjectFunc: func(projectID string, teamID string) (*mongodbatlas.Response, error) {
				return nil, nil
			},
		}
		legacyClient := &mongodbatlas.Client{
			Projects: projectsMock,
			Teams:    teamsMock,
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

		logger := zaptest.NewLogger(t).Sugar()
		reconciler := AtlasProjectReconciler{
			Client: k8sClient,
			Log:    logger,
			AtlasProvider: &atlasmocks.TestProvider{
				ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
					return legacyClient, "123", nil
				},
			},
		}
		ctx := &workflow.Context{
			Client:  legacyClient,
			Context: context.Background(),
			SdkClient: &admin.APIClient{
				PrivateEndpointServicesApi: mockPrivateEndpointAPI,
				NetworkPeeringApi:          mockPeeringEndpointAPI,
			},
			Log: logger,
		}

		assert.Equal(
			t,
			workflow.OK(),
			reconciler.handleDeletion(ctx, project),
		)
		assert.Len(t, teamsMock.RemoveTeamFromOrganizationRequests, 1)
		assert.Empty(t, project.Finalizers)
	})
}

func TestHasDependencies(t *testing.T) {
	t.Run("should return error when unable to list stream instances", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			WithInterceptorFuncs(interceptor.Funcs{List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
				return errors.New("failed to list instances")
			}}).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, project)
		require.ErrorContains(t, err, "failed to list instances")
		assert.False(t, ok)
	})

	t.Run("should return false when project has no dependencies", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project).
			WithIndex(
				instanceIndexer.Object(),
				instanceIndexer.Name(),
				instanceIndexer.Keys,
			).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, project)
		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should return true when project has dependencies", func(t *testing.T) {
		project := &akov2.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
		}
		streamsInstance := &akov2.AtlasStreamInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance-0",
				Namespace: "default",
			},
			Spec: akov2.AtlasStreamInstanceSpec{
				Project: common.ResourceRefNamespaced{
					Name:      "my-project",
					Namespace: "default",
				},
			},
		}
		instanceIndexer := indexer.NewAtlasStreamInstanceByProjectIndexer(zaptest.NewLogger(t))
		testScheme := runtime.NewScheme()
		require.NoError(t, akov2.AddToScheme(testScheme))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(project, streamsInstance).
			WithIndex(
				instanceIndexer.Object(),
				instanceIndexer.Name(),
				instanceIndexer.Keys,
			).
			Build()
		reconciler := &AtlasProjectReconciler{
			Client: k8sClient,
		}
		ctx := &workflow.Context{
			Context: context.Background(),
		}

		ok, err := reconciler.hasDependencies(ctx, project)
		require.NoError(t, err)
		assert.True(t, ok)
	})
}

func TestFindProjectsForBCP(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &akov2.AtlasBackupCompliancePolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-bcp", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "test-bcp", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns247"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "test-bcp", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						BackupCompliancePolicyRef: &common.ResourceRefNamespaced{Name: "different-test-bcp", Namespace: "ns1"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns2"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns247"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByBackupCompliancePolicyIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[akov2.AtlasBackupCompliancePolicy](indexer.AtlasProjectByBackupCompliancePolicyIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}

func TestFindProjectsForConnectionSecret(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: "test-secret", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret"},
						AlertConfigurations: []akov2.AlertConfiguration{
							{
								Notifications: []akov2.Notification{
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}}, // double entry
									{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{OpsGenieAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{ServiceKeyRef: common.ResourceRefNamespaced{Name: "test-secret"}},
									{VictorOpsSecretRef: common.ResourceRefNamespaced{Name: "test-secret"}},
								},
							},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns3"},
					Spec: akov2.AtlasProjectSpec{
						ConnectionSecret: &common.ResourceRefNamespaced{Name: "test-secret"},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-4", Namespace: "ns4"},
					Spec: akov2.AtlasProjectSpec{
						AlertConfigurations: []akov2.AlertConfiguration{
							{
								Notifications: []akov2.Notification{
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}},
									{APITokenRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}}, // double entry
									{DatadogAPIKeyRef: common.ResourceRefNamespaced{Name: "test-secret", Namespace: "ns1"}},
									{FlowdockAPITokenRef: common.ResourceRefNamespaced{Name: "test-secret2", Namespace: "ns1"}},
								},
							},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns1"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns2"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-4", Namespace: "ns4"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByConnectionSecretIndexer(zaptest.NewLogger(t))

			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[corev1.Secret](indexer.AtlasProjectBySecretsIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", tc.want, got)
			}
		})
	}
}

func TestFindProjectsForTeams(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasDeployment{},
			want: nil,
		},
		{
			name: "test",
			obj: &akov2.AtlasTeam{
				ObjectMeta: metav1.ObjectMeta{Name: "test-team", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project", Namespace: "ns1"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team"}},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-2", Namespace: "ns2"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team", Namespace: "ns1"}},
						},
					},
				},
				&akov2.AtlasProject{
					ObjectMeta: metav1.ObjectMeta{Name: "test-project-3", Namespace: "ns3"},
					Spec: akov2.AtlasProjectSpec{
						Teams: []akov2.Team{
							{TeamRef: common.ResourceRefNamespaced{Name: "test-team"}},
						},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "test-project", Namespace: "ns1"}},
				{NamespacedName: types.NamespacedName{Name: "test-project-2", Namespace: "ns2"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			projectIndexer := indexer.NewAtlasProjectByTeamIndexer(zaptest.NewLogger(t))

			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(projectIndexer.Object(), projectIndexer.Name(), projectIndexer.Keys).
				Build()

			find := newProjectsMapFunc[akov2.AtlasTeam](indexer.AtlasProjectByTeamIndex, k8sClient, zaptest.NewLogger(t).Sugar())
			got := find(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", tc.want, got)
			}
		})
	}
}
