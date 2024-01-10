package atlasproject

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/util"
)

const (
	StatusFailed      = "FAILED"
	StatusReady       = "AVAILABLE"
	StatusDeleting    = "DELETING"
	StatusTerminating = "TERMINATING"
)

type networkPeerDiff struct {
	PeersToDelete []string
	PeersToCreate []mdbv1.NetworkPeer
	PeersToUpdate []mongodbatlas.Peer
}

func ensureNetworkPeers(workflowCtx *workflow.Context, akoProject *mdbv1.AtlasProject, subobjectProtect bool) workflow.Result {
	canReconcile, err := canNetworkPeeringReconcile(workflowCtx, subobjectProtect, akoProject)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.NetworkPeerReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Network Peering due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.NetworkPeerReadyType, result)

		return result
	}

	networkPeerStatus := akoProject.Status.DeepCopy().NetworkPeers
	networkPeerSpec := akoProject.Spec.DeepCopy().NetworkPeers

	result, condition := SyncNetworkPeer(workflowCtx, akoProject.ID(), networkPeerStatus, networkPeerSpec)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(condition, result)
		return result
	}
	workflowCtx.SetConditionTrue(status.NetworkPeerReadyType)
	if len(networkPeerSpec) == 0 {
		workflowCtx.UnsetCondition(status.NetworkPeerReadyType)
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
		errMessage = fmt.Sprintf("maybe its needed to setup Azure virtual network. error: %s", errMessage)
	default:
		vpc = peer.VpcID
	}
	return status.AtlasNetworkPeer{
		Status:       StatusFailed,
		ErrorMessage: errMessage,
		VPC:          vpc,
	}
}

func SyncNetworkPeer(workflowCtx *workflow.Context, groupID string, peerStatuses []status.AtlasNetworkPeer, peerSpecs []mdbv1.NetworkPeer) (workflow.Result, status.ConditionType) {
	defer workflowCtx.EnsureStatusOption(status.AtlasProjectSetNetworkPeerOption(&peerStatuses))
	logger := workflowCtx.Log
	mongoClient := workflowCtx.Client
	logger.Debugf("syncing network peers for project %v", groupID)
	list, err := GetAllExistedNetworkPeer(workflowCtx.Context, mongoClient.Peers, groupID)
	if err != nil {
		logger.Errorf("failed to get all network peers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to get all network peers"),
			status.NetworkPeerReadyType
	}

	diff, err := sortPeers(workflowCtx.Context, list, peerSpecs, logger, mongoClient.Containers, groupID)
	if err != nil {
		logger.Errorf("failed to sort network peers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to sort network peers"),
			status.NetworkPeerReadyType
	}
	logger.Debugf("peers to create %d, peers to update %d, peers to delete %d",
		len(diff.PeersToCreate), len(diff.PeersToUpdate), len(diff.PeersToDelete))

	for _, peerToDelete := range diff.PeersToDelete {
		errDelete := deletePeerByID(workflowCtx.Context, mongoClient.Peers, groupID, peerToDelete, logger)
		if errDelete != nil {
			logger.Errorf("failed to delete network peer %s: %v", peerToDelete, errDelete)
			return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to delete network peer"),
				status.NetworkPeerReadyType
		}
	}

	peerStatuses = createNetworkPeers(workflowCtx.Context, mongoClient, groupID, diff.PeersToCreate, logger)
	peerStatuses, err = UpdateStatuses(workflowCtx.Context, mongoClient.Containers, peerStatuses, diff.PeersToUpdate, groupID, logger)
	if err != nil {
		logger.Errorf("failed to update network peer statuses: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas,
			"failed to update network peer statuses"), status.NetworkPeerReadyType
	}
	err = deleteUnusedContainers(workflowCtx.Context, mongoClient.Containers, groupID, getPeerIDs(peerStatuses))
	if err != nil {
		logger.Errorf("failed to delete unused containers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas,
			fmt.Sprintf("failed to delete unused containers: %s", err)), status.NetworkPeerReadyType
	}
	return ensurePeerStatus(peerStatuses, len(peerSpecs), logger)
}

