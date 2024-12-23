package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasPrivateEndpointCredentialsIndex = "atlasprivateendpoint.credentials"
)

func NewAtlasPrivateEndpointByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasPrivateEndpointCredentialsIndex, &akov2.AtlasPrivateEndpoint{}, logger)
}

func PrivateEndpointRequests(list *akov2.AtlasPrivateEndpointList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
