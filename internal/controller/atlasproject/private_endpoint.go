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

package atlasproject

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312013/admin"
	"golang.org/x/exp/slices"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/httputil"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/set"
)

func ensurePrivateEndpoint(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.DeprecatedResult {
	specPEs := project.Spec.DeepCopy().PrivateEndpoints

	lastAppliedPEs, err := mapLastAppliedPrivateEndpoint(project)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	atlasPEs, err := getAllPrivateEndpoints(workflowCtx.Context, workflowCtx.SdkClientSet.SdkClient20250312012, project.ID())
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	result, conditionType := syncPrivateEndpointsWithAtlas(workflowCtx, project.ID(), specPEs, atlasPEs, lastAppliedPEs)
	if !result.IsOk() {
		if conditionType == api.PrivateEndpointServiceReadyType {
			workflowCtx.UnsetCondition(api.PrivateEndpointReadyType)
		}
		workflowCtx.SetConditionFromResult(conditionType, result)
		return result
	}

	if (len(specPEs) == 0 && len(atlasPEs) == 0) || !hasManagedPrivateEndpoints(specPEs, atlasPEs, lastAppliedPEs) {
		workflowCtx.UnsetCondition(api.PrivateEndpointServiceReadyType)
		workflowCtx.UnsetCondition(api.PrivateEndpointReadyType)
		return workflow.OK()
	}

	serviceStatus := getStatusForServices(atlasPEs)
	if !serviceStatus.IsOk() {
		workflowCtx.SetConditionFromResult(api.PrivateEndpointServiceReadyType, serviceStatus)
		return serviceStatus
	}

	unconfiguredAmount := countNotConfiguredEndpoints(specPEs)
	if unconfiguredAmount != 0 {
		serviceStatus = serviceStatus.WithMessage("Interface Private Endpoint awaits configuration")
		workflowCtx.SetConditionFromResult(api.PrivateEndpointServiceReadyType, serviceStatus)

		if len(specPEs) == unconfiguredAmount {
			workflowCtx.UnsetCondition(api.PrivateEndpointReadyType)
			return serviceStatus
		} else {
			return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, errors.New("not all interface private endpoints are fully configured"))
		}
	}

	interfaceStatus := getStatusForInterfaces(workflowCtx, project.ID(), specPEs, atlasPEs)
	workflowCtx.SetConditionFromResult(api.PrivateEndpointReadyType, interfaceStatus)

	return interfaceStatus
}

func syncPrivateEndpointsWithAtlas(
	ctx *workflow.Context,
	projectID string,
	specPEs []akov2.PrivateEndpoint,
	atlasPEs []atlasPE,
	lastAppliedPEs map[string]akov2.PrivateEndpoint,
) (workflow.DeprecatedResult, api.ConditionType) {
	log := ctx.Log

	log.Debugw("PE Connections", "atlasPEs", atlasPEs, "specPEs", specPEs)
	endpointsToDelete := getEndpointsNotInSpec(specPEs, atlasPEs, lastAppliedPEs)
	log.Debugf("Number of Private Endpoints to delete: %d", len(endpointsToDelete))
	if result := deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete); !result.IsOk() {
		return result, api.PrivateEndpointServiceReadyType
	}

	endpointsToCreate, endpointCounts := getEndpointsNotInAtlas(specPEs, atlasPEs)
	log.Debugf("Number of Private Endpoints to create: %d", len(endpointsToCreate))
	newConnections, err := createPeServiceInAtlas(ctx, projectID, endpointsToCreate, endpointCounts)
	if err != nil {
		return terminateWithError(ctx, api.PrivateEndpointServiceReadyType, "Failed to create PE Service in Atlas", err)
	}

	endpointsToSync := getEndpointsIntersection(specPEs, atlasPEs)
	log.Debugf("Number of Private Endpoints to sync: %d", len(endpointsToSync))
	syncedConnections, err := syncPeInterfaceInAtlas(ctx, projectID, endpointsToSync)
	if err != nil {
		return terminateWithError(ctx, api.PrivateEndpointReadyType, "Failed to sync PE Interface in Atlas", err)
	}

	log.Debugw("PE Changes", "newConnections", newConnections, "syncedConnections", syncedConnections)
	updatePEStatusOption(ctx, projectID, newConnections, syncedConnections)

	if len(newConnections) != 0 {
		return notReadyServiceResult, api.PrivateEndpointServiceReadyType
	}

	return workflow.OK(), api.PrivateEndpointReadyType
}

