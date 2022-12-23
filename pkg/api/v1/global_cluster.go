package v1

import (
	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment/globaldeployment"
)

type CustomZoneMapping struct {
	Location string `json:"location"`
	Zone     string `json:"zone"`
}

// ManagedNamespace represents the information about managed namespace configuration.
type ManagedNamespace struct {
	Db                     string `json:"db"` //nolint:stylecheck // not changing this as is a breaking change
	Collection             string `json:"collection"`
	CustomShardKey         string `json:"customShardKey,omitempty"`
	NumInitialChunks       int    `json:"numInitialChunks,omitempty"`
	PresplitHashedZones    *bool  `json:"presplitHashedZones,omitempty"`
	IsCustomShardKeyHashed *bool  `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	IsShardKeyUnique       *bool  `json:"isShardKeyUnique,omitempty"`       // Flag that specifies whether the underlying index enforces a unique constraint.
}

func (in *ManagedNamespace) ToAtlas() globaldeployment.AtlasManagedNamespace {
	return globaldeployment.AtlasManagedNamespace{
		DB:                     in.Db,
		Collection:             in.Collection,
		CustomShardKey:         in.CustomShardKey,
		IsCustomShardKeyHashed: in.IsCustomShardKeyHashed,
		IsShardKeyUnique:       in.IsShardKeyUnique,
		NumInitialChunks:       in.NumInitialChunks,
		PresplitHashedZones:    in.PresplitHashedZones,
	}
}

func (c *CustomZoneMapping) ToAtlas() mongodbatlas.CustomZoneMapping {
	return mongodbatlas.CustomZoneMapping{
		Location: c.Location,
		Zone:     c.Zone,
	}
}
