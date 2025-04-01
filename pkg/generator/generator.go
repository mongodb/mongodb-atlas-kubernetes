package generator

import (
	"context"
	"fmt"
	"k8s.io/utils/ptr"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
)

func anyEntry[T any](source map[string]T, defaultValue T) T {
	for _, v := range source {
		return v
	}
	return defaultValue
}

type Generator struct {
	AtlasVersion string
	AtlasPath    string
	AtlasVerb    string
	GVK          string
	Categories   []string
	Spec         *openapi3.T
}

func NewGenerator(version, path, verb, gvk string, categories []string, spec *openapi3.T) *Generator {
	return &Generator{
		AtlasVersion: version,
		AtlasPath:    path,
		AtlasVerb:    verb,
		GVK:          gvk,
		Categories:   categories,
		Spec:         spec,
	}
}

func (g *Generator) Generate(ctx context.Context) (*apiextensions.CustomResourceDefinition, error) {
	var operation *openapi3.Operation
	switch g.AtlasVerb {
	case "post":
		operation = g.Spec.Paths[g.AtlasPath].Post
	default:
		return nil, fmt.Errorf("verb %q not supported", g.AtlasVerb)
	}

	crd, err := g.generateCustomResource(operation)
	if err != nil {
		return nil, err
	}

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

func guessKindToResource(kind runtimeschema.GroupVersionKind) ( /*plural*/ runtimeschema.GroupVersionResource /*singular*/, runtimeschema.GroupVersionResource) {
	kindName := kind.Kind
	if len(kindName) == 0 {
		return runtimeschema.GroupVersionResource{}, runtimeschema.GroupVersionResource{}
	}
	singularName := strings.ToLower(kindName)
	singular := kind.GroupVersion().WithResource(singularName)

	switch string(singularName[len(singularName)-1]) {
	case "s":
		return kind.GroupVersion().WithResource(singularName + "es"), singular
	case "x":
		return kind.GroupVersion().WithResource(singularName + "es"), singular
	case "y":
		return kind.GroupVersion().WithResource(strings.TrimSuffix(singularName, "y") + "ies"), singular
	}

	return kind.GroupVersion().WithResource(singularName + "s"), singular
}

func (g *Generator) generateCustomResource(operation *openapi3.Operation) (*apiextensions.CustomResourceDefinition, error) {
	gvk, _ := runtimeschema.ParseKindArg(g.GVK)
	pluralGvk, singularGvk := guessKindToResource(*gvk)

	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", pluralGvk.Resource, pluralGvk.Group),
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: pluralGvk.Group,
			Scope: apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     gvk.Kind,
				ListKind: fmt.Sprintf("%sList", gvk.Kind),
				Plural:   pluralGvk.Resource,
				Singular: singularGvk.Resource,
			},
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    gvk.Version,
					Served:  true,
					Storage: true,
				},
			},
		},
	}

	content := operation.RequestBody.Value.Content
	mediaType := anyEntry(content, nil)
	entrySchemaRef := mediaType.Schema
	entrySchema := FilterSchemaProps("", false, entrySchemaRef, func(key string, schemaRef *openapi3.SchemaRef) bool {
		return !schemaRef.Value.ReadOnly
	})
	entryProps := g.schemaPropsToJSONProps(entrySchema)
	crd.Spec.Validation = &apiextensions.CustomResourceValidation{
		OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
			Type: "object",
			Properties: map[string]apiextensions.JSONSchemaProps{
				"spec": {
					Type: "object",
					Properties: map[string]apiextensions.JSONSchemaProps{
						g.AtlasVersion: {
							Type: "object",
							Properties: map[string]apiextensions.JSONSchemaProps{
								"entry": *entryProps,
							},
						},
					},
				},
			},
		},
	}

	var statusSchemaRef *openapi3.SchemaRef
	for httpStatusCode, response := range operation.Responses {
		code, err := strconv.Atoi(httpStatusCode)
		if err != nil {
			return nil, fmt.Errorf("error converting httpStatusCode to int: %w", err)
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

	crd.Spec.Validation.OpenAPIV3Schema.Properties["status"] = apiextensions.JSONSchemaProps{
		Type: "object",
		Properties: map[string]apiextensions.JSONSchemaProps{
			"conditions": {
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
			},
		},
	}

	if statusSchemaRef != nil {
		statusSchema := FilterSchemaProps("", true, statusSchemaRef, func(key string, schemaRef *openapi3.SchemaRef) bool {
			if key == "links" {
				return false
			}
			return schemaRef.Value.ReadOnly
		})
		statusProps := g.schemaPropsToJSONProps(statusSchema)
		if statusProps != nil {
			crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[g.AtlasVersion] = *statusProps
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
			props := g.schemaPropsToJSONProps(p.Value.Schema)
			props.Description = p.Value.Description
			params.Properties[p.Value.Name] = *props
		}
	}
	crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[g.AtlasVersion].Properties["parameters"] = params

	// TODO: yaml.Marshal creates an empty status field that we should remove
	// StoredVersions is set to empty array instead of nil to bypass the following issue:
	// https://github.com/fybrik/openapi2crd/issues/1
	crd.Status.StoredVersions = []string{}

	// enable status subresource
	crd.Spec.Subresources = &apiextensions.CustomResourceSubresources{
		Status: &apiextensions.CustomResourceSubresourceStatus{},
	}

	crd.Spec.Names.Categories = g.Categories

	return crd, nil
}
