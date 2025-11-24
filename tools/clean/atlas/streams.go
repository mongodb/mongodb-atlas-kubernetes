// Copyright 2026 MongoDB Inc
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

func (c *Cleaner) listStreams(ctx context.Context, id string) []admin.StreamsTenant {
	streams, _, err := c.client.StreamsApi.ListStreamInstances(ctx, id).Execute()
	if err != nil {
		fmt.Println(text.FgRed.Sprintf("\tFailed to list stream instances for project %s: %s", id, err))
		return nil
	}
	return streams.GetResults()
}

func (c *Cleaner) deleteStreams(ctx context.Context, id string, streams []admin.StreamsTenant) {
	for _, stream := range streams {
		_, err := c.client.StreamsApi.DeleteStreamInstance(ctx, id, stream.GetName()).Execute()
		if err != nil {
			fmt.Println(text.FgRed.Sprintf("\t\t\tFailed to delete Stream instance %s", stream.GetId()), err)

			continue
		}
	}
}
