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

package atlasnetworkpeering

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/provider"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkcontainer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/networkpeering"
)

func TestApplyPeeringStatus(t *testing.T) {
	for _, tc := range []struct {
		title      string
		peer       networkpeering.NetworkPeer
		container  networkcontainer.NetworkContainer
		wantStatus status.AtlasNetworkPeeringStatus
	}{
		{
			title: "wrong provider fails",
			peer: networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: v1.AtlasNetworkPeeringConfig{
					Provider: "Azure", // Should be "AZURE"
				},
			},
			wantStatus: status.AtlasNetworkPeeringStatus{
				Status: "unsupported provider: \"Azure\"",
			},
		},

		{
			title: "Sample AWS works",
			peer: networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: v1.AtlasNetworkPeeringConfig{
					ID:       "peer-id",
					Provider: string(provider.ProviderAWS),
					AWSConfiguration: &v1.AWSNetworkPeeringConfiguration{
						AccepterRegionName:  "us-east-1",
						AWSAccountID:        "some-aws-id",
						RouteTableCIDRBlock: "10.0.0.0/18",
						VpcID:               "vpc-id-app-fake",
					},
				},
				ContainerID:  "container-id",
				Status:       "some status",
				ErrorMessage: "some error",
				AWSStatus: &status.AWSPeeringStatus{
					ConnectionID: "connection-id",
				},
			},
			container: networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider: string(provider.ProviderAWS),
					AtlasNetworkContainerConfig: v1.AtlasNetworkContainerConfig{
						ID:        "container-id",
						Region:    "us-east-2",
						CIDRBlock: "11.0.0.0/18",
					},
				},
				ID:          "container-id",
				Provisioned: true,
				AWSStatus: &networkcontainer.AWSContainerStatus{
					VpcID: "vpc-id-container-fake",
				},
			},
			wantStatus: status.AtlasNetworkPeeringStatus{
				ID:     "peer-id",
				Status: "some status",
				AWSStatus: &status.AWSPeeringStatus{
					VpcID:        "vpc-id-container-fake",
					ConnectionID: "connection-id",
				},
			},
		},

		{
			title: "Sample Azure works",
			peer: networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: v1.AtlasNetworkPeeringConfig{
					ID:       "peer-id",
					Provider: string(provider.ProviderAzure),
					AzureConfiguration: &v1.AzureNetworkPeeringConfiguration{
						AzureDirectoryID:    "azure-app-dir-id",
						AzureSubscriptionID: "azure-app-subcription-id",
						ResourceGroupName:   "resource-group",
						VNetName:            "some-net",
					},
				},
				ContainerID:  "container-id",
				Status:       "some status",
				ErrorMessage: "some error",
			},
			container: networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider: string(provider.ProviderAzure),
					AtlasNetworkContainerConfig: v1.AtlasNetworkContainerConfig{
						ID:        "container-id",
						Region:    "US_EAST_2",
						CIDRBlock: "11.0.0.0/18",
					},
				},
				ID:          "container-id",
				Provisioned: true,
				AzureStatus: &networkcontainer.AzureContainerStatus{
					AzureSubscriptionID: "azure-atlas-subcription-id",
					VnetName:            "atlas-net-name",
				},
			},
			wantStatus: status.AtlasNetworkPeeringStatus{
				ID:     "peer-id",
				Status: "some status",
				AzureStatus: &status.AzurePeeringStatus{
					AzureSubscriptionID: "azure-atlas-subcription-id",
					VnetName:            "atlas-net-name",
				},
			},
		},

		{
			title: "Sample GCP works",
			peer: networkpeering.NetworkPeer{
				AtlasNetworkPeeringConfig: v1.AtlasNetworkPeeringConfig{
					ID:       "peer-id",
					Provider: string(provider.ProviderGCP),
					GCPConfiguration: &v1.GCPNetworkPeeringConfiguration{
						GCPProjectID: "gcp-app-project",
						NetworkName:  "gcp-app-network",
					},
				},
				ContainerID:  "container-id",
				Status:       "some status",
				ErrorMessage: "some error",
			},
			container: networkcontainer.NetworkContainer{
				NetworkContainerConfig: networkcontainer.NetworkContainerConfig{
					Provider: string(provider.ProviderGCP),
					AtlasNetworkContainerConfig: v1.AtlasNetworkContainerConfig{
						ID:        "container-id",
						CIDRBlock: "11.0.0.0/18",
					},
				},
				ID:          "container-id",
				Provisioned: true,
				GCPStatus: &networkcontainer.GCPContainerStatus{
					GCPProjectID: "gcp-atlas-project",
					NetworkName:  "gcp-atlas-network",
				},
			},
			wantStatus: status.AtlasNetworkPeeringStatus{
				ID:     "peer-id",
				Status: "some status",
				GCPStatus: &status.GCPPeeringStatus{
					GCPProjectID: "gcp-atlas-project",
					NetworkName:  "gcp-atlas-network",
				},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			status := status.AtlasNetworkPeeringStatus{}
			networkpeering.ApplyPeeringStatus(&status, &tc.peer, &tc.container)
			assert.Equal(t, tc.wantStatus, status)
		})
	}
}
