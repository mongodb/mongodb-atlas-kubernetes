package datafederation

import (
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type DatafederationPrivateEndpointEntry struct {
	*akov2.DataFederationPE
	ProjectID string
}

func NewDatafederationPrivateEndpointEntry(pe *akov2.DataFederationPE, projectID string) *DatafederationPrivateEndpointEntry {
	if pe == nil {
		return nil
	}
	return &DatafederationPrivateEndpointEntry{DataFederationPE: pe, ProjectID: projectID}
}

func endpointsFromAtlas(endpoints []admin.PrivateNetworkEndpointIdEntry, projectID string) ([]*DatafederationPrivateEndpointEntry, error) {
	result := make([]*DatafederationPrivateEndpointEntry, 0, len(endpoints))
	for _, entry := range endpoints {
		result = append(result, endpointFromAtlas(&entry, projectID))
	}
	if err := cmp.Normalize(result); err != nil {
		return nil, fmt.Errorf("error normalizing data federation private endpoints: %w", err)
	}
	return result, nil
}

func endpointFromAtlas(endpoint *admin.PrivateNetworkEndpointIdEntry, projectID string) *DatafederationPrivateEndpointEntry {
	result := &DatafederationPrivateEndpointEntry{
		ProjectID: projectID,
	}
	if endpoint != nil {
		result.DataFederationPE = &akov2.DataFederationPE{
			EndpointID: endpoint.GetEndpointId(),
			Provider:   endpoint.GetProvider(),
			Type:       endpoint.GetType(),
		}
	}
	return result
}

func endpointToAtlas(ep *DatafederationPrivateEndpointEntry) *admin.PrivateNetworkEndpointIdEntry {
	if ep == nil || ep.DataFederationPE == nil {
		return nil
	}

	return &admin.PrivateNetworkEndpointIdEntry{
		EndpointId: ep.EndpointID,
		Provider:   pointer.MakePtrOrNil(ep.Provider),
		Type:       pointer.MakePtrOrNil(ep.Type),
	}
}
