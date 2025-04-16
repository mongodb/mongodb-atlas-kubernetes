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

package project

import (
	"context"

	"go.mongodb.org/atlas/mongodbatlas"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

type Integration struct {
	// Third Party Integration type such as Slack, New Relic, etc
	// +kubebuilder:validation:Enum=PAGER_DUTY;SLACK;DATADOG;NEW_RELIC;OPS_GENIE;VICTOR_OPS;FLOWDOCK;WEBHOOK;MICROSOFT_TEAMS;PROMETHEUS
	// +optional
	Type string `json:"type,omitempty"`
	// +optional
	LicenseKeyRef common.ResourceRefNamespaced `json:"licenseKeyRef,omitempty"`
	// +optional
	AccountID string `json:"accountId,omitempty"`
	// +optional
	WriteTokenRef common.ResourceRefNamespaced `json:"writeTokenRef,omitempty"`
	// +optional
	ReadTokenRef common.ResourceRefNamespaced `json:"readTokenRef,omitempty"`
	// +optional
	APIKeyRef common.ResourceRefNamespaced `json:"apiKeyRef,omitempty"`
	// +optional
	Region string `json:"region,omitempty"`
	// +optional
	ServiceKeyRef common.ResourceRefNamespaced `json:"serviceKeyRef,omitempty"`
	// +optional
	APITokenRef common.ResourceRefNamespaced `json:"apiTokenRef,omitempty"`
	// +optional
	TeamName string `json:"teamName,omitempty"`
	// +optional
	ChannelName string `json:"channelName,omitempty"`
	// +optional
	RoutingKeyRef common.ResourceRefNamespaced `json:"routingKeyRef,omitempty"`
	// +optional
	FlowName string `json:"flowName,omitempty"`
	// +optional
	OrgName string `json:"orgName,omitempty"`
	// +optional
	URL string `json:"url,omitempty"`
	// +optional
	SecretRef common.ResourceRefNamespaced `json:"secretRef,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	MicrosoftTeamsWebhookURL string `json:"microsoftTeamsWebhookUrl,omitempty"`
	// +optional
	UserName string `json:"username,omitempty"`
	// +optional
	PasswordRef common.ResourceRefNamespaced `json:"passwordRef,omitempty"`
	// +optional
	ServiceDiscovery string `json:"serviceDiscovery,omitempty"`
	// +optional
	Scheme string `json:"scheme,omitempty"`
	// +optional
	Enabled bool `json:"enabled,omitempty"`
}

func (i Integration) ToAtlas(ctx context.Context, c client.Client, defaultNS string) (result *mongodbatlas.ThirdPartyIntegration, err error) {
	result = &mongodbatlas.ThirdPartyIntegration{
		Type:                     i.Type,
		AccountID:                i.AccountID,
		Region:                   i.Region,
		TeamName:                 i.TeamName,
		ChannelName:              i.ChannelName,
		FlowName:                 i.FlowName,
		OrgName:                  i.OrgName,
		URL:                      i.URL,
		Name:                     i.Name,
		MicrosoftTeamsWebhookURL: i.MicrosoftTeamsWebhookURL,
		UserName:                 i.UserName,
		ServiceDiscovery:         i.ServiceDiscovery,
		Scheme:                   i.Scheme,
		Enabled:                  i.Enabled,
	}

	readPassword := func(passwordField common.ResourceRefNamespaced, target *string, errors *[]error) {
		if passwordField.Name == "" {
			return
		}

		*target, err = passwordField.ReadPassword(ctx, c, defaultNS)
		storeError(err, errors)
	}

	errorList := make([]error, 0)
	readPassword(i.LicenseKeyRef, &result.LicenseKey, &errorList)
	readPassword(i.WriteTokenRef, &result.WriteToken, &errorList)
	readPassword(i.ReadTokenRef, &result.ReadToken, &errorList)
	readPassword(i.APIKeyRef, &result.APIKey, &errorList)
	readPassword(i.ServiceKeyRef, &result.ServiceKey, &errorList)
	readPassword(i.APITokenRef, &result.APIToken, &errorList)
	readPassword(i.RoutingKeyRef, &result.RoutingKey, &errorList)
	readPassword(i.SecretRef, &result.Secret, &errorList)
	readPassword(i.PasswordRef, &result.Password, &errorList)

	if len(errorList) != 0 {
		firstError := (errorList)[0]
		return nil, firstError
	}

	return result, nil
}

func (i Integration) Identifier() interface{} {
	return i.Type
}

func storeError(err error, errors *[]error) {
	if err != nil {
		*errors = append(*errors, err)
	}
}
