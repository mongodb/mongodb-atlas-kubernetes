/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/config"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

type Generator struct {
	config      v1alpha1.CRDConfig
	definitions map[string]v1alpha1.OpenAPIDefinition

	// added during schemaPropsToJSONProps
	sensitiveFieldsDocs []string
}

func NewGenerator(crdConfig v1alpha1.CRDConfig, definitions []v1alpha1.OpenAPIDefinition) *Generator {
	definitionsMap := map[string]v1alpha1.OpenAPIDefinition{}
	for _, def := range definitions {
		definitionsMap[def.Name] = def
	}
	return &Generator{
		config:      crdConfig,
		definitions: definitionsMap,
	}
}

func (g *Generator) majorVersions() []string {
	var result []string
	for _, m := range g.config.Mappings {
		result = append(result, "- "+m.MajorVersion)
	}
	return result
}

func (g *Generator) Generate(ctx context.Context) (*apiextensions.CustomResourceDefinition, error) {
	pluralGvk, singularGvk := guessKindToResource(g.config.GVK)

	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", pluralGvk.Resource, pluralGvk.Group),
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: pluralGvk.Group,
			Scope: apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     g.config.GVK.Kind,
				ListKind: fmt.Sprintf("%sList", g.config.GVK.Kind),
				Plural:   pluralGvk.Resource,
				Singular: singularGvk.Resource,
			},
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    g.config.GVK.Version,
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

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource, strings.Join(g.majorVersions(), "\n")),
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
		},
	}

	for _, mapping := range g.config.Mappings {
		def, ok := g.definitions[mapping.OpenAPIRef.Name]
		if !ok {
			return nil, fmt.Errorf("no OpenAPI definition named %q found", mapping.OpenAPIRef.Name)
		}

		openApiSpec, err := config.LoadOpenAPI(def.Path)
		if err != nil {
			return nil, fmt.Errorf("error loading spec: %w", err)
		}

		err = g.generateProps(openApiSpec, crd, &mapping)
		if err != nil {
			return nil, fmt.Errorf("error generating props: %w", err)
		}
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

	// TODO: yaml.Marshal creates an empty status field that we should remove
	// StoredVersions is set to empty array instead of nil to bypass the following issue:
	// https://github.com/fybrik/openapi2crd/issues/1
	crd.Status.StoredVersions = []string{}

	// enable status subresource
	crd.Spec.Subresources = &apiextensions.CustomResourceSubresources{
		Status: &apiextensions.CustomResourceSubresourceStatus{},
	}

	crd.Spec.Names.Categories = g.config.Categories
	crd.Spec.Names.ShortNames = g.config.ShortNames

	for _, version := range crd.Spec.Versions {
		if version.Storage {
			crd.Status.StoredVersions = append(crd.Status.StoredVersions, version.Name)
		}
	}

	if err := ValidateCRD(ctx, crd); err != nil {
		log.Printf("Error validating CRD: %v", err)
	}

	return crd, nil
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

func (g *Generator) generateProps(openApiSpec *openapi3.T, crd *apiextensions.CustomResourceDefinition, mapping *v1alpha1.CRDMapping) error {
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", crd.Spec.Names.Singular, mapping.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	majorVersionSpec := crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion]

	if mapping.ParametersMapping.FieldPath.Name != "" {
		var operation *openapi3.Operation

		pathItem, ok := openApiSpec.Paths[mapping.ParametersMapping.FieldPath.Name]
		if !ok {
			return fmt.Errorf("OpenAPI path %q does not exist", mapping.ParametersMapping)
		}

		switch mapping.ParametersMapping.FieldPath.Verb {
		case "post":
			operation = pathItem.Post
		case "put":
			operation = pathItem.Put
		default:
			return fmt.Errorf("verb %q unsupported", mapping.ParametersMapping.FieldPath.Verb)
		}

		for _, p := range operation.Parameters {
			switch p.Value.Name {
			case "includeCount":
			case "itemsPerPage":
			case "pageNum":
			case "envelope":
			case "pretty":
			default:
				props := g.schemaPropsToJSONProps(p.Value.Schema, nil)
				props.Description = p.Value.Description
				props.XValidations = apiextensions.ValidationRules{
					{
						Rule:    "self == oldSelf",
						Message: fmt.Sprintf("%s cannot be modified after creation", p.Value.Name),
					},
				}
				majorVersionSpec.Properties[p.Value.Name] = *props
				majorVersionSpec.Required = append(majorVersionSpec.Required, p.Value.Name)
			}
		}
	}

	var entrySchemaRef *openapi3.SchemaRef

	if mapping.EntryMapping.Schema != "" {
		var ok bool
		entrySchemaRef, ok = openApiSpec.Components.Schemas[mapping.EntryMapping.Schema]
		if !ok {
			return fmt.Errorf("entry schema %q not found in openapi spec", mapping.EntryMapping.Schema)
		}
	}

	if mapping.EntryMapping.Path.Name != "" {
		entrySchemaRef = openApiSpec.Paths[mapping.EntryMapping.Path.Name].Operations()[strings.ToUpper(mapping.EntryMapping.Path.Verb)].RequestBody.Value.Content[mapping.EntryMapping.Path.RequestBody.MimeType].Schema
	}

	entryProps := g.schemaPropsToJSONProps(entrySchemaRef, &mapping.EntryMapping)
	entryProps.Description = fmt.Sprintf("The entry fields of the %v resource spec. These fields can be set for creating and updating %v.", crd.Spec.Names.Singular, crd.Spec.Names.Plural)
	majorVersionSpec.Properties["entry"] = *entryProps

	if mapping.StatusMapping.Schema != "" {
		statusSchemaRef, ok := openApiSpec.Components.Schemas[mapping.StatusMapping.Schema]
		if !ok {
			return fmt.Errorf("status schema %q not found in openapi spec", mapping.StatusMapping.Schema)
		}

		statusProps := g.schemaPropsToJSONProps(statusSchemaRef, &mapping.StatusMapping)
		statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", crd.Spec.Names.Singular, mapping.MajorVersion)
		if statusProps != nil {
			crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[mapping.MajorVersion] = *statusProps
		}
	}

	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = majorVersionSpec
	return nil
}
