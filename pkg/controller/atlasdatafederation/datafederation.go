package atlasdatafederation

import (
	"context"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/compat"
	"go.uber.org/zap"
)

func (r *AtlasDataFederationReconciler) ensureDataFederation(ctx *workflow.Context, project *mdbv1.AtlasProject, dataFederation *mdbv1.AtlasDataFederation) workflow.Result {
	log := ctx.Log

	clientDF := NewClient(ctx.Client)

	projectID := project.ID()
	operatorSpec := &dataFederation.Spec

	atlasSpec, resp, err := clientDF.Get(context.Background(), projectID, operatorSpec.Name)
	if err != nil {
		if resp == nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusUnauthorized {
			return workflow.Terminate(workflow.DataFederationNotCreatedInAtlas, err.Error())
		}

		_, _, err := clientDF.Create(context.Background(), projectID, operatorSpec)
		if err != nil {
			return workflow.Terminate(workflow.DataFederationNotCreatedInAtlas, err.Error())
		}

		return workflow.InProgress(workflow.DataFederationCreating, "Data Federation is being created")
	}

	if areEqual, _ := dataFederationEqual(*atlasSpec, *operatorSpec, log); areEqual {
		return workflow.OK()
	}

	_, _, err = clientDF.Update(context.Background(), projectID, operatorSpec)
	if err != nil {
		return workflow.Terminate(workflow.DataFederationNotUpdatedInAtlas, err.Error())
	}

	return workflow.InProgress(workflow.DataFederationUpdating, "Data Federation is being updated")
}

func dataFederationEqual(atlasSpec, operatorSpec mdbv1.DataFederationSpec, log *zap.SugaredLogger) (areEqual bool, diff string) {
	mergedSpec, err := getMergedSpec(atlasSpec, operatorSpec)
	if err != nil {
		log.Errorf("failed to merge Data Federation specs: %s", err.Error())
		return false, ""
	}

	d := cmp.Diff(atlasSpec, mergedSpec, cmpopts.EquateEmpty())
	if d != "" {
		log.Debugf("Data Federation diff: \n%s", d)
	}

	return d == "", d
}

func getMergedSpec(atlasSpec, operatorSpec mdbv1.DataFederationSpec) (mdbv1.DataFederationSpec, error) {
	mergedSpec := mdbv1.DataFederationSpec{}

	if err := compat.JSONCopy(&mergedSpec, atlasSpec); err != nil {
		return mergedSpec, err
	}
	if err := compat.JSONCopy(&mergedSpec, operatorSpec); err != nil {
		return mergedSpec, err
	}

	mergedSpec.Project = common.ResourceRefNamespaced{}

	return mergedSpec, nil
}
