package globaldeployment

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"
)

// TODO: remove this after Atlas client will support all Global Cluster API

type AtlasGlobalDeployment struct {
	ManagedNamespaces []AtlasManagedNamespace `json:"managedNamespaces,omitempty"`
	CustomZoneMapping map[string]string       `json:"customZoneMapping"`
}

type AtlasManagedNamespace struct {
	DB                     string `json:"db"`
	Collection             string `json:"collection"`
	CustomShardKey         string `json:"customShardKey,omitempty"`
	NumInitialChunks       int    `json:"numInitialChunks,omitempty"`
	PresplitHashedZones    *bool  `json:"presplitHashedZones,omitempty"`
	IsCustomShardKeyHashed *bool  `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	IsShardKeyUnique       *bool  `json:"isShardKeyUnique,omitempty"`       // Flag that specifies whether the underlying index enforces a unique constraint.
}

const globalDeploymentBasePath = "api/atlas/v1.5/groups/%s/clusters/%s/globalWrites/%s"

func CreateManagedNamespaceInAtlas(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string, namespace *AtlasManagedNamespace) ([]AtlasManagedNamespace, error) {
	if namespace == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	path := fmt.Sprintf(globalDeploymentBasePath, groupID, deploymentName, "managedNamespaces")

	req, err := client.NewRequest(ctx, http.MethodPost, path, namespace)
	if err != nil {
		return nil, err
	}

	root := new(AtlasGlobalDeployment)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root.ManagedNamespaces, err
}

func DeleteManagedNamespaceInAtlas(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string, namespace *AtlasManagedNamespace) ([]AtlasManagedNamespace, error) {
	if namespace == nil {
		return nil, fmt.Errorf("namespace is nil")
	}

	path := fmt.Sprintf(globalDeploymentBasePath, groupID, deploymentName, "managedNamespaces")

	req, err := client.NewRequest(ctx, http.MethodDelete, path, namespace)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("collection", namespace.Collection)
	q.Add("db", namespace.DB)
	req.URL.RawQuery = q.Encode()

	root := new(AtlasGlobalDeployment)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root.ManagedNamespaces, err
}

func GetGlobalDeploymentState(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string) ([]AtlasManagedNamespace, map[string]string, error) {
	path := fmt.Sprintf("api/atlas/v1.5/groups/%s/clusters/%s/globalWrites", groupID, deploymentName)

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AtlasGlobalDeployment)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	return root.ManagedNamespaces, root.CustomZoneMapping, err
}

func CreateCustomZoneMapping(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string, customZoneMapping *mongodbatlas.CustomZoneMappingsRequest) (map[string]string, error) {
	if customZoneMapping == nil {
		return nil, fmt.Errorf("customZoneMapping is nil")
	}
	path := fmt.Sprintf(globalDeploymentBasePath, groupID, deploymentName, "customZoneMapping")

	req, err := client.NewRequest(ctx, http.MethodPost, path, customZoneMapping)
	if err != nil {
		return nil, err
	}

	root := new(AtlasGlobalDeployment)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return nil, err
	}

	return root.CustomZoneMapping, err
}

func DeleteAllCustomZoneMapping(ctx context.Context, client mongodbatlas.Client, groupID string, deploymentName string) error {
	path := fmt.Sprintf(globalDeploymentBasePath, groupID, deploymentName, "customZoneMapping")

	req, err := client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	root := new(AtlasGlobalDeployment)
	_, err = client.Do(ctx, req, root)
	if err != nil {
		return err
	}

	return nil
}
