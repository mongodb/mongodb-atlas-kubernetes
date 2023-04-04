package atlasdatafederation

import (
	"context"
	"errors"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
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

	ctx := customresource.MarkReconciliationStarted(r.Client, dataFederation, log)
	log.Infow("-> Starting AtlasDataFederation reconciliation", "spec", dataFederation.Spec, "status", dataFederation.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, dataFederation)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(contextInt, dataFederation, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
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

	if result := r.ensureDataFederation(ctx, project, dataFederation); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationReadyType, result)
		return result.ReconcileResult(), nil
	}

	if result := r.ensurePrivateEndpoints(ctx, project, dataFederation); !result.IsOk() {
		ctx.SetConditionFromResult(status.DataFederationPEReadyType, result)
		return result.ReconcileResult(), nil
	}

	ctx.SetConditionTrue(status.ReadyType)
	return workflow.OK().ReconcileResult(), nil
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

	// Watch for Backup schedules
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupSchedule{}}, watch.NewBackupScheduleHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	// Watch for Backup policies
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasBackupPolicy{}}, watch.NewBackupPolicyHandler(r.WatchedResources))
	if err != nil {
		return err
	}

	return nil
}

// Delete implements a handler for the Delete event.
func (r *AtlasDataFederationReconciler) Delete(e event.DeleteEvent) error {
	dataFederation, ok := e.Object.(*mdbv1.AtlasDataFederation)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasDeployment{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasdeployment", kube.ObjectKeyFromObject(dataFederation))

	log.Infow("-> Starting AtlasDeployment deletion", "spec", dataFederation.Spec)

	context := context.Background()
	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(context, dataFederation, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	log = log.With("projectID", project.Status.ID, "deploymentName", dataFederation.Spec.Name)

	// We always remove the connection secrets even if the deployment is not removed from Atlas
	secrets, err := connectionsecret.ListByDeploymentName(r.Client, dataFederation.Namespace, project.ID(), dataFederation.Spec.Name)
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for i := range secrets {
		if err := r.Client.Delete(context, &secrets[i]); err != nil {
			if k8serrors.IsNotFound(err) {
				continue
			}
			log.Errorw("Failed to delete secret", "secretName", secrets[i].Name, "error", err)
		}
	}

	return nil
}
