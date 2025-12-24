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
	"go.mongodb.org/atlas-sdk/v20250312011/admin"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

type AuditingConfigType string

const (
	FilterDefault = "{}"
)

// AuditConfig represents the Atlas Project audit log config
type AuditConfig struct {
	*akov2.Auditing
}

func NewAuditConfig(auditConfig *akov2.Auditing) *AuditConfig {
	if auditConfig == nil {
		auditConfig = &akov2.Auditing{}
	}

	if auditConfig.AuditFilter == "" {
		auditConfig.AuditFilter = FilterDefault
	}

	return &AuditConfig{
		Auditing: auditConfig,
	}
}

func toAtlas(auditing *AuditConfig) *admin.AuditLog {
	auditLog := admin.NewAuditLogWithDefaults()
	auditLog.SetEnabled(auditing.Enabled)
	auditLog.SetAuditAuthorizationSuccess(auditing.AuditAuthorizationSuccess)
	auditLog.SetAuditFilter(auditing.AuditFilter)
	// ConfigurationType is not set on the PATCH operation to Atlas

	return auditLog
}

func fromAtlas(auditLog *admin.AuditLog) *AuditConfig {
	auditFilter := FilterDefault

	if auditLog.GetAuditFilter() != "" {
		auditFilter = auditLog.GetAuditFilter()
	}

	return &AuditConfig{
		Auditing: &akov2.Auditing{
			Enabled:                   auditLog.GetEnabled(),
			AuditAuthorizationSuccess: auditLog.GetAuditAuthorizationSuccess(),
			AuditFilter:               auditFilter,
		},
	}
}
