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

package atlasproject

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/atlas-sdk/v20250312009/admin"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compare"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

const (
	StatusFailed      = "FAILED"
	StatusReady       = "AVAILABLE"
	StatusDeleting    = "DELETING"
	StatusTerminating = "TERMINATING"
)

var errNortFound = errors.New("not found")

type networkPeerDiff struct {
	PeersToDelete []string
	PeersToCreate []akov2.NetworkPeer
	PeersToUpdate []admin.BaseNetworkPeeringConnectionSettings
}

func lastAppliedNetworkPeerings(atlasProject *akov2.AtlasProject) ([]akov2.NetworkPeer, error) {
	lastApplied, err := lastAppliedSpecFrom(atlasProject)
	if err != nil {
		return nil, fmt.Errorf("failed to read project last applied configuration: %w", err)
	}
	if lastApplied == nil {
		return nil, nil
	}
	return lastApplied.NetworkPeers, nil
}

func ensureNetworkPeers(workflowCtx *workflow.Context, akoProject *akov2.AtlasProject) workflow.DeprecatedResult {
	lastAppliedPeers, err := lastAppliedNetworkPeerings(akoProject)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	networkPeerStatus := akoProject.Status.DeepCopy().NetworkPeers
	networkPeerSpec := akoProject.Spec.DeepCopy().NetworkPeers

	result, condition := SyncNetworkPeer(workflowCtx, akoProject.ID(), networkPeerStatus, networkPeerSpec, lastAppliedPeers)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(condition, result)
		return result
	}
	workflowCtx.SetConditionTrue(api.NetworkPeerReadyType)
	if len(networkPeerSpec) == 0 {
		workflowCtx.UnsetCondition(api.NetworkPeerReadyType)
	}

	return result
}

func failedPeerStatus(errMessage string, peer akov2.NetworkPeer) status.AtlasNetworkPeer {
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

func SyncNetworkPeer(workflowCtx *workflow.Context, groupID string, peerStatuses []status.AtlasNetworkPeer, peerSpecs []akov2.NetworkPeer, lastAppliedPeers []akov2.NetworkPeer) (workflow.DeprecatedResult, api.ConditionType) {
	defer workflowCtx.EnsureStatusOption(status.AtlasProjectSetNetworkPeerOption(&peerStatuses))
	logger := workflowCtx.Log
	mongoClient := workflowCtx.SdkClientSet.SdkClient20250312009
	logger.Debugf("syncing network peers for project %v", groupID)
	list, err := GetAllExistedNetworkPeer(workflowCtx.Context, mongoClient.NetworkPeeringApi, groupID)
	if err != nil {
		logger.Errorf("failed to get all network peers: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("failed to get all network peers")),
			api.NetworkPeerReadyType
	}

	diff := sortPeers(workflowCtx.Context, list, lastAppliedPeers, peerSpecs, logger, mongoClient.NetworkPeeringApi, groupID)
	logger.Debugf("peers to create %d, peers to update %d, peers to delete %d",
		len(diff.PeersToCreate), len(diff.PeersToUpdate), len(diff.PeersToDelete))

	for _, peerToDelete := range diff.PeersToDelete {
		errDelete := deletePeerByID(workflowCtx.Context, mongoClient.NetworkPeeringApi, groupID, peerToDelete, logger)
		if errDelete != nil {
			logger.Errorf("failed to delete network peer %s: %v", peerToDelete, errDelete)
			return workflow.Terminate(
					workflow.ProjectNetworkPeerIsNotReadyInAtlas,
					fmt.Errorf("failed to delete network peer: %w", errDelete),
				),
				api.NetworkPeerReadyType
		}
	}

	peerStatuses = createNetworkPeers(workflowCtx.Context, mongoClient, groupID, diff.PeersToCreate, logger)
	peerStatuses, err = UpdateStatuses(workflowCtx.Context, mongoClient.NetworkPeeringApi, peerStatuses, diff.PeersToUpdate, groupID, logger)
	if err != nil {
		logger.Errorf("failed to update network peer statuses: %v", err)
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas,
			errors.New("failed to update network peer statuses")), api.NetworkPeerReadyType
	}
	if len(lastAppliedPeers) > 0 {
		err = deleteUnusedContainers(workflowCtx.Context, mongoClient.NetworkPeeringApi, groupID, getPeerIDs(peerStatuses))
		if err != nil {
			logger.Errorf("failed to delete unused containers: %v", err)
			return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas,
				fmt.Errorf("failed to delete unused containers: %w", err)), api.NetworkPeerReadyType
		}
	}
	return ensurePeerStatus(peerStatuses, len(peerSpecs), logger), api.NetworkPeerReadyType
}

