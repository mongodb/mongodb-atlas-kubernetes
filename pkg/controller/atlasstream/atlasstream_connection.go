package atlasstream

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

const kafkaConnectionAuthUsername = "username"
const kafkaConnectionAuthPassword = "password"
const kafkaConnectionSecCertificate = "certificate"

type streamConnectionMapper func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error)

// this is the dispatcher of connection registry management
func (r *InstanceReconciler) handleConnectionRegistry(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	akoStreamInstance *akov2.AtlasStreamInstance,
	atlasStreamInstance *admin.StreamsTenant,
) (ctrl.Result, error) {
	toCreate, toUpdate, toDelete, err := r.sortConnectionRegistryTasks(ctx, akoStreamInstance, atlasStreamInstance)
	if err != nil {
		return r.terminate(ctx, workflow.StreamConnectionNotConfigured, err)
	}

	// we do all operations in a single flow and only return earlier in case of failure
	if len(toCreate) > 0 {
		// if there are connection to be added to the instance
		err = createConnections(ctx, project, akoStreamInstance, toCreate, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotCreated, err)
		}
	}

	if len(toUpdate) > 0 {
		// if there are connection to be updated in the instance
		err = updateConnections(ctx, project, akoStreamInstance, toUpdate, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotUpdated, err)
		}
	}

	if len(toDelete) > 0 {
		// if there are connection to be deleted from the instance
		err = deleteConnections(ctx, project, akoStreamInstance, toDelete)
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotRemoved, err)
		}
	}

	// we can transition straight away to ready state
	return r.ready(ctx, atlasStreamInstance)
}

func (r *InstanceReconciler) sortConnectionRegistryTasks(
	ctx *workflow.Context,
	akoStreamInstance *akov2.AtlasStreamInstance,
	atlasStreamInstance *admin.StreamsTenant,
) ([]*akov2.AtlasStreamConnection, []*akov2.AtlasStreamConnection, []*admin.StreamsConnection, error) {
	toCreate := make([]*akov2.AtlasStreamConnection, 0, len(akoStreamInstance.Spec.ConnectionRegistry))
	toUpdate := make([]*akov2.AtlasStreamConnection, 0, len(akoStreamInstance.Spec.ConnectionRegistry))
	toDelete := make([]*admin.StreamsConnection, 0, len(atlasStreamInstance.GetConnections()))

	atlasConnectionByName := map[string]*admin.StreamsConnection{}
	akoConnectionByName := map[string]*akov2.AtlasStreamConnection{}
	atlasConnections := atlasStreamInstance.GetConnections()

	for i := range atlasConnections {
		atlasConnectionByName[atlasConnections[i].GetName()] = &atlasConnections[i]
	}

	for _, akoConnectionRef := range akoStreamInstance.Spec.ConnectionRegistry {
		akoConnection := akov2.AtlasStreamConnection{}
		err := r.Client.Get(ctx.Context, *akoConnectionRef.GetObject(akoStreamInstance.Namespace), &akoConnection)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to retrieve connection %v: %w", akoConnectionRef, err)
		}

		akoConnectionByName[akoConnection.Spec.Name] = &akoConnection

		// when connection doesn't exist in Atlas, we need to add it hence connection registry has changed and we return earlier
		atlasConnection, ok := atlasConnectionByName[akoConnection.Spec.Name]
		if !ok {
			toCreate = append(toCreate, &akoConnection)

			continue
		}

		hasConnectionChanged, err := hasStreamConnectionChanged(&akoConnection, *atlasConnection, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return nil, nil, nil, err
		}

		// when connection has been modified, we need to update it hence connection registry has changed and we return earlier
		if hasConnectionChanged {
			toUpdate = append(toUpdate, &akoConnection)
		}
	}

	for i := range atlasConnections {
		if _, ok := akoConnectionByName[atlasConnections[i].GetName()]; !ok {
			toDelete = append(toDelete, &atlasConnections[i])
		}
	}

	return toCreate, toUpdate, toDelete, nil
}

func streamConnectionToAtlas(ctx context.Context, k8sClient client.Client) streamConnectionMapper {
	return func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error) {
		connection := admin.StreamsConnection{
			Name: &streamConnection.Spec.Name,
			Type: &streamConnection.Spec.ConnectionType,
		}

		if streamConnection.Spec.ClusterConfig != nil {
			connection.ClusterName = &streamConnection.Spec.ClusterConfig.Name
			connection.DbRoleToExecute = &admin.DBRoleToExecute{
				Role: &streamConnection.Spec.ClusterConfig.Role.Name,
				Type: &streamConnection.Spec.ClusterConfig.Role.RoleType,
			}
		}

		if streamConnection.Spec.KafkaConfig != nil {
			authData, err := getSecretData(
				ctx,
				k8sClient,
				*streamConnection.Spec.KafkaConfig.Authentication.Credentials.GetObject(streamConnection.Namespace),
				kafkaConnectionAuthUsername,
				kafkaConnectionAuthPassword,
			)
			if err != nil {
				return nil, err
			}

			secData, err := getSecretData(
				ctx,
				k8sClient,
				*streamConnection.Spec.KafkaConfig.Security.Certificate.GetObject(streamConnection.Namespace),
				kafkaConnectionSecCertificate,
			)
			if err != nil {
				return nil, err
			}

			connection.BootstrapServers = &streamConnection.Spec.KafkaConfig.BootstrapServers
			connection.Authentication = &admin.StreamsKafkaAuthentication{
				Mechanism: &streamConnection.Spec.KafkaConfig.Authentication.Mechanism,
				Username:  pointer.MakePtr(authData[kafkaConnectionAuthUsername]),
				Password:  pointer.MakePtr(authData[kafkaConnectionAuthPassword]),
			}
			connection.Security = &admin.StreamsKafkaSecurity{
				BrokerPublicCertificate: pointer.MakePtr(secData[kafkaConnectionSecCertificate]),
				Protocol:                &streamConnection.Spec.KafkaConfig.Security.Protocol,
			}
			connection.Config = &streamConnection.Spec.KafkaConfig.Config
		}

		return &connection, nil
	}
}

func getSecretData(ctx context.Context, k8sClient client.Client, ref client.ObjectKey, keys ...string) (map[string]string, error) {
	secret := corev1.Secret{}
	err := k8sClient.Get(ctx, ref, &secret)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret %v: %w", ref, err)
	}

	data := map[string]string{}
	for _, key := range keys {
		val, ok := secret.Data[key]
		if !ok {
			return nil, fmt.Errorf("key %s is not present in the secret %v", key, ref)
		}

		data[key] = string(val)
	}

	return data, nil
}

func hasStreamConnectionChanged(
	streamConnection *akov2.AtlasStreamConnection,
	atlasStreamConnection admin.StreamsConnection,
	mapper streamConnectionMapper,
) (bool, error) {
	// Unset API metadata for comparison
	atlasStreamConnection.Links = nil
	connection, err := mapper(streamConnection)
	if err != nil {
		return false, err
	}
	if _, ok := connection.GetAuthenticationOk(); ok {
		connection.Authentication.Password = nil
	}

	return !reflect.DeepEqual(*connection, atlasStreamConnection), nil
}
