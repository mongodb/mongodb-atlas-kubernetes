package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestAtlasStreamConnectionsBySecretIndices(t *testing.T) {
	t.Run("should return nil when indexing wrong type object", func(t *testing.T) {
		core, logs := observer.New(zap.DebugLevel)
		indexer := NewAtlasStreamConnectionBySecretIndexer(zap.New(core))
		keys := indexer.Keys(&akov2.AtlasProject{})
		assert.Nil(t, keys)
		assert.Equal(t, 1, logs.Len())
		assert.Equal(t, zap.ErrorLevel, logs.All()[0].Level)
		assert.Equal(t, "expected *akov2.AtlasStreamConnection but got *v1.AtlasProject", logs.All()[0].Message)
	})

	t.Run("should return nil when connection has no secret", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "instance-0",
				ConnectionType: "Sample",
			},
		}

		indexer := NewAtlasStreamConnectionBySecretIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(connection)
		assert.Nil(t, keys)
	})

	t.Run("should return indexes slice when connection has credentials", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "connection-0",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Credentials: common.ResourceRefNamespaced{
							Name:      "connection-credentials",
							Namespace: "default",
						},
					},
				},
			},
		}

		indexer := NewAtlasStreamConnectionBySecretIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(connection)
		assert.Equal(
			t,
			[]string{
				"default/connection-credentials",
			},
			keys,
		)
	})

	t.Run("should return indexes slice when connection has certificate", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "connection-0",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Security: akov2.StreamsKafkaSecurity{
						Certificate: common.ResourceRefNamespaced{
							Name:      "connection-certificate",
							Namespace: "default",
						},
					},
				},
			},
		}

		indexer := NewAtlasStreamConnectionBySecretIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(connection)
		assert.Equal(
			t,
			[]string{
				"default/connection-certificate",
			},
			keys,
		)
	})

	t.Run("should return nil when connection has different secrets for credentials and certificate", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "connection-0",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Credentials: common.ResourceRefNamespaced{
							Name:      "connection-credentials",
							Namespace: "default",
						},
					},
					Security: akov2.StreamsKafkaSecurity{
						Certificate: common.ResourceRefNamespaced{
							Name:      "connection-certificate",
							Namespace: "default",
						},
					},
				},
			},
		}

		indexer := NewAtlasStreamConnectionBySecretIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(connection)
		assert.Equal(
			t,
			[]string{
				"default/connection-credentials",
				"default/connection-certificate",
			},
			keys,
		)
	})

	t.Run("should return nil when connection has the same secrets for credentials and certificate", func(t *testing.T) {
		connection := &akov2.AtlasStreamConnection{
			Spec: akov2.AtlasStreamConnectionSpec{
				Name:           "connection-0",
				ConnectionType: "Kafka",
				KafkaConfig: &akov2.StreamsKafkaConnection{
					Authentication: akov2.StreamsKafkaAuthentication{
						Credentials: common.ResourceRefNamespaced{
							Name:      "connection-secrets",
							Namespace: "default",
						},
					},
					Security: akov2.StreamsKafkaSecurity{
						Certificate: common.ResourceRefNamespaced{
							Name:      "connection-secrets",
							Namespace: "default",
						},
					},
				},
			},
		}

		indexer := NewAtlasStreamConnectionBySecretIndexer(zaptest.NewLogger(t))
		keys := indexer.Keys(connection)
		assert.Equal(
			t,
			[]string{
				"default/connection-secrets",
				"default/connection-secrets",
			},
			keys,
		)
	})
}
