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
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.mongodb.org/atlas-sdk/v20231115008/mockadmin"
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
	atlasmock "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/mocks/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
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
				AdvancedClusters: &atlasmock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ClusterNotFound}
					},
				},
				ServerlessInstances: &atlasmock.ServerlessInstancesClientMock{
					GetFunc: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceNotFound}
					},
				},
			}
			project := testProject(fakeNamespace)
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &atlasmock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			te := newTestDeploymentEnv(t, true, &atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.log, te.project, te.deployment)

			if tc.expectedErr == "" {
				assert.True(t, result.IsOk())
			} else {
				assert.NotNil(t, te.workflowCtx.LastCondition())
				if te.workflowCtx.LastCondition() != nil {
					assert.Equal(t, status.DeploymentReadyType, te.workflowCtx.LastCondition().Type)
				}
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
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &atlasmock.AdvancedClustersClientMock{
					GetFunc: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceFromClusterAPI}
					},
				},
				ServerlessInstances: &atlasmock.ServerlessInstancesClientMock{
					GetFunc: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := akov2.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
			te := newTestDeploymentEnv(t, true, &atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.log, te.project, te.deployment)

			if tc.expectedErr == "" {
				assert.True(t, result.IsOk())
			} else {
				assert.NotNil(t, te.workflowCtx.LastCondition())
				if te.workflowCtx.LastCondition() != nil {
					assert.Equal(t, status.DeploymentReadyType, te.workflowCtx.LastCondition().Type)
				}
				assert.Regexp(t, regexp.MustCompile(tc.expectedErr), result.GetMessage())
			}
		})
	}
}

func TestFinalizerNotFound(t *testing.T) {
	atlasClient := mongodbatlas.Client{}
	project := testProject(fakeNamespace)
	deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
			finalDeployment := &akov2.AtlasDeployment{}
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
			advancedClusterClient := &atlasmock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
			advancedClusterClient := &atlasmock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
			advancedClusterClient := &atlasmock.AdvancedClustersClientMock{
				DeleteFunc: func(groupID string, clusterName string) (*mongodbatlas.Response, error) {
					return nil, nil
				},
			}
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: advancedClusterClient,
			}
			deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
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
		d := &akov2.AtlasDeployment{} // dummy deployment

		// test cleanup
		assert.NoError(t, r.cleanupBindings(context.Background(), d))
	})

	t.Run("with unreferenced backups, still nothing happens on cleanup", func(t *testing.T) {
		r := &AtlasDeploymentReconciler{
			Log:    testLog(t),
			Client: testK8sClient(),
		}
		dn := testDeploymentName("") // deployment, schedule, policy (NOT connected)
		deployment := &akov2.AtlasDeployment{
			ObjectMeta: metav1.ObjectMeta{Name: dn.Name, Namespace: dn.Namespace},
		}
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		policy := testBackupPolicy()
		require.NoError(t, r.Client.Create(context.Background(), policy))
		schedule := testBackupSchedule("", policy)
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test cleanup
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

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
		deployment := testDeployment("", schedule)
		require.NoError(t, r.Client.Create(context.Background(), deployment))
		schedule.Status.DeploymentIDs = []string{deployment.Spec.DeploymentSpec.Name}
		require.NoError(t, r.Client.Create(context.Background(), schedule))

		// test ensureBackupPolicy and cleanup
		_, err := r.ensureBackupPolicy(&workflow.Context{Context: context.Background()}, schedule, &[]watch.WatchedObject{})
		require.NoError(t, err)
		require.NoError(t, r.cleanupBindings(context.Background(), deployment))

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

		endPolicy := &akov2.AtlasBackupPolicy{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(policy.Namespace, policy.Name), endPolicy))
		assert.NotEmpty(t, endPolicy.Finalizers, "policy should keep the finalizer")
		endSchedule := &akov2.AtlasBackupSchedule{}
		require.NoError(t, r.Client.Get(context.Background(), kube.ObjectKey(schedule.Namespace, schedule.Name), endSchedule))
		assert.Empty(t, endSchedule.Finalizers, "schedule should end up with no finalizer")
	})
}

func sameAdvancedDeployment(ns string) *mongodbatlas.AdvancedCluster {
	project := testProject(ns)
	deployment := akov2.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	advancedSpec := deployment.Spec.DeploymentSpec
	return intoAdvancedAtlasCluster(advancedSpec)
}

func sameServerlessDeployment(ns string) *mongodbatlas.Cluster {
	project := testProject(ns)
	deployment := akov2.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
	return intoServerlessAtlasCluster(deployment.Spec.ServerlessSpec)
}

type testDeploymentEnv struct {
	reconciler  *AtlasDeploymentReconciler
	workflowCtx *workflow.Context
	log         *zap.SugaredLogger
	prevResult  workflow.Result
	project     *akov2.AtlasProject
	deployment  *akov2.AtlasDeployment
}

func newTestDeploymentEnv(t *testing.T,
	protected bool,
	atlasClient *mongodbatlas.Client,
	k8sclient client.Client,
	project *akov2.AtlasProject,
	deployment *akov2.AtlasDeployment,
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

func testProject(ns string) *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeProject,
			Namespace: ns,
		},
		Status: status.AtlasProjectStatus{
			ID: fakeProjectID,
		},
	}
}

