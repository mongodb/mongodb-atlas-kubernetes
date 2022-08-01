package atlasproject

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

const (
	StatusFailed   = "FAILED"
	StatusReady    = "AVAILABLE"
	StatusDeleting = "DELETING"
)

type networkPeerDiff struct {
	PeersToDelete []string
	PeersToCreate []mdbv1.NetworkPeer
	PeersToUpdate []mongodbatlas.Peer
}

func ensureNetworkPeers(ctx *workflow.Context, groupID string, project *mdbv1.AtlasProject) workflow.Result {
	networkPeerStatus := project.Status.DeepCopy().NetworkPeers
	networkPeerSpec := project.Spec.DeepCopy().NetworkPeers

	backgroundContext := context.Background()
	result, condition := SyncNetworkPeer(backgroundContext, ctx, groupID, networkPeerStatus, networkPeerSpec)
	if !result.IsOk() {
		ctx.SetConditionFromResult(condition, result)
		return result
	}
	ctx.Log.Debugf("network peers are ready! hooray!")
	ctx.SetConditionTrue(status.NetworkPeerReadyType)
	if len(networkPeerStatus) == 0 && len(networkPeerSpec) == 0 {
		ctx.UnsetCondition(status.NetworkPeerReadyType)
	}

	return result
}

func failedPeerStatus(errMessage string, peer mdbv1.NetworkPeer) status.AtlasNetworkPeer {
	var vpc string
	switch peer.ProviderName {
	case provider.ProviderGCP:
		vpc = peer.NetworkName
	case provider.ProviderAzure:
		vpc = peer.VNetName
	default:
		vpc = peer.VpcID
	}
	return status.AtlasNetworkPeer{
		Status:       StatusFailed,
		ErrorMessage: errMessage,
		Name:         vpc,
	}
}

func SyncNetworkPeer(context context.Context, ctx *workflow.Context, groupID string, peersStatus []status.AtlasNetworkPeer, peersSpec []mdbv1.NetworkPeer) (workflow.Result, status.ConditionType) {
	defer ctx.EnsureStatusOption(status.AtlasProjectSetNetworkPeerOption(&peersStatus))
	logger := ctx.Log
	logger.Debugf("existed network peers status: %v", peersStatus)
	mongoClient := ctx.Client
	logger.Debugf("syncing network peers for project %v", groupID)
	list, err := getAllExistedNetworkPeer(context, logger, mongoClient.Peers, groupID)
	if err != nil {
		logger.Errorf("failed to get all network peers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to get all network peers"),
			status.NetworkPeerReadyType
	}

	diff, err := sortPeers(list, peersSpec, logger, mongoClient.Containers, groupID)
	if err != nil {
		logger.Errorf("failed to sort network peers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to sort network peers"),
			status.NetworkPeerReadyType
	}
	logger.Debugf("peers to create %d, peers to update %d, peers to delete %d",
		len(diff.PeersToCreate), len(diff.PeersToUpdate), len(diff.PeersToDelete))

	for _, peerToDelete := range diff.PeersToDelete {
		errDelete := deletePeerByID(context, mongoClient.Peers, groupID, peerToDelete, logger)
		if errDelete != nil {
			logger.Errorf("failed to delete network peer %s: %v", peerToDelete, errDelete)
			return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to delete network peer"),
				status.NetworkPeerReadyType
		}
	}

	peersStatus = createNetworkPeers(context, mongoClient, groupID, diff.PeersToCreate, logger)

	for _, peerToUpdate := range diff.PeersToUpdate {
		vpc := formVPC(peerToUpdate)
		peersStatus = append(peersStatus, status.FromAtlas(peerToUpdate,
			provider.ProviderName(peerToUpdate.ProviderName), vpc))
	}

	return ensurePeerStatus(peersStatus, logger)
}

func formVPC(peer mongodbatlas.Peer) string {
	switch peer.ProviderName {
	case string(provider.ProviderGCP):
		return peer.NetworkName
	case string(provider.ProviderAzure):
		return peer.VNetName
	default:
		return peer.VpcID
	}
}

func ensurePeerStatus(peersStatus []status.AtlasNetworkPeer, logger *zap.SugaredLogger) (workflow.Result, status.ConditionType) {
	for _, peerStatus := range peersStatus {
		switch peerStatus.ProviderName {
		case provider.ProviderGCP:
			if peerStatus.Status != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.Name, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
		case provider.ProviderAzure:
			// TODO: check Azure network peer. status vs statusName
		default:
			if peerStatus.StatusName != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.Name, peerStatus.StatusName)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
		}
	}
	return workflow.OK(), status.NetworkPeerReadyType
}

