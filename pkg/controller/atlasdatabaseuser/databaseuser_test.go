package atlasdatabaseuser

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	mocked "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/dbuser"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	testUserPasswordName = "password-name"

	nonExistingCluster = "non-existing-cluster"

	testDeployment = "deployment"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}

func TestFilterScopeClusters(t *testing.T) {
	scopeSpecs := []akov2.ScopeSpec{{
		Name: "dbLake",
		Type: akov2.DataLakeScopeType,
	}, {
		Name: "cluster1",
		Type: akov2.DeploymentScopeType,
	}, {
		Name: "cluster2",
		Type: akov2.DeploymentScopeType,
	}}
	clusters := []string{"cluster1", "cluster4", "cluster5"}
	scopeClusters := filterScopeDeployments(akov2.AtlasDatabaseUser{Spec: akov2.AtlasDatabaseUserSpec{Scopes: scopeSpecs}}, clusters)
	assert.Equal(t, []string{"cluster1"}, scopeClusters)
}

func TestCheckUserExpired(t *testing.T) {
	// Fake client
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	t.Run("Validate DeleteAfterDate", func(t *testing.T) {
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "", *akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("foo"))
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())

		result = checkUserExpired(context.Background(), zap.S(), fakeClient, "", *akov2.DefaultDBUser("ns", "theuser", "").WithDeleteAfterDate("2021/11/30T15:04:05"))
		assert.False(t, result.IsOk())
	})
	t.Run("User Marked Expired", func(t *testing.T) {
		data := dataForSecret()
		// Create a connection secret
		_, err := connectionsecret.Ensure(context.Background(), fakeClient, "tetNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		// The secret for the other project
		_, err = connectionsecret.Ensure(context.Background(), fakeClient, "testNs", "project2", "dsfsdf234234sdfdsf23423", "cluster1", data)
		assert.NoError(t, err)

		before := time.Now().UTC().Add(time.Minute * -1).Format("2006-01-02T15:04:05.999Z")
		user := *akov2.DefaultDBUser("testNs", data.DBUserName, "").WithDeleteAfterDate(before)
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "603e7bf38a94956835659ae5", user)
		assert.False(t, result.IsOk())
		assert.Equal(t, reconcile.Result{}, result.ReconcileResult())

		// The secret has been removed
		secret := corev1.Secret{}
		secretName := fmt.Sprintf("%s-%s-%s", "project1", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.Error(t, err)

		// The other secret still exists
		secretName = fmt.Sprintf("%s-%s-%s", "project2", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.NoError(t, err)
	})
	t.Run("No expiration happened", func(t *testing.T) {
		data := dataForSecret()
		// Create a connection secret
		_, err := connectionsecret.Ensure(context.Background(), fakeClient, "testNs", "project1", "603e7bf38a94956835659ae5", "cluster1", data)
		assert.NoError(t, err)
		after := time.Now().UTC().Add(time.Minute * 1).Format("2006-01-02T15:04:05")
		result := checkUserExpired(context.Background(), zap.S(), fakeClient, "603e7bf38a94956835659ae5", *akov2.DefaultDBUser("testNs", data.DBUserName, "").WithDeleteAfterDate(after))
		assert.True(t, result.IsOk())

		// The secret is still there
		secret := corev1.Secret{}
		secretName := fmt.Sprintf("%s-%s-%s", "project1", "cluster1", kube.NormalizeIdentifier(data.DBUserName))
		err = fakeClient.Get(context.Background(), kube.ObjectKey("testNs", secretName), &secret)
		assert.NoError(t, err)
	})
}

func TestHandleUserNameChange(t *testing.T) {
	t.Run("Only one user after name change", func(t *testing.T) {
		projectID := "project1"
		username := "theuser"
		user := *akov2.DefaultDBUser("ns", "theuser", projectID)
		user.Spec.Username = "differentuser"
		user.Status.UserName = username
		ctx := workflow.NewContext(zap.S(), []api.Condition{}, nil)
		ctx.Context = context.Background()
		testUserAPI := mockadmin.NewDatabaseUsersApi(t)
		dus := dbuser.NewAtlasUsers(testUserAPI)
		testUserAPI.EXPECT().DeleteDatabaseUser(ctx.Context, projectID, "", username).Return(
			admin.DeleteDatabaseUserApiRequest{ApiService: testUserAPI})
		testUserAPI.EXPECT().DeleteDatabaseUserExecute(mock.Anything).Return(nil, nil, nil)
		result := handleUserNameChange(ctx, dus, projectID, user)
		assert.True(t, result.IsOk())
	})
}

func dataForSecret() connectionsecret.ConnectionData {
	return connectionsecret.ConnectionData{
		DBUserName: "admin",
		ConnURL:    "mongodb://mongodb0.example.com:27017,mongodb1.example.com:27017/?authSource=admin",
		SrvConnURL: "mongodb+srv://mongodb.example.com:27017/?authSource=admin",
		Password:   "m@gick%",
	}
}

func TestEnsureDatabaseUser(t *testing.T) {
	ctx := context.Background()
	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(akov2.AddToScheme(scheme))
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	project := defaultTestProject()
	user := defaultTestUser()
	user.Spec.PasswordSecret = &common.ResourceRef{Name: testUserPasswordName}
	user.Status.PasswordVersion = "1"
	require.NoError(t, fakeClient.Create(ctx, user))
	defer fakeClient.Delete(ctx, user)
	conditions := akov2.InitCondition(user, api.FalseCondition(api.ReadyType))
	log := zap.S()
	workflowCtx := workflow.NewContext(log, conditions, ctx)
	differentUser := defaultTestUser()
	differentUser.Spec.AWSIAMType = "USER"
	r := &AtlasDatabaseUserReconciler{
		Client: fakeClient,
		Log:    log,
		AtlasProvider: &atlas.TestProvider{
			IsCloudGovFunc:           func() bool { return false },
			GlobalFallbackSecretFunc: func() *client.ObjectKey { return nil },
		},
	}
	for _, tc := range []struct {
		title           string
		password        *corev1.Secret
		dateOverride    string
		scopeOverrides  []akov2.ScopeSpec
		nameOverride    string
		dus             dbuser.AtlasUsersService
		ds              deployment.AtlasDeploymentsService
		expectedMessage string
	}{
		{
			title:           "Missing password fails",
			expectedMessage: `secrets "password-name" not found`,
		},

		{
			title:           "Wrong date format fails",
			password:        defaultTestPassword(),
			dateOverride:    "this-is-not-a-proper-date",
			expectedMessage: `failed to parse "this-is-not-a-proper-date" to an ISO date: parsing time "this-is-not-a-proper-date"`,
		},

		{
			title:           "Expired user aborts",
			password:        defaultTestPassword(),
			dateOverride:    time.Now().Add(-time.Hour).Format("2006-01-02T15:04:05-07"),
			expectedMessage: `The database user is expired and has been removed from Atlas`,
		},

		{
			title:           "Invalid user scope aborts",
			password:        defaultTestPassword(),
			scopeOverrides:  []akov2.ScopeSpec{{Type: akov2.DeploymentScopeType, Name: nonExistingCluster}},
			ds:              fakeClusterExists(ctx, testProjectID, nonExistingCluster, false, nil),
			expectedMessage: `"scopes" field references deployment named "non-existing-cluster" but such deployment doesn't exist in Atlas'`,
		},

		{
			title:           "User get fails",
			password:        defaultTestPassword(),
			dus:             fakeGetUser(ctx, nil, errRandom),
			expectedMessage: errRandom.Error(),
		},

		{
			title:    "User not found is created successfully",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, nil, dbuser.ErrorNotFound)
				return withFakeCreateUser(service, ctx, internalUser(user), nil)
			}(),
			expectedMessage: `Clusters are scheduled to handle database users updates`,
		},

		{
			title:    "User not found tries to create but fails",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, nil, dbuser.ErrorNotFound)
				return withFakeCreateUser(service, ctx, internalUser(user), errRandom)
			}(),
			expectedMessage: errRandom.Error(),
		},

		{
			title:    "User found unchanged does nothing",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(user), nil)
				return withFakeUpdateUser(service, ctx, internalUser(user), nil)
			}(),
			ds: func() deployment.AtlasDeploymentsService {
				service := fakeListClusterNames(ctx, []string{testDeployment}, nil)
				service = withFakeDeploymentIsReady(service, ctx)
				return withFakeListDeploymentConnections(service, ctx, nil)
			}(),
			expectedMessage: "",
		},

		{
			title:    "User different from Atlas is updated successfully",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(differentUser), nil)
				return withFakeUpdateUser(service, ctx, internalUser(user), nil)
			}(),
			expectedMessage: `Clusters are scheduled to handle database users updates`,
		},

		{
			title:    "User different from Atlas tries to update but fails",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(differentUser), nil)
				return withFakeUpdateUser(service, ctx, internalUser(user), errRandom)
			}(),
			expectedMessage: errRandom.Error(),
		},

		{
			title:    "User found unchanged but fails to check clusters",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(user), nil)
				return withFakeUpdateUser(service, ctx, internalUser(user), nil)
			}(),
			ds:              fakeListClusterNames(ctx, []string{testDeployment}, errRandom),
			expectedMessage: errRandom.Error(),
		},

		{
			title:    "User found unchanged but fails to check connections",
			password: defaultTestPassword(),
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(user), nil)
				return withFakeUpdateUser(service, ctx, internalUser(user), nil)
			}(),
			ds: func() deployment.AtlasDeploymentsService {
				service := fakeListClusterNames(ctx, []string{testDeployment}, nil)
				service = withFakeDeploymentIsReady(service, ctx)
				return withFakeListDeploymentConnections(service, ctx, errRandom)
			}(),
			expectedMessage: errRandom.Error(),
		},

		{
			title:        "User found unchanged but changed name succeeds",
			password:     defaultTestPassword(),
			nameOverride: "some-other-name",
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(user), nil)
				service = withFakeUpdateUser(service, ctx, internalUser(user), nil)
				return withFakeUserDeletion(service, ctx, testDatabase, testProjectID, "some-other-name", nil)
			}(),
			ds: func() deployment.AtlasDeploymentsService {
				service := fakeListClusterNames(ctx, []string{testDeployment}, nil)
				service = withFakeDeploymentIsReady(service, ctx)
				return withFakeListDeploymentConnections(service, ctx, nil)
			}(),
		},

		{
			title:        "User found unchanged but changed name but fix fails",
			password:     defaultTestPassword(),
			nameOverride: "some-other-name",
			dus: func() dbuser.AtlasUsersService {
				service := fakeGetUser(ctx, internalUser(user), nil)
				service = withFakeUpdateUser(service, ctx, internalUser(user), nil)
				return withFakeUserDeletion(service, ctx, testDatabase, testProjectID, "some-other-name", errRandom)
			}(),
			ds: func() deployment.AtlasDeploymentsService {
				service := fakeListClusterNames(ctx, []string{testDeployment}, nil)
				service = withFakeDeploymentIsReady(service, ctx)
				return withFakeListDeploymentConnections(service, ctx, nil)
			}(),
			expectedMessage: errRandom.Error(),
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			if tc.password != nil {
				require.NoError(t, fakeClient.Create(ctx, tc.password))
				defer fakeClient.Delete(ctx, tc.password)
			}
			user.Spec.DeleteAfterDate = tc.dateOverride
			user.Spec.Scopes = tc.scopeOverrides
			user.Status.UserName = tc.nameOverride
			result := r.ensureDatabaseUser(workflowCtx, tc.dus, tc.ds, *project, *user)
			if tc.expectedMessage == "" {
				assert.Equal(t, true, result.IsOk())
			} else {
				assert.Equal(t, false, result.IsOk())
				assert.Contains(t, result.GetMessage(), tc.expectedMessage)
			}
		})
	}
}

