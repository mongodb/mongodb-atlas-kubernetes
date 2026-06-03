// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type InvoicesApi interface {

	/*
		CreateCostExplorerProcess Create One Cost Explorer Query Process

		Creates a query process within the Cost Explorer for the given parameters. A token is returned that can be used to poll the status of the query and eventually retrieve the results.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param costExplorerFilterRequestBody Filter parameters for the Cost Explorer query.
		@return CreateCostExplorerProcessApiRequest
	*/
	CreateCostExplorerProcess(ctx context.Context, orgId string, costExplorerFilterRequestBody *CostExplorerFilterRequestBody) CreateCostExplorerProcessApiRequest
	/*
		CreateCostExplorerProcess Create One Cost Explorer Query Process


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param CreateCostExplorerProcessApiParams - Parameters for the request
		@return CreateCostExplorerProcessApiRequest
	*/
	CreateCostExplorerProcessWithParams(ctx context.Context, args *CreateCostExplorerProcessApiParams) CreateCostExplorerProcessApiRequest

	// Method available only for mocking purposes
	CreateCostExplorerProcessExecute(r CreateCostExplorerProcessApiRequest) (*CostExplorerFilterResponse, *http.Response, error)

	/*
		GetCostExplorerUsage Return Usage Details for One Cost Explorer Query

		Returns the usage details for a Cost Explorer query, if the query is finished and the data is ready to be viewed. If the data is not ready, a 'processing' response will indicate that another request should be sent later to view the data.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param token Unique 64 digit string that identifies the Cost Explorer query.
		@return GetCostExplorerUsageApiRequest
	*/
	GetCostExplorerUsage(ctx context.Context, orgId string, token string) GetCostExplorerUsageApiRequest
	/*
		GetCostExplorerUsage Return Usage Details for One Cost Explorer Query


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetCostExplorerUsageApiParams - Parameters for the request
		@return GetCostExplorerUsageApiRequest
	*/
	GetCostExplorerUsageWithParams(ctx context.Context, args *GetCostExplorerUsageApiParams) GetCostExplorerUsageApiRequest

	// Method available only for mocking purposes
	GetCostExplorerUsageExecute(r GetCostExplorerUsageApiRequest) (any, *http.Response, error)

	/*
			GetInvoice Return One Invoice for One Organization

			Returns one invoice that MongoDB issued to the specified organization. A unique 24-hexadecimal digit string identifies the invoice. You can choose to receive this invoice in JSON or CSV format. If you have a cross-organization setup, you can query for a linked invoice if you have the Organization Billing Admin or Organization Owner role.
		To compute the total owed amount of the invoice - sum up total owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
			@return GetInvoiceApiRequest
	*/
	GetInvoice(ctx context.Context, orgId string, invoiceId string) GetInvoiceApiRequest
	/*
		GetInvoice Return One Invoice for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetInvoiceApiParams - Parameters for the request
		@return GetInvoiceApiRequest
	*/
	GetInvoiceWithParams(ctx context.Context, args *GetInvoiceApiParams) GetInvoiceApiRequest

	// Method available only for mocking purposes
	GetInvoiceExecute(r GetInvoiceApiRequest) (*BillingInvoice, *http.Response, error)

	/*
			GetInvoiceCsv Return One Invoice as CSV for One Organization

			Returns one invoice that MongoDB issued to the specified organization in CSV format. A unique 24-hexadecimal digit string identifies the invoice. If you have a cross-organization setup, you can query for a linked invoice if you have the Organization Billing Admin or Organization Owner Role.
		 To compute the total owed amount of the invoice - sum up total owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
			@return GetInvoiceCsvApiRequest
	*/
	GetInvoiceCsv(ctx context.Context, orgId string, invoiceId string) GetInvoiceCsvApiRequest
	/*
		GetInvoiceCsv Return One Invoice as CSV for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetInvoiceCsvApiParams - Parameters for the request
		@return GetInvoiceCsvApiRequest
	*/
	GetInvoiceCsvWithParams(ctx context.Context, args *GetInvoiceCsvApiParams) GetInvoiceCsvApiRequest

	// Method available only for mocking purposes
	GetInvoiceCsvExecute(r GetInvoiceCsvApiRequest) (string, *http.Response, error)

	/*
		GetSku Return One Stock Keeping Unit

		Returns details about a single SKU (Stock Keeping Unit) by its identifier. SKUs represent different products and services offered by MongoDB.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param skuId Unique identifier of the SKU to retrieve.
		@return GetSkuApiRequest
	*/
	GetSku(ctx context.Context, skuId string) GetSkuApiRequest
	/*
		GetSku Return One Stock Keeping Unit


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param GetSkuApiParams - Parameters for the request
		@return GetSkuApiRequest
	*/
	GetSkuWithParams(ctx context.Context, args *GetSkuApiParams) GetSkuApiRequest

	// Method available only for mocking purposes
	GetSkuExecute(r GetSkuApiRequest) (*SkuResponse, *http.Response, error)

	/*
		ListInvoicePending Return All Pending Invoices for One Organization

		Returns all invoices accruing charges for the current billing cycle for the specified organization. If you have a cross-organization setup, you can view linked invoices if you have the Organization Billing Admin or Organization Owner Role.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@return ListInvoicePendingApiRequest
	*/
	ListInvoicePending(ctx context.Context, orgId string) ListInvoicePendingApiRequest
	/*
		ListInvoicePending Return All Pending Invoices for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListInvoicePendingApiParams - Parameters for the request
		@return ListInvoicePendingApiRequest
	*/
	ListInvoicePendingWithParams(ctx context.Context, args *ListInvoicePendingApiParams) ListInvoicePendingApiRequest

	// Method available only for mocking purposes
	ListInvoicePendingExecute(r ListInvoicePendingApiRequest) (*PaginatedApiInvoice, *http.Response, error)

	/*
			ListInvoices Return All Invoices for One Organization

			Returns all invoices that MongoDB issued to the specified organization. This list includes all invoices regardless of invoice status. If you have a cross-organization setup, you can view linked invoices if you have the Organization Billing Admin or Organization Owner role.
		To compute the total owed amount of the invoices - sum up total owed of each invoice. It could be computed as a sum of owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

			@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
			@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
			@return ListInvoicesApiRequest
	*/
	ListInvoices(ctx context.Context, orgId string) ListInvoicesApiRequest
	/*
		ListInvoices Return All Invoices for One Organization


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListInvoicesApiParams - Parameters for the request
		@return ListInvoicesApiRequest
	*/
	ListInvoicesWithParams(ctx context.Context, args *ListInvoicesApiParams) ListInvoicesApiRequest

	// Method available only for mocking purposes
	ListInvoicesExecute(r ListInvoicesApiRequest) (*PaginatedApiInvoiceMetadata, *http.Response, error)

	/*
		ListSkus Return All Stock Keeping Units

		Returns all available SKUs (Stock Keeping Units) that can appear on MongoDB invoices. SKUs represent different products and services offered by MongoDB.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@return ListSkusApiRequest
	*/
	ListSkus(ctx context.Context) ListSkusApiRequest
	/*
		ListSkus Return All Stock Keeping Units


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param ListSkusApiParams - Parameters for the request
		@return ListSkusApiRequest
	*/
	ListSkusWithParams(ctx context.Context, args *ListSkusApiParams) ListSkusApiRequest

	// Method available only for mocking purposes
	ListSkusExecute(r ListSkusApiRequest) (*PaginatedApiSKU, *http.Response, error)

	/*
		SearchInvoiceLineItems Return All Line Items for One Invoice by Invoice ID

		Query the `lineItems` of the specified invoice and return the result JSON. A unique 24-hexadecimal digit string identifies the invoice.

		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
		@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
		@param apiPublicUsageDetailsQueryRequest Filter parameters for the `lineItems` query. Send a request with an empty JSON body to retrieve all line items for a given `invoiceID` without applying any filters.
		@return SearchInvoiceLineItemsApiRequest
	*/
	SearchInvoiceLineItems(ctx context.Context, orgId string, invoiceId string, apiPublicUsageDetailsQueryRequest *ApiPublicUsageDetailsQueryRequest) SearchInvoiceLineItemsApiRequest
	/*
		SearchInvoiceLineItems Return All Line Items for One Invoice by Invoice ID


		@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
		@param SearchInvoiceLineItemsApiParams - Parameters for the request
		@return SearchInvoiceLineItemsApiRequest
	*/
	SearchInvoiceLineItemsWithParams(ctx context.Context, args *SearchInvoiceLineItemsApiParams) SearchInvoiceLineItemsApiRequest

	// Method available only for mocking purposes
	SearchInvoiceLineItemsExecute(r SearchInvoiceLineItemsApiRequest) (*PaginatedPublicApiUsageDetailsLineItem, *http.Response, error)
}

