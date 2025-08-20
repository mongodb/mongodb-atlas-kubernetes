package plugins

import (
	"fmt"
	"strings"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

type CrdPlugin struct{}

func (cp *CrdPlugin) Name() string {
	return "crd"
}

func (cp *CrdPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.CRDInput)

	if !ok {
		return nil // no operation performed
	}

	crd := i.CRD
	crdConfig := i.CRDConfig

	pluralGvk, singularGvk := guessKindToResource(crdConfig.GVK)

	crd.ObjectMeta = v1.ObjectMeta{
		Name: fmt.Sprintf("%s.%s", pluralGvk.Resource, pluralGvk.Group),
	}

	crd.Spec = apiextensions.CustomResourceDefinitionSpec{
		Group: pluralGvk.Group,
		Scope: apiextensions.NamespaceScoped,
		Names: apiextensions.CustomResourceDefinitionNames{
			Kind:     crdConfig.GVK.Kind,
			ListKind: fmt.Sprintf("%sList", crdConfig.GVK.Kind),
			Plural:   pluralGvk.Resource,
			Singular: singularGvk.Resource,
		},
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    crdConfig.GVK.Version,
				Served:  true,
				Storage: true,
			},
		},
		PreserveUnknownFields: ptr.To(false),
		Validation: &apiextensions.CustomResourceValidation{
			OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
				Type:        "object",
				Description: fmt.Sprintf("A %v, managed by the MongoDB Kubernetes Atlas Operator.", singularGvk.Resource),
				Properties: map[string]apiextensions.JSONSchemaProps{
					"spec": {
						Type: "object",
						Description: fmt.Sprintf(`Specification of the %v supporting the following versions:

%v

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource, strings.Join(majorVersions(crdConfig), "\n")),
						Properties: map[string]apiextensions.JSONSchemaProps{},
					},
					"status": {
						Type:        "object",
						Description: fmt.Sprintf(`Most recently observed read-only status of the %v for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource),
						Properties:  map[string]apiextensions.JSONSchemaProps{},
					},
				},
			},
		},
	}

	crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties["conditions"] = apiextensions.JSONSchemaProps{
		Type:        "array",
		Description: "Represents the latest available observations of a resource's current state.",
		Items: &apiextensions.JSONSchemaPropsOrArray{
			Schema: &apiextensions.JSONSchemaProps{
				Type:     "object",
				Required: []string{"type", "status"},
				Properties: map[string]apiextensions.JSONSchemaProps{
					"type":               {Type: "string", Description: "Type of condition."},
					"status":             {Type: "string", Description: "Status of the condition, one of True, False, Unknown."},
					"observedGeneration": {Type: "integer", Description: "observedGeneration represents the .metadata.generation that the condition was set based upon."},
					"message":            {Type: "string", Description: "A human readable message indicating details about the transition."},
					"reason":             {Type: "string", Description: "The reason for the condition's last transition."},
					"lastTransitionTime": {Type: "string", Format: "date-time", Description: "Last time the condition transitioned from one status to another."},
				},
			},
		},
		XListMapKeys: []string{
			"type",
		},
		XListType: ptr.To("map"),
	}

	crd.Status.StoredVersions = []string{}

	// enable status subresource
	crd.Spec.Subresources = &apiextensions.CustomResourceSubresources{
		Status: &apiextensions.CustomResourceSubresourceStatus{},
	}

	crd.Spec.Names.Categories = crdConfig.Categories
	crd.Spec.Names.ShortNames = crdConfig.ShortNames

	for _, version := range crd.Spec.Versions {
		if version.Storage {
			crd.Status.StoredVersions = append(crd.Status.StoredVersions, version.Name)
		}
	}

	return nil
}

func NewCrdPlugin() *CrdPlugin {
	return &CrdPlugin{}
}

func majorVersions(crdConfig *configv1alpha1.CRDConfig) []string {
	var result []string
	for _, m := range crdConfig.Mappings {
		result = append(result, "- "+m.MajorVersion)
	}
	return result
}

func guessKindToResource(gvk v1.GroupVersionKind) ( /*plural*/ runtimeschema.GroupVersionResource /*singular*/, runtimeschema.GroupVersionResource) {
	runtimeGVK := runtimeschema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
	kindName := runtimeGVK.Kind
	if len(kindName) == 0 {
		return runtimeschema.GroupVersionResource{}, runtimeschema.GroupVersionResource{}
	}
	singularName := strings.ToLower(kindName)
	singular := runtimeGVK.GroupVersion().WithResource(singularName)

	switch string(singularName[len(singularName)-1]) {
	case "s":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "x":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "y":
		return runtimeGVK.GroupVersion().WithResource(strings.TrimSuffix(singularName, "y") + "ies"), singular
	}

	return runtimeGVK.GroupVersion().WithResource(singularName + "s"), singular
}
