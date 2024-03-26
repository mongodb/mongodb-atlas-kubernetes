package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/stringutil"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const FreeTier = "M0"

func (r *AtlasDeploymentReconciler) ensureAdvancedDeploymentState(ctx *workflow.Context, project *akov2.AtlasProject, deployment *akov2.AtlasDeployment) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	advancedDeploymentSpec := deployment.Spec.DeploymentSpec

	advancedDeployment, resp, err := ctx.Client.AdvancedClusters.Get(ctx.Context, project.Status.ID, advancedDeploymentSpec.Name)

	if err != nil {
		if resp == nil {
			return advancedDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return advancedDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}

		advancedDeployment, err = advancedDeploymentSpec.ToAtlas()
		if err != nil {
			return advancedDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Advanced Deployment %s doesn't exist in Atlas - creating", advancedDeploymentSpec.Name)
		advancedDeployment, _, err = ctx.Client.AdvancedClusters.Create(ctx.Context, project.Status.ID, advancedDeployment)
		if err != nil {
			return advancedDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}
	}

	result := EnsureCustomZoneMapping(ctx, project.ID(), deployment.Spec.DeploymentSpec.CustomZoneMapping, advancedDeployment.Name)
	if !result.IsOk() {
		return advancedDeployment, result
	}

	result = EnsureManagedNamespaces(ctx, project.ID(), deployment.Spec.DeploymentSpec.ClusterType, deployment.Spec.DeploymentSpec.ManagedNamespaces, advancedDeployment.Name)
	if !result.IsOk() {
		return advancedDeployment, result
	}

	switch advancedDeployment.StateName {
	case "IDLE":
		return advancedDeploymentIdle(ctx, project, deployment, advancedDeployment)

	case "CREATING":
		return advancedDeployment, workflow.InProgress(workflow.DeploymentCreating, "deployment is provisioning")

	case "UPDATING", "REPAIRING":
		return advancedDeployment, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")

	// TODO: add "DELETING", "DELETED", handle 404 on delete

	default:
		return advancedDeployment, workflow.Terminate(workflow.Internal, fmt.Sprintf("unknown deployment state %q", advancedDeployment.StateName))
	}
}

func handleSearchNodes(ctx *workflow.Context, deployment *akov2.AtlasDeployment, projectID string) workflow.Result {
	// Get the current state from Atlas
	// If no nodes configured in atlas, but some in the operator - create them. If not - do nothing
	// If nodes already configured in atlas and in the operator - update them. If not - delete them
	ctx.Log.Debug("starting search node processing")
	defer ctx.Log.Debug("finished search node processing")

	nodesInAkoEmpty := deployment.Spec.DeploymentSpec.SearchNodes == nil || len(deployment.Spec.DeploymentSpec.SearchNodes) == 0

	nodesInAtlasEmpty := false
	currentNodesInAtlas, httpResp, err := ctx.SdkClient.AtlasSearchApi.GetAtlasSearchDeployment(ctx.Context, projectID, deployment.GetDeploymentName()).Execute()
	if err != nil {
		if httpResp == nil {
			ctx.Log.Debugf("unable to get current search nodes status: %v", err)
			return workflow.Terminate(workflow.SearchNodesReady, err.Error())
		}

		switch httpResp.StatusCode {
		// TODO: Should be only NotFound: CLOUDP-239015
		case http.StatusBadRequest, http.StatusNotFound:
			ctx.Log.Debug("no search nodes in atlas found")
			nodesInAtlasEmpty = true
		default:
			return workflow.Terminate(workflow.SearchNodesReady, err.Error())
		}
	}

	if !nodesInAtlasEmpty && currentNodesInAtlas.GetStateName() != "IDLE" {
		msg := fmt.Sprintf("search nodes are not ready: %v", currentNodesInAtlas.GetStateName())
		ctx.Log.Debug(msg)
		return workflow.InProgress(workflow.SearchNodesReady, msg)
	}

	switch {
	case !nodesInAkoEmpty && nodesInAtlasEmpty:
		ctx.Log.Debugf("creating search nodes %v", deployment.Spec.DeploymentSpec.SearchNodes)
		_, _, err = ctx.SdkClient.AtlasSearchApi.CreateAtlasSearchDeployment(ctx.Context, projectID, deployment.GetDeploymentName(), &admin.ApiSearchDeploymentRequest{
			Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas(),
		}).Execute()
		if err != nil {
			ctx.Log.Debugf("unable to create search nodes: %v", err)
			return workflow.Terminate(workflow.SearchNodesReady, err.Error())
		}
	case !nodesInAkoEmpty && !nodesInAtlasEmpty:
		ctx.Log.Debugf("updating search nodes %v", deployment.Spec.DeploymentSpec.SearchNodes)
		currentAkoNodesAsAtlas := deployment.Spec.DeploymentSpec.SearchNodesToAtlas()
		if reflect.DeepEqual(currentAkoNodesAsAtlas, currentNodesInAtlas) {
			ctx.Log.Debug("search nodes in AKO and Atlas are equal")
			return workflow.OK()
		}

		_, _, err = ctx.SdkClient.AtlasSearchApi.UpdateAtlasSearchDeployment(ctx.Context, projectID, deployment.GetDeploymentName(), &admin.ApiSearchDeploymentRequest{
			Specs: deployment.Spec.DeploymentSpec.SearchNodesToAtlas(),
		}).Execute()
		if err != nil {
			ctx.Log.Debugf("unable to update search nodes: %v", err)
			return workflow.Terminate(workflow.SearchNodesReady, err.Error())
		}
	case nodesInAkoEmpty && !nodesInAtlasEmpty:
		ctx.Log.Debug("deleting search nodes")
		_, err = ctx.SdkClient.AtlasSearchApi.DeleteAtlasSearchDeployment(ctx.Context, projectID, deployment.GetDeploymentName()).Execute()
		if err != nil {
			ctx.Log.Debugf("unable to delete search nodes: %v", err)
			return workflow.Terminate(workflow.SearchNodesReady, err.Error())
		}
	}

	return workflow.OK()
}

