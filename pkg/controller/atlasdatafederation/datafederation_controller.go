package atlasdatafederation

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"

	"go.mongodb.org/atlas/mongodbatlas"

	ctrl "sigs.k8s.io/controller-runtime"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/kube"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

// AtlasDataFederationReconciler reconciles an DataFederation object
type AtlasDataFederationReconciler struct {
	watch.ResourceWatcher
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
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	ctx := customresource.MarkReconciliationStarted(r.Client, dataFederation, log, context)
	log.Infow("-> Starting AtlasDataFederation reconciliation", "spec", dataFederation.Spec, "status", dataFederation.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, dataFederation)

	resourceVersionIsValid := customresource.ValidateResourceVersion(ctx, dataFederation, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("AtlasDataFederation validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if !r.AtlasProvider.IsResourceSupported(dataFederation) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasDataFederation is not supported by Atlas for government").
			WithoutRetry()
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	project := &akov2.AtlasProject{}
	if result := r.readProjectResource(context, dataFederation, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	atlasClient, orgID, err := r.AtlasProvider.Client(ctx.Context, project.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.OrgID = orgID
	ctx.Client = atlasClient

	// Setting protection flag to static false because ownership detection is disabled.
	owner, err := customresource.IsOwner(dataFederation, false, customresource.IsResourceManagedByOperator, managedByAtlas(context, atlasClient, project.ID(), log))
	if err != nil {
		result = workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	if !owner {
		result = workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile DataFederation: it already exists in Atlas, it was not previously managed by the operator, and the deletion protection is enabled.",
		)
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	if result = r.ensureDataFederation(ctx, project, dataFederation); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result = r.ensurePrivateEndpoints(ctx, project, dataFederation); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationPEReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result = r.ensureConnectionSecrets(ctx, project, dataFederation); !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if dataFederation.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(dataFederation, customresource.FinalizerLabel) {
			err = r.Client.Get(context, kube.ObjectKeyFromObject(dataFederation), dataFederation)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				return result.ReconcileResult(), nil
			}
			customresource.SetFinalizer(dataFederation, customresource.FinalizerLabel)
			if err = r.Client.Update(context, dataFederation); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to add finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
	}

	if !dataFederation.GetDeletionTimestamp().IsZero() {
		if customresource.HaveFinalizer(dataFederation, customresource.FinalizerLabel) {
			if customresource.IsResourcePolicyKeepOrDefault(dataFederation, r.ObjectDeletionProtection) {
				log.Info("Not removing AtlasDataFederation from Atlas as per configuration")
			} else {
				if err = r.deleteDataFederationFromAtlas(context, atlasClient, dataFederation, project, log); err != nil {
					log.Errorf("failed to remove DataFederation from Atlas: %s", err)
					result = workflow.Terminate(workflow.Internal, err.Error())
					ctx.SetConditionFromResult(status.DataFederationReadyType, result)
					return result.ReconcileResult(), nil
				}
			}
			if err = customresource.ManageFinalizer(context, r.Client, dataFederation, customresource.UnsetFinalizer); err != nil {
				result = workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		} else {
			return result.ReconcileResult(), nil
		}
	}

	err = customresource.ApplyLastConfigApplied(context, project, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	ctx.SetConditionTrue(status.ReadyType)
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasDataFederationReconciler) deleteDataFederationFromAtlas(ctx context.Context, client *mongodbatlas.Client, df *akov2.AtlasDataFederation, project *akov2.AtlasProject, log *zap.SugaredLogger) error {
	log.Infof("Deleting DataFederation instance: %s from Atlas", df.Spec.Name)

	_, err := client.DataFederation.Delete(ctx, project.ID(), df.Spec.Name)

	var apiError *mongodbatlas.ErrorResponse
	if errors.As(err, &apiError) && apiError.Error() == "DATA_LAKE_TENANT_NOT_FOUND_FOR_NAME" {
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
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasDataFederation").
		Watches(&akov2.AtlasDataFederation{}, &watch.EventHandlerWithDelete{Controller: r}, builder.WithPredicates(r.GlobalPredicates...)).
		For(&akov2.AtlasDataFederation{}, builder.WithPredicates(r.GlobalPredicates...)).
		Complete(r)
}

// Delete implements a handler for the Delete event
func (r *AtlasDataFederationReconciler) Delete(ctx context.Context, e event.DeleteEvent) error {
	dataFederation, ok := e.Object.(*akov2.AtlasDataFederation)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &akov2.AtlasDeployment{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasdatafederation", kube.ObjectKeyFromObject(dataFederation))

	log.Infow("-> Starting AtlasDataFederation deletion", "spec", dataFederation.Spec)

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

func managedByAtlas(ctx context.Context, atlasClient *mongodbatlas.Client, projectID string, log *zap.SugaredLogger) customresource.AtlasChecker {
	return func(resource akov2.AtlasCustomResource) (bool, error) {
		dataFederation, ok := resource.(*akov2.AtlasDataFederation)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasDataFederation")
		}

		atlasDataFederation, _, err := atlasClient.DataFederation.Get(ctx, projectID, dataFederation.Spec.Name)
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.DataFederationTenantNotFound || apiError.ErrorCode == atlas.ResourceNotFound) {
				return false, nil
			}
			return false, err
		}

		isSame, err := dataFederationMatchesSpec(log, atlasDataFederation, dataFederation)
		if err != nil {
			return true, nil
		}
		return !isSame, nil
	}
}
