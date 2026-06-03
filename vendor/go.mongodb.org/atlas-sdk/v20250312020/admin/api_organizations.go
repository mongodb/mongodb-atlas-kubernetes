// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type OrganizationsApi interface {

	/*
		CreateOrg Create One Organization

		Creates one organization in MongoDB Cloud and links it to the requesting Service Account's or API Key's organization. The requesting Service Account's or API Key's organization must be a paying organization. To learn more, see Configure a Paying Organization in the MongoDB Atlas documentation. Optionally, if `federationSettingsId` is provided, the new Organization will be linked to the federation. The requesting Service Account or API Key must be an Organization Owner in the federation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param createOrganizationRequest Organization that you want to create.
		@return CreateOrgApiRequest
	*/
	CreateOrg(ctx context.Context, createOrganizationRequest *CreateOrganizationRequest) CreateOrgApiRequest
	/*
		CreateOrg Create One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgApiParams - Parameters for the request
		@return CreateOrgApiRequest
	*/
	CreateOrgWithParams(ctx context.Context, args *CreateOrgApiParams) CreateOrgApiRequest

	// Method available only for mocking purposes
	CreateOrgExecute(r CreateOrgApiRequest) (*CreateOrganizationResponse, *http.Response, error)

	/*
			CreateOrgInvite Create Invitation for One MongoDB Cloud User in One Organization

			Invites one MongoDB Cloud user to join the specified organization. The user must accept the invitation to access information within the specified organization.

		**Note**: Invitation management APIs are deprecated. Use Add One MongoDB Cloud User to One Organization to invite a user.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param organizationInvitationRequest Invites one MongoDB Cloud user to join the specified organization.
			@return CreateOrgInviteApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	CreateOrgInvite(ctx context.Context, orgId string, organizationInvitationRequest *OrganizationInvitationRequest) CreateOrgInviteApiRequest
	/*
		CreateOrgInvite Create Invitation for One MongoDB Cloud User in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateOrgInviteApiParams - Parameters for the request
		@return CreateOrgInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	CreateOrgInviteWithParams(ctx context.Context, args *CreateOrgInviteApiParams) CreateOrgInviteApiRequest

	// Method available only for mocking purposes
	CreateOrgInviteExecute(r CreateOrgInviteApiRequest) (*OrganizationInvitation, *http.Response, error)

	/*
			DeleteOrg Remove One Organization

			Removes one specified organization. MongoDB Cloud imposes the following limits on this resource:

		 - Organizations with active projects cannot be removed.
		 - All projects in the organization must be removed before you can remove the organization.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return DeleteOrgApiRequest
	*/
	DeleteOrg(ctx context.Context, orgId string) DeleteOrgApiRequest
	/*
		DeleteOrg Remove One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgApiParams - Parameters for the request
		@return DeleteOrgApiRequest
	*/
	DeleteOrgWithParams(ctx context.Context, args *DeleteOrgApiParams) DeleteOrgApiRequest

	// Method available only for mocking purposes
	DeleteOrgExecute(r DeleteOrgApiRequest) (*http.Response, error)

	/*
			DeleteOrgInvite Remove One Invitation from One Organization

			Cancels one pending invitation sent to the specified MongoDB Cloud user to join an organization. You can't cancel an invitation that the user accepted.

		**Note**: Invitation management APIs are deprecated. Use Remove One MongoDB Cloud User From One Organization to remove a pending user.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
			@return DeleteOrgInviteApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	DeleteOrgInvite(ctx context.Context, orgId string, invitationId string) DeleteOrgInviteApiRequest
	/*
		DeleteOrgInvite Remove One Invitation from One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteOrgInviteApiParams - Parameters for the request
		@return DeleteOrgInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	DeleteOrgInviteWithParams(ctx context.Context, args *DeleteOrgInviteApiParams) DeleteOrgInviteApiRequest

	// Method available only for mocking purposes
	DeleteOrgInviteExecute(r DeleteOrgInviteApiRequest) (*http.Response, error)

	/*
		GetOrg Return One Organization

		Returns one organization to which the requesting Service Account or API Key has access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return GetOrgApiRequest
	*/
	GetOrg(ctx context.Context, orgId string) GetOrgApiRequest
	/*
		GetOrg Return One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgApiParams - Parameters for the request
		@return GetOrgApiRequest
	*/
	GetOrgWithParams(ctx context.Context, args *GetOrgApiParams) GetOrgApiRequest

	// Method available only for mocking purposes
	GetOrgExecute(r GetOrgApiRequest) (*AtlasOrganization, *http.Response, error)

	/*
			GetOrgGroups Return All Projects in One Organization

			Returns multiple projects in the specified organization. Each organization can have multiple projects. Use projects to:

		- Isolate different environments, such as development, test, or production environments, from each other.
		- Associate different MongoDB Cloud users or teams with different environments, or give different permission to MongoDB Cloud users in different environments.
		- Maintain separate cluster security configurations.
		- Create different alert settings.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return GetOrgGroupsApiRequest
	*/
	GetOrgGroups(ctx context.Context, orgId string) GetOrgGroupsApiRequest
	/*
		GetOrgGroups Return All Projects in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgGroupsApiParams - Parameters for the request
		@return GetOrgGroupsApiRequest
	*/
	GetOrgGroupsWithParams(ctx context.Context, args *GetOrgGroupsApiParams) GetOrgGroupsApiRequest

	// Method available only for mocking purposes
	GetOrgGroupsExecute(r GetOrgGroupsApiRequest) (*PaginatedAtlasGroup, *http.Response, error)

	/*
			GetOrgInvite Return One Invitation in One Organization by Invitation ID

			Returns the details of one pending invitation to the specified organization.

		**Note**: Invitation management APIs are deprecated. Use Return One MongoDB Cloud User in One Organization to return a pending user.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
			@return GetOrgInviteApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	GetOrgInvite(ctx context.Context, orgId string, invitationId string) GetOrgInviteApiRequest
	/*
		GetOrgInvite Return One Invitation in One Organization by Invitation ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgInviteApiParams - Parameters for the request
		@return GetOrgInviteApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	GetOrgInviteWithParams(ctx context.Context, args *GetOrgInviteApiParams) GetOrgInviteApiRequest

	// Method available only for mocking purposes
	GetOrgInviteExecute(r GetOrgInviteApiRequest) (*OrganizationInvitation, *http.Response, error)

	/*
		GetOrgSettings Return Settings for One Organization

		Returns details about the specified organization's settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return GetOrgSettingsApiRequest
	*/
	GetOrgSettings(ctx context.Context, orgId string) GetOrgSettingsApiRequest
	/*
		GetOrgSettings Return Settings for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetOrgSettingsApiParams - Parameters for the request
		@return GetOrgSettingsApiRequest
	*/
	GetOrgSettingsWithParams(ctx context.Context, args *GetOrgSettingsApiParams) GetOrgSettingsApiRequest

	// Method available only for mocking purposes
	GetOrgSettingsExecute(r GetOrgSettingsApiRequest) (*OrganizationSettings, *http.Response, error)

	/*
			ListOrgInvites Return All Invitations in One Organization

			Returns all pending invitations to the specified organization.

		**Note**: Invitation management APIs are deprecated. Use Return All MongoDB Cloud Users in One Organization and filter by `orgMembershipStatus` to return all pending users.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return ListOrgInvitesApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	ListOrgInvites(ctx context.Context, orgId string) ListOrgInvitesApiRequest
	/*
		ListOrgInvites Return All Invitations in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgInvitesApiParams - Parameters for the request
		@return ListOrgInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	ListOrgInvitesWithParams(ctx context.Context, args *ListOrgInvitesApiParams) ListOrgInvitesApiRequest

	// Method available only for mocking purposes
	ListOrgInvitesExecute(r ListOrgInvitesApiRequest) ([]OrganizationInvitation, *http.Response, error)

	/*
		ListOrgs Return All Organizations

		Returns all organizations to which the requesting Service Account or API Key has access.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListOrgsApiRequest
	*/
	ListOrgs(ctx context.Context) ListOrgsApiRequest
	/*
		ListOrgs Return All Organizations


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListOrgsApiParams - Parameters for the request
		@return ListOrgsApiRequest
	*/
	ListOrgsWithParams(ctx context.Context, args *ListOrgsApiParams) ListOrgsApiRequest

	// Method available only for mocking purposes
	ListOrgsExecute(r ListOrgsApiRequest) (*PaginatedOrganization, *http.Response, error)

	/*
		UpdateOrg Update One Organization

		Updates one organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param atlasOrganization Details to update on the specified organization.
		@return UpdateOrgApiRequest
	*/
	UpdateOrg(ctx context.Context, orgId string, atlasOrganization *AtlasOrganization) UpdateOrgApiRequest
	/*
		UpdateOrg Update One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgApiParams - Parameters for the request
		@return UpdateOrgApiRequest
	*/
	UpdateOrgWithParams(ctx context.Context, args *UpdateOrgApiParams) UpdateOrgApiRequest

	// Method available only for mocking purposes
	UpdateOrgExecute(r UpdateOrgApiRequest) (*AtlasOrganization, *http.Response, error)

	/*
			UpdateOrgInviteById Update One Invitation in One Organization by Invitation ID

			Updates the details of one pending invitation to the specified organization. To specify which invitation, provide the unique identification string for that invitation. Use the Return All Organization Invitations endpoint to retrieve IDs for all pending organization invitations.

		**Note**: Invitation management APIs are deprecated. Use Update One MongoDB Cloud User in One Organization to update a pending user.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
			@param organizationInvitationUpdateRequest Updates the details of one pending invitation to the specified organization.
			@return UpdateOrgInviteByIdApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgInviteById(ctx context.Context, orgId string, invitationId string, organizationInvitationUpdateRequest *OrganizationInvitationUpdateRequest) UpdateOrgInviteByIdApiRequest
	/*
		UpdateOrgInviteById Update One Invitation in One Organization by Invitation ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgInviteByIdApiParams - Parameters for the request
		@return UpdateOrgInviteByIdApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgInviteByIdWithParams(ctx context.Context, args *UpdateOrgInviteByIdApiParams) UpdateOrgInviteByIdApiRequest

	// Method available only for mocking purposes
	UpdateOrgInviteByIdExecute(r UpdateOrgInviteByIdApiRequest) (*OrganizationInvitation, *http.Response, error)

	/*
			UpdateOrgInvites Update One Invitation in One Organization

			Updates the details of one pending invitation to the specified organization. To specify which invitation, provide the username of the invited user.

		**Note**:  Invitation management are deprecated. Use Update One MongoDB Cloud User in One Organization to update a pending user.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param organizationInvitationRequest Updates the details of one pending invitation to the specified organization.
			@return UpdateOrgInvitesApiRequest

			Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgInvites(ctx context.Context, orgId string, organizationInvitationRequest *OrganizationInvitationRequest) UpdateOrgInvitesApiRequest
	/*
		UpdateOrgInvites Update One Invitation in One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgInvitesApiParams - Parameters for the request
		@return UpdateOrgInvitesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgInvitesWithParams(ctx context.Context, args *UpdateOrgInvitesApiParams) UpdateOrgInvitesApiRequest

	// Method available only for mocking purposes
	UpdateOrgInvitesExecute(r UpdateOrgInvitesApiRequest) (*OrganizationInvitation, *http.Response, error)

	/*
		UpdateOrgSettings Update Settings for One Organization

		Updates the organization's settings.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param organizationSettings Details to update on the specified organization's settings.
		@return UpdateOrgSettingsApiRequest
	*/
	UpdateOrgSettings(ctx context.Context, orgId string, organizationSettings *OrganizationSettings) UpdateOrgSettingsApiRequest
	/*
		UpdateOrgSettings Update Settings for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgSettingsApiParams - Parameters for the request
		@return UpdateOrgSettingsApiRequest
	*/
	UpdateOrgSettingsWithParams(ctx context.Context, args *UpdateOrgSettingsApiParams) UpdateOrgSettingsApiRequest

	// Method available only for mocking purposes
	UpdateOrgSettingsExecute(r UpdateOrgSettingsApiRequest) (*OrganizationSettings, *http.Response, error)

	/*
		UpdateOrgUserRoles Update Organization Roles for One MongoDB Cloud User

		Updates the roles of the specified user in the specified organization. To specify the user to update, provide the unique 24-hexadecimal digit string that identifies the user in the specified organization.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param userId Unique 24-hexadecimal digit string that identifies the user to modify.
		@param updateOrgRolesForUser Roles to update for the specified user.
		@return UpdateOrgUserRolesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgUserRoles(ctx context.Context, orgId string, userId string, updateOrgRolesForUser *UpdateOrgRolesForUser) UpdateOrgUserRolesApiRequest
	/*
		UpdateOrgUserRoles Update Organization Roles for One MongoDB Cloud User


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateOrgUserRolesApiParams - Parameters for the request
		@return UpdateOrgUserRolesApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for OrganizationsApi
	*/
	UpdateOrgUserRolesWithParams(ctx context.Context, args *UpdateOrgUserRolesApiParams) UpdateOrgUserRolesApiRequest

	// Method available only for mocking purposes
	UpdateOrgUserRolesExecute(r UpdateOrgUserRolesApiRequest) (*UpdateOrgRolesForUser, *http.Response, error)
}

