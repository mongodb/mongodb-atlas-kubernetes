package atlasstream

import (
	"context"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
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

const instanceNotFound = "STREAM_TENANT_NOT_FOUND_FOR_NAME"

type AtlasStreamsInstanceReconciler struct {
	Client                      client.Client
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	Log                         *zap.SugaredLogger
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreaminstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreaminstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreaminstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreaminstances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasstreamconnections,verbs=get;list
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasstreamconnections,verbs=get;list
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// https://dreampuf.github.io/GraphvizOnline/#digraph%20G%20%7B%0A%20%20%20%20subgraph%20cluster_pending%20%7B%0A%20%20%20%20%20%20%20%20skipped%3B%0A%20%20%20%20%20%20%20%20invalid%3B%0A%20%20%20%20%20%20%20%20unsupported%3B%0A%20%20%20%20%20%20%20%20terminated%3B%0A%20%20%20%20%20%20%20%20label%20%3D%20%22pending%22%3B%0A%20%20%20%20%7D%0A%0A%20%20%20%20deleted%20%5Blabel%3D%22deleted%5Cnfinalizer%20unset%22%5D%0A%0A%20%20%20%20pending%20-%3E%20pending%20%5Blabel%3D%22skip%5Cninvalidate%5Cnunsupport%5Cnterminate%22%5D%0A%20%20%20%20pending%20-%3E%20ready%20%5Blabel%3D%22create%22%5D%0A%20%20%20%20pending%20-%3E%20deleted%20%5Blabel%3D%22delete%22%5D%0A%20%20%20%20ready%20-%3E%20ready%20%5Blabel%3D%22update%22%5D%0A%20%20%20%20ready%20-%3E%20deleted%20%5Blabel%3D%22delete%22%5D%0A%20%20%20%20ready%20-%3E%20pending%20%5Blabel%3D%22terminate%22%5D%0A%7D%0A

func (r *AtlasStreamsInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasstreaminstance", req.NamespacedName)
	log.Infow("-> Starting AtlasStreamInstance reconciliation")

	akoStreamInstance := akov2.AtlasStreamInstance{}
	result := customresource.PrepareResource(ctx, r.Client, req, &akoStreamInstance, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	return r.ensureAtlasStreamsInstance(ctx, log, &akoStreamInstance)
}

// ensureAtlasStreamsInstance is the central state dispatcher
func (r *AtlasStreamsInstanceReconciler) ensureAtlasStreamsInstance(ctx context.Context, log *zap.SugaredLogger, akoStreamInstance *akov2.AtlasStreamInstance) (ctrl.Result, error) {
	// check if stream instance is in "skipped" state
	if customresource.ReconciliationShouldBeSkipped(akoStreamInstance) {
		return r.skip(ctx, log, akoStreamInstance), nil
	}

	conditions := api.InitCondition(akoStreamInstance, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx, akoStreamInstance)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, akoStreamInstance)

	// check if stream instance is in "invalid" state
	isValid := customresource.ValidateResourceVersion(workflowCtx, akoStreamInstance, r.Log)
	if !isValid.IsOk() {
		return r.invalidate(isValid), nil
	}

	// check if stream instance is in "unsupported" state
	if !r.AtlasProvider.IsResourceSupported(akoStreamInstance) {
		return r.unsupport(workflowCtx), nil
	}

	project := akov2.AtlasProject{}
	if err := r.Client.Get(ctx, akoStreamInstance.AtlasProjectObjectKey(), &project); err != nil {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	atlasClient, orgID, err := r.AtlasProvider.SdkClient(workflowCtx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		return r.terminate(workflowCtx, workflow.AtlasAPIAccessNotConfigured, err)
	}
	workflowCtx.SdkClient = atlasClient
	workflowCtx.OrgID = orgID

	atlasStreamInstance, _, err := workflowCtx.SdkClient.StreamsApi.
		GetStreamInstance(workflowCtx.Context, project.ID(), akoStreamInstance.Spec.Name).
		Execute()

	if err != nil && !admin.IsErrorCode(err, instanceNotFound) {
		return r.terminate(workflowCtx, workflow.Internal, err)
	}

	isMarkedAsDeleted := !akoStreamInstance.GetDeletionTimestamp().IsZero()
	isNotInAtlas := err != nil && admin.IsErrorCode(err, instanceNotFound)

	switch {
	case isNotInAtlas && !isMarkedAsDeleted:
		// if no streams processing instance is not in atlas and is not marked as deleted - create
		// hence, create the stream instance and transition to "ready" state
		return r.create(workflowCtx, &project, akoStreamInstance)
	case isMarkedAsDeleted:
		// if a streams processing instance is marked as deleted,
		// independently whether it exists in Atlas or not - delete
		return r.delete(workflowCtx, &project, akoStreamInstance)
	case hasChanged(akoStreamInstance, atlasStreamInstance):
		// if a streams processing instance is ready and has changed - update
		return r.update(workflowCtx, &project, akoStreamInstance)
	}

	// handle connection registry management
	return r.handleConnectionRegistry(workflowCtx, &project, akoStreamInstance, atlasStreamInstance)
}

func hasChanged(streamInstance *akov2.AtlasStreamInstance, atlasStreamInstance *admin.StreamsTenant) bool {
	config := streamInstance.Spec.Config
	dataProcessRegion := atlasStreamInstance.GetDataProcessRegion()

	return config.Provider != dataProcessRegion.GetCloudProvider() || config.Region != dataProcessRegion.GetRegion()
}

func (r *AtlasStreamsInstanceReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasStreamInstance{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasStreamsInstanceReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasStreamInstance").
		For(r.For()).
		Watches(
			&akov2.AtlasStreamConnection{},
			handler.EnqueueRequestsFromMapFunc(r.findStreamInstancesForStreamConnection),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.findStreamInstancesForSecret),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func NewAtlasStreamsInstanceReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasStreamsInstanceReconciler {
	return &AtlasStreamsInstanceReconciler{
		Scheme:                   c.GetScheme(),
		Client:                   c.GetClient(),
		EventRecorder:            c.GetEventRecorderFor("AtlasStreamsInstance"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasStreamsInstance").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasStreamsInstanceReconciler) findStreamInstancesForStreamConnection(ctx context.Context, obj client.Object) []reconcile.Request {
	streamConnection, ok := obj.(*akov2.AtlasStreamConnection)
	if !ok {
		r.Log.Warnf("watching AtlasStreamConnection but got %T", obj)
		return nil
	}

	atlasStreamInstances := &akov2.AtlasStreamInstanceList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasStreamInstanceByConnectionIndex,
			client.ObjectKeyFromObject(streamConnection).String(),
		),
	}

	err := r.Client.List(ctx, atlasStreamInstances, listOps)
	if err != nil {
		r.Log.Errorf("failed to list Atlas stream instances: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(atlasStreamInstances.Items))
	for i := range atlasStreamInstances.Items {
		item := atlasStreamInstances.Items[i]
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

func (r *AtlasStreamsInstanceReconciler) findStreamInstancesForSecret(ctx context.Context, obj client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		r.Log.Warnf("watching Secret but got %T", obj)
		return nil
	}

	connections := &akov2.AtlasStreamConnectionList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasStreamConnectionBySecretIndex,
			client.ObjectKeyFromObject(secret).String(),
		),
	}

	err := r.Client.List(ctx, connections, listOps)
	if err != nil {
		r.Log.Errorf("failed to list Atlas stream connections: %e", err)
		return []reconcile.Request{}
	}

	if len(connections.Items) == 0 {
		return []reconcile.Request{}
	}

	streamInstancesMap := make(map[string]struct{}, len(connections.Items))
	streamInstances := make([]reconcile.Request, 0, len(connections.Items))
	for i := range connections.Items {
		requests := r.findStreamInstancesForStreamConnection(ctx, &connections.Items[i])
		for j := range requests {
			key := requests[j].String()
			if _, found := streamInstancesMap[key]; !found {
				streamInstances = append(streamInstances, requests[j])
				streamInstancesMap[key] = struct{}{}
			}
		}
	}

	return streamInstances
}
