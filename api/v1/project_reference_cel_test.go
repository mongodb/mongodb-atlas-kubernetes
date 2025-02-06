package v1

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/api/v1/common"
	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/cel"
)

const (
	crdBasePath = "../../config/crd/bases/"
)

type projectReferrerObject interface {
	runtime.Object
	ProjectDualRef() *ProjectDualReference
}

var dualRefCRDs = []struct {
	obj      projectReferrerObject
	filename string
}{
	{
		obj:      &AtlasDatabaseUser{},
		filename: "atlas.mongodb.com_atlasdatabaseusers.yaml",
	},
	{
		obj:      &AtlasDeployment{},
		filename: "atlas.mongodb.com_atlasdeployments.yaml",
	},
	{
		obj:      &AtlasIPAccessList{},
		filename: "atlas.mongodb.com_atlasipaccesslists.yaml",
	},
	{
		obj:      &AtlasPrivateEndpoint{},
		filename: "atlas.mongodb.com_atlasprivateendpoints.yaml",
	},
	{
		obj:      &AtlasIPAccessList{},
		filename: "atlas.mongodb.com_atlasipaccesslists.yaml",
	},
	{
		obj: &AtlasNetworkContainer{
			Spec: AtlasNetworkContainerSpec{
				Provider: "GCP", // Avoid triggering container specific validations
			},
		},
		filename: "atlas.mongodb.com_atlasnetworkcontainers.yaml",
	},
	{
		obj: &AtlasNetworkPeering{
			Spec: AtlasNetworkPeeringSpec{ // Avoid triggering peering specific validations
				ContainerRef: ContainerDualReference{Name: "fake-ref"},
			},
		},
		filename: "atlas.mongodb.com_atlasnetworkpeerings.yaml",
	},
}

var testCases = []struct {
	title          string
	ref            *ProjectDualReference
	expectedErrors []string
}{
	{
		title:          "no project reference is set",
		ref:            &ProjectDualReference{},
		expectedErrors: []string{"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef"},
	},
	{
		title: "both project references are set",
		ref: &ProjectDualReference{
			ProjectRef: &common.ResourceRefNamespaced{
				Name: "my-project",
			},
			ExternalProjectRef: &ExternalProjectReference{
				ID: "my-project-id",
			},
		},
		expectedErrors: []string{
			"spec: Invalid value: \"object\": must define only one project reference through externalProjectRef or projectRef",
			"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
		},
	},
	{
		title: "external project references is set",
		ref: &ProjectDualReference{
			ExternalProjectRef: &ExternalProjectReference{
				ID: "my-project-id",
			},
		},
		expectedErrors: []string{
			"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
		},
	},
	{
		title: "kubernetes project references is set",
		ref: &ProjectDualReference{
			ProjectRef: &common.ResourceRefNamespaced{
				Name: "my-project",
			},
		},
	},
	{
		title: "external project references is set without connection secret",
		ref: &ProjectDualReference{
			ExternalProjectRef: &ExternalProjectReference{
				ID: "my-project-id",
			},
		},
		expectedErrors: []string{
			"spec: Invalid value: \"object\": must define a local connection secret when referencing an external project",
		},
	},
	{
		title: "external project references is set with connection secret",
		ref: &ProjectDualReference{
			ExternalProjectRef: &ExternalProjectReference{
				ID: "my-project-id",
			},
			ConnectionSecret: &api.LocalObjectReference{
				Name: "my-dbuser-connection-secret",
			},
		},
	},
	{
		title: "kubernetes project references is set without connection secret",
		ref: &ProjectDualReference{
			ProjectRef: &common.ResourceRefNamespaced{
				Name: "my-project",
			},
		},
	},
	{
		title: "kubernetes project references is set with connection secret",
		ref: &ProjectDualReference{
			ProjectRef: &common.ResourceRefNamespaced{
				Name: "my-project",
			},
			ConnectionSecret: &api.LocalObjectReference{
				Name: "my-dbuser-connection-secret",
			},
		},
	},
}

func TestProjectDualReferenceCELValidations(t *testing.T) {
	for _, dualRef := range dualRefCRDs {
		for _, tc := range testCases {
			title := fmt.Sprintf("%T %s", dualRef.obj, tc.title)
			obj := dualRef.obj.DeepCopyObject()
			pdr, ok := obj.(projectReferrerObject)
			require.True(t, ok)
			setDualRef(pdr.ProjectDualRef(), tc.ref)
			old := dualRef.obj.DeepCopyObject()
			crdPath := filepath.Join(crdBasePath, dualRef.filename)
			t.Run(title, func(t *testing.T) {
				unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&obj)
				require.NoError(t, err)

				unstructuredOldObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&old)
				require.NoError(t, err)

				validator, err := cel.VersionValidatorFromFile(t, crdPath, "v1")
				assert.NoError(t, err)
				errs := validator(unstructuredObject, unstructuredOldObject)

				for i, err := range errs {
					fmt.Printf("%s error %d: %v\n", title, i, err)
				}

				require.Equal(t, len(tc.expectedErrors), len(errs))
				for i, err := range errs {
					assert.Equal(t, tc.expectedErrors[i], err.Error())
				}
			})
		}
	}
}

func setDualRef(target, source *ProjectDualReference) {
	target.ConnectionSecret = source.ConnectionSecret
	target.ExternalProjectRef = source.ExternalProjectRef
	target.ProjectRef = source.ProjectRef
}
