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

package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
)

// AtlasDeploymentStatus defines the observed state of AtlasDeployment.
type AtlasDeploymentStatus struct {
	api.Common `json:",inline"`

	// StateName is the current state of the cluster.
	// The possible states are: IDLE, CREATING, UPDATING, DELETING, DELETED, REPAIRING
	StateName string `json:"stateName,omitempty"`

	// MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`

	// ConnectionStrings is a set of connection strings that your applications use to connect to this cluster.
	ConnectionStrings *ConnectionStrings `json:"connectionStrings,omitempty"`

	// Details that explain how MongoDB Cloud replicates data on the specified MongoDB database.
	// This array has one object per shard representing node configurations in each shard. For replica sets there is only one object representing node configurations.
	ReplicaSets []ReplicaSet `json:"replicaSets,omitempty"`

	// ServerlessPrivateEndpoints contains a list of private endpoints configured for the serverless deployment.
	ServerlessPrivateEndpoints []ServerlessPrivateEndpoint `json:"serverlessPrivateEndpoints,omitempty"`

	// List that contains key value pairs to map zones to geographic regions.
	// These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to a unique 24-hexadecimal string that identifies the custom zone.
	CustomZoneMapping *CustomZoneMapping `json:"customZoneMapping,omitempty"`

	// List that contains a namespace for a Global Cluster. MongoDB Atlas manages this cluster.
	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`

	// MongoURIUpdated is a timestamp in ISO 8601 date and time format in UTC when the connection string was last updated.
	// The connection string changes if you update any of the other values.
	MongoURIUpdated string `json:"mongoURIUpdated,omitempty"`

	// SearchIndexes contains a list of search indexes statuses configured for a project.
	SearchIndexes []DeploymentSearchIndexStatus `json:"searchIndexes,omitempty"`
}

const (
	StateIDLE      = "IDLE"
	StateCREATING  = "CREATING"
	StateUPDATING  = "UPDATING"
	StateDELETING  = "DELETING"
	StateDELETED   = "DELETED"
	StateREPAIRING = "REPAIRING"
)

type ReplicaSet struct {
	// Unique 24-hexadecimal digit string that identifies the replication object for a shard in a Cluster.
	ID string `json:"id"`
	// Human-readable label that describes the zone this shard belongs to in a Global Cluster.
	ZoneName string `json:"zoneName,omitempty"`
}

// ConnectionStrings contains configuration for applications use to connect to this cluster
type ConnectionStrings struct {
	// Public mongodb:// connection string for this cluster.
	Standard string `json:"standard,omitempty"`

	// Public mongodb+srv:// connection string for this cluster.
	StandardSrv string `json:"standardSrv,omitempty"`

	// Private endpoint connection strings.
	// Each object describes the connection strings you can use to connect to this cluster through a private endpoint.
	// Atlas returns this parameter only if you deployed a private endpoint to all regions to which you deployed this cluster's nodes.
	PrivateEndpoint []PrivateEndpoint `json:"privateEndpoint,omitempty"`

	// Network-peering-endpoint-aware mongodb:// connection strings for each interface VPC endpoint you configured to connect to this cluster.
	// Atlas returns this parameter only if you created a network peering connection to this cluster.
	Private string `json:"private,omitempty"`

	// Network-peering-endpoint-aware mongodb+srv:// connection strings for each interface VPC endpoint you configured to connect to this cluster.
	// Atlas returns this parameter only if you created a network peering connection to this cluster.
	// Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private.
	PrivateSrv string `json:"privateSrv,omitempty"`
}

// PrivateEndpoint connection strings. Each object describes the connection strings
// you can use to connect to this cluster through a private endpoint.
// Atlas returns this parameter only if you deployed a private endpoint to all regions
// to which you deployed this cluster's nodes.
type PrivateEndpoint struct {
	// Private-endpoint-aware mongodb:// connection string for this private endpoint.
	ConnectionString string `json:"connectionString,omitempty"`

	// Private endpoint through which you connect to Atlas when you use connectionStrings.privateEndpoint[n].connectionString or connectionStrings.privateEndpoint[n].srvConnectionString.
	Endpoints []Endpoint `json:"endpoints,omitempty"`

	// Private-endpoint-aware mongodb+srv:// connection string for this private endpoint.
	SRVConnectionString string `json:"srvConnectionString,omitempty"`

	// Private endpoint-aware connection string optimized for sharded clusters that uses the `mongodb+srv://` protocol to connect to MongoDB Cloud through a private endpoint.
	SRVShardOptimizedConnectionString string `json:"srvShardOptimizedConnectionString,omitempty"`

	// Type of MongoDB process that you connect to with the connection strings
	//
	// Atlas returns:
	//
	// • MONGOD for replica sets, or
	//
	// • MONGOS for sharded clusters
	Type string `json:"type,omitempty"`
}

