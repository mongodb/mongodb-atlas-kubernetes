// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasstream

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

const (
	kafkaConnectionAuthUsername   = "username"
	kafkaConnectionAuthPassword   = "password"
	kafkaConnectionSecCertificate = "certificate"
)

type streamConnectionMapper func(streamConnection *akov2.AtlasStreamConnection) (*admin.StreamsConnection, error)

type streamConnectionOperations struct {
	NoOp   []*akov2.AtlasStreamConnection
	Create []*akov2.AtlasStreamConnection
	Update []*akov2.AtlasStreamConnection
	Delete []*admin.StreamsConnection
}

func newStreamConnectionOperations(
	akoStreamInstance *akov2.AtlasStreamInstance,
	atlasStreamConnections []admin.StreamsConnection,
) *streamConnectionOperations {
	return &streamConnectionOperations{
		NoOp:   make([]*akov2.AtlasStreamConnection, 0, len(akoStreamInstance.Spec.ConnectionRegistry)),
		Create: make([]*akov2.AtlasStreamConnection, 0, len(akoStreamInstance.Spec.ConnectionRegistry)),
		Update: make([]*akov2.AtlasStreamConnection, 0, len(akoStreamInstance.Spec.ConnectionRegistry)),
		Delete: make([]*admin.StreamsConnection, 0, len(atlasStreamConnections)),
	}
}

// handleConnectionRegistry is the dispatcher of connection registry management
func (r *AtlasStreamsInstanceReconciler) handleConnectionRegistry(
	ctx *workflow.Context,
	project *akov2.AtlasProject,
	akoStreamInstance *akov2.AtlasStreamInstance,
	atlasStreamInstance *admin.StreamsTenant,
) (ctrl.Result, error) {
	streamConnections, err := paging.ListAll(ctx.Context, func(c context.Context, pageNum int) (paging.Response[admin.StreamsConnection], *http.Response, error) {
		return ctx.SdkClientSet.SdkClient20250312009.StreamsApi.
			ListStreamConnections(c, project.ID(), akoStreamInstance.Spec.Name).
			PageNum(pageNum).
			Execute()
	})
	if err != nil {
		return r.terminate(ctx, workflow.StreamConnectionNotConfigured, err)
	}

	ops, err := r.sortConnectionRegistryTasks(ctx, akoStreamInstance, streamConnections)
	if err != nil {
		return r.terminate(ctx, workflow.StreamConnectionNotConfigured, err)
	}

	// we do all operations in a single flow and only return earlier in case of failure
	if len(ops.Create) > 0 {
		// if there are connection to be added to the instance
		err = createConnections(ctx, project, akoStreamInstance, ops.Create, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotCreated, err)
		}
	}

	if len(ops.Update) > 0 {
		// if there are connection to be updated in the instance
		err = updateConnections(ctx, project, akoStreamInstance, ops.Update, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotUpdated, err)
		}
	}

	if len(ops.Delete) > 0 {
		// if there are connection to be deleted from the instance
		err = deleteConnections(ctx, project, akoStreamInstance, ops.Delete)
		if err != nil {
			return r.terminate(ctx, workflow.StreamConnectionNotRemoved, err)
		}
	}

	for i := range ops.NoOp {
		akoStreamConnection := ops.NoOp[i]
		ctx.EnsureStatusOption(
			status.AtlasStreamInstanceAddConnection(
				akoStreamConnection.Spec.Name,
				common.ResourceRefNamespaced{
					Name:      akoStreamConnection.Name,
					Namespace: akoStreamConnection.Namespace,
				},
			),
		)
	}

	if err = customresource.ManageFinalizer(ctx.Context, r.Client, akoStreamInstance, customresource.SetFinalizer); err != nil {
		return r.terminate(ctx, workflow.AtlasFinalizerNotSet, err)
	}

	// we can transition straight away to ready state
	return r.ready(ctx, atlasStreamInstance)
}

func (r *AtlasStreamsInstanceReconciler) sortConnectionRegistryTasks(
	ctx *workflow.Context,
	akoStreamInstance *akov2.AtlasStreamInstance,
	atlasStreamConnections []admin.StreamsConnection,
) (*streamConnectionOperations, error) {
	ops := newStreamConnectionOperations(akoStreamInstance, atlasStreamConnections)

	atlasConnectionByName := map[string]*admin.StreamsConnection{}
	akoConnectionByName := map[string]*akov2.AtlasStreamConnection{}

	for i := range atlasStreamConnections {
		atlasConnectionByName[atlasStreamConnections[i].GetName()] = &atlasStreamConnections[i]
	}

	for _, akoConnectionRef := range akoStreamInstance.Spec.ConnectionRegistry {
		akoConnection := akov2.AtlasStreamConnection{}
		err := r.Client.Get(ctx.Context, *akoConnectionRef.GetObject(akoStreamInstance.Namespace), &akoConnection)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve connection %v: %w", akoConnectionRef, err)
		}

		akoConnectionByName[akoConnection.Spec.Name] = &akoConnection

		// when connection doesn't exist in Atlas, we need to add it hence connection registry has changed and we return earlier
		atlasConnection, ok := atlasConnectionByName[akoConnection.Spec.Name]
		if !ok {
			ops.Create = append(ops.Create, &akoConnection)

			continue
		}

		hasConnectionChanged, err := hasStreamConnectionChanged(&akoConnection, atlasConnection, streamConnectionToAtlas(ctx.Context, r.Client))
		if err != nil {
			return nil, err
		}

		// when connection has been modified, we need to update it hence connection registry has changed and we return earlier
		if hasConnectionChanged {
			ops.Update = append(ops.Update, &akoConnection)

			continue
		}

		ops.NoOp = append(ops.NoOp, &akoConnection)
	}

	for i := range atlasStreamConnections {
		if _, ok := akoConnectionByName[atlasStreamConnections[i].GetName()]; !ok {
			ops.Delete = append(ops.Delete, &atlasStreamConnections[i])
		}
	}

	return ops, nil
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
	atlasStreamConnection *admin.StreamsConnection,
	mapper streamConnectionMapper,
) (bool, error) {
	// Unset API metadata for comparison
	atlasStreamConnection.Links = nil

	if _, ok := atlasStreamConnection.GetDbRoleToExecuteOk(); ok {
		atlasStreamConnection.DbRoleToExecute.Links = nil
	}

	if _, ok := atlasStreamConnection.GetAuthenticationOk(); ok {
		atlasStreamConnection.Authentication.Links = nil
	}

	if _, ok := atlasStreamConnection.GetSecurityOk(); ok {
		atlasStreamConnection.Security.Links = nil
	}

	connection, err := mapper(streamConnection)
	if err != nil {
		return false, err
	}
	if _, ok := connection.GetAuthenticationOk(); ok {
		connection.Authentication.Password = nil
	}

	return !reflect.DeepEqual(*connection, *atlasStreamConnection), nil
}