func getStatusForServices(atlasPEs []atlasPE) workflow.DeprecatedResult {
	allAvailable := true
	for _, conn := range atlasPEs {
		if isFailed(conn.GetStatus()) {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, errors.New(conn.GetErrorMessage()))
		}

		if !isAvailable(conn.GetStatus()) {
			allAvailable = false
		}
	}

	if !allAvailable {
		return notReadyServiceResult
	}

	return workflow.OK()
}

func getStatusForInterfaces(ctx *workflow.Context, projectID string, specPEs []akov2.PrivateEndpoint, atlasPEs []atlasPE) workflow.DeprecatedResult {
	totalInterfaceCount := 0

	for _, atlasPeService := range atlasPEs {
		interfaceEndpointIDs := atlasPeService.InterfaceEndpointIDs()
		totalInterfaceCount += len(interfaceEndpointIDs)

		for _, interfaceEndpointID := range interfaceEndpointIDs {
			if interfaceEndpointID == "" {
				return notReadyInterfaceResult
			}

			interfaceEndpoint, _, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.GetPrivateEndpointWithParams(ctx.Context, &admin.GetPrivateEndpointApiParams{
				GroupId:           projectID,
				CloudProvider:     atlasPeService.GetCloudProvider(),
				EndpointId:        interfaceEndpointID,
				EndpointServiceId: atlasPeService.GetId(),
			}).Execute()
			if err != nil {
				return workflow.Terminate(workflow.Internal, err)
			}

			interfaceIsAvailable, interfaceFailureMessage := checkIfInterfaceIsAvailable(interfaceEndpoint)
			if interfaceFailureMessage != "" {
				return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, errors.New(interfaceFailureMessage))
			}
			if !interfaceIsAvailable {
				return notReadyInterfaceResult
			}
		}
	}

	if len(specPEs) != totalInterfaceCount {
		return notReadyInterfaceResult
	}

	return workflow.OK()
}

func updatePEStatusOption(ctx *workflow.Context, projectID string, newConnections, syncedConnections []atlasPE) {
	setPEStatusOption(ctx, projectID, syncedConnections)
	addPEStatusOption(ctx, projectID, newConnections)
}

func addPEStatusOption(ctx *workflow.Context, projectID string, newPEs []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, newPEs)
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEndpointsOption(statusPEs))
}

func setPEStatusOption(ctx *workflow.Context, projectID string, atlasPeConnections []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, atlasPeConnections)
	ctx.EnsureStatusOption(status.AtlasProjectSetPrivateEndpointsOption(statusPEs))
}

type atlasPE struct {
	admin.EndpointService
}

func (a atlasPE) Identifier() interface{} {
	return a.CloudProvider + status.TransformRegionToID(a.GetRegionName())
}

func (a atlasPE) InterfaceEndpointIDs() []string {
	if len(a.GetInterfaceEndpoints()) != 0 {
		return a.GetInterfaceEndpoints()
	}

	if len(a.GetPrivateEndpoints()) != 0 {
		return a.GetPrivateEndpoints()
	}

	if len(a.GetEndpointGroupNames()) != 0 {
		return a.GetEndpointGroupNames()
	}

	return nil
}

