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

package atlascustomrole

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/customroles"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/ratelimit"
)

type AtlasCustomRoleReconciler struct {
	reconciler.AtlasReconciler
	Scheme                      *runtime.Scheme
	EventRecorder               record.EventRecorder
	GlobalPredicates            []predicate.Predicate
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	independentSyncPeriod       time.Duration
	maxConcurrentReconciles     int
}

func NewAtlasCustomRoleReconciler(c cluster.Cluster, predicates []predicate.Predicate, atlasProvider atlas.Provider, deletionProtection bool, independentSyncPeriod time.Duration, logger *zap.Logger, globalSecretRef client.ObjectKey, maxConcurrentReconciles int) *AtlasCustomRoleReconciler {
	return &AtlasCustomRoleReconciler{
		AtlasReconciler: reconciler.AtlasReconciler{
			Client:          c.GetClient(),
			Log:             logger.Named("controllers").Named("AtlasCustomRoles").Sugar(),
			GlobalSecretRef: globalSecretRef,
			AtlasProvider:   atlasProvider,
		},
		Scheme:                   c.GetScheme(),
		EventRecorder:            c.GetEventRecorderFor("AtlasCustomRoles"),
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
		maxConcurrentReconciles:  maxConcurrentReconciles,
	}
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlascustomroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlascustomroles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="events.k8s.io",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlascustomroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlascustomroles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="events.k8s.io",namespace=default,resources=events,verbs=create;patch

func (r *AtlasCustomRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	atlasCustomRole := &akov2.AtlasCustomRole{}

	result := customresource.PrepareResource(ctx, r.Client, req, atlasCustomRole, r.Log)
	if !result.IsOk() {
		return result.ReconcileResult()
	}

	if customresource.ReconciliationShouldBeSkipped(atlasCustomRole) {
		return r.skip()
	}

	r.Log.Infow("-> Starting AtlasCustomRole reconciliation", "spec", atlasCustomRole.Spec, "status",
		atlasCustomRole.GetStatus())
	conditions := akov2.InitCondition(atlasCustomRole, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, atlasCustomRole)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasCustomRole)
		r.Log.Infow("-> Finished AtlasCustomRole reconciliation", "spec", atlasCustomRole.Spec, "status",
			atlasCustomRole.GetStatus())
	}()

	valid, err := customresource.ResourceVersionIsValid(atlasCustomRole)
	if err != nil {
		return r.terminate(workflowCtx, atlasCustomRole, api.ResourceVersionStatus, workflow.AtlasResourceVersionIsInvalid, true, err)
	}

	if !valid {
		return r.terminate(workflowCtx,
			atlasCustomRole,
			api.ResourceVersionStatus,
			workflow.AtlasResourceVersionMismatch,
			true,
			fmt.Errorf("version of the resource '%s' is higher than the operator version '%s'", atlasCustomRole.GetName(), version.Version))
	}
	workflowCtx.SetConditionTrue(api.ResourceVersionStatus).SetConditionTrue(api.ValidationSucceeded)

	if !r.AtlasProvider.IsResourceSupported(atlasCustomRole) {
		return r.terminate(workflowCtx, atlasCustomRole,
			api.ProjectCustomRolesReadyType, workflow.AtlasGovUnsupported,
			false,
			fmt.Errorf("the %T is not supported by Atlas for government", atlasCustomRole))
	}

	connectionConfig, err := r.ResolveConnectionConfig(ctx, atlasCustomRole)
	if err != nil {
		return r.fail(req, err)
	}
	atlasSdkClientSet, err := r.AtlasProvider.SdkClientSet(workflowCtx.Context, connectionConfig.Credentials, workflowCtx.Log)
	if err != nil {
		return r.terminate(workflowCtx, atlasCustomRole, api.ProjectCustomRolesReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}
	service := customroles.NewCustomRoles(atlasSdkClientSet.SdkClient20250312012.CustomDatabaseRolesApi)
	project, err := r.ResolveProject(ctx, atlasSdkClientSet.SdkClient20250312012, atlasCustomRole)
	if err != nil {
		return r.terminate(workflowCtx, atlasCustomRole, api.ProjectCustomRolesReadyType, workflow.AtlasAPIAccessNotConfigured, true, err)
	}
	if res := handleCustomRole(workflowCtx, r.Client, project, service, atlasCustomRole, r.ObjectDeletionProtection); !res.IsOk() {
		return r.fail(req, fmt.Errorf("%s", res.GetMessage()))
	}
	return r.idle(workflowCtx)
}

func (r *AtlasCustomRoleReconciler) terminate(
	ctx *workflow.Context,
	object akov2.AtlasCustomResource,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	retry bool,
	err error,
) (ctrl.Result, error) {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", object, object.GetNamespace(), object.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFromResult(condition, result)

	if !retry {
		result = result.WithoutRetry()
	}

	return result.ReconcileResult()
}

func (r *AtlasCustomRoleReconciler) idle(ctx *workflow.Context) (ctrl.Result, error) {
	ctx.SetConditionTrue(api.ReadyType)
	return workflow.OK().ReconcileResult()
}

// fail terminates the reconciliation silently(no updates on conditions)
func (r *AtlasCustomRoleReconciler) fail(req ctrl.Request, err error) (ctrl.Result, error) {
	r.Log.Errorf("Failed to query object %s: %s", req.NamespacedName, err)
	return workflow.TerminateSilently(err).ReconcileResult()
}

// skip prevents the reconciliation to start and successfully return
func (r *AtlasCustomRoleReconciler) skip() (ctrl.Result, error) {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasCustomRole reconciliation as annotation %s=%s",
		customresource.ReconciliationPolicyAnnotation,
		customresource.ReconciliationPolicySkip))
	return workflow.OK().ReconcileResult()
}

func (r *AtlasCustomRoleReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasCustomRole{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasCustomRoleReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasCustomRole").
		For(r.For()).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.customRolesCredentials()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		WithOptions(controller.TypedOptions[reconcile.Request]{
			RateLimiter:             ratelimit.NewRateLimiter[reconcile.Request](),
			SkipNameValidation:      pointer.MakePtr(skipNameValidation),
			MaxConcurrentReconciles: r.maxConcurrentReconciles}).
		Complete(r)
}

func (r *AtlasCustomRoleReconciler) customRolesCredentials() handler.MapFunc {
	return indexer.CredentialsIndexMapperFunc(
		indexer.AtlasCustomRoleCredentialsIndex,
		func() *akov2.AtlasCustomRoleList { return &akov2.AtlasCustomRoleList{} },
		indexer.CustomRoleRequests,
		r.Client,
		r.Log,
	)
}
