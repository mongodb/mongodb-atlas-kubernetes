package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasCustomRoleCredentialsIndex = "atlascustomrole.credentials"
)

func NewAtlasCustomRoleByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasCustomRoleCredentialsIndex, &akov2.AtlasCustomRole{}, logger)
}

func CustomRoleRequests(list *akov2.AtlasCustomRoleList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
