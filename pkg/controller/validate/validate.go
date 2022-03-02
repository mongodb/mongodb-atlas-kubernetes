package validate

import (
	"errors"

	"github.com/hashicorp/go-multierror"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
)

func ClusterSpec(clusterSpec mdbv1.AtlasClusterSpec) error {
	var err error
	if clusterSpec.AdvancedClusterSpec == nil && clusterSpec.ClusterSpec == nil {
		err = multierror.Append(err, errors.New("expected exactly one of spec.clusterSpec or spec.advancedClusterSpec, neither were present"))
	}

	if clusterSpec.AdvancedClusterSpec != nil && clusterSpec.ClusterSpec != nil {
		err = multierror.Append(err, errors.New("expected exactly one of spec.clusterSpec or spec.advancedClusterSpec, both were present"))
	}
	return err
}

func Project(_ *mdbv1.AtlasProject) error {
	return nil
}

func DatabaseUser(_ *mdbv1.AtlasDatabaseUser) error {
	return nil
}
