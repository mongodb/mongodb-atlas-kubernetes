// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type NetworkPeeringApi interface {

	/*
		CreateGroupContainer Create One Network Peering Container

		Creates one new network peering container in the specified project. MongoDB Cloud can deploy Network Peering connections in a network peering container. GCP can have one container per project. AWS and Azure can have one container per cloud provider region.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param cloudProviderContainer Creates one new network peering container in the specified project.
		@return CreateGroupContainerApiRequest
	*/
	CreateGroupContainer(ctx context.Context, groupId string, cloudProviderContainer *CloudProviderContainer) CreateGroupContainerApiRequest
	/*
		CreateGroupContainer Create One Network Peering Container


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupContainerApiParams - Parameters for the request
		@return CreateGroupContainerApiRequest
	*/
	CreateGroupContainerWithParams(ctx context.Context, args *CreateGroupContainerApiParams) CreateGroupContainerApiRequest

	// Method available only for mocking purposes
	CreateGroupContainerExecute(r CreateGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error)

	/*
		CreateGroupPeer Create One Network Peering Connection

		Creates one new network peering connection in the specified project. Network peering allows multiple cloud-hosted applications to securely connect to the same project. To learn more about considerations and prerequisites, see the Network Peering Documentation.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param baseNetworkPeeringConnectionSettings Create one network peering connection.
		@return CreateGroupPeerApiRequest
	*/
	CreateGroupPeer(ctx context.Context, groupId string, baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings) CreateGroupPeerApiRequest
	/*
		CreateGroupPeer Create One Network Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateGroupPeerApiParams - Parameters for the request
		@return CreateGroupPeerApiRequest
	*/
	CreateGroupPeerWithParams(ctx context.Context, args *CreateGroupPeerApiParams) CreateGroupPeerApiRequest

	// Method available only for mocking purposes
	CreateGroupPeerExecute(r CreateGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error)

	/*
		DeleteGroupContainer Remove One Network Peering Container

		Removes one network peering container in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container that you want to remove.
		@return DeleteGroupContainerApiRequest
	*/
	DeleteGroupContainer(ctx context.Context, groupId string, containerId string) DeleteGroupContainerApiRequest
	/*
		DeleteGroupContainer Remove One Network Peering Container


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupContainerApiParams - Parameters for the request
		@return DeleteGroupContainerApiRequest
	*/
	DeleteGroupContainerWithParams(ctx context.Context, args *DeleteGroupContainerApiParams) DeleteGroupContainerApiRequest

	// Method available only for mocking purposes
	DeleteGroupContainerExecute(r DeleteGroupContainerApiRequest) (*http.Response, error)

	/*
		DeleteGroupPeer Remove One Network Peering Connection

		Removes one network peering connection in the specified project. If you remove the last network peering connection associated with a project, MongoDB Cloud also removes any AWS security groups from the project IP access list.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to delete.
		@return DeleteGroupPeerApiRequest
	*/
	DeleteGroupPeer(ctx context.Context, groupId string, peerId string) DeleteGroupPeerApiRequest
	/*
		DeleteGroupPeer Remove One Network Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DeleteGroupPeerApiParams - Parameters for the request
		@return DeleteGroupPeerApiRequest
	*/
	DeleteGroupPeerWithParams(ctx context.Context, args *DeleteGroupPeerApiParams) DeleteGroupPeerApiRequest

	// Method available only for mocking purposes
	DeleteGroupPeerExecute(r DeleteGroupPeerApiRequest) (any, *http.Response, error)

	/*
		DisablePeering Disable Connect via Peering-Only Mode for One Project

		Disables Connect via Peering Only mode for the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param privateIPMode Disables Connect via Peering Only mode for the specified project.
		@return DisablePeeringApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for NetworkPeeringApi
	*/
	DisablePeering(ctx context.Context, groupId string, privateIPMode *PrivateIPMode) DisablePeeringApiRequest
	/*
		DisablePeering Disable Connect via Peering-Only Mode for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param DisablePeeringApiParams - Parameters for the request
		@return DisablePeeringApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for NetworkPeeringApi
	*/
	DisablePeeringWithParams(ctx context.Context, args *DisablePeeringApiParams) DisablePeeringApiRequest

	// Method available only for mocking purposes
	DisablePeeringExecute(r DisablePeeringApiRequest) (*PrivateIPMode, *http.Response, error)

	/*
		GetGroupContainer Return One Network Peering Container

		Returns details about one network peering container in one specified project. Network peering containers contain network peering connections.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container.
		@return GetGroupContainerApiRequest
	*/
	GetGroupContainer(ctx context.Context, groupId string, containerId string) GetGroupContainerApiRequest
	/*
		GetGroupContainer Return One Network Peering Container


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupContainerApiParams - Parameters for the request
		@return GetGroupContainerApiRequest
	*/
	GetGroupContainerWithParams(ctx context.Context, args *GetGroupContainerApiParams) GetGroupContainerApiRequest

	// Method available only for mocking purposes
	GetGroupContainerExecute(r GetGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error)

	/*
		GetGroupPeer Return One Network Peering Connection in One Project

		Returns details about one specified network peering connection in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to retrieve.
		@return GetGroupPeerApiRequest
	*/
	GetGroupPeer(ctx context.Context, groupId string, peerId string) GetGroupPeerApiRequest
	/*
		GetGroupPeer Return One Network Peering Connection in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetGroupPeerApiParams - Parameters for the request
		@return GetGroupPeerApiRequest
	*/
	GetGroupPeerWithParams(ctx context.Context, args *GetGroupPeerApiParams) GetGroupPeerApiRequest

	// Method available only for mocking purposes
	GetGroupPeerExecute(r GetGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error)

	/*
		ListGroupContainerAll Return All Network Peering Containers in One Project

		Returns details about all network peering containers in the specified project. Network peering containers contain network peering connections.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupContainerAllApiRequest
	*/
	ListGroupContainerAll(ctx context.Context, groupId string) ListGroupContainerAllApiRequest
	/*
		ListGroupContainerAll Return All Network Peering Containers in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupContainerAllApiParams - Parameters for the request
		@return ListGroupContainerAllApiRequest
	*/
	ListGroupContainerAllWithParams(ctx context.Context, args *ListGroupContainerAllApiParams) ListGroupContainerAllApiRequest

	// Method available only for mocking purposes
	ListGroupContainerAllExecute(r ListGroupContainerAllApiRequest) (*PaginatedCloudProviderContainer, *http.Response, error)

	/*
		ListGroupContainers Return All Network Peering Containers in One Project for One Cloud Provider

		Returns details about all network peering containers in the specified project for the specified cloud provider. If you do not specify the cloud provider, MongoDB Cloud returns details about all network peering containers in the project for Amazon Web Services (AWS).

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupContainersApiRequest
	*/
	ListGroupContainers(ctx context.Context, groupId string) ListGroupContainersApiRequest
	/*
		ListGroupContainers Return All Network Peering Containers in One Project for One Cloud Provider


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupContainersApiParams - Parameters for the request
		@return ListGroupContainersApiRequest
	*/
	ListGroupContainersWithParams(ctx context.Context, args *ListGroupContainersApiParams) ListGroupContainersApiRequest

	// Method available only for mocking purposes
	ListGroupContainersExecute(r ListGroupContainersApiRequest) (*PaginatedCloudProviderContainer, *http.Response, error)

	/*
		ListGroupPeers Return All Network Peering Connections in One Project

		Returns details about all network peering connections in the specified project. Network peering allows multiple cloud-hosted applications to securely connect to the same project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return ListGroupPeersApiRequest
	*/
	ListGroupPeers(ctx context.Context, groupId string) ListGroupPeersApiRequest
	/*
		ListGroupPeers Return All Network Peering Connections in One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListGroupPeersApiParams - Parameters for the request
		@return ListGroupPeersApiRequest
	*/
	ListGroupPeersWithParams(ctx context.Context, args *ListGroupPeersApiParams) ListGroupPeersApiRequest

	// Method available only for mocking purposes
	ListGroupPeersExecute(r ListGroupPeersApiRequest) (*PaginatedContainerPeer, *http.Response, error)

	/*
		UpdateGroupContainer Update One Network Peering Container

		Updates the network details and labels of one specified network peering container in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container that you want to remove.
		@param cloudProviderContainer Updates the network details and labels of one specified network peering container in the specified project.
		@return UpdateGroupContainerApiRequest
	*/
	UpdateGroupContainer(ctx context.Context, groupId string, containerId string, cloudProviderContainer *CloudProviderContainer) UpdateGroupContainerApiRequest
	/*
		UpdateGroupContainer Update One Network Peering Container


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupContainerApiParams - Parameters for the request
		@return UpdateGroupContainerApiRequest
	*/
	UpdateGroupContainerWithParams(ctx context.Context, args *UpdateGroupContainerApiParams) UpdateGroupContainerApiRequest

	// Method available only for mocking purposes
	UpdateGroupContainerExecute(r UpdateGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error)

	/*
		UpdateGroupPeer Update One Network Peering Connection

		Updates one specified network peering connection in the specified project.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to update.
		@param baseNetworkPeeringConnectionSettings Modify one network peering connection.
		@return UpdateGroupPeerApiRequest
	*/
	UpdateGroupPeer(ctx context.Context, groupId string, peerId string, baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings) UpdateGroupPeerApiRequest
	/*
		UpdateGroupPeer Update One Network Peering Connection


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param UpdateGroupPeerApiParams - Parameters for the request
		@return UpdateGroupPeerApiRequest
	*/
	UpdateGroupPeerWithParams(ctx context.Context, args *UpdateGroupPeerApiParams) UpdateGroupPeerApiRequest

	// Method available only for mocking purposes
	UpdateGroupPeerExecute(r UpdateGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error)

	/*
		VerifyPrivateIpMode Verify Connect via Peering-Only Mode for One Project

		Verifies if someone set the specified project to **Connect via Peering Only** mode.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
		@return VerifyPrivateIpModeApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for NetworkPeeringApi
	*/
	VerifyPrivateIpMode(ctx context.Context, groupId string) VerifyPrivateIpModeApiRequest
	/*
		VerifyPrivateIpMode Verify Connect via Peering-Only Mode for One Project


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param VerifyPrivateIpModeApiParams - Parameters for the request
		@return VerifyPrivateIpModeApiRequest

		Deprecated: this method has been deprecated. Please check the latest resource version for NetworkPeeringApi
	*/
	VerifyPrivateIpModeWithParams(ctx context.Context, args *VerifyPrivateIpModeApiParams) VerifyPrivateIpModeApiRequest

	// Method available only for mocking purposes
	VerifyPrivateIpModeExecute(r VerifyPrivateIpModeApiRequest) (*PrivateIPMode, *http.Response, error)
}

