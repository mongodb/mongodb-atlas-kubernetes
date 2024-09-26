package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

func TestProjectReference(t *testing.T) {
	validator, err := cel.VersionValidatorFromFile(t, "../../../config/crd/bases/atlas.mongodb.com_atlasdatabaseusers.yaml", "v1")
	require.NoError(t, err)

	tests := map[string]struct {
		dbUser         *AtlasDatabaseUser
		expectedErrors []string
	}{
		"no project reference is set": {
			dbUser: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{},
			},
			expectedErrors: []string{"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef"},
		},
		"both project references are set": {
			dbUser: &AtlasDatabaseUser{
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
			dbUser: &AtlasDatabaseUser{
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
			dbUser: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			unstructuredDBUser, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.dbUser)
			require.NoError(t, err)

			errs := validator(unstructuredDBUser, nil)

			require.Equal(t, len(tt.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tt.expectedErrors[i], err.Error())
			}
		})
	}
}

func TestExternalProjectReferenceConnectionSecret(t *testing.T) {
	validator, err := cel.VersionValidatorFromFile(t, "../../../config/crd/bases/atlas.mongodb.com_atlasdatabaseusers.yaml", "v1")
	require.NoError(t, err)

	tests := map[string]struct {
		dbUser         *AtlasDatabaseUser
		expectedErrors []string
	}{
		"external project references is set without connection secret": {
			dbUser: &AtlasDatabaseUser{
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
			dbUser: &AtlasDatabaseUser{
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
			dbUser: &AtlasDatabaseUser{
				Spec: AtlasDatabaseUserSpec{
					Project: &common.ResourceRefNamespaced{
						Name: "my-project",
					},
				},
			},
		},
		"kubernetes project references is set with connection secret": {
			dbUser: &AtlasDatabaseUser{
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			unstructuredDBUser, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tt.dbUser)
			require.NoError(t, err)

			errs := validator(unstructuredDBUser, nil)

			require.Equal(t, len(tt.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tt.expectedErrors[i], err.Error())
			}
		})
	}
}
