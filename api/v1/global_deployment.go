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
	Location string `json:"location"`
	Zone     string `json:"zone"`
}

// ManagedNamespace represents the information about managed namespace configuration.
type ManagedNamespace struct {
	Db                     string `json:"db"` // not changing this as is a breaking change
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
