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

package atlasdeployment

import (
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
)

// Status transitions:
// RESERVATION_REQUESTED -> RESERVED -> INITIATING -> AVAILABLE -> DELETING
// I assume FAILED state can be reach from any other state transition
const (
	SPEStatusDeleting  = "DELETING"
	SPEStatusReserved  = "RESERVED"
	SPEStatusAvailable = "AVAILABLE"
)

func ensureServerlessPrivateEndpoints(service *workflow.Context, projectID string, deployment *akov2.AtlasDeployment) workflow.DeprecatedResult {
	if deployment == nil || deployment.Spec.ServerlessSpec == nil {
		return workflow.Terminate(workflow.Internal, errors.New("serverless deployment spec is empty"))
	}

	deploymentSpec := deployment.Spec.ServerlessSpec

	if isGCPWithPrivateEndpoints(deploymentSpec) {
		return workflow.Terminate(workflow.AtlasUnsupportedFeature, errors.New("serverless private endpoints are not supported for GCP"))
	}

	if isGCPWithoutPrivateEndpoints(deploymentSpec) {
		return workflow.OK()
	}

	finished, err := syncServerlessPrivateEndpoints(service, projectID, deploymentSpec)
	var result workflow.DeprecatedResult
	switch {
	case err != nil:
		result = workflow.Terminate(workflow.ServerlessPrivateEndpointFailed, err)
	case err == nil && !finished:
		result = workflow.InProgress(workflow.ServerlessPrivateEndpointInProgress, "Waiting serverless private endpoint to be configured")
	default:
		result = workflow.OK()
	}

	switch len(deploymentSpec.PrivateEndpoints) {
	case 0:
		service.UnsetCondition(api.ServerlessPrivateEndpointReadyType)
	default:
		service.SetConditionFromResult(api.ServerlessPrivateEndpointReadyType, result)
	}

	return result
}

func syncServerlessPrivateEndpoints(service *workflow.Context, projectID string, deployment *akov2.ServerlessSpec) (bool, error) {
	service.Log.Debugf("Syncing serverless private endpoints for deployment %s", deployment.Name)

	atlasPrivateEndpoints, err := listServerlessPrivateEndpoints(service, projectID, deployment.Name)
	// This is a shimmed flex cluster, if there are private serverless endpoints configured, we have to return an error.
	if admin.IsErrorCode(err, "NOT_SERVERLESS_TENANT_CLUSTER") {
		if len(deployment.PrivateEndpoints) > 0 {
			return false, fmt.Errorf("serverless private endpoints are not supported: %w", err)
		}
		return true, nil
	}

	if err != nil {
		return false, fmt.Errorf("unable to retrieve list of serverless private endpoints from Atlas: %w", err)
	}

	toCreate, toUpdate, toDelete := sortTasks(deployment.PrivateEndpoints, atlasPrivateEndpoints)
	speStatusMap := newSPEStatusMap(atlasPrivateEndpoints)

	service.Log.Debugf("Creating %d serverless private endpoints for deployment %s", len(toCreate), deployment.Name)
	for i := range toCreate {
		privateEndpointCreate := toCreate[i]
		service.Log.Debugf("Creating serverless private endpoint %s", privateEndpointCreate.Name)
		atlasPrivateEndpoint, err := createServerLessPrivateEndpoint(service, projectID, deployment.Name, &privateEndpointCreate)
		if err != nil {
			return false, fmt.Errorf("unable to create serverless private endpoint on Atlas: %w", err)
		}

		speStatusMap[privateEndpointCreate.Name] = speStatusFromAtlas(atlasPrivateEndpoint)
	}

	service.Log.Debugf("Connecting %d serverless private endpoints for deployment %s", len(toUpdate), deployment.Name)
	for i := range toUpdate {
		privateEndpointUpdate := toUpdate[i]
		service.Log.Debugf("Connecting serverless private endpoint %s", privateEndpointUpdate.Name)
		latestPEStatus := speStatusMap[privateEndpointUpdate.Name]
		atlasPrivateEndpoint, err := updateServerLessPrivateEndpoint(service, projectID, deployment.Name, latestPEStatus.ID, latestPEStatus.ProviderName, &privateEndpointUpdate)
		if err != nil {
			return false, fmt.Errorf("unable to update/connect serverless private endpoint on Atlas: %w", err)
		}

		speStatusMap[privateEndpointUpdate.Name] = speStatusFromAtlas(atlasPrivateEndpoint)
	}

	service.Log.Debugf("Deleting %d serverless private endpoints for deployment %s", len(toDelete), deployment.Name)
	for _, privateEndpointDelete := range toDelete {
		service.Log.Debugf("Deleting serverless private endpoint with ID %s", privateEndpointDelete)
		err = deleteServerLessPrivateEndpoint(service, projectID, deployment.Name, privateEndpointDelete.GetId())
		if err != nil {
			return false, fmt.Errorf("unable to delete serverless private endpoint on Atlas: %w", err)
		}

		// Serverless private endpoints first go through DELETING state before they are gone
		speStatus := speStatusMap[privateEndpointDelete.GetComment()]
		speStatus.Status = SPEStatusDeleting
		speStatusMap[privateEndpointDelete.GetComment()] = speStatus
	}

	speStatuses := make([]status.ServerlessPrivateEndpoint, 0, len(speStatusMap))
	for _, speStatus := range speStatusMap {
		speStatuses = append(speStatuses, speStatus)
	}
	service.EnsureStatusOption(status.AtlasDeploymentSPEOption(speStatuses))

	return areSPEsAvailable(speStatuses), nil
}

