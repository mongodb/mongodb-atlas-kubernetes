package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasDeploymentCredentialsIndex = "atlasdeployment.credentials"
)

func NewAtlasDeploymentByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasDeploymentCredentialsIndex, &akov2.AtlasDeployment{}, logger)
}

func DeploymentRequests(list *akov2.AtlasDeploymentList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
