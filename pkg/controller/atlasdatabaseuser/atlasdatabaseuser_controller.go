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
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/kube"
)

// AtlasDatabaseUserReconciler reconciles an AtlasDatabaseUser object
type AtlasDatabaseUserReconciler struct {
	watch.ResourceWatcher
	Client      client.Client
	Log         *zap.SugaredLogger
	Scheme      *runtime.Scheme
	AtlasDomain string
	OperatorPod client.ObjectKey
}

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;update;patch;delete

// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasdatabaseusers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",namespace=default,resources=secrets,verbs=create;update;patch;delete

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
	defer statushandler.Update(ctx, r.Client, databaseUser)

	project := &mdbv1.AtlasProject{}
	if result := r.readProjectResource(databaseUser, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.DatabaseUserReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.OperatorPod, project.ConnectionSecretObjectKey())
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
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasDatabaseUser{}}, &watch.EventHandlerWithDelete{Controller: r}, watch.CommonPredicates())
	if err != nil {
		return err
	}

	// Watch for DatabaseUser password Secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources))
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

	connection, err := atlas.ReadConnection(log, r.Client, r.OperatorPod, project.ConnectionSecretObjectKey())
	if err != nil {
		return err
	}

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		return fmt.Errorf("cannot build Atlas client: %w", err)
	}

	userName := dbUser.Spec.Username
	_, err = atlasClient.DatabaseUsers.Delete(context.Background(), dbUser.Spec.DatabaseName, project.ID(), userName)
	if err != nil {
		return fmt.Errorf("cannot delete Database User in Atlas: %w", err)
	}

	log.Infow("Started DatabaseUser deletion process in Atlas", "projectID", project.ID(), "userName", userName)

	secrets, err := connectionsecret.ListByUserName(r.Client, dbUser.Namespace, project.ID(), userName)
	if err != nil {
		return fmt.Errorf("failed to find connection secrets for the user: %w", err)
	}

	for _, secret := range secrets {
		// Solves the "Implicit memory aliasing in for loop" linter error
		s := secret.DeepCopy()
		err = r.Client.Delete(context.Background(), s)
		if err != nil {
			log.Errorf("Failed to remove connection Secret: %v", err)
		} else {
			log.Debugw("Removed connection Secret", "secret", kube.ObjectKeyFromObject(s))
		}
	}
	if len(secrets) > 0 {
		log.Infof("Removed %d connection secrets", len(secrets))
	}

	return nil
}