func isGCPWithPrivateEndpoints(deployment *akov2.ServerlessSpec) bool {
	if provider.ProviderName(deployment.ProviderSettings.BackingProviderName) == provider.ProviderGCP &&
		len(deployment.PrivateEndpoints) > 0 {
		return true
	}

	return false
}

func isGCPWithoutPrivateEndpoints(deployment *akov2.ServerlessSpec) bool {
	if provider.ProviderName(deployment.ProviderSettings.BackingProviderName) == provider.ProviderGCP &&
		len(deployment.PrivateEndpoints) == 0 {
		return true
	}

	return false
}

func listServerlessPrivateEndpoints(service *workflow.Context, projectID, deploymentName string) ([]admin.ServerlessTenantEndpoint, error) {
	// this endpoint does not offer paginated responses
	privateEndpoints, _, err := service.SdkClientSet.SdkClient20250312006.ServerlessPrivateEndpointsApi.
		ListServerlessPrivateEndpoints(service.Context, projectID, deploymentName).
		Execute()

	return privateEndpoints, err
}

func createServerLessPrivateEndpoint(service *workflow.Context, projectID, deploymentName string, privateEndpoint *akov2.ServerlessPrivateEndpoint) (*admin.ServerlessTenantEndpoint, error) {
	request := admin.ServerlessTenantCreateRequest{
		Comment: &privateEndpoint.Name,
	}

	atlasPrivateEndpoint, _, err := service.SdkClientSet.SdkClient20250312006.ServerlessPrivateEndpointsApi.
		CreateServerlessPrivateEndpoint(service.Context, projectID, deploymentName, &request).
		Execute()

	return atlasPrivateEndpoint, err
}

func updateServerLessPrivateEndpoint(service *workflow.Context, projectID, deploymentName, endpointID, providerName string, privateEndpoint *akov2.ServerlessPrivateEndpoint) (*admin.ServerlessTenantEndpoint, error) {
	// we don't allow update name (comment) once it's used as identifier
	request := admin.ServerlessTenantEndpointUpdate{
		ProviderName:            providerName,
		CloudProviderEndpointId: &privateEndpoint.CloudProviderEndpointID,
	}

	// when provider is Azure we expect IP Address to be set
	if privateEndpoint.PrivateEndpointIPAddress != "" {
		request.PrivateEndpointIpAddress = &privateEndpoint.PrivateEndpointIPAddress
	}

	atlasPrivateEndpoint, _, err := service.SdkClientSet.SdkClient20250312006.ServerlessPrivateEndpointsApi.
		UpdateServerlessPrivateEndpoint(service.Context, projectID, deploymentName, endpointID, &request).
		Execute()

	return atlasPrivateEndpoint, err
}

func deleteServerLessPrivateEndpoint(service *workflow.Context, projectID, deploymentName, endpointID string) error {
	_, err := service.SdkClientSet.SdkClient20250312006.ServerlessPrivateEndpointsApi.
		DeleteServerlessPrivateEndpoint(service.Context, projectID, deploymentName, endpointID).
		Execute()

	return err
}