// OrganizationsApiService OrganizationsApi service
type OrganizationsApiService service

type CreateOrgApiRequest struct {
	ctx                       context.Context
	ApiService                OrganizationsApi
	createOrganizationRequest *CreateOrganizationRequest
}

type CreateOrgApiParams struct {
	CreateOrganizationRequest *CreateOrganizationRequest
}

func (a *OrganizationsApiService) CreateOrgWithParams(ctx context.Context, args *CreateOrgApiParams) CreateOrgApiRequest {
	return CreateOrgApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		createOrganizationRequest: args.CreateOrganizationRequest,
	}
}

func (r CreateOrgApiRequest) Execute() (*CreateOrganizationResponse, *http.Response, error) {
	return r.ApiService.CreateOrgExecute(r)
}

/*
CreateOrg Create One Organization

Creates one organization in MongoDB Cloud and links it to the requesting Service Account's or API Key's organization. The requesting Service Account's or API Key's organization must be a paying organization. To learn more, see Configure a Paying Organization in the MongoDB Atlas documentation. Optionally, if `federationSettingsId` is provided, the new Organization will be linked to the federation. The requesting Service Account or API Key must be an Organization Owner in the federation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return CreateOrgApiRequest
*/
func (a *OrganizationsApiService) CreateOrg(ctx context.Context, createOrganizationRequest *CreateOrganizationRequest) CreateOrgApiRequest {
	return CreateOrgApiRequest{
		ApiService:                a,
		ctx:                       ctx,
		createOrganizationRequest: createOrganizationRequest,
	}
}

