package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasStreamConnectionBySecret = ".spec.kafkaConfig" //nolint:gosec
)

func NewAtlasStreamConnectionsBySecretIndex(ctx context.Context, logger *zap.SugaredLogger, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasStreamConnection{},
		AtlasStreamConnectionBySecret,
		AtlasStreamConnectionsBySecretIndices(logger.Named("indexers").Named(AtlasStreamConnectionBySecret)),
	)
}

func AtlasStreamConnectionsBySecretIndices(logger *zap.SugaredLogger) client.IndexerFunc {
	return func(object client.Object) []string {
		streamConnection, ok := object.(*akov2.AtlasStreamConnection)
		if !ok {
			logger.Errorf("expected *akov2.AtlasStreamConnection but got %T", object)
			return nil
		}

		var indexes []string

		key, found := credentialSecretKey(streamConnection)
		if found {
			indexes = append(indexes, key)
		}

		key, found = certificateSecretKey(streamConnection)
		if found {
			indexes = append(indexes, key)
		}

		return indexes
	}
}

func credentialSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil || connection.Spec.KafkaConfig.Authentication.Credentials.Name == "" {
		return "", false
	}

	credentialsKey := connection.Spec.KafkaConfig.Authentication.Credentials.GetObject(connection.GetNamespace())

	return credentialsKey.String(), true
}

func certificateSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil || connection.Spec.KafkaConfig.Security.Certificate.Name == "" {
		return "", false
	}

	certificateKey := connection.Spec.KafkaConfig.Security.Certificate.GetObject(connection.GetNamespace())

	return certificateKey.String(), true
}
