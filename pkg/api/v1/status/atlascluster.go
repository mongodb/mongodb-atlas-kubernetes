package status

import "go.mongodb.org/atlas/mongodbatlas"

// AtlasClusterStatus defines the observed state of AtlasCluster.
type AtlasClusterStatus struct {
	Common            `json:",inline"`
	StateName         string             `json:"stateName,omitempty"`
	MongoDBVersion    string             `json:"mongoDBVersion,omitempty"`
	ConnectionStrings *ConnectionStrings `json:"connectionStrings"`
	MongoURIUpdated   string             `json:"mongoURIUpdated,omitempty"`
}

// ConnectionStrings is a copy of mongodbatlas.ConnectionStrings for deepcopy compatibility purposes.
type ConnectionStrings struct {
	Standard          string            `json:"standard,omitempty"`
	StandardSrv       string            `json:"standardSrv,omitempty"`
	AwsPrivateLink    map[string]string `json:"awsPrivateLink,omitempty"`
	AwsPrivateLinkSrv map[string]string `json:"awsPrivateLinkSrv,omitempty"`
	Private           string            `json:"private,omitempty"`
	PrivateSrv        string            `json:"privateSrv,omitempty"`
}

// check that the two types are compatible.
var _ = ConnectionStrings(mongodbatlas.ConnectionStrings{})

// +k8s:deepcopy-gen=false

// AtlasClusterStatusOption is the option that is applied to Atlas Project Status.
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
