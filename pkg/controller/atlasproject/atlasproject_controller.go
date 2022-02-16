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

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	Log              *zap.SugaredLogger
	Scheme           *runtime.Scheme
	AtlasDomain      string
	GlobalAPISecret  client.ObjectKey
	GlobalPredicates []predicate.Predicate
	EventRecorder    record.EventRecorder
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

	if project.GetDeletionTimestamp().IsZero() {
		if !isDeletionFinalizerPresent(project) {
			log.Debugw("Add deletion finalizer", "name", getFinalizerName())
			if err := r.addDeletionFinalizer(context, project); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				ctx.SetConditionFromResult(status.ClusterReadyType, result)
				return result.ReconcileResult(), nil
			}
		}
	}

	if !project.GetDeletionTimestamp().IsZero() {
		if isDeletionFinalizerPresent(project) {
			if customresource.ResourceShouldBeLeftInAtlas(project) {
				log.Infof("Not removing the Atlas Project from Atlas as the '%s' annotation is set", customresource.ResourcePolicyAnnotation)
			} else {
				if result = DeleteAllPrivateEndpoints(ctx, atlasClient, projectID, project.Status.PrivateEndpoints, log); !result.IsOk() {
					ctx.SetConditionFromResult(status.PrivateEndpointReadyType, result)
					return result.ReconcileResult(), nil
				}

				if err = r.deleteAtlasProject(context, atlasClient, project); err != nil {
					result = workflow.Terminate(workflow.Internal, err.Error())
					ctx.SetConditionFromResult(status.ClusterReadyType, result)
					return result.ReconcileResult(), nil
				}
			}

			if err = r.removeDeletionFinalizer(context, project); err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				ctx.SetConditionFromResult(status.ClusterReadyType, result)
				return result.ReconcileResult(), nil
			}
		}

		return result.ReconcileResult(), nil
	}

	// Updating the status with "projectReady = true" and "IPAccessListReady = false" (not as separate updates!)
	ctx.SetConditionTrue(status.ProjectReadyType)
	r.EventRecorder.Event(project, "Normal", string(status.ProjectReadyType), "")

	if result = ensureIPAccessList(ctx, projectID, project); !result.IsOk() {
		ctx.SetConditionFromResult(status.IPAccessListReadyType, result)
		return result.ReconcileResult(), nil
	}

	allReady, result := r.allIPAccessListsAreReady(context, ctx, projectID)
	if !result.IsOk() {
		ctx.SetConditionFalse(status.IPAccessListReadyType)
		return result.ReconcileResult(), nil
	}

	if allReady {
		ctx.SetConditionTrue(status.IPAccessListReadyType)
		r.EventRecorder.Event(project, "Normal", string(status.IPAccessListReadyType), "")
	} else {
		ctx.SetConditionFalse(status.IPAccessListReadyType)
		return reconcile.Result{Requeue: true}, nil
	}

	if result = r.ensurePrivateEndpoint(ctx, projectID, project); !result.IsOk() {
		return result.ReconcileResult(), nil
	}
	r.EventRecorder.Event(project, "Normal", string(status.PrivateEndpointReadyType), "")

	ctx.SetConditionTrue(status.ReadyType)
	return ctrl.Result{}, nil
}

// allIPAccessListsAreReady returns true if all ipAccessLists are in the ACTIVE state.
func (r *AtlasProjectReconciler) allIPAccessListsAreReady(context context.Context, ctx *workflow.Context, projectID string) (bool, workflow.Result) {
	atlasAccess, _, err := ctx.Client.ProjectIPAccessList.List(context, projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return false, workflow.Terminate(workflow.Internal, err.Error())
	}
	for _, ipAccessList := range atlasAccess.Results {
		ipStatus, err := GetIPAccessListStatus(ctx.Client, ipAccessList)
		if err != nil {
			return false, workflow.Terminate(workflow.Internal, err.Error())
		}
		if ipStatus.Status != string(IPAccessListActive) {
			r.Log.Infof("IP Access List %v is not active", ipAccessList)
			return false, workflow.InProgress(workflow.ProjectIPAccessListNotActive, fmt.Sprintf("%s IP Access List is not yet active, current state: %s", getAccessListEntry(ipAccessList), ipStatus.Status))
		}
	}
	return true, workflow.OK()
}

func (r *AtlasProjectReconciler) deleteAtlasProject(ctx context.Context, atlasClient mongodbatlas.Client, project *mdbv1.AtlasProject) (err error) {
	log := r.Log.With("atlasproject", kube.ObjectKeyFromObject(project))
	log.Infow("-> Starting AtlasProject deletion", "spec", project.Spec)

	_, err = atlasClient.Projects.Delete(ctx, project.Status.ID)
	var apiError *mongodbatlas.ErrorResponse
	if errors.As(err, &apiError) && apiError.ErrorCode == atlas.NotInGroup {
		log.Infow("Project does not exist", "projectID", project.Status.ID)
		return nil
	}

	return err
}

func (r *AtlasProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("AtlasProject", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource AtlasProject & handle delete separately
	err = c.Watch(&source.Kind{Type: &mdbv1.AtlasProject{}}, &handler.EnqueueRequestForObject{}, r.GlobalPredicates...)
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

func (r *AtlasProjectReconciler) addDeletionFinalizer(ctx context.Context, p *mdbv1.AtlasProject) error {
	p.Finalizers = append(p.GetFinalizers(), getFinalizerName())
	if err := r.Client.Update(ctx, p); err != nil {
		return fmt.Errorf("failed to add deletion finalizer for %s: %w", p.Name, err)
	}
	return nil
}

func (r *AtlasProjectReconciler) removeDeletionFinalizer(ctx context.Context, p *mdbv1.AtlasProject) error {
	p.Finalizers = removeString(p.GetFinalizers(), getFinalizerName())
	if err := r.Client.Update(ctx, p); err != nil {
		return fmt.Errorf("failed to remove deletion finalizer from %s: %w", p.Name, err)
	}
	return nil
}

func getFinalizerName() string {
	return "mongodbatlas/finalizer"
}

func isDeletionFinalizerPresent(project *mdbv1.AtlasProject) bool {
	for _, finalizer := range project.GetFinalizers() {
		if finalizer == getFinalizerName() {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return result
}
