package atlasfederatedauth

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

// AtlasFederatedAuthReconciler reconciles an AtlasFederatedAuth object
type AtlasFederatedAuthReconciler struct {
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasfederatedauths,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasfederatedauths/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasfederatedauths,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasfederatedauths/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

func (r *AtlasFederatedAuthReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasfederatedauth", req.NamespacedName)

	fedauth := &akov2.AtlasFederatedAuth{}
	result := customresource.PrepareResource(ctx, r.Client, req, fedauth, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(fedauth) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasFederatedAuth reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", fedauth.Spec)
		if !fedauth.GetDeletionTimestamp().IsZero() {
			if err := customresource.ManageFinalizer(ctx, r.Client, fedauth, customresource.UnsetFinalizer); err != nil {
				result = workflow.Terminate(workflow.Internal, err)
				log.Errorw("Failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	conditions := akov2.InitCondition(fedauth, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx, fedauth)
	log.Infow("-> Starting AtlasFederatedAuth reconciliation")

	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, fedauth)

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, fedauth, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("federated auth validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if !r.AtlasProvider.IsResourceSupported(fedauth) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, errors.New("the AtlasFederatedAuth is not supported by Atlas for government")).
			WithoutRetry()
		setCondition(workflowCtx, api.FederatedAuthReadyType, result)
		return result.ReconcileResult(), nil
	}

	atlasClient, orgID, err := r.AtlasProvider.SdkClient(workflowCtx.Context, fedauth.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err)
		setCondition(workflowCtx, api.FederatedAuthReadyType, result)
		return result.ReconcileResult(), nil
	}

	atlasClientSet, _, err := r.AtlasProvider.SdkClientSet(workflowCtx.Context, fedauth.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err)
		setCondition(workflowCtx, api.FederatedAuthReadyType, result)
		return result.ReconcileResult(), nil
	}

	workflowCtx.SdkClient = atlasClient
	workflowCtx.SdkClientSet = atlasClientSet
	workflowCtx.OrgID = orgID

	result = r.ensureFederatedAuth(workflowCtx, fedauth)
	workflowCtx.SetConditionFromResult(api.FederatedAuthReadyType, result)
	workflowCtx.SetConditionFromResult(api.ReadyType, result)

	return result.ReconcileResult(), nil
}

func (r *AtlasFederatedAuthReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasFederatedAuth{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasFederatedAuthReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasFederatedAuth").
		For(r.For()).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.findAtlasFederatedAuthForSecret),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasFederatedAuthReconciler) findAtlasFederatedAuthForSecret(ctx context.Context, obj client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		r.Log.Warnf("watching Secret but got %T", obj)
		return nil
	}

	auths := &akov2.AtlasFederatedAuthList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasFederatedAuthBySecretsIndex,
			client.ObjectKeyFromObject(secret).String(),
		),
	}
	err := r.Client.List(ctx, auths, listOps)
	if err != nil {
		r.Log.Errorf("failed to list AtlasFederatedAuth: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(auths.Items))
	for i := range auths.Items {
		item := auths.Items[i]
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

func NewAtlasFederatedAuthReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasFederatedAuthReconciler {
	return &AtlasFederatedAuthReconciler{
		Scheme:                   c.GetScheme(),
		Client:                   c.GetClient(),
		EventRecorder:            c.GetEventRecorderFor("AtlasFederatedAuth"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasFederatedAuth").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func setCondition(ctx *workflow.Context, condition api.ConditionType, result workflow.Result) {
	ctx.SetConditionFromResult(condition, result)
	logIfWarning(ctx, result)
}

func logIfWarning(ctx *workflow.Context, result workflow.Result) {
	if result.IsWarning() {
		ctx.Log.Warnw(result.GetMessage())
	}
}
