package atlasproject

import (
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func ensureAuditing(workflowCtx *workflow.Context, project *akov2.AtlasProject) workflow.Result {
	result := createOrDeleteAuditing(workflowCtx, project.ID(), project)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(api.AuditingReadyType, result)
		return result
	}

	if isAuditingEmpty(project.Spec.Auditing) {
		workflowCtx.UnsetCondition(api.AuditingReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(api.AuditingReadyType)
	return workflow.OK()
}

func createOrDeleteAuditing(ctx *workflow.Context, projectID string, project *akov2.AtlasProject) workflow.Result {
	atlas, err := fetchAuditing(ctx, projectID)
	if err != nil {
		return workflow.Terminate(workflow.ProjectAuditingReady, err.Error())
	}

	if !auditingInSync(atlas, project.Spec.Auditing) {
		err := patchAuditing(ctx, projectID, prepareAuditingSpec(project.Spec.Auditing))
		if err != nil {
			return workflow.Terminate(workflow.ProjectAuditingReady, err.Error())
		}
	}

	return workflow.OK()
}

func prepareAuditingSpec(spec *akov2.Auditing) *mongodbatlas.Auditing {
	if isAuditingEmpty(spec) {
		return &mongodbatlas.Auditing{
			Enabled: pointer.MakePtr(false),
		}
	}

	return spec.ToAtlas()
}

func auditingInSync(atlas *mongodbatlas.Auditing, spec *akov2.Auditing) bool {
	if isAuditingEmpty(atlas) && isAuditingEmpty(spec) {
		return true
	}

	specAsAtlas := &mongodbatlas.Auditing{
		AuditAuthorizationSuccess: pointer.MakePtr(false),
		Enabled:                   pointer.MakePtr(false),
	}

	if !isAuditingEmpty(spec) {
		specAsAtlas = spec.ToAtlas()
	}

	if isAuditingEmpty(atlas) {
		atlas = &mongodbatlas.Auditing{
			AuditAuthorizationSuccess: pointer.MakePtr(false),
			Enabled:                   pointer.MakePtr(false),
		}
	}

	removeConfigurationType(atlas)

	return reflect.DeepEqual(atlas, specAsAtlas)
}

func isAuditingEmpty[Auditing mongodbatlas.Auditing | akov2.Auditing](auditing *Auditing) bool {
	return auditing == nil
}

func removeConfigurationType(atlas *mongodbatlas.Auditing) {
	atlas.ConfigurationType = ""
}

func fetchAuditing(ctx *workflow.Context, projectID string) (*mongodbatlas.Auditing, error) {
	res, _, err := ctx.Client.Auditing.Get(ctx.Context, projectID)
	return res, err
}

func patchAuditing(ctx *workflow.Context, projectID string, auditing *mongodbatlas.Auditing) error {
	_, _, err := ctx.Client.Auditing.Configure(ctx.Context, projectID, auditing)
	return err
}
