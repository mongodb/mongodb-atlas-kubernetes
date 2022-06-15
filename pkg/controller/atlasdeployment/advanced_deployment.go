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

func advancedDeploymentIdle(ctx *workflow.Context, project *mdbv1.AtlasProject, deployment *mdbv1.AtlasDeployment, advancedDeployment *mongodbatlas.AdvancedCluster) (*mongodbatlas.AdvancedCluster, workflow.Result) {
	resultingDeployment, err := MergedAdvancedDeployment(*advancedDeployment, deployment.Spec)
	if err != nil {
		return advancedDeployment, workflow.Terminate(workflow.Internal, err.Error())
	}

	if done := AdvancedDeploymentsEqual(ctx.Log, *advancedDeployment, resultingDeployment); done {
		return advancedDeployment, workflow.OK()
	}

	if deployment.Spec.AdvancedDeploymentSpec.Paused != nil {
		if advancedDeployment.Paused == nil || *advancedDeployment.Paused != *deployment.Spec.AdvancedDeploymentSpec.Paused {
			// paused is different from Atlas
			// we need to first send a special (un)pause request before reconciling everything else
			resultingDeployment = mongodbatlas.AdvancedCluster{
				Paused: deployment.Spec.AdvancedDeploymentSpec.Paused,
			}
		} else {
			// otherwise, don't send the paused field
			resultingDeployment.Paused = nil
		}
	}

	resultingDeployment = cleanupAdvancedDeployment(resultingDeployment)

	advancedDeployment, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, deployment.Spec.AdvancedDeploymentSpec.Name, &resultingDeployment)
	if err != nil {
		return advancedDeployment, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

func cleanupAdvancedDeployment(deployment mongodbatlas.AdvancedCluster) mongodbatlas.AdvancedCluster {
	deployment.ID = ""
	deployment.MongoDBVersion = ""
	deployment.StateName = ""
	deployment.ConnectionStrings = nil
	return deployment
}

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

// AdvancedDeploymentsEqual compares two Atlas Advanced Deployments
func AdvancedDeploymentsEqual(log *zap.SugaredLogger, deploymentAtlas mongodbatlas.AdvancedCluster, deploymentOperator mongodbatlas.AdvancedCluster) bool {
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
