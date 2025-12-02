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
	"go.mongodb.org/atlas-sdk/v20250312009/admin"
)

func (c *Cleaner) listFederatedDatabases(ctx context.Context, projectID string) []admin.DataLakeTenant {
	federatedDBs, _, err := c.client.DataFederationApi.
		ListDataFederation(ctx, projectID).
		Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list federated databases for project %s: %s", projectID, err))

		return nil
	}

	return federatedDBs
}

func (c *Cleaner) deleteFederatedDatabases(ctx context.Context, projectID string, dbs []admin.DataLakeTenant) {
	for _, fedDB := range dbs {
		if fedDB.GetState() == "DELETED" {
			fmt.Println(text.FgHiBlue.Sprintf("\t\t\tFederated Database %s is being deleted...", fedDB.GetName()))

			continue
		}

		_, err := c.client.DataFederationApi.DeleteDataFederation(ctx, projectID, fedDB.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to request deletion of Federated Database %s: %s", fedDB.GetName(), err))
		}

		fmt.Println(text.FgBlue.Sprintf("\t\t\tRequested deletion of Federated Database %s", fedDB.GetName()))
	}
}
