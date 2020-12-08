package atlasproject

import (
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"go.uber.org/zap"
)

func configureIpAccessList(connection atlas.Connection, project *mdbv1.AtlasProject, log *zap.SugaredLogger) error {
	client, err := atlas.AtlasClient(connection, log)
	if err != nil {
		return err
	}

}
