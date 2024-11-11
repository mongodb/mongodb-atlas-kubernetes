package v1

import "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"

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

func NewFailedToCreateManagedNamespaceStatus(namespace ManagedNamespace, err error) status.ManagedNamespace {
	return status.ManagedNamespace{
		Db:                     namespace.Db,
		Collection:             namespace.Collection,
		CustomShardKey:         namespace.CustomShardKey,
		IsCustomShardKeyHashed: namespace.IsCustomShardKeyHashed,
		IsShardKeyUnique:       namespace.IsShardKeyUnique,
		NumInitialChunks:       namespace.NumInitialChunks,
		PresplitHashedZones:    namespace.PresplitHashedZones,
		Status:                 status.StatusFailed,
		ErrMessage:             err.Error(),
	}
}

func NewCreatedManagedNamespaceStatus(namespace ManagedNamespace) status.ManagedNamespace {
	return status.ManagedNamespace{
		Db:                     namespace.Db,
		Collection:             namespace.Collection,
		CustomShardKey:         namespace.CustomShardKey,
		IsCustomShardKeyHashed: namespace.IsCustomShardKeyHashed,
		IsShardKeyUnique:       namespace.IsShardKeyUnique,
		NumInitialChunks:       namespace.NumInitialChunks,
		PresplitHashedZones:    namespace.PresplitHashedZones,
		Status:                 status.StatusReady,
	}
}