func UpdateStatuses(context context.Context, containerService mongodbatlas.ContainersService,
	peerStatuses []status.AtlasNetworkPeer, peersToUpdate []mongodbatlas.Peer, groupID string, logger *zap.SugaredLogger) ([]status.AtlasNetworkPeer, error) {
	for _, peerToUpdate := range peersToUpdate {
		vpc := formVPC(peerToUpdate)
		switch peerToUpdate.ProviderName {
		case string(provider.ProviderGCP), string(provider.ProviderAzure):
			container, errGet := getContainer(context, containerService, peerToUpdate, groupID, logger)
			if errGet != nil {
				return nil, errGet
			}
			peerStatuses = append(peerStatuses, status.NewNetworkPeerStatus(peerToUpdate,
				provider.ProviderName(peerToUpdate.ProviderName), vpc, container))
		default:
			peerStatuses = append(peerStatuses, status.NewNetworkPeerStatus(peerToUpdate,
				provider.ProviderName(peerToUpdate.ProviderName), vpc, mongodbatlas.Container{}))
		}
	}
	return peerStatuses, nil
}

func getPeerIDs(statuses []status.AtlasNetworkPeer) []string {
	ids := make([]string, 0, len(statuses))
	for _, networkPeer := range statuses {
		ids = append(ids, networkPeer.ContainerID)
	}
	return ids
}

func deleteUnusedContainers(context context.Context, containerService mongodbatlas.ContainersService, groupID string, doNotDelete []string) error {
	containers, _, err := containerService.List(context, groupID, nil)
	if err != nil {
		return err
	}
	for _, container := range containers {
		if !util.Contains(doNotDelete, container.ID) {
			response, errDelete := containerService.Delete(context, groupID, container.ID)
			if errDelete != nil && response.StatusCode != http.StatusConflict { // AWS peer does not contain container id
				return errDelete
			}
		}
	}
	return nil
}

func getContainer(context context.Context, containerService mongodbatlas.ContainersService,
	peerToUpdate mongodbatlas.Peer, groupID string, logger *zap.SugaredLogger) (mongodbatlas.Container, error) {
	var container mongodbatlas.Container

	if peerToUpdate.ContainerID != "" {
		atlasContainer, _, err := containerService.Get(context, groupID, peerToUpdate.ContainerID)
		if err != nil {
			logger.Errorf("failed to get container for gcp status %s: %v", peerToUpdate.ContainerID, err)
			return container, fmt.Errorf("failed to get container for gcp status %s: %w", peerToUpdate.ContainerID, err)
		}
		if atlasContainer != nil {
			container = *atlasContainer
		}
	} else if peerToUpdate.AtlasCIDRBlock != "" {
		list, _, err := containerService.List(context, groupID, &mongodbatlas.ContainersListOptions{
			ProviderName: string(provider.ProviderGCP),
		})
		if err != nil {
			logger.Errorf("failed to list containers for gcp status %v", err)
			return container, fmt.Errorf("failed to list containers for gcp status %w", err)
		}
		for _, cont := range list {
			if cont.AtlasCIDRBlock == peerToUpdate.AtlasCIDRBlock {
				return cont, nil
			}
		}
	}
	return container, nil
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

func ensurePeerStatus(peerStatuses []status.AtlasNetworkPeer, lenOfSpec int, logger *zap.SugaredLogger) (workflow.Result, status.ConditionType) {
	if len(peerStatuses) != lenOfSpec {
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
			status.NetworkPeerReadyType
	}

	for _, peerStatus := range peerStatuses {
		switch peerStatus.ProviderName {
		case provider.ProviderGCP:
			if peerStatus.Status != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
			if peerStatus.AtlasNetworkName == "" || peerStatus.AtlasGCPProjectID == "" { // We need this information to create the network peer connection
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
		case provider.ProviderAzure:
			if peerStatus.Status != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
		default:
			if peerStatus.StatusName != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.StatusName)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "not all network peers are ready"),
					status.NetworkPeerReadyType
			}
		}
	}
	return workflow.OK(), status.NetworkPeerReadyType
}

