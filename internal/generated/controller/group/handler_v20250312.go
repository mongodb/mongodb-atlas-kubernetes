// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package group

import (
	"context"
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/translate"
	akov2generated "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1"
	ctrlstate "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/state"
	result "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/result"
	state "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/state"
	"go.mongodb.org/atlas-sdk/v20250312008/admin"
	controllerruntime "sigs.k8s.io/controller-runtime"
	builder "sigs.k8s.io/controller-runtime/pkg/builder"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	controller "sigs.k8s.io/controller-runtime/pkg/controller"
	reconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type GroupHandlerv20250312 struct {
	client             client.Client
	atlasClient        *admin.APIClient
	translationRequest *translate.Request
}

func NewGroupHandlerv20250312(client client.Client, atlasClient *admin.APIClient, translatorRequest *translate.Request) *GroupHandlerv20250312 {
	return &GroupHandlerv20250312{
		client:             client,
		atlasClient:        atlasClient,
		translationRequest: translatorRequest,
	}
}

// HandleInitial handles the initial state for version v20250312
func (h *GroupHandlerv20250312) HandleInitial(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	atlasGroup := &admin.Group{}
	err := translate.ToAPI(h.translationRequest, atlasGroup, group)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate group to Atlas: %w", err))
	}

	params := &admin.CreateGroupApiParams{Group: atlasGroup, ProjectOwnerId: group.Spec.V20250312.ProjectOwnerId}
	response, _, err := h.atlasClient.ProjectsApi.CreateGroupWithParams(ctx, params).Execute()
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to create project: %w", err))
	}

	objsFromAtlas, err := translate.FromAPI(h.translationRequest, group.DeepCopy(), response)
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to translate group from Atlas: %w", err))
	}

	groupFromAtlas := objsFromAtlas[0].(*akov2generated.Group)
	err = h.client.Status().Patch(ctx, groupFromAtlas, client.MergeFrom(group))
	if err != nil {
		return result.Error(state.StateInitial, fmt.Errorf("failed to patch group status: %w", err))
	}

	return result.NextState(state.StateCreated, "Project created.")
}

// HandleImportRequested handles the importrequested state for version v20250312
func (h *GroupHandlerv20250312) HandleImportRequested(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	// TODO: Implement importrequested state logic
	// TODO: Use h.atlasProvider.SdkClientSet(ctx, h.globalSecretRef, h.log) to get Atlas SDK client
	return result.NextState(state.StateImported, "Import completed")
}

// HandleImported handles the imported state for version v20250312
func (h *GroupHandlerv20250312) HandleImported(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Ready")
}

// HandleCreating handles the creating state for version v20250312
func (h *GroupHandlerv20250312) HandleCreating(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return result.NextState(state.StateCreated, "Resource created")
}

// HandleCreated handles the created state for version v20250312
func (h *GroupHandlerv20250312) HandleCreated(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Ready")
}

// HandleUpdating handles the updating state for version v20250312
func (h *GroupHandlerv20250312) HandleUpdating(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return result.NextState(state.StateUpdated, "Update completed")
}

// HandleUpdated handles the updated state for version v20250312
func (h *GroupHandlerv20250312) HandleUpdated(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	//atlasClient, err := h.atlasSDKClient(ctx, group)
	//if err != nil {
	//	return result.Error(state.StateUpdated, fmt.Errorf("failed to setup Atlas SDK client: %w", err))
	//}
	//response, _, err := atlasClient.ProjectsApi.GetGroup(ctx, *group.Status.V20250312.Id).Execute()
	//if err != nil {
	//	return result.Error(state.StateUpdated, fmt.Errorf("failed to get group: %w", err))
	//}
	//
	//
	//generationChanged := meta.FindStatusCondition(*group.Status.Conditions, state.StateCondition).ObservedGeneration != group.GetGeneration()
	//hasError := meta.FindStatusCondition(*group.Status.Conditions, state.ReadyCondition).Reason == ctrlstate.ReadyReasonError
	//shouldReapply, err := reapply.ShouldReapply(obj)
	//if err != nil {
	//	return result.Error(currentState, err)
	//}
	//
	//switch {
	//case generationChanged:
	//case shouldReapply:
	//case hasError:
	//default:
	//	return result.NextState(currentState, "Upserted group.")
	//}
	//
	//params := &atlas20250312002.UpdateProjectApiParams{GroupUpdate: &atlas20250312002.GroupUpdate{}}
	//json.MustUnmarshal(json.MustMarshal(obj.Spec.V20250312.Entry), params.GroupUpdate)
	//params.GroupId = *obj.Status.V20250312.Id
	//
	//response, _, err = atlasClients.SdkClient20250312002.ProjectsApi.UpdateProjectWithParams(ctx, params).Execute()
	//if err != nil {
	//	return result.Error(currentState, fmt.Errorf("failed to update group: %w", err))
	//}
	//
	//json.MustUnmarshal(json.MustMarshal(response), obj.Status.V20250312)
	//err = r.client.Status().Patch(ctx, obj, client.RawPatch(types.MergePatchType, json.MustMarshal(obj)))
	//if err != nil {
	//	return result.Error(currentState, fmt.Errorf("failed to patch group: %w", err))
	//}

	return result.NextState(state.StateUpdated, "Ready")
}

// HandleDeletionRequested handles the deletionrequested state for version v20250312
func (h *GroupHandlerv20250312) HandleDeletionRequested(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	if group.Status.V20250312 == nil || group.Status.V20250312.Id == nil {
		return result.NextState(state.StateDeletionRequested, "Project deleted.")
	}

	_, err := h.atlasClient.ProjectsApi.DeleteGroup(ctx, *group.Status.V20250312.Id).Execute()
	if admin.IsErrorCode(err, "GROUP_NOT_FOUND") {
		return result.NextState(state.StateDeletionRequested, "Project deleted.")
	}
	if err != nil {
		return result.Error(state.StateDeletionRequested, fmt.Errorf("failed to delete project: %w", err))
	}

	return result.NextState(state.StateDeleted, "Deleted")
}

// HandleDeleting handles the deleting state for version v20250312
func (h *GroupHandlerv20250312) HandleDeleting(ctx context.Context, group *akov2generated.Group) (ctrlstate.Result, error) {
	return result.NextState(state.StateDeleted, "Deleted")
}

// For returns the resource and predicates for the controller
func (h *GroupHandlerv20250312) For() (client.Object, builder.Predicates) {
	return &akov2generated.Group{}, builder.WithPredicates()
}

// SetupWithManager sets up the controller with the Manager
func (h *GroupHandlerv20250312) SetupWithManager(mgr controllerruntime.Manager, rec reconcile.Reconciler, defaultOptions controller.Options) error {
	// This method is not used for version-specific handlers but required by StateHandler interface
	return nil
}