// InvoicesApiService InvoicesApi service
type InvoicesApiService service

type CreateCostExplorerProcessApiRequest struct {
	ctx                           context.Context
	ApiService                    InvoicesApi
	orgId                         string
	costExplorerFilterRequestBody *CostExplorerFilterRequestBody
}

type CreateCostExplorerProcessApiParams struct {
	OrgId                         string
	CostExplorerFilterRequestBody *CostExplorerFilterRequestBody
}

func (a *InvoicesApiService) CreateCostExplorerProcessWithParams(ctx context.Context, args *CreateCostExplorerProcessApiParams) CreateCostExplorerProcessApiRequest {
	return CreateCostExplorerProcessApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         args.OrgId,
		costExplorerFilterRequestBody: args.CostExplorerFilterRequestBody,
	}
}

func (r CreateCostExplorerProcessApiRequest) Execute() (*CostExplorerFilterResponse, *http.Response, error) {
	return r.ApiService.CreateCostExplorerProcessExecute(r)
}

/*
CreateCostExplorerProcess Create One Cost Explorer Query Process

Creates a query process within the Cost Explorer for the given parameters. A token is returned that can be used to poll the status of the query and eventually retrieve the results.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return CreateCostExplorerProcessApiRequest
*/
func (a *InvoicesApiService) CreateCostExplorerProcess(ctx context.Context, orgId string, costExplorerFilterRequestBody *CostExplorerFilterRequestBody) CreateCostExplorerProcessApiRequest {
	return CreateCostExplorerProcessApiRequest{
		ApiService:                    a,
		ctx:                           ctx,
		orgId:                         orgId,
		costExplorerFilterRequestBody: costExplorerFilterRequestBody,
	}
}

