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

package httputil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

type MockTransport struct {
	mock.Mock
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func newResponse(body string) *http.Response {
	return &http.Response{
		Body: io.NopCloser(strings.NewReader(body)),
	}
}

func newRequest(method, url, body string) *http.Request {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, _ := http.NewRequestWithContext(context.Background(), method, url, bodyReader)
	return req
}

func TestNewTransportWithDiff(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()

	transport := NewTransportWithDiff(mockTransport, logger)

	assert.NotNil(t, transport)
	assert.Equal(t, mockTransport, transport.transport)
	assert.Equal(t, logger, transport.log)
}

func TestTransportWithDiff_RoundTrip_GET_Request(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	req := newRequest(http.MethodGet, "https://cloud-qa.mongodb.com/v1/groups", "")
	expectedResp := newResponse("test response")

	mockTransport.On("RoundTrip", req).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_RoundTrip_PUT_Request_Success(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	originalData := map[string]any{
		"id":      1,
		"name":    "old name",
		"links":   []string{"link1", "link2"},
		"created": "2023-01-01",
	}
	originalJSON, _ := json.Marshal(originalData)

	updatedData := map[string]any{
		"id":   1,
		"name": "new name",
	}
	updatedJSON, _ := json.Marshal(updatedData)

	putReq := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", string(updatedJSON))

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(string(originalJSON)), nil)

	expectedResp := newResponse("updated")
	mockTransport.On("RoundTrip", putReq).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(putReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_RoundTrip_PATCH_Request_Success(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	originalData := map[string]any{
		"id":   1,
		"name": "old name",
	}
	originalJSON, _ := json.Marshal(originalData)

	updatedData := map[string]any{
		"id":   1,
		"name": "new name",
	}
	updatedJSON, _ := json.Marshal(updatedData)

	patchReq := newRequest(http.MethodPatch, "https://cloud-qa.mongodb.com/api/v1/groups", string(updatedJSON))

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(string(originalJSON)), nil)

	expectedResp := newResponse("updated")
	mockTransport.On("RoundTrip", patchReq).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(patchReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_RoundTrip_PUT_Request_GET_Error(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	updatedData := map[string]any{
		"id":   1,
		"name": "new name",
	}
	updatedJSON, _ := json.Marshal(updatedData)

	putReq := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", string(updatedJSON))

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return((*http.Response)(nil), fmt.Errorf("network error"))

	expectedResp := newResponse("updated")
	mockTransport.On("RoundTrip", putReq).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(putReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_RoundTrip_PUT_Request_Invalid_JSON(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	putReq := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", `{"invalid": json}`)

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(`{"invalid": json}`), nil)

	expectedResp := newResponse("updated")
	mockTransport.On("RoundTrip", putReq).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(putReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_RoundTrip_PUT_Request_No_Body(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	putReq := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", "")

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(`{"id": 1}`), nil)

	expectedResp := newResponse("updated")
	mockTransport.On("RoundTrip", putReq).Return(expectedResp, nil)

	resp, err := transport.RoundTrip(putReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedResp, resp)
	mockTransport.AssertExpectations(t)
}

func TestCleanLinksField(t *testing.T) {
	data := map[string]any{
		"id":    1,
		"name":  "test",
		"links": []string{"link1", "link2"},
	}

	cleanLinksField(data)

	_, exists := data["links"]
	assert.False(t, exists)
	assert.Equal(t, 1, data["id"])
	assert.Equal(t, "test", data["name"])
}

func TestCleanLinksField_NoLinksField(t *testing.T) {
	data := map[string]any{
		"id":   1,
		"name": "test",
	}

	cleanLinksField(data)

	assert.Equal(t, 1, data["id"])
	assert.Equal(t, "test", data["name"])
}

func TestCleanCreatedField(t *testing.T) {
	data := map[string]any{
		"id":      1,
		"name":    "test",
		"created": "2023-01-01",
	}

	cleanCreatedField(data)

	_, exists := data["created"]
	assert.False(t, exists)
	assert.Equal(t, 1, data["id"])
	assert.Equal(t, "test", data["name"])
}

func TestCleanCreatedField_NoCreatedField(t *testing.T) {
	data := map[string]any{
		"id":   1,
		"name": "test",
	}

	cleanCreatedField(data)

	assert.Equal(t, 1, data["id"])
	assert.Equal(t, "test", data["name"])
}

func TestTransportWithDiff_tryCalculateDiff_Success(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	originalData := map[string]any{
		"id":      1,
		"name":    "old name",
		"links":   []string{"link1", "link2"},
		"created": "2023-01-01",
	}
	originalJSON, _ := json.Marshal(originalData)

	updatedData := map[string]any{
		"id":   1,
		"name": "new name",
	}
	updatedJSON, _ := json.Marshal(updatedData)

	req := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", string(updatedJSON))

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(string(originalJSON)), nil)

	diffString, err := transport.tryCalculateDiff(req, cleanLinksField, cleanCreatedField)

	assert.NoError(t, err)
	assert.NotEmpty(t, diffString)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_tryCalculateDiff_GET_Error(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	req := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", `{"id": 1}`)

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return((*http.Response)(nil), fmt.Errorf("network error"))

	diffString, err := transport.tryCalculateDiff(req, cleanLinksField, cleanCreatedField)

	assert.Error(t, err)
	assert.Empty(t, diffString)
	assert.Contains(t, err.Error(), "failed to GET original resource")
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_tryCalculateDiff_Invalid_JSON_Response(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	req := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", `{"id": 1}`)

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(`{"invalid": json}`), nil)

	diffString, err := transport.tryCalculateDiff(req, cleanLinksField, cleanCreatedField)

	assert.Error(t, err)
	assert.Empty(t, diffString)
	assert.Contains(t, err.Error(), "failed to unmarshal payloadFromGetParsed JSON")
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_tryCalculateDiff_No_Body(t *testing.T) {
	mockTransport := &MockTransport{}
	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(mockTransport, logger)

	req := newRequest(http.MethodPut, "https://cloud-qa.mongodb.com/api/v1/groups", "")

	mockTransport.On("RoundTrip", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == http.MethodGet
	})).Return(newResponse(`{"id": 1}`), nil)

	diffString, err := transport.tryCalculateDiff(req, cleanLinksField, cleanCreatedField)

	assert.Contains(t, err.Error(), "failed to compare JSON payloads")
	assert.Empty(t, diffString)
	mockTransport.AssertExpectations(t)
}

func TestTransportWithDiff_Integration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			response := map[string]any{
				"id":      1,
				"name":    "original name",
				"links":   []string{"link1", "link2"},
				"created": "2023-01-01",
			}
			json.NewEncoder(w).Encode(response)
		} else if r.Method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("updated"))
		}
	}))
	defer server.Close()

	logger := zaptest.NewLogger(t).Sugar()
	transport := NewTransportWithDiff(http.DefaultTransport, logger)

	updatedData := map[string]any{
		"id":   1,
		"name": "new name",
	}
	updatedJSON, _ := json.Marshal(updatedData)

	req := newRequest(http.MethodPut, server.URL+"/api/1", string(updatedJSON))

	resp, err := transport.RoundTrip(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, StatusCode(resp))

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "updated", string(body))
}
