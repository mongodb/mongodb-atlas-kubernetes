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

type CustomZoneMapping struct {
	// List that contains key value pairs to map zones to geographic regions.
	// These pairs map an ISO 3166-1a2 location code, with an ISO 3166-2 subdivision code when possible, to a unique 24-hexadecimal string that identifies the custom zone.
	CustomZoneMapping map[string]string `json:"customZoneMapping,omitempty"`
	// Status of the Custom Zone Mapping.
	ZoneMappingState string `json:"zoneMappingState,omitempty"`
	// Error message for failed Custom Zone Mapping.
	ZoneMappingErrMessage string `json:"zoneMappingErrMessage,omitempty"`
}

type ManagedNamespace struct {
	// Human-readable label of the database to manage for this Global Cluster.
	Db string `json:"db"` // not changing this as is a breaking change
	// Human-readable label of the collection to manage for this Global Cluster.
	Collection string `json:"collection"`
	// Database parameter used to divide the collection into shards. Global clusters require a compound shard key.
	// This compound shard key combines the location parameter and the user-selected custom key.
	CustomShardKey string `json:"customShardKey,omitempty"`
	// Minimum number of chunks to create initially when sharding an empty collection with a hashed shard key.
	NumInitialChunks int `json:"numInitialChunks,omitempty"`
	// Flag that indicates whether someone hashed the custom shard key for the specified collection.
	// If you set this value to false, MongoDB Atlas uses ranged sharding.
	IsCustomShardKeyHashed *bool `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	// Flag that indicates whether someone hashed the custom shard key. If this parameter returns false, this cluster uses ranged sharding.
	IsShardKeyUnique *bool `json:"isShardKeyUnique,omitempty"` // Flag that specifies whether the underlying index enforces a unique constraint.
	// Status of the Managed Namespace.
	Status string `json:"status,omitempty"`
	// Flag that indicates whether MongoDB Cloud should create and distribute initial chunks for an empty or non-existing collection.
	// MongoDB Atlas distributes data based on the defined zones and zone ranges for the collection.
	PresplitHashedZones *bool `json:"presplitHashedZones,omitempty"`
	// Error message for a failed Managed Namespace.
	ErrMessage string `json:"errMessage,omitempty"`
}