func getAllPrivateEndpoints(ctx context.Context, client *admin.APIClient, projectID string) (result []atlasPE, err error) {
	providers := []string{"AWS", "AZURE", "GCP"}
	for _, p := range providers {
		// this endpoint does not offer paginated resources.
		atlasPeConnections, _, err := client.PrivateEndpointServicesApi.ListPrivateEndpointService(ctx, projectID, p).Execute()
		if err != nil {
			return nil, err
		}

		for connIdx := range atlasPeConnections {
			atlasPeConnections[connIdx].CloudProvider = p
		}

		for _, atlasPeConnection := range atlasPeConnections {
			result = append(result, atlasPE{atlasPeConnection})
		}
	}

	return
}

func createPeServiceInAtlas(ctx *workflow.Context, projectID string, endpointsToCreate []akov2.PrivateEndpoint, endpointCounts []int) (newConnections []atlasPE, err error) {
	newConnections = make([]atlasPE, 0)
	for idx, pe := range endpointsToCreate {
		conn, _, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.CreatePrivateEndpointService(ctx.Context, projectID, &admin.CloudProviderEndpointServiceRequest{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		}).Execute()
		if err != nil {
			return newConnections, err
		}

		conn.SetCloudProvider(string(pe.Provider))
		conn.SetRegionName(string(pe.Provider))
		newConn := atlasPE{*conn}

		for i := 0; i < endpointCounts[idx]; i++ {
			newConnections = append(newConnections, newConn)
		}
	}

	return newConnections, nil
}

func syncPeInterfaceInAtlas(ctx *workflow.Context, projectID string, endpointsToUpdate []intersectionPair) (syncedEndpoints []atlasPE, err error) {
	syncedEndpoints = make([]atlasPE, 0)
	for _, pair := range endpointsToUpdate {
		specPeService := pair.spec
		atlasPeService := pair.atlas

		ctx.Log.Debugw("endpointNeedsUpdating", "specPeService", specPeService, "atlasPeService", atlasPeService, "endpointNeedsUpdating", endpointNeedsUpdating(specPeService, atlasPeService))
		if endpointNeedsUpdating(specPeService, atlasPeService) {
			interfaceConn := &admin.CreateEndpointRequest{
				Id:                       pointer.SetOrNil(specPeService.ID, ""),
				PrivateEndpointIPAddress: pointer.SetOrNil(specPeService.IP, ""),
				EndpointGroupName:        pointer.SetOrNil(specPeService.EndpointGroupName, ""),
				GcpProjectId:             pointer.SetOrNil(specPeService.GCPProjectID, ""),
			}
			interfaceConn.Endpoints = specPeService.Endpoints.ConvertToAtlas()

			privateEndpoint, response, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.CreatePrivateEndpointWithParams(ctx.Context, &admin.CreatePrivateEndpointApiParams{
				GroupId:               projectID,
				CloudProvider:         string(specPeService.Provider),
				EndpointServiceId:     atlasPeService.GetId(),
				CreateEndpointRequest: interfaceConn,
			}).Execute()

			ctx.Log.Debugw("CreatePrivateEndpoint Reply", "privateEndpoint", privateEndpoint, "err", err)
			if err != nil {
				ctx.Log.Debugw("failed to create PE Interface", "error", err)
				statusCode := httputil.StatusCode(response)
				if statusCode == http.StatusBadRequest || statusCode == http.StatusConflict {
					return syncedEndpoints, err
				}
			}
		}

		atlasPeService.SetCloudProvider(string(specPeService.Provider))
		atlasPeService.SetRegionName(specPeService.Region)
		syncedEndpoints = append(syncedEndpoints, atlasPeService)
	}

	return
}

func endpointNeedsUpdating(specPeService akov2.PrivateEndpoint, atlasPeService atlasPE) bool {
	if isAvailable(atlasPeService.GetStatus()) && endpointDefinedInSpec(specPeService) {
		switch specPeService.Provider {
		case provider.ProviderAWS, provider.ProviderAzure:
			return !slices.Contains(atlasPeService.InterfaceEndpointIDs(), specPeService.ID)
		case provider.ProviderGCP:
			return !slices.Contains(atlasPeService.InterfaceEndpointIDs(), specPeService.EndpointGroupName) || len(atlasPeService.GetServiceAttachmentNames()) != len(specPeService.Endpoints)
		}
	}

	return false
}

