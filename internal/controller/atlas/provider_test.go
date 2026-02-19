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

package atlas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

func TestProvider_IsCloudGov(t *testing.T) {
	t.Run("should return false for invalid domain", func(t *testing.T) {
		p := NewProductionProvider("http://x:namedport", false, false)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return false for commercial Atlas domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodb.com/", false, false)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return true for Atlas for government domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodbgov.com/", false, false)
		assert.True(t, p.IsCloudGov())
	})
}

func TestProvider_IsResourceSupported(t *testing.T) {
	dataProvider := map[string]struct {
		domain      string
		resource    api.AtlasCustomResource
		expectation bool
	}{
		"should return true when it's commercial Atlas": {
			domain:      "https://cloud.mongodb.com",
			resource:    &akov2.AtlasDataFederation{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Project": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasProject{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is Team": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasTeam{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupSchedule": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupSchedule{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is BackupPolicy": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is DatabaseUser": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasBackupPolicy{},
			expectation: true,
		},
		"should return true when it's Atlas Gov and resource is regular Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{},
			},
			expectation: true,
		},
		"should return false when it's Atlas Gov and resource is DataFederation": {
			domain:      "https://cloud.mongodbgov.com",
			resource:    &akov2.AtlasDataFederation{},
			expectation: false,
		},
		"should return false when it's Atlas Gov and resource is Serverless Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					ServerlessSpec: &akov2.ServerlessSpec{},
				},
			},
			expectation: false,
		},
		"should return false when it's Atlas Gov and resource is Flex Deployment": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					FlexSpec: &akov2.FlexSpec{},
				},
			},
			expectation: false,
		},
		"should return false when it's Atlas Gov and resource is a Deployment with search nodes": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{
						SearchNodes: []akov2.SearchNode{
							{
								InstanceSize: "M10",
								NodeCount:    3,
							},
						},
					},
				},
			},
			expectation: false,
		},
		"should return true when it's Atlas Gov and resource is a Deployment with no search nodes": {
			domain: "https://cloud.mongodbgov.com",
			resource: &akov2.AtlasDeployment{
				Spec: akov2.AtlasDeploymentSpec{
					DeploymentSpec: &akov2.AdvancedDeploymentSpec{},
				},
			},
			expectation: true,
		},
	}

	for desc, data := range dataProvider {
		t.Run(desc, func(t *testing.T) {
			p := NewProductionProvider(data.domain, false, false)
			assert.Equal(t, data.expectation, p.IsResourceSupported(data.resource))
		})
	}
}

func TestProvider_SdkClientSet_NilCredentials(t *testing.T) {
	p := NewProductionProvider("https://cloud.mongodb.com", false, false)
	_, err := p.SdkClientSet(context.Background(), &Credentials{}, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no credentials provided")
}

func TestProvider_SdkClientSet_APIKeys(t *testing.T) {
	p := NewProductionProvider("https://cloud.mongodb.com", false, false)
	creds := &Credentials{APIKeys: &APIKeys{PublicKey: "pub", PrivateKey: "priv"}}
	cs, err := p.SdkClientSet(context.Background(), creds, zap.NewNop().Sugar())
	require.NoError(t, err)
	assert.NotNil(t, cs)
	assert.NotNil(t, cs.SdkClient20250312013)
}

func TestProvider_SdkClientSet_ServiceAccount(t *testing.T) {
	p := NewProductionProvider("https://cloud.mongodb.com", false, false)
	creds := &Credentials{ServiceAccount: &ServiceAccountToken{BearerToken: "test-token"}}
	cs, err := p.SdkClientSet(context.Background(), creds, zap.NewNop().Sugar())
	require.NoError(t, err)
	assert.NotNil(t, cs)
	assert.NotNil(t, cs.SdkClient20250312013)
}

func TestBearerTokenTransport(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer my-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	transport := &bearerTokenTransport{token: "my-token"}
	httpClient := &http.Client{Transport: transport}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOperatorUserAgent(t *testing.T) {
	userAgent := operatorUserAgent()

	require.Contains(t, userAgent, "MongoDBAtlasKubernetesOperator")
	require.Contains(t, userAgent, version.Version)
}
