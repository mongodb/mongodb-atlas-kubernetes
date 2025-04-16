// Copyright 2024 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atlasprivateendpoint

import (
	"context"
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
)

func (r *AtlasPrivateEndpointReconciler) handlePrivateEndpointService(
	ctx *workflow.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
) (ctrl.Result, error) {
	akoPEService := privateendpoint.NewPrivateEndpoint(akoPrivateEndpoint)
	atlasPEService, err := r.getOrMatchPrivateEndpointService(ctx.Context, privateEndpointService, projectID, akoPrivateEndpoint)
	if err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.ReadyType, workflow.Internal, err)
	}

	wasDeleted := !akoPrivateEndpoint.GetDeletionTimestamp().IsZero()
	existInAtlas := atlasPEService != nil

	switch {
	case !existInAtlas && !wasDeleted:
		return r.createPrivateEndpointService(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService)
	case !existInAtlas && wasDeleted:
		return r.unmanage(ctx, akoPrivateEndpoint)
	case existInAtlas && wasDeleted:
		return r.deletePrivateEndpointService(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
	}

	return r.watchServiceState(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
}

func (r *AtlasPrivateEndpointReconciler) createPrivateEndpointService(
	ctx *workflow.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	akoPEService privateendpoint.EndpointService,
) (ctrl.Result, error) {
	atlasPEService, err := privateEndpointService.CreatePrivateEndpointService(ctx.Context, projectID, akoPEService)
	if err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceFailedToCreate, err)
	}

	return r.watchServiceState(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
}

func (r *AtlasPrivateEndpointReconciler) deletePrivateEndpointService(
	ctx *workflow.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	akoPEService privateendpoint.EndpointService,
	atlasPEService privateendpoint.EndpointService,
) (ctrl.Result, error) {
	if customresource.IsResourcePolicyKeepOrDefault(akoPrivateEndpoint, r.ObjectDeletionProtection) {
		return r.unmanage(ctx, akoPrivateEndpoint)
	}

	atlasPEService, err := r.deletePrivateEndpoint(ctx.Context, privateEndpointService, projectID, atlasPEService)
	if err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointFailedToDelete, err)
	}

	return r.watchServiceState(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
}

func (r *AtlasPrivateEndpointReconciler) watchServiceState(
	ctx *workflow.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	akoPEService privateendpoint.EndpointService,
	atlasPEService privateendpoint.EndpointService,
) (ctrl.Result, error) {
	switch atlasPEService.Status() {
	case privateendpoint.StatusInitiating:
		return r.inProgress(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceInitializing, "Private Endpoint is being initialized")
	case privateendpoint.StatusPending, privateendpoint.StatusPendingAcceptance, privateendpoint.StatusWaitingForUser, privateendpoint.StatusVerified:
		return r.inProgress(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceInitializing, "Private Endpoint is waiting for human action")
	case privateendpoint.StatusFailed, privateendpoint.StatusRejected:
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceFailedToConfigure, errors.New(atlasPEService.ErrorMessage()))
	case privateendpoint.StatusDeleting:
		return r.inProgress(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceDeleting, "Private Endpoint is being deleted")
	}

	ctx.SetConditionTrue(api.PrivateEndpointServiceReady)
	r.EventRecorder.Event(akoPrivateEndpoint, "Normal", string(workflow.PrivateEndpointServiceCreated), "Private Endpoint Service is available")

	return r.handlePrivateEndpointInterface(ctx, privateEndpointService, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
}

