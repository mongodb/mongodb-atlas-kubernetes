package indexer

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasFederatedAuthBySecretsIndex = "atlasfederatedauth.spec.connectionSecret"
)

type AtlasFederatedAuthBySecretsIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasFederatedAuthBySecretsIndexer(logger *zap.Logger) *AtlasFederatedAuthBySecretsIndexer {
	return &AtlasFederatedAuthBySecretsIndexer{
		logger: logger.Named(AtlasFederatedAuthBySecretsIndex).Sugar(),
	}
}

func (*AtlasFederatedAuthBySecretsIndexer) Object() client.Object {
	return &akov2.AtlasFederatedAuth{}
}

func (*AtlasFederatedAuthBySecretsIndexer) Name() string {
	return AtlasFederatedAuthBySecretsIndex
}

func (a *AtlasFederatedAuthBySecretsIndexer) Keys(object client.Object) []string {
	fedAuth, ok := object.(*akov2.AtlasFederatedAuth)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasFederatedAuth but got %T", object)
		return nil
	}

	if fedAuth.Spec.ConnectionSecretRef.IsEmpty() {
		return nil
	}

	return []string{fedAuth.ConnectionSecretObjectKey().String()}
}
