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
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	atlas_mock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/toptr"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const (
	fakeDomain     = "atlas-unit-test.local"
	fakeProject    = "test-project"
	fakeProjectID  = "fake-test-project-id"
	fakeDeployment = "fake-cluster"
	fakeNamespace  = "fake-namespace"
)

func TestDeploymentManaged(t *testing.T) {
	testCases := []struct {
		title      string
		protected  bool
		managedTag bool
	}{
		{
			title:      "unprotected means managed",
			protected:  false,
			managedTag: false,
		},
		{
			title:      "protected and tagged means managed",
			protected:  true,
			managedTag: true,
		},
		{
			title:      "protected not tagged and missing in Atlas means managed",
			protected:  true,
			managedTag: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &atlas_mock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ClusterNotFound}
					},
				},
				ServerlessInstances: &atlas_mock.ServerlessInstancesClientMock{
					GetFunc: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceNotFound}
					},
				},
			}
			project := testProject(fakeNamespace)
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			te := newTestDeploymentEnv(t, tc.protected, &atlasClient, testK8sClient(), project, deployment)
			if tc.managedTag {
				customresource.SetAnnotation(te.deployment, customresource.AnnotationLastAppliedConfiguration, "")
			}

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.log, te.project, te.deployment)

			assert.True(t, result.IsOk())
		})
	}
}

func TestProtectedAdvancedDeploymentManagedInAtlas(t *testing.T) {
	testCases := []struct {
		title       string
		inAtlas     *mongodbatlas.AdvancedCluster
		expectedErr string
	}{
		{
			title:       "advanced deployment not tagged and same in Atlas STILL means managed",
			inAtlas:     sameAdvancedDeployment(fakeDomain),
			expectedErr: "",
		},
		{
			title:       "advanced deployment not tagged and different in Atlas means unmanaged",
			inAtlas:     differentAdvancedDeployment(fakeDomain),
			expectedErr: "unable to reconcile Deployment due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &atlas_mock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			te := newTestDeploymentEnv(t, true, &atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.log, te.project, te.deployment)

			if tc.expectedErr == "" {
				assert.True(t, result.IsOk())
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedErr), result.GetMessage())
			}
		})
	}
}

func TestProtectedServerlessManagedInAtlas(t *testing.T) {
	testCases := []struct {
		title       string
		inAtlas     *mongodbatlas.Cluster
		expectedErr string
	}{
		{
			title:       "serverless deployment not tagged and same in Atlas STILL means managed",
			inAtlas:     sameServerlessDeployment(fakeDomain),
			expectedErr: "",
		},
		{
			title:       "serverless deployment not tagged and different in Atlas means unmanaged",
			inAtlas:     differentServerlessDeployment(fakeDomain),
			expectedErr: "unable to reconcile Deployment due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &atlas_mock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceFromClusterAPI}
					},
				},
				ServerlessInstances: &atlas_mock.ServerlessInstancesClientMock{
					GetFunc: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := v1.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
			te := newTestDeploymentEnv(t, true, &atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.log, te.project, te.deployment)

			if tc.expectedErr == "" {
				assert.True(t, result.IsOk())
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedErr), result.GetMessage())
			}
		})
	}
}

func TestFinalizerNotFound(t *testing.T) {
	atlasClient := mongodbatlas.Client{}
	project := testProject(fakeNamespace)
	deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	k8sclient := testK8sClient()
	te := newTestDeploymentEnv(t, false, &atlasClient, k8sclient, project, deployment)

	deletionRequest, result := te.reconciler.handleDeletion(
		te.workflowCtx,
		te.log,
		te.prevResult,
		te.project,
		te.deployment,
	)

	require.True(t, deletionRequest)
	assert.Regexp(t, regexp.MustCompile("not found"), result.GetMessage())
}

func TestFinalizerGetsSet(t *testing.T) {
	testCases := []struct {
		title         string
		haveFinalizer bool
	}{
		{
			title:         "with a finalizer, it remains set",
			haveFinalizer: true,
		},
		{
			title:         "without a finalizer, it gets set",
			haveFinalizer: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			atlasClient := mongodbatlas.Client{}
			project := testProject(fakeNamespace)
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			if tc.haveFinalizer {
				customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			}
			k8sclient := testK8sClient()
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			te := newTestDeploymentEnv(t, false, &atlasClient, k8sclient, project, deployment)

			deletionRequest, _ := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.False(t, deletionRequest)
			finalDeployment := &v1.AtlasDeployment{}
			require.NoError(t, te.reconciler.Client.Get(context.Background(), client.ObjectKeyFromObject(te.deployment), finalDeployment))
			assert.True(t, customresource.HaveFinalizer(finalDeployment, customresource.FinalizerLabel))
		})
	}
}