func countNotConfiguredEndpoints(endpoints []akov2.PrivateEndpoint) (count int) {
	for _, pe := range endpoints {
		if !endpointDefinedInSpec(pe) {
			count++
		}
	}

	return count
}

func endpointDefinedInSpec(specEndpoint akov2.PrivateEndpoint) bool {
	return specEndpoint.ID != "" || specEndpoint.EndpointGroupName != ""
}

func DeleteAllPrivateEndpoints(ctx *workflow.Context, atlasProject *akov2.AtlasProject) workflow.DeprecatedResult {
	atlasPEs, err := getAllPrivateEndpoints(ctx.Context, ctx.SdkClientSet.SdkClient20250312012, atlasProject.ID())
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	lastAppliedSpecPEs, err := mapLastAppliedPrivateEndpoint(atlasProject)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	endpointsToDelete := getEndpointsNotInSpec([]akov2.PrivateEndpoint{}, atlasPEs, lastAppliedSpecPEs)
	return deletePrivateEndpointsFromAtlas(ctx, atlasProject.ID(), endpointsToDelete)
}

func deletePrivateEndpointsFromAtlas(ctx *workflow.Context, projectID string, listsToRemove []atlasPE) workflow.DeprecatedResult {
	if len(listsToRemove) == 0 {
		return workflow.OK()
	}

	for _, peService := range listsToRemove {
		if isDeleting(peService.GetStatus()) {
			ctx.Log.Debugf("%s Private Endpoint Service for the region %s is being deleted", peService.GetCloudProvider(), peService.GetRegionName())
			continue
		}

		interfaceEndpointIDs := peService.InterfaceEndpointIDs()
		if len(interfaceEndpointIDs) != 0 {
			for _, interfaceEndpointID := range interfaceEndpointIDs {
				_, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.DeletePrivateEndpointWithParams(ctx.Context, &admin.DeletePrivateEndpointApiParams{
					GroupId:           projectID,
					CloudProvider:     peService.GetCloudProvider(),
					EndpointId:        interfaceEndpointID,
					EndpointServiceId: peService.GetId(),
				}).Execute()
				if err != nil {
					return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, errors.New("failed to delete Private Endpoint"))
				}
			}

			continue
		}

		_, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.DeletePrivateEndpointServiceWithParams(ctx.Context, &admin.DeletePrivateEndpointServiceApiParams{
			GroupId:           projectID,
			CloudProvider:     peService.GetCloudProvider(),
			EndpointServiceId: peService.GetId(),
		}).Execute()
		if err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, errors.New("failed to delete Private Endpoint Service"))
		}

		ctx.Log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "provider", peService.GetCloudProvider(), "regionName", peService.RegionName)
	}

	return workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint is deleting")
}

func convertAllToStatus(ctx *workflow.Context, projectID string, peList []atlasPE) (result []status.ProjectPrivateEndpoint) {
	for _, endpoint := range peList {
		result = append(result, convertOneServiceToStatus(ctx, projectID, endpoint)...)
	}

	return result
}

func convertOneServiceToStatus(ctx *workflow.Context, projectID string, conn atlasPE) []status.ProjectPrivateEndpoint {
	interfaceEndpointIDs := conn.InterfaceEndpointIDs()

	if len(interfaceEndpointIDs) == 0 {
		return []status.ProjectPrivateEndpoint{
			filledPEStatus(ctx, projectID, conn, ""),
		}
	}

	result := make([]status.ProjectPrivateEndpoint, 0)
	for _, interfaceEndpointID := range interfaceEndpointIDs {
		result = append(result, filledPEStatus(ctx, projectID, conn, interfaceEndpointID))
	}

	return result
}