// sortTasks Build and return all the operations pending to reconcile
// There are 3 possible operations:
// CREATE: A private endpoint with a given name doesn't exist on Atlas
// UPDATE: A private endpoint with a given name exists on Atlas with status RESERVED (waiting to be connected)
// DELETE: A private endpoint with a given name exists on Atlas, but it's not describe in Kubernetes resource
//
// A private endpoint is not expected to have duplicated name as validation happens at very beginning of the
// reconciliation. See validate.serverlessPrivateEndpoints
func sortTasks(
	privateEndpoints []akov2.ServerlessPrivateEndpoint,
	atlasPrivateEndpoints []admin.ServerlessTenantEndpoint,
) ([]akov2.ServerlessPrivateEndpoint, []akov2.ServerlessPrivateEndpoint, []admin.ServerlessTenantEndpoint) {
	toCreate := make([]akov2.ServerlessPrivateEndpoint, 0, len(privateEndpoints))
	toUpdate := make([]akov2.ServerlessPrivateEndpoint, 0, len(privateEndpoints))
	toDelete := make([]admin.ServerlessTenantEndpoint, 0, len(atlasPrivateEndpoints))

	privateEndpointsByName := map[string]*akov2.ServerlessPrivateEndpoint{}
	for i := range privateEndpoints {
		privateEndpoint := privateEndpoints[i]
		privateEndpointsByName[privateEndpoint.Name] = &privateEndpoint
	}

	atlasPrivateEndpointsByComment := map[string]*admin.ServerlessTenantEndpoint{}
	for i := range atlasPrivateEndpoints {
		atlasPrivateEndpoint := atlasPrivateEndpoints[i]
		atlasPrivateEndpointsByComment[atlasPrivateEndpoint.GetComment()] = &atlasPrivateEndpoint
	}

	// Collect all endpoints to create and update (connect)
	for i := range privateEndpoints {
		privateEndpoint := privateEndpoints[i]
		atlasPrivateEndpoint, ok := atlasPrivateEndpointsByComment[privateEndpoint.Name]

		// If a private endpoint with a given name doesn't exist on Atlas, add to creation list
		if !ok {
			toCreate = append(toCreate, privateEndpoint)
		}

		// If a private endpoint with a given name exists on Atlas with status RESERVED, add to update list (need to be connected)
		if isReadyToConnect(&privateEndpoint, atlasPrivateEndpoint) {
			toUpdate = append(toUpdate, privateEndpoint)
		}
	}

	for _, atlasPrivateEndpoint := range atlasPrivateEndpoints {
		// If an existing Atlas private endpoint is not present in kubernetes resource, add to deletion list
		if _, ok := privateEndpointsByName[atlasPrivateEndpoint.GetComment()]; !ok {
			toDelete = append(toDelete, atlasPrivateEndpoint)
		}
	}

	return toCreate, toUpdate, toDelete
}

func newSPEStatusMap(atlasPrivateEndpoints []admin.ServerlessTenantEndpoint) map[string]status.ServerlessPrivateEndpoint {
	statuses := map[string]status.ServerlessPrivateEndpoint{}

	for i := range atlasPrivateEndpoints {
		atlasPrivateEndpoint := atlasPrivateEndpoints[i]
		statuses[atlasPrivateEndpoint.GetComment()] = speStatusFromAtlas(&atlasPrivateEndpoint)
	}

	return statuses
}

func isReadyToConnect(privateEndpoint *akov2.ServerlessPrivateEndpoint, atlasPrivateEndpoint *admin.ServerlessTenantEndpoint) bool {
	if atlasPrivateEndpoint.GetStatus() != SPEStatusReserved {
		return false
	}

	switch provider.ProviderName(atlasPrivateEndpoint.GetProviderName()) {
	case provider.ProviderAWS:
		return privateEndpoint.CloudProviderEndpointID != ""
	case provider.ProviderAzure:
		return privateEndpoint.CloudProviderEndpointID != "" && privateEndpoint.PrivateEndpointIPAddress != ""
	}

	return false
}

func areSPEsAvailable(pe []status.ServerlessPrivateEndpoint) bool {
	for _, p := range pe {
		if p.Status != SPEStatusAvailable {
			return false
		}
	}

	return true
}

func speStatusFromAtlas(in *admin.ServerlessTenantEndpoint) status.ServerlessPrivateEndpoint {
	return status.ServerlessPrivateEndpoint{
		ID: in.GetId(),
		// Comment property is internally used as name to identify and match items on the operator against their peers on Atlas
		Name:                         in.GetComment(),
		ProviderName:                 in.GetProviderName(),
		CloudProviderEndpointID:      in.GetCloudProviderEndpointId(),
		PrivateEndpointIPAddress:     in.GetPrivateEndpointIpAddress(),
		EndpointServiceName:          in.GetEndpointServiceName(),
		PrivateLinkServiceResourceID: in.GetPrivateLinkServiceResourceId(),
		Status:                       in.GetStatus(),
		ErrorMessage:                 in.GetErrorMessage(),
	}
}
