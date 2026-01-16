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

package networkpeering

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312012/admin"

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
	return NewNetworkPeeringService(clientSet.SdkClient20250312011.NetworkPeeringApi)
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
	newAtlasConn, _, err := np.peeringAPI.CreateGroupPeer(ctx, projectID, atlasConnRequest).Execute()
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
	atlasConn, _, err := np.peeringAPI.GetGroupPeer(ctx, projectID, peerID).Execute()
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
	newAtlasConn, _, err := np.peeringAPI.UpdateGroupPeer(ctx, projectID, peerID, atlasConnRequest).Execute()
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
	_, _, err := np.peeringAPI.DeleteGroupPeer(ctx, projectID, peerID).Execute()
	if admin.IsErrorCode(err, "PEER_ALREADY_REQUESTED_DELETION") || admin.IsErrorCode(err, "PEER_NOT_FOUND") {
		return errors.Join(err, ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("failed to delete peering connection for peer %s: %w", peerID, err)
	}
	return nil
}
