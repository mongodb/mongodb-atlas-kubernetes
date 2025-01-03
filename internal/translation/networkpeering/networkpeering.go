package networkpeering

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
)

type PeerConnectionsService interface {
	CreatePeer(ctx context.Context, projectID string, conn *NetworkPeer) (*NetworkPeer, error)
	ListPeers(ctx context.Context, projectID string) ([]NetworkPeer, error)
	DeletePeer(ctx context.Context, projectID, containerID string) error
}

type PeeringContainerService interface {
	CreateContainer(ctx context.Context, projectID string, container *ProviderContainer) (*ProviderContainer, error)
	GetContainer(ctx context.Context, projectID, containerID string) (*ProviderContainer, error)
	ListContainers(ctx context.Context, projectID, providerName string) ([]ProviderContainer, error)
	DeleteContainer(ctx context.Context, projectID, containerID string) error
}

type NetworkPeeringService interface {
	PeerConnectionsService
	PeeringContainerService
}

type networkPeeringService struct {
	peeringAPI admin.NetworkPeeringApi
}

func NewNetworkPeeringService(peeringAPI admin.NetworkPeeringApi) NetworkPeeringService {
	return &networkPeeringService{peeringAPI: peeringAPI}
}

func (np *networkPeeringService) CreatePeer(ctx context.Context, projectID string, conn *NetworkPeer) (*NetworkPeer, error) {
	atlasConnRequest, err := toAtlasConnection(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer to Atlas: %w", err)
	}
	newAtlasConn, _, err := np.peeringAPI.CreatePeeringConnection(ctx, projectID, atlasConnRequest).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create network peer %v: %w", conn, err)
	}
	newConn, err := fromAtlasConnection(newAtlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return newConn, nil
}

func (np *networkPeeringService) ListPeers(ctx context.Context, projectID string) ([]NetworkPeer, error) {
	var peersList []NetworkPeer
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

func (np *networkPeeringService) listPeersForProvider(ctx context.Context, projectID string, providerName provider.ProviderName) ([]NetworkPeer, error) {
	results, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.BaseNetworkPeeringConnectionSettings], *http.Response, error) {
		p := &admin.ListPeeringConnectionsApiParams{
			GroupId:      projectID,
			ProviderName: admin.PtrString(string(providerName)),
		}
		return np.peeringAPI.ListPeeringConnectionsWithParams(ctx, p).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers: %w", err)
	}

	return fromAtlasConnectionList(results)
}

func (np *networkPeeringService) DeletePeer(ctx context.Context, projectID, peerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringConnection(ctx, projectID, peerID).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete peering connection for peer %s: %w", peerID, err)
	}
	return nil
}

func (np *networkPeeringService) CreateContainer(ctx context.Context, projectID string, container *ProviderContainer) (*ProviderContainer, error) {
	newContainer, _, err := np.peeringAPI.CreatePeeringContainer(ctx, projectID, toAtlasContainer(container)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create peering container %s: %w", container.ID, err)
	}
	return fromAtlasContainer(newContainer), nil
}

func (np *networkPeeringService) GetContainer(ctx context.Context, projectID, containerID string) (*ProviderContainer, error) {
	container, _, err := np.peeringAPI.GetPeeringContainer(ctx, projectID, containerID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get container for gcp status %s: %w", containerID, err)
	}
	return fromAtlasContainer(container), nil
}

func (np *networkPeeringService) ListContainers(ctx context.Context, projectID, providerName string) ([]ProviderContainer, error) {
	results := []ProviderContainer{}
	pageNum := 1
	listOpts := &admin.ListPeeringContainerByCloudProviderApiParams{
		GroupId:      projectID,
		ProviderName: pointer.SetOrNil(providerName, ""),
		PageNum:      pointer.MakePtr(pageNum),
	}
	for {
		page, _, err := np.peeringAPI.ListPeeringContainerByCloudProviderWithParams(ctx, listOpts).Execute()
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

func (np *networkPeeringService) DeleteContainer(ctx context.Context, projectID, containerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringContainer(ctx, projectID, containerID).Execute()
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}
	return nil
}
