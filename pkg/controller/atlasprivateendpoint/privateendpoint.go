/*
Copyright 2024 MongoDB.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package atlasprivateendpoint

import (
	"context"
	"errors"
	"reflect"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/privateendpoint"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasPrivateEndpointReconciler) handlePrivateEndpointService(
	ctx *workflow.Context,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
) (ctrl.Result, error) {
	akoPEService := privateendpoint.NewPrivateEndpoint(akoPrivateEndpoint)
	atlasPEService, err := r.getPrivateEndpointService(ctx.Context, projectID, akoPrivateEndpoint)
	if err != nil {
		return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.ReadyType, workflow.Internal, err)
	}

	wasDeleted := !akoPrivateEndpoint.GetDeletionTimestamp().IsZero()
	existInAtlas := atlasPEService != nil

	switch {
	case !existInAtlas && !wasDeleted:
		atlasPEService, err = r.privateEndpointService.CreatePrivateEndpointService(ctx.Context, projectID, akoPEService)
		if err != nil {
			return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointServiceFailedToCreate, err)
		}
	case !existInAtlas && wasDeleted:
		return r.unmanage(ctx, akoPrivateEndpoint)
	case existInAtlas && wasDeleted:
		if customresource.IsResourcePolicyKeepOrDefault(akoPrivateEndpoint, r.ObjectDeletionProtection) {
			return r.unmanage(ctx, akoPrivateEndpoint)
		}

		atlasPEService, err = r.deletePrivateEndpoint(ctx.Context, projectID, atlasPEService)
		if err != nil {
			return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointServiceReady, workflow.PrivateEndpointFailedToDelete, err)
		}
	}

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

	return r.handlePrivateEndpointInterface(ctx, projectID, akoPrivateEndpoint, akoPEService, atlasPEService)
}

func (r *AtlasPrivateEndpointReconciler) handlePrivateEndpointInterface(
	ctx *workflow.Context,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
	akoPEService privateendpoint.EndpointService,
	atlasPEService privateendpoint.EndpointService,
) (ctrl.Result, error) {
	if len(akoPEService.EndpointInterfaces()) == 0 && len(atlasPEService.EndpointInterfaces()) == 0 {
		return r.waitForConfiguration(ctx, akoPrivateEndpoint, atlasPEService)
	}

	inProgress := false
	for _, akoPEInterface := range akoPEService.EndpointInterfaces() {
		atlasPEInterfaces := atlasPEService.EndpointInterfaces()
		atlasPEInterface := atlasPEInterfaces.Get(akoPEInterface.InterfaceID())
		existInAtlas := atlasPEInterface != nil
		inProgress = isInterfaceInProgress(akoPEInterface.Status())
		var err error

		if !existInAtlas {
			gcpProjectID := getGCPProjectID(akoPrivateEndpoint, akoPEInterface.InterfaceID())
			_, err = r.privateEndpointService.CreatePrivateEndpointInterface(
				ctx.Context,
				projectID,
				akoPrivateEndpoint.Spec.Provider,
				akoPrivateEndpoint.Status.ServiceID,
				gcpProjectID,
				akoPEInterface,
			)
			if err != nil {
				return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointFailedToCreate, err)
			}

			inProgress = true
		}
	}

	for _, atlasPEInterface := range atlasPEService.EndpointInterfaces() {
		akoPEInterfaces := akoPEService.EndpointInterfaces()
		akoPEInterface := akoPEInterfaces.Get(atlasPEInterface.InterfaceID())
		wasDeleted := akoPEInterface == nil
		inProgress = isInterfaceInProgress(atlasPEInterface.Status())
		var err error

		if wasDeleted && atlasPEInterface.Status() != privateendpoint.StatusDeleting {
			err = r.privateEndpointService.DeleteEndpointInterface(ctx.Context, projectID, akoPrivateEndpoint.Spec.Provider, akoPrivateEndpoint.Status.ServiceID, atlasPEInterface.InterfaceID())
			if err != nil {
				return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointFailedToDelete, err)
			}

			inProgress = true

			continue
		}

		if hasInterfaceFailed(atlasPEInterface.Status()) {
			return r.terminate(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointFailedToConfigure, errors.New(atlasPEInterface.ErrorMessage()))
		}
	}

	if !inProgress && reflect.DeepEqual(akoPEService.EndpointInterfaces(), atlasPEService.EndpointInterfaces()) {
		return r.ready(ctx, akoPrivateEndpoint, atlasPEService)
	}

	return r.inProgress(ctx, akoPrivateEndpoint, atlasPEService, api.PrivateEndpointReady, workflow.PrivateEndpointUpdating, "Private Endpoints are being updated")
}

func (r *AtlasPrivateEndpointReconciler) getPrivateEndpointService(
	ctx context.Context,
	projectID string,
	akoPrivateEndpoint *akov2.AtlasPrivateEndpoint,
) (privateendpoint.EndpointService, error) {
	if akoPrivateEndpoint.Status.ServiceID != "" {
		return r.privateEndpointService.GetPrivateEndpoint(ctx, projectID, akoPrivateEndpoint.Spec.Provider, akoPrivateEndpoint.Status.ServiceID)
	}

	endpointServices, err := r.privateEndpointService.ListPrivateEndpoints(ctx, projectID, akoPrivateEndpoint.Spec.Provider)
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
	projectID string,
	atlasPEService privateendpoint.EndpointService,
) (privateendpoint.EndpointService, error) {
	if len(atlasPEService.EndpointInterfaces()) == 0 && atlasPEService.Status() != privateendpoint.StatusDeleting {
		err := r.privateEndpointService.DeleteEndpointService(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID())
		if err != nil {
			return atlasPEService, err
		}
	}

	for _, i := range atlasPEService.EndpointInterfaces() {
		if i.Status() != privateendpoint.StatusDeleting {
			err := r.privateEndpointService.DeleteEndpointInterface(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID(), i.InterfaceID())
			if err != nil {
				return atlasPEService, err
			}
		}
	}

	return r.privateEndpointService.GetPrivateEndpoint(ctx, projectID, atlasPEService.Provider(), atlasPEService.ServiceID())
}

func isInterfaceInProgress(status string) bool {
	return status == privateendpoint.StatusInitiating ||
		status == privateendpoint.StatusPending ||
		status == privateendpoint.StatusPendingAcceptance ||
		status == privateendpoint.StatusWaitingForUser ||
		status == privateendpoint.StatusVerified ||
		status == privateendpoint.StatusDeleting
}

func hasInterfaceFailed(status string) bool {
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
