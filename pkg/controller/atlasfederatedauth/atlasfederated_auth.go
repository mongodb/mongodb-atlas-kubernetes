package atlasfederatedauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/atlas-sdk/v20231115004/admin"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasFederatedAuthReconciler) ensureFederatedAuth(service *workflow.Context, fedauth *mdbv1.AtlasFederatedAuth) workflow.Result {
	// If disabled, skip with no error
	if !fedauth.Spec.Enabled {
		return workflow.OK().WithMessage(string(workflow.FederatedAuthIsNotEnabledInCR))
	}

	// Get current IDP for the ORG
	atlasFedSettings, _, err := service.SdkClient.FederatedAuthenticationApi.
		GetFederationSettings(service.Context, service.OrgID).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthNotAvailable, err.Error())
	}

	identityProvider, err := GetIdentityProviderForFederatedSettings(service.Context, service.SdkClient, atlasFedSettings)
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthNotAvailable, err.Error())
	}

	// Get current Org config
	orgConfig, _, err := service.SdkClient.FederatedAuthenticationApi.
		GetConnectedOrgConfig(service.Context, atlasFedSettings.GetId(), service.OrgID).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.FederatedAuthOrgNotConnected, err.Error())
	}

	projectList, err := prepareProjectList(service.Context, service.SdkClient)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintf("Can not list projects for org ID %s. %s", service.OrgID, err.Error()))
	}

	operatorConf, err := fedauth.Spec.ToAtlas(service.OrgID, identityProvider.GetOktaIdpId(), projectList)
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintln("Can not convert Federated Auth spec to Atlas", err.Error()))
	}

	if result := r.ensureIDPSettings(service.Context, atlasFedSettings.GetId(), identityProvider, fedauth, service.SdkClient); !result.IsOk() {
		return result
	}

	if federatedSettingsAreEqual(operatorConf, orgConfig) {
		return workflow.OK()
	}

	updatedSettings, _, err := service.SdkClient.FederatedAuthenticationApi.
		UpdateConnectedOrgConfig(service.Context, atlasFedSettings.GetId(), service.OrgID, operatorConf).
		Execute()
	if err != nil {
		return workflow.Terminate(workflow.Internal, fmt.Sprintln("Can not update federation settings", err.Error()))
	}

	if updatedSettings.UserConflicts != nil && len(*updatedSettings.UserConflicts) != 0 {
		users := make([]string, 0, len(*updatedSettings.UserConflicts))
		for i := range *updatedSettings.UserConflicts {
			users = append(users, (*updatedSettings.UserConflicts)[i].EmailAddress)
		}

		return workflow.Terminate(workflow.FederatedAuthUsersConflict,
			fmt.Sprintln("The following users are in conflict", users))
	}

	return workflow.OK()
}

func prepareProjectList(ctx context.Context, client *admin.APIClient) (map[string]string, error) {
	if client == nil {
		return nil, errors.New("client is not created")
	}

	projects, _, err := client.ProjectsApi.ListProjects(ctx).Execute()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(projects.GetResults()))
	for _, p := range projects.GetResults() {
		result[p.GetName()] = p.GetId()
	}

	return result, nil
}

func (r *AtlasFederatedAuthReconciler) ensureIDPSettings(ctx context.Context, federationSettingsID string, idp *admin.FederationIdentityProvider, fedauth *mdbv1.AtlasFederatedAuth, client *admin.APIClient) workflow.Result {
	if fedauth.Spec.SSODebugEnabled != nil {
		if idp.GetSsoDebugEnabled() == *fedauth.Spec.SSODebugEnabled {
			return workflow.OK()
		}

		idpUpdate := admin.IdentityProviderUpdate{
			DisplayName:     idp.DisplayName,
			IssuerUri:       idp.IssuerUri,
			SsoUrl:          idp.SsoUrl,
			SsoDebugEnabled: fedauth.Spec.SSODebugEnabled,
		}
		_, _, err := client.FederatedAuthenticationApi.UpdateIdentityProvider(ctx, federationSettingsID, idp.GetId(), &idpUpdate).Execute()
		if err != nil {
			return workflow.Terminate(workflow.Internal, err.Error())
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