// CreateCostExplorerProcessExecute executes the request
//
//	@return CostExplorerFilterResponse
func (a *InvoicesApiService) CreateCostExplorerProcessExecute(r CreateCostExplorerProcessApiRequest) (*CostExplorerFilterResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *CostExplorerFilterResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.CreateCostExplorerProcess")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/billing/costExplorer/usage"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.costExplorerFilterRequestBody == nil {
		return localVarReturnValue, nil, reportError("costExplorerFilterRequestBody is required and must be specified")
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
	localVarPostBody = r.costExplorerFilterRequestBody
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

type GetCostExplorerUsageApiRequest struct {
	ctx        context.Context
	ApiService InvoicesApi
	orgId      string
	token      string
}

type GetCostExplorerUsageApiParams struct {
	OrgId string
	Token string
}

func (a *InvoicesApiService) GetCostExplorerUsageWithParams(ctx context.Context, args *GetCostExplorerUsageApiParams) GetCostExplorerUsageApiRequest {
	return GetCostExplorerUsageApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		token:      args.Token,
	}
}

func (r GetCostExplorerUsageApiRequest) Execute() (any, *http.Response, error) {
	return r.ApiService.GetCostExplorerUsageExecute(r)
}

/*
GetCostExplorerUsage Return Usage Details for One Cost Explorer Query

Returns the usage details for a Cost Explorer query, if the query is finished and the data is ready to be viewed. If the data is not ready, a 'processing' response will indicate that another request should be sent later to view the data.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param token Unique 64 digit string that identifies the Cost Explorer query.
	@return GetCostExplorerUsageApiRequest
*/
func (a *InvoicesApiService) GetCostExplorerUsage(ctx context.Context, orgId string, token string) GetCostExplorerUsageApiRequest {
	return GetCostExplorerUsageApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		token:      token,
	}
}

