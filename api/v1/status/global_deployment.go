package status

type CustomZoneMapping struct {
	CustomZoneMapping     map[string]string `json:"customZoneMapping,omitempty"`
	ZoneMappingState      string            `json:"zoneMappingState,omitempty"`
	ZoneMappingErrMessage string            `json:"zoneMappingErrMessage,omitempty"`
}

type ManagedNamespace struct {
	Db                     string `json:"db"` //nolint:stylecheck // not changing this as is a breaking change
	Collection             string `json:"collection"`
	CustomShardKey         string `json:"customShardKey,omitempty"`
	NumInitialChunks       int    `json:"numInitialChunks,omitempty"`
	IsCustomShardKeyHashed *bool  `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	IsShardKeyUnique       *bool  `json:"isShardKeyUnique,omitempty"`       // Flag that specifies whether the underlying index enforces a unique constraint.
	Status                 string `json:"status,omitempty"`
	PresplitHashedZones    *bool  `json:"presplitHashedZones,omitempty"`
	ErrMessage             string `json:"errMessage,omitempty"`
}
