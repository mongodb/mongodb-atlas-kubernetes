package atlasdatabaseuser

import (
	"fmt"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
)

func (r *AtlasDatabaseUserReconciler) ensureDatabaseUser(ctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject, dbUser mdbv1.AtlasDatabaseUser) workflow.Result {
	fmt.Printf("%v %v %v %v", ctx, connection, project, dbUser)
	return workflow.OK()
}
