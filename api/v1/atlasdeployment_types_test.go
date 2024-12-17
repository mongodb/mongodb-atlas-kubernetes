package v1 // nolint: dupl

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
)

func TestDeploymentProjectReference(t *testing.T) {
	tests := celTestCase{
		"no project reference is set": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef"},
		},
		"both project references are set": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
						ExternalProjectRef: &ExternalProjectReference{
							ID: "my-project-id",
						},
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef",
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"external project references is set": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ExternalProjectRef: &ExternalProjectReference{
							ID: "my-project-id",
						},
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"kubernetes project references is set": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
				},
			},
		},
	}

	assertCELValidation(t, "../../config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml", tests)
}

func TestDeploymentExternalProjectReferenceConnectionSecret(t *testing.T) {
	tests := celTestCase{
		"external project references is set without connection secret": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ExternalProjectRef: &ExternalProjectReference{
							ID: "my-project-id",
						},
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"external project references is set with connection secret": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ExternalProjectRef: &ExternalProjectReference{
							ID: "my-project-id",
						},
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-dbuser-connection-secret",
						},
					},
				},
			},
		},
		"kubernetes project references is set without connection secret": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
					},
				},
			},
		},
		"kubernetes project references is set with connection secret": {
			object: &AtlasDeployment{
				Spec: AtlasDeploymentSpec{
					ProjectDualReference: ProjectDualReference{
						ProjectRef: &common.ResourceRefNamespaced{
							Name: "my-project",
						},
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-dbuser-connection-secret",
						},
					},
				},
			},
		},
	}

	assertCELValidation(t, "../../config/crd/bases/atlas.mongodb.com_atlasdeployments.yaml", tests)
}
