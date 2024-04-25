package deployment

import (
	"context"
	"errors"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translayer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/atlas"
)

type Service struct {
	mongodbatlas.AdvancedClustersService
	mongodbatlas.ClustersService
}

func NewService(ctx context.Context, provider atlas.Provider, secretRef *types.NamespacedName, log *zap.SugaredLogger) (*Service, error) {
	client, err := translayer.NewLegacyClient(ctx, provider, secretRef, log)
	if err != nil {
		return nil, err
	}
	return &Service{AdvancedClustersService: client.AdvancedClusters, ClustersService: client.Clusters}, nil
}

func (ds *Service) Exists(ctx context.Context, projectID, clusterName string) (bool, error) {
	var apiError *mongodbatlas.ErrorResponse
	_, _, err := ds.AdvancedClustersService.Get(ctx, projectID, clusterName)
	if errors.As(err, &apiError) && apiError.ErrorCode == atlas.ClusterNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (ds *Service) IsReady(ctx context.Context, projectID, deploymentName string) (bool, error) {
	resourceStatus, _, err := ds.ClustersService.Status(ctx, projectID, deploymentName)
	if err != nil {
		return false, err
	}
	return resourceStatus.ChangeStatus == mongodbatlas.ChangeStatusApplied, nil
}
