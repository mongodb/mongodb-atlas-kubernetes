package data

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e/utils"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/project"
)

const ProjectName = "my-project"

func DefaultProject() *akov2.AtlasProject {
	return &akov2.AtlasProject{
		ObjectMeta: metav1.ObjectMeta{
			Name: ProjectName,
		},
		Spec: akov2.AtlasProjectSpec{
			Name: utils.RandomName("Test Atlas Operator Project"),
			ProjectIPAccessList: []project.IPAccessList{
				{
					CIDRBlock: "0.0.0.0/1",
					Comment:   "Everyone has access. For the test purpose only.",
				},
				{
					CIDRBlock: "128.0.0.0/1",
					Comment:   "Everyone has access. For the test purpose only.",
				},
			},
		},
	}
}
