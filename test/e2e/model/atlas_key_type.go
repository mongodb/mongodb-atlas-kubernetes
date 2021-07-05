package model

type AtlasRoles = string

const (
	OrgOwner        AtlasRoles = "ORG_OWNER"
	OrgMember       AtlasRoles = "ORG_MEMBER"
	OrgGroupCreator AtlasRoles = "ORG_GROUP_CREATOR"
	OrgBillingAdmin AtlasRoles = "ORG_BILLING_ADMIN"
	OrgReadOnly     AtlasRoles = "ORG_READ_ONLY"

	GroupClusterManager      AtlasRoles = "GROUP_CLUSTER_MANAGER"
	GroupDataAccessAdmin     AtlasRoles = "GROUP_DATA_ACCESS_ADMIN"
	GropuDataAccessReadOnly  AtlasRoles = "GROUP_DATA_ACCESS_READ_ONLY"
	GroupDataAccessReadWrite AtlasRoles = "GROUP_DATA_ACCESS_READ_WRITE"
	GroupOwner               AtlasRoles = "GROUP_OWNER"
	GroupReadOnly            AtlasRoles = "GROUP_READ_ONLY"
)

type AtlasKeyType struct {
	GlobalLevelKey       bool         // if true, tests create "<operator-deployment-name>-api-key"
	DefaultFullAccessKey bool         // use full access key provided with github secrets
	Roles                []AtlasRoles // specify role for non default
	Whitelist            []string
}

func NewAtlasKeyType(r []AtlasRoles, wl []string) *AtlasKeyType {
	return &AtlasKeyType{
		DefaultFullAccessKey: false,
		Roles:                r,
		Whitelist:            wl,
	}
}

func NewEmptyAtlasKeyType() *AtlasKeyType {
	return &AtlasKeyType{}
}

func (a *AtlasKeyType) GetRole() []AtlasRoles {
	return a.Roles
}

func (a *AtlasKeyType) UseDefaulFullAccess() *AtlasKeyType {
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
