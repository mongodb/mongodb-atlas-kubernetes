package atlasdatafederation

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"

	"sigs.k8s.io/controller-runtime/pkg/handler"

	"go.mongodb.org/atlas/mongodbatlas"

	ctrl "sigs.k8s.io/controller-runtime"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasDataFederationReconciler reconciles an DataFederation object
type AtlasDataFederationReconciler struct {
	watch.ResourceWatcher
	Client           client.Client
	Log              *zap.SugaredLogger
	Scheme           *runtime.Scheme
	AtlasDomain      string
	GlobalAPISecret  client.ObjectKey
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatafederations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatafederations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *AtlasDataFederationReconciler) Reconcile(contextInt context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasdatafederation", req.NamespacedName)

	dataFederation := &mdbv1.AtlasDataFederation{}
	result := customresource.PrepareResource(r.Client, req, dataFederation, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if shouldSkip := customresource.ReconciliationShouldBeSkipped(dataFederation); shouldSkip {
		log.Infow(fmt.Sprintf("-> Skipping AtlasDataFederation reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", dataFederation.Spec)
		if !dataFederation.GetDeletionTimestamp().IsZero() {
			err := r.removeDeletionFinalizer(contextInt, dataFederation)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	ctx := customresource.MarkReconciliationStarted(r.Client, dataFederation, log)
	log.Infow("-> Starting AtlasDataFederation reconciliation", "spec", dataFederation.Spec, "status", dataFederation.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, dataFederation)

	resourceVersionIsValid := customresource.ValidateResourceVersion(ctx, dataFederation, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("AtlasDataFederation validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(contextInt, dataFederation, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		if errRm := r.removeDeletionFinalizer(contextInt, dataFederation); errRm != nil {
			result = workflow.Terminate(workflow.Internal, errRm.Error())
			return result.ReconcileResult(), nil
		}
		return result.ReconcileResult(), nil
	}
	ctx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Client = atlasClient

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
			err = r.Client.Get(contextInt, kube.ObjectKeyFromObject(dataFederation), dataFederation)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				return result.ReconcileResult(), nil
			}
			customresource.SetFinalizer(dataFederation, customresource.FinalizerLabel)
			if err = r.Client.Update(contextInt, dataFederation); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to add finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
	}

	if !dataFederation.GetDeletionTimestamp().IsZero() {
		if customresource.HaveFinalizer(dataFederation, customresource.FinalizerLabel) {
			if customresource.ResourceShouldBeLeftInAtlas(dataFederation) {
				log.Infof("Not removing AtlasDataFederation from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
			} else {
				if err = r.deleteDataFederationFromAtlas(contextInt, &atlasClient, dataFederation, project, log); err != nil {
					log.Errorf("failed to remove DataFederation from Atlas: %s", err)
					result = workflow.Terminate(workflow.Internal, err.Error())
					ctx.SetConditionFromResult(status.DataFederationReadyType, result)
					return result.ReconcileResult(), nil
				}
			}
			if err = r.removeDeletionFinalizer(contextInt, dataFederation); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		} else {
			return result.ReconcileResult(), nil
		}
	}

	ctx.SetConditionTrue(status.ReadyType)
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasDataFederationReconciler) removeDeletionFinalizer(context context.Context, df *mdbv1.AtlasDataFederation) error {
	err := r.Client.Get(context, kube.ObjectKeyFromObject(df), df)
	if err != nil {
		return fmt.Errorf("cannot get AtlasDeployment while adding finalizer: %w", err)
	}

	customresource.UnsetFinalizer(df, customresource.FinalizerLabel)
	if err = r.Client.Update(context, df); err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", df.GetName(), err)
	}
	return nil
}

func (r *AtlasDataFederationReconciler) deleteDataFederationFromAtlas(ctx context.Context, client *mongodbatlas.Client, df *mdbv1.AtlasDataFederation, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	log.Infof("Deleting DataFederation instance: %s from Atlas", df.Spec.Name)

	_, err := client.DataFederation.Delete(ctx, project.ID(), df.Spec.Name)

	var apiError *mongodbatlas.ErrorResponse
	if errors.As(err, &apiError) && apiError.Error() == "DATA_LAKE_TENANT_NOT_FOUND_FOR_NAME" {
		log.Info("DataFederation doesn't exist or is already deleted")
		return nil
	}

	if err != nil {
		log.Errorw("Can not delete Atlas data federation", "error", err)
		return err
	}

	return nil
}

func (r *AtlasDataFederationReconciler) readProjectResource(ctx context.Context, dataFederation *mdbv1.AtlasDataFederation, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(ctx, dataFederation.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("AtlasDataFederation", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DataFederation & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDataFederation{}}, &watch.EventHandlerWithDelete{Controller: r}, r.GlobalPredicates...)
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDataFederation{}}, &handler.EnqueueRequestForObject{}, r.GlobalPredicates...)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements a handler for the Delete event
func (r *AtlasDataFederationReconciler) Delete(e event.DeleteEvent) error {
	dataFederation, ok := e.Object.(*mdbv1.AtlasDataFederation)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasDeployment{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasdatafederation", kube.ObjectKeyFromObject(dataFederation))

	log.Infow("-> Starting AtlasDataFederation deletion", "spec", dataFederation.Spec)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(context.Background(), dataFederation, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	log = log.With("projectID", project.Status.ID, "dataFederationName", dataFederation.Spec.Name)

	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(r.Client, dataFederation.Namespace, project.ID(), dataFederation.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(context.Background(), &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
	}

	return nil
}
