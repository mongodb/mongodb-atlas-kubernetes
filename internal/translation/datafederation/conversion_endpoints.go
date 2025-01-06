package datafederation

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/cmp"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
)

type DatafederationPrivateEndpointEntry struct {
	*akov2.DataFederationPE
	ProjectID string
}

func NewDataFederationPrivateEndpointEntry(projectID string, pe *akov2.DataFederationPE) *DatafederationPrivateEndpointEntry {
	if pe == nil {
		return nil
	}
	return &DatafederationPrivateEndpointEntry{DataFederationPE: pe, ProjectID: projectID}
}

func (e *DatafederationPrivateEndpointEntry) EqualsTo(target *DatafederationPrivateEndpointEntry) bool {
	return reflect.DeepEqual(e.DataFederationPE, target.DataFederationPE)
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

type DataFederationPrivateEndpoint struct {
	AKO, Atlas, LastApplied *DatafederationPrivateEndpointEntry
}

func MapDatafederationPrivateEndpoints(projectID string, df *akov2.AtlasDataFederation, atlasEndpoints []*DatafederationPrivateEndpointEntry) (map[string]*DataFederationPrivateEndpoint, error) {
	var lastApplied akov2.AtlasDataFederation
	if ann, ok := df.GetAnnotations()[customresource.AnnotationLastAppliedConfiguration]; ok {
		err := json.Unmarshal([]byte(ann), &lastApplied.Spec)
		if err != nil {
			return nil, fmt.Errorf("error reading data federation from last applied annotation: %w", err)
		}
	}

	result := make(map[string]*DataFederationPrivateEndpoint)
	addEndpoint := func(endpointID string) {
		if _, ok := result[endpointID]; !ok {
			result[endpointID] = &DataFederationPrivateEndpoint{}
		}
	}

	for _, pe := range atlasEndpoints {
		addEndpoint(pe.EndpointID)
		result[pe.EndpointID].Atlas = pe
	}
	for _, pe := range df.Spec.PrivateEndpoints {
		addEndpoint(pe.EndpointID)
		result[pe.EndpointID].AKO = NewDataFederationPrivateEndpointEntry(projectID, &pe)
	}
	for _, pe := range lastApplied.Spec.PrivateEndpoints {
		addEndpoint(pe.EndpointID)
		result[pe.EndpointID].LastApplied = NewDataFederationPrivateEndpointEntry(projectID, &pe)
	}

	return result, nil
}