func createNetworkPeers(context context.Context, mongoClient mongodbatlas.Client, groupID string, peers []mdbv1.NetworkPeer, logger *zap.SugaredLogger) []status.AtlasNetworkPeer {
	var newPeerStatuses []status.AtlasNetworkPeer
	for _, peer := range peers {
		err := validateInitNetworkPeer(peer)
		if err != nil {
			newPeerStatuses = append(newPeerStatuses,
				failedPeerStatus(fmt.Errorf("failed to validate network peer %w", err).Error(), peer))
			logger.Errorf("failed to validate network peer: %s", err)
			continue
		}
		// if containerId is empty, then we need to create a container first by using the spec
		if peer.ContainerID == "" {
			containerID, errCreate := createContainer(context, mongoClient.Containers, groupID, peer, logger)
			if errCreate != nil {
				newPeerStatuses = append(newPeerStatuses,
					failedPeerStatus(fmt.Errorf("failed to create container for network peer %w", errCreate).Error(), peer))
				logger.Errorf("failed to create container for network peer: %s", errCreate)
				continue
			}
			peer.ContainerID = containerID
		}

		atlasPeer, err := createNetworkPeer(context, groupID, mongoClient.Peers, peer, logger)
		if err != nil {
			logger.Errorf("failed to create network peer: %v", err)
			newPeerStatuses = append(newPeerStatuses,
				failedPeerStatus(fmt.Errorf("failed to create network peer: %w", err).Error(), peer))
			continue
		}
		if atlasPeer != nil {
			vpc := formVPC(*atlasPeer)
			if atlasPeer.AccepterRegionName == "" {
				atlasPeer.AccepterRegionName = peer.AccepterRegionName
			}
			newPeerStatuses = append(newPeerStatuses, status.FromAtlas(*atlasPeer, peer.ProviderName, vpc))
		}
	}
	return newPeerStatuses
}

func getAllExistedNetworkPeer(ctx context.Context, logger *zap.SugaredLogger, peerService mongodbatlas.PeersService, groupID string) ([]mongodbatlas.Peer, error) {
	var peersList []mongodbatlas.Peer
	listAWS, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{})
	if err != nil {
		logger.Errorf("failed to list network peers: %v", err)
		return nil, err
	}
	logger.Debugf("got %d aws peers", len(listAWS))
	peersList = append(peersList, listAWS...)

	listGCP, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{
		ProviderName: string(provider.ProviderGCP),
	})
	if err != nil {
		logger.Errorf("failed to list network peers: %v", err)
		return nil, err
	}
	logger.Debugf("got %d gcp peers", len(listGCP))
	peersList = append(peersList, listGCP...)

	listAzure, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{
		ProviderName: string(provider.ProviderAzure),
	})
	if err != nil {
		logger.Errorf("failed to list network peers: %v", err)
		return nil, err
	}
	logger.Debugf("got %d azure peers", len(listAzure))
	peersList = append(peersList, listAzure...)
	return peersList, nil
}

func sortPeers(existedPeers []mongodbatlas.Peer, expectedPeers []mdbv1.NetworkPeer, logger *zap.SugaredLogger, containerService mongodbatlas.ContainersService, groupID string) (*networkPeerDiff, error) {
	var diff networkPeerDiff
	var peersToUpdate []mdbv1.NetworkPeer
	for _, existedPeer := range existedPeers {
		logger.Debugf("existed peer %v", existedPeer)
		needToDelete := true
		for _, expectedPeer := range expectedPeers {
			logger.Debugf("expected peer %v", expectedPeer)
			if comparePeersPair(existedPeer, expectedPeer, containerService, groupID) {
				logger.Debugf("peer %v is equal to expected peer %v", existedPeer, expectedPeer)
				//expectedPeer.ID = existedPeer.ID
				existedPeer.ProviderName = string(expectedPeer.ProviderName)
				existedPeer.AccepterRegionName = expectedPeer.AccepterRegionName
				diff.PeersToUpdate = append(diff.PeersToUpdate, existedPeer)
				peersToUpdate = append(peersToUpdate, expectedPeer)
				needToDelete = false
			}
		}

		if needToDelete {
			if existedPeer.Status != StatusDeleting || existedPeer.StatusName != StatusDeleting {
				logger.Debugf("peer %v will be deleted", existedPeer)
				diff.PeersToDelete = append(diff.PeersToDelete, existedPeer.ID)
			}
		}
	}

	for _, expectedPeer := range expectedPeers {
		needToCreate := true
		for _, peerToUpdate := range peersToUpdate {
			opPeer, err := peerToUpdate.ToAtlas()
			if err != nil {
				return nil, err
			}
			if comparePeersPair(*opPeer, expectedPeer, containerService, groupID) {
				needToCreate = false
			}
		}
		if needToCreate {
			diff.PeersToCreate = append(diff.PeersToCreate, expectedPeer)
		}
	}
	return &diff, nil
}

