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

	"sigs.k8s.io/controller-runtime/pkg/builder"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/watch"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/kube"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	Client client.Client
	watch.ResourceWatcher
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	AtlasDomain                 string
	GlobalAPISecret             client.ObjectKey
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool
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

// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasteams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,resources=atlasteams/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasteams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.mongodb.com,namespace=default,resources=atlasteams/status,verbs=get;update;patch

func (r *AtlasProjectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.With("atlasproject", req.NamespacedName)

	project := &mdbv1.AtlasProject{}
	result := customresource.PrepareResource(r.Client, req, project, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(project) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasProject reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", project.Spec)
		if !project.GetDeletionTimestamp().IsZero() {
			err := customresource.ManageFinalizer(ctx, r.Client, project, customresource.UnsetFinalizer)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	workflowCtx := customresource.MarkReconciliationStarted(r.Client, project, log, ctx)
	log.Infow("-> Starting AtlasProject reconciliation", "spec", project.Spec)

	if project.ConnectionSecretObjectKey() != nil {
		// Note, that we are not watching the global connection secret - seems there is no point in reconciling all
		// the projects once that secret is changed
		workflowCtx.AddResourcesToWatch(watch.WatchedObject{ResourceKind: "Secret", Resource: *project.ConnectionSecretObjectKey()})
	}

	// This update will make sure the status is always updated in case of any errors or successful result
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, project)
		r.EnsureMultiplesResourcesAreWatched(req.NamespacedName, log, workflowCtx.ListResourcesToWatch()...)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, project, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("project validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if err := validate.Project(project, customresource.IsGov(r.AtlasDomain)); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		setCondition(workflowCtx, status.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(status.ValidationSucceeded)

	if !customresource.IsResourceSupportedInDomain(project, r.AtlasDomain) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasProject is not supported by Atlas for government").
			WithoutRetry()
		setCondition(workflowCtx, status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}

	connection, err := atlas.ReadConnection(log, r.Client, r.GlobalAPISecret, project.ConnectionSecretObjectKey())
	if err != nil {
		result = workflow.Terminate(workflow.AtlasCredentialsNotProvided, err.Error())
		setCondition(workflowCtx, status.ProjectReadyType, result)
		if errRm := customresource.ManageFinalizer(ctx, r.Client, project, customresource.UnsetFinalizer); errRm != nil {
			result = workflow.Terminate(workflow.Internal, errRm.Error())
			return result.ReconcileResult(), nil
		}
		return result.ReconcileResult(), nil
	}
	workflowCtx.Connection = connection

	atlasClient, err := atlas.Client(r.AtlasDomain, connection, log)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		setCondition(workflowCtx, status.DeploymentReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.Client = atlasClient

	owner, err := customresource.IsOwner(project, r.ObjectDeletionProtection, customresource.IsResourceManagedByOperator, managedByAtlas(workflowCtx))
	if err != nil {
		result = workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.ProjectReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	if !owner {
		result = workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Project due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.ProjectReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	projectID, result := r.ensureProjectExists(workflowCtx, project)
	if !result.IsOk() {
		setCondition(workflowCtx, status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}

	workflowCtx.EnsureStatusOption(status.AtlasProjectIDOption(projectID))

	if result = r.ensureDeletionFinalizer(workflowCtx, atlasClient, project); !result.IsOk() {
		setCondition(workflowCtx, status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}

	if project.ID() == "" {
		err = customresource.ApplyLastConfigApplied(ctx, project, r.Client)
		if err != nil {
			result = workflow.Terminate(workflow.Internal, err.Error())
			workflowCtx.SetConditionFromResult(status.ProjectReadyType, result)
			log.Error(result.GetMessage())

			return result.ReconcileResult(), nil
		}

		return result.WithRetry(workflow.DefaultRetry).ReconcileResult(), nil
	}

	var authModes authmode.AuthModes
	if authModes, result = r.ensureX509(workflowCtx, projectID, project); !result.IsOk() {
		setCondition(workflowCtx, status.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	authModes.AddAuthMode(authmode.Scram) // add the default auth method
	workflowCtx.EnsureStatusOption(status.AtlasProjectAuthModesOption(authModes))

	// Updating the status with "projectReady = true" and "IPAccessListReady = false" (not as separate updates!)
	workflowCtx.SetConditionTrue(status.ProjectReadyType)
	r.EventRecorder.Event(project, "Normal", string(status.ProjectReadyType), "")

	results := r.ensureProjectResources(workflowCtx, project)
	for i := range results {
		if !results[i].IsOk() {
			logIfWarning(workflowCtx, result)
			return results[i].ReconcileResult(), nil
		}
	}

	err = customresource.ApplyLastConfigApplied(ctx, project, r.Client)
	if err != nil {
		result = workflow.Terminate(workflow.Internal, err.Error())
		workflowCtx.SetConditionFromResult(status.ProjectReadyType, result)
		log.Error(result.GetMessage())

		return result.ReconcileResult(), nil
	}

	workflowCtx.SetConditionTrue(status.ReadyType)
	return workflow.OK().ReconcileResult(), nil
}

func (r *AtlasProjectReconciler) ensureDeletionFinalizer(workflowCtx *workflow.Context, atlasClient mongodbatlas.Client, project *mdbv1.AtlasProject) (result workflow.Result) {
	log := workflowCtx.Log

	if project.GetDeletionTimestamp().IsZero() {
		if !customresource.HaveFinalizer(project, customresource.FinalizerLabel) {
			log.Debugw("Add deletion finalizer", "name", customresource.FinalizerLabel)
			if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, project, customresource.SetFinalizer); err != nil {
				return workflow.Terminate(workflow.AtlasFinalizerNotSet, err.Error())
			}
		}
	}

	if !project.GetDeletionTimestamp().IsZero() {
		if customresource.HaveFinalizer(project, customresource.FinalizerLabel) {
			if customresource.IsResourceProtected(project, r.ObjectDeletionProtection) {
				log.Info("Not removing Project from Atlas as per configuration")
				result = workflow.OK()
			} else {
				if result = DeleteAllPrivateEndpoints(workflowCtx, project.ID()); !result.IsOk() {
					setCondition(workflowCtx, status.PrivateEndpointReadyType, result)
					return result
				}
				if result = DeleteAllNetworkPeers(workflowCtx.Context, project.ID(), workflowCtx.Client.Peers, workflowCtx.Log); !result.IsOk() {
					setCondition(workflowCtx, status.NetworkPeerReadyType, result)
					return result
				}

				if err := r.deleteAtlasProject(workflowCtx.Context, atlasClient, project); err != nil {
					result = workflow.Terminate(workflow.Internal, err.Error())
					setCondition(workflowCtx, status.DeploymentReadyType, result)
					return result
				}
			}

			if err := customresource.ManageFinalizer(workflowCtx.Context, r.Client, project, customresource.UnsetFinalizer); err != nil {
				return workflow.Terminate(workflow.AtlasFinalizerNotRemoved, err.Error())
			}
		}
		return result
	}

	return workflow.OK()
}

// ensureProjectResources ensures IP Access List, Private Endpoints, Integrations, Maintenance Window and Encryption at Rest
func (r *AtlasProjectReconciler) ensureProjectResources(workflowCtx *workflow.Context, project *mdbv1.AtlasProject) (results []workflow.Result) {
	for k, v := range project.Annotations {
		workflowCtx.Log.Debugf(k)
		workflowCtx.Log.Debugf(v)
	}

	var result workflow.Result
	if result = ensureIPAccessList(workflowCtx, atlas.CustomIPAccessListStatus(&workflowCtx.Client), project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.IPAccessListReadyType), "")
	}
	results = append(results, result)

	if result = ensurePrivateEndpoint(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.PrivateEndpointReadyType), "")
	}
	results = append(results, result)

	if result = ensureProviderAccessStatus(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.CloudProviderAccessReadyType), "")
	}
	results = append(results, result)

	if result = ensureNetworkPeers(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.NetworkPeerReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureAlertConfigurations(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.AlertConfigurationReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureIntegration(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.IntegrationReadyType), "")
	}
	results = append(results, result)

	if result = ensureMaintenanceWindow(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.MaintenanceWindowReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureEncryptionAtRest(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.EncryptionAtRestReadyType), "")
	}
	results = append(results, result)

	if result = ensureAuditing(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.AuditingReadyType), "")
	}
	results = append(results, result)

	if result = ensureProjectSettings(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.ProjectSettingsReadyType), "")
	}
	results = append(results, result)

	if result = ensureCustomRoles(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.ProjectCustomRolesReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureAssignedTeams(workflowCtx, project, r.SubObjectDeletionProtection); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(status.ProjectTeamsReadyType), "")
	}
	results = append(results, result)

	return results
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
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasProject").
		For(&mdbv1.AtlasProject{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(&source.Kind{Type: &corev1.Secret{}}, watch.NewSecretHandler(r.WatchedResources)).
		Watches(&source.Kind{Type: &mdbv1.AtlasTeam{}}, watch.NewAtlasTeamHandler(r.WatchedResources)).
		Complete(r)
}

// setCondition sets the condition from the result and logs the warnings
func setCondition(ctx *workflow.Context, condition status.ConditionType, result workflow.Result) {
	ctx.SetConditionFromResult(condition, result)
	logIfWarning(ctx, result)
}

func logIfWarning(ctx *workflow.Context, result workflow.Result) {
	if result.IsWarning() {
		ctx.Log.Warnw(result.GetMessage())
	}
}

func managedByAtlas(workflowCtx *workflow.Context) customresource.AtlasChecker {
	return func(resource mdbv1.AtlasCustomResource) (bool, error) {
		project, ok := resource.(*mdbv1.AtlasProject)
		if !ok {
			return false, errors.New("failed to match resource type as AtlasProject")
		}

		if project.ID() == "" {
			return false, nil
		}

		atlasProject, _, err := workflowCtx.Client.Projects.GetOneProject(workflowCtx.Context, project.ID())
		if err != nil {
			var apiError *mongodbatlas.ErrorResponse
			if errors.As(err, &apiError) && (apiError.ErrorCode == atlas.NotInGroup || apiError.ErrorCode == atlas.ResourceNotFound) {
				return false, nil
			}

			return false, err
		}

		if project.Spec.Name == atlasProject.Name {
			return false, err
		}

		return true, nil
	}
}
