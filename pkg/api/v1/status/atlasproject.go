package status

import (
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/authmode"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

// +k8s:deepcopy-gen=false

// AtlasProjectStatusOption is the option that is applied to Atlas Project Status
type AtlasProjectStatusOption func(s *AtlasProjectStatus)

func AtlasProjectIDOption(id string) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.ID = id
	}
}

func AtlasProjectExpiredIPAccessOption(lists []project.IPAccessList) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.ExpiredIPAccessList = lists
	}
}

func AtlasProjectAddPrivateEndpointsOption(privateEndpoints []ProjectPrivateEndpoint) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.PrivateEndpoints = append(s.PrivateEndpoints, privateEndpoints...)
	}
}

func AtlasProjectSetPrivateEndpointsOption(privateEndpoints []ProjectPrivateEndpoint) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.PrivateEndpoints = privateEndpoints
	}
}

func AtlasProjectSetNetworkPeerOption(networkPeers *[]AtlasNetworkPeer) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.NetworkPeers = *networkPeers
	}
}

func AtlasProjectAuthModesOption(authModes []authmode.AuthMode) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.AuthModes = authModes
	}
}

func AtlasProjectSetAlertConfigOption(alertConfigs *[]AlertConfiguration) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.AlertConfigurations = *alertConfigs
	}
}

func AtlasProjectCloudIntegrationsOption(cloudIntegrations []CloudProviderIntegration) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.CloudProviderIntegrations = cloudIntegrations
	}
}

func AtlasProjectSetCustomRolesOption(customRoles *[]CustomRole) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.CustomRoles = *customRoles
	}
}

func AtlasProjectSetTeamsOption(teams *[]ProjectTeamStatus) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		if teams == nil {
			s.Teams = nil

			return
		}

		s.Teams = *teams
	}
}

func AtlasProjectPrometheusOption(prometheus *Prometheus) AtlasProjectStatusOption {
	return func(s *AtlasProjectStatus) {
		s.Prometheus = prometheus
	}
}

// AtlasProjectStatus defines the observed state of AtlasProject
type AtlasProjectStatus struct {
	api.Common `json:",inline"`

	// The ID of the Atlas Project
	// +optional
	ID string `json:"id,omitempty"`

	// The list of IP Access List entries that are expired due to 'deleteAfterDate' being less than the current date.
	// Note, that this field is updated by the Atlas Operator only after specification changes
	ExpiredIPAccessList []project.IPAccessList `json:"expiredIpAccessList,omitempty"`

	// The list of private endpoints configured for current project
	PrivateEndpoints []ProjectPrivateEndpoint `json:"privateEndpoints,omitempty"`

	// The list of network peers that are configured for current project
	NetworkPeers []AtlasNetworkPeer `json:"networkPeers,omitempty"`

	// AuthModes contains a list of configured authentication modes
	// "SCRAM" is default authentication method and requires a password for each user
	// "X509" signifies that self-managed X.509 authentication is configured
	AuthModes authmode.AuthModes `json:"authModes,omitempty"`

	// AlertConfigurations contains a list of alert configuration statuses
	AlertConfigurations []AlertConfiguration `json:"alertConfigurations,omitempty"`

	// CloudProviderIntegrations contains a list of configured cloud provider access roles. AWS support only
	CloudProviderIntegrations []CloudProviderIntegration `json:"cloudProviderIntegrations,omitempty"`

	// CustomRoles contains a list of custom roles statuses
	CustomRoles []CustomRole `json:"customRoles,omitempty"`

	// Teams contains a list of teams assignment statuses
	Teams []ProjectTeamStatus `json:"teams,omitempty"`

	// Prometheus contains the status for Prometheus integration
	// including the prometheusDiscoveryURL
	// +optional
	Prometheus *Prometheus `json:"prometheus,omitempty"`
}
