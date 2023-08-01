package project

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/atlas/mongodbatlas"
)

const (
	ProviderGCP   = "GCP"
	ProviderAWS   = "AWS"
	ProviderAzure = "AZURE"
)

func DeleteAllNetworkPeers(ctx context.Context, peerService mongodbatlas.PeersService, projectID string) error {
	peerList, err := GetAllNetworkPeers(ctx, peerService, projectID)
	if err != nil {
		return fmt.Errorf("error getting network peers: %w", err)
	}

	var allErr error
	for _, peer := range peerList {
		_, err = peerService.Delete(ctx, projectID, peer.ID)
		log.Printf("Deleting network peer %s", peer.ID)
		if err != nil {
			allErr = errors.Join(allErr, fmt.Errorf("error deleting network peer: %w", err))
			continue
		}
		log.Printf("Deleted network peer %s", peer.ID)
	}
	return allErr
}

func GetAllNetworkPeers(ctx context.Context, peerService mongodbatlas.PeersService, projectID string) ([]mongodbatlas.Peer, error) {
	var peersList []mongodbatlas.Peer
	listAWS, _, err := peerService.List(ctx, projectID, &mongodbatlas.ContainersListOptions{})
	if err != nil {
		return nil, err
	}
	peersList = append(peersList, listAWS...)

	listGCP, _, err := peerService.List(ctx, projectID, &mongodbatlas.ContainersListOptions{
		ProviderName: ProviderGCP,
	})
	if err != nil {
		return nil, err
	}
	peersList = append(peersList, listGCP...)

	listAzure, _, err := peerService.List(ctx, projectID, &mongodbatlas.ContainersListOptions{
		ProviderName: ProviderAzure,
	})
	if err != nil {
		return nil, err
	}
	peersList = append(peersList, listAzure...)
	return peersList, nil
}
