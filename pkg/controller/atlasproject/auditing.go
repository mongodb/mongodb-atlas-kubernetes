package atlasproject

import (
	"encoding/json"
	"fmt"
	"reflect"

	"go.mongodb.org/atlas/mongodbatlas"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/pointer"
	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/status"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/customresource"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/controller/workflow"
)

func ensureAuditing(workflowCtx *workflow.Context, project *v1.AtlasProject, protected bool) workflow.Result {
	canReconcile, err := canAuditingReconcile(workflowCtx, protected, project)
	if err != nil {
		result := workflow.Terminate(workflow.Internal, fmt.Sprintf("unable to resolve ownership for deletion protection: %s", err))
		workflowCtx.SetConditionFromResult(status.AuditingReadyType, result)

		return result
	}

	if !canReconcile {
		result := workflow.Terminate(
			workflow.AtlasDeletionProtection,
			"unable to reconcile Auditing due to deletion protection being enabled. see https://dochub.mongodb.org/core/ako-deletion-protection for further information",
		)
		workflowCtx.SetConditionFromResult(status.AuditingReadyType, result)

		return result
	}

	result := createOrDeleteAuditing(workflowCtx, project.ID(), project)
	if !result.IsOk() {
		workflowCtx.SetConditionFromResult(status.AuditingReadyType, result)
		return result
	}

	if isAuditingEmpty(project.Spec.Auditing) {
		workflowCtx.UnsetCondition(status.AuditingReadyType)
		return workflow.OK()
	}

	workflowCtx.SetConditionTrue(status.AuditingReadyType)
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
			Enabled: pointer.MakePtr(false),
		}
	}

	return spec.ToAtlas()
}

func auditingInSync(atlas *mongodbatlas.Auditing, spec *v1.Auditing) bool {
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

func isAuditingEmpty[Auditing mongodbatlas.Auditing | v1.Auditing](auditing *Auditing) bool {
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

func canAuditingReconcile(workflowCtx *workflow.Context, protected bool, akoProject *v1.AtlasProject) (bool, error) {
	if !protected {
		return true, nil
	}

	latestConfig := &v1.AtlasProjectSpec{}
	latestConfigString, ok := akoProject.Annotations[customresource.AnnotationLastAppliedConfiguration]
	if ok {
		if err := json.Unmarshal([]byte(latestConfigString), latestConfig); err != nil {
			return false, err
		}
	}

	auditing, _, err := workflowCtx.Client.Auditing.Get(workflowCtx.Context, akoProject.ID())
	if err != nil {
		return false, err
	}

	if isAuditingEmpty(auditing) {
		return true, nil
	}

	return auditingInSync(auditing, latestConfig.Auditing) ||
		auditingInSync(auditing, akoProject.Spec.Auditing), nil
}
