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

package atlasdeployment

import (
	"context"
	"errors"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compare"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/deployment"
)

func (r *AtlasDeploymentReconciler) ensureManagedNamespaces(service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, groupID string, clusterType string, managedNamespace []akov2.ManagedNamespace, deploymentName string) workflow.Result {
	if clusterType != string(akov2.TypeGeoSharded) && managedNamespace != nil {
		return workflow.Terminate(workflow.ManagedNamespacesReady, errors.New("managed namespace is only supported by GeoSharded clusters"))
	}

	result := r.syncManagedNamespaces(service, deploymentService, groupID, deploymentName, managedNamespace)
	if !result.IsOk() {
		service.SetConditionFromResult(api.ManagedNamespacesReadyType, result)
		return result
	}

	if managedNamespace == nil {
		service.UnsetCondition(api.ManagedNamespacesReadyType)
		service.EnsureStatusOption(status.AtlasDeploymentManagedNamespacesOption(nil))
	} else {
		service.SetConditionTrue(api.ManagedNamespacesReadyType)
	}
	return result
}

func (r *AtlasDeploymentReconciler) syncManagedNamespaces(service *workflow.Context, deploymentService deployment.AtlasDeploymentsService, groupID string, deploymentName string, managedNamespaces []akov2.ManagedNamespace) workflow.Result {
	logger := service.Log
	existingManagedNamespaces, err := deploymentService.GetManagedNamespaces(service.Context, groupID, deploymentName)
	logger.Debugf("Syncing managed namespaces %s", deploymentName)
	if err != nil {
		return workflow.Terminate(workflow.ManagedNamespacesReady, fmt.Errorf("failed to get managed namespaces: %w", err))
	}
	diff := sortManagedNamespaces(existingManagedNamespaces, managedNamespaces)
	logger.Debugw("diff", "To create: %v", diff.ToCreate, "To delete: %v", diff.ToDelete, "To update status: %v", diff.ToUpdateStatus)
	err = deleteManagedNamespaces(service.Context, deploymentService, groupID, deploymentName, diff.ToDelete)
	if err != nil {
		logger.Errorf("failed to delete managed namespaces: %v", err)
		return workflow.Terminate(workflow.ManagedNamespacesReady, fmt.Errorf("failed to delete managed namespaces: %w", err))
	}
	nsStatuses := createManagedNamespaces(service.Context, deploymentService, groupID, deploymentName, diff.ToCreate)
	for _, ns := range diff.ToUpdateStatus {
		nsStatuses = append(nsStatuses, akov2.NewCreatedManagedNamespaceStatus(ns))
	}
	logger.Debugw("Managed namespaces statuses", "statuses", nsStatuses)

	service.EnsureStatusOption(status.AtlasDeploymentManagedNamespacesOption(nsStatuses))
	return checkManagedNamespaceStatus(nsStatuses)
}

func checkManagedNamespaceStatus(managedNamespaces []status.ManagedNamespace) workflow.Result {
	for _, ns := range managedNamespaces {
		if ns.Status != status.StatusReady {
			return workflow.Terminate(workflow.ManagedNamespacesReady, errors.New("managed namespaces are not ready"))
		}
	}
	return workflow.OK()
}

func sortManagedNamespaces(existing, desired []akov2.ManagedNamespace) NamespaceDiff {
	var result NamespaceDiff
	for _, d := range desired {
		found := false
		for _, e := range existing {
			if isManagedNamespaceStateMatch(e, d) {
				found = true
				result.ToUpdateStatus = append(result.ToUpdateStatus, d)
				break
			}
		}
		if !found {
			result.ToCreate = append(result.ToCreate, d)
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

func isManagedNamespaceStateMatch(existing, desired akov2.ManagedNamespace) bool {
	if existing.Db == desired.Db &&
		existing.Collection == desired.Collection &&
		existing.CustomShardKey == desired.CustomShardKey &&
		compare.PtrValuesEqual(existing.IsShardKeyUnique, desired.IsShardKeyUnique) &&
		compare.PtrValuesEqual(existing.IsCustomShardKeyHashed, desired.IsCustomShardKeyHashed) &&
		existing.NumInitialChunks == desired.NumInitialChunks {
		return true
	}
	return false
}

func deleteManagedNamespaces(ctx context.Context, client deployment.AtlasDeploymentsService, id string, name string, namespaces []akov2.ManagedNamespace) error {
	for i := range namespaces {
		err := client.DeleteManagedNamespace(ctx, id, name, &namespaces[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func createManagedNamespaces(ctx context.Context, client deployment.AtlasDeploymentsService, id string, name string, namespaces []akov2.ManagedNamespace) []status.ManagedNamespace {
	var newStatuses []status.ManagedNamespace
	for i := range namespaces {
		ns := namespaces[i]
		err := client.CreateManagedNamespace(ctx, id, name, &ns)

		if err != nil {
			newStatuses = append(newStatuses, akov2.NewFailedToCreateManagedNamespaceStatus(ns, err))
		} else {
			newStatuses = append(newStatuses, akov2.NewCreatedManagedNamespaceStatus(ns))
		}
	}
	return newStatuses
}

type NamespaceDiff struct {
	ToCreate       []akov2.ManagedNamespace
	ToDelete       []akov2.ManagedNamespace
	ToUpdateStatus []akov2.ManagedNamespace
}
