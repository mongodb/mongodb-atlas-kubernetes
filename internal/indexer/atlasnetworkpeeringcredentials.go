package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasNetworkPeeringCredentialsIndex = "atlasnetworkpeering.credentials"
)

func NewAtlasNetworkPeeringByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasNetworkPeeringCredentialsIndex, &akov2.AtlasNetworkPeering{}, logger)
}

func NetworkPeeringRequests(list *akov2.AtlasNetworkPeeringList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
