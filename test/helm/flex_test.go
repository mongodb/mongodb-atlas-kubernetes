package helm

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	akov2 "github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cmd"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/decoder"
)

func TestFlexSpec(t *testing.T) {
	output := cmd.RunCommand(t, "helm", "template", "--values=flex_values.yaml", "../../helm-charts/atlas-deployment")
	objects := decoder.DecodeAll(t, output)

	var gotDeployment *akov2.AtlasDeployment
	for _, obj := range objects {
		if d, ok := obj.(*akov2.AtlasDeployment); ok {
			if gotDeployment != nil {
				t.Errorf("Expect one deployment only but also got: %v", d)
				continue
			}
			gotDeployment = d
		}
	}

	// ignore
	gotDeployment.Kind = ""
	gotDeployment.APIVersion = ""
	gotDeployment.Labels = nil

	wantDeployment := &akov2.AtlasDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flex-instance",
			Namespace: "default",
		},
		Spec: akov2.AtlasDeploymentSpec{
			FlexSpec: &akov2.FlexSpec{
				Name: "flex-instance",
				ProviderSettings: &akov2.FlexProviderSettings{
					BackingProviderName: "AWS",
					RegionName:          "US_EAST_1",
				},
			},
			ProjectDualReference: akov2.ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "release-name-my-project",
				},
				ExternalProjectRef: nil,
				ConnectionSecret:   nil,
			},
		},
	}

	require.Equal(t, wantDeployment, gotDeployment)
}
