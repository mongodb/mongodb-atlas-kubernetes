package audit

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

// AuditLogService is the interface exposed by this translation layer over
// the Atlas AuditLog
type AuditLogService interface {
	Get(ctx context.Context, projectID string) (*AuditConfig, error)
	Set(ctx context.Context, projectID string, auditing *AuditConfig) error
}

// AuditLog is the default implementation of the AuditLogService using the Atlas SDK
type AuditLog struct {
	auditAPI admin.AuditingApi
}

// NewAuditLogService creates an AuditLog from credentials and the atlas provider
func NewAuditLogService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*AuditLog, error) {
	client, err := translation.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewAuditLog(client.AuditingApi), nil
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
	return fromAtlas(auditLog)
}

// Set an Atlas Project audit log configuration
func (s *AuditLog) Set(ctx context.Context, projectID string, auditing *AuditConfig) error {
	_, _, err := s.auditAPI.UpdateAuditingConfiguration(ctx, projectID, toAtlas(auditing)).Execute()
	return err
}
