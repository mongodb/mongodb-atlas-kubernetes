package atlascustomrole

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/version"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type AtlasCustomRoleReconciler struct {
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	GlobalPredicates            []predicate.Predicate
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
	independentSyncPeriod       time.Duration
}

func NewAtlasCustomRoleReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	independentSyncPeriod time.Duration,
	logger *zap.Logger,
) *AtlasCustomRoleReconciler {
	return &AtlasCustomRoleReconciler{
		Client:                   mgr.GetClient(),
		Log:                      logger.Named("controllers").Named("AtlasCustomRoles").Sugar(),
		Scheme:                   mgr.GetScheme(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasCustomRoles"),
		AtlasProvider:            atlasProvider,
		GlobalPredicates:         predicates,
		ObjectDeletionProtection: deletionProtection,
		independentSyncPeriod:    independentSyncPeriod,
	}
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlascustomroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlascustomroles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlascustomroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlascustomroles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasCustomRoleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	atlasCustomRole := &akov2.AtlasCustomRole{}

	err := r.Client.Get(ctx, req.NamespacedName, atlasCustomRole)
	objectNotFound := err != nil && apiErrors.IsNotFound(err)
	failedToRetrieve := err != nil && !objectNotFound

	switch {
	case failedToRetrieve:
		return r.fail(req, err), nil
	case objectNotFound:
		return r.notFound(req), nil
	}

	if customresource.ReconciliationShouldBeSkipped(atlasCustomRole) {
		return r.skip(), nil
	}

	r.Log.Infow("-> Starting AtlasCustomRole reconciliation", "spec", atlasCustomRole.Spec, "status",
		atlasCustomRole.GetStatus())
	conditions := akov2.InitCondition(atlasCustomRole, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx)
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasCustomRole)
		r.Log.Infow("-> Finished AtlasCustomRole reconciliation", "spec", atlasCustomRole.Spec, "status",
			atlasCustomRole.GetStatus())
	}()

	valid, err := customresource.ResourceVersionIsValid(atlasCustomRole)
	if err != nil {
		return r.terminate(workflowCtx, atlasCustomRole, api.ResourceVersionStatus, workflow.AtlasResourceVersionIsInvalid, true, err), nil
	}

	if !valid {
		return r.terminate(workflowCtx,
			atlasCustomRole,
			api.ResourceVersionStatus,
			workflow.AtlasResourceVersionMismatch,
			true,
			fmt.Errorf("version of the resource '%s' is higher than the operator version '%s'", atlasCustomRole.GetName(), version.Version)), nil
	}
	workflowCtx.SetConditionTrue(api.ResourceVersionStatus).SetConditionTrue(api.ValidationSucceeded)

	if !r.AtlasProvider.IsResourceSupported(atlasCustomRole) {
		return r.terminate(workflowCtx, atlasCustomRole,
			api.ProjectCustomRolesReadyType, workflow.AtlasGovUnsupported,
			false,
			fmt.Errorf("the %T is not supported by Atlas for government", atlasCustomRole)), nil
	}

	credentials, err := selectCredentials(workflowCtx.Context, r.Client, atlasCustomRole)

	atlasSdkClient, _, err := r.AtlasProvider.SdkClient(workflowCtx.Context, credentials, workflowCtx.Log)
	if err != nil {
		return r.terminate(workflowCtx,
			atlasCustomRole,
			api.ProjectCustomRolesReadyType,
			workflow.AtlasAPIAccessNotConfigured,
			true,
			fmt.Errorf("unable to create atlas client: %s", err.Error())), nil
	}
	workflowCtx.SdkClient = atlasSdkClient
	if res := handleCustomRole(workflowCtx, r.Client, atlasCustomRole, r.ObjectDeletionProtection); !res.IsOk() {
		return r.fail(req, fmt.Errorf("%s", res.GetMessage())), nil
	}
	return r.idle(workflowCtx), nil
}