// GetCostExplorerUsageExecute executes the request
//
//	@return any
func (a *InvoicesApiService) GetCostExplorerUsageExecute(r GetCostExplorerUsageApiRequest) (any, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue any
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.GetCostExplorerUsage")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/billing/costExplorer/usage/{token}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.token == "" {
		return localVarReturnValue, nil, reportError("token is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"token"+"}", url.PathEscape(r.token), -1)

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

type GetInvoiceApiRequest struct {
	ctx        context.Context
	ApiService InvoicesApi
	orgId      string
	invoiceId  string
}

type GetInvoiceApiParams struct {
	OrgId     string
	InvoiceId string
}

func (a *InvoicesApiService) GetInvoiceWithParams(ctx context.Context, args *GetInvoiceApiParams) GetInvoiceApiRequest {
	return GetInvoiceApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		invoiceId:  args.InvoiceId,
	}
}

func (r GetInvoiceApiRequest) Execute() (*BillingInvoice, *http.Response, error) {
	return r.ApiService.GetInvoiceExecute(r)
}

/*
GetInvoice Return One Invoice for One Organization

Returns one invoice that MongoDB issued to the specified organization. A unique 24-hexadecimal digit string identifies the invoice. You can choose to receive this invoice in JSON or CSV format. If you have a cross-organization setup, you can query for a linked invoice if you have the Organization Billing Admin or Organization Owner role.
To compute the total owed amount of the invoice - sum up total owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
	@return GetInvoiceApiRequest
*/
func (a *InvoicesApiService) GetInvoice(ctx context.Context, orgId string, invoiceId string) GetInvoiceApiRequest {
	return GetInvoiceApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		invoiceId:  invoiceId,
	}
}

// GetInvoiceExecute executes the request
//
//	@return BillingInvoice
func (a *InvoicesApiService) GetInvoiceExecute(r GetInvoiceApiRequest) (*BillingInvoice, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *BillingInvoice
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.GetInvoice")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invoices/{invoiceId}"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invoiceId == "" {
		return localVarReturnValue, nil, reportError("invoiceId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invoiceId"+"}", url.PathEscape(r.invoiceId), -1)

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

type GetInvoiceCsvApiRequest struct {
	ctx        context.Context
	ApiService InvoicesApi
	orgId      string
	invoiceId  string
}

type GetInvoiceCsvApiParams struct {
	OrgId     string
	InvoiceId string
}

func (a *InvoicesApiService) GetInvoiceCsvWithParams(ctx context.Context, args *GetInvoiceCsvApiParams) GetInvoiceCsvApiRequest {
	return GetInvoiceCsvApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
		invoiceId:  args.InvoiceId,
	}
}

func (r GetInvoiceCsvApiRequest) Execute() (string, *http.Response, error) {
	return r.ApiService.GetInvoiceCsvExecute(r)
}

/*
GetInvoiceCsv Return One Invoice as CSV for One Organization

Returns one invoice that MongoDB issued to the specified organization in CSV format. A unique 24-hexadecimal digit string identifies the invoice. If you have a cross-organization setup, you can query for a linked invoice if you have the Organization Billing Admin or Organization Owner Role.

	To compute the total owed amount of the invoice - sum up total owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
	@return GetInvoiceCsvApiRequest
*/
func (a *InvoicesApiService) GetInvoiceCsv(ctx context.Context, orgId string, invoiceId string) GetInvoiceCsvApiRequest {
	return GetInvoiceCsvApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
		invoiceId:  invoiceId,
	}
}

