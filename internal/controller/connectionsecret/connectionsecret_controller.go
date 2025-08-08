// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package connectionsecret

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

type ConnectionSecretReconciler struct {
	reconciler.AtlasReconciler
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
}

func (r *ConnectionSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Parses the request name and fills up the identifiers: ProjectID, ClusterName, DatabaseUsername
	ids, err := LoadRequestNameParts(ctx, r.Client, req.NamespacedName)
	if err != nil {
		r.Log.Errorw("Failed to parse request name", "error", err)
		return workflow.Terminate("InvalidConnectionSecretName", err).ReconcileResult()
	}

	// Loads the pair of AtlasDeployment and AtlasDatabaseUser via the indexers
	pair, err := LoadPairedResources(ctx, r.Client, ids, req.Namespace)
	if err != nil {
		r.Log.Errorw("Failed to fetch deployment and user", "error", err)
		return workflow.Terminate("InvalidConnectionResources", err).ReconcileResult()
	}

	// Checks is the pair is ready
	if ready, notReady := pair.AreResourcesReady(); !ready {
		return workflow.InProgress("ConnectionSecretNotReady", fmt.Sprintf("Not ready: %s", strings.Join(notReady, ", "))).ReconcileResult()
	}

	// Check if we have a ProjectName which is required for the K8s naming format
	if ids.ProjectName == "" {
		// If both AtlasDeployment and AtlasDatabaseUser use externalRef, resolve via SDK
		if pair.NeedsSDKProjectResolution() {
			connectionConfig, err := r.ResolveConnectionConfig(ctx, pair.Deployment)
			if err != nil {
				r.Log.Errorw("Failed to resolve Atlas connection config", "error", err)
				return workflow.Terminate("FailedToResolveConnectionConfig", err).ReconcileResult()
			}

			sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, r.Log)
			if err != nil {
				r.Log.Errorw("Failed to create SDK client", "error", err)
				return workflow.Terminate("FailedToCreateAtlasClient", err).ReconcileResult()
			}

			projectService := project.NewProjectAPIService(sdkClientSet.SdkClient20250312002.ProjectsApi)
			err = pair.ResolveProjectNameSDK(ctx, projectService)
			if err != nil {
				r.Log.Errorw("Failed to fetch project from Atlas API", "projectID", ids.ProjectID, "error", err)
				return workflow.Terminate("FailedToFetchProjectFromAtlas", err).ReconcileResult()
			}
		} else { // Otherwise, fetch AtlasProject from K8s
			err = pair.ResolveProjectNameK8s(ctx, r.Client, req.Namespace)
			if err != nil {
				r.Log.Errorw("Failed to retrieve AtlasProject to resolve ProjectName", "projectID", ids.ProjectID, "error", err)
				return workflow.Terminate("FailedToResolveProjectName", err).ReconcileResult()
			}
		}
	}

	// Retrives the connection data
	data, err := pair.BuildConnectionData(ctx, r.Client)
	if err != nil {
		r.Log.Errorw("Failed to build connection data", "error", err)
		return workflow.Terminate("FailedToBuildConnectionData", err).ReconcileResult()
	}

	// Creates or Updates the connection secret
	if err := pair.HandleSecret(ctx, r.Client, data); err != nil {
		r.Log.Errorw("Failed to create or update connection Secret", "error", err)
		return workflow.Terminate("FailedToEnsureConnectionSecret", err).ReconcileResult()
	}

	r.Log.Infow("Successfully found associated resources",
		"deployment", pair.Deployment.GetDeploymentName(),
		"user", pair.User.Spec.Username,
		"projectID", ids.ProjectID,
	)

	return reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *ConnectionSecretReconciler) DeploymentReadyPredicate() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		d, ok := obj.(*akov2.AtlasDeployment)
		return ok && IsDeploymentReady(d)
	})
}

func (r *ConnectionSecretReconciler) DatabaseUserReadyPredicate() predicate.Predicate {
	return predicate.NewPredicateFuncs(func(obj client.Object) bool {
		u, ok := obj.(*akov2.AtlasDatabaseUser)
		return ok && IsDatabaseUserReady(u)
	})
}

