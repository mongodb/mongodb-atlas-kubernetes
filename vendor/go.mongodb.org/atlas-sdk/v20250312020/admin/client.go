// Code based on the AtlasAPI V2 OpenAPI file

package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	jsonCheck       = regexp.MustCompile(`(?i:(?:application|text)/(?:vnd\.[^;]+\+)?json)`)
	xmlCheck        = regexp.MustCompile(`(?i:(?:application|text)/xml)`)
	queryParamSplit = regexp.MustCompile(`(^|&)([^&]+)`)
	queryDescape    = strings.NewReplacer("%5B", "[", "%5D", "]")
)

// APIClient manages communication with the MongoDB Atlas Administration API API
// In most cases there should be only one, shared, APIClient.
type APIClient struct {
	cfg           *Configuration
	common        service       // Reuse a single struct instead of allocating one for each service on the heap.
	UntypedClient UntypedClient // Make API calls without using a typed model.

	// API Services

	AWSClustersDNSApi AWSClustersDNSApi

	AccessTrackingApi AccessTrackingApi

	ActivityFeedApi ActivityFeedApi

	AlertConfigurationsApi AlertConfigurationsApi

	AlertsApi AlertsApi

	AtlasSearchApi AtlasSearchApi

	AuditingApi AuditingApi

	CloudBackupsApi CloudBackupsApi

	CloudMigrationServiceApi CloudMigrationServiceApi

	CloudProviderAccessApi CloudProviderAccessApi

	ClusterOutageSimulationApi ClusterOutageSimulationApi

	ClustersApi ClustersApi

	CollectionLevelMetricsApi CollectionLevelMetricsApi

	CustomDatabaseRolesApi CustomDatabaseRolesApi

	DataFederationApi DataFederationApi

	DataLakePipelinesApi DataLakePipelinesApi

	DatabaseUsersApi DatabaseUsersApi

	EncryptionAtRestUsingCustomerKeyManagementApi EncryptionAtRestUsingCustomerKeyManagementApi

	EventsApi EventsApi

	FederatedAuthenticationApi FederatedAuthenticationApi

	FlexClustersApi FlexClustersApi

	FlexRestoreJobsApi FlexRestoreJobsApi

	FlexSnapshotsApi FlexSnapshotsApi

	GlobalClustersApi GlobalClustersApi

	InvoicesApi InvoicesApi

	LDAPConfigurationApi LDAPConfigurationApi

	LegacyBackupApi LegacyBackupApi

	MaintenanceWindowsApi MaintenanceWindowsApi

	MongoDBCloudUsersApi MongoDBCloudUsersApi

	MonitoringAndLogsApi MonitoringAndLogsApi

	NetworkPeeringApi NetworkPeeringApi

	OnlineArchiveApi OnlineArchiveApi

	OrganizationsApi OrganizationsApi

	PerformanceAdvisorApi PerformanceAdvisorApi

	PrivateEndpointServicesApi PrivateEndpointServicesApi

	ProgrammaticAPIKeysApi ProgrammaticAPIKeysApi

	ProjectIPAccessListApi ProjectIPAccessListApi

	ProjectsApi ProjectsApi

	PushBasedLogExportApi PushBasedLogExportApi

	QueryShapeInsightsApi QueryShapeInsightsApi

	RateLimitingApi RateLimitingApi

	ResourcePoliciesApi ResourcePoliciesApi

	RollingIndexApi RollingIndexApi

	RootApi RootApi

	ServerlessInstancesApi ServerlessInstancesApi

	ServerlessPrivateEndpointsApi ServerlessPrivateEndpointsApi

	ServiceAccountsApi ServiceAccountsApi

	StreamsApi StreamsApi

	TeamsApi TeamsApi

	ThirdPartyIntegrationsApi ThirdPartyIntegrationsApi

	X509AuthenticationApi X509AuthenticationApi
}

type service struct {
	client *APIClient
}

