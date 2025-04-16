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

package atlasdatafederation

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

func (r *AtlasDataFederationReconciler) ensurePrivateEndpoints(ctx *workflow.Context, service datafederation.DatafederationPrivateEndpointService, project *akov2.AtlasProject, dataFederation *akov2.AtlasDataFederation) workflow.Result {
	projectID := project.ID()
	fromAtlas, err := service.List(ctx.Context, projectID)
	if err != nil {
		return r.privateEndpointsFailed(ctx, err)
	}

	m, err := datafederation.MapDatafederationPrivateEndpoints(projectID, dataFederation, fromAtlas)
	if err != nil {
		return r.privateEndpointsFailed(ctx, err)
	}

	for _, endpoint := range m {
		endpointReconciler := &PrivateEndpointReconciler{service, endpoint}
		if err := endpointReconciler.Reconcile(ctx); err != nil {
			return r.privateEndpointsFailed(ctx, err)
		}
	}

	if len(dataFederation.Spec.PrivateEndpoints) == 0 {
		return r.privateEndpointsUnmanage(ctx)
	}

	return r.privateEndpointsIdle(ctx)
}

func (r *AtlasDataFederationReconciler) privateEndpointsFailed(ctx *workflow.Context, err error) workflow.Result {
	ctx.Log.Errorw("getAllDataFederationPEs error", "err", err.Error())
	result := workflow.Terminate(workflow.Internal, err)
	ctx.SetConditionFromResult(api.DataFederationPEReadyType, result)
	return result
}

func (r *AtlasDataFederationReconciler) privateEndpointsIdle(ctx *workflow.Context) workflow.Result {
	ctx.SetConditionTrue(api.DataFederationPEReadyType)
	return workflow.OK()
}

func (r *AtlasDataFederationReconciler) privateEndpointsUnmanage(ctx *workflow.Context) workflow.Result {
	ctx.UnsetCondition(api.DataFederationPEReadyType)
	return workflow.OK()
}

type PrivateEndpointReconciler struct {
	service  datafederation.DatafederationPrivateEndpointService
	endpoint *datafederation.DataFederationPrivateEndpoint
}

func (r *PrivateEndpointReconciler) Reconcile(ctx *workflow.Context) error {
	inAKO := r.endpoint.AKO != nil
	inAtlas := r.endpoint.Atlas != nil
	inLastApplied := r.endpoint.LastApplied != nil

	switch {
	case inAKO && !inAtlas:
		return r.create(ctx)
	case inAKO:
		return r.update(ctx)
	case inAtlas && inLastApplied:
		// only delete private endpoints that used to be tracked in AKO
		return r.delete(ctx)
	}

	return nil
}

func (r *PrivateEndpointReconciler) create(ctx *workflow.Context) error {
	if err := r.service.Create(ctx.Context, r.endpoint.AKO); err != nil {
		return fmt.Errorf("error creating private endpoint: %w", err)
	}
	return nil
}

func (r *PrivateEndpointReconciler) delete(ctx *workflow.Context) error {
	if err := r.service.Delete(ctx.Context, r.endpoint.Atlas); err != nil {
		return fmt.Errorf("error deleting private endpoint: %w", err)
	}
	return nil
}

func (r *PrivateEndpointReconciler) update(ctx *workflow.Context) error {
	if r.endpoint.AKO.EqualsTo(r.endpoint.Atlas) {
		return nil
	}
	if err := r.delete(ctx); err != nil {
		return err
	}
	if err := r.create(ctx); err != nil {
		return err
	}
	return nil
}
