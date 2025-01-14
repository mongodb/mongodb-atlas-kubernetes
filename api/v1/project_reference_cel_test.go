package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

type CELTestObjectGenFunc func(*ProjectDualReference) AtlasCustomResource

type celValidationTestCase struct {
	title          string
	obj, old       AtlasCustomResource
	expectedErrors []string
}

func launchProjectRefCELTests(t *testing.T, genObjFunc CELTestObjectGenFunc, crdPath string) {
	testCases := testCasesFor(genObjFunc)
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.obj)
			require.NoError(t, err)

			unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&tc.old)
			require.NoError(t, err)

			validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
			assert.NoError(t, err)
			errs := validator(unstructuredObject, unstructuredOldObject)

			require.Equal(t, len(tc.expectedErrors), len(errs))

			for i, err := range errs {
				assert.Equal(t, tc.expectedErrors[i], err.Error())
			}
		})
	}
}

func testCasesFor(genObjFunc CELTestObjectGenFunc) []celValidationTestCase {
	return []celValidationTestCase{
		{
			title:          "no project reference is set",
			obj:            genObjFunc(nil),
			expectedErrors: []string{"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef"},
		},
		{
			title: "both project references are set",
			obj: genObjFunc(&ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
				ExternalProjectRef: &ExternalProjectReference{
					ID: "my-project-id",
				},
			}),
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef",
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		{
			title: "external project references is set",
			obj: genObjFunc(&ProjectDualReference{
				ExternalProjectRef: &ExternalProjectReference{
					ID: "my-project-id",
				},
			}),
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		{
			title: "kubernetes project references is set",
			obj: genObjFunc(&ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
			}),
		},
		{
			title: "external project references is set without connection secret",
			obj: genObjFunc(&ProjectDualReference{
				ExternalProjectRef: &ExternalProjectReference{
					ID: "my-project-id",
				},
			}),
			expectedErrors: []string{
				"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
			},
		},
		{
			title: "external project references is set with connection secret",
			obj: genObjFunc(&ProjectDualReference{
				ExternalProjectRef: &ExternalProjectReference{
					ID: "my-project-id",
				},
				ConnectionSecret: &api.LocalObjectReference{
					Name: "my-dbuser-connection-secret",
				},
			}),
		},
		{
			title: "kubernetes project references is set without connection secret",
			obj: genObjFunc(&ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
			}),
		},
		{
			title: "kubernetes project references is set with connection secret",
			obj: genObjFunc(&ProjectDualReference{
				ProjectRef: &common.ResourceRefNamespaced{
					Name: "my-project",
				},
				ConnectionSecret: &api.LocalObjectReference{
					Name: "my-dbuser-connection-secret",
				},
			}),
		},
	}
}

func setDualRef(target, source *ProjectDualReference) {
	target.ConnectionSecret = source.ConnectionSecret
	target.ExternalProjectRef = source.ExternalProjectRef
	target.ProjectRef = source.ProjectRef
}