// NetworkPeeringApiService NetworkPeeringApi service
type NetworkPeeringApiService service

type CreateGroupContainerApiRequest struct {
	ctx                    context.Context
	ApiService             NetworkPeeringApi
	groupId                string
	cloudProviderContainer *CloudProviderContainer
}

type CreateGroupContainerApiParams struct {
	GroupId                string
	CloudProviderContainer *CloudProviderContainer
}

func (a *NetworkPeeringApiService) CreateGroupContainerWithParams(ctx context.Context, args *CreateGroupContainerApiParams) CreateGroupContainerApiRequest {
	return CreateGroupContainerApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		cloudProviderContainer: args.CloudProviderContainer,
	}
}

func (r CreateGroupContainerApiRequest) Execute() (*CloudProviderContainer, *http.Response, error) {
	return r.ApiService.CreateGroupContainerExecute(r)
}

/*
CreateGroupContainer Create One Network Peering Container

Creates one new network peering container in the specified project. MongoDB Cloud can deploy Network Peering connections in a network peering container. GCP can have one container per project. AWS and Azure can have one container per cloud provider region.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupContainerApiRequest
*/
func (a *NetworkPeeringApiService) CreateGroupContainer(ctx context.Context, groupId string, cloudProviderContainer *CloudProviderContainer) CreateGroupContainerApiRequest {
	return CreateGroupContainerApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		cloudProviderContainer: cloudProviderContainer,
	}
}