func TestDeploymentDeletionProtection(t *testing.T) {
	testCases := []struct {
		title         string
		protected     bool
		expectRemoval int
	}{
		{
			title:         "Deployment with protection ON and no annotations is kept",
			protected:     true,
			expectRemoval: 0,
		},
		{
			title:         "Deployment with protection OFF and no annotations is removed",
			protected:     false,
			expectRemoval: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			advancedClusterClient := &atlas_mock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			// set deletion timestamp after creation in k8s
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			te := newTestDeploymentEnv(t, tc.protected, &atlasClient, k8sclient, project, deployment)
			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Len(t, advancedClusterClient.DeleteRequests, tc.expectRemoval)
		})
	}
}

func TestKeepAnnotatedDeploymentAlwaysRemain(t *testing.T) {
	testCases := []struct {
		title     string
		protected bool
	}{
		{
			title:     "Deployment with protection ON and 'keep' annotation is kept",
			protected: true,
		},
		{
			title:     "Deployment with protection OFF but 'keep' annotation is kept",
			protected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			advancedClusterClient := &atlas_mock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			customresource.SetAnnotation(deployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyKeep,
			)
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			// set deletion timestamp after creation in k8s, otherwise the creation would reset the deletion timestamp
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			te := newTestDeploymentEnv(t, tc.protected, &atlasClient, k8sclient, project, deployment)

			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Len(t, advancedClusterClient.DeleteRequests, 0)
		})
	}
}

func TestDeleteAnnotatedDeploymentGetRemoved(t *testing.T) {
	testCases := []struct {
		title     string
		protected bool
	}{
		{
			title:     "Deployment with protection ON but 'delete' annotation is removed",
			protected: true,
		},
		{
			title:     "Deployment with protection OFF and 'delete' annotation is removed",
			protected: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			advancedClusterClient := &atlas_mock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			customresource.SetAnnotation(deployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyDelete,
			)
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			// set deletion timestamp after creation in k8s, otherwise the creation would reset the deletion timestamp
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			te := newTestDeploymentEnv(t, tc.protected, &atlasClient, k8sclient, project, deployment)

			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Len(t, advancedClusterClient.DeleteRequests, 1)
		})
	}
}

