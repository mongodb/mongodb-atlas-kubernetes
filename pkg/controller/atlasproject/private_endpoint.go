package atlasproject

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/customresource"

	"golang.org/x/exp/slices"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/set"
)

func ensurePrivateEndpoint(workflowCtx *workflow.Context, project *mdbv1.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canPrivateEndpointReconcile(workflowCtx.Client, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, "the operator could not validate ownership for deletion protection")
		workflowCtx.SetConditionFromResult(status.PrivateEndpointReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Private Endpoint(s) due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.PrivateEndpointReadyType, result)

		return result
	}

	specPEs := project.Spec.DeepCopy().PrivateEndpoints

	atlasPEs, err := getAllPrivateEndpoints(workflowCtx.Client, project.ID())
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	result, conditionType := syncPrivateEndpointsWithAtlas(workflowCtx, project.ID(), specPEs, atlasPEs)
	if !result.IsOk() {
		if conditionType == status.PrivateEndpointServiceReadyType {
			workflowCtx.UnsetCondition(status.PrivateEndpointReadyType)
		}
		workflowCtx.SetConditionFromResult(conditionType, result)
		return result
	}

	if len(specPEs) == 0 && len(atlasPEs) == 0 {
		workflowCtx.UnsetCondition(status.PrivateEndpointServiceReadyType)
		workflowCtx.UnsetCondition(status.PrivateEndpointReadyType)
		return workflow.OK()
	}

	serviceStatus := getStatusForServices(workflowCtx, atlasPEs)
	if !serviceStatus.IsOk() {
		workflowCtx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, serviceStatus)
		return serviceStatus
	}

	unconfiguredAmount := countNotConfiguredEndpoints(specPEs)
	if unconfiguredAmount != 0 {
		serviceStatus = serviceStatus.WithMessage("Interface Private Endpoint awaits configuration")
		workflowCtx.SetConditionFromResult(status.PrivateEndpointServiceReadyType, serviceStatus)

		if len(specPEs) == unconfiguredAmount {
			workflowCtx.UnsetCondition(status.PrivateEndpointReadyType)
			return serviceStatus
		} else {
			return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "Not All Interface Private Endpoint are fully configured")
		}
	}

	interfaceStatus := getStatusForInterfaces(workflowCtx, project.ID(), specPEs, atlasPEs)
	workflowCtx.SetConditionFromResult(status.PrivateEndpointReadyType, interfaceStatus)

	return interfaceStatus
}

func syncPrivateEndpointsWithAtlas(ctx *workflow.Context, projectID string, specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) (workflow.Result, status.ConditionType) {
	log := ctx.Log

	log.Debugw("PE Connections", "atlasPEs", atlasPEs, "specPEs", specPEs)
	endpointsToDelete := getEndpointsNotInSpec(specPEs, atlasPEs)
	log.Debugf("Number of Private Endpoints to delete: %d", len(endpointsToDelete))
	if result := deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete); !result.IsOk() {
		return result, status.PrivateEndpointServiceReadyType
	}

	endpointsToCreate, endpointCounts := getEndpointsNotInAtlas(specPEs, atlasPEs)
	log.Debugf("Number of Private Endpoints to create: %d", len(endpointsToCreate))
	newConnections, err := createPeServiceInAtlas(ctx, projectID, endpointsToCreate, endpointCounts)
	if err != nil {
		return terminateWithError(ctx, status.PrivateEndpointServiceReadyType, "Failed to create PE Service in Atlas", err)
	}

	endpointsToSync := getEndpointsIntersection(specPEs, atlasPEs)
	log.Debugf("Number of Private Endpoints to sync: %d", len(endpointsToSync))
	syncedConnections, err := syncPeInterfaceInAtlas(ctx, projectID, endpointsToSync)
	if err != nil {
		return terminateWithError(ctx, status.PrivateEndpointReadyType, "Failed to sync PE Interface in Atlas", err)
	}

	log.Debugw("PE Changes", "newConnections", newConnections, "syncedConnections", syncedConnections)
	updatePEStatusOption(ctx, projectID, newConnections, syncedConnections)

	if len(newConnections) != 0 {
		return notReadyServiceResult, status.PrivateEndpointServiceReadyType
	}

	return workflow.OK(), status.PrivateEndpointReadyType
}

