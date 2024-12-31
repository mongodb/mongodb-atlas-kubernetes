package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasDatabaseUserCredentialsIndex = "atlasdatabaseuser.credentials"
)

func NewAtlasDatabaseUserByCredentialIndexer(logger *zap.Logger) *LocalCredentialIndexer {
	return NewLocalCredentialsIndexer(AtlasDatabaseUserCredentialsIndex, &akov2.AtlasDatabaseUser{}, logger)
}

func DatabaseUserRequests(list *akov2.AtlasDatabaseUserList) []reconcile.Request {
	requests := make([]reconcile.Request, 0, len(list.Items))
	for _, item := range list.Items {
		requests = append(requests, toRequest(&item))
	}
	return requests
}
