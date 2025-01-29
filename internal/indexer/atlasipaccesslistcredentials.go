package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasIPAccessListCredentialsIndex = "atlasipaccesslist.credentials"
)

func NewAtlasIPAccessListCredentialsByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasIPAccessListCredentialsIndex, &akov2.AtlasIPAccessList{}, logger)
}

func IPAccessListRequests(list *akov2.AtlasIPAccessListList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
