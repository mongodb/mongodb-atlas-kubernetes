/*
Copyright 2021 MongoDB.

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

package atlasinventory

import (
	"context"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	dbaasv1beta1 "github.com/RHEcosystemAppEng/dbaas-operator/api/v1beta1"
	"go.mongodb.org/atlas/mongodbatlas"

	dbaas "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/dbaas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

// MongoDBAtlasInventoryReconciler reconciles a MongoDBAtlasInventory object
type MongoDBAtlasInventoryReconciler struct {
	Client      client.Client
	AtlasClient *mongodbatlas.Client
	watch.ResourceWatcher
	Log             *zap.SugaredLogger
	Scheme          *runtime.Scheme
	AtlasDomain     string
	GlobalAPISecret client.ObjectKey
	EventRecorder   record.EventRecorder
}

// Dev note: duplicate the permissions in both sections below to generate both Role and ClusterRoles

// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasinventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,resources=mongodbatlasinventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasinventories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbaas.redhat.com,namespace=default,resources=mongodbatlasinventories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch

func (r *MongoDBAtlasInventoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = ctx
	log := r.Log.With("MongoDBAtlasInventory", req.NamespacedName)
	log.Info("Reconciling MongoDBAtlasInventory")
	inventory := &dbaas.MongoDBAtlasInventory{}
	if err := r.Client.Get(ctx, req.NamespacedName, inventory); err != nil {
		if errors.IsNotFound(err) {
			// CR deleted since request queued, child objects getting GC'd, no requeue
			log.Info("MongoDBAtlasInventory resource not found, has been deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching MongoDBAtlasInventory for reconcile")
		return ctrl.Result{}, err
	}

	log.Infow("-> Starting MongoDBAtlasInventory reconciliation", "spec", inventory.Spec)

	// This update will make sure the status is always updated in case of any errors or successful result
	defer func(inv *dbaas.MongoDBAtlasInventory) {
		err := r.Client.Status().Update(ctx, inv)
		if err != nil {
			log.Infow("Could not update resource status:%v", err)
		}
	}(inventory)

	secretKey := inventory.ConnectionSecretObjectKey()

	if secretKey == nil {
		result := workflow.Terminate(workflow.MongoDBAtlasInventoryInputError, "Secret name for atlas credentials is missing")
		dbaas.SetInventoryCondition(inventory, dbaasv1beta1.DBaaSInventoryProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return result.ReconcileResult(), nil
	} else {
		// Note, that we are not watching the global connection secret - seems there is no point in reconciling all
		// the services once that secret is changed
		r.EnsureResourcesAreWatched(req.NamespacedName, "Secret", log, *secretKey)
	}

	atlasConn, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, inventory.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.MongoDBAtlasInventoryInputError, err.Error())
		dbaas.SetInventoryCondition(inventory, dbaasv1beta1.DBaaSInventoryProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return result.ReconcileResult(), nil
	}

	atlasClient := r.AtlasClient
	if atlasClient == nil {
		cl, err := atlas.Client(r.AtlasDomain, atlasConn, log)
		if err != nil {
			result := workflow.Terminate(workflow.MongoDBAtlasConnectionBackendError, err.Error())
			dbaas.SetInventoryCondition(inventory, dbaasv1beta1.DBaaSInventoryProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
			return result.ReconcileResult(), nil
		}
		atlasClient = &cl
	}

	inventoryList, result := discoverInstances(atlasClient)
	if !result.IsOk() {
		dbaas.SetInventoryCondition(inventory, dbaasv1beta1.DBaaSInventoryProviderSyncType, metav1.ConditionFalse, string(result.Reason()), result.Message())
		return result.ReconcileResult(), nil
	}

	// Update the status
	dbaas.SetInventoryCondition(inventory, dbaasv1beta1.DBaaSInventoryProviderSyncType, metav1.ConditionTrue, string(workflow.MongoDBAtlasInventorySyncOK), "Spec sync OK")
	inventory.Status.DatabaseServices = inventoryList
	return ctrl.Result{}, nil
}

// Delete is a no-op
func (r *MongoDBAtlasInventoryReconciler) Delete(e event.DeleteEvent) error {
	return nil
}

func (r *MongoDBAtlasInventoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("MongoDBAtlasInventory", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MongoDBAtlasInventory & handle delete separately
	err = c.Watch(&source.Kind{Type: &dbaas.MongoDBAtlasInventory{}},
		&watch.EventHandlerWithDelete{Controller: r},
		watch.CommonPredicates())
	if err != nil {
		return err
	}

	// Watch for changes to other resource MongoDBAtlasInstance
	err = c.Watch(&source.Kind{Type: &dbaas.MongoDBAtlasInstance{}},
		handler.EnqueueRequestsFromMapFunc(instanceMapFunc))
	if err != nil {
		return err
	}

	// Watch for Connection Secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources))
	if err != nil {
		return err
	}
	return nil
}

// instanceMapFunc defines a function for EnqueueRequestsFromMapFunc so that when a MongoDBAtlasInstance
// CR status is updated or the CR is deleted, the corresponding inventory is enqueued in order to refresh
// the instance list
func instanceMapFunc(a client.Object) []ctrl.Request {
	if instance, ok := a.(*dbaas.MongoDBAtlasInstance); ok {
		return []ctrl.Request{
			{NamespacedName: types.NamespacedName{
				Name:      instance.Spec.InventoryRef.Name,
				Namespace: instance.Spec.InventoryRef.Namespace,
			},
			}}
	}
	return nil
}
