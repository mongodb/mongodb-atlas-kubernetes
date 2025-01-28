package networkcontainer

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

var (
	// ErrNotFound means an resource is missing
	ErrNotFound = errors.New("not found")

	// ErrContainerInUse is a failure to remove a containe still in use
	ErrContainerInUse = errors.New("container still in use")
)

type NetworkContainerService interface {
	Create(ctx context.Context, projectID string, container *NetworkContainer) (*NetworkContainer, error)
	Get(ctx context.Context, projectID, containerID string) (*NetworkContainer, error)
	Update(ctx context.Context, projectID string, container *NetworkContainer) (*NetworkContainer, error)
	Delete(ctx context.Context, projectID, containerID string) error
}

type networkContainerService struct {
	peeringAPI admin.NetworkPeeringApi
}

func NewNetworkPeeringServiceFromClientSet(clientSet *atlas.ClientSet) NetworkContainerService {
	return NewNetworkContainerService(clientSet.SdkClient20231115008.NetworkPeeringApi)
}

func NewNetworkContainerService(peeringAPI admin.NetworkPeeringApi) NetworkContainerService {
	return &networkContainerService{peeringAPI: peeringAPI}
}

func (np *networkContainerService) Create(ctx context.Context, projectID string, container *NetworkContainer) (*NetworkContainer, error) {
	newContainer, _, err := np.peeringAPI.CreatePeeringContainer(ctx, projectID, toAtlas(container)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create peering container %s: %w", container.ID, err)
	}
	return fromAtlas(newContainer), nil
}

func (np *networkContainerService) Get(ctx context.Context, projectID, containerID string) (*NetworkContainer, error) {
	container, _, err := np.peeringAPI.GetPeeringContainer(ctx, projectID, containerID).Execute()
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return nil, errors.Join(err, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get container %s: %w", containerID, err)
	}
	return fromAtlas(container), nil
}

func (np *networkContainerService) Update(ctx context.Context, projectID string, container *NetworkContainer) (*NetworkContainer, error) {
	updatedContainer, _, err := np.peeringAPI.UpdatePeeringContainer(ctx, projectID, container.ID, toAtlas(container)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update peering container %s: %w", container.ID, err)
	}
	return fromAtlas(updatedContainer), nil
}

func (np *networkContainerService) Delete(ctx context.Context, projectID, containerID string) error {
	_, _, err := np.peeringAPI.DeletePeeringContainer(ctx, projectID, containerID).Execute()
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return errors.Join(err, ErrNotFound)
	}
	if admin.IsErrorCode(err, "CONTAINERS_IN_USE") {
		return fmt.Errorf("failed to remove container %s as it is still in use: %w", containerID, ErrContainerInUse)
	}
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}
	return nil
}
