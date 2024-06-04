package translation

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

func NewVersionedClient(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*admin.APIClient, error) {
	apiClient, _, err := provider.SdkClient(ctx, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Versioned Atlas client: %w", err)
	}
	return apiClient, nil
}