func filledPEStatus(ctx *workflow.Context, projectID string, conn atlasPE, InterfaceEndpointID string) status.ProjectPrivateEndpoint {
	pe := status.ProjectPrivateEndpoint{
		ID:                  conn.GetId(),
		Provider:            provider.ProviderName(conn.GetCloudProvider()),
		Region:              conn.GetRegionName(),
		InterfaceEndpointID: InterfaceEndpointID,
	}

	switch pe.Provider {
	case provider.ProviderAWS:
		pe.ServiceName = conn.GetEndpointServiceName()
		pe.ServiceResourceID = conn.GetId()
	case provider.ProviderAzure:
		pe.ServiceName = conn.GetPrivateLinkServiceName()
		pe.ServiceResourceID = conn.GetPrivateLinkServiceResourceId()
	case provider.ProviderGCP:
		pe.ServiceAttachmentNames = conn.GetServiceAttachmentNames()
		if InterfaceEndpointID != "" {
			var err error
			pe.Endpoints, err = getGCPInterfaceEndpoint(ctx, projectID, pe)
			if err != nil {
				ctx.Log.Warnw("failed to get Interface Endpoint Data for GCP", "err", err, "pe", pe)
			}
		}
	}

	ctx.Log.Debugw("Converted One Status", "connection", conn, "private endpoint", pe)

	return pe
}

// getGCPInterfaceEndpoint returns an InterfaceEndpointID and a list of GCP endpoints
func getGCPInterfaceEndpoint(ctx *workflow.Context, projectID string, endpoint status.ProjectPrivateEndpoint) ([]status.GCPEndpoint, error) {
	log := ctx.Log
	if endpoint.InterfaceEndpointID == "" {
		return nil, errors.New("InterfaceEndpointID is empty")
	}
	interfaceEndpointConn, _, err := ctx.SdkClientSet.SdkClient20250312012.PrivateEndpointServicesApi.GetPrivateEndpointWithParams(ctx.Context, &admin.GetPrivateEndpointApiParams{
		GroupId:           projectID,
		CloudProvider:     string(provider.ProviderGCP),
		EndpointId:        endpoint.InterfaceEndpointID,
		EndpointServiceId: endpoint.ID,
	}).Execute()
	if err != nil {
		return nil, err
	}

	interfaceConns := interfaceEndpointConn.GetEndpoints()
	listOfInterfaces := make([]status.GCPEndpoint, 0)
	for _, e := range interfaceConns {
		endpoint := status.GCPEndpoint{
			Status:       e.GetStatus(),
			EndpointName: e.GetEndpointName(),
			IPAddress:    e.GetIpAddress(),
		}
		listOfInterfaces = append(listOfInterfaces, endpoint)
	}
	log.Debugw("DeprecatedResult of getGCPEndpointData", "endpoint.ID", endpoint.ID, "listOfInterfaces", listOfInterfaces)

	return listOfInterfaces, nil
}

// checkIfInterfaceIsAvailable checks if an interface and all of its nested endpoints are available and also returns an error message
func checkIfInterfaceIsAvailable(interfaceEndpointConn *admin.PrivateLinkEndpoint) (allAvailable bool, failureMessage string) {
	allAvailable = true

	if isFailed(interfaceEndpointConn.GetStatus()) {
		return false, interfaceEndpointConn.GetErrorMessage()
	}
	if !isAvailable(interfaceEndpointConn.GetStatus()) && !isAvailable(interfaceEndpointConn.GetConnectionStatus()) {
		allAvailable = false
	}

	for _, endpoint := range interfaceEndpointConn.GetEndpoints() {
		if isFailed(endpoint.GetStatus()) {
			return false, interfaceEndpointConn.GetErrorMessage()
		}
		if !isAvailable(endpoint.GetStatus()) {
			allAvailable = false
		}
	}

	return
}

func isAvailable(status string) bool {
	return status == "AVAILABLE"
}

func isDeleting(status string) bool {
	return status == "DELETING"
}

