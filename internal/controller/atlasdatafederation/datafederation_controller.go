package atlasdatafederation

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/indexer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

// AtlasDataFederationReconciler reconciles an DataFederation object
type AtlasDataFederationReconciler struct {
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *AtlasDataFederationReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdatafederation", req.NamespacedName)

	dataFederation := &akov2.AtlasDataFederation{}
	result := customresource.PrepareResource(context, r.Client, req, dataFederation, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(dataFederation) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDataFederation reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", dataFederation.Spec)
		if !dataFederation.GetDeletionTimestamp().IsZero() {
			err := customresource.ManageFinalizer(context, r.Client, dataFederation, customresource.UnsetFinalizer)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err)
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	conditions := akov2.InitCondition(dataFederation, api.FalseCondition(api.ReadyType))
	ctx := workflow.NewContext(log, conditions, context, dataFederation)
	log.Infow("-> Starting AtlasDataFederation reconciliation", "spec", dataFederation.Spec, "status", dataFederation.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, dataFederation)

	resourceVersionIsValid := customresource.ValidateResourceVersion(ctx, dataFederation, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("AtlasDataFederation validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if !r.AtlasProvider.IsResourceSupported(dataFederation) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, errors.New("the AtlasDataFederation is not supported by Atlas for government")).
			WithoutRetry()
		ctx.SetConditionFromResult(api.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	project := &akov2.AtlasProject{}
	if result := r.readProjectResource(context, dataFederation, project); !result.IsOk() {
		ctx.SetConditionFromResult(api.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	endpointService, err := datafederation.NewDatafederationPrivateEndpointService(ctx.Context, r.AtlasProvider, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err)
		ctx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	dataFederationService, err := datafederation.NewAtlasDataFederationService(ctx.Context, r.AtlasProvider, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result = workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err)
		ctx.SetConditionFromResult(api.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result = r.ensureDataFederation(ctx, project, dataFederation, dataFederationService); !result.IsOk() {
		ctx.SetConditionFromResult(api.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result = r.ensurePrivateEndpoints(ctx, endpointService, project, dataFederation); !result.IsOk() {
		ctx.SetConditionFromResult(api.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result = r.ensureConnectionSecrets(ctx, dataFederationService, project, dataFederation); !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if dataFederation.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(dataFederation, customresource.FinalizerLabel) {
			err = r.Client.Get(context, kube.ObjectKeyFromObject(dataFederation), dataFederation)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err)
				return result.ReconcileResult(), nil
			}
			customresource.SetFinalizer(dataFederation, customresource.FinalizerLabel)
			if err = r.Client.Update(context, dataFederation); err != nil {
				result = workflow.Terminate(workflow.Internal, err)
				log.Errorw("failed to add finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
	}

	if !dataFederation.GetDeletionTimestamp().IsZero() {
		return r.handleDelete(ctx, log, dataFederation, project, dataFederationService).ReconcileResult(), nil
	}

	err = customresource.ApplyLastConfigApplied(context, dataFederation, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err)
		ctx.SetConditionFromResult(api.DataFederationReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	ctx.SetConditionTrue(api.DataFederationReadyType)
	ctx.SetConditionTrue(api.ReadyType)
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasDataFederationReconciler) handleDelete(ctx *workflow.Context, log *zap.SugaredLogger, dataFederation *akov2.AtlasDataFederation, project *akov2.AtlasProject, service datafederation.DataFederationService) workflow.Result {
	if customresource.HaveFinalizer(dataFederation, customresource.FinalizerLabel) {
		if customresource.IsResourcePolicyKeepOrDefault(dataFederation, r.ObjectDeletionProtection) {
			log.Info("Not removing AtlasDataFederation from Atlas as per configuration")
		} else {
			if err := r.deleteConnectionSecrets(ctx.Context, dataFederation); err != nil {
				log.Errorf("failed to remove DataFederation connection secrets from Atlas: %s", err)
				result := workflow.Terminate(workflow.Internal, err)
				ctx.SetConditionFromResult(api.DataFederationReadyType, result)
				return result
			}
			if err := r.deleteDataFederationFromAtlas(ctx.Context, service, dataFederation, project, log); err != nil {
				log.Errorf("failed to remove DataFederation from Atlas: %s", err)
				result := workflow.Terminate(workflow.Internal, err)
				ctx.SetConditionFromResult(api.DataFederationReadyType, result)
				return result
			}
		}
		if err := customresource.ManageFinalizer(ctx.Context, r.Client, dataFederation, customresource.UnsetFinalizer); err != nil {
			result := workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err)
			log.Errorw("failed to remove finalizer", "error", err)
			return result
		}
	}

	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) deleteDataFederationFromAtlas(ctx context.Context, service datafederation.DataFederationService, df *akov2.AtlasDataFederation, project *akov2.AtlasProject, log *zap.SugaredLogger) error {
	log.Infof("Deleting DataFederation instance: %s from Atlas", df.Spec.Name)

	err := service.Delete(ctx, project.ID(), df.Spec.Name)

	if errors.Is(err, datafederation.ErrorNotFound) {
		log.Info("DataFederation doesn't exist or is already deleted")
		return nil
	}

	if err != nil {
		log.Errorw("Can not delete Atlas DataFederation", "error", err)
		return err
	}

	return nil
}

func (r *AtlasDataFederationReconciler) readProjectResource(ctx context.Context, dataFederation *akov2.AtlasDataFederation, project *akov2.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, dataFederation.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}
	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) For() (client.Object, builder.Predicates) {
	return &akov2.AtlasDataFederation{}, builder.WithPredicates(r.GlobalPredicates...)
}

func (r *AtlasDataFederationReconciler) SetupWithManager(mgr ctrl.Manager, skipNameValidation bool) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDataFederation").
		For(r.For()).
		Watches(
			&akov2.AtlasProject{},
			handler.EnqueueRequestsFromMapFunc(r.findAtlasDataFederationForProjects),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		WithOptions(controller.TypedOptions[reconcile.Request]{SkipNameValidation: pointer.MakePtr(skipNameValidation)}).
		Complete(r)
}

func (r *AtlasDataFederationReconciler) findAtlasDataFederationForProjects(ctx context.Context, obj client.Object) []reconcile.Request {
	project, ok := obj.(*akov2.AtlasProject)
	if !ok {
		r.Log.Warnf("watching AtlasProject but got %T", obj)
		return nil
	}

	datafederationList := &akov2.AtlasDataFederationList{}
	listOps := &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector(
			indexer.AtlasDataFederationByProject,
			client.ObjectKeyFromObject(project).String(),
		),
	}
	err := r.Client.List(ctx, datafederationList, listOps)
	if err != nil {
		r.Log.Errorf("failed to list AtlasDataFederation: %e", err)
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(datafederationList.Items))
	for _, item := range datafederationList.Items {
		requests = append(requests, reconcile.Request{NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace}})
	}
	return requests
}

func NewAtlasDataFederationReconciler(
	c cluster.Cluster,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasDataFederationReconciler {
	return &AtlasDataFederationReconciler{
		Scheme:                   c.GetScheme(),
		Client:                   c.GetClient(),
		EventRecorder:            c.GetEventRecorderFor("AtlasDataFederation"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasDataFederation").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func (r *AtlasDataFederationReconciler) deleteConnectionSecrets(ctx context.Context, dataFederation *akov2.AtlasDataFederation) error {
	log := r.Log.With("atlasdatafederation", kube.ObjectKeyFromObject(dataFederation))

	project := &akov2.AtlasProject{}

	if result := r.readProjectResource(ctx, dataFederation, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	log = log.With("projectID", project.Status.ID, "dataFederationName", dataFederation.Spec.Name)

	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(ctx, r.Client, dataFederation.Namespace, project.ID(), dataFederation.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(ctx, &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
	}

	return nil
}