func TestCleanupBindings(t *testing.T) {
	t.Run("without backup references, nothing happens on cleanup", func(t *testing.T) {
		r := &AtlasDeploymentReconciler{
			Log:    testLog(t),
			Client: testK8sClient(),
		}
		d := &v1.AtlasDeployment{} // dummy deployment

		// test cleanup
		assert.NoError(t, r.cleanupBindings(context.Background(), d))
	})

	t.Run("with unreferenced backups, still nothing happens on cleanup", func(t *testing.T) {
		r := &AtlasDeploymentReconciler{
			Log:    testLog(t),
			Client: testK8sClient(),
		}
		dn := testDeploymentName("") // deployment, schedule, policy (NOT connected)
		deployment := &v1.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: dn.Name, Namespace: dn.Namespace},
		}
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		policy := testBackupPolicy()
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test cleanup
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

		endPolicy := &v1.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.Equal(t, []string{customresource.FinalizerLabel}, endPolicy.Finalizers)
		endSchedule := &v1.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.Equal(t, []string{customresource.FinalizerLabel}, endSchedule.Finalizers)
	})

	t.Run("last deployment's referenced backups finalizers are cleaned up", func(t *testing.T) {
		atlasProvider := &atlas_mock.TestProvider{
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
		deployment := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		schedule.Status.DeploymentIDs = []string{deployment.Spec.DeploymentSpec.Name}
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test ensureBackupPolicy and cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule, &[]watch.WatchedObject{})
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

		endPolicy := &v1.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.Empty(t, endPolicy.Finalizers, "policy should end up with no finalizer")
		endSchedule := &v1.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.Empty(t, endSchedule.Finalizers, "schedule should end up with no finalizer")
	})

	t.Run("referenced backups finalizers are NOT cleaned up if reachable by other deployment", func(t *testing.T) {
		atlasProvider := &atlas_mock.TestProvider{
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
		deployment := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		deployment2 := testDeployment("2", schedule)
		require.NoError(t, r.Client.Create(context.Background(), deployment2))
		schedule.Status.DeploymentIDs = []string{
			deployment.Spec.DeploymentSpec.Name,
			deployment2.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule, &[]watch.WatchedObject{})
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

		endPolicy := &v1.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(policy), endPolicy))
		assert.NotEmpty(t, endPolicy.Finalizers, "policy should keep the finalizer")
		endSchedule := &v1.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKeyFromObject(schedule), endSchedule))
		assert.NotEmpty(t, endSchedule.Finalizers, "schedule should keep the finalizer")
	})

	t.Run("policy finalizer stays if still referenced", func(t *testing.T) {
		atlasProvider := &atlas_mock.TestProvider{
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
		deployment := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		deployment2 := testDeployment("2", schedule2)
		require.NoError(t, r.Client.Create(context.Background(), deployment2))
		schedule.Status.DeploymentIDs = []string{
			deployment.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule))
		schedule2.Status.DeploymentIDs = []string{
			deployment2.Spec.DeploymentSpec.Name,
		}
		require.NoError(t, r.Client.Create(context.Background(), schedule2))
		policy.Status.BackupScheduleIDs = []string{
			fmt.Sprintf("%s/%s", schedule.Namespace, schedule.Name),
			fmt.Sprintf("%s/%s", schedule2.Namespace, schedule2.Name),
		}

		// test cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule, &[]watch.WatchedObject{})
		require.NoError(t, err)
		_, err = r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule2, &[]watch.WatchedObject{})
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

		endPolicy := &v1.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(policy.Namespace, policy.Name), endPolicy))
		assert.NotEmpty(t, endPolicy.Finalizers, "policy should keep the finalizer")
		endSchedule := &v1.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(schedule.Namespace, schedule.Name), endSchedule))
		assert.Empty(t, endSchedule.Finalizers, "schedule should end up with no finalizer")
	})
}

func differentAdvancedDeployment(ns string) *mongodbatlas.AdvancedCluster {
	project := testProject(ns)
	deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	deployment.Spec.DeploymentSpec.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize = "M2"
	advancedSpec := deployment.Spec.DeploymentSpec
	return intoAdvancedAtlasCluster(advancedSpec)
}

func sameAdvancedDeployment(ns string) *mongodbatlas.AdvancedCluster {
	project := testProject(ns)
	deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	advancedSpec := deployment.Spec.DeploymentSpec
	return intoAdvancedAtlasCluster(advancedSpec)
}

func differentServerlessDeployment(ns string) *mongodbatlas.Cluster {
	project := testProject(ns)
	deployment := v1.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
	deployment.Spec.ServerlessSpec.ProviderSettings.RegionName = "US_EAST_2"
	return intoServerlessAtlasCluster(deployment.Spec.ServerlessSpec)
}

func sameServerlessDeployment(ns string) *mongodbatlas.Cluster {
	project := testProject(ns)
	deployment := v1.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
	return intoServerlessAtlasCluster(deployment.Spec.ServerlessSpec)
}

type testDeploymentEnv struct {
	reconciler  *AtlasDeploymentReconciler
	workflowCtx *workflow.Context
	log         *zap.SugaredLogger
	prevResult  workflow.Result
	project     *v1.AtlasProject
	deployment  *v1.AtlasDeployment
}

func newTestDeploymentEnv(t *testing.T,
	protected bool,
	atlasClient *mongodbatlas.Client,
	k8sclient client.Client,
	project *v1.AtlasProject,
	deployment *v1.AtlasDeployment,
) *testDeploymentEnv {
	t.Helper()

	logger := testLog(t)
	r := testDeploymentReconciler(logger, k8sclient, protected)

	prevResult := testPrevResult()
	workflowCtx := customresource.MarkReconciliationStarted(r.Client, deployment, logger, context.Background())
	workflowCtx.Client = atlasClient
	return &testDeploymentEnv{
		reconciler:  r,
		workflowCtx: workflowCtx,
		log:         r.Log.With("atlasdeployment", "test-namespace"),
		prevResult:  prevResult,
		deployment:  deployment,
		project:     project,
	}
}

