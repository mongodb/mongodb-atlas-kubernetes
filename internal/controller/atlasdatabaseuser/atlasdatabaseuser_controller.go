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
	"fmt"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
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

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

//nolint:stylecheck
var ErrOIDCNotEnabled = fmt.Errorf("'OIDCAuthType' field is set but OIDC authentication is disabled")

// AtlasDatabaseUserReconciler reconciles an AtlasDatabaseUser object
type AtlasDatabaseUserReconciler struct {
	reconciler.AtlasReconciler
	AtlasProvider               atlas.Provider
	Scheme                      *runtime.Scheme
	EventRecorder               record.EventRecorder
	GlobalPredicates            []predicate.Predicate
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	independentSyncPeriod       time.Duration
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
	atlasDatabaseUser := &akov2.AtlasDatabaseUser{}

	err := r.Client.Get(ctx, req.NamespacedName, atlasDatabaseUser)
	objectNotFound := err != nil && apiErrors.IsNotFound(err)
	failedToRetrieve := err != nil && !objectNotFound

	switch {
	case failedToRetrieve:
		return r.fail(req, err), nil
	case objectNotFound:
		return r.notFound(req), nil
	}

	if customresource.ReconciliationShouldBeSkipped(atlasDatabaseUser) {
		return r.skip(), nil
	}

	r.Log.Infow("-> Starting AtlasDatabaseUser reconciliation", "spec", atlasDatabaseUser.Spec, "status", atlasDatabaseUser.GetStatus())
	conditions := akov2.InitCondition(atlasDatabaseUser, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, atlasDatabaseUser)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasDatabaseUser)
	}()

	return r.handleDatabaseUser(workflowCtx, atlasDatabaseUser), nil
}

// notFound terminates the reconciliation silently(no updates on conditions) and without retry
func (r *AtlasDatabaseUserReconciler) notFound(req ctrl.Request) ctrl.Result {
	err := fmt.Errorf("object %s doesn't exist, was it deleted after reconcile request?", req.NamespacedName)
	r.Log.Infof(err.Error())
	return workflow.TerminateSilently(err).WithoutRetry().ReconcileResult()
}

// fail terminates the reconciliation silently(no updates on conditions)
func (r *AtlasDatabaseUserReconciler) fail(req ctrl.Request, err error) ctrl.Result {
	r.Log.Errorf("Failed to query object %s: %s", req.NamespacedName, err)
	return workflow.TerminateSilently(err).ReconcileResult()
}

// skip prevents the reconciliation to start and successfully return
func (r *AtlasDatabaseUserReconciler) skip() ctrl.Result {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasDatabaseUser reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip))
	return workflow.OK().ReconcileResult()
}

// terminate interrupts the reconciliation and update the conditions with a reason and error message
func (r *AtlasDatabaseUserReconciler) terminate(
	ctx *workflow.Context,
	object akov2.AtlasCustomResource,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	retry bool,
	err error,
) ctrl.Result {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", object, object.GetNamespace(), object.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFromResult(condition, result)

	if !retry {
		result = result.WithoutRetry()
	}

	return result.ReconcileResult()
}

// unmanage remove finalizer and release resource
func (r *AtlasDatabaseUserReconciler) unmanage(ctx *workflow.Context, projectID string, atlasDatabaseUser *akov2.AtlasDatabaseUser) ctrl.Result {
	err := connectionsecret.RemoveStaleSecretsByUserName(ctx.Context, r.Client, projectID, atlasDatabaseUser.Spec.Username, *atlasDatabaseUser, r.Log)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.DatabaseUserConnectionSecretsNotDeleted, true, err)
	}

	if customresource.HaveFinalizer(atlasDatabaseUser, customresource.FinalizerLabel) {
		err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasDatabaseUser, customresource.UnsetFinalizer)
		if err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasFinalizerNotRemoved, true, err)
		}
	}

	return workflow.OK().ReconcileResult()
}

