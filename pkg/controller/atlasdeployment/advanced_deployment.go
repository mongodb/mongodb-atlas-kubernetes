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

const FreeTier = "M0"

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

	result := EnsureCustomZoneMapping(ctx, project.ID(), deployment.Spec.AdvancedDeploymentSpec.CustomZoneMapping, advancedDeployment.Name)
	if !result.IsOk() {
		return advancedDeployment, result
	}

	result = EnsureManagedNamespaces(ctx, project.ID(), deployment.Spec.AdvancedDeploymentSpec.ClusterType, deployment.Spec.AdvancedDeploymentSpec.ManagedNamespaces, advancedDeployment.Name)
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
	err := handleAutoscaling(ctx, deployment.Spec.AdvancedDeploymentSpec, atlasDeploymentAsAtlas)
	if err != nil {
		return atlasDeploymentAsAtlas, workflow.Terminate(workflow.Internal, err.Error())
	}

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

func cleanupTheSpec(specMerged *mdbv1.AdvancedDeploymentSpec) {
	specMerged.MongoDBVersion = ""

	globalInstanceSize := ""
	for i, replicationSpec := range specMerged.ReplicationSpecs {
		for k := range replicationSpec.RegionConfigs {
			regionConfig := specMerged.ReplicationSpecs[i].RegionConfigs[k]

			specs := []*mdbv1.Specs{
				regionConfig.AnalyticsSpecs,
				regionConfig.ElectableSpecs,
				regionConfig.ReadOnlySpecs,
			}

			forEachSpec := func(f func(spec *mdbv1.Specs)) {
				for _, spec := range specs {
					if spec != nil {
						f(spec)
					}
				}
			}

			forEachSpec(func(spec *mdbv1.Specs) {
				if globalInstanceSize == "" && spec.NodeCount != nil && *spec.NodeCount != 0 {
					globalInstanceSize = spec.InstanceSize
				}
			})

			forEachSpec(func(spec *mdbv1.Specs) {
				if spec.NodeCount == nil || *spec.NodeCount == 0 {
					spec.InstanceSize = globalInstanceSize
				}
			})
		}
	}
}

// This will prevent from setting diskSizeGB if at least one region config has enabled disk size autoscaling
// It will also prevent from setting ANY of (electable | analytics | readonly) specs
//
//	if region config has enabled compute autoscaling
func handleAutoscaling(ctx *workflow.Context, desiredDeployment *mdbv1.AdvancedDeploymentSpec, currentDeployment *mongodbatlas.AdvancedCluster) error {
	isDiskAutoScaled := false
	syncInstanceSize := func(s *mdbv1.Specs, as *mdbv1.AdvancedAutoScalingSpec) error {
		if s != nil {
			size, err := normalizeInstanceSize(ctx, s.InstanceSize, as)
			if err != nil {
				return err
			}

			if isInstanceSizeTheSame(currentDeployment, size) {
				size = ""
			}

			s.InstanceSize = size
		}

		return nil
	}
	for _, repSpec := range desiredDeployment.ReplicationSpecs {
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
					err := syncInstanceSize(regConfig.ElectableSpecs, regConfig.AutoScaling)
					if err != nil {
						return err
					}

					err = syncInstanceSize(regConfig.AnalyticsSpecs, regConfig.AutoScaling)
					if err != nil {
						return err
					}

					err = syncInstanceSize(regConfig.ReadOnlySpecs, regConfig.AutoScaling)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	if isDiskAutoScaled {
		desiredDeployment.DiskSizeGB = nil
	}

	return nil
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

func normalizeInstanceSize(ctx *workflow.Context, currentInstanceSize string, autoscaling *mdbv1.AdvancedAutoScalingSpec) (string, error) {
	currentSize, err := NewFromInstanceSizeName(currentInstanceSize)
	if err != nil {
		return "", err
	}

	minSize, err := NewFromInstanceSizeName(autoscaling.Compute.MinInstanceSize)
	if err != nil {
		return "", err
	}

	maxSize, err := NewFromInstanceSizeName(autoscaling.Compute.MaxInstanceSize)
	if err != nil {
		return "", err
	}

	if CompareInstanceSizes(currentSize, minSize) == -1 {
		ctx.Log.Warnf("The instance size is below the minimum autoscaling configuration. Setting it to %s. Consider update your CRD", autoscaling.Compute.MinInstanceSize)
		return autoscaling.Compute.MinInstanceSize, nil
	}

	if CompareInstanceSizes(currentSize, maxSize) == 1 {
		ctx.Log.Warnf("The instance size is above the maximum autoscaling configuration. Setting it to %s. Consider update your CRD", autoscaling.Compute.MaxInstanceSize)
		return autoscaling.Compute.MaxInstanceSize, nil
	}

	return currentInstanceSize, nil
}

func isInstanceSizeTheSame(currentDeployment *mongodbatlas.AdvancedCluster, desiredInstanceSize string) bool {
	return desiredInstanceSize == currentDeployment.ReplicationSpecs[0].RegionConfigs[0].ElectableSpecs.InstanceSize
}
