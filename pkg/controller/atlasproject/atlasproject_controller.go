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
	"fmt"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/project"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/validate"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/indexer"
)

// AtlasProjectReconciler reconciles a AtlasProject object
type AtlasProjectReconciler struct {
	Client                      client.Client
	Log                         *zap.SugaredLogger
	Scheme                      *runtime.Scheme
	GlobalPredicates            []predicate.Predicate
	EventRecorder               record.EventRecorder
	AtlasProvider               atlas.Provider
	ObjectDeletionProtection    bool
	SubObjectDeletionProtection bool

	projectService project.ProjectService
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

	atlasProject := &akov2.AtlasProject{}
	result := customresource.PrepareResource(ctx, r.Client, req, atlasProject, log)
	if !result.IsOk() {
		return result.ReconcileResult(), nil
	}

	if customresource.ReconciliationShouldBeSkipped(atlasProject) {
		log.Infow(fmt.Sprintf("-> Skipping AtlasProject reconciliation as annotation %s=%s", customresource.ReconciliationPolicyAnnotation, customresource.ReconciliationPolicySkip), "spec", atlasProject.Spec)
		if !atlasProject.GetDeletionTimestamp().IsZero() {
			err := customresource.ManageFinalizer(ctx, r.Client, atlasProject, customresource.UnsetFinalizer)
			if err != nil {
				result = workflow.Terminate(workflow.Internal, err.Error())
				log.Errorw("Failed to remove finalizer", "error", err)
				return result.ReconcileResult(), nil
			}
		}
		return workflow.OK().ReconcileResult(), nil
	}

	conditions := akov2.InitCondition(atlasProject, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(log, conditions, ctx)
	log.Infow("-> Starting AtlasProject reconciliation", "spec", atlasProject.Spec)

	// This update will make sure the status is always updated in case of any errors or successful result
	defer func() {
		statushandler.Update(workflowCtx, r.Client, r.EventRecorder, atlasProject)
	}()

	resourceVersionIsValid := customresource.ValidateResourceVersion(workflowCtx, atlasProject, r.Log)
	if !resourceVersionIsValid.IsOk() {
		r.Log.Debugf("project validation result: %v", resourceVersionIsValid)
		return resourceVersionIsValid.ReconcileResult(), nil
	}

	if err := validate.Project(atlasProject, r.AtlasProvider.IsCloudGov()); err != nil {
		result := workflow.Terminate(workflow.Internal, err.Error())
		setCondition(workflowCtx, api.ValidationSucceeded, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SetConditionTrue(api.ValidationSucceeded)

	if !r.AtlasProvider.IsResourceSupported(atlasProject) {
		result := workflow.Terminate(workflow.AtlasGovUnsupported, "the AtlasProject is not supported by Atlas for government").
			WithoutRetry()
		setCondition(workflowCtx, api.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}

	atlasSdkClient, orgID, err := r.AtlasProvider.SdkClient(workflowCtx.Context, atlasProject.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		setCondition(workflowCtx, api.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.SdkClient = atlasSdkClient
	r.projectService = project.NewProjectAPIService(atlasSdkClient.ProjectsApi)

	atlasClient, _, err := r.AtlasProvider.Client(workflowCtx.Context, atlasProject.ConnectionSecretObjectKey(), log)
	if err != nil {
		result := workflow.Terminate(workflow.AtlasAPIAccessNotConfigured, err.Error())
		setCondition(workflowCtx, api.ProjectReadyType, result)
		return result.ReconcileResult(), nil
	}
	workflowCtx.OrgID = orgID
	workflowCtx.Client = atlasClient

	return r.handleProject(workflowCtx, orgID, atlasProject)
}

// ensureProjectResources ensures IP Access List, Private Endpoints, Integrations, Maintenance Window and Encryption at Rest
func (r *AtlasProjectReconciler) ensureProjectResources(workflowCtx *workflow.Context, project *akov2.AtlasProject) (results []workflow.Result) {
	for k, v := range project.Annotations {
		workflowCtx.Log.Debugf(k)
		workflowCtx.Log.Debugf(v)
	}

	var result workflow.Result
	if result = handleIPAccessList(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.IPAccessListReadyType), "")
	}
	results = append(results, result)

	if result = ensurePrivateEndpoint(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.PrivateEndpointReadyType), "")
	}
	results = append(results, result)

	if result = ensureCloudProviderIntegration(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.CloudProviderIntegrationReadyType), "")
	}
	results = append(results, result)

	if result = ensureNetworkPeers(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.NetworkPeerReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureAlertConfigurations(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.AlertConfigurationReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureIntegration(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.IntegrationReadyType), "")
	}
	results = append(results, result)

	if result = ensureMaintenanceWindow(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.MaintenanceWindowReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureEncryptionAtRest(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.EncryptionAtRestReadyType), "")
	}
	results = append(results, result)

	if result = handleAudit(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.AuditingReadyType), "")
	}
	results = append(results, result)

	if result = ensureProjectSettings(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.ProjectSettingsReadyType), "")
	}
	results = append(results, result)

	if result = ensureCustomRoles(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.ProjectCustomRolesReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureAssignedTeams(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.ProjectTeamsReadyType), "")
	}
	results = append(results, result)

	if result = r.ensureBackupCompliance(workflowCtx, project); result.IsOk() {
		r.EventRecorder.Event(project, "Normal", string(api.BackupComplianceReadyType), "")
	}
	results = append(results, result)

	return results
}

