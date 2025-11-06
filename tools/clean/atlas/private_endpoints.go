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
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
	"go.mongodb.org/atlas-sdk/v20250312006/admin"
)

func (c *Cleaner) listPrivateEndpoints(ctx context.Context, projectID, cloudProvider string) []admin.EndpointService {
	endpoints, _, err := c.client.PrivateEndpointServicesApi.
		ListPrivateEndpointServices(ctx, projectID, cloudProvider).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list %s private endpoints for project %s: %s", cloudProvider, projectID, err))

		return nil
	}

	return endpoints
}

func (c *Cleaner) deletePrivateEndpoints(ctx context.Context, projectID, provider string, endpointServices []admin.EndpointService) {
	var endpointIDs []string

	for _, endpointService := range endpointServices {
		switch provider {
		case CloudProviderAWS:
			if endpointService.GetEndpointServiceName() != "" {
				err := c.aws.DeleteEndpoint(ctx, endpointService.GetEndpointServiceName(), endpointService.GetRegionName())
				if err != nil {
					fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC Endpoint %s at region %s from AWS: %s", endpointService.GetEndpointServiceName(), endpointService.GetRegionName(), err))

					continue
				}
			}

			endpointIDs = endpointService.GetInterfaceEndpoints()
		case CloudProviderGCP:
			if len(endpointService.GetEndpointGroupNames()) > 0 {
				err := c.gcp.DeletePrivateEndpoint(ctx, endpointService.GetEndpointGroupNames()[0], endpointService.GetServiceAttachmentNames()[0])
				if err != nil {
					fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC Endpoint at region %s from GCP: %s", endpointService.GetRegionName(), err))

					continue
				}

				endpointIDs = endpointService.GetEndpointGroupNames()
			}
		case CloudProviderAZURE:
			if len(endpointService.GetPrivateEndpoints()) > 0 {
				err := c.azure.DeletePrivateEndpoint(ctx, endpointService.GetPrivateEndpoints())
				if err != nil {
					fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete VPC Endpoint from Azure: %s", err))

					continue
				}

				endpointIDs = endpointService.GetPrivateEndpoints()
			}
		}

		for _, endpoint := range endpointIDs {
			_, err := c.client.PrivateEndpointServicesApi.DeletePrivateEndpoint(ctx, projectID, provider, endpoint, endpointService.GetId()).Execute()
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of private endpoint %s: %s", endpoint, err))

				continue
			}
		}

		if len(endpointIDs) == 0 {
			_, err := c.client.PrivateEndpointServicesApi.DeletePrivateEndpointService(ctx, projectID, provider, endpointService.GetId()).Execute()
			if err != nil {
				fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of private endpoint %s: %s", endpointService.GetId(), err))

				continue
			}
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of private endpoint %s", endpointService.GetId()))
	}
}
