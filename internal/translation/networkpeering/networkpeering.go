package networkpeering

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

var (
	// ErrNotFound means an resource is missing
	ErrNotFound = errors.New("not found")
)

type NetworkPeeringService interface {
	Create(ctx context.Context, projectID, containerID string, cfg *akov2.AtlasNetworkPeeringConfig) (*NetworkPeer, error)
	Get(ctx context.Context, projectID, peerID string) (*NetworkPeer, error)
	Update(ctx context.Context, pojectID, peerID, containerID string, cfg *akov2.AtlasNetworkPeeringConfig) (*NetworkPeer, error)
	Delete(ctx context.Context, projectID, peerID string) error
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

func (np *networkPeeringService) Create(ctx context.Context, projectID, containerID string, cfg *akov2.AtlasNetworkPeeringConfig) (*NetworkPeer, error) {
	atlasConnRequest, err := toAtlas(&NetworkPeer{
		AtlasNetworkPeeringConfig: *cfg,
		ContainerID:               containerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer to Atlas: %w", err)
	}
	newAtlasConn, _, err := np.peeringAPI.CreatePeeringConnection(ctx, projectID, atlasConnRequest).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create network peer from config %v: %w", cfg, err)
	}
	newPeer, err := fromAtlas(newAtlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return newPeer, nil
}

func (np *networkPeeringService) Get(ctx context.Context, projectID, peerID string) (*NetworkPeer, error) {
	atlasConn, _, err := np.peeringAPI.GetPeeringConnection(ctx, projectID, peerID).Execute()
	if err != nil {
		if admin.IsErrorCode(err, "PEER_NOT_FOUND") {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get network peer for peer id %v: %w", peerID, err)
	}
	peer, err := fromAtlas(atlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return peer, nil
}

func (np *networkPeeringService) Update(ctx context.Context, projectID, peerID, containerID string, cfg *akov2.AtlasNetworkPeeringConfig) (*NetworkPeer, error) {
	atlasConnRequest, err := toAtlas(&NetworkPeer{
		AtlasNetworkPeeringConfig: *cfg,
		ContainerID:               containerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer to Atlas: %w", err)
	}
	newAtlasConn, _, err := np.peeringAPI.UpdatePeeringConnection(ctx, projectID, peerID, atlasConnRequest).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update network peer from config %v: %w", cfg, err)
	}
	newPeer, err := fromAtlas(newAtlasConn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert peer from Atlas: %w", err)
	}
	return newPeer, nil
}

func (np *networkPeeringService) Delete(ctx context.Context, projectID, peerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringConnection(ctx, projectID, peerID).Execute()
	if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") || admin.IsErrorCode(err, "PEER_NOT_FOUND") {
		return errors.Join(err, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to delete peering connection for peer %s: %w", peerID, err)
	}
	return nil
}