func (r *AtlasProjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("AtlasProject").
		For(&akov2.AtlasProject{}, builder.WithPredicates(r.GlobalPredicates...)).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(newProjectsMapFunc[corev1.Secret](indexer.AtlasProjectBySecretsIndex, r.Client, r.Log)),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&akov2.AtlasTeam{},
			handler.EnqueueRequestsFromMapFunc(newProjectsMapFunc[akov2.AtlasTeam](indexer.AtlasProjectByTeamIndex, r.Client, r.Log)),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Watches(
			&akov2.AtlasBackupCompliancePolicy{},
			handler.EnqueueRequestsFromMapFunc(newProjectsMapFunc[akov2.AtlasBackupCompliancePolicy](indexer.AtlasProjectByBackupCompliancePolicyIndex, r.Client, r.Log)),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Complete(r)
}

func NewAtlasProjectReconciler(
	mgr manager.Manager,
	predicates []predicate.Predicate,
	atlasProvider atlas.Provider,
	deletionProtection bool,
	logger *zap.Logger,
) *AtlasProjectReconciler {
	return &AtlasProjectReconciler{
		Scheme:                   mgr.GetScheme(),
		Client:                   mgr.GetClient(),
		EventRecorder:            mgr.GetEventRecorderFor("AtlasProject"),
		GlobalPredicates:         predicates,
		Log:                      logger.Named("controllers").Named("AtlasProject").Sugar(),
		AtlasProvider:            atlasProvider,
		ObjectDeletionProtection: deletionProtection,
	}
}

func newProjectsMapFunc[T any](indexName string, kubeClient client.Client, logger *zap.SugaredLogger) handler.MapFunc {
	return func(ctx context.Context, obj client.Object) []reconcile.Request {
		_, ok := any(obj).(*T)
		if !ok {
			var watchedObject T
			logger.Warnf("watching %T but got %T", &watchedObject, obj)
			return nil
		}

		projects := &akov2.AtlasProjectList{}
		listOpts := &client.ListOptions{
			FieldSelector: fields.OneTermEqualSelector(
				indexName,
				client.ObjectKeyFromObject(obj).String(),
			),
		}
		err := kubeClient.List(ctx, projects, listOpts)
		if err != nil {
			logger.Errorf("failed to list Atlas projects: %e", err)
			return []reconcile.Request{}
		}

		requests := make([]reconcile.Request, 0, len(projects.Items))
		for i := range projects.Items {
			item := projects.Items[i]
			requests = append(
				requests,
				reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      item.Name,
						Namespace: item.Namespace,
					},
				},
			)
		}
		return requests
	}
}

// setCondition sets the condition from the result and logs the warnings
func setCondition(ctx *workflow.Context, condition api.ConditionType, result workflow.Result) {
	ctx.SetConditionFromResult(condition, result)
	logIfWarning(ctx, result)
}

func logIfWarning(ctx *workflow.Context, result workflow.Result) {
	if result.IsWarning() {
		ctx.Log.Warnw(result.GetMessage())
	}
}