func UpdateStatuses(context context.Context, containerService admin.NetworkPeeringApi,
	peerStatuses []status.AtlasNetworkPeer, peersToUpdate []admin.BaseNetworkPeeringConnectionSettings, groupID string, logger *zap.SugaredLogger) ([]status.AtlasNetworkPeer, error) {
	for _, peerToUpdate := range peersToUpdate {
		vpc := formVPC(peerToUpdate)
		switch peerToUpdate.GetProviderName() {
		case string(provider.ProviderGCP), string(provider.ProviderAzure):
			container, errGet := getContainer(context, containerService, peerToUpdate, groupID, logger)
			if errGet != nil {
				return nil, errGet
			}
			peerStatuses = append(peerStatuses, status.NewNetworkPeerStatus(peerToUpdate,
				provider.ProviderName(peerToUpdate.GetProviderName()), vpc, container))
		default:
			peerStatuses = append(peerStatuses, status.NewNetworkPeerStatus(peerToUpdate,
				provider.ProviderName(peerToUpdate.GetProviderName()), vpc, admin.CloudProviderContainer{}))
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

func deleteUnusedContainers(context context.Context, containerService admin.NetworkPeeringApi, groupID string, doNotDelete []string) error {
	containers, _, err := containerService.ListGroupContainers(context, groupID).Execute()
	if err != nil {
		return err
	}
	for _, container := range containers.GetResults() {
		if container.GetProvisioned() { // a provisioned container is in use, should not be removed
			continue
		}
		if !compare.Contains(doNotDelete, container.GetId()) {
			response, errDelete := containerService.DeleteGroupContainer(context, groupID, container.GetId()).Execute()
			if errDelete != nil && response.StatusCode != http.StatusConflict { // AWS peer does not contain container id
				return errDelete
			}
		}
	}
	return nil
}

func getContainer(context context.Context, containerService admin.NetworkPeeringApi,
	peerToUpdate admin.BaseNetworkPeeringConnectionSettings, groupID string, logger *zap.SugaredLogger) (admin.CloudProviderContainer, error) {
	var container admin.CloudProviderContainer

	if peerToUpdate.GetContainerId() != "" {
		atlasContainer, _, err := containerService.GetGroupContainer(context, groupID, peerToUpdate.GetContainerId()).Execute()
		if err != nil {
			logger.Errorf("failed to get container for gcp status %s: %v", peerToUpdate.GetContainerId(), err)
			return container, fmt.Errorf("failed to get container for gcp status %s: %w", peerToUpdate.GetContainerId(), err)
		}
		if atlasContainer != nil {
			container = *atlasContainer
		}
	}

	return container, nil
}

func formVPC(peer admin.BaseNetworkPeeringConnectionSettings) string {
	switch peer.GetProviderName() {
	case string(provider.ProviderGCP):
		return peer.GetNetworkName()
	case string(provider.ProviderAzure):
		return peer.GetVnetName()
	default:
		return peer.GetVpcId()
	}
}

func ensurePeerStatus(peerStatuses []status.AtlasNetworkPeer, lenOfSpec int, logger *zap.SugaredLogger) workflow.DeprecatedResult {
	if len(peerStatuses) != lenOfSpec {
		return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("not all network peers are ready"))
	}

	for _, peerStatus := range peerStatuses {
		switch peerStatus.ProviderName {
		case provider.ProviderGCP:
			if peerStatus.Status != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("not all network peers are ready"))
			}
			if peerStatus.AtlasNetworkName == "" || peerStatus.AtlasGCPProjectID == "" { // We need this information to create the network peer connection
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("not all network peers are ready"))
			}
		case provider.ProviderAzure:
			if peerStatus.Status != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.Status)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("not all network peers are ready"))
			}
		default:
			if peerStatus.StatusName != StatusReady {
				logger.Debugf("network peer %s is not ready .%s.", peerStatus.VPC, peerStatus.StatusName)
				return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas, errors.New("not all network peers are ready"))
			}
		}
	}
	return workflow.OK()
}

