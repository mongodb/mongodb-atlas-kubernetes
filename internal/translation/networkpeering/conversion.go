package networkpeering

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type AWSStatus struct {
	ConnectionID string
}

type NetworkPeer struct {
	akov2.AtlasNetworkPeeringConfig
	ID           string
	Status       string
	ErrorMessage string
	AWSStatus    *AWSStatus
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
		np.AWSStatus = &AWSStatus{
			ConnectionID: atlas.AWSStatus.ConnectionID,
		}
	}
}

// NewNetworkPeer creates a network peering from the given config
func NewNetworkPeer(id string, cfg *akov2.AtlasNetworkPeeringConfig) *NetworkPeer {
	return &NetworkPeer{
		AtlasNetworkPeeringConfig: *cfg,
		ID:                        id,
	}
}

// NewNetworkPeeringSpec creates an spec for network peering from the given config
func NewNetworkPeeringSpec(cfg *akov2.AtlasNetworkPeeringConfig) *NetworkPeer {
	return NewNetworkPeer("", cfg)
}

type ProviderContainer struct {
	akov2.AtlasProviderContainerConfig
	ID           string
	Provider     string
	Provisioned  bool
	AWSStatus    *AWSContainerStatus
	AzureStatus  *AzureContainerStatus
	GoogleStatus *GoogleContainerStatus
}

type AWSContainerStatus struct {
	VpcID string
}

type AzureContainerStatus struct {
	AzureSubscriptionID string
	VnetName            string
}

type GoogleContainerStatus struct {
	GCPProjectID string
	NetworkName  string
}

func NewProviderContainer(id string, provider string, cfg *akov2.AtlasProviderContainerConfig) *ProviderContainer {
	return &ProviderContainer{
		AtlasProviderContainerConfig: *cfg,
		ID:                           id,
		Provider:                     provider,
	}
}

func (pc *ProviderContainer) UpdateStatus(atlas *ProviderContainer) {
	pc.ID = atlas.ID
	pc.Provisioned = atlas.Provisioned
	switch provider.ProviderName(pc.Provider) {
	case provider.ProviderAWS:
		pc.AWSStatus = atlas.AWSStatus
	case provider.ProviderAzure:
		pc.AzureStatus = atlas.AzureStatus
	case provider.ProviderGCP:
		pc.GoogleStatus = atlas.GoogleStatus
	}
}

func (pc *ProviderContainer) String() string {
	return fmt.Sprintf("ProviderContainer for %s ID=%s\nConfig:%v\nStatus:%v",
		pc.Provider, pc.ID, pc.configString(), pc.statusString())
}

func (pc *ProviderContainer) configString() string {
	return fmt.Sprintf("{ ContainerRegion=%s AtlasCIDRBlock=%s }", pc.ContainerRegion, pc.AtlasCIDRBlock)
}

func (pc *ProviderContainer) statusString() string {
	aws := ""
	if pc.AWSStatus != nil {
		status := pc.AWSStatus
		aws = fmt.Sprintf("AWSStatus:{ VpcID=%s } ", status.VpcID)
	}
	azure := ""
	if pc.AzureStatus != nil {
		status := pc.AzureStatus
		azure = fmt.Sprintf("AzureStatus:{ AzureSubscriptionID=%s VnetName=%s } ",
			status.AzureSubscriptionID, status.VnetName)
	}
	google := ""
	if pc.GoogleStatus != nil {
		status := pc.GoogleStatus
		google = fmt.Sprintf("GoogleStatus:{GCPProjectID=%s NetworkName=%s } ",
			status.GCPProjectID, status.NetworkName)
	}
	return fmt.Sprintf("{ Provisioned=%v %s%s%s}", pc.Provisioned, aws, azure, google)
}

func toAtlasConnection(peer *NetworkPeer) (*admin.BaseNetworkPeeringConnectionSettings, error) {
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
		return nil, fmt.Errorf("unsupported provider %q", peer.Provider)
	}
}

func fromAtlasConnection(conn *admin.BaseNetworkPeeringConnectionSettings) (*NetworkPeer, error) {
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
		return nil, fmt.Errorf("unsupported provider %q", conn.GetProviderName())
	}
	return networkPeer, nil
}

func fromAtlasAWSStatus(conn *admin.BaseNetworkPeeringConnectionSettings) *AWSStatus {
	if conn.ConnectionId == nil {
		return nil
	}
	return &AWSStatus{
		ConnectionID: conn.GetConnectionId(),
	}
}