func testK8sClient() client.Client {
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	sch := runtime.NewScheme()
	sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasDeployment{})
	sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupSchedule{})
	sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupScheduleList{})
	sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupPolicy{})
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
	return fake.NewClientBuilder().WithScheme(sch).
		WithStatusSubresource(&v1.AtlasBackupSchedule{}, &v1.AtlasBackupPolicy{}).
		Build()
}

func testLog(t *testing.T) *zap.SugaredLogger {
	t.Helper()

	return zaptest.NewLogger(t).Sugar()
}

func testPrevResult() workflow.Result {
	return workflow.Result{}.WithMessage("unchanged")
}

func testDeploymentReconciler(log *zap.SugaredLogger, k8sclient client.Client, protected bool) *AtlasDeploymentReconciler {
	return &AtlasDeploymentReconciler{
		Client:                   k8sclient,
		Log:                      log,
		ObjectDeletionProtection: protected,
	}
}

func testProject(ns string) *v1.AtlasProject {
	return &v1.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeProject,
			Namespace: ns,
		},
		Status: status.AtlasProjectStatus{
			ID: fakeProjectID,
		},
	}
}

func intoAdvancedAtlasCluster(advancedSpec *v1.AdvancedDeploymentSpec) *mongodbatlas.AdvancedCluster {
	ac, err := advancedSpec.ToAtlas()
	if err != nil {
		log.Fatalf("failed to convert advanced deployment to atlas: %v", err)
	}
	return ac
}

func intoServerlessAtlasCluster(serverlessSpec *v1.ServerlessSpec) *mongodbatlas.Cluster {
	ac, err := serverlessSpec.ToAtlas()
	if err != nil {
		log.Fatalf("failed to convert serverless deployment to atlas: %v", err)
	}
	return ac
}

func testDeploymentName(suffix string) types.NamespacedName {
	return types.NamespacedName{
		Name:      fmt.Sprintf("test-deployment%s", suffix),
		Namespace: "test-namespace",
	}
}

func testDeployment(suffix string, schedule *v1.AtlasBackupSchedule) *v1.AtlasDeployment {
	dn := testDeploymentName(suffix)
	return &v1.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: dn.Name, Namespace: dn.Namespace},
		Spec: v1.AtlasDeploymentSpec{
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Name: fmt.Sprintf("atlas-%s", dn.Name),
			},
			BackupScheduleRef: common.ResourceRefNamespaced{
				Name:      schedule.Name,
				Namespace: schedule.Namespace,
			},
		},
	}
}

func testBackupSchedule(suffix string, policy *v1.AtlasBackupPolicy) *v1.AtlasBackupSchedule {
	return &v1.AtlasBackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:       fmt.Sprintf("test-backup-schedule%s", suffix),
			Namespace:  "test-namespace",
			Finalizers: []string{customresource.FinalizerLabel},
		},
		Spec: v1.AtlasBackupScheduleSpec{
			PolicyRef: common.ResourceRefNamespaced{Name: policy.Name, Namespace: policy.Namespace},
		},
	}
}