func getStatusForServices(ctx *workflow.Context, atlasPEs []atlasPE) workflow.Result {
	allAvailable, failureMessage := areServicesAvailableOrFailed(atlasPEs)
	ctx.Log.Debugw("Get Status for Services", "allAvailable", allAvailable, "failureMessage", failureMessage)
	if failureMessage != "" {
		return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, failureMessage)
	}
	if !allAvailable {
		return notReadyServiceResult
	}

	return workflow.OK()
}

func getStatusForInterfaces(ctx *workflow.Context, projectID string, specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) workflow.Result {
	totalInterfaceCount := 0

	for _, atlasPeService := range atlasPEs {
		interfaceEndpointIDs := atlasPeService.InterfaceEndpointIDs()
		totalInterfaceCount += len(interfaceEndpointIDs)

		for _, interfaceEndpointID := range interfaceEndpointIDs {
			if interfaceEndpointID == "" {
				return notReadyInterfaceResult
			}

			interfaceEndpoint, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, atlasPeService.ProviderName, atlasPeService.ID, interfaceEndpointID)
			if err != nil {
				return workflow.Terminate(workflow.Internal, err.Error())
			}

			interfaceIsAvailable, interfaceFailureMessage := checkIfInterfaceIsAvailable(interfaceEndpoint)
			if interfaceFailureMessage != "" {
				return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, interfaceFailureMessage)
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

func areServicesAvailableOrFailed(atlasPeConnections []atlasPE) (allAvailable bool, failureMessage string) {
	allAvailable = true

	for _, conn := range atlasPeConnections {
		if isFailed(conn.Status) {
			failureMessage = conn.ErrorMessage
			return
		}
		if !isAvailable(conn.Status) {
			allAvailable = false
		}
	}

	return
}

func updatePEStatusOption(ctx *workflow.Context, projectID string, newConnections, syncedConnections []atlasPE) {
	setPEStatusOption(ctx, projectID, syncedConnections)
	addPEStatusOption(ctx, projectID, newConnections)
}

func addPEStatusOption(ctx *workflow.Context, projectID string, newPEs []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, newPEs)
	ctx.EnsureStatusOption(status.AtlasProjectAddPrivateEnpointsOption(statusPEs))
}

func setPEStatusOption(ctx *workflow.Context, projectID string, atlasPeConnections []atlasPE) {
	statusPEs := convertAllToStatus(ctx, projectID, atlasPeConnections)
	ctx.EnsureStatusOption(status.AtlasProjectSetPrivateEnpointsOption(statusPEs))
}

type atlasPE mongodbatlas.PrivateEndpointConnection

func (a atlasPE) Identifier() interface{} {
	return a.ProviderName + status.TransformRegionToID(a.RegionName)
}

func (a atlasPE) InterfaceEndpointIDs() []string {
	if len(a.InterfaceEndpoints) != 0 {
		return a.InterfaceEndpoints
	}

	if len(a.PrivateEndpoints) != 0 {
		return a.PrivateEndpoints
	}

	if len(a.EndpointGroupNames) != 0 {
		return a.EndpointGroupNames
	}

	return nil
}

func getAllPrivateEndpoints(client mongodbatlas.Client, projectID string) (result []atlasPE, err error) {
	providers := []string{"AWS", "AZURE", "GCP"}
	for _, provider := range providers {
		atlasPeConnections, _, err := client.PrivateEndpoints.List(context.Background(), projectID, provider, &mongodbatlas.ListOptions{})
		if err != nil {
			return nil, err
		}

		for connIdx := range atlasPeConnections {
			atlasPeConnections[connIdx].ProviderName = provider
		}

		for _, atlasPeConnection := range atlasPeConnections {
			result = append(result, atlasPE(atlasPeConnection))
		}
	}

	return
}

func createPeServiceInAtlas(ctx *workflow.Context, projectID string, endpointsToCreate []mdbv1.PrivateEndpoint, endpointCounts []int) (newConnections []atlasPE, err error) {
	newConnections = make([]atlasPE, 0)
	for idx, pe := range endpointsToCreate {
		conn, _, err := ctx.Client.PrivateEndpoints.Create(context.Background(), projectID, &mongodbatlas.PrivateEndpointConnection{
			ProviderName: string(pe.Provider),
			Region:       pe.Region,
		})
		if err != nil {
			return newConnections, err
		}

		conn.ProviderName = string(pe.Provider)
		conn.Region = pe.Region
		newConn := atlasPE(*conn)

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
			interfaceConn := &mongodbatlas.InterfaceEndpointConnection{
				ID:                       specPeService.ID,
				PrivateEndpointIPAddress: specPeService.IP,
				EndpointGroupName:        specPeService.EndpointGroupName,
				GCPProjectID:             specPeService.GCPProjectID,
			}
			if gcpEndpoints, err := specPeService.Endpoints.ConvertToAtlas(); err == nil {
				interfaceConn.Endpoints = gcpEndpoints
			}

			interfaceConn, response, err := ctx.Client.PrivateEndpoints.AddOnePrivateEndpoint(context.Background(), projectID, string(specPeService.Provider), atlasPeService.ID, interfaceConn)
			ctx.Log.Debugw("AddOnePrivateEndpoint Reply", "interfaceConn", interfaceConn, "err", err)
			if err != nil {
				ctx.Log.Debugw("failed to create PE Interface", "error", err)
				if response.StatusCode == http.StatusBadRequest || response.StatusCode == http.StatusConflict {
					return syncedEndpoints, err
				}
			}
		}

		atlasPeService.ProviderName = string(specPeService.Provider)
		atlasPeService.Region = specPeService.Region
		syncedEndpoints = append(syncedEndpoints, atlasPeService)
	}

	return
}

