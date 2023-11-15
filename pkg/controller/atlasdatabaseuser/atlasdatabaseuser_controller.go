/*
Copyright 2020 MongoDB.

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

package atlasdatabaseuser

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

// AtlasDatabaseUserReconciler reconciles an AtlasDatabaseUser object
type AtlasDatabaseUserReconciler struct {
	watch.ResourceWatcher
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	AtlasDomain                 string
	GlobalAPISecret             client.ObjectKey
	EventRecorder               record.EventRecorder
	GlobalPredicates            []predicate.Predicate
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasDatabaseUserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdatabaseuser", req.NamespacedName)

	databaseUser := &mdbv1.AtlasDatabaseUser{}
	result := customresource.PrepareResource(r.Client, req, databaseUser, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(databaseUser) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDatabaseUser reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", databaseUser.Spec)
		return workflow.OK().ReconcileResult(), nil
	}

	workflowCtx := customresource.MarkReconciliationStarted(r.Client, databaseUser, log, ctx)
	log.Infow("-> Starting AtlasDatabaseUser reconciliation", "spec", databaseUser.Spec, "status", databaseUser.Status)
	if databaseUser.Spec.PasswordSecret != nil {
		workflowCtx.AddResourcesToWatch(watch.WatchedObject{ResourceKind: "Secret", Resource: *databaseUser.PasswordSecretObjectKey()})
	}
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, databaseUser)
		r.EnsureMultiplesResourcesAreWatched(req.NamespacedName, log, workflowCtx.ListResourcesToWatch()...)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, databaseUser, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("database user validation result: %v", resourceVersionIsValid)

		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if err := validate.DatabaseUser(databaseUser); err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.ValidationSucceeded, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(status.ValidationSucceeded)

	if !customresource.IsResourceSupportedInDomain(databaseUser, r.AtlasDomain) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasDatabaseUser is not supported by Atlas for government").
			WithoutRetry()
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	project := &mdbv1.AtlasProject{}
	if result = r.readProjectResource(databaseUser, project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result = workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.Client = atlasClient

	owner, err := customresource.IsOwner(databaseUser, r.ObjectDeletionProtection, customresource.IsResourceManagedByOperator, managedByAtlas(ctx, atlasClient, project.ID(), log))
	if err != nil {
		result = workflow.Terminate(workflow.Internal, fmt.Sprintf("enable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	if !owner {
		result = workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile database user: it already exists in Atlas, it was not previously managed by the operator, and the deletion protection is enabled.",
		)
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	deletionRequest, result := r.handleDeletion(ctx, databaseUser, project, atlasClient, log)
	if deletionRequest {
		return result.ReconcileResult(), nil
	}

	err = customresource.ApplyLastConfigApplied(ctx, databaseUser, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	result = r.ensureDatabaseUser(workflowCtx, *project, *databaseUser)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}

	err = customresource.ManageFinalizer(ctx, r.Client, databaseUser, customresource.SetFinalizer)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasFinalizerNotSet, err.Error())
		workflowCtx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	workflowCtx.SetConditionTrue(status.DatabaseUserReadyType)
	workflowCtx.SetConditionTrue(status.ReadyType)

	return result.ReconcileResult(), nil
}

func (r *AtlasDatabaseUserReconciler) readProjectResource(user *mdbv1.AtlasDatabaseUser, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(context.Background(), user.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDatabaseUserReconciler) handleDeletion(
	ctx context.Context,
	dbUser *mdbv1.AtlasDatabaseUser,
	project *mdbv1.AtlasProject,
	atlasClient mongodbatlas.Client,
	log *zap.SugaredLogger,
) (bool, workflow.Result) {
	if dbUser.GetDeletionTimestamp().IsZero() {
		return false, workflow.OK()
	}

	if customresource.HaveFinalizer(dbUser, customresource.FinalizerLabel) {
		err := connectionsecret.RemoveStaleSecretsByUserName(r.Client, project.ID(), dbUser.Spec.Username, *dbUser, log)
		if err != nil {
			return true, workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotDeleted, err.Error())
		}
	}

	if customresource.IsResourceProtected(dbUser, r.ObjectDeletionProtection) {
		log.Info("Not removing Atlas database user from Atlas as per configuration")

		err := customresource.ManageFinalizer(ctx, r.Client, dbUser, customresource.UnsetFinalizer)
		if err != nil {
			return true, workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
		}

		return true, workflow.OK()
	}

	_, err := atlasClient.DatabaseUsers.Delete(context.Background(), dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
	if err != nil {
		var apiError *mongodbatlas.ErrorResponse
		if errors.As(err, &apiError) && apiError.ErrorCode != atlas.UsernameNotFound {
			return true, workflow.Terminate(workflow.DatabaseUserNotDeletedInAtlas, err.Error())
		}

		log.Info("Database user doesn't exist or is already deleted")
	}

	err = customresource.ManageFinalizer(ctx, r.Client, dbUser, customresource.UnsetFinalizer)
	if err != nil {
		return true, workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
	}

	return true, workflow.OK()
}

func (r *AtlasDatabaseUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDatabaseUser").
		For(&mdbv1.AtlasDatabaseUser{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources)).
		Complete(r)
}

func managedByAtlas(ctx context.Context, atlasClient mongodbatlas.Client, projectID string, log *zap.SugaredLogger) customresource.AtlasChecker {
	return func(resource mdbv1.AtlasCustomResource) (bool, error) {
		dbUser, ok := resource.(*mdbv1.AtlasDatabaseUser)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasDatabaseUser")
		}

		atlasDBUser, _, err := atlasClient.DatabaseUsers.Get(ctx, dbUser.Spec.DatabaseName, projectID, dbUser.Spec.Username)
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
				return false, nil
			}

			return false, err
		}

		isSame, err := userMatchesSpec(log, atlasDBUser, dbUser.Spec)
		if err != nil {
			return true, err
		}

		return !isSame, nil
	}
}