// CreateOrgExecute executes the request
//
//	@return CreateOrganizationResponse
func (a *OrganizationsApiService) CreateOrgExecute(r CreateOrgApiRequest) (*CreateOrganizationResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CreateOrganizationResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.CreateOrg")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.createOrganizationRequest == nil {
		return localVarReturnValue, nil, reportError("createOrganizationRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.createOrganizationRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type CreateOrgInviteApiRequest struct {
	ctx                           context.Context
	ApiService                    OrganizationsApi
	orgId                         string
	organizationInvitationRequest *OrganizationInvitationRequest
}

type CreateOrgInviteApiParams struct {
	OrgId                         string
	OrganizationInvitationRequest *OrganizationInvitationRequest
}

func (a *OrganizationsApiService) CreateOrgInviteWithParams(ctx context.Context, args *CreateOrgInviteApiParams) CreateOrgInviteApiRequest {
	return CreateOrgInviteApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         args.OrgId,
		organizationInvitationRequest: args.OrganizationInvitationRequest,
	}
}

func (r CreateOrgInviteApiRequest) Execute() (*OrganizationInvitation, *http.Response, error) {
	return r.ApiService.CreateOrgInviteExecute(r)
}

/*
CreateOrgInvite Create Invitation for One MongoDB Cloud User in One Organization

Invites one MongoDB Cloud user to join the specified organization. The user must accept the invitation to access information within the specified organization.

**Note**: Invitation management APIs are deprecated. Use Add One MongoDB Cloud User to One Organization to invite a user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateOrgInviteApiRequest

Deprecated
*/
func (a *OrganizationsApiService) CreateOrgInvite(ctx context.Context, orgId string, organizationInvitationRequest *OrganizationInvitationRequest) CreateOrgInviteApiRequest {
	return CreateOrgInviteApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         orgId,
		organizationInvitationRequest: organizationInvitationRequest,
	}
}

