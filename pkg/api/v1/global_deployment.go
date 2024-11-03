package v1

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