func comparePeersPair(existedPeer mongodbatlas.Peer, expectedPeer mdbv1.NetworkPeer, containerService mongodbatlas.ContainersService, groupID string) bool {
	if expectedPeer.ProviderName == "" {
		expectedPeer.ProviderName = provider.ProviderAWS
	}

	if existedPeer.AWSAccountID != "" {
		existedPeer.ProviderName = string(provider.ProviderAWS)
	} else if existedPeer.AzureSubscriptionID != "" {
		existedPeer.ProviderName = string(provider.ProviderAzure)
	} else if existedPeer.GCPProjectID != "" {
		existedPeer.ProviderName = string(provider.ProviderGCP)
	}

	if expectedPeer.ContainerID != "" {
		if existedPeer.ContainerID != expectedPeer.ContainerID {
			return false
		}
	}

	if expectedPeer.AtlasCIDRBlock != "" {
		if existedPeer.AtlasCIDRBlock == "" {
			// existed peer doesn't contain AtlasCIDRBlock. so we have to get it by containerID
			get, _, err := containerService.Get(context.Background(), groupID, existedPeer.ContainerID)
			if err != nil {
				return false
			}
			existedPeer.AtlasCIDRBlock = get.AtlasCIDRBlock
		}
		if existedPeer.AtlasCIDRBlock != expectedPeer.AtlasCIDRBlock {
			return false
		}
	}

	switch expectedPeer.ProviderName {
	case provider.ProviderAWS:
		if existedPeer.VpcID == expectedPeer.VpcID &&
			expectedPeer.AWSAccountID == existedPeer.AWSAccountID &&
			expectedPeer.RouteTableCIDRBlock == existedPeer.RouteTableCIDRBlock {
			return true
		}
		return false
	case provider.ProviderGCP:
		if existedPeer.GCPProjectID == expectedPeer.GCPProjectID &&
			existedPeer.NetworkName == expectedPeer.NetworkName {
			return true
		}
		return false
	case provider.ProviderAzure:

		if existedPeer.AzureSubscriptionID == expectedPeer.AzureSubscriptionID &&
			existedPeer.AzureDirectoryID == expectedPeer.AzureDirectoryID &&
			existedPeer.ResourceGroupName == expectedPeer.ResourceGroupName &&
			existedPeer.VNetName == expectedPeer.VNetName {
			return true
		}
		return false
	default:
		return false
	}
}

func deletePeerByID(ctx context.Context, peerService mongodbatlas.PeersService, groupID string, peerID string, logger *zap.SugaredLogger) error {
	_, err := peerService.Delete(ctx, groupID, peerID)
	if err != nil {
		logger.Errorf("failed to delete peer %s: %v", peerID, err)
		return err
	}
	return nil
}

func containerRegionMatcher(regionName string, providerName provider.ProviderName) string {
	switch providerName {
	case provider.ProviderAWS:
		return awsRegionMatcher(regionName)
	case provider.ProviderGCP:
		return "" // GCP doesn't have region
	case provider.ProviderAzure:
		return "" // TODO: add Azure region matcher
	default:
		return ""
	}
}

func awsRegionMatcher(regionName string) string {
	result := strings.Replace(regionName, "-", "_", -1)
	return strings.ToUpper(result)
}

