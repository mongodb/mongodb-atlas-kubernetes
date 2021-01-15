package status

import "go.mongodb.org/atlas/mongodbatlas"

// AtlasClusterStatus defines the observed state of AtlasCluster.
type AtlasClusterStatus struct {
	Common `json:",inline"`

	// StateName is the current state of the cluster.
	// The possible states are: IDLE, CREATING, UPDATING, DELETING, DELETED, REPAIRING
	StateName string `json:"stateName,omitempty"`

	// MongoDBVersion is the version of MongoDB the cluster runs, in <major version>.<minor version> format.
	MongoDBVersion string `json:"mongoDBVersion,omitempty"`

	// ConnectionStrings is a set of connection strings that your applications use to connect to this cluster.
	ConnectionStrings *ConnectionStrings `json:"connectionStrings,omitempty"`

	// MongoURIUpdated is a timestamp in ISO 8601 date and time format in UTC when the connection string was last updated.
	// The connection string changes if you update any of the other values.
	MongoURIUpdated string `json:"mongoURIUpdated,omitempty"`
}

// ConnectionStrings is a copy of mongodbatlas.ConnectionStrings for deepcopy compatibility purposes.
type ConnectionStrings struct {
	// Public mongodb:// connection string for this cluster.
	Standard string `json:"standard,omitempty"`

	// Public mongodb+srv:// connection string for this cluster.
	StandardSrv string `json:"standardSrv,omitempty"`

	// TODO: update when the replacement is implemented in mongodbatlas

	// Deprecated: use connectionStrings.privateEndpoint[n].connectionString instead.
	// Private-endpoint-aware mongodb://connection strings for each AWS PrivateLink private endpoint.
	// Atlas returns this parameter only if you deployed a AWS PrivateLink private endpoint to the same regions as all of this cluster's nodes.
	AwsPrivateLink map[string]string `json:"awsPrivateLink,omitempty"`

	// Deprecated: use connectionStrings.privateEndpoint[n].srvConnectionString instead.
	// Private-endpoint-aware mongodb+srv:// connection strings for each AWS PrivateLink private endpoint.
	AwsPrivateLinkSrv map[string]string `json:"awsPrivateLinkSrv,omitempty"`

	// Network-peering-endpoint-aware mongodb://connection strings for each interface VPC endpoint you configured to connect to this cluster.
	// Atlas returns this parameter only if you created a network peering connection to this cluster.
	Private string `json:"private,omitempty"`

	// Network-peering-endpoint-aware mongodb+srv:// connection strings for each interface VPC endpoint you configured to connect to this cluster.
	// Atlas returns this parameter only if you created a network peering connection to this cluster.
	// Use this URI format if your driver supports it. If it doesn't, use connectionStrings.private.
	PrivateSrv string `json:"privateSrv,omitempty"`
}

// Check compatibility with library type.
var _ = ConnectionStrings(mongodbatlas.ConnectionStrings{})

// +k8s:deepcopy-gen=false

// AtlasClusterStatusOption is the option that is applied to Atlas Cluster Status.
type AtlasClusterStatusOption func(s *AtlasClusterStatus)

func AtlasClusterStateNameOption(stateName string) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		s.StateName = stateName
	}
}

func AtlasClusterMongoDBVersionOption(mongoDBVersion string) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		s.MongoDBVersion = mongoDBVersion
	}
}

func AtlasClusterConnectionStringsOption(connectionStrings *mongodbatlas.ConnectionStrings) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		cs := ConnectionStrings(*connectionStrings)
		s.ConnectionStrings = &cs
	}
}

func AtlasClusterMongoURIUpdatedOption(mongoURIUpdated string) AtlasClusterStatusOption {
	return func(s *AtlasClusterStatus) {
		s.MongoURIUpdated = mongoURIUpdated
	}
}
