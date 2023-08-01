package fixtest

import (
	"context"
	"io"
	"sort"

	"go.mongodb.org/atlas/mongodbatlas"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLoggerFrom builds a zap.SugaredLogger from an IO Writer
func ZapLoggerFrom(w io.Writer) *zap.SugaredLogger {
	zcore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(zapcore.AddSync(w)),
		zap.NewAtomicLevel(),
	)
	return zap.New(zcore).Sugar()
}

// EnsureNoDuplicates removes projects with same name but different ID.
// Atlas sometimes creates duplicate projects, we need our tests to defend
// against that to avoid flaky tests
func EnsureNoDuplicates(client *mongodbatlas.Client, logger *zap.SugaredLogger, projectName string) error {
	found, err := listProjectsByName(client, projectName)
	if err != nil || len(found) <= 1 {
		return err
	}
	logger.Warnf("Found more than one project with name %q", projectName)
	keep, rest := selectProject(found)
	logger.Warnf("Will keep project ID %s as %s and remove the rest %v", keep.ID, projectName, ids(rest))
	return removeProjects(client, rest)
}

func listProjectsByName(client *mongodbatlas.Client, projectName string) ([]*mongodbatlas.Project, error) {
	projects, _, err := client.Projects.GetAllProjects(
		context.Background(),
		&mongodbatlas.ListOptions{},
	)
	found := []*mongodbatlas.Project{}
	if err != nil {
		return found, err
	}
	for _, project := range projects.Results {
		if project.Name == projectName {
			found = append(found, project)
		}
	}
	return found, nil
}

func selectProject(projects []*mongodbatlas.Project) (*mongodbatlas.Project, []*mongodbatlas.Project) {
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].ID < projects[j].ID
	})
	return projects[0], projects[1:]
}

func removeProjects(client *mongodbatlas.Client, projects []*mongodbatlas.Project) error {
	for _, project := range projects {
		_, err := client.Projects.Delete(context.Background(), project.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func ids(projects []*mongodbatlas.Project) []string {
	ids := []string{}
	for _, prj := range projects {
		ids = append(ids, prj.ID)
	}
	return ids
}
