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

package model

import (
	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

type AtlasRoles = string

const (
	GroupOwner AtlasRoles = "GROUP_OWNER"
)

type AtlasKeyType struct {
	DefaultFullAccessKey bool         // use full access key provided with github secrets
	Roles                []AtlasRoles // specify role for non default
	Whitelist            []string

	GlobalLevelKey    bool // if true, tests create "<operator-deployment-name>-api-key"
	GlobalKeyAttached *admin.ApiKeyUserDetails
}

func NewEmptyAtlasKeyType() *AtlasKeyType {
	return &AtlasKeyType{}
}

func (a *AtlasKeyType) GetRole() []AtlasRoles {
	return a.Roles
}

func (a *AtlasKeyType) UseDefaultFullAccess() *AtlasKeyType {
	a.DefaultFullAccessKey = true
	return a
}

func (a *AtlasKeyType) CreateAsGlobalLevelKey() *AtlasKeyType {
	a.GlobalLevelKey = true
	return a
}

func (a *AtlasKeyType) WithRoles(r []AtlasRoles) *AtlasKeyType {
	a.DefaultFullAccessKey = false
	a.Roles = r
	return a
}

func (a *AtlasKeyType) WithWhiteList(wl []string) *AtlasKeyType {
	a.DefaultFullAccessKey = false
	a.Whitelist = wl
	return a
}

func (a *AtlasKeyType) IsFullAccess() bool {
	return a.DefaultFullAccessKey
}
