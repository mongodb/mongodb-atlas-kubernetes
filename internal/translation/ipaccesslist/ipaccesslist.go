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

package ipaccesslist

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/atlas-sdk/v20250312012/admin"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/paging"
)

type IPAccessListService interface {
	List(ctx context.Context, projectID string) (IPAccessEntries, error)
	Add(ctx context.Context, projectID string, entries IPAccessEntries) error
	Delete(ctx context.Context, projectID string, entry *IPAccessEntry) error
	Status(ctx context.Context, projectID string, entry *IPAccessEntry) (string, error)
}

type IPAccessList struct {
	ipAccessListAPI admin.ProjectIPAccessListApi
}

func (i *IPAccessList) List(ctx context.Context, projectID string) (IPAccessEntries, error) {
	netPermResult, err := paging.ListAll(ctx, func(ctx context.Context, pageNum int) (paging.Response[admin.NetworkPermissionEntry], *http.Response, error) {
		return i.ipAccessListAPI.ListAccessListEntries(ctx, projectID).PageNum(pageNum).Execute()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ip access list from Atlas: %w", err)
	}

	return fromAtlas(netPermResult), nil
}

func (i *IPAccessList) Add(ctx context.Context, projectID string, entries IPAccessEntries) error {
	_, _, err := i.ipAccessListAPI.CreateAccessListEntry(ctx, projectID, toAtlas(entries)).Execute()
	if err != nil {
		return fmt.Errorf("failed to create ip access list from Atlas: %w", err)
	}

	return nil
}

func (i *IPAccessList) Delete(ctx context.Context, projectID string, entry *IPAccessEntry) error {
	_, err := i.ipAccessListAPI.DeleteAccessListEntry(ctx, projectID, entry.ID()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete ip access list from Atlas: %w", err)
	}

	return nil
}

func (i *IPAccessList) Status(ctx context.Context, projectID string, entry *IPAccessEntry) (string, error) {
	result, _, err := i.ipAccessListAPI.GetAccessListStatus(ctx, projectID, entry.ID()).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to get status of ip access list from Atlas: %w", err)
	}

	return result.GetSTATUS(), nil
}

func NewIPAccessList(api admin.ProjectIPAccessListApi) *IPAccessList {
	return &IPAccessList{
		ipAccessListAPI: api,
	}
}