func (r *AtlasPrivateEndpointReconciler) handlePrivateEndpointInterface(
	ctx *workflow.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	akoPEService privateendpoint.EndpointService,
	atlasPEService privateendpoint.EndpointService,
) (ctrl.Result, error) {
	compositeInterfacesMap := privateendpoint.MapPrivateEndpoints(akoPEService.EndpointInterfaces(), atlasPEService.EndpointInterfaces())

	if len(compositeInterfacesMap) == 0 {
		return r.waitForConfiguration(ctx, akoPrivateEndpoint, atlasPEService)
	}

	// we want to sync all interface, if any of them is in progress, after all we transition to in progress
	pendingResources := false
	for _, compositeInterfaceMap := range compositeInterfacesMap {
		// The interface can be in 4 state:
		// * It doesn't exist and need to be created, and it's expected to be in progress on next reconciliation
		// * It exist and need to be deleted, and it's expected to in progress on next reconciliation
		// * It's in progress, we skip it to wait it to be ready
		// * It's failed to be provisioned, we terminate the reconciliation
		inAKO := compositeInterfaceMap.AKO != nil
		inAtlas := compositeInterfaceMap.Atlas != nil
		inProgress := isInterfaceInProgress(compositeInterfaceMap.Atlas)
		failed := hasInterfaceFailed(compositeInterfaceMap.Atlas)

		switch {
		case failed:
			return r.terminate(
				ctx,
				akoPrivateEndpoint,
				atlasPEService,
				api.PrivateEndpointReady,
				workflow.PrivateEndpointFailedToConfigure,
				errors.New(compositeInterfaceMap.Atlas.ErrorMessage()),
			)
		case inProgress:
			pendingResources = true
			continue
		case inAKO && !inAtlas:
			gcpProjectID := getGCPProjectID(akoPrivateEndpoint, compositeInterfaceMap.AKO.InterfaceID())
			_, err := privateEndpointService.CreatePrivateEndpointInterface(
				ctx.Context,
				projectID,
				akoPrivateEndpoint.Spec.Provider,
				akoPrivateEndpoint.Status.ServiceID,
				gcpProjectID,
				compositeInterfaceMap.AKO,
			)
			if err != nil {
				return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointFailedToCreate, err)
			}
			pendingResources = true
		case !inAKO && inAtlas:
			err := privateEndpointService.DeleteEndpointInterface(
				ctx.Context,
				projectID,
				akoPrivateEndpoint.Spec.Provider,
				akoPrivateEndpoint.Status.ServiceID,
				compositeInterfaceMap.Atlas.InterfaceID(),
			)
			if err != nil {
				return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointFailedToDelete, err)
			}
			pendingResources = true
		}
	}

	if pendingResources {
		return r.inProgress(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointUpdating, "Private Endpoints are being updated")
	}

	return r.ready(ctx, akoPrivateEndpoint, atlasPEService)
}

// getOrMatchPrivateEndpointService retrieve the project by ID if one is set or try to match by provider/region
// only one private endpoint service per provider/region is allowed
func (r *AtlasPrivateEndpointReconciler) getOrMatchPrivateEndpointService(
	ctx context.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
) (privateendpoint.EndpointService, error) {
	if akoPrivateEndpoint.Status.ServiceID != "" {
		return privateEndpointService.GetPrivateEndpoint(ctx, projectID, akoPrivateEndpoint.Spec.Provider, akoPrivateEndpoint.Status.ServiceID)
	}

	endpointServices, err := privateEndpointService.ListPrivateEndpoints(ctx, projectID, akoPrivateEndpoint.Spec.Provider)
	if err != nil {
		return nil, err
	}

	for _, endpointService := range endpointServices {
		if endpointService.Region() == akoPrivateEndpoint.Spec.Region {
			return endpointService, err
		}
	}

	return nil, nil
}

func (r *AtlasPrivateEndpointReconciler) deletePrivateEndpoint(
	ctx context.Context,
	privateEndpointService privateendpoint.PrivateEndpointService,
	projectID string,
	atlasPEService privateendpoint.EndpointService,
) (privateendpoint.EndpointService, error) {
	if len(atlasPEService.EndpointInterfaces()) == 0 && atlasPEService.Status() != privateendpoint.StatusDeleting {
		err := privateEndpointService.DeleteEndpointService(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID())
		if err != nil {
			return nil, err
		}
	}

	for _, i := range atlasPEService.EndpointInterfaces() {
		if i.Status() != privateendpoint.StatusDeleting {
			err := privateEndpointService.DeleteEndpointInterface(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID(), i.InterfaceID())
			if err != nil {
				return nil, err
			}
		}
	}

	return privateEndpointService.GetPrivateEndpoint(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID())
}

func isInterfaceInProgress(ep privateendpoint.EndpointInterface) bool {
	if ep == nil {
		return false
	}

	status := ep.Status()

	return status == privateendpoint.StatusInitiating ||
		status == privateendpoint.StatusPending ||
		status == privateendpoint.StatusPendingAcceptance ||
		status == privateendpoint.StatusWaitingForUser ||
		status == privateendpoint.StatusVerified ||
		status == privateendpoint.StatusDeleting
}

func hasInterfaceFailed(ep privateendpoint.EndpointInterface) bool {
	if ep == nil {
		return false
	}

	status := ep.Status()

	return status == privateendpoint.StatusFailed || status == privateendpoint.StatusRejected
}

func getGCPProjectID(akoPrivateEndpoint *akov2.AtlasPrivateEndpoint, interfaceID string) string {
	if akoPrivateEndpoint.Spec.Provider == privateendpoint.ProviderGCP && len(akoPrivateEndpoint.Spec.GCPConfiguration) > 0 {
		for _, config := range akoPrivateEndpoint.Spec.GCPConfiguration {
			if config.GroupName == interfaceID {
				return config.ProjectID
			}
		}
	}

	return ""
}
