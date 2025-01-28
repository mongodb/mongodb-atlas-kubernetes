package networkpeering

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

var (
	// ErrNotFound means an resource is missing
	ErrNotFound = errors.New("not found")

	// ErrContainerInUse is a failure to remove a containe still in use
	ErrContainerInUse = errors.New("container still in use")
)

type PeerConnectionsService interface {
	CreatePeer(ctx context.Context, projectID string, conn *NetworkPeer) (*NetworkPeer, error)
	GetPeer(ctx context.Context, projectID, containerID string) (*NetworkPeer, error)
	DeletePeer(ctx context.Context, projectID, containerID string) error
}

type PeeringContainerService interface {
	CreateContainer(ctx context.Context, projectID string, container *ProviderContainer) (*ProviderContainer, error)
	GetContainer(ctx context.Context, projectID, containerID string) (*ProviderContainer, error)
	FindContainer(ctx context.Context, projectID, provider, cidrBlock string) (*ProviderContainer, error)
	DeleteContainer(ctx context.Context, projectID, containerID string) error
}

type NetworkPeeringService interface {
	PeerConnectionsService
	PeeringContainerService
}

type networkPeeringService struct {
	peeringAPI admin.NetworkPeeringApi
}

func NewNetworkPeeringServiceFromClientSet(clientSet *atlas.ClientSet) NetworkPeeringService {
	return NewNetworkPeeringService(clientSet.SdkClient20231115008.NetworkPeeringApi)
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
	newPeer, err := fromAtlasConnection(newAtlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return newPeer, nil
}

func (np *networkPeeringService) GetPeer(ctx context.Context, projectID, peerID string) (*NetworkPeer, error) {
	atlasConn, _, err := np.peeringAPI.GetPeeringConnection(ctx, projectID, peerID).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "PEER_NOT_FOUND") {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get network peer for peer id %v: %w", peerID, err)
	}
	peer, err := fromAtlasConnection(atlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return peer, nil
}

func (np *networkPeeringService) DeletePeer(ctx context.Context, projectID, peerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringConnection(ctx, projectID, peerID).Execute()
	if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") || admin.IsErrorCode(err, "PEER_NOT_FOUND") {
		return nil // if it was already removed or being removed it is also fine
	}
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
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get container %s: %w", containerID, err)
	}
	return fromAtlasContainer(container), nil
}

func (np *networkPeeringService) FindContainer(ctx context.Context, projectID, provider, cidrBlock string) (*ProviderContainer, error) {
	containers, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.CloudProviderContainer], *http.Response, error) {
		return np.peeringAPI.ListPeeringContainerByCloudProvider(ctx, projectID).ProviderName(provider).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers at project %s: %w", projectID, err)
	}
	for _, container := range containers {
		if container.GetAtlasCidrBlock() == cidrBlock {
			return fromAtlasContainer(&container), nil
		}
	}
	return nil, ErrNotFound
}

func (np *networkPeeringService) DeleteContainer(ctx context.Context, projectID, containerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringContainer(ctx, projectID, containerID).Execute()
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return ErrNotFound
	}
	if admin.IsErrorCode(err, "CONTAINERS_IN_USE") {
		return fmt.Errorf("failed to remove container %s as it is still in use: %w", containerID, ErrContainerInUse)
	}
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}
	return nil
}