func (r *ConnectionSecretReconciler) For() (client.Object, builder.Predicates) {
	labelPredicates := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		labels := obj.GetLabels()
		_, hasType := labels[TypeLabelKey]
		_, hasProject := labels[ProjectLabelKey]
		_, hasCluster := labels[ClusterLabelKey]
		return hasType && hasProject && hasCluster
	})

	predicates := append(r.GlobalPredicates, labelPredicates)
	return &corev1.Secret{}, builder.WithPredicates(predicates...)
}

func (r *ConnectionSecretReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("ConnectionSecret").
		For(r.For()).
		Watches(
			&akov2.AtlasDeployment{},
			handler.EnqueueRequestsFromMapFunc(r.newDeploymentMapFunc),
			builder.WithPredicates(r.DeploymentReadyPredicate()),
		).
		Watches(
			&akov2.AtlasDatabaseUser{},
			handler.EnqueueRequestsFromMapFunc(r.newDatabaseUserMapFunc),
			builder.WithPredicates(r.DatabaseUserReadyPredicate()),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:        ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation: pointer.MakePtr(skipNameValidation),
		}).
		Complete(r)
}

func (r *ConnectionSecretReconciler) generateConnectionSecretRequests(
	projectID string,
	deployments []akov2.AtlasDeployment,
	users []akov2.AtlasDatabaseUser,
) []reconcile.Request {
	var requests []reconcile.Request
	for _, d := range deployments {
		for _, u := range users {
			scopes := u.GetScopes(akov2.DeploymentScopeType)
			if len(scopes) != 0 && !stringutil.Contains(scopes, d.GetDeploymentName()) {
				continue
			}

			requestName := CreateInternalFormat(projectID, d.GetDeploymentName(), u.Spec.Username)
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: u.Namespace,
					Name:      requestName,
				},
			})
		}
	}
	return requests
}

func (r *ConnectionSecretReconciler) newDeploymentMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	deployment, ok := obj.(*akov2.AtlasDeployment)
	if !ok {
		r.Log.Warnf("watching AtlasDeployment but got %T", obj)
		return nil
	}

	projectID, err := ResolveProjectIDFromDeployment(ctx, r.Client, deployment)
	if err != nil {
		r.Log.Errorw("Unable to resolve projectID for deployment", "error", err)
		return nil
	}

	users := &akov2.AtlasDatabaseUserList{}
	if err := r.Client.List(ctx, users, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDatabaseUserByProject, projectID),
	}); err != nil {
		r.Log.Errorf("failed to list AtlasDatabaseUsers: %e", err)
		return nil
	}

	return r.generateConnectionSecretRequests(projectID, []akov2.AtlasDeployment{*deployment}, users.Items)
}

func (r *ConnectionSecretReconciler) newDatabaseUserMapFunc(ctx context.Context, obj client.Object) []reconcile.Request {
	user, ok := obj.(*akov2.AtlasDatabaseUser)
	if !ok {
		r.Log.Warnf("watching AtlasDatabaseUser but got %T", obj)
		return nil
	}

	projectID, err := ResolveProjectIDFromDatabaseUser(ctx, r.Client, user)
	if err != nil {
		r.Log.Errorw("Unable to resolve projectID for user", "error", err)
		return nil
	}

	deployments := &akov2.AtlasDeploymentList{}
	if err := r.Client.List(ctx, deployments, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(indexer.AtlasDeploymentByProject, projectID),
	}); err != nil {
		r.Log.Errorf("failed to list AtlasDeployments: %e", err)
		return nil
	}

	return r.generateConnectionSecretRequests(projectID, deployments.Items, []akov2.AtlasDatabaseUser{*user})
}

func NewConnectionSecretReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	logger *zap.Logger,
	globalSecretRef types.NamespacedName,
) *ConnectionSecretReconciler {
	return &ConnectionSecretReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             logger.Named("controllers").Named("ConnectionSecret").Sugar(),
			GlobalSecretRef: globalSecretRef,
			AtlasProvider:   atlasProvider,
		},
		EventRecorder:    c.GetEventRecorderFor("ConnectionSecret"),
		GlobalPredicates: predicates,
	}
}
