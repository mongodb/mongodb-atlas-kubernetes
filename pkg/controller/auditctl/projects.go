package auditctl

import (
	"context"
	"errors"
	"reflect"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1alpha1/status"
)

func (r *AtlasAuditingReconciler) reconcileProjectsAuditing(ctx context.Context, auditing *v1alpha1.AtlasAuditing) (*v1alpha1.AtlasAuditing, error) {
	var errs error
	for _, projectID := range projectIDs(auditing) {
		var err error
		auditing, err = r.reconcileProjectAuditing(ctx, auditing, projectID)
		errs = errors.Join(errs, err)
	}
	return auditing, errs
}

// reconcileProjectAuditing will try to synchronize the auditing spec on the given project
//
// State machine summary:
//
// UNKNOWN --(project not found)--> MISSING
// UNKNOWN --(failure)--> ERROR
// UNKNOWN --(different)--> UPDATE [--(failure)--> ERROR]
// UNKNOWN --(in sync)--> IDLE
//
//	NOTE: IDLE is the only state that does not update the status
func (r *AtlasAuditingReconciler) reconcileProjectAuditing(ctx context.Context, auditing *v1alpha1.AtlasAuditing, projectID string) (*v1alpha1.AtlasAuditing, error) {
	atlasAuditing, err := r.AuditService.Get(ctx, projectID)
	if err != nil {
		return failure(auditing, projectID, err)
	}
	if reflect.DeepEqual(auditing.Spec, atlasAuditing) {
		return idle(auditing)
	}
	return r.update(ctx, auditing, projectID)
}

func (r *AtlasAuditingReconciler) update(ctx context.Context, auditing *v1alpha1.AtlasAuditing, projectID string) (*v1alpha1.AtlasAuditing, error) {
	resultAuditing := auditing.DeepCopy()
	if err := r.AuditService.Set(ctx, projectID, &auditing.Spec); err != nil {
		resultAuditing.UpdateStatus([]api.Condition{}, status.WithProjectFailure(projectID, err))
		return resultAuditing, err
	}
	resultAuditing.UpdateStatus([]api.Condition{}, status.WithSuccess(projectID))
	return resultAuditing, nil
}

func failure(auditing *v1alpha1.AtlasAuditing, projectID string, err error) (*v1alpha1.AtlasAuditing, error) {
	// TODO: Distinguish MISSING and ERROR?
	resultAuditing := auditing.DeepCopy()
	resultAuditing.UpdateStatus([]api.Condition{}, status.WithProjectFailure(projectID, err))
	return resultAuditing, err
}

func idle(auditing *v1alpha1.AtlasAuditing) (*v1alpha1.AtlasAuditing, error) {
	return auditing, nil
}