func endpointNeedsUpdating(specPeService mdbv1.PrivateEndpoint, atlasPeService atlasPE) bool {
	if isAvailable(atlasPeService.Status) && endpointDefinedInSpec(specPeService) {
		switch specPeService.Provider {
		case provider.ProviderAWS, provider.ProviderAzure:
			return !slices.Contains(atlasPeService.InterfaceEndpointIDs(), specPeService.ID)
		case provider.ProviderGCP:
			return !slices.Contains(atlasPeService.InterfaceEndpointIDs(), specPeService.EndpointGroupName) || len(atlasPeService.ServiceAttachmentNames) != len(specPeService.Endpoints)
		}
	}

	return false
}

func countNotConfiguredEndpoints(endpoints []mdbv1.PrivateEndpoint) (count int) {
	for _, pe := range endpoints {
		if !endpointDefinedInSpec(pe) {
			count++
		}
	}

	return count
}

func endpointDefinedInSpec(specEndpoint mdbv1.PrivateEndpoint) bool {
	return specEndpoint.ID != "" || specEndpoint.EndpointGroupName != ""
}

func DeleteAllPrivateEndpoints(ctx *workflow.Context, projectID string) workflow.Result {
	atlasPEs, err := getAllPrivateEndpoints(ctx.Client, projectID)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err.Error())
	}

	endpointsToDelete := getEndpointsNotInSpec([]mdbv1.PrivateEndpoint{}, atlasPEs)
	return deletePrivateEndpointsFromAtlas(ctx, projectID, endpointsToDelete)
}

func deletePrivateEndpointsFromAtlas(ctx *workflow.Context, projectID string, listsToRemove []atlasPE) workflow.Result {
	if len(listsToRemove) == 0 {
		return workflow.OK()
	}

	for _, peService := range listsToRemove {
		if isDeleting(peService.Status) {
			ctx.Log.Debugf("%s Private Endpoint Service for the region %s is being deleted", peService.ProviderName, peService.RegionName)
			continue
		}

		interfaceEndpointIDs := peService.InterfaceEndpointIDs()
		if len(interfaceEndpointIDs) != 0 {
			for _, interfaceEndpointID := range interfaceEndpointIDs {
				if _, err := ctx.Client.PrivateEndpoints.DeleteOnePrivateEndpoint(context.Background(), projectID, peService.ProviderName, peService.ID, interfaceEndpointID); err != nil {
					return workflow.Terminate(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "failed to delete Private Endpoint")
				}
			}

			continue
		}

		if _, err := ctx.Client.PrivateEndpoints.Delete(context.Background(), projectID, peService.ProviderName, peService.ID); err != nil {
			return workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, "failed to delete Private Endpoint Service")
		}

		ctx.Log.Debugw("Removed Private Endpoint Service from Atlas as it's not specified in current AtlasProject", "provider", peService.ProviderName, "regionName", peService.RegionName)
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
		ID:                  conn.ID,
		Provider:            provider.ProviderName(conn.ProviderName),
		Region:              conn.Region,
		InterfaceEndpointID: InterfaceEndpointID,
	}

	switch pe.Provider {
	case provider.ProviderAWS:
		pe.ServiceName = conn.EndpointServiceName
		pe.ServiceResourceID = conn.ID
	case provider.ProviderAzure:
		pe.ServiceName = conn.PrivateLinkServiceName
		pe.ServiceResourceID = conn.PrivateLinkServiceResourceID
	case provider.ProviderGCP:
		pe.ServiceAttachmentNames = conn.ServiceAttachmentNames
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
	interfaceEndpointConn, _, err := ctx.Client.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), projectID, string(provider.ProviderGCP), endpoint.ID, endpoint.InterfaceEndpointID)
	if err != nil {
		return nil, err
	}

	interfaceConns := interfaceEndpointConn.Endpoints
	listOfInterfaces := make([]status.GCPEndpoint, 0)
	for _, e := range interfaceConns {
		endpoint := status.GCPEndpoint{
			Status:       e.Status,
			EndpointName: e.EndpointName,
			IPAddress:    e.IPAddress,
		}
		listOfInterfaces = append(listOfInterfaces, endpoint)
	}
	log.Debugw("Result of getGCPEndpointData", "endpoint.ID", endpoint.ID, "listOfInterfaces", listOfInterfaces)

	return listOfInterfaces, nil
}

