package translation

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas-sdk/v20231115008/admin"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/atlas"
)

func NewVersionedClient(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*admin.APIClient, error) {
	apiClientSet, _, err := provider.SdkClientSet(ctx, secretRef, log)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Versioned Atlas client: %w", err)
	}
	return apiClientSet.SdkClient20231115008, nil
}
