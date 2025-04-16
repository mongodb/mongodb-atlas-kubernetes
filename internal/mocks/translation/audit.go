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

package translation

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
)

type AuditLogMock struct {
	GetFunc    func(projectID string) (*audit.AuditConfig, error)
	UpdateFunc func(projectID string, auditing *audit.AuditConfig) error
}

func (c *AuditLogMock) Get(_ context.Context, projectID string) (*audit.AuditConfig, error) {
	return c.GetFunc(projectID)
}
func (c *AuditLogMock) Update(_ context.Context, projectID string, auditing *audit.AuditConfig) error {
	return c.UpdateFunc(projectID, auditing)
}
