package atlasproject

import (
	"context"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/toptr"
)

func ensureAuditing(ctx *workflow.Context, projectID string, project *v1.AtlasProject) workflow.Result {
	result := createOrDeleteAuditing(ctx, projectID, project)
	if !result.IsOk() {
		ctx.SetConditionFromResult(status.AuditingReadyType, result)
		return result
	}

	if isAuditingEmpty(project.Spec.Auditing) {
		ctx.UnsetCondition(status.AuditingReadyType)
		return workflow.OK()
	}

	ctx.SetConditionTrue(status.AuditingReadyType)
	return workflow.OK()
}

func createOrDeleteAuditing(ctx *workflow.Context, projectID string, project *v1.AtlasProject) workflow.Result {
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

func prepareAuditingSpec(spec *v1.Auditing) *mongodbatlas.Auditing {
	if isAuditingEmpty(spec) {
		return &mongodbatlas.Auditing{
			Enabled: toptr.MakePtr(false),
		}
	}

	return spec.ToAtlas()
}

func auditingInSync(atlas *mongodbatlas.Auditing, spec *v1.Auditing) bool {
	if isAuditingEmpty(atlas) && isAuditingEmpty(spec) {
		return true
	}

	if isAuditingEmpty(atlas) || isAuditingEmpty(spec) {
		return false
	}

	specAsAtlas := spec.ToAtlas()
	removeConfigurationType(atlas)
	return reflect.DeepEqual(atlas, specAsAtlas)
}

func isAuditingEmpty[Auditing mongodbatlas.Auditing | v1.Auditing](auditing *Auditing) bool {
	return auditing == nil
}

func removeConfigurationType(atlas *mongodbatlas.Auditing) {
	atlas.ConfigurationType = ""
}

func fetchAuditing(ctx *workflow.Context, projectID string) (*mongodbatlas.Auditing, error) {
	res, _, err := ctx.Client.Auditing.Get(context.Background(), projectID)
	return res, err
}

func patchAuditing(ctx *workflow.Context, projectID string, auditing *mongodbatlas.Auditing) error {
	_, _, err := ctx.Client.Auditing.Configure(context.Background(), projectID, auditing)
	return err
}