func selectCredentials(ctx context.Context, k8sClient client.Client, akoRole *akov2.AtlasCustomRole) (*client.ObjectKey, error) {
	switch {
	// First, try the externalProjectID and it's credentials
	case akoRole.Spec.ExternalProjectIDRef != nil:
		if akoRole.Spec.ConnectionSecret == nil {
			return nil, errors.New("the 'externalProjectIDRef' is set but the 'connectionSecret' is missing")
		}
		return &client.ObjectKey{Name: akoRole.Spec.ConnectionSecret.Name, Namespace: akoRole.GetNamespace()}, nil
	// Try the external project ref
	case akoRole.Spec.ProjectRef != nil:
		// if the local credentials are set, use them
		if akoRole.Spec.ConnectionSecret != nil {
			return &client.ObjectKey{Name: akoRole.Spec.ConnectionSecret.Name, Namespace: akoRole.GetNamespace()}, nil
		}
		// otherwise, use those attached to the AtlasProject that is referenced by the externalProjectRef
		project := &akov2.AtlasProject{}
		err := k8sClient.Get(ctx,
			client.ObjectKey{Name: akoRole.Spec.ProjectRef.Name, Namespace: akoRole.Spec.ProjectRef.Namespace}, project)
		if err != nil {
			return nil, errors.Wrap(err, "can not read credentials from AtlasProject")
		}
		if project.Spec.ConnectionSecret == nil || project.Spec.ConnectionSecret.Name == "" {
			return nil, errors.Wrapf(err, "credentials for AtlasProject '%s' are not configured", project.GetName())
		}
		return &client.ObjectKey{Name: project.Spec.ConnectionSecret.Name, Namespace: project.Spec.ConnectionSecret.Namespace}, nil
	}
	return nil, errors.New("either 'externalProjectIDRef' or 'projectRef' must be set for the AtlasCustomRole resource")
}

func (r *AtlasCustomRoleReconciler) terminate(
	ctx *workflow.Context,
	object akov2.AtlasCustomResource,
	condition api.ConditionType,
	reason workflow.ConditionReason,
	retry bool,
	err error,
) ctrl.Result {
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s", object, object.GetNamespace(), object.GetName(), condition, err)
	result := workflow.Terminate(reason, err.Error())
	ctx.SetConditionFromResult(condition, result)

	if !retry {
		result = result.WithoutRetry()
	}

	return result.ReconcileResult()
}

func (r *AtlasCustomRoleReconciler) idle(ctx *workflow.Context) ctrl.Result {
	ctx.SetConditionTrue(api.ReadyType)
	return workflow.OK().ReconcileResult()
}

// fail terminates the reconciliation silently(no updates on conditions)
func (r *AtlasCustomRoleReconciler) fail(req ctrl.Request, err error) ctrl.Result {
	r.Log.Errorf("Failed to query object %s: %s", req.NamespacedName, err)
	return workflow.TerminateSilently().ReconcileResult()
}

// skip prevents the reconciliation to start and successfully return
func (r *AtlasCustomRoleReconciler) skip() ctrl.Result {
	r.Log.Infow(fmt.Sprintf("-> Skipping AtlasCustomRole reconciliation as annotation %s=%s",
		customresource.ReconciliationPolicyAnnotation,
		customresource.ReconciliationPolicySkip))
	return workflow.OK().ReconcileResult()
}

// notFound terminates the reconciliation silently(no updates on conditions) and without retry
func (r *AtlasCustomRoleReconciler) notFound(req ctrl.Request) ctrl.Result {
	r.Log.Infof("Object %s doesn't exist, was it deleted after reconcile request?", req.NamespacedName)
	return workflow.TerminateSilently().WithoutRetry().ReconcileResult()
}

func (r *AtlasCustomRoleReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasCustomRole").
		For(&akov2.AtlasCustomRole{}).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.customRolesCredentials()),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{})).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
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
