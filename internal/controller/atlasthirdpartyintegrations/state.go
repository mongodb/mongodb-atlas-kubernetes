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
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/statushandler"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	akov2next "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/v1"
)

const (
	AtlasThirdPartyIntegration = "AtlasThirdPartyIntegration"
)

func (r *AtlasThirdPartyIntegrationsReconciler) handleCustomResource(ctx context.Context, integration *akov2next.AtlasThirdPartyIntegration) (ctrl.Result, error) {
	if customresource.ReconciliationShouldBeSkipped(integration) {
		return r.Skip(ctx, AtlasThirdPartyIntegration, integration, &integration.Spec)
	}

	conditions := api.InitCondition(integration, api.FalseCondition(api.ReadyType))
	workflowCtx := workflow.NewContext(r.Log, conditions, ctx, integration)
	defer statushandler.Update(workflowCtx, r.Client, r.EventRecorder, integration)

	isValid := customresource.ValidateResourceVersion(workflowCtx, integration, r.Log)
	if !isValid.IsOk() {
		return r.Invalidate(AtlasThirdPartyIntegration, isValid)
	}

	if !r.AtlasProvider.IsResourceSupported(integration) {
		return r.Unsupport(workflowCtx, AtlasThirdPartyIntegration)
	}

	connectionConfig, err := r.ResolveConnectionConfig(ctx, integration)
	if err != nil {
		return r.release(workflowCtx, integration, err)
	}
	sdkClientSet, err := r.AtlasProvider.SdkClientSet(ctx, connectionConfig.Credentials, r.Log)
	if err != nil {
		return r.terminate(workflowCtx, integration, workflow.AtlasAPIAccessNotConfigured, err)
	}
	_, err = r.ResolveProject(ctx, sdkClientSet.SdkClient20231115008, integration)
	if err != nil {
		return r.release(workflowCtx, integration, err)
	}
	return ctrl.Result{}, nil
}