// GetInvoiceCsvExecute executes the request
//
//	@return string
func (a *InvoicesApiService) GetInvoiceCsvExecute(r GetInvoiceCsvApiRequest) (string, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.GetInvoiceCsv")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invoices/{invoiceId}/csv"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invoiceId == "" {
		return localVarReturnValue, nil, reportError("invoiceId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invoiceId"+"}", url.PathEscape(r.invoiceId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2023-01-01+csv"}

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

type GetSkuApiRequest struct {
	ctx        context.Context
	ApiService InvoicesApi
	skuId      string
}

type GetSkuApiParams struct {
	SkuId string
}

func (a *InvoicesApiService) GetSkuWithParams(ctx context.Context, args *GetSkuApiParams) GetSkuApiRequest {
	return GetSkuApiRequest{
		ApiService: a,
		ctx:        ctx,
		skuId:      args.SkuId,
	}
}

func (r GetSkuApiRequest) Execute() (*SkuResponse, *http.Response, error) {
	return r.ApiService.GetSkuExecute(r)
}

/*
GetSku Return One Stock Keeping Unit

Returns details about a single SKU (Stock Keeping Unit) by its identifier. SKUs represent different products and services offered by MongoDB.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param skuId Unique identifier of the SKU to retrieve.
	@return GetSkuApiRequest
*/
func (a *InvoicesApiService) GetSku(ctx context.Context, skuId string) GetSkuApiRequest {
	return GetSkuApiRequest{
		ApiService: a,
		ctx:        ctx,
		skuId:      skuId,
	}
}

// GetSkuExecute executes the request
//
//	@return SkuResponse
func (a *InvoicesApiService) GetSkuExecute(r GetSkuApiRequest) (*SkuResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *SkuResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.GetSku")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/skus/{skuId}"
	if r.skuId == "" {
		return localVarReturnValue, nil, reportError("skuId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"skuId"+"}", url.PathEscape(r.skuId), -1)

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

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

type ListInvoicePendingApiRequest struct {
	ctx        context.Context
	ApiService InvoicesApi
	orgId      string
}

type ListInvoicePendingApiParams struct {
	OrgId string
}

func (a *InvoicesApiService) ListInvoicePendingWithParams(ctx context.Context, args *ListInvoicePendingApiParams) ListInvoicePendingApiRequest {
	return ListInvoicePendingApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      args.OrgId,
	}
}

func (r ListInvoicePendingApiRequest) Execute() (*PaginatedApiInvoice, *http.Response, error) {
	return r.ApiService.ListInvoicePendingExecute(r)
}

/*
ListInvoicePending Return All Pending Invoices for One Organization

Returns all invoices accruing charges for the current billing cycle for the specified organization. If you have a cross-organization setup, you can view linked invoices if you have the Organization Billing Admin or Organization Owner Role.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListInvoicePendingApiRequest
*/
func (a *InvoicesApiService) ListInvoicePending(ctx context.Context, orgId string) ListInvoicePendingApiRequest {
	return ListInvoicePendingApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListInvoicePendingExecute executes the request
//
//	@return PaginatedApiInvoice
func (a *InvoicesApiService) ListInvoicePendingExecute(r ListInvoicePendingApiRequest) (*PaginatedApiInvoice, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiInvoice
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.ListInvoicePending")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invoices/pending"
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

type ListInvoicesApiRequest struct {
	ctx                context.Context
	ApiService         InvoicesApi
	orgId              string
	includeCount       *bool
	itemsPerPage       *int
	pageNum            *int
	viewLinkedInvoices *bool
	statusNames        *[]string
	fromDate           *string
	toDate             *string
	sortBy             *string
	orderBy            *string
}

type ListInvoicesApiParams struct {
	OrgId              string
	IncludeCount       *bool
	ItemsPerPage       *int
	PageNum            *int
	ViewLinkedInvoices *bool
	StatusNames        *[]string
	FromDate           *string
	ToDate             *string
	SortBy             *string
	OrderBy            *string
}

func (a *InvoicesApiService) ListInvoicesWithParams(ctx context.Context, args *ListInvoicesApiParams) ListInvoicesApiRequest {
	return ListInvoicesApiRequest{
		ApiService:         a,
		ctx:                ctx,
		orgId:              args.OrgId,
		includeCount:       args.IncludeCount,
		itemsPerPage:       args.ItemsPerPage,
		pageNum:            args.PageNum,
		viewLinkedInvoices: args.ViewLinkedInvoices,
		statusNames:        args.StatusNames,
		fromDate:           args.FromDate,
		toDate:             args.ToDate,
		sortBy:             args.SortBy,
		orderBy:            args.OrderBy,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListInvoicesApiRequest) IncludeCount(includeCount bool) ListInvoicesApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListInvoicesApiRequest) ItemsPerPage(itemsPerPage int) ListInvoicesApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListInvoicesApiRequest) PageNum(pageNum int) ListInvoicesApiRequest {
	r.pageNum = &pageNum
	return r
}

// Flag that indicates whether to return linked invoices in the &#x60;linkedInvoices&#x60; field.
func (r ListInvoicesApiRequest) ViewLinkedInvoices(viewLinkedInvoices bool) ListInvoicesApiRequest {
	r.viewLinkedInvoices = &viewLinkedInvoices
	return r
}

// Statuses of the invoice to be retrieved. Omit to return invoices of all statuses.
func (r ListInvoicesApiRequest) StatusNames(statusNames []string) ListInvoicesApiRequest {
	r.statusNames = &statusNames
	return r
}

// Retrieve the invoices the &#x60;startDates&#x60; of which are greater than or equal to the &#x60;fromDate&#x60;. If omit, the invoices return will go back to earliest &#x60;startDate&#x60;.
func (r ListInvoicesApiRequest) FromDate(fromDate string) ListInvoicesApiRequest {
	r.fromDate = &fromDate
	return r
}

// Retrieve the invoices the &#x60;endDates&#x60; of which are smaller than or equal to the &#x60;toDate&#x60;. If omit, the invoices return will go further to latest &#x60;endDate&#x60;.
func (r ListInvoicesApiRequest) ToDate(toDate string) ListInvoicesApiRequest {
	r.toDate = &toDate
	return r
}

// Field used to sort the returned invoices by. Use in combination with &#x60;orderBy&#x60; parameter to control the order of the result.
func (r ListInvoicesApiRequest) SortBy(sortBy string) ListInvoicesApiRequest {
	r.sortBy = &sortBy
	return r
}

// Field used to order the returned invoices by. Use in combination of &#x60;sortBy&#x60; parameter to control the order of the result.
func (r ListInvoicesApiRequest) OrderBy(orderBy string) ListInvoicesApiRequest {
	r.orderBy = &orderBy
	return r
}

func (r ListInvoicesApiRequest) Execute() (*PaginatedApiInvoiceMetadata, *http.Response, error) {
	return r.ApiService.ListInvoicesExecute(r)
}

/*
ListInvoices Return All Invoices for One Organization

Returns all invoices that MongoDB issued to the specified organization. This list includes all invoices regardless of invoice status. If you have a cross-organization setup, you can view linked invoices if you have the Organization Billing Admin or Organization Owner role.
To compute the total owed amount of the invoices - sum up total owed of each invoice. It could be computed as a sum of owed amount of each payment included into the invoice. To compute payment's owed amount - use formula `totalBilledCents` * `unitPrice` + `salesTax` - `startingBalanceCents`.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@return ListInvoicesApiRequest
*/
func (a *InvoicesApiService) ListInvoices(ctx context.Context, orgId string) ListInvoicesApiRequest {
	return ListInvoicesApiRequest{
		ApiService: a,
		ctx:        ctx,
		orgId:      orgId,
	}
}

// ListInvoicesExecute executes the request
//
//	@return PaginatedApiInvoiceMetadata
func (a *InvoicesApiService) ListInvoicesExecute(r ListInvoicesApiRequest) (*PaginatedApiInvoiceMetadata, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiInvoiceMetadata
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.ListInvoices")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invoices"
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
	if r.viewLinkedInvoices != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "viewLinkedInvoices", r.viewLinkedInvoices, "")
	} else {
		var defaultValue bool = true
		r.viewLinkedInvoices = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "viewLinkedInvoices", r.viewLinkedInvoices, "")
	}
	if r.statusNames != nil {
		t := *r.statusNames
		// Workaround for unused import
		_ = reflect.Append
		parameterAddToHeaderOrQuery(localVarQueryParams, "statusNames", t, "multi")

	}
	if r.fromDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "fromDate", r.fromDate, "")
	}
	if r.toDate != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "toDate", r.toDate, "")
	}
	if r.sortBy != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "sortBy", r.sortBy, "")
	} else {
		var defaultValue string = "END_DATE"
		r.sortBy = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "sortBy", r.sortBy, "")
	}
	if r.orderBy != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orderBy", r.orderBy, "")
	} else {
		var defaultValue string = "desc"
		r.orderBy = &defaultValue
		parameterAddToHeaderOrQuery(localVarQueryParams, "orderBy", r.orderBy, "")
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

type ListSkusApiRequest struct {
	ctx          context.Context
	ApiService   InvoicesApi
	includeCount *bool
	itemsPerPage *int
	pageNum      *int
}

type ListSkusApiParams struct {
	IncludeCount *bool
	ItemsPerPage *int
	PageNum      *int
}

func (a *InvoicesApiService) ListSkusWithParams(ctx context.Context, args *ListSkusApiParams) ListSkusApiRequest {
	return ListSkusApiRequest{
		ApiService:   a,
		ctx:          ctx,
		includeCount: args.IncludeCount,
		itemsPerPage: args.ItemsPerPage,
		pageNum:      args.PageNum,
	}
}

// Flag that indicates whether the response returns the total number of items (&#x60;totalCount&#x60;) in the response.
func (r ListSkusApiRequest) IncludeCount(includeCount bool) ListSkusApiRequest {
	r.includeCount = &includeCount
	return r
}

// Number of items that the response returns per page.
func (r ListSkusApiRequest) ItemsPerPage(itemsPerPage int) ListSkusApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r ListSkusApiRequest) PageNum(pageNum int) ListSkusApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r ListSkusApiRequest) Execute() (*PaginatedApiSKU, *http.Response, error) {
	return r.ApiService.ListSkusExecute(r)
}

