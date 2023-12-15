package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/connectionsecret"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/compat"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/stringutil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util/toptr"
)

const FreeTier = "M0"

func (r *AtlasDeploymentReconciler) ensureAdvancedDeploymentState(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	advancedDeploymentSpec := deployment.Spec.DeploymentSpec

	advancedDeployment, resp, err := ctx.Client.AdvancedClusters.Get(context.Background(), project.Status.ID, advancedDeploymentSpec.Name)

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
		advancedDeployment, _, err = ctx.Client.AdvancedClusters.Create(context.Background(), project.Status.ID, advancedDeployment)
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

func advancedDeploymentIdle(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, atlasDeploymentAsAtlas *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	specDeployment, atlasDeployment, err := MergedAdvancedDeployment(*atlasDeploymentAsAtlas, *deployment.Spec.DeploymentSpec)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	if areEqual, _ := AdvancedDeploymentsEqual(ctx.Log, specDeployment, atlasDeployment); areEqual {
		return atlasDeploymentAsAtlas, workflow.OK()
	}

	if specDeployment.Paused != nil {
		if atlasDeployment.Paused == nil || *atlasDeployment.Paused != *specDeployment.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			specDeployment = mdbv1.AdvancedDeploymentSpec{
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

	atlasDeploymentAsAtlas, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, deployment.Spec.DeploymentSpec.Name, deploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

// MergedAdvancedDeployment will return the result of merging AtlasDeploymentSpec with Atlas Advanced Deployment
func MergedAdvancedDeployment(atlasDeploymentAsAtlas mongodbatlas.AdvancedCluster, specDeployment mdbv1.AdvancedDeploymentSpec) (mergedDeployment mdbv1.AdvancedDeploymentSpec, atlasDeployment mdbv1.AdvancedDeploymentSpec, err error) {
	if IsFreeTierAdvancedDeployment(&atlasDeploymentAsAtlas) {
		atlasDeploymentAsAtlas.DiskSizeGB = nil
	}

	atlasDeployment, err = AdvancedDeploymentFromAtlas(atlasDeploymentAsAtlas)
	if err != nil {
		return
	}

	normalizeSpecs(specDeployment.ReplicationSpecs[0].RegionConfigs)

	mergedDeployment = mdbv1.AdvancedDeploymentSpec{}

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

func AdvancedDeploymentFromAtlas(advancedDeployment mongodbatlas.AdvancedCluster) (mdbv1.AdvancedDeploymentSpec, error) {
	result := mdbv1.AdvancedDeploymentSpec{}

	convertDiskSizeField(&result, &advancedDeployment)
	if err := compat.JSONCopy(&result, advancedDeployment); err != nil {
		return result, err
	}

	return result, nil
}

func convertDiskSizeField(result *mdbv1.AdvancedDeploymentSpec, atlas *mongodbatlas.AdvancedCluster) {
	var value *int
	if atlas.DiskSizeGB != nil && *atlas.DiskSizeGB >= 1 {
		value = toptr.MakePtr(int(*atlas.DiskSizeGB))
	}
	result.DiskSizeGB = value
	atlas.DiskSizeGB = nil
}

// AdvancedDeploymentsEqual compares two Atlas Advanced Deployments
func AdvancedDeploymentsEqual(log *zap.SugaredLogger, deploymentOperator mdbv1.AdvancedDeploymentSpec, deploymentAtlas mdbv1.AdvancedDeploymentSpec) (areEqual bool, diff string) {
	deploymentAtlas = cleanupFieldsToCompare(deploymentAtlas, deploymentOperator)

	d := cmp.Diff(deploymentAtlas, deploymentOperator, cmpopts.EquateEmpty(), cmpopts.SortSlices(mdbv1.LessAD))
	if d != "" {
		log.Debugf("Deployments are different: %s", d)
	}

	return d == "", d
}

func cleanupFieldsToCompare(atlas, operator mdbv1.AdvancedDeploymentSpec) mdbv1.AdvancedDeploymentSpec {
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
		}
	}

	return atlas
}

// GetAllDeploymentNames returns all deployment names including regular and advanced deployment.
func GetAllDeploymentNames(client *mongodbatlas.Client, projectID string) ([]string, error) {
	var deploymentNames []string

	advancedDeployments, _, err := client.AdvancedClusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, d := range advancedDeployments.Results {
		deploymentNames = append(deploymentNames, d.Name)
	}

	return deploymentNames, nil
}

func (r *AtlasDeploymentReconciler) ensureConnectionSecrets(ctx *workflow.Context, project *mdbv1.AtlasProject, name string, connectionStrings *mongodbatlas.ConnectionStrings, deploymentResource *mdbv1.AtlasDeployment) workflow.Result {
	databaseUsers := mdbv1.AtlasDatabaseUserList{}
	err := r.Client.List(context.TODO(), &databaseUsers, &client.ListOptions{})
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

		scopes := dbUser.GetScopes(mdbv1.DeploymentScopeType)
		if len(scopes) != 0 && !stringutil.Contains(scopes, name) {
			continue
		}

		password, err := dbUser.ReadPassword(r.Client)
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

		secretName, err := connectionsecret.Ensure(r.Client, dbUser.Namespace, project.Spec.Name, project.ID(), name, data)
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

func dbUserBelongsToProject(dbUser *mdbv1.AtlasDatabaseUser, project *mdbv1.AtlasProject) bool {
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