func createNetworkPeers(context context.Context, mongoClient *mongodbatlas.Client, groupID string, peers []mdbv1.NetworkPeer, logger *zap.SugaredLogger) []status.AtlasNetworkPeer {
	var newPeerStatuses []status.AtlasNetworkPeer
	for _, peer := range peers {
		err := validateInitNetworkPeer(peer)
		if err != nil {
			newPeerStatuses = append(newPeerStatuses,
				failedPeerStatus(fmt.Errorf("failed to validate network peer %w", err).Error(), peer))
			logger.Errorf("failed to validate network peer: %s", err)
			continue
		}
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
			switch peer.ProviderName {
			case provider.ProviderGCP, provider.ProviderAzure:
				var container mongodbatlas.Container

				atlasContainer, _, err := mongoClient.Containers.Get(context, groupID, peer.ContainerID)
				if err != nil {
					logger.Errorf("failed to get container for gcp status %s: %v", peer.ContainerID, err)
					newPeerStatuses = append(newPeerStatuses,
						failedPeerStatus(fmt.Errorf("failed to get container for gcp status %w", err).Error(), peer))
					continue
				}
				if atlasContainer != nil {
					container = *atlasContainer
				}

				newPeerStatuses = append(newPeerStatuses, status.NewNetworkPeerStatus(*atlasPeer, peer.ProviderName, vpc,
					container))
			default:
				newPeerStatuses = append(newPeerStatuses, status.NewNetworkPeerStatus(*atlasPeer, peer.ProviderName, vpc,
					mongodbatlas.Container{}))
			}
		}
	}
	return newPeerStatuses
}

func GetAllExistedNetworkPeer(ctx context.Context, peerService mongodbatlas.PeersService, groupID string) ([]mongodbatlas.Peer, error) {
	var peersList []mongodbatlas.Peer
	listAWS, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for AWS: %w", err)
	}
	peersList = append(peersList, listAWS...)

	listGCP, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{
		ProviderName: string(provider.ProviderGCP),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for GCP: %w", err)
	}
	peersList = append(peersList, listGCP...)

	listAzure, _, err := peerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{
		ProviderName: string(provider.ProviderAzure),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for Azure: %w", err)
	}
	peersList = append(peersList, listAzure...)
	return peersList, nil
}