// CreateOrgInviteExecute executes the request
//
//	@return OrganizationInvitation
//
// Deprecated
func (a *OrganizationsApiService) CreateOrgInviteExecute(r CreateOrgInviteApiRequest) (*OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.CreateOrgInvite")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.organizationInvitationRequest == nil {
		return localVarReturnValue, nil, reportError("organizationInvitationRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.organizationInvitationRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type DeleteOrgApiRequest struct {
	ctx        context.Context
	ApiService OrganizationsApi
	orgId      string
}

type DeleteOrgApiParams struct {
	OrgId string
}

func (a *OrganizationsApiService) DeleteOrgWithParams(ctx context.Context, args *DeleteOrgApiParams) DeleteOrgApiRequest {
	return DeleteOrgApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r DeleteOrgApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgExecute(r)
}

/*
DeleteOrg Remove One Organization

Removes one specified organization. MongoDB Cloud imposes the following limits on this resource:

  - Organizations with active projects cannot be removed.

  - All projects in the organization must be removed before you can remove the organization.

    @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
    @param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
    @return DeleteOrgApiRequest
*/
func (a *OrganizationsApiService) DeleteOrg(ctx context.Context, orgId string) DeleteOrgApiRequest {
	return DeleteOrgApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// DeleteOrgExecute executes the request
func (a *OrganizationsApiService) DeleteOrgExecute(r DeleteOrgApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.DeleteOrg")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type DeleteOrgInviteApiRequest struct {
	ctx          context.Context
	ApiService   OrganizationsApi
	orgId        string
	invitationId string
}

type DeleteOrgInviteApiParams struct {
	OrgId        string
	InvitationId string
}

func (a *OrganizationsApiService) DeleteOrgInviteWithParams(ctx context.Context, args *DeleteOrgInviteApiParams) DeleteOrgInviteApiRequest {
	return DeleteOrgInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		invitationId: args.InvitationId,
	}
}

func (r DeleteOrgInviteApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteOrgInviteExecute(r)
}

/*
DeleteOrgInvite Remove One Invitation from One Organization

Cancels one pending invitation sent to the specified MongoDB Cloud user to join an organization. You can't cancel an invitation that the user accepted.

**Note**: Invitation management APIs are deprecated. Use Remove One MongoDB Cloud User From One Organization to remove a pending user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return DeleteOrgInviteApiRequest

Deprecated
*/
func (a *OrganizationsApiService) DeleteOrgInvite(ctx context.Context, orgId string, invitationId string) DeleteOrgInviteApiRequest {
	return DeleteOrgInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        orgId,
		invitationId: invitationId,
	}
}

