package indexer

import (
	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
)

const (
	AtlasProjectByTeamIndex = "atlasproject.spec.teams"
)

type AtlasProjectByTeamIndexer struct {
	logger *zap.SugaredLogger
}

func NewAtlasProjectByTeamIndexer(logger *zap.Logger) *AtlasProjectByTeamIndexer {
	return &AtlasProjectByTeamIndexer{
		logger: logger.Named(AtlasProjectByTeamIndex).Sugar(),
	}
}

func (*AtlasProjectByTeamIndexer) Object() client.Object {
	return &akov2.AtlasProject{}
}

func (*AtlasProjectByTeamIndexer) Name() string {
	return AtlasProjectByTeamIndex
}

func (a *AtlasProjectByTeamIndexer) Keys(object client.Object) []string {
	project, ok := object.(*akov2.AtlasProject)
	if !ok {
		a.logger.Errorf("expected *akov2.AtlasProject but got %T", object)
		return nil
	}

	result := sets.New[string]()
	for _, team := range project.Spec.Teams {
		if team.TeamRef.IsEmpty() {
			continue
		}
		result.Insert(team.TeamRef.GetObject(project.Namespace).String())
	}

	return result.UnsortedList()
}
