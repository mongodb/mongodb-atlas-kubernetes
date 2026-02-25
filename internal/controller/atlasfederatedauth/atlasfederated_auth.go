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

package atlasfederatedauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas-sdk/v20250312014/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

func (r *AtlasFederatedAuthReconciler) ensureFederatedAuth(service *workflow.Context, fedauth *akov2.AtlasFederatedAuth) workflow.DeprecatedResult {
	// If disabled, skip with no error
	if !fedauth.Spec.Enabled {
		return workflow.OK().WithMessage(string(workflow.FederatedAuthIsNotEnabledInCR))
	}

	// Get current IDP for the ORG
	atlasFedSettings, _, err := service.SdkClientSet.SdkClient20250312013.FederatedAuthenticationApi.
		GetFederationSettings(service.Context, service.OrgID).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthNotAvailable, err)
	}

	identityProvider, err := GetIdentityProviderForFederatedSettings(service.Context, service.SdkClientSet.SdkClient20250312013, atlasFedSettings)
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthNotAvailable, err)
	}

	// Get current Org config
	orgConfig, _, err := service.SdkClientSet.SdkClient20250312013.FederatedAuthenticationApi.
		GetConnectedOrgConfig(service.Context, atlasFedSettings.GetId(), service.OrgID).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthOrgNotConnected, err)
	}

	projectList, err := prepareProjectList(service.Context, service.SdkClientSet.SdkClient20250312013)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Errorf("cannot list projects for org ID %s: %w", service.OrgID, err))
	}

	operatorConf, err := fedauth.Spec.ToAtlas(service.OrgID, identityProvider.GetOktaIdpId(), projectList)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Errorf("cannot convert Federated Auth spec to Atlas: %w", err))
	}

	if result := r.ensureIDPSettings(service.Context, atlasFedSettings.GetId(), identityProvider, fedauth, service.SdkClientSet.SdkClient20250312013); !result.IsOk() {
		return result
	}

	if federatedSettingsAreEqual(operatorConf, orgConfig) {
		return workflow.OK()
	}

	updatedSettings, _, err := service.SdkClientSet.SdkClient20250312013.FederatedAuthenticationApi.
		UpdateConnectedOrgConfig(service.Context, atlasFedSettings.GetId(), service.OrgID, operatorConf).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Errorf("cannot update federation settings: %w", err))
	}

	if updatedSettings.UserConflicts != nil && len(*updatedSettings.UserConflicts) != 0 {
		users := make([]string, 0, len(*updatedSettings.UserConflicts))
		for i := range *updatedSettings.UserConflicts {
			users = append(users, (*updatedSettings.UserConflicts)[i].EmailAddress)
		}

		return workflow.Terminate(workflow.FederatedAuthUsersConflict,
			fmt.Errorf("the following users are in conflict: %v", users))
	}

	return workflow.OK()
}

func prepareProjectList(ctx context.Context, client *admin.APIClient) (map[string]string, error) {
	if client == nil {
		return nil, errors.New("client is not created")
	}

	projects, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.Group], *http.Response, error) {
		return client.ProjectsApi.ListGroups(ctx).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(projects))
	for _, p := range projects {
		result[p.GetName()] = p.GetId()
	}

	return result, nil
}

func (r *AtlasFederatedAuthReconciler) ensureIDPSettings(ctx context.Context, federationSettingsID string, idp *admin.FederationIdentityProvider, fedauth *akov2.AtlasFederatedAuth, client *admin.APIClient) workflow.DeprecatedResult {
	if fedauth.Spec.SSODebugEnabled != nil {
		if idp.GetSsoDebugEnabled() == *fedauth.Spec.SSODebugEnabled {
			return workflow.OK()
		}

		idpUpdate := admin.FederationIdentityProviderUpdate{
			DisplayName:     idp.DisplayName,
			IssuerUri:       idp.IssuerUri,
			SsoUrl:          idp.SsoUrl,
			SsoDebugEnabled: fedauth.Spec.SSODebugEnabled,
		}
		_, _, err := client.FederatedAuthenticationApi.UpdateIdentityProvider(ctx, federationSettingsID, idp.GetId(), &idpUpdate).Execute()
		if err != nil {
			return workflow.Terminate(workflow.Internal, err)
		}
	}

	// TODO: Add more IDP settings
	return workflow.OK()
}

func federatedSettingsAreEqual(operator, atlas *admin.ConnectedOrgConfig) bool {
	operator.UserConflicts = nil
	atlas.UserConflicts = nil
	return cmp.Diff(operator, atlas) == ""
}

func GetIdentityProviderForFederatedSettings(ctx context.Context, atlasClient *admin.APIClient, fedSettings *admin.OrgFederationSettings) (*admin.FederationIdentityProvider, error) {
	identityProviders, _, err := atlasClient.FederatedAuthenticationApi.ListIdentityProviders(ctx, fedSettings.GetId()).Execute()
	if err != nil {
		return nil, err
	}

	for _, identityProvider := range identityProviders.GetResults() {
		if identityProvider.GetOktaIdpId() == fedSettings.GetIdentityProviderId() {
			return &identityProvider, nil
		}
	}

	return nil, fmt.Errorf("identity provider for Org Federation Settings %s not found", fedSettings.GetId())
}
