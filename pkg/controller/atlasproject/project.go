package atlasproject

import (
	"context"

	mdbv1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/atlas"
	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/controller/workflow"
	"go.mongodb.org/atlas/mongodbatlas"
)

// ensureProjectExists creates the project if it doesn't exist yet. Returns the project ID
func ensureProjectExists(ctx *workflow.Context, connection atlas.Connection, project *mdbv1.AtlasProject) (string, workflow.Result) {
	client, err := atlas.Client(connection, ctx.Log)
	if err != nil {
		return "", workflow.Terminate(workflow.Internal, err.Error())
	}
	// Try to find the project
	projectID, err := findProject(connection, project, client)
	if err != nil {
		return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
	}
	if projectID != "" {
		ctx.Log.Debugw("Found Atlas Project", "id", projectID)
		return projectID, workflow.OK()
	}

	// Otherwise try to create it
	p := &mongodbatlas.Project{
		OrgID: connection.OrgID,
		Name:  project.Spec.Name,
	}

	if p, _, err = client.Projects.Create(context.Background(), p); err != nil {
		return "", workflow.Terminate(workflow.ProjectNotCreatedInAtlas, err.Error())
	}
	ctx.Log.Infow("Created Atlas Project", "name", project.Spec.Name, "id", p.ID)

	return p.ID, workflow.OK()
}

func findProject(connection atlas.Connection, project *mdbv1.AtlasProject, client *mongodbatlas.Client) (string, error) {
	var projectID string
	err := atlas.TraversePages(func(pageNum int) (atlas.Paginated, error) {
		return getProjectsForOrganizations(client, connection.OrgID, pageNum)
	}, func(entity interface{}) bool {
		p := entity.(*mongodbatlas.Project)
		if p.Name == project.Spec.Name {
			projectID = p.ID
			return true
		}
		return false
	})
	return projectID, err
}

func getProjectsForOrganizations(client *mongodbatlas.Client, orgID string, pageNum int) (atlas.Paginated, error) {
	// TODO test if the project level API key allows to find the project
	projects, response, err := client.Organizations.Projects(context.Background(), orgID, atlas.DefaultListOptions(pageNum))
	if err != nil {
		return nil, err
	}
	return atlas.NewAtlasPaginated(response, projects.Results), err
}
