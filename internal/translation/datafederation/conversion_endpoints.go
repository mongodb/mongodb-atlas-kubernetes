package datafederation

import (
	"fmt"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
)

type PrivateEndpoint struct {
	*akov2.DataFederationPE
	ProjectID string
}

func NewPrivateEndpoints(projectID string, privateEndpoints []akov2.DataFederationPE) ([]PrivateEndpoint, error) {
	entries := make([]PrivateEndpoint, 0, len(privateEndpoints))
	for _, privateEndpoint := range privateEndpoints {
		entries = append(entries, NewPrivateEndpoint(projectID, privateEndpoint.DeepCopy()))
	}

	if err := cmp.Normalize(entries); err != nil {
		return nil, fmt.Errorf("error normalizing data federation private endpoints: %w", err)
	}

	return entries, nil
}

func NewPrivateEndpoint(projectID string, pe *akov2.DataFederationPE) PrivateEndpoint {
	if pe == nil {
		return PrivateEndpoint{}
	}

	return PrivateEndpoint{DataFederationPE: pe, ProjectID: projectID}
}

func endpointsFromAtlas(projectID string, endpoints []admin.PrivateNetworkEndpointIdEntry) ([]PrivateEndpoint, error) {
	entries := make([]PrivateEndpoint, 0, len(endpoints))
	for _, entry := range endpoints {
		entries = append(entries, endpointFromAtlas(projectID, &entry))
	}

	if err := cmp.Normalize(entries); err != nil {
		return nil, fmt.Errorf("error normalizing data federation private endpoints: %w", err)
	}

	return entries, nil
}

func endpointFromAtlas(projectID string, endpoint *admin.PrivateNetworkEndpointIdEntry) PrivateEndpoint {
	if endpoint == nil {
		return PrivateEndpoint{}
	}

	return PrivateEndpoint{
		ProjectID: projectID,
		DataFederationPE: &akov2.DataFederationPE{
			EndpointID: endpoint.GetEndpointId(),
			Provider:   endpoint.GetProvider(),
			Type:       endpoint.GetType(),
		},
	}
}

func endpointToAtlas(ep *PrivateEndpoint) *admin.PrivateNetworkEndpointIdEntry {
	if ep == nil || ep.DataFederationPE == nil {
		return nil
	}

	return &admin.PrivateNetworkEndpointIdEntry{
		EndpointId: ep.EndpointID,
		Provider:   pointer.MakePtr(ep.Provider),
		Type:       pointer.MakePtr(ep.Type),
	}
}
