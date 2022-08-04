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
	specDeployment, atlasDeployment, err := MergedAdvancedDeployment(*atlasDeploymentAsAtlas, *deployment.Spec.AdvancedDeploymentSpec)
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

	deploymentAsAtlas, err := cleanupTheSpec(specDeployment).ToAtlas()
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	atlasDeploymentAsAtlas, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, deployment.Spec.AdvancedDeploymentSpec.Name, deploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

func cleanupTheSpec(deployment mdbv1.AdvancedDeploymentSpec) *mdbv1.AdvancedDeploymentSpec {
	deployment.MongoDBVersion = ""
	deployment = removeAdvancedDeploymentOutdatedFields(&deployment, nil)
	return &deployment
}

// removeOutdatedFields unsets fields which are should be empty based on flags
func removeAdvancedDeploymentOutdatedFields(removeFrom *mdbv1.AdvancedDeploymentSpec, lookAt *mdbv1.AdvancedDeploymentSpec) mdbv1.AdvancedDeploymentSpec {
	if lookAt == nil {
		lookAt = removeFrom
	}

	result := *removeFrom

	for i := range lookAt.ReplicationSpecs {
		for j := range lookAt.ReplicationSpecs[i].RegionConfigs {
			regionConfig := lookAt.ReplicationSpecs[i].RegionConfigs[j]
			if regionConfig.AutoScaling != nil && regionConfig.AutoScaling.Compute != nil {
				if *regionConfig.AutoScaling.Compute.Enabled {
					// unset InstanceSize from all specs when compute autoscaling is enabled
					resultRegionConfig := result.ReplicationSpecs[i].RegionConfigs[j]

					unsetInstanceSize(resultRegionConfig.ElectableSpecs)
					unsetInstanceSize(resultRegionConfig.ReadOnlySpecs)
					unsetInstanceSize(resultRegionConfig.AnalyticsSpecs)
				} else {
					// Shouldn't be able to set max, min InstanceSize if autoscaling is not enabled
					result.ReplicationSpecs[i].RegionConfigs[j].AutoScaling.Compute = &mdbv1.ComputeSpec{}
				}
			}
			if regionConfig.AutoScaling.DiskGB != nil && *regionConfig.AutoScaling.DiskGB.Enabled {
				// unset diskSizeGB when if disk autoscaling is enabled
				result.DiskSizeGB = nil
			}
		}
	}
	return result
}

func unsetInstanceSize(spec *mdbv1.Specs) {
	if spec != nil && spec.InstanceSize != "" {
		spec.InstanceSize = ""
	}
}

// MergedAdvancedDeployment will return the result of merging AtlasDeploymentSpec with Atlas Advanced Deployment
func MergedAdvancedDeployment(atlasDeploymentAsAtlas mongodbatlas.AdvancedCluster, specDeployment mdbv1.AdvancedDeploymentSpec) (mergedDeployment mdbv1.AdvancedDeploymentSpec, atlasDeployment mdbv1.AdvancedDeploymentSpec, err error) {
	atlasDeployment, err = AdvancedDeploymentFromAtlas(atlasDeploymentAsAtlas)
	if err != nil {
		return
	}

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
			if v.BackingProviderName == "" {
				mergedDeployment.ReplicationSpecs[i].RegionConfigs[k].BackingProviderName = ""
			}
		}
	}
	return
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
	deploymentAtlas = removeAdvancedDeploymentOutdatedFields(&deploymentAtlas, &deploymentOperator)
	deploymentOperator = removeAdvancedDeploymentOutdatedFields(&deploymentOperator, nil)

	d := cmp.Diff(deploymentOperator, deploymentAtlas, cmpopts.EquateEmpty())
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
