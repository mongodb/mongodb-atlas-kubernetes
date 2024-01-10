package status

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
)

// AtlasDeploymentStatus defines the observed state of AtlasDeployment.
type AtlasDeploymentStatus struct {
	Common `json:",inline"`

	// StateName is the current state of the cluster.
	// The possible states are: IDLE, CREATING, UPDATING, DELETING, DELETED, REPAIRING
	StateName string `json:"stateName,omitempty"`

	// MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`

	// ConnectionStrings is a set of connection strings that your applications use to connect to this cluster.
	ConnectionStrings *ConnectionStrings `json:"connectionStrings,omitempty"`

	ReplicaSets []ReplicaSet `json:"replicaSets,omitempty"`

	ServerlessPrivateEndpoints []ServerlessPrivateEndpoint `json:"serverlessPrivateEndpoints,omitempty"`

	CustomZoneMapping *CustomZoneMapping `json:"customZoneMapping,omitempty"`

	ManagedNamespaces []ManagedNamespace `json:"managedNamespaces,omitempty"`

	// MongoURIUpdated is a timestamp in ISO 8601 date and time format in UTC when the connection string was last updated.
	// The connection string changes if you update any of the other values.
	MongoURIUpdated string `json:"mongoURIUpdated,omitempty"`
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
	ID       string `json:"id"`
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

func AtlasDeploymentConnectionStringsOption(connectionStrings *mongodbatlas.ConnectionStrings) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		cs := ConnectionStrings{}
		err := compat.JSONCopy(&cs, connectionStrings)
		if err != nil {
			return
		}
		s.ConnectionStrings = &cs
	}
}

func AtlasDeploymentMongoURIUpdatedOption(mongoURIUpdated string) AtlasDeploymentStatusOption {
	return func(s *AtlasDeploymentStatus) {
		s.MongoURIUpdated = mongoURIUpdated
	}
}
