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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/version"
)

func TestProvider_IsCloudGov(t *testing.T) {
	t.Run("should return false for invalid domain", func(t *testing.T) {
		p := NewProductionProvider("http://x:namedport", false)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return false for commercial Atlas domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodb.com/", false)
		assert.False(t, p.IsCloudGov())
	})

	t.Run("should return true for Atlas for government domain", func(t *testing.T) {
		p := NewProductionProvider("https://cloud.mongodbgov.com/", false)
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
			p := NewProductionProvider(data.domain, false)
			assert.Equal(t, data.expectation, p.IsResourceSupported(data.resource))
		})
	}
}

func TestOperatorUserAgent(t *testing.T) {
	userAgent := operatorUserAgent()

	require.Contains(t, userAgent, "MongoDBAtlasKubernetesOperator")
	require.Contains(t, userAgent, version.Version)
}