func createNetworkPeers(context context.Context, mongoClient *admin.APIClient, groupID string, peers []akov2.NetworkPeer, logger *zap.SugaredLogger) []status.AtlasNetworkPeer {
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
			containerID, errCreate := createContainer(context, mongoClient.NetworkPeeringApi, groupID, peer, logger)
			if errCreate != nil {
				newPeerStatuses = append(newPeerStatuses,
					failedPeerStatus(fmt.Errorf("failed to create container for network peer %w", errCreate).Error(), peer))
				logger.Errorf("failed to create container for network peer: %s", errCreate)
				continue
			}
			peer.ContainerID = containerID
		}

		atlasPeer, err := createNetworkPeer(context, groupID, mongoClient.NetworkPeeringApi, peer, logger)
		if err != nil {
			logger.Errorf("failed to create network peer: %v", err)
			newPeerStatuses = append(newPeerStatuses,
				failedPeerStatus(fmt.Errorf("failed to create network peer: %w", err).Error(), peer))
			continue
		}
		if atlasPeer != nil {
			vpc := formVPC(*atlasPeer)
			atlasPeer.AccepterRegionName = pointer.SetOrNil(peer.AccepterRegionName, "")
			switch peer.ProviderName {
			case provider.ProviderGCP, provider.ProviderAzure:
				var container admin.CloudProviderContainer

				atlasContainer, _, err := mongoClient.NetworkPeeringApi.GetGroupContainer(context, groupID, peer.ContainerID).Execute()
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
					admin.CloudProviderContainer{}))
			}
		}
	}
	return newPeerStatuses
}

func GetAllExistedNetworkPeer(ctx context.Context, peerService admin.NetworkPeeringApi, groupID string) ([]admin.BaseNetworkPeeringConnectionSettings, error) {
	var peersList []admin.BaseNetworkPeeringConnectionSettings
	listAWS, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.BaseNetworkPeeringConnectionSettings], *http.Response, error) {
		return peerService.ListGroupPeersWithParams(ctx, &admin.ListGroupPeersApiParams{
			GroupId:      groupID,
			ProviderName: admin.PtrString(string(provider.ProviderAWS)),
		}).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for AWS: %w", err)
	}
	peersList = append(peersList, listAWS...)

	listGCP, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.BaseNetworkPeeringConnectionSettings], *http.Response, error) {
		return peerService.ListGroupPeersWithParams(ctx, &admin.ListGroupPeersApiParams{
			GroupId:      groupID,
			ProviderName: admin.PtrString(string(provider.ProviderGCP)),
		}).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for GCP: %w", err)
	}
	peersList = append(peersList, listGCP...)

	listAzure, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.BaseNetworkPeeringConnectionSettings], *http.Response, error) {
		return peerService.ListGroupPeersWithParams(ctx, &admin.ListGroupPeersApiParams{
			GroupId:      groupID,
			ProviderName: admin.PtrString(string(provider.ProviderAzure)),
		}).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list network peers for Azure: %w", err)
	}
	peersList = append(peersList, listAzure...)
	return peersList, nil
}