// DeleteOrgInviteExecute executes the request
// Deprecated
func (a *OrganizationsApiService) DeleteOrgInviteExecute(r DeleteOrgInviteApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.DeleteOrgInvite")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites/{invitationId}"
	if r.orgId == "" {
		return nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invitationId == "" {
		return nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarHTTPResponse, newErr
	}

	return localVarHTTPResponse, nil
}

type GetOrgApiRequest struct {
	ctx        context.Context
	ApiService OrganizationsApi
	orgId      string
}

type GetOrgApiParams struct {
	OrgId string
}

func (a *OrganizationsApiService) GetOrgWithParams(ctx context.Context, args *GetOrgApiParams) GetOrgApiRequest {
	return GetOrgApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r GetOrgApiRequest) Execute() (*AtlasOrganization, *http.Response, error) {
	return r.ApiService.GetOrgExecute(r)
}

/*
GetOrg Return One Organization

Returns one organization to which the requesting Service Account or API Key has access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return GetOrgApiRequest
*/
func (a *OrganizationsApiService) GetOrg(ctx context.Context, orgId string) GetOrgApiRequest {
	return GetOrgApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// GetOrgExecute executes the request
//
//	@return AtlasOrganization
func (a *OrganizationsApiService) GetOrgExecute(r GetOrgApiRequest) (*AtlasOrganization, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AtlasOrganization
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.GetOrg")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type GetOrgGroupsApiRequest struct {
	ctx          context.Context
	ApiService   OrganizationsApi
	orgId        string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	name         *string
}

type GetOrgGroupsApiParams struct {
	OrgId        string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Name         *string
}

func (a *OrganizationsApiService) GetOrgGroupsWithParams(ctx context.Context, args *GetOrgGroupsApiParams) GetOrgGroupsApiRequest {
	return GetOrgGroupsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		name:         args.Name,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r GetOrgGroupsApiRequest) IncludeCount(includeCount bool) GetOrgGroupsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r GetOrgGroupsApiRequest) ItemsPerPage(itemsPerPage int) GetOrgGroupsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r GetOrgGroupsApiRequest) PageNum(pageNum int) GetOrgGroupsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Human-readable label of the project to use to filter the returned list. Performs a case-insensitive search for a project within the organization which is prefixed by the specified name.
func (r GetOrgGroupsApiRequest) Name(name string) GetOrgGroupsApiRequest {
	r.name = &name
	return r
}

func (r GetOrgGroupsApiRequest) Execute() (*PaginatedAtlasGroup, *http.Response, error) {
	return r.ApiService.GetOrgGroupsExecute(r)
}

/*
GetOrgGroups Return All Projects in One Organization

Returns multiple projects in the specified organization. Each organization can have multiple projects. Use projects to:

- Isolate different environments, such as development, test, or production environments, from each other.
- Associate different MongoDB Cloud users or teams with different environments, or give different permission to MongoDB Cloud users in different environments.
- Maintain separate cluster security configurations.
- Create different alert settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return GetOrgGroupsApiRequest
*/
func (a *OrganizationsApiService) GetOrgGroups(ctx context.Context, orgId string) GetOrgGroupsApiRequest {
	return GetOrgGroupsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// GetOrgGroupsExecute executes the request
//
//	@return PaginatedAtlasGroup
func (a *OrganizationsApiService) GetOrgGroupsExecute(r GetOrgGroupsApiRequest) (*PaginatedAtlasGroup, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedAtlasGroup
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.GetOrgGroups")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/groups"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeCount != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	} else {
		var defaultValue bool = true
		r.includeCount = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	}
	if r.itemsPerPage != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	} else {
		var defaultValue int = 100
		r.itemsPerPage = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	}
	if r.pageNum != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	} else {
		var defaultValue int = 1
		r.pageNum = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	}
	if r.name != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name", r.name, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type GetOrgInviteApiRequest struct {
	ctx          context.Context
	ApiService   OrganizationsApi
	orgId        string
	invitationId string
}

type GetOrgInviteApiParams struct {
	OrgId        string
	InvitationId string
}

func (a *OrganizationsApiService) GetOrgInviteWithParams(ctx context.Context, args *GetOrgInviteApiParams) GetOrgInviteApiRequest {
	return GetOrgInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        args.OrgId,
		invitationId: args.InvitationId,
	}
}

func (r GetOrgInviteApiRequest) Execute() (*OrganizationInvitation, *http.Response, error) {
	return r.ApiService.GetOrgInviteExecute(r)
}

/*
GetOrgInvite Return One Invitation in One Organization by Invitation ID

Returns the details of one pending invitation to the specified organization.

**Note**: Invitation management APIs are deprecated. Use Return One MongoDB Cloud User in One Organization to return a pending user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return GetOrgInviteApiRequest

Deprecated
*/
func (a *OrganizationsApiService) GetOrgInvite(ctx context.Context, orgId string, invitationId string) GetOrgInviteApiRequest {
	return GetOrgInviteApiRequest{
		ApiService:   a,
		ctx:          ctx,
		orgId:        orgId,
		invitationId: invitationId,
	}
}