// NewAPIClient creates a new API client. Requires a userAgent string describing your application.
// optionally a custom http.Client to allow for advanced features such as caching.
func NewAPIClient(cfg *Configuration) *APIClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}

	c := &APIClient{}
	c.cfg = cfg
	c.common.client = c
	c.UntypedClient.client = c

	// API Services
	c.AWSClustersDNSApi = (*AWSClustersDNSApiService)(&c.common)
	c.AccessTrackingApi = (*AccessTrackingApiService)(&c.common)
	c.ActivityFeedApi = (*ActivityFeedApiService)(&c.common)
	c.AlertConfigurationsApi = (*AlertConfigurationsApiService)(&c.common)
	c.AlertsApi = (*AlertsApiService)(&c.common)
	c.AtlasSearchApi = (*AtlasSearchApiService)(&c.common)
	c.AuditingApi = (*AuditingApiService)(&c.common)
	c.CloudBackupsApi = (*CloudBackupsApiService)(&c.common)
	c.CloudMigrationServiceApi = (*CloudMigrationServiceApiService)(&c.common)
	c.CloudProviderAccessApi = (*CloudProviderAccessApiService)(&c.common)
	c.ClusterOutageSimulationApi = (*ClusterOutageSimulationApiService)(&c.common)
	c.ClustersApi = (*ClustersApiService)(&c.common)
	c.CollectionLevelMetricsApi = (*CollectionLevelMetricsApiService)(&c.common)
	c.CustomDatabaseRolesApi = (*CustomDatabaseRolesApiService)(&c.common)
	c.DataFederationApi = (*DataFederationApiService)(&c.common)
	c.DataLakePipelinesApi = (*DataLakePipelinesApiService)(&c.common)
	c.DatabaseUsersApi = (*DatabaseUsersApiService)(&c.common)
	c.EncryptionAtRestUsingCustomerKeyManagementApi = (*EncryptionAtRestUsingCustomerKeyManagementApiService)(&c.common)
	c.EventsApi = (*EventsApiService)(&c.common)
	c.FederatedAuthenticationApi = (*FederatedAuthenticationApiService)(&c.common)
	c.FlexClustersApi = (*FlexClustersApiService)(&c.common)
	c.FlexRestoreJobsApi = (*FlexRestoreJobsApiService)(&c.common)
	c.FlexSnapshotsApi = (*FlexSnapshotsApiService)(&c.common)
	c.GlobalClustersApi = (*GlobalClustersApiService)(&c.common)
	c.InvoicesApi = (*InvoicesApiService)(&c.common)
	c.LDAPConfigurationApi = (*LDAPConfigurationApiService)(&c.common)
	c.LegacyBackupApi = (*LegacyBackupApiService)(&c.common)
	c.MaintenanceWindowsApi = (*MaintenanceWindowsApiService)(&c.common)
	c.MongoDBCloudUsersApi = (*MongoDBCloudUsersApiService)(&c.common)
	c.MonitoringAndLogsApi = (*MonitoringAndLogsApiService)(&c.common)
	c.NetworkPeeringApi = (*NetworkPeeringApiService)(&c.common)
	c.OnlineArchiveApi = (*OnlineArchiveApiService)(&c.common)
	c.OrganizationsApi = (*OrganizationsApiService)(&c.common)
	c.PerformanceAdvisorApi = (*PerformanceAdvisorApiService)(&c.common)
	c.PrivateEndpointServicesApi = (*PrivateEndpointServicesApiService)(&c.common)
	c.ProgrammaticAPIKeysApi = (*ProgrammaticAPIKeysApiService)(&c.common)
	c.ProjectIPAccessListApi = (*ProjectIPAccessListApiService)(&c.common)
	c.ProjectsApi = (*ProjectsApiService)(&c.common)
	c.PushBasedLogExportApi = (*PushBasedLogExportApiService)(&c.common)
	c.QueryShapeInsightsApi = (*QueryShapeInsightsApiService)(&c.common)
	c.RateLimitingApi = (*RateLimitingApiService)(&c.common)
	c.ResourcePoliciesApi = (*ResourcePoliciesApiService)(&c.common)
	c.RollingIndexApi = (*RollingIndexApiService)(&c.common)
	c.RootApi = (*RootApiService)(&c.common)
	c.ServerlessInstancesApi = (*ServerlessInstancesApiService)(&c.common)
	c.ServerlessPrivateEndpointsApi = (*ServerlessPrivateEndpointsApiService)(&c.common)
	c.ServiceAccountsApi = (*ServiceAccountsApiService)(&c.common)
	c.StreamsApi = (*StreamsApiService)(&c.common)
	c.TeamsApi = (*TeamsApiService)(&c.common)
	c.ThirdPartyIntegrationsApi = (*ThirdPartyIntegrationsApiService)(&c.common)
	c.X509AuthenticationApi = (*X509AuthenticationApiService)(&c.common)

	return c
}

// selectHeaderContentType select a content type from the available list.
func selectHeaderContentType(contentTypes []string) string {
	if len(contentTypes) == 0 {
		return ""
	}
	if contains(contentTypes, "application/json") {
		return "application/json"
	}
	return contentTypes[0] // use the first content type specified in 'consumes'
}