func testBackupPolicy() *v1.AtlasBackupPolicy {
	return &v1.AtlasBackupPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "test-backup-policy",
			Namespace:  "test-namespace",
			Finalizers: []string{customresource.FinalizerLabel},
		},
		Spec: v1.AtlasBackupPolicySpec{
			Items: []v1.AtlasBackupPolicyItem{
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

func TestUniqueKey(t *testing.T) {
	t.Run("Test duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &v1.AtlasDeploymentSpec{
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Tags: []*v1.TagSpec{{Key: "foo", Value: "true"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &v1.AtlasDeploymentSpec{
			DeploymentSpec: &v1.AdvancedDeploymentSpec{
				Tags: []*v1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foobar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
	t.Run("Test duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &v1.AtlasDeploymentSpec{
			ServerlessSpec: &v1.ServerlessSpec{
				Tags: []*v1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &v1.AtlasDeploymentSpec{
			ServerlessSpec: &v1.ServerlessSpec{
				Tags: []*v1.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
}

func TestReconciliation(t *testing.T) {
	t.Run("should reconcile with existing cluster", func(t *testing.T) {
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
		project := &v1.AtlasProject{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-project",
				Namespace: "default",
			},
			Spec: v1.AtlasProjectSpec{
				Name: "MyProject",
				ConnectionSecret: &common.ResourceRefNamespaced{
					Name:      secret.Name,
					Namespace: secret.Namespace,
				},
			},
			Status: status.AtlasProjectStatus{ID: "abc123"},
		}
		bPolicy := &v1.AtlasBackupPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-policy",
				Namespace: project.Namespace,
			},
			Spec: v1.AtlasBackupPolicySpec{
				Items: []v1.AtlasBackupPolicyItem{
					{
						FrequencyType:     "days",
						FrequencyInterval: 1,
						RetentionUnit:     "weekly",
						RetentionValue:    1,
					},
				},
			},
		}
		bSchedule := &v1.AtlasBackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "my-schedule",
				Namespace: project.Namespace,
			},
			Spec: v1.AtlasBackupScheduleSpec{
				PolicyRef: common.ResourceRefNamespaced{
					Name:      bPolicy.Name,
					Namespace: bPolicy.Namespace,
				},
				ReferenceHourOfDay:    20,
				ReferenceMinuteOfHour: 30,
				RestoreWindowDays:     7,
			},
		}
		deployment := v1.DefaultAwsAdvancedDeployment(project.Namespace, project.Name)
		deployment.Spec.DeploymentSpec.BackupEnabled = toptr.MakePtr(true)
		deployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
			Name:      bSchedule.Name,
			Namespace: bSchedule.Namespace,
		}

		sch := runtime.NewScheme()
		sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{})
		sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasProject{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasDeployment{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupSchedule{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupScheduleList{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasBackupPolicy{})
		sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasDatabaseUserList{})
		// Subresources need to be explicitly set now since controller-runtime 1.15
		// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
		k8sClient := fake.NewClientBuilder().
			WithScheme(sch).
			WithObjects(secret, project, bPolicy, bSchedule, deployment).
			WithStatusSubresource(bPolicy, bSchedule).
			Build()

		logger := zaptest.NewLogger(t).Sugar()
		atlasProvider := &atlas_mock.TestProvider{
			ClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*mongodbatlas.Client, string, error) {
				return &mongodbatlas.Client{
					AdvancedClusters: &atlas_mock.AdvancedClustersClientMock{
						GetFunc: func(projectID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
							return &mongodbatlas.AdvancedCluster{
								ID:            "123789",
								Name:          clusterName,
								GroupID:       projectID,
								BackupEnabled: toptr.MakePtr(true),
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
													NodeCount:    toptr.MakePtr(3),
												},
												Priority: toptr.MakePtr(7),
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
					GlobalClusters: &atlas_mock.GlobalClustersClientMock{
						GetFunc: func(projectID string, clusterName string) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
							return &mongodbatlas.GlobalCluster{}, nil, nil
						},
					},
					CloudProviderSnapshotBackupPolicies: &atlas_mock.CloudProviderSnapshotBackupPoliciesClientMock{
						GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
							return &mongodbatlas.CloudProviderSnapshotBackupPolicy{
								ClusterID:             "123789",
								ClusterName:           deployment.GetDeploymentName(),
								ReferenceHourOfDay:    toptr.MakePtr(int64(20)),
								ReferenceMinuteOfHour: toptr.MakePtr(int64(30)),
								RestoreWindowDays:     toptr.MakePtr(int64(7)),
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
								AutoExportEnabled:                 toptr.MakePtr(false),
								UseOrgAndGroupNamesInExportPrefix: toptr.MakePtr(false),
							}, nil, nil
						},
					},
				}, "0987654321", nil
			},
			IsCloudGovFunc: func() bool {
				return false
			},
			IsSupportedFunc: func() bool {
				return true
			},
		}

		reconciler := &AtlasDeploymentReconciler{
			ResourceWatcher:             watch.NewResourceWatcher(),
			Client:                      k8sClient,
			Log:                         logger,
			AtlasProvider:               atlasProvider,
			EventRecorder:               record.NewFakeRecorder(10),
			ObjectDeletionProtection:    false,
			SubObjectDeletionProtection: false,
		}

		result, err := reconciler.Reconcile(
			context.Background(),
			ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: deployment.Namespace,
					Name:      deployment.Name,
				},
			},
		)
		assert.NoError(t, err)
		assert.Equal(t, ctrl.Result{Requeue: false, RequeueAfter: 0}, result)
	})
}
