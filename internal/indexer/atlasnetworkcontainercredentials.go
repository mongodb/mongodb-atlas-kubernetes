package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasNetworkContainerCredentialsIndex = "atlasnetworkcontainer.credentials"
)

func NewAtlasNetworkContainerByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasNetworkContainerCredentialsIndex, &akov2.AtlasNetworkContainer{}, logger)
}

func NetworkContainerRequests(list *akov2.AtlasNetworkContainerList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