// inProgress set finalizer and requeue the reconciliation
func (r *AtlasDatabaseUserReconciler) inProgress(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser, passwordVersion, msg string) ctrl.Result {
	if !customresource.HaveFinalizer(atlasDatabaseUser, customresource.FinalizerLabel) {
		if err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasDatabaseUser, customresource.SetFinalizer); err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasFinalizerNotSet, true, err)
		}
	}

	err := customresource.ApplyLastConfigApplied(ctx.Context, atlasDatabaseUser, r.Client)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	result := workflow.InProgress(workflow.DatabaseUserDeploymentAppliedChanges, msg)
	ctx.SetConditionFromResult(api.DatabaseUserReadyType, result).
		EnsureStatusOption(status.AtlasDatabaseUserNameOption(atlasDatabaseUser.Spec.Username)).
		EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(passwordVersion))

	return result.ReconcileResult()
}

// ready set finalizer and put the resource in ready state
func (r *AtlasDatabaseUserReconciler) ready(ctx *workflow.Context, atlasDatabaseUser *akov2.AtlasDatabaseUser, passwordVersion string) ctrl.Result {
	if !customresource.HaveFinalizer(atlasDatabaseUser, customresource.FinalizerLabel) {
		if err := customresource.ManageFinalizer(ctx.Context, r.Client, atlasDatabaseUser, customresource.SetFinalizer); err != nil {
			return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.AtlasFinalizerNotSet, true, err)
		}
	}

	err := customresource.ApplyLastConfigApplied(ctx.Context, atlasDatabaseUser, r.Client)
	if err != nil {
		return r.terminate(ctx, atlasDatabaseUser, api.DatabaseUserReadyType, workflow.Internal, true, err)
	}

	ctx.SetConditionTrue(api.ReadyType).
		SetConditionTrue(api.DatabaseUserReadyType).
		EnsureStatusOption(status.AtlasDatabaseUserNameOption(atlasDatabaseUser.Spec.Username)).
		EnsureStatusOption(status.AtlasDatabaseUserPasswordVersion(passwordVersion))

	if atlasDatabaseUser.Spec.ExternalProjectRef != nil {
		return workflow.Requeue(r.independentSyncPeriod).ReconcileResult()
	}

	return workflow.OK().ReconcileResult()
}

func (r *AtlasDatabaseUserReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasDatabaseUser{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasDatabaseUserReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDatabaseUser").
		For(r.For()).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.findAtlasDatabaseUserForSecret),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.databaseUsersForCredentialMapFunc()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasDatabaseUserReconciler) findAtlasDatabaseUserForSecret(ctx context.Context, obj client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		r.Log.Warnf("watching Secret but got %T", obj)
		return nil
	}

	users := &akov2.AtlasDatabaseUserList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDatabaseUserBySecretsIndex,
			client.ObjectKeyFromObject(secret).String(),
		),
	}
	err := r.Client.List(ctx, users, listOps)
	if err != nil {
		r.Log.Errorf("failed to list AtlasDatabaseUser: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(users.Items))
	for i := range users.Items {
		item := users.Items[i]
		requests = append(
			requests,
			reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      item.Name,
					Namespace: item.Namespace,
				},
			},
		)
	}

	return requests
}

func (r *AtlasDatabaseUserReconciler) databaseUsersForCredentialMapFunc() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasDatabaseUserCredentialsIndex,
		func() *akov2.AtlasDatabaseUserList { return &akov2.AtlasDatabaseUserList{} },
		indexer.DatabaseUserRequests,
		r.Client,
		r.Log,
	)
}

func NewAtlasDatabaseUserReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	independentSyncPeriod time.Duration,
	featureFlags *featureflags.FeatureFlags,
	logger *zap.Logger,
) *AtlasDatabaseUserReconciler {
	return &AtlasDatabaseUserReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client: c.GetClient(),
			Log:    logger.Named("controllers").Named("AtlasDatabaseUser").Sugar(),
		},
		AtlasProvider:            atlasProvider,
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasDatabaseUser"),
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}
