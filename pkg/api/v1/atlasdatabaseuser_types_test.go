package v1 // nolint: dupl

import (
	"testing"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
)

func TestProjectReference(t *testing.T) {
	tests := celTestCase{
		"no project reference is set": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef"},
		},
		"both project references are set": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
					ExternalProjectRef: &ExternalProjectReference{
						ID: "my-project-id",
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef",
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"external project references is set": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					ExternalProjectRef: &ExternalProjectReference{
						ID: "my-project-id",
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"kubernetes project references is set": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
				},
			},
		},
	}

	assertCELValidation(t, "../../../config/crd/bases/atlas.mongodb.com_atlasdatabaseusers.yaml", tests)
}

func TestExternalProjectReferenceConnectionSecret(t *testing.T) {
	tests := celTestCase{
		"external project references is set without connection secret": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					ExternalProjectRef: &ExternalProjectReference{
						ID: "my-project-id",
					},
				},
			},
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		"external project references is set with connection secret": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					ExternalProjectRef: &ExternalProjectReference{
						ID: "my-project-id",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-dbuser-connection-secret",
						},
					},
				},
			},
		},
		"kubernetes project references is set without connection secret": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
				},
			},
		},
		"kubernetes project references is set with connection secret": {
			object: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
					LocalCredentialHolder: api.LocalCredentialHolder{
						ConnectionSecret: &api.LocalObjectReference{
							Name: "my-dbuser-connection-secret",
						},
					},
				},
			},
		},
	}

	assertCELValidation(t, "../../../config/crd/bases/atlas.mongodb.com_atlasdatabaseusers.yaml", tests)
}