func createContainer(ctx context.Context, containerService mongodbatlas.ContainersService, groupID string, peer mdbv1.NetworkPeer, logger *zap.SugaredLogger) (string, error) {
	create, response, err := containerService.Create(ctx, groupID, &mongodbatlas.Container{
		AtlasCIDRBlock: peer.AtlasCIDRBlock,
		ProviderName:   string(peer.ProviderName),
		RegionName:     containerRegionMatcher(peer.GetContainerRegion(), peer.ProviderName),
	})
	if err != nil {
		if response.StatusCode == http.StatusConflict {
			list, _, errList := containerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{ProviderName: string(peer.ProviderName)})
			if errList != nil {
				logger.Errorf("failed to list containers: %v", errList)
				return "", errList
			}
			for _, container := range list {
				if container.AtlasCIDRBlock == peer.AtlasCIDRBlock &&
					container.RegionName == containerRegionMatcher(peer.GetContainerRegion(), peer.ProviderName) &&
					container.GCPProjectID == peer.GCPProjectID { //TODO: check if its work with azure
					return container.ID, nil
				}
			}
		}
		logger.Errorf("failed to create container: %v", err)
		return "", err
	}
	return create.ID, nil
}

func createNetworkPeer(ctx context.Context, groupID string, service mongodbatlas.PeersService, peer mdbv1.NetworkPeer, logger *zap.SugaredLogger) (*mongodbatlas.Peer, error) {
	peerToCreate, err := peer.ToAtlas()
	if err != nil {
		return nil, err
	}
	p, _, err := service.Create(ctx, groupID, peerToCreate)
	if err != nil {
		logger.Errorf("failed to create network peer %v: %v", peer, err)
		return peerToCreate, err
	}
	return p, nil
}

// validateInitNetworkPeer is validation according https://www.mongodb.com/docs/atlas/reference/api/vpc-create-peering-connection/
func validateInitNetworkPeer(peer mdbv1.NetworkPeer) error {
	if peer.ProviderName == "" {
		peer.ProviderName = provider.ProviderAWS
	}

	if peer.ContainerID == "" && peer.AtlasCIDRBlock == "" {
		return fmt.Errorf("containerID or AtlasCIDRBlock must be specified")
	}

	switch peer.ProviderName {
	case provider.ProviderAWS:
		if peer.AccepterRegionName == "" {
			return fmt.Errorf("accepterRegionName is required for AWS")
		}
		if peer.AWSAccountID == "" {
			return fmt.Errorf("awsAccountId is required for AWS")
		}
		if peer.RouteTableCIDRBlock == "" {
			return fmt.Errorf("routeTableCIDRBlock is required for AWS")
		}
		if peer.VpcID == "" {
			return fmt.Errorf("vpcId is required for AWS")
		}
		return nil
	case provider.ProviderGCP:
		if peer.GCPProjectID == "" {
			return fmt.Errorf("gcpProjectId is required for GCP")
		}
		if peer.AccepterRegionName == "" {
			return fmt.Errorf("accepterRegionName is required for GCP")
		}
		return nil
	case provider.ProviderAzure:
		if peer.AzureDirectoryID == "" {
			return fmt.Errorf("azureDirectoryId is required for Azure")
		}
		if peer.AzureSubscriptionID == "" {
			return fmt.Errorf("azureSubscriptionId is required for Azure")
		}
		if peer.ResourceGroupName == "" {
			return fmt.Errorf("resourceGroupName is required for Azure")
		}
		if peer.VNetName == "" {
			return fmt.Errorf("vNetName is required for Azure")
		}
	}
	return fmt.Errorf("unsupported provider %s", peer.ProviderName)
}

func DeleteAllNetworkPeers(ctx context.Context, groupID string, service mongodbatlas.PeersService, logger *zap.SugaredLogger) workflow.Result {
	result := workflow.OK()
	err := deleteAllNetworkPeers(ctx, groupID, service, logger)
	if err != nil {
		workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to delete NetworkPeers")
	}
	return result
}

func deleteAllNetworkPeers(ctx context.Context, groupID string, service mongodbatlas.PeersService, logger *zap.SugaredLogger) error {
	peers, err := getAllExistedNetworkPeer(ctx, logger, service, groupID)
	if err != nil {
		logger.Errorf("failed to list network peers for project %s: %v", groupID, err)
		return err
	}
	for _, peer := range peers {
		errDelete := deletePeerByID(ctx, service, groupID, peer.ID, logger)
		if errDelete != nil {
			logger.Errorf("failed to delete network peer %s: %v", peer.ID, errDelete)
			return errDelete
		}
	}
	return nil
}
