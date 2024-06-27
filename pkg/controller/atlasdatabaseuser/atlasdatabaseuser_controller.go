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
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

//nolint:stylecheck
var ErrOIDCNotEnabled = fmt.Errorf("'OIDCAuthType' field is set but OIDC authentication is disabled")

// AtlasDatabaseUserReconciler reconciles an AtlasDatabaseUser object
type AtlasDatabaseUserReconciler struct {
	watch.DeprecatedResourceWatcher
	Client                        client.Client
	Log                           *zap.SugaredLogger
	Scheme                        *runtime.Scheme
	EventRecorder                 record.EventRecorder
	AtlasProvider                 atlas.Provider
	GlobalPredicates              []predicate.Predicate
	ObjectDeletionProtection      bool
	SubObjectDeletionProtection   bool
	FeaturePreviewOIDCAuthEnabled bool
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

	databaseUser := &akov2.AtlasDatabaseUser{}
	result := customresource.PrepareResource(ctx, r.Client, req, databaseUser, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(databaseUser) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDatabaseUser reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", databaseUser.Spec)
		return workflow.OK().ReconcileResult(), nil
	}

	conditions := akov2.InitCondition(databaseUser, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx)
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
		workflowCtx.SetConditionFromResult(api.ValidationSucceeded, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(api.ValidationSucceeded)

	if !r.AtlasProvider.IsResourceSupported(databaseUser) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasDatabaseUser is not supported by Atlas for government").
			WithoutRetry()
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	project := &akov2.AtlasProject{}
	if result = r.readProjectResource(ctx, databaseUser, project); !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}

	atlasClient, orgID, err := r.AtlasProvider.Client(ctx, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}
	workflowCtx.OrgID = orgID
	workflowCtx.Client = atlasClient

	deletionRequest, result := r.handleDeletion(ctx, databaseUser, project, atlasClient, log)
	if deletionRequest {
		return result.ReconcileResult(), nil
	}

	err = customresource.ApplyLastConfigApplied(ctx, databaseUser, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	err = r.handleFeatureFlags(databaseUser)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(api.ReadyType, result)
		log.Error(result.GetMessage())
		return result.ReconcileResult(), nil
	}

	result = r.ensureDatabaseUser(workflowCtx, *project, *databaseUser)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)

		return result.ReconcileResult(), nil
	}

	err = customresource.ManageFinalizer(ctx, r.Client, databaseUser, customresource.SetFinalizer)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasFinalizerNotSet, err.Error())
		workflowCtx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	workflowCtx.SetConditionTrue(api.DatabaseUserReadyType)
	workflowCtx.SetConditionTrue(api.ReadyType)

	return result.ReconcileResult(), nil
}

func (r *AtlasDatabaseUserReconciler) handleFeatureFlags(dbuser *akov2.AtlasDatabaseUser) error {
	err := handleOIDCPreview(r.FeaturePreviewOIDCAuthEnabled, dbuser)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Remove after the OIDC feature becomes stable
func handleOIDCPreview(OIDCEnabled bool, dbuser *akov2.AtlasDatabaseUser) error {
	if dbuser == nil {
		return nil
	}

	if !OIDCEnabled && dbuser.Spec.OIDCAuthType == "IDP_GROUP" {
		return ErrOIDCNotEnabled
	}

	return nil
}

func (r *AtlasDatabaseUserReconciler) readProjectResource(ctx context.Context, user *akov2.AtlasDatabaseUser, project *akov2.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, user.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDatabaseUserReconciler) handleDeletion(
	ctx context.Context,
	dbUser *akov2.AtlasDatabaseUser,
	project *akov2.AtlasProject,
	atlasClient *mongodbatlas.Client,
	log *zap.SugaredLogger,
) (bool, workflow.Result) {
	if dbUser.GetDeletionTimestamp().IsZero() {
		return false, workflow.OK()
	}

	if customresource.HaveFinalizer(dbUser, customresource.FinalizerLabel) {
		err := connectionsecret.RemoveStaleSecretsByUserName(ctx, r.Client, project.ID(), dbUser.Spec.Username, *dbUser, log)
		if err != nil {
			return true, workflow.Terminate(workflow.DatabaseUserConnectionSecretsNotDeleted, err.Error())
		}
	}

	if customresource.IsResourcePolicyKeepOrDefault(dbUser, r.ObjectDeletionProtection) {
		log.Info("Not removing Atlas database user from Atlas as per configuration")

		err := customresource.ManageFinalizer(ctx, r.Client, dbUser, customresource.UnsetFinalizer)
		if err != nil {
			return true, workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
		}

		return true, workflow.OK()
	}

	_, err := atlasClient.DatabaseUsers.Delete(ctx, dbUser.Spec.DatabaseName, project.ID(), dbUser.Spec.Username)
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
		For(&akov2.AtlasDatabaseUser{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(&corev1.Secret{}, watch.NewSecretHandler(&r.DeprecatedResourceWatcher)).
		Complete(r)
}

func NewAtlasDatabaseUserReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	featureFlags *featureflags.FeatureFlags,
	logger *zap.Logger,
) *AtlasDatabaseUserReconciler {
	return &AtlasDatabaseUserReconciler{
		Scheme:                        mgr.GetScheme(),
		Client:                        mgr.GetClient(),
		EventRecorder:                 mgr.GetEventRecorderFor("AtlasDatabaseUser"),
		DeprecatedResourceWatcher:     watch.NewDeprecatedResourceWatcher(),
		GlobalPredicates:              predicates,
		Log:                           logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		AtlasProvider:                 atlasProvider,
		ObjectDeletionProtection:      deletionProtection,
		FeaturePreviewOIDCAuthEnabled: featureFlags.IsFeaturePresent(featureflags.FeatureOIDC),
	}
}