func intoAdvancedAtlasCluster(advancedSpec *akov2.AdvancedDeploymentSpec) *mongodbatlas.AdvancedCluster {
	ac, err := advancedSpec.ToAtlas()
	if err != nil {
		log.Fatalf("failed to convert advanced deployment to atlas: %v", err)
	}
	return ac
}

func intoServerlessAtlasCluster(serverlessSpec *akov2.ServerlessSpec) *mongodbatlas.Cluster {
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

func testDeployment(suffix string, schedule *akov2.AtlasBackupSchedule) *akov2.AtlasDeployment {
	dn := testDeploymentName(suffix)
	return &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{Name: dn.Name, Namespace: dn.Namespace},
		Spec: akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Name: fmt.Sprintf("atlas-%s", dn.Name),
			},
			BackupScheduleRef: common.ResourceRefNamespaced{
				Name:      schedule.Name,
				Namespace: schedule.Namespace,
			},
		},
	}
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

func TestUniqueKey(t *testing.T) {
	t.Run("Test duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Tags: []*akov2.TagSpec{{Key: "foo", Value: "true"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Advanced Deployment", func(t *testing.T) {
		deploymentSpec := &akov2.AtlasDeploymentSpec{
			DeploymentSpec: &akov2.AdvancedDeploymentSpec{
				Tags: []*akov2.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foobar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
	t.Run("Test duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &akov2.AtlasDeploymentSpec{
			ServerlessSpec: &akov2.ServerlessSpec{
				Tags: []*akov2.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}, {Key: "foo", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.Error(t, err)
	})
	t.Run("Test no duplicates in Serverless Instance", func(t *testing.T) {
		deploymentSpec := &akov2.AtlasDeploymentSpec{
			ServerlessSpec: &akov2.ServerlessSpec{
				Tags: []*akov2.TagSpec{{Key: "foo", Value: "true"}, {Key: "bar", Value: "false"}},
			},
		}
		err := uniqueKey(deploymentSpec)
		assert.NoError(t, err)
	})
}

func TestReconciliation(t *testing.T) {
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
	deployment := akov2.DefaultAwsAdvancedDeployment(project.Namespace, project.Name)
	deployment.Spec.DeploymentSpec.BackupEnabled = pointer.MakePtr(true)
	deployment.Spec.BackupScheduleRef = common.ResourceRefNamespaced{
		Name:      bSchedule.Name,
		Namespace: bSchedule.Namespace,
	}
	deployment.Spec.DeploymentSpec.SearchNodes = searchNodes

	sch := runtime.NewScheme()
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Secret{})
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasProject{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasDeployment{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupSchedule{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupScheduleList{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasBackupPolicy{})
	sch.AddKnownTypes(akov2.GroupVersion, &akov2.AtlasDatabaseUserList{})
	// Subresources need to be explicitly set now since controller-runtime 1.15
	// https://github.com/kubernetes-sigs/controller-runtime/issues/2362#issuecomment-1698194188
	k8sClient := fake.NewClientBuilder().
		WithScheme(sch).
		WithObjects(secret, project, bPolicy, bSchedule, deployment).
		WithStatusSubresource(bPolicy, bSchedule).
		Build()

	searchAPI := mockadmin.NewAtlasSearchApi(t)
	searchAPI.EXPECT().GetAtlasSearchDeployment(context.Background(), project.ID(), deployment.Spec.DeploymentSpec.Name).
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

	orgID := "0987654321"
	logger := zaptest.NewLogger(t).Sugar()
	atlasProvider := &atlasmock.TestProvider{
		SdkClientFunc: func(secretRef *client.ObjectKey, log *zap.SugaredLogger) (*admin.APIClient, string, error) {
			return &admin.APIClient{
				AtlasSearchApi: searchAPI,
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
				GlobalClusters: &atlasmock.GlobalClustersClientMock{
					GetFunc: func(projectID string, clusterName string) (*mongodbatlas.GlobalCluster, *mongodbatlas.Response, error) {
						return &mongodbatlas.GlobalCluster{}, nil, nil
					},
				},
				CloudProviderSnapshotBackupPolicies: &atlasmock.CloudProviderSnapshotBackupPoliciesClientMock{
					GetFunc: func(projectID string, clusterName string) (*mongodbatlas.CloudProviderSnapshotBackupPolicy, *mongodbatlas.Response, error) {
						return &mongodbatlas.CloudProviderSnapshotBackupPolicy{
							ClusterID:             "123789",
							ClusterName:           deployment.GetDeploymentName(),
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
		DeprecatedResourceWatcher:   watch.NewDeprecatedResourceWatcher(),
		Client:                      k8sClient,
		Log:                         logger,
		AtlasProvider:               atlasProvider,
		EventRecorder:               record.NewFakeRecorder(10),
		ObjectDeletionProtection:    false,
		SubObjectDeletionProtection: false,
	}

	t.Run("should reconcile with existing cluster", func(t *testing.T) {
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