// GetOrgInviteExecute executes the request
//
//	@return OrganizationInvitation
//
// Deprecated
func (a *OrganizationsApiService) GetOrgInviteExecute(r GetOrgInviteApiRequest) (*OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.GetOrgInvite")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites/{invitationId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invitationId == "" {
		return localVarReturnValue, nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type GetOrgSettingsApiRequest struct {
	ctx        context.Context
	ApiService OrganizationsApi
	orgId      string
}

type GetOrgSettingsApiParams struct {
	OrgId string
}

func (a *OrganizationsApiService) GetOrgSettingsWithParams(ctx context.Context, args *GetOrgSettingsApiParams) GetOrgSettingsApiRequest {
	return GetOrgSettingsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r GetOrgSettingsApiRequest) Execute() (*OrganizationSettings, *http.Response, error) {
	return r.ApiService.GetOrgSettingsExecute(r)
}

/*
GetOrgSettings Return Settings for One Organization

Returns details about the specified organization's settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return GetOrgSettingsApiRequest
*/
func (a *OrganizationsApiService) GetOrgSettings(ctx context.Context, orgId string) GetOrgSettingsApiRequest {
	return GetOrgSettingsApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// GetOrgSettingsExecute executes the request
//
//	@return OrganizationSettings
func (a *OrganizationsApiService) GetOrgSettingsExecute(r GetOrgSettingsApiRequest) (*OrganizationSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.GetOrgSettings")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/settings"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ListOrgInvitesApiRequest struct {
	ctx        context.Context
	ApiService OrganizationsApi
	orgId      string
	username   *string
}

type ListOrgInvitesApiParams struct {
	OrgId    string
	Username *string
}

func (a *OrganizationsApiService) ListOrgInvitesWithParams(ctx context.Context, args *ListOrgInvitesApiParams) ListOrgInvitesApiRequest {
	return ListOrgInvitesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		username:   args.Username,
	}
}

// Email address of the user account invited to this organization. If you exclude this parameter, this resource returns all pending invitations.
func (r ListOrgInvitesApiRequest) Username(username string) ListOrgInvitesApiRequest {
	r.username = &username
	return r
}

func (r ListOrgInvitesApiRequest) Execute() ([]OrganizationInvitation, *http.Response, error) {
	return r.ApiService.ListOrgInvitesExecute(r)
}

/*
ListOrgInvites Return All Invitations in One Organization

Returns all pending invitations to the specified organization.

**Note**: Invitation management APIs are deprecated. Use Return All MongoDB Cloud Users in One Organization and filter by `orgMembershipStatus` to return all pending users.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListOrgInvitesApiRequest

Deprecated
*/
func (a *OrganizationsApiService) ListOrgInvites(ctx context.Context, orgId string) ListOrgInvitesApiRequest {
	return ListOrgInvitesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListOrgInvitesExecute executes the request
//
//	@return []OrganizationInvitation
//
// Deprecated
func (a *OrganizationsApiService) ListOrgInvitesExecute(r ListOrgInvitesApiRequest) ([]OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue []OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.ListOrgInvites")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.username != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "username", r.username, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ListOrgsApiRequest struct {
	ctx          context.Context
	ApiService   OrganizationsApi
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	name         *string
}

type ListOrgsApiParams struct {
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	Name         *string
}

func (a *OrganizationsApiService) ListOrgsWithParams(ctx context.Context, args *ListOrgsApiParams) ListOrgsApiRequest {
	return ListOrgsApiRequest{
		ApiService:   a,
		ctx:          ctx,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		name:         args.Name,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListOrgsApiRequest) IncludeCount(includeCount bool) ListOrgsApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListOrgsApiRequest) ItemsPerPage(itemsPerPage int) ListOrgsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListOrgsApiRequest) PageNum(pageNum int) ListOrgsApiRequest {
	r.pageNum = &pageNum
	return r
}

// Human-readable label of the organization to use to filter the returned list. Performs a case-insensitive search for an organization that starts with the specified name.
func (r ListOrgsApiRequest) Name(name string) ListOrgsApiRequest {
	r.name = &name
	return r
}

func (r ListOrgsApiRequest) Execute() (*PaginatedOrganization, *http.Response, error) {
	return r.ApiService.ListOrgsExecute(r)
}

/*
ListOrgs Return All Organizations

Returns all organizations to which the requesting Service Account or API Key has access.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListOrgsApiRequest
*/
func (a *OrganizationsApiService) ListOrgs(ctx context.Context) ListOrgsApiRequest {
	return ListOrgsApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListOrgsExecute executes the request
//
//	@return PaginatedOrganization
func (a *OrganizationsApiService) ListOrgsExecute(r ListOrgsApiRequest) (*PaginatedOrganization, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedOrganization
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.ListOrgs")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.includeCount != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	} else {
		var defaultValue bool = true
		r.includeCount = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "includeCount", r.includeCount, "")
	}
	if r.itemsPerPage != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	} else {
		var defaultValue int = 100
		r.itemsPerPage = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "itemsPerPage", r.itemsPerPage, "")
	}
	if r.pageNum != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	} else {
		var defaultValue int = 1
		r.pageNum = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageNum", r.pageNum, "")
	}
	if r.name != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "name", r.name, "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateOrgApiRequest struct {
	ctx               context.Context
	ApiService        OrganizationsApi
	orgId             string
	atlasOrganization *AtlasOrganization
}

type UpdateOrgApiParams struct {
	OrgId             string
	AtlasOrganization *AtlasOrganization
}

func (a *OrganizationsApiService) UpdateOrgWithParams(ctx context.Context, args *UpdateOrgApiParams) UpdateOrgApiRequest {
	return UpdateOrgApiRequest{
		ApiService:        a,
		ctx:               ctx,
		orgId:             args.OrgId,
		atlasOrganization: args.AtlasOrganization,
	}
}

func (r UpdateOrgApiRequest) Execute() (*AtlasOrganization, *http.Response, error) {
	return r.ApiService.UpdateOrgExecute(r)
}

