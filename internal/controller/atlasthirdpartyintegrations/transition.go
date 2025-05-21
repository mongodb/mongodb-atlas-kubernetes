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

package integrations

import (
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/reconciler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"

	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

func (r *AtlasThirdPartyIntegrationsReconciler) release(workflowCtx *workflow.Context, integration *akov2next.AtlasThirdPartyIntegration, err error) (ctrl.Result, error) {
	if errors.Is(err, reconciler.ErrMissingKubeProject) {
		if finalizerErr := customresource.ManageFinalizer(workflowCtx.Context, r.Client, integration, customresource.UnsetFinalizer); finalizerErr != nil {
			err = errors.Join(err, finalizerErr)
		}
	}
	return r.terminate(workflowCtx, integration, workflow.NetworkPeeringNotConfigured, err)
}

func (r *AtlasThirdPartyIntegrationsReconciler) terminate(
	ctx *workflow.Context,
	resource api.AtlasCustomResource,
	reason workflow.ConditionReason,
	err error,
) (ctrl.Result, error) {
	condition := api.ReadyType
	r.Log.Errorf("resource %T(%s/%s) failed on condition %s: %s",
		resource, resource.GetNamespace(), resource.GetName(), condition, err)
	result := workflow.Terminate(reason, err)
	ctx.SetConditionFalse(api.ReadyType).SetConditionFromResult(condition, result)

	return result.ReconcileResult(), nil
}