// CreateGroupContainerExecute executes the request
//
//	@return CloudProviderContainer
func (a *NetworkPeeringApiService) CreateGroupContainerExecute(r CreateGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderContainer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.CreateGroupContainer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudProviderContainer == nil {
		return localVarReturnValue, nil, reportError("cloudProviderContainer is required and must be specified")
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
	localVarPostBody = r.cloudProviderContainer
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

type CreateGroupPeerApiRequest struct {
	ctx                                  context.Context
	ApiService                           NetworkPeeringApi
	groupId                              string
	baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings
}

type CreateGroupPeerApiParams struct {
	GroupId                              string
	BaseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings
}

func (a *NetworkPeeringApiService) CreateGroupPeerWithParams(ctx context.Context, args *CreateGroupPeerApiParams) CreateGroupPeerApiRequest {
	return CreateGroupPeerApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              args.GroupId,
		baseNetworkPeeringConnectionSettings: args.BaseNetworkPeeringConnectionSettings,
	}
}

func (r CreateGroupPeerApiRequest) Execute() (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	return r.ApiService.CreateGroupPeerExecute(r)
}

/*
CreateGroupPeer Create One Network Peering Connection

Creates one new network peering connection in the specified project. Network peering allows multiple cloud-hosted applications to securely connect to the same project. To learn more about considerations and prerequisites, see the Network Peering Documentation.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return CreateGroupPeerApiRequest
*/
func (a *NetworkPeeringApiService) CreateGroupPeer(ctx context.Context, groupId string, baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings) CreateGroupPeerApiRequest {
	return CreateGroupPeerApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              groupId,
		baseNetworkPeeringConnectionSettings: baseNetworkPeeringConnectionSettings,
	}
}

