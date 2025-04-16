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
	CustomZoneMapping     map[string]string `json:"customZoneMapping,omitempty"`
	ZoneMappingState      string            `json:"zoneMappingState,omitempty"`
	ZoneMappingErrMessage string            `json:"zoneMappingErrMessage,omitempty"`
}

type ManagedNamespace struct {
	Db                     string `json:"db"` // not changing this as is a breaking change
	Collection             string `json:"collection"`
	CustomShardKey         string `json:"customShardKey,omitempty"`
	NumInitialChunks       int    `json:"numInitialChunks,omitempty"`
	IsCustomShardKeyHashed *bool  `json:"isCustomShardKeyHashed,omitempty"` // Flag that specifies whether the custom shard key for the collection is hashed.
	IsShardKeyUnique       *bool  `json:"isShardKeyUnique,omitempty"`       // Flag that specifies whether the underlying index enforces a unique constraint.
	Status                 string `json:"status,omitempty"`
	PresplitHashedZones    *bool  `json:"presplitHashedZones,omitempty"`
	ErrMessage             string `json:"errMessage,omitempty"`
}