func sortPeers(ctx context.Context, existedPeers []admin.BaseNetworkPeeringConnectionSettings, lastApplied, expectedPeers []akov2.NetworkPeer, logger *zap.SugaredLogger, containerService admin.NetworkPeeringApi, groupID string) *networkPeerDiff {
	var diff networkPeerDiff
	var peersToUpdate []akov2.NetworkPeer
	for _, existedPeer := range existedPeers {
		needToDelete := true
		for _, expectedPeer := range expectedPeers {
			if comparePeersPair(ctx, *akov2.NewNetworkPeerFromAtlas(existedPeer), expectedPeer, containerService, groupID) {
				existedPeer.AccepterRegionName = pointer.SetOrNil(expectedPeer.AccepterRegionName, "")
				diff.PeersToUpdate = append(diff.PeersToUpdate, existedPeer)
				peersToUpdate = append(peersToUpdate, expectedPeer)
				needToDelete = false
			}
		}
		if !needToDelete || !isAtlasPeerManaged(ctx, lastApplied, existedPeer, containerService, groupID) {
			continue
		}

		logger.Debugf("peer %v will be deleted", existedPeer)
		if !isPeerDeleting(existedPeer) {
			logger.Debugf("peer %v will be deleted", existedPeer)
			diff.PeersToDelete = append(diff.PeersToDelete, existedPeer.GetId())
		}
	}

	for _, expectedPeer := range expectedPeers {
		needToCreate := true
		for _, peerToUpdate := range peersToUpdate {
			if comparePeersPair(ctx, peerToUpdate, expectedPeer, containerService, groupID) {
				needToCreate = false
			}
		}
		if needToCreate {
			diff.PeersToCreate = append(diff.PeersToCreate, expectedPeer)
		}
	}
	return &diff
}

func isAtlasPeerManaged(ctx context.Context, lastApplied []akov2.NetworkPeer, atlasPeer admin.BaseNetworkPeeringConnectionSettings, containerService admin.NetworkPeeringApi, groupID string) bool {
	for _, lastAppliedPeer := range lastApplied {
		if comparePeersPair(ctx, *akov2.NewNetworkPeerFromAtlas(atlasPeer), lastAppliedPeer, containerService, groupID) {
			return true
		}
	}
	return false
}

func isPeerDeleting(peer admin.BaseNetworkPeeringConnectionSettings) bool {
	return peer.GetStatus() == StatusDeleting || peer.GetStatusName() == StatusDeleting || peer.GetStatusName() == StatusTerminating
}

