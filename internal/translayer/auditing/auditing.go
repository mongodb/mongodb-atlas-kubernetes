package auditing

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type AtlasAuditService interface {
	Get(ctx context.Context, projectID string) (*Auditing, error)
	Set(ctx context.Context, projectID string, auditing *Auditing) error
}

type ProductionAtlasAudit struct {
	auditAPI admin.AuditingApi
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*ProductionAtlasAudit, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewProductionAtlasAudit(client.AuditingApi), nil
}

func NewProductionAtlasAudit(api admin.AuditingApi) *ProductionAtlasAudit {
	return &ProductionAtlasAudit{auditAPI: api}
}

func (s *ProductionAtlasAudit) Get(ctx context.Context, projectID string) (*Auditing, error) {
	auditLog, _, err := s.auditAPI.GetAuditingConfiguration(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log from Atlas: %w", err)
	}
	return fromAtlas(auditLog)
}

func (s *ProductionAtlasAudit) Set(ctx context.Context, projectID string, auditing *Auditing) error {
	_, _, err := s.auditAPI.UpdateAuditingConfiguration(ctx, projectID, toAtlas(auditing)).Execute()
	return err
}
