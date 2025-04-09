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
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/config"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

func anyEntry[T any](source map[string]T, defaultValue T) T {
	for _, v := range source {
		return v
	}
	return defaultValue
}

type Generator struct {
	config      v1alpha1.CRDConfig
	definitions map[string]v1alpha1.OpenAPIDefinition
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
			Validation: &apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Type: "object",
					Properties: map[string]apiextensions.JSONSchemaProps{
						"spec": {
							Type:       "object",
							Properties: map[string]apiextensions.JSONSchemaProps{},
						},
						"status": {
							Type:       "object",
							Properties: map[string]apiextensions.JSONSchemaProps{},
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

		spec, err := config.LoadOpenAPI(def.Path)
		if err != nil {
			return nil, fmt.Errorf("error loading spec: %w", err)
		}

		// Fix known types (ref: https://github.com/kubernetes/kubernetes/issues/62329)
		spec.Components.Schemas["k8s.io/apimachinery/pkg/util/intstr.IntOrString"] = openapi3.NewSchemaRef("", &openapi3.Schema{
			AnyOf: openapi3.SchemaRefs{
				{
					Value: openapi3.NewStringSchema(),
				},
				{
					Value: openapi3.NewIntegerSchema(),
				},
			},
		})

		var operation *openapi3.Operation
		switch mapping.Verb {
		case "post":
			pathItem, ok := spec.Paths[mapping.Path]
			if !ok {
				log.Printf("WARNING: OpenAPI path %q does not exist in %q\n", mapping.Path, def.Path)
				continue
			}
			operation = pathItem.Post
		case "put":
			pathItem, ok := spec.Paths[mapping.Path]
			if !ok {
				log.Printf("WARNING: OpenAPI path %q does not exist in %q\n", mapping.Path, def.Path)
				continue
			}
			operation = pathItem.Put
		default:
			return nil, fmt.Errorf("verb %q not supported", mapping.Verb)
		}

		err = g.generateProps(crd, &mapping, operation)
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

func (g *Generator) generateProps(crd *apiextensions.CustomResourceDefinition, mapping *v1alpha1.Mapping, operation *openapi3.Operation) error {
	content := operation.RequestBody.Value.Content
	mediaType := anyEntry(content, nil)
	entrySchemaRef := mediaType.Schema
	entrySchema := FilterSchemaProps("", false, entrySchemaRef, func(key string, schemaRef *openapi3.SchemaRef) bool {
		return !schemaRef.Value.ReadOnly
	})
	entryProps := g.schemaPropsToJSONProps(entrySchema, mapping)
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion] = apiextensions.JSONSchemaProps{
		Type: "object",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"entry": *entryProps,
		},
	}

	var statusSchemaRef *openapi3.SchemaRef
	for httpStatusCode, response := range operation.Responses {
		code, err := strconv.Atoi(httpStatusCode)
		if err != nil {
			return fmt.Errorf("error converting httpStatusCode to int: %w", err)
		}

		switch code {
		case http.StatusOK:
		case http.StatusCreated:
		default:
			continue
		}

		statusSchemaRef = anyEntry(response.Value.Content, nil).Schema
		break
	}

	if statusSchemaRef != nil {
		statusSchema := FilterSchemaProps("", true, statusSchemaRef, func(key string, schemaRef *openapi3.SchemaRef) bool {
			if key == "links" {
				return false
			}
			return schemaRef.Value.ReadOnly
		})
		statusProps := g.schemaPropsToJSONProps(statusSchema, mapping)
		if statusProps != nil {
			crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[mapping.MajorVersion] = *statusProps
		}
	}

	var params apiextensions.JSONSchemaProps
	params.Type = "object"
	params.Properties = make(map[string]apiextensions.JSONSchemaProps)
	for _, p := range operation.Parameters {
		switch p.Value.Name {
		case "includeCount":
		case "itemsPerPage":
		case "pageNum":
		case "envelope":
		case "pretty":
		default:
			props := g.schemaPropsToJSONProps(p.Value.Schema, mapping)
			props.Description = p.Value.Description
			params.Properties[p.Value.Name] = *props
		}
	}

	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion].Properties["parameters"] = params

	return nil
}
