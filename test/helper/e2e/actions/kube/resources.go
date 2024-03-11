package kube

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/model"
)

func GetProjectResource(data *model.TestDataProvider) (akov2.AtlasProject, error) {
	project := akov2.AtlasProject{}
	err := data.K8SClient.Get(data.Context, client.ObjectKey{Namespace: data.Resources.Namespace,
		Name: data.Resources.Project.ObjectMeta.GetName()}, &project)
	if err != nil {
		return akov2.AtlasProject{}, err
	}
	return project, nil
}
