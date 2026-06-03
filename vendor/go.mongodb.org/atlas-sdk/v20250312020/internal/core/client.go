// Copyright 2023 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Response is a MongoDBAtlas response. This wraps the standard http.Response returned from MongoDBAtlas API.
type Response struct {
	*http.Response

	// Links that were returned with the response.
	Links []*Link `json:"links"`

	// Raw data from server response
	Raw []byte `json:"-"`
}

// Link is the link to sub-resources and/or related resources.
type Link struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

// ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// Response that caused this error
	Response *http.Response
	// ErrorCode is the code as specified in https://docs.atlas.mongodb.com/reference/api/api-errors/
	ErrorCode string `json:"errorCode"`
	// HTTPCode status code.
	HTTPCode int `json:"error"`
	// Reason is short description of the error, which is simply the HTTP status phrase.
	Reason string `json:"reason"`
	// Detail is more detailed description of the error.
	Detail string `json:"detail,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d (request %q) %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.ErrorCode, r.Detail)
}

func (r *ErrorResponse) Is(target error) bool {
	var v *ErrorResponse

	return errors.As(target, &v) &&
		r.ErrorCode == v.ErrorCode &&
		r.HTTPCode == v.HTTPCode &&
		r.Reason == v.Reason &&
		r.Detail == v.Detail
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func (resp *Response) CheckResponse(body io.ReadCloser) error {
	if c := resp.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: resp.Response}
	data, err := io.ReadAll(body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			log.Printf("[DEBUG] unmarshal error response: %s", err)
			errorResponse.Reason = string(data)
		}
	}

	return errorResponse
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req) //nolint:gosec // G704 - HTTP SDK makes requests to user-configured servers
}
