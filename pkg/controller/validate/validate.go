package validate

import (
	"errors"
	"reflect"

	"github.com/hashicorp/go-multierror"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func ClusterSpec(clusterSpec mdbv1.AtlasDeploymentSpec) error {
	var err error

	if allAreNil(clusterSpec.AdvancedDeploymentSpec, clusterSpec.ServerlessSpec, clusterSpec.DeploymentSpec) {
		err = multierror.Append(err, errors.New("expected exactly one of spec.clusterSpec or spec.advancedClusterSpec or spec.serverlessSpec to be present, but none were"))
	}

	if moreThanOneIsNonNil(clusterSpec.AdvancedDeploymentSpec, clusterSpec.ServerlessSpec, clusterSpec.DeploymentSpec) {
		err = multierror.Append(err, errors.New("expected exactly one of spec.clusterSpec, spec.advancedClusterSpec or spec.serverlessSpec, more than one were present"))
	}

	if clusterSpec.DeploymentSpec != nil {
		if clusterSpec.DeploymentSpec.ProviderSettings != nil && (clusterSpec.DeploymentSpec.ProviderSettings.InstanceSizeName == "" && clusterSpec.DeploymentSpec.ProviderSettings.ProviderName != "SERVERLESS") {
			err = multierror.Append(err, errors.New("must specify instanceSizeName if provider name is not SERVERLESS"))
		}
		if clusterSpec.DeploymentSpec.ProviderSettings != nil && (clusterSpec.DeploymentSpec.ProviderSettings.InstanceSizeName != "" && clusterSpec.DeploymentSpec.ProviderSettings.ProviderName == "SERVERLESS") {
			err = multierror.Append(err, errors.New("must not specify instanceSizeName if provider name is SERVERLESS"))
		}
	}

	return err
}

func Project(_ *mdbv1.AtlasProject) error {
	return nil
}

func DatabaseUser(_ *mdbv1.AtlasDatabaseUser) error {
	return nil
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