// selectHeaderAccept join all accept types and return
func selectHeaderAccept(accepts []string) string {
	if len(accepts) == 0 {
		return ""
	}
	return accepts[0]
}

// contains is a case insensitive match, finding needle in a haystack
func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if strings.EqualFold(a, needle) {
			return true
		}
	}
	return false
}

// parameterAddToHeaderOrQuery adds the provided object to the request header or url query
// supporting deep object syntax
func parameterAddToHeaderOrQuery(headerOrQueryParams any, keyPrefix string, obj any, collectionType string) {
	var v = reflect.ValueOf(obj)
	var value = ""
	if v == reflect.ValueOf(nil) {
		value = "null"
	} else {
		switch v.Kind() {
		case reflect.Invalid:
			value = "invalid"

		case reflect.Struct:
			if t, ok := obj.(time.Time); ok {
				parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, t.Format(time.RFC3339), collectionType)
				return
			}
			value = v.Type().String() + " value"
		case reflect.Slice:
			var indValue = reflect.ValueOf(obj)
			if indValue == reflect.ValueOf(nil) {
				return
			}
			var lenIndValue = indValue.Len()
			for i := 0; i < lenIndValue; i++ {
				var arrayValue = indValue.Index(i)
				parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, arrayValue.Interface(), collectionType)
			}
			return

		case reflect.Map:
			var indValue = reflect.ValueOf(obj)
			if indValue == reflect.ValueOf(nil) {
				return
			}
			iter := indValue.MapRange()
			for iter.Next() {
				mapKey, mapValue := iter.Key(), iter.Value()
				parameterAddToHeaderOrQuery(headerOrQueryParams, fmt.Sprintf("%s[%s]", keyPrefix, mapKey.String()), mapValue.Interface(), collectionType)
			}
			return

		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			parameterAddToHeaderOrQuery(headerOrQueryParams, keyPrefix, v.Elem().Interface(), collectionType)
			return

		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(v.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = strconv.FormatUint(v.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(v.Float(), 'g', -1, 32)
		case reflect.Bool:
			value = strconv.FormatBool(v.Bool())
		case reflect.String:
			value = v.String()
		default:
			value = v.Type().String() + " value"
		}
	}

	switch valuesMap := headerOrQueryParams.(type) {
	case url.Values:
		if collectionType == "csv" && valuesMap.Get(keyPrefix) != "" {
			valuesMap.Set(keyPrefix, valuesMap.Get(keyPrefix)+","+value)
		} else {
			valuesMap.Add(keyPrefix, value)
		}
	case map[string]string:
		valuesMap[keyPrefix] = value
	}
}

// callAPI do the request.
func (c *APIClient) callAPI(request *http.Request) (*http.Response, error) {
	if c.cfg.Debug {
		dump, err := httputil.DumpRequestOut(request, true)
		if err != nil {
			return nil, err
		}
		log.Printf("\n%s\n", string(dump)) //nolint:gosec // G706 - debug-only logging of raw HTTP traffic
	}

	resp, err := c.cfg.HTTPClient.Do(request) //nolint:gosec // G704 - HTTP SDK makes requests to user-configured servers
	if err != nil {
		return resp, err
	}

	if c.cfg.Debug {
		dump, err1 := httputil.DumpResponse(resp, true)
		if err1 != nil {
			return resp, err
		}
		log.Printf("\n%s\n", string(dump)) //nolint:gosec // G706 - debug-only logging of raw HTTP traffic
	}
	return resp, err
}

// Allow modification of underlying config for alternate implementations and testing
// Caution: modifying the configuration while live can cause data races and potentially unwanted behavior
func (c *APIClient) GetConfig() *Configuration {
	return c.cfg
}

type formFile struct {
	fileBytes    []byte
	fileName     string
	formFileName string
}