func internalUser(user *akov2.AtlasDatabaseUser) *dbuser.User {
	return &dbuser.User{
		AtlasDatabaseUserSpec: &user.Spec,
		Password:              "some-secret-here",
		ProjectID:             testProjectID,
	}
}

func defaultTestPassword() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: testUserPasswordName},
		Data:       map[string][]byte{"password": []byte("some-secret-here")},
	}
}

func fakeGetUser(ctx context.Context, usr *dbuser.User, err error) *mocked.AtlasUsersServiceMock {
	service := mocked.AtlasUsersServiceMock{}
	service.EXPECT().Get(ctx, testDatabase, testProjectID, testUsername).Return(usr, err)
	return &service
}

func withFakeCreateUser(service *mocked.AtlasUsersServiceMock, ctx context.Context, usr *dbuser.User, err error) *mocked.AtlasUsersServiceMock {
	service.EXPECT().Create(ctx, usr).Return(err)
	return service
}

func withFakeUpdateUser(service *mocked.AtlasUsersServiceMock, ctx context.Context, usr *dbuser.User, err error) *mocked.AtlasUsersServiceMock {
	service.EXPECT().Update(ctx, usr).Return(err)
	return service
}

func fakeClusterExists(ctx context.Context, projectID, clusterName string, exists bool, err error) *mocked.AtlasDeploymentsServiceMock {
	service := mocked.AtlasDeploymentsServiceMock{}
	service.EXPECT().ClusterExists(ctx, projectID, clusterName).Return(exists, err)
	return &service
}

func fakeListClusterNames(ctx context.Context, names []string, err error) *mocked.AtlasDeploymentsServiceMock {
	service := mocked.AtlasDeploymentsServiceMock{}
	service.EXPECT().ListClusterNames(ctx, testProjectID).Return(names, err)
	return &service
}

func withFakeDeploymentIsReady(service *mocked.AtlasDeploymentsServiceMock, ctx context.Context) *mocked.AtlasDeploymentsServiceMock {
	service.EXPECT().DeploymentIsReady(ctx, testProjectID, testDeployment).Return(true, nil)
	return service
}

func withFakeListDeploymentConnections(service *mocked.AtlasDeploymentsServiceMock, ctx context.Context, err error) *mocked.AtlasDeploymentsServiceMock {
	service.EXPECT().ListDeploymentConnections(ctx, testProjectID).Return(nil, err)
	return service
}