// CreateGroupPeerExecute executes the request
//
//	@return BaseNetworkPeeringConnectionSettings
func (a *NetworkPeeringApiService) CreateGroupPeerExecute(r CreateGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BaseNetworkPeeringConnectionSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.CreateGroupPeer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/peers"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.baseNetworkPeeringConnectionSettings == nil {
		return localVarReturnValue, nil, reportError("baseNetworkPeeringConnectionSettings is required and must be specified")
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
	localVarPostBody = r.baseNetworkPeeringConnectionSettings
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

type DeleteGroupContainerApiRequest struct {
	ctx         context.Context
	ApiService  NetworkPeeringApi
	groupId     string
	containerId string
}

type DeleteGroupContainerApiParams struct {
	GroupId     string
	ContainerId string
}

func (a *NetworkPeeringApiService) DeleteGroupContainerWithParams(ctx context.Context, args *DeleteGroupContainerApiParams) DeleteGroupContainerApiRequest {
	return DeleteGroupContainerApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		containerId: args.ContainerId,
	}
}

func (r DeleteGroupContainerApiRequest) Execute() (*http.Response, error) {
	return r.ApiService.DeleteGroupContainerExecute(r)
}

/*
DeleteGroupContainer Remove One Network Peering Container

Removes one network peering container in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container that you want to remove.
	@return DeleteGroupContainerApiRequest
*/
func (a *NetworkPeeringApiService) DeleteGroupContainer(ctx context.Context, groupId string, containerId string) DeleteGroupContainerApiRequest {
	return DeleteGroupContainerApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		containerId: containerId,
	}
}

