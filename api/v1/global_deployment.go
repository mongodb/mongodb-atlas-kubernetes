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

package v1

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

type CustomZoneMapping struct {
	// Code that represents a location that maps to a zone in your global cluster.
	// MongoDB Atlas represents this location with a ISO 3166-2 location and subdivision codes when possible.
	Location string `json:"location"`
	// Human-readable label that identifies the zone in your global cluster. This zone maps to a location code.
	Zone string `json:"zone"`
}

// ManagedNamespace represents the information about managed namespace configuration.
type ManagedNamespace struct {
	// Human-readable label of the database to manage for this Global Cluster.
	Db string `json:"db"` // not changing this as is a breaking change
	// Human-readable label of the collection to manage for this Global Cluster.
	Collection string `json:"collection"`
	// Database parameter used to divide the collection into shards. Global clusters require a compound shard key.
	// This compound shard key combines the location parameter and the user-selected custom key.
	CustomShardKey string `json:"customShardKey,omitempty"`
	// Minimum number of chunks to create initially when sharding an empty collection with a hashed shard key.
	// Maximum value is 8192.
	NumInitialChunks int `json:"numInitialChunks,omitempty"`
	// Flag that indicates whether MongoDB Cloud should create and distribute initial chunks for an empty or non-existing collection.
	// MongoDB Cloud distributes data based on the defined zones and zone ranges for the collection.
	PresplitHashedZones *bool `json:"presplitHashedZones,omitempty"`
	// Flag that indicates whether someone hashed the custom shard key for the specified collection.
	// If you set this value to false, MongoDB Cloud uses ranged sharding.
	IsCustomShardKeyHashed *bool `json:"isCustomShardKeyHashed,omitempty"`
	// Flag that indicates whether someone hashed the custom shard key.
	// If this parameter returns false, this cluster uses ranged sharding.
	IsShardKeyUnique *bool `json:"isShardKeyUnique,omitempty"`
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
