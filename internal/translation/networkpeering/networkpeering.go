package networkpeering

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

type PeerConnectionsService interface {
	CreatePeer(ctx context.Context, projectID string, conn *akov2.NetworkPeer) (*akov2.NetworkPeer, error)
	ListPeers(ctx context.Context, projectID string) ([]akov2.NetworkPeer, error)
	DeletePeer(ctx context.Context, projectID, containerID string) error
}

type PeeringContainerService interface {
	CreateContainer(ctx context.Context, projectID string, container *ProviderContainer) (*ProviderContainer, error)
	GetContainer(ctx context.Context, projectID, containerID string) (*ProviderContainer, error)
	ListContainers(ctx context.Context, projectID, providerName string) ([]ProviderContainer, error)
	DeleteContainers(ctx context.Context, projectID, containerID string) error
}

type NetworkPeeringService interface {
	PeerConnectionsService
	PeeringContainerService
}

type NetworkPeering struct {
	connsService      admin.NetworkPeeringApi
	containersService admin.NetworkPeeringApi
}

func (np *NetworkPeering) CreatePeer(ctx context.Context, projectID string, conn *akov2.NetworkPeer) (*akov2.NetworkPeer, error) {
	newConn, _, err := np.connsService.CreatePeeringConnection(ctx, projectID, toAtlasConnection(conn)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create network peer %v: %w", conn, err)
	}
	return fromAtlasConnection(newConn), nil
}

func (np *NetworkPeering) ListPeers(ctx context.Context, projectID string) ([]akov2.NetworkPeer, error) {
	var peersList []akov2.NetworkPeer
	providers := []provider.ProviderName{provider.ProviderAWS, provider.ProviderAzure, provider.ProviderGCP}
	for _, providerName := range providers {
		peers, err := np.listPeersForProvider(ctx, projectID, providerName)
		if err != nil {
			return nil, fmt.Errorf("failed to list network peers for %s: %w", string(providerName), err)
		}
		peersList = append(peersList, peers...)
	}
	return peersList, nil
}

func (np *NetworkPeering) listPeersForProvider(ctx context.Context, projectID string, providerName provider.ProviderName) ([]akov2.NetworkPeer, error) {
	results := []akov2.NetworkPeer{}
	pageNum := 1
	listOpts := &admin.ListPeeringConnectionsApiParams{
		GroupId:      projectID,
		ProviderName: admin.PtrString(string(providerName)),
		PageNum:      pointer.MakePtr(pageNum),
	}
	for {
		page, _, err := np.connsService.ListPeeringConnectionsWithParams(ctx, listOpts).Execute()
		if err != nil {
			return nil, fmt.Errorf("failed to list network peers: %w", err)
		}
		results = append(results, fromAtlasConnectionList(page.GetResults())...)
		if len(results) >= page.GetTotalCount() {
			return results, nil
		}
		pageNum += 1
	}
}

func (np *NetworkPeering) DeletePeer(ctx context.Context, projectID, containerID string) error {
	_, _, err := np.connsService.DeletePeeringConnection(ctx, projectID, containerID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete peering connection for container %s: %w", containerID, err)
	}
	return nil
}

func (np *NetworkPeering) CreateContainer(ctx context.Context, projectID string, container *ProviderContainer) (*ProviderContainer, error) {
	newContainer, _, err := np.containersService.CreatePeeringContainer(ctx, projectID, toAtlasContainer(container)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create peering container %s: %w", container.ID, err)
	}
	return fromAtlasContainer(newContainer), nil
}

func (np *NetworkPeering) GetContainer(ctx context.Context, projectID, containerID string) (*ProviderContainer, error) {
	container, _, err := np.containersService.GetPeeringContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get container for gcp status %s: %w", containerID, err)
	}
	return fromAtlasContainer(container), nil
}

func (np *NetworkPeering) ListContainers(ctx context.Context, projectID, providerName string) ([]ProviderContainer, error) {
	results := []ProviderContainer{}
	pageNum := 1
	listOpts := &admin.ListPeeringContainerByCloudProviderApiParams{
		GroupId:      projectID,
		ProviderName: pointer.SetOrNil(providerName, ""),
		PageNum:      pointer.MakePtr(pageNum),
	}
	for {
		page, _, err := np.containersService.ListPeeringContainerByCloudProviderWithParams(ctx, listOpts).Execute()
		if err != nil {
			return nil, fmt.Errorf("failed to list containers: %w", err)
		}
		results = append(results, fromAtlasContainerList(page.GetResults())...)
		if len(results) >= page.GetTotalCount() {
			return results, nil
		}
		pageNum += 1
	}
}

func (np *NetworkPeering) DeleteContainers(ctx context.Context, projectID, containerID string) error {
	_, _, err := np.connsService.DeletePeeringContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}
	return nil
}
