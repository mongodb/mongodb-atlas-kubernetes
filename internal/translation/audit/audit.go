package audit

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
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
	auditLog, _, err := s.auditAPI.GetAuditingConfiguration(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log from Atlas: %w", err)
	}

	return fromAtlas(auditLog), nil
}

// Update an Atlas Project audit log configuration
func (s *AuditLog) Update(ctx context.Context, projectID string, auditing *AuditConfig) error {
	_, _, err := s.auditAPI.UpdateAuditingConfiguration(ctx, projectID, toAtlas(auditing)).Execute()
	return err
}
