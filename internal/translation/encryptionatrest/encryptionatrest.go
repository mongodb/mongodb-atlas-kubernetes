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

package encryptionatrest

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312010/admin"
)

type EncryptionAtRestService interface {
	Get(context.Context, string) (*EncryptionAtRest, error)
	Update(context.Context, string, EncryptionAtRest) error
}

type EncryptionAtRestAPI struct {
	encryptionAtRestAPI admin.EncryptionAtRestUsingCustomerKeyManagementApi
}

func NewEncryptionAtRestAPI(api admin.EncryptionAtRestUsingCustomerKeyManagementApi) *EncryptionAtRestAPI {
	return &EncryptionAtRestAPI{encryptionAtRestAPI: api}
}

func (e *EncryptionAtRestAPI) Get(ctx context.Context, projectID string) (*EncryptionAtRest, error) {
	result, _, err := e.encryptionAtRestAPI.GetEncryptionAtRest(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get encryption at rest from Atlas: %w", err)
	}

	return fromAtlas(result), nil
}

func (e *EncryptionAtRestAPI) Update(ctx context.Context, projectID string, ear EncryptionAtRest) error {
	a := toAtlas(&ear)
	_, _, err := e.encryptionAtRestAPI.UpdateEncryptionAtRest(ctx, projectID, a).Execute()
	if err != nil {
		return fmt.Errorf("failed to update encryption at rest in Atlas: %w", err)
	}

	return nil
}