func sortPeers(ctx context.Context, existedPeers []mongodbatlas.Peer, expectedPeers []mdbv1.NetworkPeer, logger *zap.SugaredLogger, containerService mongodbatlas.ContainersService, groupID string) (*networkPeerDiff, error) {
	var diff networkPeerDiff
	var peersToUpdate []mdbv1.NetworkPeer
	for _, existedPeer := range existedPeers {
		needToDelete := true
		for _, expectedPeer := range expectedPeers {
			if comparePeersPair(ctx, existedPeer, expectedPeer, containerService, groupID) {
				existedPeer.ProviderName = string(expectedPeer.ProviderName)
				existedPeer.AccepterRegionName = expectedPeer.AccepterRegionName
				existedPeer.ContainerID = expectedPeer.ContainerID
				existedPeer.AtlasCIDRBlock = expectedPeer.AtlasCIDRBlock
				diff.PeersToUpdate = append(diff.PeersToUpdate, existedPeer)
				peersToUpdate = append(peersToUpdate, expectedPeer)
				needToDelete = false
			}
		}

		if needToDelete {
			logger.Debugf("peer %v will be deleted", existedPeer)
			if !isPeerDeleting(existedPeer) {
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
			if comparePeersPair(ctx, *opPeer, expectedPeer, containerService, groupID) {
				needToCreate = false
			}
		}
		if needToCreate {
			diff.PeersToCreate = append(diff.PeersToCreate, expectedPeer)
		}
	}
	return &diff, nil
}

func isPeerDeleting(peer mongodbatlas.Peer) bool {
	return peer.Status == StatusDeleting || peer.StatusName == StatusDeleting || peer.StatusName == StatusTerminating
}

func comparePeersPair(ctx context.Context, existedPeer mongodbatlas.Peer, expectedPeer mdbv1.NetworkPeer, containerService mongodbatlas.ContainersService, groupID string) bool {
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
			get, _, err := containerService.Get(ctx, groupID, existedPeer.ContainerID)
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

// containerRegionNameMatcher is a matcher that matches a container's region name with a given region name. AWS only
func containerRegionNameMatcher(regionName string, providerName provider.ProviderName) string {
	switch providerName {
	case provider.ProviderAWS:
		return awsRegionMatcher(regionName)
	default:
		return ""
	}
}

// containerRegionMatcher is a matcher that matches a container's region name with a given region name. Azure only
func containerRegionMatcher(regionName string, providerName provider.ProviderName) string {
	switch providerName {
	case provider.ProviderAzure:
		return regionName
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
		RegionName:     containerRegionNameMatcher(peer.GetContainerRegion(), peer.ProviderName),
		Region:         containerRegionMatcher(peer.GetContainerRegion(), peer.ProviderName),
	})
	if err != nil {
		if response.StatusCode == http.StatusConflict {
			list, _, errList := containerService.List(ctx, groupID, &mongodbatlas.ContainersListOptions{ProviderName: string(peer.ProviderName)})
			if errList != nil {
				logger.Errorf("failed to list containers: %v", errList)
				return "", errList
			}
			for _, container := range list {
				switch peer.ProviderName {
				case provider.ProviderAzure: // Azure network peer container use Region field to store region name
					if container.AtlasCIDRBlock == peer.AtlasCIDRBlock &&
						container.Region == peer.GetContainerRegion() {
						return container.ID, nil
					}
				case provider.ProviderAWS: // AWS network peer container use RegionName field to store region name.
					if container.AtlasCIDRBlock == peer.AtlasCIDRBlock &&
						container.RegionName == containerRegionNameMatcher(peer.GetContainerRegion(), peer.ProviderName) {
						return container.ID, nil
					}
				case provider.ProviderGCP:
					if container.AtlasCIDRBlock == peer.AtlasCIDRBlock {
						return container.ID, nil
					}
				}
			}
		}
		logger.Errorf("failed to create container: %v", err)
		return "", err
	}
	return create.ID, nil
}

func createNetworkPeer(ctx context.Context, groupID string, service mongodbatlas.PeersService, peer mdbv1.NetworkPeer, logger *zap.SugaredLogger) (*mongodbatlas.Peer, error) {
	p, _, err := service.Create(ctx, groupID, peer.ToAtlasPeer())
	if err != nil {
		logger.Errorf("failed to create network peer %v: %v", peer, err)
		return p, err
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
		if peer.NetworkName == "" {
			return fmt.Errorf("networkName is required for GCP")
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
		return nil
	}
	return fmt.Errorf("unsupported provider: %s", peer.ProviderName)
}

func DeleteAllNetworkPeers(ctx context.Context, groupID string, service mongodbatlas.PeersService, logger *zap.SugaredLogger) workflow.Result {
	result := workflow.OK()
	err := deleteAllNetworkPeers(ctx, groupID, service, logger)
	if err != nil {
		result = workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, "failed to delete NetworkPeers")
	}
	return result
}

func deleteAllNetworkPeers(ctx context.Context, groupID string, service mongodbatlas.PeersService, logger *zap.SugaredLogger) error {
	peers, err := GetAllExistedNetworkPeer(ctx, service, groupID)
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

func canNetworkPeeringReconcile(workflowCtx *workflow.Context, protected bool, akoProject *mdbv1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &mdbv1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	containers, _, err := workflowCtx.Client.Containers.List(workflowCtx.Context, akoProject.ID(), &mongodbatlas.ContainersListOptions{})
	if err != nil {
		return false, err
	}

	if len(containers) > 0 && !areContainersEqual(latestConfig.NetworkPeers, containers) && !areContainersEqual(akoProject.Spec.NetworkPeers, containers) {
		return false, nil
	}

	peers, _, err := workflowCtx.Client.Peers.List(workflowCtx.Context, akoProject.ID(), &mongodbatlas.ContainersListOptions{})
	if err != nil {
		return false, err
	}

	if len(peers) == 0 {
		return true, nil
	}

	if !arePeersEqual(latestConfig.NetworkPeers, peers) && !arePeersEqual(akoProject.Spec.NetworkPeers, peers) {
		return false, nil
	}

	return true, err
}

func areContainersEqual(operatorContainers []mdbv1.NetworkPeer, atlasContainers []mongodbatlas.Container) bool {
	if len(operatorContainers) != len(atlasContainers) {
		return false
	}

	atlasContainersIDs := map[string]struct{}{}
	for _, container := range atlasContainers {
		switch container.ProviderName {
		case string(provider.ProviderAWS):
			atlasContainersIDs[fmt.Sprintf("%s.%s.%s", container.ProviderName, container.RegionName, container.AtlasCIDRBlock)] = struct{}{}
		case string(provider.ProviderGCP):
			atlasContainersIDs[fmt.Sprintf("%s.%s", container.ProviderName, container.AtlasCIDRBlock)] = struct{}{}
		case string(provider.ProviderAzure):
			atlasContainersIDs[fmt.Sprintf("%s.%s.%s", container.ProviderName, container.Region, container.AtlasCIDRBlock)] = struct{}{}
		}
	}

	for _, container := range operatorContainers {
		switch container.ProviderName {
		case provider.ProviderAWS:
			delete(atlasContainersIDs, fmt.Sprintf("%s.%s.%s", container.ProviderName, containerRegionNameMatcher(container.GetContainerRegion(), container.ProviderName), container.AtlasCIDRBlock))
		case provider.ProviderGCP:
			delete(atlasContainersIDs, fmt.Sprintf("%s.%s", container.ProviderName, container.AtlasCIDRBlock))
		case provider.ProviderAzure:
			delete(atlasContainersIDs, fmt.Sprintf("%s.%s.%s", container.ProviderName, containerRegionMatcher(container.GetContainerRegion(), container.ProviderName), container.AtlasCIDRBlock))
		}
	}

	return len(atlasContainersIDs) == 0
}

func arePeersEqual(operatorPeers []mdbv1.NetworkPeer, atlasPeers []mongodbatlas.Peer) bool {
	if len(operatorPeers) != len(atlasPeers) {
		return false
	}

	atlasPeersIDs := map[string]struct{}{}
	for _, peer := range atlasPeers {
		switch peer.ProviderName {
		case string(provider.ProviderAWS):
			atlasPeersIDs[fmt.Sprintf("%s.%s.%s", peer.AWSAccountID, peer.VpcID, peer.RouteTableCIDRBlock)] = struct{}{}
		case string(provider.ProviderGCP):
			atlasPeersIDs[fmt.Sprintf("%s.%s", peer.GCPProjectID, peer.NetworkName)] = struct{}{}
		case string(provider.ProviderAzure):
			atlasPeersIDs[fmt.Sprintf("%s.%s.%s.%s", peer.AzureSubscriptionID, peer.AzureDirectoryID, peer.ResourceGroupName, peer.VNetName)] = struct{}{}
		}
	}

	for _, peer := range operatorPeers {
		switch peer.ProviderName {
		case provider.ProviderAWS:
			delete(atlasPeersIDs, fmt.Sprintf("%s.%s.%s", peer.AWSAccountID, peer.VpcID, peer.RouteTableCIDRBlock))
		case provider.ProviderGCP:
			delete(atlasPeersIDs, fmt.Sprintf("%s.%s", peer.GCPProjectID, peer.NetworkName))
		case provider.ProviderAzure:
			delete(atlasPeersIDs, fmt.Sprintf("%s.%s.%s.%s", peer.AzureSubscriptionID, peer.AzureDirectoryID, peer.ResourceGroupName, peer.VNetName))
		}
	}

	return len(atlasPeersIDs) == 0
}
