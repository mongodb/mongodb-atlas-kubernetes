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

package networkcontainer

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

var (
	// ErrNotFound means an resource is missing
	ErrNotFound = errors.New("not found")

	// ErrContainerInUse is a failure to remove a containe still in use
	ErrContainerInUse = errors.New("container still in use")

	// ErrAmbigousFind fails when a find result is ambiguous,
	// usually more than one result was found when either one or noe was expected
	ErrAmbigousFind = errors.New("ambigous find results")
)

type NetworkContainerService interface {
	Create(ctx context.Context, projectID string, cfg *NetworkContainerConfig) (*NetworkContainer, error)
	Get(ctx context.Context, projectID, containerID string) (*NetworkContainer, error)
	Find(ctx context.Context, projectID string, cfg *NetworkContainerConfig) (*NetworkContainer, error)
	Update(ctx context.Context, projectID, containerID string, cfg *NetworkContainerConfig) (*NetworkContainer, error)
	Delete(ctx context.Context, projectID, containerID string) error
}

type networkContainerService struct {
	peeringAPI admin.NetworkPeeringApi
}

func NewNetworkContainerServiceFromClientSet(clientSet *atlas.ClientSet) NetworkContainerService {
	return NewNetworkContainerService(clientSet.SdkClient20250312012.NetworkPeeringApi)
}

func NewNetworkContainerService(peeringAPI admin.NetworkPeeringApi) NetworkContainerService {
	return &networkContainerService{peeringAPI: peeringAPI}
}

func (np *networkContainerService) Create(ctx context.Context, projectID string, cfg *NetworkContainerConfig) (*NetworkContainer, error) {
	newContainer, _, err := np.peeringAPI.CreateGroupContainer(ctx, projectID, toAtlasConfig(cfg)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create peering container at project %s: %w", projectID, err)
	}
	return fromAtlas(newContainer), nil
}

func (np *networkContainerService) Get(ctx context.Context, projectID, containerID string) (*NetworkContainer, error) {
	container, _, err := np.peeringAPI.GetGroupContainer(ctx, projectID, containerID).Execute()
	if admin.IsErrorCode(err, "CLOUD_PROVIDER_CONTAINER_NOT_FOUND") {
		return nil, errors.Join(err, ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get container %s: %w", containerID, err)
	}
	return fromAtlas(container), nil
}

func (np *networkContainerService) Find(ctx context.Context, projectID string, cfg *NetworkContainerConfig) (*NetworkContainer, error) {
	atlasContainers, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.CloudProviderContainer], *http.Response, error) {
		return np.peeringAPI.ListGroupContainers(ctx, projectID).ProviderName(cfg.Provider).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers at project %s: %w", projectID, err)
	}
	containers := []*NetworkContainer{}
	for _, atlasContainer := range atlasContainers {
		container := fromAtlas(&atlasContainer)
		switch provider.ProviderName(cfg.Provider) {
		case provider.ProviderGCP:
			if container.CIDRBlock == cfg.CIDRBlock {
				containers = append(containers, container)
			}
		default:
			if container.CIDRBlock == cfg.CIDRBlock && container.Region == cfg.Region {
				containers = append(containers, container)
			}
		}
	}
	if len(containers) < 1 {
		return nil, ErrNotFound
	}
	if len(containers) > 1 {
		return nil, ErrAmbigousFind
	}
	return containers[0], nil
}

func (np *networkContainerService) Update(ctx context.Context, projectID, containerID string, cfg *NetworkContainerConfig) (*NetworkContainer, error) {
	updatedContainer, _, err := np.peeringAPI.UpdateGroupContainer(ctx, projectID, containerID, toAtlasConfig(cfg)).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to update peering container %s: %w", containerID, err)
	}
	return fromAtlas(updatedContainer), nil
}

func (np *networkContainerService) Delete(ctx context.Context, projectID, containerID string) error {
	_, err := np.peeringAPI.DeleteGroupContainer(ctx, projectID, containerID).Execute()
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
