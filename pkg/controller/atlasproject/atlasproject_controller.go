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

package atlasproject

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	Client client.Client
	watch.ResourceWatcher
	Log         *zap.SugaredLogger
	Scheme      *runtime.Scheme
	AtlasDomain string
	OperatorPod client.ObjectKey
}

// Dev note: duplicate the permissions in both sections below to generate both Role and ClusterRoles

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch

func (r *AtlasProjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasproject", req.NamespacedName)

	project := &mdbv1.AtlasProject{}
	result := customresource.PrepareResource(r.Client, req, project, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	if project.ConnectionSecretObjectKey() != nil {
		r.EnsureResourcesAreWatched(req.NamespacedName, "Secret", log, *project.ConnectionSecretObjectKey())
		// TODO CLOUDP-80516: the "global" connection secret also needs to be watched
	}
	ctx := customresource.MarkReconciliationStarted(r.Client, project, log)

	log.Infow("-> Starting AtlasProject reconciliation", "spec", project.Spec)

	// This update will make sure the status is always updated in case of any errors or successful result
	defer statushandler.Update(ctx, r.Client, project)

	connection, result := atlas.ReadConnection(log, r.Client, r.OperatorPod, project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		// merge result into ctx
		ctx.SetConditionFromResult(status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}

	var projectID string
	if projectID, result = r.ensureProjectExists(ctx, connection, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))

	// Updating the status with "projectReady = true" and "IPAccessListReady = false" (not as separate updates!)
	ctx.SetConditionTrue(status.ProjectReadyType)

	if result = r.ensureIPAccessList(ctx, connection, projectID, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.IPAccessListReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.SetConditionTrue(status.IPAccessListReadyType)
	ctx.SetConditionTrue(status.ReadyType)
	return ctrl.Result{}, nil
}

func (r *AtlasProjectReconciler) Delete(e event.DeleteEvent) error {
	project, ok := e.Object.(*mdbv1.AtlasProject)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasProject{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasproject", kube.ObjectKeyFromObject(project))

	log.Infow("-> Starting AtlasProject deletion", "spec", project.Spec)

	connection, result := atlas.ReadConnection(log, r.Client, r.OperatorPod, project.ConnectionSecretObjectKey())
	if !result.IsOk() {
		return errors.New("cannot read Atlas connection")
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	_, err = atlasClient.Projects.Delete(context.Background(), project.Status.ID)
	if err != nil {
		return fmt.Errorf("cannot delete Atlas project: %w", err)
	}

	log.Infow("Successfully deleted Atlas project", "projectID", project.Status.ID)

	return nil
}

func (r *AtlasProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("AtlasProject", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AtlasProject & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasProject{}}, &watch.EventHandlerWithDelete{Controller: r}, watch.CommonPredicates())
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
