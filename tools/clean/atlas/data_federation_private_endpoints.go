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
	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

func (c *Cleaner) listFederatedDBPrivateEndpoints(ctx context.Context, projectID string) []admin.PrivateNetworkEndpointIdEntry {
	federatedDBPEs, _, err := c.client.DataFederationApi.ListPrivateEndpointIds(ctx, projectID).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list federated databases for project %s: %s", projectID, err))

		return nil
	}

	return federatedDBPEs.GetResults()
}

func (c *Cleaner) deleteFederatedDBPrivateEndpoints(ctx context.Context, projectID string, dbpes []admin.PrivateNetworkEndpointIdEntry) {
	for _, fedDBPE := range dbpes {
		_, err := c.client.DataFederationApi.DeletePrivateEndpointId(ctx, projectID, fedDBPE.GetEndpointId()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of Federated DB private endpoint %s: %s", fedDBPE.GetEndpointId(), err))
		}
		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of Federated DB private endpoint %s", fedDBPE.GetEndpointId()))
	}
}
