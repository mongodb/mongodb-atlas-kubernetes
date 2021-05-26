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
	"time"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	Client client.Client
	watch.ResourceWatcher
	Log             *zap.SugaredLogger
	Scheme          *runtime.Scheme
	AtlasDomain     string
	GlobalAPISecret client.ObjectKey
	EventRecorder   record.EventRecorder
}

// Dev note: duplicate the permissions in both sections below to generate both Role and ClusterRoles

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasprojects/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasProjectReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context
	log := r.Log.With("atlasproject", req.NamespacedName)

	project := &mdbv1.AtlasProject{}
	result := customresource.PrepareResource(r.Client, req, project, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	if project.ConnectionSecretObjectKey() != nil {
		// Note, that we are not watching the global connection secret - seems there is no point in reconciling all
		// the projects once that secret is changed
		r.EnsureResourcesAreWatched(req.NamespacedName, "Secret", log, *project.ConnectionSecretObjectKey())
	}
	ctx := customresource.MarkReconciliationStarted(r.Client, project, log)

	log.Infow("-> Starting AtlasProject reconciliation", "spec", project.Spec)

	// This update will make sure the status is always updated in case of any errors or successful result
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, project)

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result = workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		ctx.SetConditionFromResult(status.ClusterReadyType, workflow.Terminate(workflow.Internal, err.Error()))
		return result.ReconcileResult(), nil
	}
	ctx.Client = atlasClient

	var projectID string
	if projectID, result = r.ensureProjectExists(ctx, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))

	// Updating the status with "projectReady = true" and "IPAccessListReady = false" (not as separate updates!)
	ctx.SetConditionTrue(status.ProjectReadyType)
	r.EventRecorder.Event(project, "Normal", string(status.ProjectReadyType), "")

	if result = r.ensureIPAccessList(ctx, projectID, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.IPAccessListReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.SetConditionTrue(status.IPAccessListReadyType)
	r.EventRecorder.Event(project, "Normal", string(status.IPAccessListReadyType), "")

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

	if project.Status.ID == "" {
		log.Infof("Project does not exist in Atlas, nothing to remove")
		return nil
	}

	if customresource.ResourceShouldBeLeftInAtlas(project) {
		log.Infof("Not removing the Atlas Project from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
		return nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	go func() {
		timeout := time.Now().Add(workflow.DefaultTimeout)

		for time.Now().Before(timeout) {
			_, err = atlasClient.Projects.Delete(context.Background(), project.Status.ID)
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && apiError.ErrorCode == atlas.NotInGroup {
				log.Infow("Project is already deleted", "projectID", project.Status.ID)
				return
			}

			if err != nil {
				log.Errorw("cannot delete Atlas project", "error", err)
				time.Sleep(workflow.DefaultRetry)
				continue
			}

			log.Infow("Successfully deleted Atlas project", "projectID", project.Status.ID)
			return
		}

		log.Errorw("Failed to delete Atlas project in time", "projectID", project.Status.ID)
	}()

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
