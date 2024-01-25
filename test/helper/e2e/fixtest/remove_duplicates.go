package fixtest

import (
	"context"
	"io"
	"sort"

	"go.mongodb.org/atlas-sdk/v20231115004/admin"
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
func EnsureNoDuplicates(client *admin.APIClient, logger *zap.SugaredLogger, projectName string) error {
	found, err := listProjectsByName(client, projectName)
	if err != nil || len(found) <= 1 {
		return err
	}

	logger.Warnf("Found more than one project with name %q", projectName)

	keep, rest := selectProject(found)

	logger.Warnf("Will keep project ID %s as %s and remove the rest %v", keep.GetId(), projectName, ids(rest))

	return removeProjects(client, rest)
}

func listProjectsByName(client *admin.APIClient, projectName string) ([]admin.Group, error) {
	projects, _, err := client.ProjectsApi.ListProjects(context.Background()).Execute()
	if err != nil {
		return nil, err
	}

	found := make([]admin.Group, 0, projects.GetTotalCount())
	for _, project := range *projects.Results {
		if project.Name == projectName {
			found = append(found, project)
		}
	}

	return found, nil
}

func selectProject(projects []admin.Group) (admin.Group, []admin.Group) {
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].GetId() < projects[j].GetId()
	})

	return projects[0], projects[1:]
}

func removeProjects(client *admin.APIClient, projects []admin.Group) error {
	for _, project := range projects {
		_, _, err := client.ProjectsApi.DeleteProject(context.Background(), project.GetId()).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func ids(projects []admin.Group) []string {
	items := make([]string, 0, len(projects))
	for _, prj := range projects {
		items = append(items, prj.GetId())
	}

	return items
}