/*
ListSkus Return All Stock Keeping Units

Returns all available SKUs (Stock Keeping Units) that can appear on MongoDB invoices. SKUs represent different products and services offered by MongoDB.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ListSkusApiRequest
*/
func (a *InvoicesApiService) ListSkus(ctx context.Context) ListSkusApiRequest {
	return ListSkusApiRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// ListSkusExecute executes the request
//
//	@return PaginatedApiSKU
func (a *InvoicesApiService) ListSkusExecute(r ListSkusApiRequest) (*PaginatedApiSKU, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedApiSKU
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.ListSkus")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/skus"

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
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2025-03-12+json"}

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

type SearchInvoiceLineItemsApiRequest struct {
	ctx                               context.Context
	ApiService                        InvoicesApi
	orgId                             string
	invoiceId                         string
	apiPublicUsageDetailsQueryRequest *ApiPublicUsageDetailsQueryRequest
	itemsPerPage                      *int
	pageNum                           *int
}

type SearchInvoiceLineItemsApiParams struct {
	OrgId                             string
	InvoiceId                         string
	ApiPublicUsageDetailsQueryRequest *ApiPublicUsageDetailsQueryRequest
	ItemsPerPage                      *int
	PageNum                           *int
}

func (a *InvoicesApiService) SearchInvoiceLineItemsWithParams(ctx context.Context, args *SearchInvoiceLineItemsApiParams) SearchInvoiceLineItemsApiRequest {
	return SearchInvoiceLineItemsApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		orgId:                             args.OrgId,
		invoiceId:                         args.InvoiceId,
		apiPublicUsageDetailsQueryRequest: args.ApiPublicUsageDetailsQueryRequest,
		itemsPerPage:                      args.ItemsPerPage,
		pageNum:                           args.PageNum,
	}
}

// Number of items that the response returns per page.
func (r SearchInvoiceLineItemsApiRequest) ItemsPerPage(itemsPerPage int) SearchInvoiceLineItemsApiRequest {
	r.itemsPerPage = &itemsPerPage
	return r
}

// Number of the page that displays the current set of the total objects that the response returns.
func (r SearchInvoiceLineItemsApiRequest) PageNum(pageNum int) SearchInvoiceLineItemsApiRequest {
	r.pageNum = &pageNum
	return r
}

func (r SearchInvoiceLineItemsApiRequest) Execute() (*PaginatedPublicApiUsageDetailsLineItem, *http.Response, error) {
	return r.ApiService.SearchInvoiceLineItemsExecute(r)
}

/*
SearchInvoiceLineItems Return All Line Items for One Invoice by Invoice ID

Query the `lineItems` of the specified invoice and return the result JSON. A unique 24-hexadecimal digit string identifies the invoice.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param orgId Unique 24-hexadecimal digit string that identifies the organization that contains your projects. Use the [`/orgs`](#tag/Organizations/operation/listOrganizations) endpoint to retrieve all organizations to which the authenticated user has access.
	@param invoiceId Unique 24-hexadecimal digit string that identifies the invoice submitted to the specified organization. Charges typically post the next day.
	@return SearchInvoiceLineItemsApiRequest
*/
func (a *InvoicesApiService) SearchInvoiceLineItems(ctx context.Context, orgId string, invoiceId string, apiPublicUsageDetailsQueryRequest *ApiPublicUsageDetailsQueryRequest) SearchInvoiceLineItemsApiRequest {
	return SearchInvoiceLineItemsApiRequest{
		ApiService:                        a,
		ctx:                               ctx,
		orgId:                             orgId,
		invoiceId:                         invoiceId,
		apiPublicUsageDetailsQueryRequest: apiPublicUsageDetailsQueryRequest,
	}
}

// SearchInvoiceLineItemsExecute executes the request
//
//	@return PaginatedPublicApiUsageDetailsLineItem
func (a *InvoicesApiService) SearchInvoiceLineItemsExecute(r SearchInvoiceLineItemsApiRequest) (*PaginatedPublicApiUsageDetailsLineItem, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    any
		formFiles           []formFile
		localVarReturnValue *PaginatedPublicApiUsageDetailsLineItem
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "InvoicesApiService.SearchInvoiceLineItems")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/atlas/v2/orgs/{orgId}/invoices/{invoiceId}/lineItems:search"
	if r.orgId == "" {
		return localVarReturnValue, nil, reportError("orgId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"orgId"+"}", url.PathEscape(r.orgId), -1)
	if r.invoiceId == "" {
		return localVarReturnValue, nil, reportError("invoiceId is empty and must be specified")
	}
	localVarPath = strings.Replace(localVarPath, "{"+"invoiceId"+"}", url.PathEscape(r.invoiceId), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.apiPublicUsageDetailsQueryRequest == nil {
		return localVarReturnValue, nil, reportError("apiPublicUsageDetailsQueryRequest is required and must be specified")
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
	localVarHTTPContentTypes := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header (only first one)
	localVarHTTPHeaderAccepts := []string{"application/vnd.atlas.2024-08-05+json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.apiPublicUsageDetailsQueryRequest
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
