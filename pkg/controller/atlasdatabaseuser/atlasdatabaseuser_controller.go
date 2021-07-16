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

package atlasdatabaseuser

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

// AtlasDatabaseUserReconciler reconciles an AtlasDatabaseUser object
type AtlasDatabaseUserReconciler struct {
	watch.ResourceWatcher
	Client           client.Client
	Log              *zap.SugaredLogger
	Scheme           *runtime.Scheme
	AtlasDomain      string
	GlobalAPISecret  client.ObjectKey
	EventRecorder    record.EventRecorder
	GlobalPredicates []predicate.Predicate
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=create;update;patch;delete
// +kubebuilder:rbac:groups="",namespace=default,resources=events,verbs=create;patch

func (r *AtlasDatabaseUserReconciler) Reconcile(context context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = context
	log := r.Log.With("atlasdatabaseuser", req.NamespacedName)

	databaseUser := &mdbv1.AtlasDatabaseUser{}
	result := customresource.PrepareResource(r.Client, req, databaseUser, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	if databaseUser.Spec.PasswordSecret != nil {
		r.EnsureResourcesAreWatched(req.NamespacedName, "Secret", log, *databaseUser.PasswordSecretObjectKey())
	}
	ctx := customresource.MarkReconciliationStarted(r.Client, databaseUser, log)

	log.Infow("-> Starting AtlasDatabaseUser reconciliation", "spec", databaseUser.Spec, "status", databaseUser.Status)
	defer statushandler.Update(ctx, r.Client, r.EventRecorder, databaseUser)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(databaseUser, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result := workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		ctx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		ctx.SetConditionFromResult(status.ClusterReadyType, result)
		return result.ReconcileResult(), nil
	}
	ctx.Client = atlasClient

	result = r.ensureDatabaseUser(ctx, *project, *databaseUser)
	if !result.IsOk() {
		ctx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	ctx.SetConditionTrue(status.DatabaseUserReadyType)
	ctx.SetConditionTrue(status.ReadyType)
	return result.ReconcileResult(), nil
}

func (r *AtlasDatabaseUserReconciler) readProjectResource(user *mdbv1.AtlasDatabaseUser, project *mdbv1.AtlasProject) workflow.Result {
	if err := r.Client.Get(context.Background(), user.AtlasProjectObjectKey(), project); err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}
	return workflow.OK()
}

func (r *AtlasDatabaseUserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("AtlasDatabaseUser", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AtlasDatabaseUser & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDatabaseUser{}}, &watch.EventHandlerWithDelete{Controller: r}, r.GlobalPredicates...)
	if err != nil {
		return err
	}

	// Watch for DatabaseUser password Secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources), r.GlobalPredicates...)
	if err != nil {
		return err
	}

	return nil
}

func (r AtlasDatabaseUserReconciler) Delete(e event.DeleteEvent) error {
	dbUser, ok := e.Object.(*mdbv1.AtlasDatabaseUser)
	if !ok {
		r.Log.Errorf("Ignoring malformed Delete() call (expected type %T, got %T)", &mdbv1.AtlasDatabaseUser{}, e.Object)
		return nil
	}

	log := r.Log.With("atlasdatabaseuser", kube.ObjectKeyFromObject(dbUser))

	log.Infow("-> Starting AtlasDatabaseUser deletion", "spec", dbUser.Spec)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(dbUser, project); !result.IsOk() {
		return errors.New("cannot read project resource")
	}

	if customresource.ResourceShouldBeLeftInAtlas(dbUser) {
		log.Infof("Not removing Atlas database user from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
	} else if err := r.deleteUserFromAtlas(dbUser, project, log); err != nil {
		log.Error("Failed to remove database user from Atlas: %s", err)
	}

	// We ignore the error as it will be printed by the function
	_ = removeStaleSecretsByUserName(r.Client, project.ID(), dbUser.Spec.Username, *dbUser, log)

	return nil
}

func (r AtlasDatabaseUserReconciler) deleteUserFromAtlas(dbUser *mdbv1.AtlasDatabaseUser, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	userName := dbUser.Spec.Username

	go func() {
		timeout := time.Now().Add(workflow.DefaultTimeout)

		for time.Now().Before(timeout) {
			_, err = atlasClient.DatabaseUsers.Delete(context.Background(), dbUser.Spec.DatabaseName, project.ID(), userName)
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && apiError.ErrorCode == atlas.UsernameNotFound {
				log.Info("Database user doesn't exist or is already deleted")
				return
			}

			if err != nil {
				log.Errorw("Cannot delete Atlas database user", "error", err)
				time.Sleep(workflow.DefaultRetry)
				continue
			}

			log.Infow("Started DatabaseUser deletion process in Atlas", "projectID", project.ID(), "userName", userName)
			return
		}
	}()

	return nil
}
