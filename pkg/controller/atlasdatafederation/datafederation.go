package atlasdatafederation

import (
	"net/http"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/compat"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func (r *AtlasDataFederationReconciler) ensureDataFederation(ctx *workflow.Context, project *akov2.AtlasProject, dataFederation *akov2.AtlasDataFederation) workflow.Result {
	log := ctx.Log

	projectID := project.ID()
	operatorSpec := &dataFederation.Spec

	dataFederationToAtlas, err := dataFederation.ToAtlas()
	if err != nil {
		return workflow.Terminate(workflow.Internal, "can not convert DataFederation (operator -> atlas)")
	}

	atlasSpec, resp, err := ctx.Client.DataFederation.Get(ctx.Context, projectID, operatorSpec.Name)
	if err != nil {
		if resp == nil {
			return workflow.Terminate(workflow.Internal, err.Error())
		}

		if resp.StatusCode != http.StatusNotFound {
			return workflow.Terminate(workflow.DataFederationNotCreatedInAtlas, err.Error())
		}

		_, _, err = ctx.Client.DataFederation.Create(ctx.Context, projectID, dataFederationToAtlas)
		if err != nil {
			return workflow.Terminate(workflow.DataFederationNotCreatedInAtlas, err.Error())
		}

		return workflow.InProgress(workflow.DataFederationCreating, "Data Federation is being created")
	}

	dfFromAtlas, err := DataFederationFromAtlas(atlasSpec)
	if err != nil {
		return workflow.Terminate(workflow.Internal, "can not convert DataFederation (atlas -> operator)")
	}

	if areEqual, _ := dataFederationEqual(*dfFromAtlas, *operatorSpec, log); areEqual {
		return workflow.OK()
	}

	_, _, err = ctx.Client.DataFederation.Update(ctx.Context, projectID, dataFederation.Spec.Name, dataFederationToAtlas, nil)
	if err != nil {
		return workflow.Terminate(workflow.DataFederationNotUpdatedInAtlas, err.Error())
	}

	return workflow.InProgress(workflow.DataFederationUpdating, "Data Federation is being updated")
}

func DataFederationFromAtlas(atlasDF *mongodbatlas.DataFederationInstance) (*akov2.DataFederationSpec, error) {
	dfSpec := &akov2.DataFederationSpec{}
	err := compat.JSONCopy(dfSpec, atlasDF)
	return dfSpec, err
}

func dataFederationEqual(atlasSpec, operatorSpec akov2.DataFederationSpec, log *zap.SugaredLogger) (areEqual bool, diff string) {
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

func getMergedSpec(atlasSpec, operatorSpec akov2.DataFederationSpec) (akov2.DataFederationSpec, error) {
	mergedSpec := akov2.DataFederationSpec{}

	operatorSpec.PrivateEndpoints = []akov2.DataFederationPE{}

	if err := compat.JSONCopy(&mergedSpec, atlasSpec); err != nil {
		return mergedSpec, err
	}
	if err := compat.JSONCopy(&mergedSpec, operatorSpec); err != nil {
		return mergedSpec, err
	}

	mergedSpec.Project = common.ResourceRefNamespaced{}

	return mergedSpec, nil
}

func dataFederationMatchesSpec(log *zap.SugaredLogger, atlasSpec *mongodbatlas.DataFederationInstance, operatorSpec *akov2.AtlasDataFederation) (bool, error) {
	newAtlasSpec, err := DataFederationFromAtlas(atlasSpec)
	if err != nil {
		return false, err
	}

	equal, diff := dataFederationEqual(*newAtlasSpec, operatorSpec.Spec, log)
	if !equal {
		log.Debugf("DataFederation differs from spec: %s", diff)
	}

	return equal, nil
}
