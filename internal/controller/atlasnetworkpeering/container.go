package atlasnetworkpeering

import (
	"errors"
	"fmt"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

func (r *AtlasNetworkPeeringReconciler) handleContainer(req *reconcileRequest) (*networkpeering.ProviderContainer, error) {
	atlasContainer, err := discoverContainer(req)
	if err != nil && !errors.Is(err, networkpeering.ErrNotFound) {
		return nil, fmt.Errorf("failed to discover container: %w", err)
	}
	inAtlas := atlasContainer != nil
	deleted := req.networkPeering.DeletionTimestamp != nil
	switch {
	case !deleted && !inAtlas:
		return r.createContainer(req)
	case !deleted && inAtlas:
		return atlasContainer, nil
	case deleted && inAtlas:
		return r.deleteContainer(req)
	default:
		return r.unmanageContainer()
	}
}

func discoverContainer(req *reconcileRequest) (*networkpeering.ProviderContainer, error) {
	containerID := containerID(req.networkPeering)
	if containerID != "" {
		container, err := req.service.GetContainer(req.workflowCtx.Context, req.projectID, containerID)
		if err != nil {
			return nil, fmt.Errorf("failed to get container from project %s and container %s: %w",
				req.projectID, containerID, err)
		}
		return container, nil
	}
	container, err := req.service.FindContainer(
		req.workflowCtx.Context,
		req.projectID,
		req.networkPeering.Spec.Provider,
		req.networkPeering.Spec.AtlasCIDRBlock,
	)
	if errors.Is(err, networkpeering.ErrNotFound) {
		return nil, nil // no existing container
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get find container with CIDR %q for provider %s at project %s: %w",
			req.networkPeering.Spec.AtlasCIDRBlock,
			req.networkPeering.Spec.Provider,
			req.projectID,
			err,
		)
	}
	return container, nil
}

func (r *AtlasNetworkPeeringReconciler) createContainer(req *reconcileRequest) (*networkpeering.ProviderContainer, error) {
	specContainer := networkpeering.NewProviderContainer(
		"",
		req.networkPeering.Spec.Provider,
		&req.networkPeering.Spec.AtlasProviderContainerConfig,
	)
	fixContainerRegion(specContainer, req)
	createdContainer, err :=
		req.service.CreateContainer(req.workflowCtx.Context, req.projectID, specContainer)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	return createdContainer, nil
}

func fixContainerRegion(specContainer *networkpeering.ProviderContainer, req *reconcileRequest) {
	// in AWS it is assumed not specifying the reagion means it is the same as the accepter region
	if req.networkPeering.Spec.Provider == string(provider.ProviderAWS) &&
		specContainer.ContainerRegion == "" &&
		req.networkPeering.Spec.AWSConfiguration != nil {
		specContainer.ContainerRegion = req.networkPeering.Spec.AWSConfiguration.AccepterRegionName
	}
}

func (r *AtlasNetworkPeeringReconciler) deleteContainer(req *reconcileRequest) (*networkpeering.ProviderContainer, error) {
	containerID := containerID(req.networkPeering)
	err := req.service.DeleteContainer(req.workflowCtx.Context, req.projectID, containerID)
	if err == nil || errors.Is(err, networkpeering.ErrNotFound) {
		return r.unmanageContainer()
	}
	return nil, fmt.Errorf("failed to delete container %s: %w", containerID, err)
}

func (r *AtlasNetworkPeeringReconciler) unmanageContainer() (*networkpeering.ProviderContainer, error) {
	return nil, nil
}

func containerID(peer *akov2.AtlasNetworkPeering) string {
	if peer.Spec.ContainerID != "" {
		return peer.Spec.ContainerID
	}
	return peer.Status.ContainerID
}
