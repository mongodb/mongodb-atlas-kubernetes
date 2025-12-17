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

package audit

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20250312011/admin"
)

// AuditLogService is the interface exposed by this translation layer over
// the Atlas AuditLog
type AuditLogService interface {
	Get(ctx context.Context, projectID string) (*AuditConfig, error)
	Update(ctx context.Context, projectID string, auditing *AuditConfig) error
}

// AuditLog is the default implementation of the AuditLogService using the Atlas SDK
type AuditLog struct {
	auditAPI admin.AuditingApi
}

// NewAuditLog wraps the SDK AuditingApi as an AuditLog
func NewAuditLog(api admin.AuditingApi) *AuditLog {
	return &AuditLog{auditAPI: api}
}

// Get an Atlas Project audit log configuration
func (s *AuditLog) Get(ctx context.Context, projectID string) (*AuditConfig, error) {
	auditLog, _, err := s.auditAPI.GetGroupAuditLog(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log from Atlas: %w", err)
	}

	return fromAtlas(auditLog), nil
}

// Update an Atlas Project audit log configuration
func (s *AuditLog) Update(ctx context.Context, projectID string, auditing *AuditConfig) error {
	_, _, err := s.auditAPI.UpdateAuditLog(ctx, projectID, toAtlas(auditing)).Execute()
	return err
}