// checkIfInterfaceIsAvailable checks if an interface and all of its nested endpoints are available and also returns an error message
func checkIfInterfaceIsAvailable(interfaceEndpointConn *mongodbatlas.InterfaceEndpointConnection) (allAvailable bool, failureMessage string) {
	allAvailable = true

	if isFailed(interfaceEndpointConn.Status) {
		return false, interfaceEndpointConn.ErrorMessage
	}
	if !isAvailable(interfaceEndpointConn.Status) && !isAvailable(interfaceEndpointConn.AWSConnectionStatus) {
		allAvailable = false
	}

	for _, endpoint := range interfaceEndpointConn.Endpoints {
		if isFailed(endpoint.Status) {
			return false, interfaceEndpointConn.ErrorMessage
		}
		if !isAvailable(endpoint.Status) {
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

func terminateWithError(ctx *workflow.Context, conditionType status.ConditionType, message string, err error) (workflow.Result, status.ConditionType) {
	ctx.Log.Debugw(message, "error", err)
	result := workflow.Terminate(workflow.ProjectPEServiceIsNotReadyInAtlas, err.Error()).WithoutRetry()
	return result, conditionType
}

var notReadyServiceResult = workflow.InProgress(workflow.ProjectPEServiceIsNotReadyInAtlas, "Private Endpoint Service is not ready")
var notReadyInterfaceResult = workflow.InProgress(workflow.ProjectPEInterfaceIsNotReadyInAtlas, "Interface Private Endpoint is not ready")

func getEndpointsNotInSpec(specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) []atlasPE {
	uniqueItems, _ := getUniqueDifference(atlasPEs, specPEs)
	return uniqueItems
}

func getEndpointsNotInAtlas(specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) (toCreate []mdbv1.PrivateEndpoint, counts []int) {
	return getUniqueDifference(specPEs, atlasPEs)
}

func getUniqueDifference[ResultType interface{}, OtherType interface{}](left []ResultType, right []OtherType) (uniques []ResultType, counts []int) {
	difference := set.Difference(left, right)

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

func getEndpointsIntersection(specPEs []mdbv1.PrivateEndpoint, atlasPEs []atlasPE) []intersectionPair {
	intersection := set.Intersection(specPEs, atlasPEs)
	result := []intersectionPair{}
	for _, item := range intersection {
		pair := intersectionPair{}
		pair.spec = item[0].(mdbv1.PrivateEndpoint)
		pair.atlas = item[1].(atlasPE)
		result = append(result, pair)
	}
	return result
}

type intersectionPair struct {
	spec  mdbv1.PrivateEndpoint
	atlas atlasPE
}

func canPrivateEndpointReconcile(atlasClient mongodbatlas.Client, protected bool, akoProject *mdbv1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &mdbv1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	list, err := getAllPrivateEndpoints(atlasClient, akoProject.ID())
	if err != nil {
		return false, err
	}

	if len(list) == 0 {
		return true, nil
	}

	diff, _ := getUniqueDifference(list, latestConfig.PrivateEndpoints)

	if len(diff) == 0 {
		return true, nil
	}

	diff, _ = getUniqueDifference(list, akoProject.Spec.PrivateEndpoints)

	return len(diff) == 0, nil
}