// DeleteGroupContainerExecute executes the request
func (a *NetworkPeeringApiService) DeleteGroupContainerExecute(r DeleteGroupContainerApiRequest) (*http.Response, error) {
	var (
		localVarHTTPMethod = http.MethodDelete
		localVarPostBody   any
		formFiles          []formFile
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.DeleteGroupContainer")
	if err != nil {
		return nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers/{containerId}"
	if r.groupId == "" {
		return nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.containerId == "" {
		return nil, reportError("containerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"containerId"+"}", url.PathEscape(r.containerId), -1)

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

type DeleteGroupPeerApiRequest struct {
	ctx        context.Context
	ApiService NetworkPeeringApi
	groupId    string
	peerId     string
}

type DeleteGroupPeerApiParams struct {
	GroupId string
	PeerId  string
}

func (a *NetworkPeeringApiService) DeleteGroupPeerWithParams(ctx context.Context, args *DeleteGroupPeerApiParams) DeleteGroupPeerApiRequest {
	return DeleteGroupPeerApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		peerId:     args.PeerId,
	}
}

func (r DeleteGroupPeerApiRequest) Execute() (any, *http.Response, error) {
	return r.ApiService.DeleteGroupPeerExecute(r)
}

/*
DeleteGroupPeer Remove One Network Peering Connection

Removes one network peering connection in the specified project. If you remove the last network peering connection associated with a project, MongoDB Cloud also removes any AWS security groups from the project IP access list.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to delete.
	@return DeleteGroupPeerApiRequest
*/
func (a *NetworkPeeringApiService) DeleteGroupPeer(ctx context.Context, groupId string, peerId string) DeleteGroupPeerApiRequest {
	return DeleteGroupPeerApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		peerId:     peerId,
	}
}

// DeleteGroupPeerExecute executes the request
//
//	@return any
func (a *NetworkPeeringApiService) DeleteGroupPeerExecute(r DeleteGroupPeerApiRequest) (any, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue any
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.DeleteGroupPeer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/peers/{peerId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.peerId == "" {
		return localVarReturnValue, nil, reportError("peerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"peerId"+"}", url.PathEscape(r.peerId), -1)

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

type DisablePeeringApiRequest struct {
	ctx           context.Context
	ApiService    NetworkPeeringApi
	groupId       string
	privateIPMode *PrivateIPMode
}

type DisablePeeringApiParams struct {
	GroupId       string
	PrivateIPMode *PrivateIPMode
}

func (a *NetworkPeeringApiService) DisablePeeringWithParams(ctx context.Context, args *DisablePeeringApiParams) DisablePeeringApiRequest {
	return DisablePeeringApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       args.GroupId,
		privateIPMode: args.PrivateIPMode,
	}
}

func (r DisablePeeringApiRequest) Execute() (*PrivateIPMode, *http.Response, error) {
	return r.ApiService.DisablePeeringExecute(r)
}

/*
DisablePeering Disable Connect via Peering-Only Mode for One Project

Disables Connect via Peering Only mode for the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return DisablePeeringApiRequest

Deprecated
*/
func (a *NetworkPeeringApiService) DisablePeering(ctx context.Context, groupId string, privateIPMode *PrivateIPMode) DisablePeeringApiRequest {
	return DisablePeeringApiRequest{
		ApiService:    a,
		ctx:           ctx,
		groupId:       groupId,
		privateIPMode: privateIPMode,
	}
}

// DisablePeeringExecute executes the request
//
//	@return PrivateIPMode
//
// Deprecated
func (a *NetworkPeeringApiService) DisablePeeringExecute(r DisablePeeringApiRequest) (*PrivateIPMode, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PrivateIPMode
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.DisablePeering")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateIpMode"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.privateIPMode == nil {
		return localVarReturnValue, nil, reportError("privateIPMode is required and must be specified")
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
	localVarPostBody = r.privateIPMode
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

type GetGroupContainerApiRequest struct {
	ctx         context.Context
	ApiService  NetworkPeeringApi
	groupId     string
	containerId string
}

type GetGroupContainerApiParams struct {
	GroupId     string
	ContainerId string
}

func (a *NetworkPeeringApiService) GetGroupContainerWithParams(ctx context.Context, args *GetGroupContainerApiParams) GetGroupContainerApiRequest {
	return GetGroupContainerApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     args.GroupId,
		containerId: args.ContainerId,
	}
}

func (r GetGroupContainerApiRequest) Execute() (*CloudProviderContainer, *http.Response, error) {
	return r.ApiService.GetGroupContainerExecute(r)
}

/*
GetGroupContainer Return One Network Peering Container

Returns details about one network peering container in one specified project. Network peering containers contain network peering connections.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container.
	@return GetGroupContainerApiRequest
*/
func (a *NetworkPeeringApiService) GetGroupContainer(ctx context.Context, groupId string, containerId string) GetGroupContainerApiRequest {
	return GetGroupContainerApiRequest{
		ApiService:  a,
		ctx:         ctx,
		groupId:     groupId,
		containerId: containerId,
	}
}

// GetGroupContainerExecute executes the request
//
//	@return CloudProviderContainer
func (a *NetworkPeeringApiService) GetGroupContainerExecute(r GetGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderContainer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.GetGroupContainer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers/{containerId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.containerId == "" {
		return localVarReturnValue, nil, reportError("containerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"containerId"+"}", url.PathEscape(r.containerId), -1)

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

type GetGroupPeerApiRequest struct {
	ctx        context.Context
	ApiService NetworkPeeringApi
	groupId    string
	peerId     string
}

type GetGroupPeerApiParams struct {
	GroupId string
	PeerId  string
}

func (a *NetworkPeeringApiService) GetGroupPeerWithParams(ctx context.Context, args *GetGroupPeerApiParams) GetGroupPeerApiRequest {
	return GetGroupPeerApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
		peerId:     args.PeerId,
	}
}

func (r GetGroupPeerApiRequest) Execute() (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	return r.ApiService.GetGroupPeerExecute(r)
}

/*
GetGroupPeer Return One Network Peering Connection in One Project

Returns details about one specified network peering connection in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to retrieve.
	@return GetGroupPeerApiRequest
*/
func (a *NetworkPeeringApiService) GetGroupPeer(ctx context.Context, groupId string, peerId string) GetGroupPeerApiRequest {
	return GetGroupPeerApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
		peerId:     peerId,
	}
}

// GetGroupPeerExecute executes the request
//
//	@return BaseNetworkPeeringConnectionSettings
func (a *NetworkPeeringApiService) GetGroupPeerExecute(r GetGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BaseNetworkPeeringConnectionSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.GetGroupPeer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/peers/{peerId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.peerId == "" {
		return localVarReturnValue, nil, reportError("peerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"peerId"+"}", url.PathEscape(r.peerId), -1)

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

type ListGroupContainerAllApiRequest struct {
	ctx          context.Context
	ApiService   NetworkPeeringApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupContainerAllApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *NetworkPeeringApiService) ListGroupContainerAllWithParams(ctx context.Context, args *ListGroupContainerAllApiParams) ListGroupContainerAllApiRequest {
	return ListGroupContainerAllApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupContainerAllApiRequest) IncludeCount(includeCount bool) ListGroupContainerAllApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupContainerAllApiRequest) ItemsPerPage(itemsPerPage int) ListGroupContainerAllApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupContainerAllApiRequest) PageNum(pageNum int) ListGroupContainerAllApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupContainerAllApiRequest) Execute() (*PaginatedCloudProviderContainer, *http.Response, error) {
	return r.ApiService.ListGroupContainerAllExecute(r)
}

/*
ListGroupContainerAll Return All Network Peering Containers in One Project

Returns details about all network peering containers in the specified project. Network peering containers contain network peering connections.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupContainerAllApiRequest
*/
func (a *NetworkPeeringApiService) ListGroupContainerAll(ctx context.Context, groupId string) ListGroupContainerAllApiRequest {
	return ListGroupContainerAllApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupContainerAllExecute executes the request
//
//	@return PaginatedCloudProviderContainer
func (a *NetworkPeeringApiService) ListGroupContainerAllExecute(r ListGroupContainerAllApiRequest) (*PaginatedCloudProviderContainer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedCloudProviderContainer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.ListGroupContainerAll")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers/all"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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

type ListGroupContainersApiRequest struct {
	ctx          context.Context
	ApiService   NetworkPeeringApi
	groupId      string
	providerName *string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListGroupContainersApiParams struct {
	GroupId      string
	ProviderName *string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *NetworkPeeringApiService) ListGroupContainersWithParams(ctx context.Context, args *ListGroupContainersApiParams) ListGroupContainersApiRequest {
	return ListGroupContainersApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		providerName: args.ProviderName,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Cloud service provider that serves the desired network peering containers.
func (r ListGroupContainersApiRequest) ProviderName(providerName string) ListGroupContainersApiRequest {
	r.providerName = &providerName
	return r
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupContainersApiRequest) IncludeCount(includeCount bool) ListGroupContainersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupContainersApiRequest) ItemsPerPage(itemsPerPage int) ListGroupContainersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupContainersApiRequest) PageNum(pageNum int) ListGroupContainersApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListGroupContainersApiRequest) Execute() (*PaginatedCloudProviderContainer, *http.Response, error) {
	return r.ApiService.ListGroupContainersExecute(r)
}

/*
ListGroupContainers Return All Network Peering Containers in One Project for One Cloud Provider

Returns details about all network peering containers in the specified project for the specified cloud provider. If you do not specify the cloud provider, MongoDB Cloud returns details about all network peering containers in the project for Amazon Web Services (AWS).

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupContainersApiRequest
*/
func (a *NetworkPeeringApiService) ListGroupContainers(ctx context.Context, groupId string) ListGroupContainersApiRequest {
	return ListGroupContainersApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupContainersExecute executes the request
//
//	@return PaginatedCloudProviderContainer
func (a *NetworkPeeringApiService) ListGroupContainersExecute(r ListGroupContainersApiRequest) (*PaginatedCloudProviderContainer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedCloudProviderContainer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.ListGroupContainers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.providerName == nil {
		return localVarReturnValue, nil, reportError("providerName is required and must be specified")
	}

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
	parameterAddToHeaderOrQuery(localVarQueryParams, "providerName", r.providerName, "")
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

type ListGroupPeersApiRequest struct {
	ctx          context.Context
	ApiService   NetworkPeeringApi
	groupId      string
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
	providerName *string
}

type ListGroupPeersApiParams struct {
	GroupId      string
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
	ProviderName *string
}

func (a *NetworkPeeringApiService) ListGroupPeersWithParams(ctx context.Context, args *ListGroupPeersApiParams) ListGroupPeersApiRequest {
	return ListGroupPeersApiRequest{
		ApiService:   a,
		ctx:          ctx,
		groupId:      args.GroupId,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
		providerName: args.ProviderName,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListGroupPeersApiRequest) IncludeCount(includeCount bool) ListGroupPeersApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListGroupPeersApiRequest) ItemsPerPage(itemsPerPage int) ListGroupPeersApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListGroupPeersApiRequest) PageNum(pageNum int) ListGroupPeersApiRequest {
	r.pageNum = &pageNum
	return r
}

// Cloud service provider to use for this VPC peering connection.
func (r ListGroupPeersApiRequest) ProviderName(providerName string) ListGroupPeersApiRequest {
	r.providerName = &providerName
	return r
}

func (r ListGroupPeersApiRequest) Execute() (*PaginatedContainerPeer, *http.Response, error) {
	return r.ApiService.ListGroupPeersExecute(r)
}

/*
ListGroupPeers Return All Network Peering Connections in One Project

Returns details about all network peering connections in the specified project. Network peering allows multiple cloud-hosted applications to securely connect to the same project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return ListGroupPeersApiRequest
*/
func (a *NetworkPeeringApiService) ListGroupPeers(ctx context.Context, groupId string) ListGroupPeersApiRequest {
	return ListGroupPeersApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// ListGroupPeersExecute executes the request
//
//	@return PaginatedContainerPeer
func (a *NetworkPeeringApiService) ListGroupPeersExecute(r ListGroupPeersApiRequest) (*PaginatedContainerPeer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedContainerPeer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.ListGroupPeers")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/peers"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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
	if r.providerName != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "providerName", r.providerName, "")
	} else {
		var defaultValue string = "AWS"
		r.providerName = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "providerName", r.providerName, "")
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

type UpdateGroupContainerApiRequest struct {
	ctx                    context.Context
	ApiService             NetworkPeeringApi
	groupId                string
	containerId            string
	cloudProviderContainer *CloudProviderContainer
}

type UpdateGroupContainerApiParams struct {
	GroupId                string
	ContainerId            string
	CloudProviderContainer *CloudProviderContainer
}

func (a *NetworkPeeringApiService) UpdateGroupContainerWithParams(ctx context.Context, args *UpdateGroupContainerApiParams) UpdateGroupContainerApiRequest {
	return UpdateGroupContainerApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                args.GroupId,
		containerId:            args.ContainerId,
		cloudProviderContainer: args.CloudProviderContainer,
	}
}

func (r UpdateGroupContainerApiRequest) Execute() (*CloudProviderContainer, *http.Response, error) {
	return r.ApiService.UpdateGroupContainerExecute(r)
}

/*
UpdateGroupContainer Update One Network Peering Container

Updates the network details and labels of one specified network peering container in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param containerId Unique 24-hexadecimal digit string that identifies the MongoDB Cloud network container that you want to remove.
	@return UpdateGroupContainerApiRequest
*/
func (a *NetworkPeeringApiService) UpdateGroupContainer(ctx context.Context, groupId string, containerId string, cloudProviderContainer *CloudProviderContainer) UpdateGroupContainerApiRequest {
	return UpdateGroupContainerApiRequest{
		ApiService:             a,
		ctx:                    ctx,
		groupId:                groupId,
		containerId:            containerId,
		cloudProviderContainer: cloudProviderContainer,
	}
}

// UpdateGroupContainerExecute executes the request
//
//	@return CloudProviderContainer
func (a *NetworkPeeringApiService) UpdateGroupContainerExecute(r UpdateGroupContainerApiRequest) (*CloudProviderContainer, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CloudProviderContainer
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.UpdateGroupContainer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/containers/{containerId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.containerId == "" {
		return localVarReturnValue, nil, reportError("containerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"containerId"+"}", url.PathEscape(r.containerId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.cloudProviderContainer == nil {
		return localVarReturnValue, nil, reportError("cloudProviderContainer is required and must be specified")
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
	localVarPostBody = r.cloudProviderContainer
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

type UpdateGroupPeerApiRequest struct {
	ctx                                  context.Context
	ApiService                           NetworkPeeringApi
	groupId                              string
	peerId                               string
	baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings
}

type UpdateGroupPeerApiParams struct {
	GroupId                              string
	PeerId                               string
	BaseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings
}

func (a *NetworkPeeringApiService) UpdateGroupPeerWithParams(ctx context.Context, args *UpdateGroupPeerApiParams) UpdateGroupPeerApiRequest {
	return UpdateGroupPeerApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              args.GroupId,
		peerId:                               args.PeerId,
		baseNetworkPeeringConnectionSettings: args.BaseNetworkPeeringConnectionSettings,
	}
}

func (r UpdateGroupPeerApiRequest) Execute() (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	return r.ApiService.UpdateGroupPeerExecute(r)
}

/*
UpdateGroupPeer Update One Network Peering Connection

Updates one specified network peering connection in the specified project.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@param peerId Unique 24-hexadecimal digit string that identifies the network peering connection that you want to update.
	@return UpdateGroupPeerApiRequest
*/
func (a *NetworkPeeringApiService) UpdateGroupPeer(ctx context.Context, groupId string, peerId string, baseNetworkPeeringConnectionSettings *BaseNetworkPeeringConnectionSettings) UpdateGroupPeerApiRequest {
	return UpdateGroupPeerApiRequest{
		ApiService:                           a,
		ctx:                                  ctx,
		groupId:                              groupId,
		peerId:                               peerId,
		baseNetworkPeeringConnectionSettings: baseNetworkPeeringConnectionSettings,
	}
}

// UpdateGroupPeerExecute executes the request
//
//	@return BaseNetworkPeeringConnectionSettings
func (a *NetworkPeeringApiService) UpdateGroupPeerExecute(r UpdateGroupPeerApiRequest) (*BaseNetworkPeeringConnectionSettings, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPatch
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BaseNetworkPeeringConnectionSettings
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.UpdateGroupPeer")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/peers/{peerId}"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)
	if r.peerId == "" {
		return localVarReturnValue, nil, reportError("peerId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"peerId"+"}", url.PathEscape(r.peerId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.baseNetworkPeeringConnectionSettings == nil {
		return localVarReturnValue, nil, reportError("baseNetworkPeeringConnectionSettings is required and must be specified")
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
	localVarPostBody = r.baseNetworkPeeringConnectionSettings
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

type VerifyPrivateIpModeApiRequest struct {
	ctx        context.Context
	ApiService NetworkPeeringApi
	groupId    string
}

type VerifyPrivateIpModeApiParams struct {
	GroupId string
}

func (a *NetworkPeeringApiService) VerifyPrivateIpModeWithParams(ctx context.Context, args *VerifyPrivateIpModeApiParams) VerifyPrivateIpModeApiRequest {
	return VerifyPrivateIpModeApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    args.GroupId,
	}
}

func (r VerifyPrivateIpModeApiRequest) Execute() (*PrivateIPMode, *http.Response, error) {
	return r.ApiService.VerifyPrivateIpModeExecute(r)
}

/*
VerifyPrivateIpMode Verify Connect via Peering-Only Mode for One Project

Verifies if someone set the specified project to **Connect via Peering Only** mode.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param groupId Unique 24-hexadecimal digit string that identifies your project. Use the [/groups](#tag/Projects/operation/listProjects) endpoint to retrieve all projects to which the authenticated user has access.  **NOTE**: Groups and projects are synonymous terms. Your group id is the same as your project id. For existing groups, your group/project id remains the same. The resource and corresponding endpoints use the term groups.
	@return VerifyPrivateIpModeApiRequest

Deprecated
*/
func (a *NetworkPeeringApiService) VerifyPrivateIpMode(ctx context.Context, groupId string) VerifyPrivateIpModeApiRequest {
	return VerifyPrivateIpModeApiRequest{
		ApiService: a,
		ctx:        ctx,
		groupId:    groupId,
	}
}

// VerifyPrivateIpModeExecute executes the request
//
//	@return PrivateIPMode
//
// Deprecated
func (a *NetworkPeeringApiService) VerifyPrivateIpModeExecute(r VerifyPrivateIpModeApiRequest) (*PrivateIPMode, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PrivateIPMode
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "NetworkPeeringApiService.VerifyPrivateIpMode")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/groups/{groupId}/privateIpMode"
	if r.groupId == "" {
		return localVarReturnValue, nil, reportError("groupId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"groupId"+"}", url.PathEscape(r.groupId), -1)

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
