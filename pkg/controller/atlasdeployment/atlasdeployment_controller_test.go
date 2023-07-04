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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
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
				AdvancedClusters: &advancedClustersClientMock{
					GetFn: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ClusterNotFound}
					},
				},
				ServerlessInstances: &serverlessClientMock{
					GetFn: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceNotFound}
					},
				},
			}
			project := testProject(fakeNamespace)
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			te := newTestDeploymentEnv(t, tc.protected, atlasClient, testK8sClient(), project, deployment)
			if tc.managedTag {
				customresource.SetAnnotation(te.deployment, customresource.AnnotationLastAppliedConfiguration, "")
			}

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.context, te.log, te.project, te.deployment)

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
			protected := true
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &advancedClustersClientMock{
					GetFn: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			te := newTestDeploymentEnv(t, protected, atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.context, te.log, te.project, te.deployment)

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
			protected := true
			project := testProject(fakeNamespace)
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &advancedClustersClientMock{
					GetFn: func(groupID string, clusterName string) (*mongodbatlas.AdvancedCluster, *mongodbatlas.Response, error) {
						return nil, nil, &mongodbatlas.ErrorResponse{ErrorCode: atlas.ServerlessInstanceFromClusterAPI}
					},
				},
				ServerlessInstances: &serverlessClientMock{
					GetFn: func(groupID string, name string) (*mongodbatlas.Cluster, *mongodbatlas.Response, error) {
						return tc.inAtlas, nil, nil
					},
				},
			}
			deployment := v1.NewDefaultAWSServerlessInstance(project.Namespace, project.Name)
			te := newTestDeploymentEnv(t, protected, atlasClient, testK8sClient(), project, deployment)

			result := te.reconciler.checkDeploymentIsManaged(te.workflowCtx, te.context, te.log, te.project, te.deployment)

			if tc.expectedErr == "" {
				assert.True(t, result.IsOk())
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedErr), result.GetMessage())
			}
		})
	}
}

func TestFinalizerNotFound(t *testing.T) {
	protected := false
	atlasClient := mongodbatlas.Client{}
	project := testProject(fakeNamespace)
	deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	k8sclient := testK8sClient()
	te := newTestDeploymentEnv(t, protected, atlasClient, k8sclient, project, deployment)

	deletionRequest, result := te.reconciler.handleDeletion(
		te.workflowCtx,
		te.context,
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
			protected := false
			atlasClient := mongodbatlas.Client{}
			project := testProject(fakeNamespace)
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			if tc.haveFinalizer {
				customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			}
			k8sclient := testK8sClient()
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			te := newTestDeploymentEnv(t, protected, atlasClient, k8sclient, project, deployment)

			deletionRequest, _ := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.context,
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
		expectRemoval bool
	}{
		{
			title:         "Deployment with protection ON and no annotations is kept",
			protected:     true,
			expectRemoval: false,
		},
		{
			title:         "Deployment with protection OFF and no annotations is removed",
			protected:     false,
			expectRemoval: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			project := testProject(fakeNamespace)
			called := false
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &advancedClustersClientMock{
					DeleteFn: func(groupID string, clusterName string, options *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error) {
						called = true
						return nil, nil
					},
				},
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			te := newTestDeploymentEnv(t, tc.protected, atlasClient, k8sclient, project, deployment)

			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.context,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Equal(t, tc.expectRemoval, called)
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
			project := testProject(fakeNamespace)
			called := false
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &advancedClustersClientMock{
					DeleteFn: func(groupID string, clusterName string, options *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error) {
						called = true
						return nil, nil
					},
				},
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			customresource.SetAnnotation(deployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyKeep,
			)
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			te := newTestDeploymentEnv(t, tc.protected, atlasClient, k8sclient, project, deployment)

			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.context,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Equal(t, false, called)
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
			project := testProject(fakeNamespace)
			called := false
			atlasClient := mongodbatlas.Client{
				AdvancedClusters: &advancedClustersClientMock{
					DeleteFn: func(groupID string, clusterName string, options *mongodbatlas.DeleteAdvanceClusterOptions) (*mongodbatlas.Response, error) {
						called = true
						return nil, nil
					},
				},
			}
			deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
			deployment.SetDeletionTimestamp(&metav1.Time{Time: time.Now()})
			customresource.SetAnnotation(deployment,
				customresource.ResourcePolicyAnnotation,
				customresource.ResourcePolicyDelete,
			)
			k8sclient := testK8sClient()
			customresource.SetFinalizer(deployment, customresource.FinalizerLabel)
			require.NoError(t, k8sclient.Create(context.Background(), deployment))
			te := newTestDeploymentEnv(t, tc.protected, atlasClient, k8sclient, project, deployment)

			deletionRequest, result := te.reconciler.handleDeletion(
				te.workflowCtx,
				te.context,
				te.log,
				te.prevResult,
				te.project,
				te.deployment,
			)

			require.True(t, deletionRequest)
			require.True(t, result.IsOk())
			assert.Equal(t, true, called)
		})
	}
}

func differentAdvancedDeployment(ns string) *mongodbatlas.AdvancedCluster {
	project := testProject(ns)
	deployment := v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment)
	deployment.Spec.DeploymentSpec.ProviderSettings.InstanceSizeName = "M2"
	advancedSpec := asAdvanced(deployment).Spec.AdvancedDeploymentSpec
	return intoAdvancedAtlasCluster(advancedSpec)
}

func sameAdvancedDeployment(ns string) *mongodbatlas.AdvancedCluster {
	project := testProject(ns)
	deployment := asAdvanced(v1.NewDeployment(project.Namespace, fakeDeployment, fakeDeployment))
	advancedSpec := asAdvanced(deployment).Spec.AdvancedDeploymentSpec
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
	context     context.Context
	log         *zap.SugaredLogger
	prevResult  workflow.Result
	project     *v1.AtlasProject
	deployment  *v1.AtlasDeployment
}

func newTestDeploymentEnv(t *testing.T,
	protected bool,
	atlasClient mongodbatlas.Client,
	k8sclient client.Client,
	project *v1.AtlasProject,
	deployment *v1.AtlasDeployment,
) *testDeploymentEnv {
	t.Helper()

	log := testLog(t)
	r := testDeploymentReconciler(log, k8sclient, protected)

	prevResult := testPrevResult()
	workflowCtx := customresource.MarkReconciliationStarted(r.Client, deployment, log)
	workflowCtx.Client = atlasClient
	return &testDeploymentEnv{
		reconciler:  r,
		workflowCtx: workflowCtx,
		context:     context.Background(),
		log:         r.Log.With("atlasdeployment", "test-namespace"),
		prevResult:  prevResult,
		deployment:  deployment,
		project:     project,
	}
}

func testK8sClient() client.Client {
	sch := runtime.NewScheme()
	sch.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.SecretList{})
	sch.AddKnownTypes(v1.GroupVersion, &v1.AtlasDeployment{})
	return fake.NewClientBuilder().WithScheme(sch).Build()
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

func asAdvanced(deployment *v1.AtlasDeployment) *v1.AtlasDeployment {
	if err := ConvertLegacyDeployment(&deployment.Spec); err != nil {
		log.Fatalf("failed to convert legacy deployment: %v", err)
	}
	deployment.Spec.DeploymentSpec = nil
	return deployment
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
