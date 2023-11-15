package kube

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/e2e/model"
)

func GetProjectResource(data *model.TestDataProvider) (v1.AtlasProject, error) {
	project := v1.AtlasProject{}
	err := data.K8SClient.Get(data.Context, client.ObjectKey{Namespace: data.Resources.Namespace,
		Name: data.Resources.Project.ObjectMeta.GetName()}, &project)
	if err != nil {
		return v1.AtlasProject{}, err
	}
	return project, nil
}