func isFailed(status string) bool {
	return status == "FAILED"
}

func terminateWithError(ctx *workflow.Context, conditionType api.ConditionType, message string, err error) (workflow.DeprecatedResult, api.ConditionType) {
	ctx.Log.Debugw(message, "error", err)
	result := workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, err).WithoutRetry()
	return result, conditionType
}

var notReadyServiceResult = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is not ready")
var notReadyInterfaceResult = workflow.InProgress(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "Interface Private Endpoint is not ready")

func getEndpointsNotInSpec(specPEs []akov2.PrivateEndpoint, atlasPEs []atlasPE, lastAppliedPEs map[string]akov2.PrivateEndpoint) []atlasPE {
	notInSpecItems, _ := getUniqueDifference(atlasPEs, specPEs)
	toDelete := make([]atlasPE, 0, len(notInSpecItems))
	for _, item := range notInSpecItems {
		if _, ok := lastAppliedPEs[item.Identifier().(string)]; ok {
			toDelete = append(toDelete, item)
		}
	}

	return toDelete
}

func getEndpointsNotInAtlas(specPEs []akov2.PrivateEndpoint, atlasPEs []atlasPE) (toCreate []akov2.PrivateEndpoint, counts []int) {
	return getUniqueDifference(specPEs, atlasPEs)
}

func getUniqueDifference[ResultType interface{}, OtherType interface{}](left []ResultType, right []OtherType) (uniques []ResultType, counts []int) {
	difference := set.DeprecatedDifference(left, right)

	uniqueItems := make(map[string]itemCount)
	for _, item := range difference {
		key := item.Identifier().(string)
		if uniqueItem, found := uniqueItems[key]; found {
			uniqueItem.Count += 1
			uniqueItems[key] = uniqueItem
		} else {
			uniqueItems[key] = itemCount{
				Item:  item,
				Count: 1,
			}
		}
	}

	for _, value := range uniqueItems {
		uniques = append(uniques, value.Item.(ResultType))
		counts = append(counts, value.Count)
	}

	return
}

type itemCount struct {
	Item  interface{}
	Count int
}

func getEndpointsIntersection(specPEs []akov2.PrivateEndpoint, atlasPEs []atlasPE) []intersectionPair {
	intersection := set.DeprecatedIntersection(specPEs, atlasPEs)
	result := []intersectionPair{}
	for _, item := range intersection {
		pair := intersectionPair{}
		pair.spec = item[0].(akov2.PrivateEndpoint)
		pair.atlas = item[1].(atlasPE)
		result = append(result, pair)
	}
	return result
}

type intersectionPair struct {
	spec  akov2.PrivateEndpoint
	atlas atlasPE
}

func mapLastAppliedPrivateEndpoint(atlasProject *akov2.AtlasProject) (map[string]akov2.PrivateEndpoint, error) {
	lastApplied, err := lastAppliedSpecFrom(atlasProject)
	if err != nil {
		return nil, err
	}

	if lastApplied == nil || len(lastApplied.PrivateEndpoints) == 0 {
		return nil, nil
	}

	result := map[string]akov2.PrivateEndpoint{}
	for _, pe := range lastApplied.PrivateEndpoints {
		result[pe.Identifier().(string)] = pe
	}

	return result, nil
}

func hasManagedPrivateEndpoints(specPEs []akov2.PrivateEndpoint, atlasPEs []atlasPE, lastAppliedPEs map[string]akov2.PrivateEndpoint) bool {
	for _, pe := range atlasPEs {
		// if any of the PE in atlas was previously managed, return true
		if _, ok := lastAppliedPEs[pe.Identifier().(string)]; ok {
			return true
		}
	}

	// if any of the PE in atlas is specified in the spec, return true
	a := set.DeprecatedIntersection(specPEs, atlasPEs)
	if len(specPEs) > 0 && len(a) > 0 {
		return true
	}

	// if there are not managed or previously managed entries, return false
	return false
}