// Endpoint through which you connect to Atlas
type Endpoint struct {
	// Unique identifier of the private endpoint.
	EndpointID string `json:"endpointId,omitempty"`

	// Cloud provider to which you deployed the private endpoint. Atlas returns AWS or AZURE.
	ProviderName string `json:"providerName,omitempty"`

	// Region to which you deployed the private endpoint.
	Region string `json:"region,omitempty"`

	// Private IP address of the private endpoint network interface you created in your Azure VNet.
	// +optional
	IP string `json:"ip,omitempty"`
}

// +k8s:deepcopy-gen=false

// AtlasDeploymentStatusOption is the option that is applied to Atlas Deployment Status.
type AtlasDeploymentStatusOption func(s *AtlasDeploymentStatus)

func AtlasDeploymentStateNameOption(stateName string) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.StateName = stateName
	}
}

func AtlasDeploymentReplicaSet(replicas []ReplicaSet) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.ReplicaSets = replicas
	}
}

func AtlasDeploymentSPEOption(pe []ServerlessPrivateEndpoint) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.ServerlessPrivateEndpoints = pe
	}
}

func AtlasDeploymentCustomZoneMappingOption(czm *CustomZoneMapping) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.CustomZoneMapping = czm
	}
}

func AtlasDeploymentManagedNamespacesOption(namespaces []ManagedNamespace) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.ManagedNamespaces = namespaces
	}
}

func AtlasDeploymentMongoDBVersionOption(mongoDBVersion string) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.MongoDBVersion = mongoDBVersion
	}
}

func AtlasDeploymentConnectionStringsOption(connStr *ConnectionStrings) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.ConnectionStrings = connStr
	}
}

func AtlasDeploymentRemoveStatusesWithEmptyIDs() AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		var result []DeploymentSearchIndexStatus
		for i := range s.SearchIndexes {
			if s.SearchIndexes[i].ID != "" {
				result = append(result, s.SearchIndexes[i])
			}
		}
		s.SearchIndexes = result
	}
}

// AtlasDeploymentSetSearchIndexStatus set the status for one SearchIndex
func AtlasDeploymentSetSearchIndexStatus(indexStatus DeploymentSearchIndexStatus) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		for i := range s.SearchIndexes {
			if s.SearchIndexes[i].Name == indexStatus.Name {
				s.SearchIndexes[i].Status = indexStatus.Status
				s.SearchIndexes[i].Message = indexStatus.Message
				if indexStatus.ID != "" {
					s.SearchIndexes[i].ID = indexStatus.ID
				}
				return
			}
		}
		s.SearchIndexes = append(s.SearchIndexes, indexStatus)
	}
}

// AtlasDeploymentUnsetSearchIndexStatus removes the status for one SearchIndex
func AtlasDeploymentUnsetSearchIndexStatus(indexStatus DeploymentSearchIndexStatus) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		for i := range s.SearchIndexes {
			if s.SearchIndexes[i].Name == indexStatus.Name {
				s.SearchIndexes = append(s.SearchIndexes[:i], s.SearchIndexes[i+1:]...)
				return
			}
		}
	}
}
