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

package datafederation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/atlas-sdk/v20250312021/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
)

func TestCloudProviderConfigRoundtrip(t *testing.T) {
	tests := map[string]struct {
		input    *akov2.CloudProviderConfig
		wantSpec *akov2.CloudProviderConfig
	}{
		"nil config": {},
		"AWS only": {
			input:    &akov2.CloudProviderConfig{AWS: &akov2.AWSProviderConfig{RoleID: "role1", TestS3Bucket: "bucket1"}},
			wantSpec: &akov2.CloudProviderConfig{AWS: &akov2.AWSProviderConfig{RoleID: "role1", TestS3Bucket: "bucket1"}},
		},
		"Azure only": {
			input:    &akov2.CloudProviderConfig{Azure: &akov2.AzureProviderConfig{RoleID: "az-role"}},
			wantSpec: &akov2.CloudProviderConfig{Azure: &akov2.AzureProviderConfig{RoleID: "az-role"}},
		},
		"GCP only": {
			input:    &akov2.CloudProviderConfig{GCP: &akov2.GCPProviderConfig{RoleID: "gcp-role"}},
			wantSpec: &akov2.CloudProviderConfig{GCP: &akov2.GCPProviderConfig{RoleID: "gcp-role"}},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			atlasOut := cloudProviderConfigToAtlas(tt.input)
			got := cloudProviderConfigFromAtlas(atlasOut)
			assert.Equal(t, tt.wantSpec, got)
		})
	}
}

func TestCloudProviderConfigStatusFromAtlas(t *testing.T) {
	tests := map[string]struct {
		atlasIn *admin.DataLakeCloudProviderConfig
		want    *status.DataFederationCloudProviderConfigStatus
	}{
		"nil config": {},
		"Azure - read-only fields mapped to status": {
			atlasIn: &admin.DataLakeCloudProviderConfig{
				Azure: &admin.DataFederationAzureCloudProviderConfig{
					RoleId:             "az-role",
					AtlasAppId:         new("app-id"),
					ServicePrincipalId: new("sp-id"),
					TenantId:           new("tenant-id"),
				},
			},
			want: &status.DataFederationCloudProviderConfigStatus{
				Azure: &status.AzureProviderConfigStatus{
					AtlasAppID:         "app-id",
					ServicePrincipalID: "sp-id",
					TenantID:           "tenant-id",
				},
			},
		},
		"GCP - read-only service account mapped to status": {
			atlasIn: &admin.DataLakeCloudProviderConfig{
				Gcp: &admin.DataFederationGCPCloudProviderConfig{
					RoleId:            "gcp-role",
					GcpServiceAccount: new("sa@project.iam.gserviceaccount.com"),
				},
			},
			want: &status.DataFederationCloudProviderConfigStatus{
				GCP: &status.GCPProviderConfigStatus{
					GCPServiceAccount: "sa@project.iam.gserviceaccount.com",
				},
			},
		},
		"AWS - IAM fields mapped to status": {
			atlasIn: &admin.DataLakeCloudProviderConfig{
				Aws: &admin.DataLakeAWSCloudProviderConfig{
					RoleId:            "role1",
					TestS3Bucket:      "bucket1",
					ExternalId:        new("ext-id"),
					IamAssumedRoleARN: new("arn:aws:iam::123:role/role"),
					IamUserARN:        new("arn:aws:iam::123:user/user"),
				},
			},
			want: &status.DataFederationCloudProviderConfigStatus{
				AWS: &status.AWSProviderConfigStatus{
					ExternalID:        "ext-id",
					IAMAssumedRoleARN: "arn:aws:iam::123:role/role",
					IAMUserARN:        "arn:aws:iam::123:user/user",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := cloudProviderConfigStatusFromAtlas(tt.atlasIn)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAtlasClusterStoreRoundtrip(t *testing.T) {
	input := &akov2.Storage{
		Stores: []akov2.Store{
			{
				Name:        "atlas-store",
				Provider:    "atlas",
				ClusterName: "my-cluster",
				ReadConcern: &akov2.ReadConcern{Level: "local"},
				ReadPreference: &akov2.ReadPreference{
					Mode:                "secondary",
					MaxStalenessSeconds: 120,
					TagSets: [][]akov2.ReadPreferenceTag{
						{{Name: "region", Value: "us-east-1"}},
						{{Name: "dc", Value: "nyc"}},
					},
				},
			},
		},
	}

	atlasOut := storageToAtlas(input)
	require.NotNil(t, atlasOut)
	require.Len(t, atlasOut.GetStores(), 1)

	st := atlasOut.GetStores()[0]
	assert.Equal(t, "atlas-store", st.GetName())
	assert.Equal(t, "atlas", st.Provider)
	assert.Equal(t, "my-cluster", st.GetClusterName())
	assert.Equal(t, "local", st.ReadConcern.GetLevel())
	assert.Equal(t, "secondary", st.ReadPreference.GetMode())
	assert.Equal(t, 120, st.ReadPreference.GetMaxStalenessSeconds())
	assert.Len(t, st.ReadPreference.GetTagSets(), 2)

	got := storageFromAtlas(atlasOut)
	require.NotNil(t, got)
	require.Len(t, got.Stores, 1)
	assert.Equal(t, input.Stores[0].ClusterName, got.Stores[0].ClusterName)
	assert.Equal(t, input.Stores[0].ReadConcern, got.Stores[0].ReadConcern)
	assert.Equal(t, input.Stores[0].ReadPreference, got.Stores[0].ReadPreference)
}

func TestRoundtrip_DataFederation(t *testing.T) {
	f := fuzz.New()

	for range 100 {
		fuzzed := &DataFederation{}
		f.Fuzz(fuzzed)
		fuzzed, err := NewDataFederation(fuzzed.DataFederationSpec, fuzzed.ProjectID, fuzzed.Hostnames)
		require.NoError(t, err)

		toAtlasResult := toAtlas(fuzzed)
		fromAtlasResult, err := fromAtlas(toAtlasResult)
		require.NoError(t, err)

		equals := fuzzed.SpecEqualsTo(fromAtlasResult)
		if !equals {
			t.Log(cmp.Diff(fuzzed, fromAtlasResult))
		}
		require.True(t, equals)
	}
}
