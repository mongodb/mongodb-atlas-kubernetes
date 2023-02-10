package validate

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/hashicorp/go-multierror"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func DeploymentSpec(deploymentSpec mdbv1.AtlasDeploymentSpec) error {
	var err error

	if allAreNil(deploymentSpec.AdvancedDeploymentSpec, deploymentSpec.ServerlessSpec, deploymentSpec.DeploymentSpec) {
		err = multierror.Append(err, errors.New("expected exactly one of spec.deploymentSpec or spec.advancedDepploymentSpec or spec.serverlessSpec to be present, but none were"))
	}

	if moreThanOneIsNonNil(deploymentSpec.AdvancedDeploymentSpec, deploymentSpec.ServerlessSpec, deploymentSpec.DeploymentSpec) {
		err = multierror.Append(err, errors.New("expected exactly one of spec.deploymentSpec, spec.advancedDepploymentSpec or spec.serverlessSpec, more than one were present"))
	}

	if deploymentSpec.DeploymentSpec != nil {
		if deploymentSpec.DeploymentSpec.ProviderSettings != nil && (deploymentSpec.DeploymentSpec.ProviderSettings.InstanceSizeName == "" && deploymentSpec.DeploymentSpec.ProviderSettings.ProviderName != "SERVERLESS") {
			err = multierror.Append(err, errors.New("must specify instanceSizeName if provider name is not SERVERLESS"))
		}
		if deploymentSpec.DeploymentSpec.ProviderSettings != nil && (deploymentSpec.DeploymentSpec.ProviderSettings.InstanceSizeName != "" && deploymentSpec.DeploymentSpec.ProviderSettings.ProviderName == "SERVERLESS") {
			err = multierror.Append(err, errors.New("must not specify instanceSizeName if provider name is SERVERLESS"))
		}
	}

	if deploymentSpec.AdvancedDeploymentSpec != nil {
		instanceSizeErr := instanceSizeForAdvancedDeployment(deploymentSpec.AdvancedDeploymentSpec.ReplicationSpecs)
		if instanceSizeErr != nil {
			err = multierror.Append(err, instanceSizeErr)
		}

		autoscalingErr := autoscalingForAdvancedDeployment(deploymentSpec.AdvancedDeploymentSpec.ReplicationSpecs)
		if autoscalingErr != nil {
			err = multierror.Append(err, autoscalingErr)
		}
	}

	return err
}

func Project(project *mdbv1.AtlasProject) error {
	if err := projectCustomRoles(project.Spec.CustomRoles); err != nil {
		return err
	}

	return nil
}

func DatabaseUser(_ *mdbv1.AtlasDatabaseUser) error {
	return nil
}

func BackupSchedule(bSchedule *mdbv1.AtlasBackupSchedule, deployment *mdbv1.AtlasDeployment) error {
	var err error

	if bSchedule.Spec.Export == nil && bSchedule.Spec.AutoExportEnabled {
		err = multierror.Append(err, errors.New("you must specify export policy when auto export is enabled"))
	}

	replicaSets := map[string]struct{}{}
	if deployment.Status.ReplicaSets != nil {
		for _, replicaSet := range deployment.Status.ReplicaSets {
			replicaSets[replicaSet.ID] = struct{}{}
		}
	}

	for position, copySetting := range bSchedule.Spec.CopySettings {
		if copySetting.RegionName == nil {
			err = multierror.Append(err, fmt.Errorf("copy setting at position %d: you must set a region name", position))
		}

		if copySetting.ReplicationSpecID == nil {
			err = multierror.Append(err, fmt.Errorf("copy setting at position %d: you must set a valid ReplicationSpecID", position))
		} else if _, ok := replicaSets[*copySetting.ReplicationSpecID]; !ok {
			err = multierror.Append(err, fmt.Errorf("copy setting at position %d: referenced ReplicationSpecID is invalid", position))
		}

		if copySetting.ShouldCopyOplogs != nil && *copySetting.ShouldCopyOplogs {
			if deployment.Spec.AdvancedDeploymentSpec != nil &&
				(deployment.Spec.AdvancedDeploymentSpec.PitEnabled == nil ||
					!*deployment.Spec.AdvancedDeploymentSpec.PitEnabled) {
				err = multierror.Append(err, fmt.Errorf("copy setting at position %d: you must enable pit before enable copyOplogs", position))
			}

			if deployment.Spec.DeploymentSpec != nil &&
				(deployment.Spec.DeploymentSpec.PitEnabled == nil ||
					!*deployment.Spec.DeploymentSpec.PitEnabled) {
				err = multierror.Append(err, fmt.Errorf("copy setting at position %d: you must enable pit before enable copyOplogs", position))
			}
		}
	}

	return err
}

func getNonNilCount(values ...interface{}) int {
	nonNilCount := 0
	for _, v := range values {
		if !reflect.ValueOf(v).IsNil() {
			nonNilCount += 1
		}
	}
	return nonNilCount
}

// allAreNil returns true if all elements are nil.
func allAreNil(values ...interface{}) bool {
	return getNonNilCount(values...) == 0
}

// moreThanOneIsNil returns true if there are more than one non nil elements.
func moreThanOneIsNonNil(values ...interface{}) bool {
	return getNonNilCount(values...) > 1
}

func instanceSizeForAdvancedDeployment(replicationSpecs []*mdbv1.AdvancedReplicationSpec) error {
	var instanceSize string
	err := errors.New("instance size must be the same for all nodes in all regions and across all replication specs for advanced deployment ")

	isInstanceSizeEqual := func(nodeInstanceType string) bool {
		if instanceSize == "" {
			instanceSize = nodeInstanceType
		}

		return nodeInstanceType == instanceSize
	}

	for _, replicationSpec := range replicationSpecs {
		for _, regionSpec := range replicationSpec.RegionConfigs {
			if instanceSize == "" {
				instanceSize = regionSpec.ElectableSpecs.InstanceSize
			}

			if regionSpec.ElectableSpecs != nil && !isInstanceSizeEqual(regionSpec.ElectableSpecs.InstanceSize) {
				return err
			}

			if regionSpec.ReadOnlySpecs != nil && !isInstanceSizeEqual(regionSpec.ReadOnlySpecs.InstanceSize) {
				return err
			}

			if regionSpec.AnalyticsSpecs != nil && !isInstanceSizeEqual(regionSpec.AnalyticsSpecs.InstanceSize) {
				return err
			}
		}
	}

	return nil
}

func autoscalingForAdvancedDeployment(replicationSpecs []*mdbv1.AdvancedReplicationSpec) error {
	var autoscaling *mdbv1.AdvancedAutoScalingSpec
	first := true

	for _, replicationSpec := range replicationSpecs {
		for _, regionSpec := range replicationSpec.RegionConfigs {
			if first {
				autoscaling = regionSpec.AutoScaling
				first = false
			}

			if cmp.Diff(autoscaling, regionSpec.AutoScaling, cmpopts.EquateEmpty()) != "" {
				return errors.New("autoscaling must be the same for all regions and across all replication specs for advanced deployment ")
			}
		}
	}

	return nil
}

func projectCustomRoles(customRoles []mdbv1.CustomRole) error {
	if len(customRoles) == 0 {
		return nil
	}

	var err error
	customRolesMap := map[string]struct{}{}

	for _, customRole := range customRoles {
		if _, ok := customRolesMap[customRole.Name]; ok {
			err = multierror.Append(err, fmt.Errorf("the custom rone \"%s\" is duplicate. custom role name must be unique", customRole.Name))
		}

		customRolesMap[customRole.Name] = struct{}{}
	}

	return err
}