// prepareRequest build the request
func (c *APIClient) prepareRequest(
	ctx context.Context,
	path string,
	method string,
	postBody any,
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	formFiles []formFile) (localVarRequest *http.Request, err error) {
	var body *bytes.Buffer

	// Detect postBody type and post.
	if postBody != nil {
		contentType := headerParams["Content-Type"]
		if contentType == "" {
			contentType = detectContentType(postBody)
			headerParams["Content-Type"] = contentType
		}

		body, err = setBody(postBody, contentType)
		if err != nil {
			return nil, err
		}
	}

	// add form parameters and file if available.
	if strings.HasPrefix(headerParams["Content-Type"], "multipart/form-data") && len(formParams) > 0 || (len(formFiles) > 0) {
		if body != nil {
			return nil, errors.New("cannot specify postBody and multipart form at the same time")
		}
		body = &bytes.Buffer{}
		w := multipart.NewWriter(body)

		for k, v := range formParams {
			for _, iv := range v {
				if strings.HasPrefix(k, "@") { // file
					err = addFile(w, k[1:], iv)
					if err != nil {
						return nil, err
					}
				} else { // form value
					err = w.WriteField(k, iv)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		for _, formFile := range formFiles {
			if len(formFile.fileBytes) > 0 && formFile.fileName != "" {
				w.Boundary()
				part, err1 := w.CreateFormFile(formFile.formFileName, filepath.Base(formFile.fileName))
				if err1 != nil {
					return nil, err
				}
				_, err1 = part.Write(formFile.fileBytes)
				if err1 != nil {
					return nil, err
				}
			}
		}

		// Set the Boundary in the Content-Type
		headerParams["Content-Type"] = w.FormDataContentType()

		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
		w.Close()
	}

	if strings.HasPrefix(headerParams["Content-Type"], "application/x-www-form-urlencoded") && len(formParams) > 0 {
		if body != nil {
			return nil, errors.New("cannot specify postBody and x-www-form-urlencoded form at the same time")
		}
		body = &bytes.Buffer{}
		body.WriteString(formParams.Encode())
		// Set Content-Length
		headerParams["Content-Length"] = fmt.Sprintf("%d", body.Len())
	}

	// Setup path and query parameters
	urlData, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Override request host, if applicable
	if c.cfg.Host != "" {
		urlData.Host = c.cfg.Host
	}

	// Override request scheme, if applicable
	if c.cfg.Scheme != "" {
		urlData.Scheme = c.cfg.Scheme
	}

	// Adding Query Param
	query := urlData.Query()
	for k, v := range queryParams {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}

	// Encode the parameters.
	urlData.RawQuery = queryParamSplit.ReplaceAllStringFunc(query.Encode(), func(s string) string {
		pieces := strings.Split(s, "=")
		pieces[0] = queryDescape.Replace(pieces[0])
		return strings.Join(pieces, "=")
	})

	// Generate a new request
	if body != nil {
		localVarRequest, err = http.NewRequest(method, urlData.String(), body)
	} else {
		localVarRequest, err = http.NewRequest(method, urlData.String(), http.NoBody)
	}
	if err != nil {
		return nil, err
	}

	// add header parameters, if any
	if len(headerParams) > 0 {
		headers := http.Header{}
		for h, v := range headerParams {
			headers[h] = []string{v}
		}
		localVarRequest.Header = headers
	}

	// Add the user agent to the request.
	localVarRequest.Header.Add("User-Agent", c.cfg.UserAgent)

	if ctx != nil {
		// add context to the request
		localVarRequest = localVarRequest.WithContext(ctx)
	}

	for header, value := range c.cfg.DefaultHeader {
		localVarRequest.Header.Add(header, value)
	}
	return localVarRequest, nil
}

func (c *APIClient) makeApiError(res *http.Response, httpMethod, httpPath string) error {
	defer res.Body.Close()

	newErr := &GenericOpenAPIError{
		error: res.Status,
	}

	localVarBody, err := io.ReadAll(res.Body)
	if err != nil {
		newErr.error = fmt.Sprintf("(%s) failed to read response body: %s", res.Status, err.Error())
		return newErr
	}
	newErr.body = localVarBody

	var v ApiError
	err = c.decode(&v, io.NopCloser(bytes.NewBuffer(localVarBody)), res.Header.Get("Content-Type"))
	if err != nil {
		newErr.error = fmt.Sprintf("(%s) failed to decode response body: %s", res.Status, err.Error())
		return newErr
	}
	newErr.error = FormatErrorMessageWithDetails(res.Status, httpMethod, httpPath, v)
	newErr.model = v
	return newErr
}

func (c *APIClient) decode(v any, b io.ReadCloser, contentType string) (err error) {
	switch r := v.(type) {
	case *string:
		buf, err := io.ReadAll(b)
		_ = b.Close()
		if err != nil {
			return err
		}
		*r = string(buf)
		return nil
	case *io.ReadCloser:
		*r = b
		return nil
	case **io.ReadCloser:
		*r = &b
		return nil
	default:
		buf, err := io.ReadAll(b)
		_ = b.Close()
		if err != nil {
			return err
		}
		if len(buf) == 0 {
			return nil
		}
		if xmlCheck.MatchString(contentType) {
			return xml.Unmarshal(buf, v)
		}
		if jsonCheck.MatchString(contentType) {
			if actualObj, ok := v.(interface{ GetActualInstance() any }); ok { // oneOf, anyOf schemas
				if unmarshalObj, ok := actualObj.(interface{ UnmarshalJSON([]byte) error }); ok { // make sure it has UnmarshalJSON defined
					if err = unmarshalObj.UnmarshalJSON(buf); err != nil {
						return err
					}
				} else {
					return errors.New("unknown type with GetActualInstance but no unmarshalObj.UnmarshalJSON defined")
				}
			} else {
				// UseNumber preserves large integers in any/[]any/map[string]any fields
				// as json.Number instead of float64, preventing silent precision loss above 2^53.
				dec := json.NewDecoder(bytes.NewReader(buf))
				dec.UseNumber()
				if err = dec.Decode(v); err != nil {
					return err
				}
			}
			return nil
		}
		return errors.New("undefined response type")
	}
}

// Add a file to the multipart request
func addFile(w *multipart.Writer, fieldName, path string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	part, err := w.CreateFormFile(fieldName, filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	return err
}

// Prevent trying to import "fmt"
func reportError(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// Set request body from an any
func setBody(body any, contentType string) (bodyBuf *bytes.Buffer, err error) {
	bodyBuf = &bytes.Buffer{}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if fp, ok := body.(*io.ReadCloser); ok {
		_, err = bodyBuf.ReadFrom(*fp)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else if jsonCheck.MatchString(contentType) {
		err = json.NewEncoder(bodyBuf).Encode(body)
	} else if xmlCheck.MatchString(contentType) {
		err = xml.NewEncoder(bodyBuf).Encode(body)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("invalid body type %s", contentType)
		return nil, err
	}
	return bodyBuf, nil
}

// detectContentType method is used to figure out `Request.Body` content type for request header
func detectContentType(body any) string {
	contentType := "text/plain; charset=utf-8"
	kind := reflect.TypeOf(body).Kind()

	switch kind {
	case reflect.Struct, reflect.Map, reflect.Ptr:
		contentType = "application/json; charset=utf-8"
	case reflect.String:
		contentType = "text/plain; charset=utf-8"
	default:
		if b, ok := body.([]byte); ok {
			contentType = http.DetectContentType(b)
		} else if kind == reflect.Slice {
			contentType = "application/json; charset=utf-8"
		}
	}

	return contentType
}

// GenericOpenAPIError Provides access to the body, error and model on returned errors.
type GenericOpenAPIError struct {
	body  []byte
	error string
	model ApiError
}

// Error returns non-empty string if there was an error.
func (e GenericOpenAPIError) Error() string {
	return e.error
}

// Body returns the raw bytes of the response
func (e GenericOpenAPIError) Body() []byte {
	return e.body
}

// Model returns the unpacked model of the error
func (e GenericOpenAPIError) Model() ApiError {
	return e.model
}

// SetModel sets model instance: Should be only used for testing
func (e *GenericOpenAPIError) SetModel(errorModel ApiError) {
	e.model = errorModel
}

// SetError sets error string: Should be only used for testing
func (e *GenericOpenAPIError) SetError(errorString string) {
	e.error = errorString
}

// FormatErrorMessageWithDetails formats error message using error struct fields. It should be only used for testing.
func FormatErrorMessageWithDetails(status, path, method string, v ApiError) string {
	badRequestDetailString := ""
	if v.BadRequestDetail != nil {
		badRequestDetail, _ := json.Marshal(v.GetBadRequestDetail())
		badRequestDetailString = string(badRequestDetail)
	}

	return fmt.Sprintf("%v %v: HTTP %v (Error code: %q) Detail: %v Reason: %v. Params: %v, BadRequestDetail: %v",
		method, path, status, v.GetErrorCode(),
		v.GetDetail(), v.GetReason(), v.GetParameters(), badRequestDetailString)
}

type UntypedClient struct {
	client *APIClient
}

func (u *UntypedClient) PrepareRequest(
	ctx context.Context,
	path string,
	method string,
	postBody any,
	headerParams map[string]string,
	queryParams url.Values,
	formParams url.Values,
	formFiles []formFile) (localVarRequest *http.Request, err error) {
	return u.client.prepareRequest(ctx, path, method, postBody, headerParams, queryParams, formParams, formFiles)
}

func (u *UntypedClient) CallAPI(request *http.Request) (*http.Response, error) {
	return u.client.callAPI(request)
}

func (u *UntypedClient) MakeApiError(res *http.Response, httpMethod, httpPath string) error {
	return u.client.makeApiError(res, httpMethod, httpPath)
}
