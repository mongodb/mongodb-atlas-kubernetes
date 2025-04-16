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
	"errors"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/translation/datafederation"
)

func (r *AtlasDataFederationReconciler) ensureDataFederation(ctx *workflow.Context, project *akov2.AtlasProject, dataFederation *akov2.AtlasDataFederation, federationService datafederation.DataFederationService) workflow.Result {
	projectID := project.ID()
	operatorSpec := &dataFederation.Spec

	akoDataFederation, err := datafederation.NewDataFederation(&dataFederation.Spec, projectID, nil)
	if err != nil {
		return workflow.Terminate(workflow.Internal, err)
	}

	atlasDataFederation, err := federationService.Get(ctx.Context, projectID, operatorSpec.Name)
	if err != nil {
		if !errors.Is(err, datafederation.ErrorNotFound) {
			return workflow.Terminate(workflow.Internal, err)
		}

		err = federationService.Create(ctx.Context, akoDataFederation)
		if err != nil {
			return workflow.Terminate(workflow.DataFederationNotCreatedInAtlas, err)
		}

		return workflow.InProgress(workflow.DataFederationCreating, "Data Federation is being created")
	}

	if akoDataFederation.SpecEqualsTo(atlasDataFederation) {
		return workflow.OK()
	}

	err = federationService.Update(ctx.Context, akoDataFederation)
	if err != nil {
		return workflow.Terminate(workflow.DataFederationNotUpdatedInAtlas, err)
	}

	return workflow.InProgress(workflow.DataFederationUpdating, "Data Federation is being updated")
}