func advancedDeploymentIdle(ctx *workflow.Context, project *akov2.AtlasProject, deployment *akov2.AtlasDeployment, atlasDeploymentAsAtlas *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	specDeployment, atlasDeployment, err := MergedAdvancedDeployment(*atlasDeploymentAsAtlas, *deployment.Spec.DeploymentSpec)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	searchNodeResult := handleSearchNodes(ctx, deployment, project.ID())
	if !searchNodeResult.IsOk() {
		return atlasDeploymentAsAtlas, searchNodeResult
	}

	if areEqual, _ := AdvancedDeploymentsEqual(ctx.Log, &specDeployment, &atlasDeployment); areEqual {
		return atlasDeploymentAsAtlas, workflow.OK()
	}

	if specDeployment.Paused != nil {
		if atlasDeployment.Paused == nil || *atlasDeployment.Paused != *specDeployment.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			specDeployment = akov2.AdvancedDeploymentSpec{
				Paused: deployment.Spec.DeploymentSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			specDeployment.Paused = nil
		}
	}

	syncRegionConfiguration(&specDeployment, atlasDeploymentAsAtlas)

	deploymentAsAtlas, err := specDeployment.ToAtlas()
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	// TODO: Potential bug with disabling autoscaling if it was previously enabled

	atlasDeploymentAsAtlas, _, err = ctx.Client.AdvancedClusters.Update(ctx.Context, project.Status.ID, deployment.Spec.DeploymentSpec.Name, deploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

// MergedAdvancedDeployment will return the result of merging AtlasDeploymentSpec with Atlas Advanced Deployment
func MergedAdvancedDeployment(atlasDeploymentAsAtlas mongodbatlas.AdvancedCluster, specDeployment akov2.AdvancedDeploymentSpec) (mergedDeployment akov2.AdvancedDeploymentSpec, atlasDeployment akov2.AdvancedDeploymentSpec, err error) {
	if IsFreeTierAdvancedDeployment(&atlasDeploymentAsAtlas) {
		atlasDeploymentAsAtlas.DiskSizeGB = nil
	}

	atlasDeployment, err = AdvancedDeploymentFromAtlas(atlasDeploymentAsAtlas)
	if err != nil {
		return
	}

	normalizeSpecs(specDeployment.ReplicationSpecs[0].RegionConfigs)

	mergedDeployment = akov2.AdvancedDeploymentSpec{}

	if err = compat.JSONCopy(&mergedDeployment, atlasDeployment); err != nil {
		return
	}

	if err = compat.JSONCopy(&mergedDeployment, specDeployment); err != nil {
		return
	}

	for i, replicationSpec := range atlasDeployment.ReplicationSpecs {
		for k, v := range replicationSpec.RegionConfigs {
			// the response does not return backing provider names in some situations.
			// if this is the case, we want to strip these fields so they do not cause a bad comparison.
			if v.BackingProviderName == "" && k < len(mergedDeployment.ReplicationSpecs[i].RegionConfigs) {
				mergedDeployment.ReplicationSpecs[i].RegionConfigs[k].BackingProviderName = ""
			}
		}
	}

	atlasDeployment.MongoDBVersion = ""
	mergedDeployment.MongoDBVersion = ""

	return
}

func IsFreeTierAdvancedDeployment(deployment *mongodbatlas.AdvancedCluster) bool {
	if deployment != nil && deployment.ReplicationSpecs != nil {
		for _, replicationSpec := range deployment.ReplicationSpecs {
			if replicationSpec.RegionConfigs != nil {
				for _, regionConfig := range replicationSpec.RegionConfigs {
					if regionConfig != nil &&
						regionConfig.ElectableSpecs != nil &&
						regionConfig.ElectableSpecs.InstanceSize == FreeTier {
						return true
					}
				}
			}
		}
	}
	return false
}

func AdvancedDeploymentFromAtlas(advancedDeployment mongodbatlas.AdvancedCluster) (akov2.AdvancedDeploymentSpec, error) {
	result := akov2.AdvancedDeploymentSpec{}

	convertDiskSizeField(&result, &advancedDeployment)
	if err := compat.JSONCopy(&result, advancedDeployment); err != nil {
		return result, err
	}

	return result, nil
}

func convertDiskSizeField(result *akov2.AdvancedDeploymentSpec, atlas *mongodbatlas.AdvancedCluster) {
	var value *int
	if atlas.DiskSizeGB != nil && *atlas.DiskSizeGB >= 1 {
		value = pointer.MakePtr(int(*atlas.DiskSizeGB))
	}
	result.DiskSizeGB = value
	atlas.DiskSizeGB = nil
}

// AdvancedDeploymentsEqual compares two Atlas Advanced Deployments
func AdvancedDeploymentsEqual(log *zap.SugaredLogger, deploymentOperator *akov2.AdvancedDeploymentSpec, deploymentAtlas *akov2.AdvancedDeploymentSpec) (areEqual bool, diff string) {
	expected := deploymentOperator.DeepCopy()
	actualCleaned := cleanupFieldsToCompare(deploymentAtlas.DeepCopy(), expected)

	// Ignore Atlas Search and Search Nodes
	expected.SearchNodes = nil
	// Ignore differences on auto-scaled region configs
	for _, rs := range expected.ReplicationSpecs {
		for _, rc := range rs.RegionConfigs {
			if isComputeAutoScalingEnabled(rc.AutoScaling) {
				ignoreInstanceSize(rc)
			}
		}
	}
	d := cmp.Diff(actualCleaned, expected, cmpopts.EquateEmpty(), cmpopts.SortSlices(akov2.LessAD))
	if d != "" {
		log.Debugf("Deployments are different: %s", d)
	}

	return d == "", d
}

func cleanupFieldsToCompare(atlas, operator *akov2.AdvancedDeploymentSpec) *akov2.AdvancedDeploymentSpec {
	if atlas.ReplicationSpecs == nil {
		return atlas
	}

	for specIdx, replicationSpec := range atlas.ReplicationSpecs {
		if replicationSpec.RegionConfigs == nil {
			continue
		}

		for configIdx, regionConfig := range replicationSpec.RegionConfigs {
			if regionConfig.AnalyticsSpecs == nil || regionConfig.AnalyticsSpecs.NodeCount == nil || *regionConfig.AnalyticsSpecs.NodeCount == 0 {
				if configIdx < len(operator.ReplicationSpecs[specIdx].RegionConfigs) {
					regionConfig.AnalyticsSpecs = operator.ReplicationSpecs[specIdx].RegionConfigs[configIdx].AnalyticsSpecs
				}
			}

			if regionConfig.ElectableSpecs == nil || regionConfig.ElectableSpecs.NodeCount == nil || *regionConfig.ElectableSpecs.NodeCount == 0 {
				if configIdx < len(operator.ReplicationSpecs[specIdx].RegionConfigs) {
					regionConfig.ElectableSpecs = operator.ReplicationSpecs[specIdx].RegionConfigs[configIdx].ElectableSpecs
				}
			}

			if regionConfig.ReadOnlySpecs == nil || regionConfig.ReadOnlySpecs.NodeCount == nil || *regionConfig.ReadOnlySpecs.NodeCount == 0 {
				if configIdx < len(operator.ReplicationSpecs[specIdx].RegionConfigs) {
					regionConfig.ReadOnlySpecs = operator.ReplicationSpecs[specIdx].RegionConfigs[configIdx].ReadOnlySpecs
				}
			}

			if isComputeAutoScalingEnabled(regionConfig.AutoScaling) {
				ignoreInstanceSize(regionConfig)
			}
		}
	}

	return atlas
}

// GetAllDeploymentNames returns all deployment names including regular and advanced deployment.
func GetAllDeploymentNames(ctx context.Context, client *mongodbatlas.Client, projectID string) ([]string, error) {
	var deploymentNames []string

	advancedDeployments, _, err := client.AdvancedClusters.List(ctx, projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, d := range advancedDeployments.Results {
		deploymentNames = append(deploymentNames, d.Name)
	}

	return deploymentNames, nil
}

func (r *AtlasDeploymentReconciler) ensureConnectionSecrets(ctx *workflow.Context, project *akov2.AtlasProject, name string, connectionStrings *mongodbatlas.ConnectionStrings, deploymentResource *akov2.AtlasDeployment) workflow.Result {
	databaseUsers := akov2.AtlasDatabaseUserList{}
	err := r.Client.List(ctx.Context, &databaseUsers, &client.ListOptions{})
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	secrets := make([]string, 0)
	for i := range databaseUsers.Items {
		dbUser := databaseUsers.Items[i]

		if !dbUserBelongsToProject(&dbUser, project) {
			continue
		}

		found := false
		for _, c := range dbUser.Status.Conditions {
			if c.Type == status.ReadyType && c.Status == v1.ConditionTrue {
				found = true
				break
			}
		}

		if !found {
			ctx.Log.Debugw("AtlasDatabaseUser not ready - not creating connection secret", "user.name", dbUser.Name)
			continue
		}

		scopes := dbUser.GetScopes(akov2.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, name) {
			continue
		}

		password, err := dbUser.ReadPassword(ctx.Context, r.Client)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}

		data := connectionsecret.ConnectionData{
			DBUserName: dbUser.Spec.Username,
			Password:   password,
			ConnURL:    connectionStrings.Standard,
			SrvConnURL: connectionStrings.StandardSrv,
		}
		connectionsecret.FillPrivateConnStrings(connectionStrings, &data)

		ctx.Log.Debugw("Creating a connection Secret", "data", data)

		secretName, err := connectionsecret.Ensure(ctx.Context, r.Client, dbUser.Namespace, project.Spec.Name, project.ID(), name, data)
		if err != nil {
			return workflow.Terminate(workflow.DeploymentConnectionSecretsNotCreated, err.Error())
		}
		secrets = append(secrets, secretName)
	}

	if len(secrets) > 0 {
		r.EventRecorder.Eventf(deploymentResource, "Normal", "ConnectionSecretsEnsured", "Connection Secrets were created/updated: %s", strings.Join(secrets, ", "))
	}

	return workflow.OK()
}

func dbUserBelongsToProject(dbUser *akov2.AtlasDatabaseUser, project *akov2.AtlasProject) bool {
	if dbUser.Spec.Project.Name != project.Name {
		return false
	}

	if dbUser.Spec.Project.Namespace == "" && dbUser.Namespace != project.Namespace {
		return false
	}

	if dbUser.Spec.Project.Namespace != "" && dbUser.Spec.Project.Namespace != project.Namespace {
		return false
	}

	return true
}
