package model

import (
	"go.mongodb.org/atlas-sdk/v20231001002/admin"
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
