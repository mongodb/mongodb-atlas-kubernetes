package atlasdeployment

import (
	"context"
	"fmt"

	"go.mongodb.org/atlas/mongodbatlas"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlasdeployment/globaldeployment"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util"
)

func EnsureManagedNamespaces(service *workflow.Context, groupID string, clusterType string, managedNamespace []mdbv1.ManagedNamespace, deploymentName string) workflow.Result {
	if clusterType != string(mdbv1.TypeGeoSharded) && managedNamespace != nil {
		return workflow.Terminate(workflow.ManagedNamespacesReady, "Managed namespace is only supported by GeoSharded clusters")
	}

	result := syncManagedNamespaces(context.Background(), service, groupID, deploymentName, managedNamespace)
	if !result.IsOk() {
		service.SetConditionFromResult(status.ManagedNamespacesReadyType, result)
		return result
	}

	if managedNamespace == nil {
		service.UnsetCondition(status.ManagedNamespacesReadyType)
		service.EnsureStatusOption(status.AtlasDeploymentManagedNamespacesOption(nil))
	} else {
		service.SetConditionTrue(status.ManagedNamespacesReadyType)
	}
	return result
}

func syncManagedNamespaces(ctx context.Context, service *workflow.Context, groupID string, deploymentName string, managedNamespaces []mdbv1.ManagedNamespace) workflow.Result {
	logger := service.Log
	existingManagedNamespaces, _, err := globaldeployment.GetGlobalDeploymentState(ctx, service.Client, groupID, deploymentName)
	logger.Debugf("Syncing managed namespaces %s", deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.ManagedNamespacesReady, fmt.Sprintf("Failed to get managed namespaces: %v", err))
	}
	diff := sortManagedNamespaces(existingManagedNamespaces, managedNamespaces)
	logger.Debugw("diff", "To create: %v", diff.ToCreate, "To delete: %v", diff.ToDelete, "To update status: %v", diff.ToUpdateStatus)
	err = deleteManagedNamespaces(ctx, service.Client, groupID, deploymentName, diff.ToDelete)
	if err != nil {
		logger.Errorf("failed to delete managed namespaces: %v", err)
		return workflow.Terminate(workflow.ManagedNamespacesReady, fmt.Sprintf("Failed to delete managed namespaces: %v", err))
	}
	nsStatuses := createManagedNamespaces(ctx, service.Client, groupID, deploymentName, diff.ToCreate)
	for _, ns := range diff.ToUpdateStatus {
		nsStatuses = append(nsStatuses, status.NewCreatedManagedNamespaceStatus(ns))
	}
	logger.Debugw("Managed namespaces statuses", "statuses", nsStatuses)

	service.EnsureStatusOption(status.AtlasDeploymentManagedNamespacesOption(nsStatuses))
	return checkManagedNamespaceStatus(nsStatuses)
}

func checkManagedNamespaceStatus(managedNamespaces []status.ManagedNamespace) workflow.Result {
	for _, ns := range managedNamespaces {
		if ns.Status != status.StatusReady {
			return workflow.Terminate(workflow.ManagedNamespacesReady, "Managed namespaces are not ready")
		}
	}
	return workflow.OK()
}

func sortManagedNamespaces(existing []globaldeployment.AtlasManagedNamespace, desired []mdbv1.ManagedNamespace) NamespaceDiff {
	var result NamespaceDiff
	for _, d := range desired {
		found := false
		for _, e := range existing {
			if isManagedNamespaceStateMatch(e, d) {
				found = true
				result.ToUpdateStatus = append(result.ToUpdateStatus, d.ToAtlas())
				break
			}
		}
		if !found {
			result.ToCreate = append(result.ToCreate, d.ToAtlas())
		}
	}

	for _, e := range existing {
		found := false
		for _, d := range desired {
			if isManagedNamespaceStateMatch(e, d) {
				found = true
				break
			}
		}
		if !found {
			result.ToDelete = append(result.ToDelete, e)
		}
	}

	return result
}

func isManagedNamespaceStateMatch(existing globaldeployment.AtlasManagedNamespace, desired mdbv1.ManagedNamespace) bool {
	if existing.DB == desired.Db &&
		existing.Collection == desired.Collection &&
		existing.CustomShardKey == desired.CustomShardKey &&
		util.PtrValuesEqual(existing.IsShardKeyUnique, desired.IsShardKeyUnique) &&
		util.PtrValuesEqual(existing.IsCustomShardKeyHashed, desired.IsCustomShardKeyHashed) &&
		existing.NumInitialChunks == desired.NumInitialChunks {
		return true
	}
	return false
}

func deleteManagedNamespaces(ctx context.Context, client mongodbatlas.Client, id string, name string, namespaces []globaldeployment.AtlasManagedNamespace) error {
	for i := range namespaces {
		_, err := globaldeployment.DeleteManagedNamespaceInAtlas(ctx, client, id, name, &namespaces[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func createManagedNamespaces(ctx context.Context, client mongodbatlas.Client, id string, name string, namespaces []globaldeployment.AtlasManagedNamespace) []status.ManagedNamespace {
	var newStatuses []status.ManagedNamespace
	for i := range namespaces {
		ns := namespaces[i]
		_, err := globaldeployment.CreateManagedNamespaceInAtlas(ctx, client, id, name, &ns)
		if err != nil {
			newStatuses = append(newStatuses, status.NewFailedToCreateManagedNamespaceStatus(ns, err))
		} else {
			newStatuses = append(newStatuses, status.NewCreatedManagedNamespaceStatus(ns))
		}
	}
	return newStatuses
}

type NamespaceDiff struct {
	ToCreate       []globaldeployment.AtlasManagedNamespace
	ToDelete       []globaldeployment.AtlasManagedNamespace
	ToUpdateStatus []globaldeployment.AtlasManagedNamespace
}