func fromAtlasConnectionNoStatus(conn *admin.BaseNetworkPeeringConnectionSettings) (*NetworkPeer, error) {
	switch provider.ProviderName(conn.GetProviderName()) {
	case provider.ProviderAWS:
		return &NetworkPeer{
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
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
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
				GCPConfiguration: &akov2.GCPNetworkPeeringConfiguration{
					GCPProjectID: conn.GetGcpProjectId(),
					NetworkName:  conn.GetNetworkName(),
				},
			},
		}, nil
	case provider.ProviderAzure:
		return &NetworkPeer{
			ID: conn.GetId(),
			AtlasNetworkPeeringConfig: akov2.AtlasNetworkPeeringConfig{
				ContainerID: conn.GetContainerId(),
				Provider:    conn.GetProviderName(),
				AzureConfiguration: &akov2.AzureNetworkPeeringConfiguration{
					AzureDirectoryID:    conn.GetAzureDirectoryId(),
					AzureSubscriptionID: conn.GetAzureSubscriptionId(),
					ResourceGroupName:   conn.GetResourceGroupName(),
					VNetName:            conn.GetVnetName(),
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", conn.GetProviderName())
	}
}

func fromAtlasConnectionList(list []admin.BaseNetworkPeeringConnectionSettings) ([]NetworkPeer, error) {
	if list == nil {
		return nil, nil
	}
	peers := make([]NetworkPeer, 0, len(list))
	for i, conn := range list {
		c, err := fromAtlasConnection(&conn)
		if err != nil {
			return nil, fmt.Errorf("failed to convert connection list item %d: %w", i, err)
		}
		peers = append(peers, *c)
	}
	return peers, nil
}

func toAtlasContainer(container *ProviderContainer) *admin.CloudProviderContainer {
	cpc := &admin.CloudProviderContainer{
		Id:             pointer.SetOrNil(container.ID, ""),
		ProviderName:   pointer.SetOrNil(container.Provider, ""),
		AtlasCidrBlock: pointer.SetOrNil(container.AtlasCIDRBlock, ""),
	}
	if cpc.GetProviderName() == string(provider.ProviderAWS) {
		cpc.RegionName = pointer.SetOrNil(container.ContainerRegion, "")
	} else {
		cpc.Region = pointer.SetOrNil(container.ContainerRegion, "")
	}
	return cpc
}

func fromAtlasContainer(container *admin.CloudProviderContainer) *ProviderContainer {
	pc := fromAtlasContainerNoStatus(container)
	pc.Provisioned = container.GetProvisioned()
	switch provider.ProviderName(pc.Provider) {
	case provider.ProviderAWS:
		pc.AWSStatus = fromAtlasAWSContainerStatus(container)
	case provider.ProviderAzure:
		pc.AzureStatus = fromAtlasAzureContainerStatus(container)
	case provider.ProviderGCP:
		pc.GoogleStatus = fromAtlasGoogleContainerStatus(container)
	}
	return pc
}

func fromAtlasAWSContainerStatus(container *admin.CloudProviderContainer) *AWSContainerStatus {
	if container.VpcId == nil {
		return nil
	}
	return &AWSContainerStatus{
		VpcID: container.GetVpcId(),
	}
}

func fromAtlasAzureContainerStatus(container *admin.CloudProviderContainer) *AzureContainerStatus {
	if container.AzureSubscriptionId == nil && container.VnetName == nil {
		return nil
	}
	return &AzureContainerStatus{
		AzureSubscriptionID: container.GetAzureSubscriptionId(),
		VnetName:            container.GetVnetName(),
	}
}

func fromAtlasGoogleContainerStatus(container *admin.CloudProviderContainer) *GoogleContainerStatus {
	if container.GcpProjectId == nil && container.NetworkName == nil {
		return nil
	}
	return &GoogleContainerStatus{
		GCPProjectID: container.GetGcpProjectId(),
		NetworkName:  container.GetNetworkName(),
	}
}

func fromAtlasContainerNoStatus(container *admin.CloudProviderContainer) *ProviderContainer {
	region := container.GetRegion()
	if container.GetProviderName() == string(provider.ProviderAWS) {
		region = container.GetRegionName()
	}
	return &ProviderContainer{
		ID:       container.GetId(),
		Provider: container.GetProviderName(),
		AtlasProviderContainerConfig: akov2.AtlasProviderContainerConfig{
			AtlasCIDRBlock:  container.GetAtlasCidrBlock(),
			ContainerRegion: region,
		},
	}
}

func fromAtlasContainerList(list []admin.CloudProviderContainer) []ProviderContainer {
	if list == nil {
		return nil
	}
	containers := make([]ProviderContainer, 0, len(list))
	for _, container := range list {
		containers = append(containers, *fromAtlasContainer(&container))
	}
	return containers
}
