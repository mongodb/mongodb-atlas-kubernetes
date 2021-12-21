package kube

import (
	"encoding/json"

	v1 "github.com/mongodb/mongodb-atlas-kubernetes/pkg/api/v1"
	kubecli "github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/cli/kubecli"
	"github.com/mongodb/mongodb-atlas-kubernetes/test/e2e/model"
)

func GetProjectResource(data *model.TestDataProvider) (v1.AtlasProject, error) {
	rawData := kubecli.GetProjectResource(data.Resources.Namespace, data.Resources.GetAtlasProjectFullKubeName())
	var project v1.AtlasProject
	err := json.Unmarshal(rawData, &project)
	if err != nil {
		return v1.AtlasProject{}, err
	}
	// ExpectWithOffset(1, json.Unmarshal(rawData, &project)).ShouldNot(HaveOccurred())
	return project, nil
}
