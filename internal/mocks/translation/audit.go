package translation

import (
	"context"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/audit"
)

type AuditLogMock struct {
	GetFunc func(projectID string) (*audit.AuditConfig, error)
	SetFunc func(projectID string, auditing *audit.AuditConfig) error
}

func (c *AuditLogMock) Get(_ context.Context, projectID string) (*audit.AuditConfig, error) {
	return c.GetFunc(projectID)
}
func (c *AuditLogMock) Set(_ context.Context, projectID string, auditing *audit.AuditConfig) error {
	return c.SetFunc(projectID, auditing)
}
