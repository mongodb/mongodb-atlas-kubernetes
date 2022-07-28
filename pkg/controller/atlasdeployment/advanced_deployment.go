package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
)

func (r *AtlasDeploymentReconciler) ensureAdvancedDeploymentState(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	advancedDeploymentSpec := deployment.Spec.AdvancedDeploymentSpec

	advancedDeployment, resp, err := ctx.Client.AdvancedClusters.Get(context.Background(), project.Status.ID, advancedDeploymentSpec.Name)

	if err != nil {
		if resp == nil {
			return advancedDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return advancedDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}

		advancedDeployment, err = advancedDeploymentSpec.AdvancedDeployment()
		if err != nil {
			return advancedDeployment, workflow.Terminate(workflow.Internal, err.Error())
		}

		ctx.Log.Infof("Advanced Deployment %s doesn't exist in Atlas - creating", advancedDeploymentSpec.Name)
		advancedDeployment, _, err = ctx.Client.AdvancedClusters.Create(context.Background(), project.Status.ID, advancedDeployment)
		if err != nil {
			return advancedDeployment, workflow.Terminate(workflow.DeploymentNotCreatedInAtlas, err.Error())
		}
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
	specDeployment := *deployment.Spec.AdvancedDeploymentSpec
	atlasDeployment, err := AdvancedDeploymentFromAtlas(*atlasDeploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	if areEqual := AdvancedDeploymentsEqual(ctx.Log, specDeployment, atlasDeployment); areEqual {
		return atlasDeploymentAsAtlas, workflow.OK()
	}

	if specDeployment.Paused != nil {
		if atlasDeployment.Paused == nil || *atlasDeployment.Paused != *specDeployment.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			specDeployment = mdbv1.AdvancedDeploymentSpec{
				Paused: deployment.Spec.AdvancedDeploymentSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			specDeployment.Paused = nil
		}
	}

	deploymentAsAtlas, err := specDeployment.AdvancedDeployment()
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	atlasDeploymentAsAtlas, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, deployment.Spec.AdvancedDeploymentSpec.Name, deploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

// func cleanupAdvancedDeployment(deployment mongodbatlas.AdvancedCluster) mongodbatlas.AdvancedCluster {
// 	deployment.ID = ""
// 	deployment.GroupID = ""
// 	deployment.MongoDBVersion = ""
// 	deployment.CreateDate = ""
// 	deployment.StateName = ""
// 	deployment.ConnectionStrings = nil
// 	for i := range deployment.ReplicationSpecs {
// 		deployment.ReplicationSpecs[i].ID = ""
// 	}
// 	return deployment
// }

// MergedAdvancedDeployment will return the result of merging AtlasDeploymentSpec with Atlas Advanced Deployment
func MergedAdvancedDeployment(advancedDeployment mongodbatlas.AdvancedCluster, spec mdbv1.AtlasDeploymentSpec) (mongodbatlas.AdvancedCluster, error) {
	result := mongodbatlas.AdvancedCluster{}
	if err := compat.JSONCopy(&result, advancedDeployment); err != nil {
		return result, err
	}

	if err := compat.JSONCopy(&result, spec.AdvancedDeploymentSpec); err != nil {
		return result, err
	}

	for i, replicationSpec := range advancedDeployment.ReplicationSpecs {
		for k, v := range replicationSpec.RegionConfigs {
			// the response does not return backing provider names in some situations.
			// if this is the case, we want to strip these fields so they do not cause a bad comparison.
			if v.BackingProviderName == "" {
				result.ReplicationSpecs[i].RegionConfigs[k].BackingProviderName = ""
			}
		}
	}
	return result, nil
}

func AdvancedDeploymentFromAtlas(advancedDeployment mongodbatlas.AdvancedCluster) (mdbv1.AdvancedDeploymentSpec, error) {
	result := mdbv1.AdvancedDeploymentSpec{}
	if err := compat.JSONCopy(&result, advancedDeployment); err != nil {
		return result, err
	}

	return result, nil
}

// AdvancedDeploymentsEqual compares two Atlas Advanced Deployments
func AdvancedDeploymentsEqual(log *zap.SugaredLogger, deploymentAtlas mdbv1.AdvancedDeploymentSpec, deploymentOperator mdbv1.AdvancedDeploymentSpec) bool {
	d := cmp.Diff(deploymentAtlas, deploymentOperator, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Deployments are different: %s", d)
	}

	return d == ""
}

// GetAllDeploymentNames returns all deployment names including regular and advanced deployment.
func GetAllDeploymentNames(client mongodbatlas.Client, projectID string) ([]string, error) {
	var deploymentNames []string
	deployment, _, err := client.Clusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	advancedDeployments, _, err := client.AdvancedClusters.List(context.Background(), projectID, &mongodbatlas.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, c := range deployment {
		deploymentNames = append(deploymentNames, c.Name)
	}

	for _, d := range advancedDeployments.Results {
		// based on configuration settings, some advanced deployment also show up in the regular deployments API.
		// For these deployments, we don't want to duplicate the secret so we skip them.
		found := false
		for _, regularDeployment := range deployment {
			if regularDeployment.Name == d.Name {
				found = true
				break
			}
		}

		// we only include deployment names which have not been handled by the regular deployment API.
		if !found {
			deploymentNames = append(deploymentNames, d.Name)
		}
	}
	return deploymentNames, nil
}
