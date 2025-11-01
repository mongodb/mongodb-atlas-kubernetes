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

// CustomRole lets you create and change a custom role in your cluster.
// Use custom roles to specify custom sets of actions that the Atlas built-in roles can't describe.
// Deprecated: Migrate to the AtlasCustomRoles custom resource in accordance with the migration guide
// at https://www.mongodb.com/docs/atlas/operator/current/migrate-parameter-to-resource/#std-label-ak8so-migrate-ptr
type CustomRole struct {
	// Human-readable label that identifies the role. This name must be unique for this custom role in this project.
	Name string `json:"name"`
	// List of the built-in roles that this custom role inherits.
	// +optional
	InheritedRoles []Role `json:"inheritedRoles,omitempty"`
	// List of the individual privilege actions that the role grants.
	// +optional
	Actions []Action `json:"actions,omitempty"`
}

type Role struct {
	// Human-readable label that identifies the role inherited.
	Name string `json:"name"`
	// Human-readable label that identifies the database on which someone grants the action to one MongoDB user.
	Database string `json:"database"`
}

type Action struct {
	// Human-readable label that identifies the privilege action.
	Name string `json:"name"`
	// List of resources on which you grant the action.
	Resources []Resource `json:"resources"`
}

type Resource struct {
	// Flag that indicates whether to grant the action on the cluster resource. If true, MongoDB Cloud ignores Database and Collection parameters.
	Cluster *bool `json:"cluster,omitempty"`
	// Human-readable label that identifies the database on which you grant the action to one MongoDB user.
	Database *string `json:"database,omitempty"`
	// Human-readable label that identifies the collection on which you grant the action to one MongoDB user.
	Collection *string `json:"collection,omitempty"`
}
