package indexer

import (
	"context"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

const (
	AtlasStreamConnectionByCredentialsSecret = ".spec.kafkaConfig.authentication.credentials" //nolint:gosec
	AtlasStreamConnectionByCertificateSecret = ".spec.kafkaConfig.security.certificate"       //nolint:gosec
)

func NewAtlasStreamConnectionsByCredentialSecretIndex(ctx context.Context, logger *zap.SugaredLogger, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasStreamConnection{},
		AtlasStreamConnectionByCredentialsSecret,
		AtlasStreamConnectionsBySecretIndices(logger.Named("indexers").Named(AtlasStreamConnectionByCredentialsSecret), CredentialSecretKey),
	)
}

func NewAtlasStreamConnectionsByCertificateSecretIndex(ctx context.Context, logger *zap.SugaredLogger, idx client.FieldIndexer) error {
	return idx.IndexField(ctx,
		&akov2.AtlasStreamConnection{},
		AtlasStreamConnectionByCertificateSecret,
		AtlasStreamConnectionsBySecretIndices(logger.Named("indexers").Named(AtlasStreamConnectionByCertificateSecret), CertificateSecretKey),
	)
}

func AtlasStreamConnectionsBySecretIndices(logger *zap.SugaredLogger, keyReader ConnectionSecretKeyReader) client.IndexerFunc {
	return func(object client.Object) []string {
		streamConnection, ok := object.(*akov2.AtlasStreamConnection)
		if !ok {
			logger.Errorf("expected *akov2.AtlasStreamConnection but got %T", object)
			return nil
		}

		key, found := keyReader(streamConnection)
		if !found {
			return nil
		}

		return []string{key}
	}
}

type ConnectionSecretKeyReader func(connection *akov2.AtlasStreamConnection) (string, bool)

func CredentialSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil {
		return "", false
	}

	credentialsKey := connection.Spec.KafkaConfig.Authentication.Credentials.GetObject(connection.GetNamespace())

	return credentialsKey.String(), true
}

func CertificateSecretKey(connection *akov2.AtlasStreamConnection) (string, bool) {
	if connection == nil || connection.Spec.KafkaConfig == nil || connection.Spec.KafkaConfig.Security.Certificate.Name == "" {
		return "", false
	}

	certificateKey := connection.Spec.KafkaConfig.Security.Certificate.GetObject(connection.GetNamespace())

	return certificateKey.String(), true
}
