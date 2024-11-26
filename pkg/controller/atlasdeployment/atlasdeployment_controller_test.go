/*
Copyright 2023 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"
	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"
	"net/http"
	"reflect"
	"testing"

	adminv20241113001 "go.mongodb.org/atlas-sdk/v20241113001/admin"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

func TestCleanupBindings(t *testing.T) {
	t.Run("without backup references, nothing happens on cleanup", func(t *testing.T) {
		r := &AtlasDeploymentReconciler{
			Log:    testLog(t),
			Client: testK8sClient(),
		}
		d := testDeployment("cluster", nil)

		// test cleanup
		assert.NoError(t, r.cleanupBindings(context.Background(), deployment.NewDeployment("project-id", d)))
	})

	t.Run("with unreferenced backups, still nothing happens on cleanup", func(t *testing.T) {
		r := &AtlasDeploymentReconciler{
			Log:    testLog(t),
			Client: testK8sClient(),
		}
		d := testDeployment("cluster", nil)
		require.NoError(t, r.Client.Create(context.Background(), d))
		policy := testBackupPolicy()
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test cleanup
		require.NoError(t, r.cleanupBindings(context.Background(), deployment.NewDeployment("project-id", d)))

		endPolicy := &akov2.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.Equal(t, []string{customresource.FinalizerLabel}, endPolicy.Finalizers)
		endSchedule := &akov2.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.Equal(t, []string{customresource.FinalizerLabel}, endSchedule.Finalizers)
	})

	t.Run("last deployment's referenced backups finalizers are cleaned up", func(t *testing.T) {
		atlasProvider := &atlasmock.TestProvider{
			IsSupportedFunc: func() bool {
				return true
			},
		}
		r := &AtlasDeploymentReconciler{
			Log:           testLog(t),
			Client:        testK8sClient(),
			AtlasProvider: atlasProvider,
		}
		policy := testBackupPolicy() // deployment -> schedule -> policy
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		d := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), d))
		schedule.Status.DeploymentIDs = []string{d.Spec.DeploymentSpec.Name}
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test ensureBackupPolicy and cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule)
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment.NewDeployment("project-id", d)))

		endPolicy := &akov2.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.Empty(t, endPolicy.Finalizers, "policy should end up with no finalizer")
		endSchedule := &akov2.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.Empty(t, endSchedule.Finalizers, "schedule should end up with no finalizer")
	})

	t.Run("referenced backups finalizers are NOT cleaned up if reachable by other deployment", func(t *testing.T) {
		atlasProvider := &atlasmock.TestProvider{
			IsSupportedFunc: func() bool {
				return true
			},
		}
		r := &AtlasDeploymentReconciler{
			Log:           testLog(t),
			Client:        testK8sClient(),
			AtlasProvider: atlasProvider,
		}
		policy := testBackupPolicy() // deployment + deployment2 -> schedule -> policy
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		d := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), d))
		d2 := testDeployment("2", schedule)
		require.NoError(t, r.Client.Create(context.Background(), d2))
		schedule.Status.DeploymentIDs = []string{
			d.Spec.DeploymentSpec.Name,
			d2.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule)
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment.NewDeployment("project-id", d)))

		endPolicy := &akov2.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.NotEmpty(t, endPolicy.Finalizers, "policy should keep the finalizer")
		endSchedule := &akov2.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.NotEmpty(t, endSchedule.Finalizers, "schedule should keep the finalizer")
	})

	t.Run("policy finalizer stays if still referenced", func(t *testing.T) {
		atlasProvider := &atlasmock.TestProvider{
			IsSupportedFunc: func() bool {
				return true
			},
		}
		r := &AtlasDeploymentReconciler{
			Log:           testLog(t),
			Client:        testK8sClient(),
			AtlasProvider: atlasProvider,
		}
		policy := testBackupPolicy() // deployment -> schedule + schedule2 -> policy
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		schedule2 := testBackupSchedule("2", policy)
		d := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), d))
		d2 := testDeployment("2", schedule2)
		require.NoError(t, r.Client.Create(context.Background(), d2))
		schedule.Status.DeploymentIDs = []string{
			d.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule))
		schedule2.Status.DeploymentIDs = []string{
			d2.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule2))
		policy.Status.BackupScheduleIDs = []string{
			fmt.Sprintf("%s/%s", schedule.Namespace, schedule.Name),
			fmt.Sprintf("%s/%s", schedule2.Namespace, schedule2.Name),
		}

		// test cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule)
		require.NoError(t, err)
		_, err = r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule2)
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment.NewDeployment("project-id", d)))

		endPolicy := &akov2.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(policy.Namespace, policy.Name), endPolicy))
		assert.NotEmpty(t, endPolicy.Finalizers, "policy should keep the finalizer")
		endSchedule := &akov2.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(schedule.Namespace, schedule.Name), endSchedule))
		assert.Empty(t, endSchedule.Finalizers, "schedule should end up with no finalizer")
	})
}

func testK8sClient() client.Client {
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	sch := runtime.NewScheme()
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasDeployment{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupSchedule{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupScheduleList{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupPolicy{})
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
	return fake.NewClientBuilder().WithScheme(sch).
		WithStatusSubresource(&akov2.AtlasBackupSchedule{}, &akov2.AtlasBackupPolicy{}).
		Build()
}

func testLog(t *testing.T) *zap.SugaredLogger {
	t.Helper()

	return zaptest.NewLogger(t).Sugar()
}

func testDeploymentName(suffix string) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("test-deployment%s", suffix),
		Namespace: "test-namespace",
	}
}

func testDeployment(suffix string, schedule *akov2.AtlasBackupSchedule) *akov2.AtlasDeployment {
	dn := testDeploymentName(suffix)
	d := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: dn.Name, Namespace: dn.Namespace},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name: fmt.Sprintf("atlas-%s", dn.Name),
			},
		},
	}

	if schedule != nil {
		d.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
			Name:      schedule.Name,
			Namespace: schedule.Namespace,
		}
	}

	return d
}

func testBackupSchedule(suffix string, policy *akov2.AtlasBackupPolicy) *akov2.AtlasBackupSchedule {
	return &akov2.AtlasBackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-backup-schedule%s", suffix),
			Namespace:  "test-namespace",
			Finalizers: []string{customresource.FinalizerLabel},
		},
		Spec: akov2.AtlasBackupScheduleSpec{
			PolicyRef: common.ResourceRefNamespaced{Name: policy.Name, Namespace: policy.Namespace},
		},
	}
}

func testBackupPolicy() *akov2.AtlasBackupPolicy {
	return &akov2.AtlasBackupPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-backup-policy",
			Namespace:  "test-namespace",
			Finalizers: []string{customresource.FinalizerLabel},
		},
		Spec: akov2.AtlasBackupPolicySpec{
			Items: []akov2.AtlasBackupPolicyItem{
				{
					FrequencyType:     "weekly",
					FrequencyInterval: 1,
					RetentionUnit:     "days",
					RetentionValue:    7,
				},
			},
		},
	}
}

func TestRegularClusterReconciliation(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-secret",
			Namespace: "default",
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"orgId":         []byte("1234567890"),
			"publicApiKey":  []byte("a1b2c3"),
			"privateApiKey": []byte("abcdef123456"),
		},
		Type: "Opaque",
	}
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-project",
			Namespace: "default",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "MyProject",
			ConnectionSecret: &common.ResourceRefNamespaced{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
		},
		Status: status.AtlasProjectStatus{ID: "abc123"},
	}
	bPolicy := &akov2.AtlasBackupPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-policy",
			Namespace: project.Namespace,
		},
		Spec: akov2.AtlasBackupPolicySpec{
			Items: []akov2.AtlasBackupPolicyItem{
				{
					FrequencyType:     "days",
					FrequencyInterval: 1,
					RetentionUnit:     "weekly",
					RetentionValue:    1,
				},
			},
		},
	}
	bSchedule := &akov2.AtlasBackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-schedule",
			Namespace: project.Namespace,
		},
		Spec: akov2.AtlasBackupScheduleSpec{
			PolicyRef: common.ResourceRefNamespaced{
				Name:      bPolicy.Name,
				Namespace: bPolicy.Namespace,
			},
			ReferenceHourOfDay:    20,
			ReferenceMinuteOfHour: 30,
			RestoreWindowDays:     7,
		},
	}
	searchNodes := []akov2.SearchNode{
		{
			InstanceSize: "S100_LOWCPU_NVME",
			NodeCount:    4,
		},
	}
	d := akov2.DefaultAwsAdvancedDeployment(project.Namespace, project.Name)
	d.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
	d.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
		Name:      bSchedule.Name,
		Namespace: bSchedule.Namespace,
	}
	d.Spec.DeploymentSpec.SearchNodes = searchNodes

	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	sch := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(sch))
	require.NoError(t, corev1.AddToScheme(sch))
	dbUserProjectIndexer := indexer.NewAtlasDatabaseUserByProjectIndexer(ctx, nil, logger)
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secret, project, bPolicy, bSchedule, d).
		WithStatusSubresource(bPolicy, bSchedule).
		WithIndex(dbUserProjectIndexer.Object(), dbUserProjectIndexer.Name(), dbUserProjectIndexer.Keys).
		Build()

	orgID := "0987654321"
	atlasProvider := &atlasmock.TestProvider{
		SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
			return &atlas.ClientSet{SdkClient20241113001: &adminv20241113001.APIClient{}}, "", nil
		},
		SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
			clusterAPI := mockadmin.NewClustersApi(t)
			clusterAPI.EXPECT().GetCluster(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
			clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
				Return(
					&admin.AdvancedClusterDescription{
						GroupId:       pointer.MakePtr(project.ID()),
						Name:          pointer.MakePtr(d.GetDeploymentName()),
						ClusterType:   pointer.MakePtr(d.Spec.DeploymentSpec.ClusterType),
						BackupEnabled: pointer.MakePtr(true),
						StateName:     pointer.MakePtr("IDLE"),
						ReplicationSpecs: &[]admin.ReplicationSpec{
							{
								ZoneName:  pointer.MakePtr("Zone 1"),
								NumShards: pointer.MakePtr(1),
								RegionConfigs: &[]admin.CloudRegionConfig{
									{
										ProviderName: pointer.MakePtr("AWS"),
										RegionName:   pointer.MakePtr("US_EAST_1"),
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &admin.HardwareSpec{
											InstanceSize: pointer.MakePtr("M10"),
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
					&http.Response{},
					nil,
				)
			clusterAPI.EXPECT().GetClusterAdvancedConfiguration(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.GetClusterAdvancedConfigurationApiRequest{ApiService: clusterAPI})
			clusterAPI.EXPECT().GetClusterAdvancedConfigurationExecute(mock.AnythingOfType("admin.GetClusterAdvancedConfigurationApiRequest")).
				Return(
					&admin.ClusterDescriptionProcessArgs{},
					&http.Response{},
					nil,
				)

			searchAPI := mockadmin.NewAtlasSearchApi(t)
			searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), project.ID(), d.Spec.DeploymentSpec.Name).
				Return(admin.GetAtlasSearchDeploymentApiRequest{ApiService: searchAPI})
			searchAPI.EXPECT().GetAtlasSearchDeploymentExecute(mock.Anything).
				Return(
					&admin.ApiSearchDeploymentResponse{
						GroupId:   pointer.MakePtr(project.ID()),
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

			globalAPI := mockadmin.NewGlobalClustersApi(t)
			globalAPI.EXPECT().GetManagedNamespace(context.Background(), project.ID(), d.Spec.DeploymentSpec.Name).
				Return(admin.GetManagedNamespaceApiRequest{ApiService: globalAPI})
			globalAPI.EXPECT().GetManagedNamespaceExecute(mock.Anything).
				Return(&admin.GeoSharding{}, nil, nil)

			return &admin.APIClient{
				ClustersApi:            clusterAPI,
				AtlasSearchApi:         searchAPI,
				ServerlessInstancesApi: mockadmin.NewServerlessInstancesApi(t),
				GlobalClustersApi:      globalAPI,
			}, orgID, nil
		},
		ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
			return &mongodbatlas.Client{
				AdvancedClusters: &atlasmock.AdvancedClustersClientMock{
					GetFunc: func(projectID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return &mongodbatlas.AdvancedCluster{
							ID:            "123789",
							Name:          clusterName,
							GroupID:       projectID,
							BackupEnabled: pointer.MakePtr(true),
							ClusterType:   "REPLICASET",
							ReplicationSpecs: []*mongodbatlas.AdvancedReplicationSpec{
								{
									ID:       "789123",
									ZoneName: "Zone 1",
									RegionConfigs: []*mongodbatlas.AdvancedRegionConfig{
										{
											ProviderName: "AWS",
											RegionName:   "US_EAST_1",
											ElectableSpecs: &mongodbatlas.Specs{
												InstanceSize: "M10",
												NodeCount:    pointer.MakePtr(3),
											},
											Priority: pointer.MakePtr(7),
										},
									},
								},
							},
							StateName: "IDLE",
						}, nil, nil
					},
					DeleteFunc: func(projectID string, clusterName string) (*mongodbatlas.Response, error) {
						return nil, nil
					},
				},
				CloudProviderSnapshotBackupPolicies: &atlasmock.CloudProviderSnapshotBackupPoliciesClientMock{
					GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
						return &mongodbatlas.CloudProviderSnapshotBackupPolicy{
							ClusterID:             "123789",
							ClusterName:           d.GetDeploymentName(),
							ReferenceHourOfDay:    pointer.MakePtr(int64(20)),
							ReferenceMinuteOfHour: pointer.MakePtr(int64(30)),
							RestoreWindowDays:     pointer.MakePtr(int64(7)),
							Policies: []mongodbatlas.Policy{
								{
									ID: "456987",
									PolicyItems: []mongodbatlas.PolicyItem{
										{
											ID:                "987654",
											FrequencyInterval: 1,
											FrequencyType:     "days",
											RetentionUnit:     "weekly",
											RetentionValue:    1,
										},
									},
								},
							},
							AutoExportEnabled:                 pointer.MakePtr(false),
							UseOrgAndGroupNamesInExportPrefix: pointer.MakePtr(false),
						}, nil, nil
					},
				},
			}, orgID, nil
		},
		IsCloudGovFunc: func() bool {
			return false
		},
		IsSupportedFunc: func() bool {
			return true
		},
	}

	reconciler := &AtlasDeploymentReconciler{
		Client:                      k8sClient,
		Log:                         logger.Sugar(),
		AtlasProvider:               atlasProvider,
		EventRecorder:               record.NewFakeRecorder(10),
		ObjectDeletionProtection:    false,
		SubObjectDeletionProtection: false,
	}

	t.Run("should reconcile with existing cluster", func(t *testing.T) {
		result, err := reconciler.Reconcile(
			ctx,
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: d.Namespace,
					Name:      d.Name,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: 0}, result)
	})
}

func TestServerlessInstanceReconciliation(t *testing.T) {
	ctx := context.Background()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-secret",
			Namespace: "default",
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"orgId":         []byte("1234567890"),
			"publicApiKey":  []byte("a1b2c3"),
			"privateApiKey": []byte("abcdef123456"),
		},
		Type: "Opaque",
	}
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-project",
			Namespace: "default",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "MyProject",
			ConnectionSecret: &common.ResourceRefNamespaced{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
		},
		Status: status.AtlasProjectStatus{ID: "abc123"},
	}
	d := akov2.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)

	logger := zaptest.NewLogger(t)

	sch := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(sch))
	require.NoError(t, corev1.AddToScheme(sch))
	dbUserProjectIndexer := indexer.NewAtlasDatabaseUserByProjectIndexer(ctx, nil, logger)
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secret, project, d).
		WithIndex(dbUserProjectIndexer.Object(), dbUserProjectIndexer.Name(), dbUserProjectIndexer.Keys).
		Build()

	orgID := "0987654321"
	atlasProvider := &atlasmock.TestProvider{
		SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
			return &atlas.ClientSet{SdkClient20241113001: &adminv20241113001.APIClient{}}, "", nil
		},
		SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
			err := &admin.GenericOpenAPIError{}
			err.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr(atlas.ServerlessInstanceFromClusterAPI)})
			clusterAPI := mockadmin.NewClustersApi(t)
			clusterAPI.EXPECT().GetCluster(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
			clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
				Return(nil, nil, err)

			serverlessAPI := mockadmin.NewServerlessInstancesApi(t)
			serverlessAPI.EXPECT().GetServerlessInstance(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessAPI})
			serverlessAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
				Return(
					&admin.ServerlessInstanceDescription{
						GroupId: pointer.MakePtr(project.ID()),
						Name:    pointer.MakePtr(d.GetDeploymentName()),
						ProviderSettings: admin.ServerlessProviderSettings{
							BackingProviderName: "AWS",
							ProviderName:        pointer.MakePtr("SERVERLESS"),
							RegionName:          "US_EAST_1",
						},
						ServerlessBackupOptions: &admin.ClusterServerlessBackupOptions{
							ServerlessContinuousBackupEnabled: pointer.MakePtr(false),
						},
						StateName:                    pointer.MakePtr("IDLE"),
						TerminationProtectionEnabled: pointer.MakePtr(false),
					},
					&http.Response{},
					nil,
				)

			speClient := mockadmin.NewServerlessPrivateEndpointsApi(t)
			speClient.EXPECT().ListServerlessPrivateEndpoints(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.ListServerlessPrivateEndpointsApiRequest{ApiService: speClient})
			speClient.EXPECT().ListServerlessPrivateEndpointsExecute(mock.AnythingOfType("admin.ListServerlessPrivateEndpointsApiRequest")).
				Return(nil, &http.Response{}, nil)

			return &admin.APIClient{
				ClustersApi:                   clusterAPI,
				ServerlessInstancesApi:        serverlessAPI,
				ServerlessPrivateEndpointsApi: speClient,
			}, orgID, nil
		},
		ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
			return &mongodbatlas.Client{}, orgID, nil
		},
		IsCloudGovFunc: func() bool {
			return false
		},
		IsSupportedFunc: func() bool {
			return true
		},
	}

	reconciler := &AtlasDeploymentReconciler{
		Client:                      k8sClient,
		Log:                         logger.Sugar(),
		AtlasProvider:               atlasProvider,
		EventRecorder:               record.NewFakeRecorder(10),
		ObjectDeletionProtection:    false,
		SubObjectDeletionProtection: false,
	}

	t.Run("should reconcile with existing serverless instance", func(t *testing.T) {
		result, err := reconciler.Reconcile(
			ctx,
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: d.Namespace,
					Name:      d.Name,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: 0}, result)
	})
}

func TestDeletionReconciliation(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-secret",
			Namespace: "default",
			Labels: map[string]string{
				"atlas.mongodb.com/type": "credentials",
			},
		},
		Data: map[string][]byte{
			"orgId":         []byte("1234567890"),
			"publicApiKey":  []byte("a1b2c3"),
			"privateApiKey": []byte("abcdef123456"),
		},
		Type: "Opaque",
	}
	project := &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-project",
			Namespace: "default",
		},
		Spec: akov2.AtlasProjectSpec{
			Name: "MyProject",
			ConnectionSecret: &common.ResourceRefNamespaced{
				Name:      secret.Name,
				Namespace: secret.Namespace,
			},
		},
		Status: status.AtlasProjectStatus{ID: "abc123"},
	}
	bPolicy := &akov2.AtlasBackupPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-policy",
			Namespace: project.Namespace,
		},
		Spec: akov2.AtlasBackupPolicySpec{
			Items: []akov2.AtlasBackupPolicyItem{
				{
					FrequencyType:     "days",
					FrequencyInterval: 1,
					RetentionUnit:     "weekly",
					RetentionValue:    1,
				},
			},
		},
	}
	bSchedule := &akov2.AtlasBackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-schedule",
			Namespace: project.Namespace,
		},
		Spec: akov2.AtlasBackupScheduleSpec{
			PolicyRef: common.ResourceRefNamespaced{
				Name:      bPolicy.Name,
				Namespace: bPolicy.Namespace,
			},
			ReferenceHourOfDay:    20,
			ReferenceMinuteOfHour: 30,
			RestoreWindowDays:     7,
		},
	}
	searchNodes := []akov2.SearchNode{
		{
			InstanceSize: "S100_LOWCPU_NVME",
			NodeCount:    4,
		},
	}
	d := akov2.DefaultAwsAdvancedDeployment(project.Namespace, project.Name)
	d.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
	d.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
		Name:      bSchedule.Name,
		Namespace: bSchedule.Namespace,
	}
	d.Spec.DeploymentSpec.SearchNodes = searchNodes
	d.Finalizers = []string{customresource.FinalizerLabel}

	sch := runtime.NewScheme()
	require.NoError(t, akov2.AddToScheme(sch))
	require.NoError(t, corev1.AddToScheme(sch))
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secret, project, bPolicy, bSchedule, d).
		WithStatusSubresource(bPolicy, bSchedule, d).
		Build()

	orgID := "0987654321"
	logger := zaptest.NewLogger(t).Sugar()
	atlasProvider := &atlasmock.TestProvider{
		SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
			return &atlas.ClientSet{SdkClient20241113001: &adminv20241113001.APIClient{}}, "", nil
		},
		SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
			clusterAPI := mockadmin.NewClustersApi(t)
			clusterAPI.EXPECT().GetCluster(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
			clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
				Return(
					&admin.AdvancedClusterDescription{
						GroupId:       pointer.MakePtr(project.ID()),
						Name:          pointer.MakePtr(d.GetDeploymentName()),
						ClusterType:   pointer.MakePtr(d.Spec.DeploymentSpec.ClusterType),
						BackupEnabled: pointer.MakePtr(true),
						StateName:     pointer.MakePtr("IDLE"),
						ReplicationSpecs: &[]admin.ReplicationSpec{
							{
								ZoneName:  pointer.MakePtr("Zone 1"),
								NumShards: pointer.MakePtr(1),
								RegionConfigs: &[]admin.CloudRegionConfig{
									{
										ProviderName: pointer.MakePtr("AWS"),
										RegionName:   pointer.MakePtr("US_EAST_1"),
										Priority:     pointer.MakePtr(7),
										ElectableSpecs: &admin.HardwareSpec{
											InstanceSize: pointer.MakePtr("M10"),
											NodeCount:    pointer.MakePtr(3),
										},
									},
								},
							},
						},
					},
					&http.Response{},
					nil,
				)
			clusterAPI.EXPECT().DeleteCluster(context.Background(), project.ID(), d.GetDeploymentName()).
				Return(admin.DeleteClusterApiRequest{ApiService: clusterAPI})
			clusterAPI.EXPECT().DeleteClusterExecute(mock.AnythingOfType("admin.DeleteClusterApiRequest")).
				Return(&http.Response{}, nil)

			return &admin.APIClient{
				ClustersApi:            clusterAPI,
				ServerlessInstancesApi: mockadmin.NewServerlessInstancesApi(t),
			}, orgID, nil
		},
		ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
			return &mongodbatlas.Client{}, orgID, nil
		},
		IsCloudGovFunc: func() bool {
			return false
		},
		IsSupportedFunc: func() bool {
			return true
		},
	}

	reconciler := &AtlasDeploymentReconciler{
		Client:                      k8sClient,
		Log:                         logger,
		AtlasProvider:               atlasProvider,
		EventRecorder:               record.NewFakeRecorder(10),
		ObjectDeletionProtection:    false,
		SubObjectDeletionProtection: false,
	}

	t.Run("should reconcile deletion of existing cluster", func(t *testing.T) {
		require.NoError(t, k8sClient.Delete(context.Background(), d))
		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: d.Namespace,
					Name:      d.Name,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: 0}, result)
	})
}

func TestFindDeploymentsForSearchIndexConfig(t *testing.T) {
	t.Run("should fail when watching wrong object", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		reconciler := &AtlasDeploymentReconciler{
			Log: zap.New(core).Sugar(),
		}

		assert.Nil(t, reconciler.findDeploymentsForSearchIndexConfig(context.Background(), &akov2.AtlasProject{}))
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zap.WarnLevel, logs.All()[0].Level)
		assert.Equal(t, "watching AtlasSearchIndexConfig but got *v1.AtlasProject", logs.All()[0].Message)
	})

	t.Run("should return slice of requests for instances", func(t *testing.T) {
		instance1 := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance1",
				Namespace: "default",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					SearchIndexes: []akov2.SearchIndex{
						{
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      "index1",
									Namespace: "default",
								},
							},
						},
					},
				},
			},
		}
		instance2 := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "instance2",
				Namespace: "other-ns",
			},
			Spec: akov2.AtlasDeploymentSpec{
				DeploymentSpec: &akov2.AdvancedDeploymentSpec{
					SearchIndexes: []akov2.SearchIndex{
						{
							Search: &akov2.Search{
								SearchConfigurationRef: common.ResourceRefNamespaced{
									Name:      "index1",
									Namespace: "default",
								},
							},
						},
					},
				},
			},
		}
		connection := &akov2.AtlasSearchIndexConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "index1",
				Namespace: "default",
			},
		}
		testScheme := runtime.NewScheme()
		assert.NoError(t, akov2.AddToScheme(testScheme))
		deploymentIndexer := indexer.NewAtlasDeploymentBySearchIndexIndexer(zaptest.NewLogger(t))
		k8sClient := fake.NewClientBuilder().
			WithScheme(testScheme).
			WithObjects(connection, instance1, instance2).
			WithIndex(
				deploymentIndexer.Object(),
				deploymentIndexer.Name(),
				deploymentIndexer.Keys,
			).
			Build()
		reconciler := &AtlasDeploymentReconciler{
			Client: k8sClient,
			Log:    zaptest.NewLogger(t).Sugar(),
		}

		requests := reconciler.findDeploymentsForSearchIndexConfig(context.Background(), connection)
		assert.Equal(
			t,
			[]ctrl.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: "default",
						Name:      "instance1",
					},
				},
				{
					NamespacedName: types.NamespacedName{
						Namespace: "other-ns",
						Name:      "instance2",
					},
				},
			},
			requests,
		)
	})
}

func TestFindDeploymentsForBackupPolicy(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasProject{},
			want: nil,
		},
		{
			name: "transitive dependency",
			obj: &akov2.AtlasBackupPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "some-policy", Namespace: "ns1"},
			},
			initObjs: []client.Object{
				&akov2.AtlasBackupSchedule{
					ObjectMeta: metav1.ObjectMeta{Name: "some-schedule", Namespace: "ns2"},
					Spec: akov2.AtlasBackupScheduleSpec{
						PolicyRef: common.ResourceRefNamespaced{Name: "some-policy", Namespace: "ns1"},
					},
				},
				&akov2.AtlasBackupSchedule{
					ObjectMeta: metav1.ObjectMeta{Name: "some-schedule2", Namespace: "ns2"},
					Spec: akov2.AtlasBackupScheduleSpec{
						PolicyRef: common.ResourceRefNamespaced{Name: "some-policy", Namespace: "ns1"},
					},
				},
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment", Namespace: "ns3"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule", Namespace: "ns2"},
					},
				},
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment2", Namespace: "ns4"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule", Namespace: "ns2"},
					},
				},
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment3", Namespace: "ns5"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule2", Namespace: "ns2"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "some-deployment", Namespace: "ns3"}},
				{NamespacedName: types.NamespacedName{Name: "some-deployment2", Namespace: "ns4"}},
				{NamespacedName: types.NamespacedName{Name: "some-deployment3", Namespace: "ns5"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			backupScheduleIndexer := indexer.NewAtlasBackupScheduleByBackupPolicyIndexer(zaptest.NewLogger(t))
			deploymentIndexer := indexer.NewAtlasDeploymentByBackupScheduleIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(backupScheduleIndexer.Object(), backupScheduleIndexer.Name(), backupScheduleIndexer.Keys).
				WithIndex(deploymentIndexer.Object(), deploymentIndexer.Name(), deploymentIndexer.Keys).
				Build()
			reconciler := &AtlasDeploymentReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findDeploymentsForBackupPolicy(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}

func TestFindDeploymentsForBackupSchedule(t *testing.T) {
	for _, tc := range []struct {
		name     string
		obj      client.Object
		initObjs []client.Object
		want     []reconcile.Request
	}{
		{
			name: "wrong type",
			obj:  &akov2.AtlasProject{},
			want: nil,
		},
		{
			name: "transitive dependency",
			obj: &akov2.AtlasBackupSchedule{
				ObjectMeta: metav1.ObjectMeta{Name: "some-schedule", Namespace: "ns2"},
			},
			initObjs: []client.Object{
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment", Namespace: "ns3"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule", Namespace: "ns2"},
					},
				},
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment2", Namespace: "ns4"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule", Namespace: "ns2"},
					},
				},
				&akov2.AtlasDeployment{
					ObjectMeta: metav1.ObjectMeta{Name: "some-deployment3", Namespace: "ns5"},
					Spec: akov2.AtlasDeploymentSpec{
						BackupScheduleRef: common.ResourceRefNamespaced{Name: "some-schedule2", Namespace: "ns2"},
					},
				},
			},
			want: []reconcile.Request{
				{NamespacedName: types.NamespacedName{Name: "some-deployment", Namespace: "ns3"}},
				{NamespacedName: types.NamespacedName{Name: "some-deployment2", Namespace: "ns4"}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			backupScheduleIndexer := indexer.NewAtlasBackupScheduleByBackupPolicyIndexer(zaptest.NewLogger(t))
			deploymentIndexer := indexer.NewAtlasDeploymentByBackupScheduleIndexer(zaptest.NewLogger(t))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tc.initObjs...).
				WithIndex(backupScheduleIndexer.Object(), backupScheduleIndexer.Name(), backupScheduleIndexer.Keys).
				WithIndex(deploymentIndexer.Object(), deploymentIndexer.Name(), deploymentIndexer.Keys).
				Build()
			reconciler := &AtlasDeploymentReconciler{
				Log:    zaptest.NewLogger(t).Sugar(),
				Client: k8sClient,
			}
			got := reconciler.findDeploymentsForBackupSchedule(context.Background(), tc.obj)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("want reconcile requests: %v, got %v", got, tc.want)
			}
		})
	}
}

func TestGetProjectFromAtlas(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment  *akov2.AtlasDeployment
		deploymentSecret *corev1.Secret
		atlasProvider    atlas.Provider
		expectedErr      error
	}{
		"failed to create atlas client": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "project-id",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "project-creds",
						},
					},
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to create client")
				},
			},
			expectedErr: errors.New("failed to create client"),
		},
		"failed to get project from atlas": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "project-id",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "project-creds",
						},
					},
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return false
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{}, "", nil
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProject(context.Background(), "project-id").
						Return(admin.GetProjectApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(nil, nil, errors.New("failed to get project"))

					return &admin.APIClient{ProjectsApi: projectAPI}, "", nil
				},
			},
			expectedErr: errors.New("failed to get project"),
		},
		"get project": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					ExternalProjectRef: &akov2.ExternalProjectReference{
						ID: "project-id",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "project-creds",
						},
					},
				},
			},
			deploymentSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user-pass",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"password": []byte("Passw0rd!"),
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return false
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{}, "", nil
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					projectAPI := mockadmin.NewProjectsApi(t)
					projectAPI.EXPECT().GetProject(context.Background(), "project-id").
						Return(admin.GetProjectApiRequest{ApiService: projectAPI})
					projectAPI.EXPECT().GetProjectExecute(mock.AnythingOfType("admin.GetProjectApiRequest")).
						Return(&admin.Group{Id: pointer.MakePtr("project-id")}, nil, nil)

					return &admin.APIClient{ProjectsApi: projectAPI}, "", nil
				},
				ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
					return &mongodbatlas.Client{}, "", nil
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasDeployment).
				WithStatusSubresource(tt.atlasDeployment)

			if tt.deploymentSecret != nil {
				k8sClient.WithObjects(tt.deploymentSecret)
			}

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDeploymentReconciler{
				Client:        k8sClient.Build(),
				AtlasProvider: tt.atlasProvider,
				Log:           logger,
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			_, err := r.getProjectFromAtlas(ctx, tt.atlasDeployment)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestGetProjectFromKube(t *testing.T) {
	tests := map[string]struct {
		atlasDeployment *akov2.AtlasDeployment
		project         *akov2.AtlasProject
		projectSecret   *corev1.Secret
		atlasProvider   atlas.Provider
		expectedErr     error
	}{
		"failed to get project": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "2.4.1",
					},
				},
				Spec: akov2.AtlasDeploymentSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
			},
			expectedErr: &k8serrors.StatusError{
				ErrStatus: metav1.Status{
					Status:  "Failure",
					Message: "atlasprojects.atlas.mongodb.com \"my-project\" not found",
					Reason:  "NotFound",
					Code:    404,
					Details: &metav1.StatusDetails{
						Group: "atlas.mongodb.com",
						Kind:  "atlasprojects",
						Name:  "my-project",
					},
				},
			},
		},
		"failed to create atlas sdk": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
					Labels: map[string]string{
						"mongodb.com/atlas-resource-version": "2.4.1",
					},
				},
				Spec: akov2.AtlasDeploymentSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name: "project-secret",
					},
				},
			},
			projectSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "project-secret",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"publicKey":  []byte("publicKey"),
					"privateKey": []byte("privateKey"),
					"orgID":      []byte("orgID"),
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsSupportedFunc: func() bool {
					return true
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return nil, "", errors.New("failed to create client")
				},
			},
			expectedErr: errors.New("failed to create client"),
		},
		"get project": {
			atlasDeployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "user1",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "project-creds",
						},
					},
				},
			},
			project: &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "my-project",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name: "project-secret",
					},
				},
			},
			projectSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "project-secret",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"publicKey":  []byte("publicKey"),
					"privateKey": []byte("privateKey"),
					"orgID":      []byte("orgID"),
				},
			},
			atlasProvider: &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return false
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{}, "", nil
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					return &admin.APIClient{}, "", nil
				},
				ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
					return &mongodbatlas.Client{}, "", nil
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testScheme := runtime.NewScheme()
			assert.NoError(t, akov2.AddToScheme(testScheme))
			assert.NoError(t, corev1.AddToScheme(testScheme))
			k8sClient := fake.NewClientBuilder().
				WithScheme(testScheme).
				WithObjects(tt.atlasDeployment).
				WithStatusSubresource(tt.atlasDeployment)

			if tt.project != nil {
				k8sClient.WithObjects(tt.project)
			}

			if tt.projectSecret != nil {
				k8sClient.WithObjects(tt.projectSecret)
			}

			logger := zaptest.NewLogger(t).Sugar()
			r := AtlasDeploymentReconciler{
				Client:        k8sClient.Build(),
				AtlasProvider: tt.atlasProvider,
				Log:           logger,
			}
			ctx := &workflow.Context{
				Context: context.Background(),
				Log:     logger,
			}

			_, err := r.getProjectFromKube(ctx, tt.atlasDeployment)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestChangeDeploymentType(t *testing.T) {
	tests := map[string]struct {
		deployment *akov2.AtlasDeployment
	}{
		"should fail when existing cluster is regular but manifest defines a serverless instance": {
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					ServerlessSpec: &akov2.ServerlessSpec{
						Name: "cluster0",
						ProviderSettings: &akov2.ServerlessProviderSettingsSpec{
							ProviderName:        "SERVERLESS",
							BackingProviderName: "AWS",
						},
					},
				},
				Status: status.AtlasDeploymentStatus{
					StateName: "IDLE",
				},
			},
		},
		"should fail when existing cluster is serverless instance but manifest defines a regular deployment": {
			deployment: &akov2.AtlasDeployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cluster0",
					Namespace: "default",
				},
				Spec: akov2.AtlasDeploymentSpec{
					Project: &common.ResourceRefNamespaced{
						Name:      "my-project",
						Namespace: "default",
					},
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						Name: "cluster0",
					},
				},
				Status: status.AtlasDeploymentStatus{
					StateName: "IDLE",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "api-secret",
					Namespace: "default",
					Labels: map[string]string{
						"atlas.mongodb.com/type": "credentials",
					},
				},
				Data: map[string][]byte{
					"orgId":         []byte("1234567890"),
					"publicApiKey":  []byte("a1b2c3"),
					"privateApiKey": []byte("abcdef123456"),
				},
				Type: "Opaque",
			}
			project := &akov2.AtlasProject{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-project",
					Namespace: "default",
				},
				Spec: akov2.AtlasProjectSpec{
					Name: "MyProject",
					ConnectionSecret: &common.ResourceRefNamespaced{
						Name:      secret.Name,
						Namespace: secret.Namespace,
					},
				},
				Status: status.AtlasProjectStatus{ID: "abc123"},
			}

			ctx := context.Background()
			logger := zaptest.NewLogger(t)

			sch := runtime.NewScheme()
			require.NoError(t, akov2.AddToScheme(sch))
			require.NoError(t, corev1.AddToScheme(sch))
			dbUserProjectIndexer := indexer.NewAtlasDatabaseUserByProjectIndexer(ctx, nil, logger)
			k8sClient := fake.NewClientBuilder().
				WithScheme(sch).
				WithObjects(secret, project, tt.deployment).
				WithStatusSubresource(project, tt.deployment).
				WithIndex(dbUserProjectIndexer.Object(), dbUserProjectIndexer.Name(), dbUserProjectIndexer.Keys).
				Build()

			atlasProvider := &atlasmock.TestProvider{
				IsCloudGovFunc: func() bool {
					return false
				},
				IsSupportedFunc: func() bool {
					return true
				},
				SdkSetClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*atlas.ClientSet, string, error) {
					return &atlas.ClientSet{}, "", nil
				},
				ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
					return &mongodbatlas.Client{}, "org-id", nil
				},
				SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
					clusterAPI := mockadmin.NewClustersApi(t)
					clusterAPI.EXPECT().GetCluster(ctx, "abc123", "cluster0").
						Return(admin.GetClusterApiRequest{ApiService: clusterAPI})
					clusterAPI.EXPECT().GetClusterExecute(mock.AnythingOfType("admin.GetClusterApiRequest")).
						RunAndReturn(
							func(request admin.GetClusterApiRequest) (*admin.AdvancedClusterDescription, *http.Response, error) {
								if !tt.deployment.IsServerless() {
									err := &admin.GenericOpenAPIError{}
									err.SetModel(admin.ApiError{ErrorCode: pointer.MakePtr(atlas.ServerlessInstanceFromClusterAPI)})
									return nil, nil, err
								}
								return &admin.AdvancedClusterDescription{Name: pointer.MakePtr("cluster0")}, nil, nil
							},
						)

					serverlessAPI := mockadmin.NewServerlessInstancesApi(t)
					if !tt.deployment.IsServerless() {
						serverlessAPI.EXPECT().GetServerlessInstance(ctx, "abc123", "cluster0").
							Return(admin.GetServerlessInstanceApiRequest{ApiService: serverlessAPI})
						serverlessAPI.EXPECT().GetServerlessInstanceExecute(mock.AnythingOfType("admin.GetServerlessInstanceApiRequest")).
							Return(&admin.ServerlessInstanceDescription{Name: pointer.MakePtr("cluster0")}, nil, nil)
					}

					return &admin.APIClient{ClustersApi: clusterAPI, ServerlessInstancesApi: serverlessAPI}, "org-id", nil
				},
			}

			r := &AtlasDeploymentReconciler{
				Client:        k8sClient,
				AtlasProvider: atlasProvider,
				Log:           logger.Sugar(),
				EventRecorder: record.NewFakeRecorder(1),
			}
			result, err := r.Reconcile(
				ctx,
				ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: tt.deployment.Namespace,
						Name:      tt.deployment.Name,
					},
				},
			)

			assert.NoError(t, err)
			assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: workflow.DefaultRetry}, result)
			assert.NoError(t, k8sClient.Get(ctx, client.ObjectKeyFromObject(tt.deployment), tt.deployment))
			assert.True(
				t,
				cmp.Equal(
					[]api.Condition{
						api.FalseCondition(api.ReadyType),
						api.TrueCondition(api.ResourceVersionStatus),
						api.TrueCondition(api.ValidationSucceeded),
						api.FalseCondition(api.DeploymentReadyType).
							WithReason(string(workflow.Internal)).
							WithMessageRegexp("regular deployment cannot be converted into a serverless deployment and vice-versa"),
					},
					tt.deployment.Status.Conditions,
					cmpopts.IgnoreFields(api.Condition{}, "LastTransitionTime"),
				),
			)
		})
	}
}
