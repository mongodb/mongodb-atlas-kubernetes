package audit

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type Service interface {
	Get(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingConfig, error)
	Set(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingConfig) error
}

type service struct {
	admin.AuditingApi
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (Service, error) {
	client, err := translayer.NewVersionedClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return NewFromAuditingAPI(client.AuditingApi), nil
}

func NewFromAuditingAPI(api admin.AuditingApi) *service {
	return &service{AuditingApi: api}
}

func (s *service) Get(ctx context.Context, projectID string) (*v1alpha1.AtlasAuditingConfig, error) {
	auditLog, _, err := s.AuditingApi.GetAuditingConfiguration(ctx, projectID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log from Atlas: %w", err)
	}
	return fromAtlas(auditLog), nil
}

func (s *service) Set(ctx context.Context, projectID string, auditing *v1alpha1.AtlasAuditingConfig) error {
	_, _, err := s.AuditingApi.UpdateAuditingConfiguration(ctx, projectID, toAtlas(auditing)).Execute()
	return err
}
