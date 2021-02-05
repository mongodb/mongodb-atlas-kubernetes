package atlasproject

import (
	"fmt"
	"time"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/util/timeutil"
)

func (r *AtlasProjectReconciler) ensureIPAccessList(ctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject) workflow.Result {
	active, expired, result := filterActiveIPAccessLists(project)
	if !result.IsOk() {
		return result
	}
	fmt.Println(active, expired, ctx, connection)
	return workflow.OK()
}

func filterActiveIPAccessLists(project *mdbv1.AtlasProject) ([]mdbv1.ProjectIPAccessList, []mdbv1.ProjectIPAccessList, workflow.Result) {
	active := make([]mdbv1.ProjectIPAccessList, 0)
	expired := make([]mdbv1.ProjectIPAccessList, 0)
	for _, list := range project.Spec.ProjectIPAccessList {
		if list.DeleteAfterDate != "" {
			iso8601, err := timeutil.ParseISO8601(list.DeleteAfterDate)
			if err != nil {
				return active, expired, workflow.Terminate(workflow.ProjectIPAccessBadFormatted, err.Error())
			}
			if iso8601.Before(time.Now()) {
				expired = append(expired, list)
				continue
			}
		}
		// Either 'deleteAfterDate' field is not specified or it's higher than the current time
		active = append(active, list)
	}
	return active, expired, workflow.OK()
}