/*
UpdateOrg Update One Organization

Updates one organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return UpdateOrgApiRequest
*/
func (a *OrganizationsApiService) UpdateOrg(ctx context.Context, orgId string, atlasOrganization *AtlasOrganization) UpdateOrgApiRequest {
	return UpdateOrgApiRequest{
		ApiService:        a,
		ctx:               ctx,
		orgId:             orgId,
		atlasOrganization: atlasOrganization,
	}
}

// UpdateOrgExecute executes the request
//
//	@return AtlasOrganization
func (a *OrganizationsApiService) UpdateOrgExecute(r UpdateOrgApiRequest) (*AtlasOrganization, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *AtlasOrganization
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.UpdateOrg")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.atlasOrganization == nil {
		return localVarReturnValue, nil, reportError("atlasOrganization is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.atlasOrganization
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateOrgInviteByIdApiRequest struct {
	ctx                                 context.Context
	ApiService                          OrganizationsApi
	orgId                               string
	invitationId                        string
	organizationInvitationUpdateRequest *OrganizationInvitationUpdateRequest
}

type UpdateOrgInviteByIdApiParams struct {
	OrgId                               string
	InvitationId                        string
	OrganizationInvitationUpdateRequest *OrganizationInvitationUpdateRequest
}

func (a *OrganizationsApiService) UpdateOrgInviteByIdWithParams(ctx context.Context, args *UpdateOrgInviteByIdApiParams) UpdateOrgInviteByIdApiRequest {
	return UpdateOrgInviteByIdApiRequest{
		ApiService:                          a,
		ctx:                                 ctx,
		orgId:                               args.OrgId,
		invitationId:                        args.InvitationId,
		organizationInvitationUpdateRequest: args.OrganizationInvitationUpdateRequest,
	}
}

func (r UpdateOrgInviteByIdApiRequest) Execute() (*OrganizationInvitation, *http.Response, error) {
	return r.ApiService.UpdateOrgInviteByIdExecute(r)
}

/*
UpdateOrgInviteById Update One Invitation in One Organization by Invitation ID

Updates the details of one pending invitation to the specified organization. To specify which invitation, provide the unique identification string for that invitation. Use the Return All Organization Invitations endpoint to retrieve IDs for all pending organization invitations.

**Note**: Invitation management APIs are deprecated. Use Update One MongoDB Cloud User in One Organization to update a pending user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invitationId Unique 24-hexadecimal digit string that identifies the invitation.
	@return UpdateOrgInviteByIdApiRequest

Deprecated
*/
func (a *OrganizationsApiService) UpdateOrgInviteById(ctx context.Context, orgId string, invitationId string, organizationInvitationUpdateRequest *OrganizationInvitationUpdateRequest) UpdateOrgInviteByIdApiRequest {
	return UpdateOrgInviteByIdApiRequest{
		ApiService:                          a,
		ctx:                                 ctx,
		orgId:                               orgId,
		invitationId:                        invitationId,
		organizationInvitationUpdateRequest: organizationInvitationUpdateRequest,
	}
}

// UpdateOrgInviteByIdExecute executes the request
//
//	@return OrganizationInvitation
//
// Deprecated
func (a *OrganizationsApiService) UpdateOrgInviteByIdExecute(r UpdateOrgInviteByIdApiRequest) (*OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.UpdateOrgInviteById")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites/{invitationId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invitationId == "" {
		return localVarReturnValue, nil, reportError("invitationId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invitationId"+"}", url.PathEscape(r.invitationId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.organizationInvitationUpdateRequest == nil {
		return localVarReturnValue, nil, reportError("organizationInvitationUpdateRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.organizationInvitationUpdateRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateOrgInvitesApiRequest struct {
	ctx                           context.Context
	ApiService                    OrganizationsApi
	orgId                         string
	organizationInvitationRequest *OrganizationInvitationRequest
}

type UpdateOrgInvitesApiParams struct {
	OrgId                         string
	OrganizationInvitationRequest *OrganizationInvitationRequest
}

func (a *OrganizationsApiService) UpdateOrgInvitesWithParams(ctx context.Context, args *UpdateOrgInvitesApiParams) UpdateOrgInvitesApiRequest {
	return UpdateOrgInvitesApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         args.OrgId,
		organizationInvitationRequest: args.OrganizationInvitationRequest,
	}
}

func (r UpdateOrgInvitesApiRequest) Execute() (*OrganizationInvitation, *http.Response, error) {
	return r.ApiService.UpdateOrgInvitesExecute(r)
}

/*
UpdateOrgInvites Update One Invitation in One Organization

Updates the details of one pending invitation to the specified organization. To specify which invitation, provide the username of the invited user.

**Note**:  Invitation management are deprecated. Use Update One MongoDB Cloud User in One Organization to update a pending user.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return UpdateOrgInvitesApiRequest

Deprecated
*/
func (a *OrganizationsApiService) UpdateOrgInvites(ctx context.Context, orgId string, organizationInvitationRequest *OrganizationInvitationRequest) UpdateOrgInvitesApiRequest {
	return UpdateOrgInvitesApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         orgId,
		organizationInvitationRequest: organizationInvitationRequest,
	}
}

// UpdateOrgInvitesExecute executes the request
//
//	@return OrganizationInvitation
//
// Deprecated
func (a *OrganizationsApiService) UpdateOrgInvitesExecute(r UpdateOrgInvitesApiRequest) (*OrganizationInvitation, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationInvitation
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.UpdateOrgInvites")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invites"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.organizationInvitationRequest == nil {
		return localVarReturnValue, nil, reportError("organizationInvitationRequest is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.organizationInvitationRequest
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateOrgSettingsApiRequest struct {
	ctx                  context.Context
	ApiService           OrganizationsApi
	orgId                string
	organizationSettings *OrganizationSettings
}

type UpdateOrgSettingsApiParams struct {
	OrgId                string
	OrganizationSettings *OrganizationSettings
}

func (a *OrganizationsApiService) UpdateOrgSettingsWithParams(ctx context.Context, args *UpdateOrgSettingsApiParams) UpdateOrgSettingsApiRequest {
	return UpdateOrgSettingsApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		orgId:                args.OrgId,
		organizationSettings: args.OrganizationSettings,
	}
}

func (r UpdateOrgSettingsApiRequest) Execute() (*OrganizationSettings, *http.Response, error) {
	return r.ApiService.UpdateOrgSettingsExecute(r)
}

/*
UpdateOrgSettings Update Settings for One Organization

Updates the organization's settings.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return UpdateOrgSettingsApiRequest
*/
func (a *OrganizationsApiService) UpdateOrgSettings(ctx context.Context, orgId string, organizationSettings *OrganizationSettings) UpdateOrgSettingsApiRequest {
	return UpdateOrgSettingsApiRequest{
		ApiService:           a,
		ctx:                  ctx,
		orgId:                orgId,
		organizationSettings: organizationSettings,
	}
}

// UpdateOrgSettingsExecute executes the request
//
//	@return OrganizationSettings
func (a *OrganizationsApiService) UpdateOrgSettingsExecute(r UpdateOrgSettingsApiRequest) (*OrganizationSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *OrganizationSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.UpdateOrgSettings")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/settings"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.organizationSettings == nil {
		return localVarReturnValue, nil, reportError("organizationSettings is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.organizationSettings
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type UpdateOrgUserRolesApiRequest struct {
	ctx                   context.Context
	ApiService            OrganizationsApi
	orgId                 string
	userId                string
	updateOrgRolesForUser *UpdateOrgRolesForUser
}

type UpdateOrgUserRolesApiParams struct {
	OrgId                 string
	UserId                string
	UpdateOrgRolesForUser *UpdateOrgRolesForUser
}

func (a *OrganizationsApiService) UpdateOrgUserRolesWithParams(ctx context.Context, args *UpdateOrgUserRolesApiParams) UpdateOrgUserRolesApiRequest {
	return UpdateOrgUserRolesApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 args.OrgId,
		userId:                args.UserId,
		updateOrgRolesForUser: args.UpdateOrgRolesForUser,
	}
}

func (r UpdateOrgUserRolesApiRequest) Execute() (*UpdateOrgRolesForUser, *http.Response, error) {
	return r.ApiService.UpdateOrgUserRolesExecute(r)
}

/*
UpdateOrgUserRoles Update Organization Roles for One MongoDB Cloud User

Updates the roles of the specified user in the specified organization. To specify the user to update, provide the unique 24-hexadecimal digit string that identifies the user in the specified organization.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param userId Unique 24-hexadecimal digit string that identifies the user to modify.
	@return UpdateOrgUserRolesApiRequest

Deprecated
*/
func (a *OrganizationsApiService) UpdateOrgUserRoles(ctx context.Context, orgId string, userId string, updateOrgRolesForUser *UpdateOrgRolesForUser) UpdateOrgUserRolesApiRequest {
	return UpdateOrgUserRolesApiRequest{
		ApiService:            a,
		ctx:                   ctx,
		orgId:                 orgId,
		userId:                userId,
		updateOrgRolesForUser: updateOrgRolesForUser,
	}
}

// UpdateOrgUserRolesExecute executes the request
//
//	@return UpdateOrgRolesForUser
//
// Deprecated
func (a *OrganizationsApiService) UpdateOrgUserRolesExecute(r UpdateOrgUserRolesApiRequest) (*UpdateOrgRolesForUser, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPut
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *UpdateOrgRolesForUser
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "OrganizationsApiService.UpdateOrgUserRoles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/users/{userId}/roles"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.userId == "" {
		return localVarReturnValue, nil, reportError("userId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"userId"+"}", url.PathEscape(r.userId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.updateOrgRolesForUser == nil {
		return localVarReturnValue, nil, reportError("updateOrgRolesForUser is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.updateOrgRolesForUser
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := a.client.makeApiError(localVarHTTPResponse, localVarHTTPMethod, localVarPath)
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarHTTPResponse.Body, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		defer localVarHTTPResponse.Body.Close()
		buf, readErr := io.ReadAll(localVarHTTPResponse.Body)
		if readErr != nil {
			err = readErr
		}
		newErr := &GenericOpenAPIError{
			body:  buf,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