func comparePeersPair(ctx context.Context, existedPeer, expectedPeer akov2.NetworkPeer, containerService admin.NetworkPeeringApi, groupID string) bool {
	if expectedPeer.ProviderName == "" {
		expectedPeer.ProviderName = provider.ProviderAWS
	}

	if existedPeer.AWSAccountID != "" {
		existedPeer.ProviderName = provider.ProviderAWS
	} else if existedPeer.AzureSubscriptionID != "" {
		existedPeer.ProviderName = provider.ProviderAzure
	} else if existedPeer.GCPProjectID != "" {
		existedPeer.ProviderName = provider.ProviderGCP
	}

	if expectedPeer.ContainerID != "" {
		if existedPeer.ContainerID != expectedPeer.ContainerID {
			return false
		}
	}

	if expectedPeer.AtlasCIDRBlock != "" {
		if existedPeer.AtlasCIDRBlock == "" {
			get, _, err := containerService.GetGroupContainer(ctx, groupID, existedPeer.ContainerID).Execute()
			if err != nil {
				return false
			}
			existedPeer.AtlasCIDRBlock = get.GetAtlasCidrBlock()
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

func deletePeerByID(ctx context.Context, peerService admin.NetworkPeeringApi, groupID string, containerID string, logger *zap.SugaredLogger) error {
	_, response, err := peerService.DeleteGroupPeer(ctx, groupID, containerID).Execute()
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			return errors.Join(err, errNortFound)
		}
		logger.Errorf("failed to delete peering container %s: %v", containerID, err)
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

func createContainer(ctx context.Context, containerService admin.NetworkPeeringApi, groupID string, peer akov2.NetworkPeer, logger *zap.SugaredLogger) (string, error) {
	container := &admin.CloudProviderContainer{
		AtlasCidrBlock: pointer.SetOrNil(peer.AtlasCIDRBlock, ""),
		ProviderName:   pointer.SetOrNil(string(peer.ProviderName), ""),
	}
	if regionName := containerRegionNameMatcher(peer.GetContainerRegion(), peer.ProviderName); regionName != "" {
		container.SetRegionName(regionName)
	}
	if region := containerRegionMatcher(peer.GetContainerRegion(), peer.ProviderName); region != "" {
		container.SetRegion(region)
	}

	create, response, err := containerService.CreateGroupContainer(ctx, groupID, container).Execute()
	if err != nil {
		if response.StatusCode == http.StatusConflict {
			list, _, errList := containerService.ListGroupContainers(ctx, groupID).ProviderName(string(peer.ProviderName)).Execute()
			if errList != nil {
				logger.Errorf("failed to list containers: %v", errList)
				return "", errList
			}
			for _, container := range list.GetResults() {
				switch peer.ProviderName {
				case provider.ProviderAzure: // Azure network peer container use Region field to store region name
					if container.GetAtlasCidrBlock() == peer.AtlasCIDRBlock &&
						container.GetRegion() == peer.GetContainerRegion() {
						return container.GetId(), nil
					}
				case provider.ProviderAWS: // AWS network peer container use RegionName field to store region name.
					if container.GetAtlasCidrBlock() == peer.AtlasCIDRBlock &&
						container.GetRegionName() == containerRegionNameMatcher(peer.GetContainerRegion(), peer.ProviderName) {
						return container.GetId(), nil
					}
				case provider.ProviderGCP:
					if container.GetAtlasCidrBlock() == peer.AtlasCIDRBlock {
						return container.GetId(), nil
					}
				}
			}
		}
		logger.Errorf("failed to create container: %v", err)
		return "", err
	}
	return create.GetId(), nil
}

func createNetworkPeer(ctx context.Context, groupID string, service admin.NetworkPeeringApi, peer akov2.NetworkPeer, logger *zap.SugaredLogger) (*admin.BaseNetworkPeeringConnectionSettings, error) {
	p, _, err := service.CreateGroupPeer(ctx, groupID, peer.ToAtlasPeer()).Execute()
	if err != nil {
		logger.Errorf("failed to create network peer %v: %v", peer, err)
		return p, err
	}
	return p, nil
}

// validateInitNetworkPeer is validation according https://www.mongodb.com/docs/atlas/reference/api/vpc-create-peering-connection/
func validateInitNetworkPeer(peer akov2.NetworkPeer) error {
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

func DeleteOwnedNetworkPeers(ctx context.Context, project *akov2.AtlasProject, service admin.NetworkPeeringApi, logger *zap.SugaredLogger) workflow.DeprecatedResult {
	for _, peerStatus := range project.Status.NetworkPeers {
		errDelete := deletePeerByID(ctx, service, project.ID(), peerStatus.ID, logger)
		if errDelete != nil && !errors.Is(errDelete, errNortFound) {
			logger.Errorf("failed to delete network peer %s: %v", peerStatus.ID, errDelete)
			return workflow.Terminate(workflow.ProjectNetworkPeerIsNotReadyInAtlas,
				fmt.Errorf("failed to delete NetworkPeers: %w", errDelete))
		}
	}
	return workflow.OK()
}
