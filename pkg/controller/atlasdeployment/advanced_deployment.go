package atlasdeployment

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	handleAutoscaling(deployment.Spec.AdvancedDeploymentSpec)

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

	// Prevent changing of instanceSize and diskSizeGB if autoscaling is enabled

	cleanupTheSpec(&specDeployment)

	deploymentAsAtlas, err := specDeployment.ToAtlas()
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

	atlasDeploymentAsAtlas, _, err = ctx.Client.AdvancedClusters.Update(context.Background(), project.Status.ID, deployment.Spec.AdvancedDeploymentSpec.Name, deploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.DeploymentNotUpdatedInAtlas, err.Error())
	}

	return nil, workflow.InProgress(workflow.DeploymentUpdating, "deployment is updating")
}

func cleanupTheSpec(deployment *mdbv1.AdvancedDeploymentSpec) {
	deployment.MongoDBVersion = ""
}

// This will prevent from setting diskSizeGB if at least one region config has enabled disk size autoscaling
// It will also prevent from setting ANY of (electable | analytics | readonly) specs
//
//	if region config has enabled compute autoscaling
func handleAutoscaling(kubeDeployment *mdbv1.AdvancedDeploymentSpec) {
	isDiskAutoScaled := false
	syncInstanceSize := func(s *mdbv1.Specs, as *mdbv1.AdvancedAutoScalingSpec) {
		if s != nil {
			s.InstanceSize = normalizeInstanceSize(s.InstanceSize, as)
		}
	}
	for _, repSpec := range kubeDeployment.ReplicationSpecs {
		for _, regConfig := range repSpec.RegionConfigs {
			if regConfig.AutoScaling != nil {
				if regConfig.AutoScaling.DiskGB != nil &&
					regConfig.AutoScaling.DiskGB.Enabled != nil &&
					*regConfig.AutoScaling.DiskGB.Enabled {
					isDiskAutoScaled = true
				}

				if regConfig.AutoScaling.Compute != nil &&
					regConfig.AutoScaling.Compute.Enabled != nil &&
					*regConfig.AutoScaling.Compute.Enabled {
					syncInstanceSize(regConfig.ElectableSpecs, regConfig.AutoScaling)
					syncInstanceSize(regConfig.AnalyticsSpecs, regConfig.AutoScaling)
					syncInstanceSize(regConfig.ReadOnlySpecs, regConfig.AutoScaling)
				}
			}
		}
	}

	if isDiskAutoScaled {
		kubeDeployment.DiskSizeGB = nil
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

func normalizeInstanceSize(currentInstanceSize string, autoscaling *mdbv1.AdvancedAutoScalingSpec) string {
	currentSize := extractNumberFromInstanceTypeName(currentInstanceSize)
	minSize := extractNumberFromInstanceTypeName(autoscaling.Compute.MinInstanceSize)
	maxSize := extractNumberFromInstanceTypeName(autoscaling.Compute.MaxInstanceSize)

	if currentSize < minSize {
		return autoscaling.Compute.MinInstanceSize
	}

	if currentSize > maxSize {
		return autoscaling.Compute.MaxInstanceSize
	}

	return currentInstanceSize
}

// extractNumberFromInstanceTypeName get the existing number from a given instance type name, fail when the name is incorrect
func extractNumberFromInstanceTypeName(instanceTypeName string) int {
	name := strings.TrimPrefix(instanceTypeName, "M")
	name = strings.TrimPrefix(name, "R")
	name = strings.TrimSuffix(name, "_NVME")

	number, _ := strconv.Atoi(name)

	return number
}
