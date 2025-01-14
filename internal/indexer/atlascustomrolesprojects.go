//nolint:dupl
package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasCustomRoleByProject = "atlascustomrole.spec.projectRef"
)

type AtlasCustomRoleByProjectIndexer struct {
	AtlasReferrerByProjectIndexerBase
}

func NewAtlasCustomRoleByProjectIndexer(logger *zap.Logger) *AtlasCustomRoleByProjectIndexer {
	return &AtlasCustomRoleByProjectIndexer{
		AtlasReferrerByProjectIndexerBase: *NewAtlasReferrerByProjectIndexer(
			logger,
			AtlasCustomRoleByProject,
		),
	}
}

func (*AtlasCustomRoleByProjectIndexer) Object() client.Object {
	return &akov2.AtlasCustomRole{}
}
