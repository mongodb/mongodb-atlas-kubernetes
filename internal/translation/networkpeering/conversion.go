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
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
)

var (
	// ErrUnsupportedProvider marks an error when parsing an invalid provider input
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

type NetworkPeer struct {
	akov2.AtlasNetworkPeeringConfig
	ContainerID  string
	Status       string
	ErrorMessage string
	AWSStatus    *status.AWSPeeringStatus
}

func (np *NetworkPeer) Failed() bool {
	return np.ErrorMessage != ""
}

func (np *NetworkPeer) AWSConnectionID() string {
	if np.AWSStatus == nil {
		return ""
	}
	return np.AWSStatus.ConnectionID
}

func (np *NetworkPeer) String() string {
	return fmt.Sprintf("NetworkPeer for %s ID=%s ContainerID=%s\nConfig:%v\nStatus:%v",
		np.Provider, np.ID, np.ContainerID, np.configString(), np.statusString())
}

func (np *NetworkPeer) configString() string {
	aws := ""
	if np.AWSConfiguration != nil {
		cfg := np.AWSConfiguration
		aws = fmt.Sprintf("AWSCfg:{ AccepterRegionName=%s AccountID=%s RouteTableCIDRBlock=%s VpcID=%s } ",
			cfg.AccepterRegionName, cfg.AWSAccountID, cfg.RouteTableCIDRBlock, cfg.VpcID)
	}
	azure := ""
	if np.AzureConfiguration != nil {
		cfg := np.AzureConfiguration
		azure = fmt.Sprintf("AzureCfg:{ AzureDirectoryID=%s AzureSubscriptionID=%s ResourceGroupName=%s VnetName=%s } ",
			cfg.AzureDirectoryID, cfg.AzureSubscriptionID, cfg.ResourceGroupName, cfg.VNetName)
	}
	google := ""
	if np.GCPConfiguration != nil {
		cfg := np.GCPConfiguration
		google = fmt.Sprintf("GoogleCfg:{ GCPProjectID=%s NetworkName=%s } ",
			cfg.GCPProjectID, cfg.NetworkName)
	}
	return fmt.Sprintf("{%s%s%s}", aws, azure, google)
}

func (np *NetworkPeer) statusString() string {
	tail := ""
	if np.AWSStatus != nil {
		tail = fmt.Sprintf(" AWSStatus:{ConnectionId=%s}", np.AWSStatus.ConnectionID)
	}
	return fmt.Sprintf("{Status=%q ErrorMessage=%q%s}", np.Status, np.ErrorMessage, tail)
}

// Available returns whether or not the Network Peering is connected and ready to use
func (np *NetworkPeer) Available() bool {
	return np.Status == "AVAILABLE"
}

// Closing returns whether or not the Network Peering is being shut down
func (np *NetworkPeer) Closing() bool {
	// GCP DELETING AWS TERMINATING AZURE ?
	return np.Status == "DELETING" || np.Status == "TERMINATING"
}

// UpdateStatus copies the network peering status fields only from the given peer input
func (np *NetworkPeer) UpdateStatus(atlas *NetworkPeer) {
	np.Status = atlas.Status
	np.ErrorMessage = atlas.ErrorMessage
	if np.Provider == string(provider.ProviderAWS) && atlas.AWSStatus != nil {
		np.AWSStatus = atlas.AWSStatus.DeepCopy()
	}
}

// NewNetworkPeer creates a network peering from the given config
func NewNetworkPeer(id string, cfg *akov2.AtlasNetworkPeeringConfig) *NetworkPeer {
	peer := &NetworkPeer{
		AtlasNetworkPeeringConfig: *cfg,
	}
	peer.ID = id
	return peer
}

func toAtlas(peer *NetworkPeer) (*admin.BaseNetworkPeeringConnectionSettings, error) {
	switch peer.Provider {
	case string(provider.ProviderAWS):
		if peer.AWSConfiguration == nil {
			return nil, fmt.Errorf("unsupported AWS peer with AWSConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(peer.Provider, ""),
			AccepterRegionName:  pointer.SetOrNil(peer.AWSConfiguration.AccepterRegionName, ""),
			AwsAccountId:        pointer.SetOrNil(peer.AWSConfiguration.AWSAccountID, ""),
			RouteTableCidrBlock: pointer.SetOrNil(peer.AWSConfiguration.RouteTableCIDRBlock, ""),
			VpcId:               pointer.SetOrNil(peer.AWSConfiguration.VpcID, ""),
		}, nil
	case string(provider.ProviderGCP):
		if peer.GCPConfiguration == nil {
			return nil, fmt.Errorf("unsupported Google peer with GCPConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:  peer.ContainerID,
			ProviderName: pointer.SetOrNil(peer.Provider, ""),
			GcpProjectId: pointer.SetOrNil(peer.GCPConfiguration.GCPProjectID, ""),
			NetworkName:  pointer.SetOrNil(peer.GCPConfiguration.NetworkName, ""),
		}, nil
	case string(provider.ProviderAzure):
		if peer.AzureConfiguration == nil {
			return nil, fmt.Errorf("unsupported Azure peer with AzureConfiguration unset")
		}
		return &admin.BaseNetworkPeeringConnectionSettings{
			ContainerId:         peer.ContainerID,
			ProviderName:        pointer.SetOrNil(peer.Provider, ""),
			AzureDirectoryId:    pointer.SetOrNil(peer.AzureConfiguration.AzureDirectoryID, ""),
			AzureSubscriptionId: pointer.SetOrNil(peer.AzureConfiguration.AzureSubscriptionID, ""),
			ResourceGroupName:   pointer.SetOrNil(peer.AzureConfiguration.ResourceGroupName, ""),
			VnetName:            pointer.SetOrNil(peer.AzureConfiguration.VNetName, ""),
		}, nil
	default:
		return nil, fmt.Errorf("%w %q", ErrUnsupportedProvider, peer.Provider)
	}
}

func fromAtlas(conn *admin.BaseNetworkPeeringConnectionSettings) (*NetworkPeer, error) {
	networkPeer, err := fromAtlasConnectionNoStatus(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to convert BaseNetworkPeeringConnectionSettings to NetworkPeer: %w", err)
	}
	switch provider.ProviderName(conn.GetProviderName()) {
	case provider.ProviderAWS:
		networkPeer.Status = conn.GetStatusName()
		networkPeer.ErrorMessage = conn.GetErrorStateName()
		networkPeer.AWSStatus = fromAtlasAWSStatus(conn)
	case provider.ProviderGCP:
		networkPeer.Status = conn.GetStatus()
		networkPeer.ErrorMessage = conn.GetErrorMessage()
	case provider.ProviderAzure:
		networkPeer.Status = conn.GetStatus()
		networkPeer.ErrorMessage = conn.GetErrorState()
	default:
		return nil, fmt.Errorf("%w %q", ErrUnsupportedProvider, conn.GetProviderName())
	}
	return networkPeer, nil
}

func fromAtlasAWSStatus(conn *admin.BaseNetworkPeeringConnectionSettings) *status.AWSPeeringStatus {
	if conn.ConnectionId == nil {
		return nil
	}
	return &status.AWSPeeringStatus{
		ConnectionID: conn.GetConnectionId(),
	}
}

func fromAtlasConnectionNoStatus(conn *admin.BaseNetworkPeeringConnectionSettings) (*NetworkPeer, error) {
	switch provider.ProviderName(conn.GetProviderName()) {
	case provider.ProviderAWS:
		return &NetworkPeer{
			ContainerID: conn.GetContainerId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ID:       conn.GetId(),
				Provider: conn.GetProviderName(),
				AWSConfiguration: &akov2.AWSNetworkPeeringConfiguration{
					AccepterRegionName:  conn.GetAccepterRegionName(),
					AWSAccountID:        conn.GetAwsAccountId(),
					RouteTableCIDRBlock: conn.GetRouteTableCidrBlock(),
					VpcID:               conn.GetVpcId(),
				},
			},
		}, nil
	case provider.ProviderGCP:
		return &NetworkPeer{
			ContainerID: conn.GetContainerId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ID:       conn.GetId(),
				Provider: conn.GetProviderName(),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: conn.GetGcpProjectId(),
					NetworkName:  conn.GetNetworkName(),
				},
			},
		}, nil
	case provider.ProviderAzure:
		return &NetworkPeer{
			ContainerID: conn.GetContainerId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ID:       conn.GetId(),
				Provider: conn.GetProviderName(),
				AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
					AzureDirectoryID:    conn.GetAzureDirectoryId(),
					AzureSubscriptionID: conn.GetAzureSubscriptionId(),
					ResourceGroupName:   conn.GetResourceGroupName(),
					VNetName:            conn.GetVnetName(),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("%w %q", ErrUnsupportedProvider, conn.GetProviderName())
	}
}

func fromAtlasConnectionList(list []admin.BaseNetworkPeeringConnectionSettings) ([]NetworkPeer, error) {
	if list == nil {
		return nil, nil
	}
	peers := make([]NetworkPeer, 0, len(list))
	for i, conn := range list {
		c, err := fromAtlas(&conn)
		if err != nil {
			return nil, fmt.Errorf("failed to convert connection list item %d: %w", i, err)
		}
		peers = append(peers, *c)
	}
	return peers, nil
}

func CompareConfigs(a, b *NetworkPeer) bool {
	aCopy := a.DeepCopy()
	bCopy := b.DeepCopy()
	// accepter region cannot be updated and it might be empty when it matches the container region
	// so we clear it here to avoid finding bogus differences
	if aCopy.AWSConfiguration != nil {
		aCopy.AWSConfiguration.AccepterRegionName = ""
	}
	if bCopy.AWSConfiguration != nil {
		bCopy.AWSConfiguration.AccepterRegionName = ""
	}
	return reflect.DeepEqual(aCopy, bCopy)
}

func ApplyPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus, peer *NetworkPeer, container *networkcontainer.NetworkContainer) {
	peeringStatus.ID = peer.ID
	peeringStatus.Status = peer.Status
	switch provider.ProviderName(peer.Provider) {
	case provider.ProviderAWS:
		if container.AWSStatus == nil && peer.AWSStatus == nil {
			break
		}
		if peeringStatus.AWSStatus == nil {
			peeringStatus.AWSStatus = &status.AWSPeeringStatus{}
		}
		if container.AWSStatus != nil {
			peeringStatus.AWSStatus.VpcID = container.AWSStatus.VpcID
		}
		if peer.AWSStatus != nil {
			peeringStatus.AWSStatus.ConnectionID = peer.AWSStatus.ConnectionID
		}
	case provider.ProviderAzure:
		if container.AzureStatus != nil {
			if peeringStatus.AzureStatus == nil {
				peeringStatus.AzureStatus = &status.AzurePeeringStatus{}
			}
			peeringStatus.AzureStatus.AzureSubscriptionID = container.AzureStatus.AzureSubscriptionID
			peeringStatus.AzureStatus.VnetName = container.AzureStatus.VnetName
		}
	case provider.ProviderGCP:
		if container.GCPStatus != nil {
			if peeringStatus.GCPStatus == nil {
				peeringStatus.GCPStatus = &status.GCPPeeringStatus{}
			}
			peeringStatus.GCPStatus.GCPProjectID = container.GCPStatus.GCPProjectID
			peeringStatus.GCPStatus.NetworkName = container.GCPStatus.NetworkName
		}
	default:
		peeringStatus.Status = fmt.Sprintf("unsupported provider: %q", peer.Provider)
		return
	}
}

func ClearPeeringStatus(peeringStatus *status.AtlasNetworkPeeringStatus) {
	peeringStatus.ID = ""
	peeringStatus.Status = ""
	peeringStatus.AWSStatus = nil
	peeringStatus.AzureStatus = nil
	peeringStatus.GCPStatus = nil
}
