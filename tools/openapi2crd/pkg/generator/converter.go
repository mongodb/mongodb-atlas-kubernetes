// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package generator

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/utils/ptr"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/plugins"
)

func (g *Generator) Convert(input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
	if input.Depth == 10 {
		return nil
	}

	if input.Schema == nil {
		return nil
	}

	if len(input.Path) == 0 {
		input.Path = []string{"$"}
	}

	propertySchema := input.Schema.Value
	if propertySchema == nil {
		return nil
	}

	if _, seen := input.Visited[propertySchema]; seen {
		return &apiextensions.JSONSchemaProps{
			XPreserveUnknownFields: ptr.To(true),
		}
	}
	// Mark this schema as visited for the duration of this call (DFS backtracking).
	// All child recursions see it as an ancestor; defer removes it on return so
	// sibling branches at the caller's level do not falsely detect a cycle.
	input.Visited[propertySchema] = struct{}{}
	defer delete(input.Visited, propertySchema)

	typ := ""
	if propertySchema.Type != nil && len(propertySchema.Type.Slice()) > 0 {
		typ = (*propertySchema.Type)[0]
	}
	example := apiextensions.JSON(propertySchema.Example)
	extensionSchemaRef := input.ExtensionsSchemaRef
	props := &apiextensions.JSONSchemaProps{
		//ID:               schemaProps.ID,
		//Schema:           apiextensions.JSONSchemaURL(string(schemaRef.Ref.)),
		//Ref:              ref,
		Description: propertySchema.Description,
		Type:        typ,
		//Format:      schemaProps.Format,
		Title: propertySchema.Title,
		//Maximum:          schemaProps.Max,
		//ExclusiveMaximum: schemaProps.ExclusiveMax,
		//Minimum:          schemaProps.Min,
		//ExclusiveMinimum: schemaProps.ExclusiveMin,
		//MaxLength:        castUInt64P(schemaProps.MaxLength),
		//MinLength:        castUInt64(schemaProps.MinLength),
		// patterns seem to be incompatible in Atlas OpenAPI
		//Pattern:              schemaProps.Pattern,
		//MaxItems:             castUInt64P(schemaProps.MaxItems),
		//MinItems:             castUInt64(schemaProps.MinItems),
		UniqueItems: false, // The field uniqueItems cannot be set to true.
		MultipleOf:  propertySchema.MultipleOf,
		//Enum:        enumJSON(schemaProps.Enum),
		//MaxProperties:        castUInt64P(schemaProps.MaxProps),
		//MinProperties:        castUInt64(schemaProps.MinProps),
		Required: propertySchema.Required,
		Items: g.convertPropertyOrArray(
			input.Child(propertySchema.Items, extensionSchemaRef, input.Depth, append(input.Path, "[*]")),
		),
		AllOf: g.convertPropertySlice(propertySchema.AllOf, input),
		OneOf: g.convertPropertySlice(propertySchema.OneOf, input),
		AnyOf: g.convertPropertySlice(propertySchema.AnyOf, input),
		Not: g.Convert(
			input.Child(propertySchema.Not, extensionSchemaRef, input.Depth+1, input.Path),
		),
		Properties: g.convertPropertyMap(propertySchema.Properties, input),
		AdditionalProperties: g.convertPropertyOrBool(
			input.Child(propertySchema.AdditionalProperties.Schema, extensionSchemaRef, input.Depth, input.Path),
		),
		Example: &example,
	}

	for _, p := range g.pluginSet.Property {
		req := &plugins.PropertyProcessorRequest{
			Property:         props,
			PropertyConfig:   input.PropertyConfig,
			OpenAPISchema:    propertySchema,
			ExtensionsSchema: extensionSchemaRef,
			Path:             input.Path,
		}
		err := p.Process(req)
		// Currently, property plugins are not expected to return an error.
		// If an error case is introduced in the future, we should handle it appropriately.
		if err != nil {
			return nil
		}

		if req.Property == nil {
			return nil
		}
	}

	if props.Type == "" {
		props.Type = "object"
	}

	if props.Type == "object" && props.Items == nil && len(props.Properties) == 0 && props.AdditionalProperties == nil {
		props.XPreserveUnknownFields = ptr.To(true)
	}

	// Apply custom transformations
	props = g.transformations(
		props,
		input.Child(input.Schema, extensionSchemaRef, input.Depth, input.Path),
	)

	return props
}

func (g *Generator) convertPropertyOrBool(input converter.PropertyConvertInput) *apiextensions.JSONSchemaPropsOrBool {
	if input.Schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema:  g.Convert(input.Child(input.Schema, input.ExtensionsSchemaRef, input.Depth+1, input.Path)),
		Allows: true,
	}
}

func (g *Generator) convertPropertyOrArray(input converter.PropertyConvertInput) *apiextensions.JSONSchemaPropsOrArray {
	if input.Schema == nil {
		return nil
	}

	input.ExtensionsSchemaRef.Value.Items = openapi3.NewSchemaRef("", openapi3.NewSchema())

	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: g.Convert(input.Child(input.Schema, input.ExtensionsSchemaRef.Value.Items, input.Depth+1, input.Path)),
	}
}

func (g *Generator) convertPropertySlice(schemas openapi3.SchemaRefs, input converter.PropertyConvertInput) []apiextensions.JSONSchemaProps {
	if len(schemas) == 0 {
		return nil
	}

	s := make([]apiextensions.JSONSchemaProps, 0, len(schemas))

	for _, schema := range schemas {
		input.Depth++
		result := g.Convert(
			input.Child(schema, input.ExtensionsSchemaRef, input.Depth+1, input.Path),
		)
		if result == nil {
			continue
		}
		s = append(s, *result)
	}

	return s
}

func (g *Generator) convertPropertyMap(schemaMap openapi3.Schemas, input converter.PropertyConvertInput) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		childExtensionsSchema := openapi3.NewSchemaRef("", openapi3.NewSchema())
		result := g.Convert(
			input.Child(schema, childExtensionsSchema, input.Depth+1, append(input.Path, key)),
		)
		if result == nil {
			continue
		}

		propName := key
		if result.ID != "" { // workaround for the fact that CRD props do not let us specify its own property name
			propName = result.ID
			result.ID = ""
		}

		if input.ExtensionsSchemaRef.Value.Properties == nil {
			input.ExtensionsSchemaRef.Value.Properties = make(openapi3.Schemas)
		}
		input.ExtensionsSchemaRef.Value.Properties[propName] = childExtensionsSchema

		m[propName] = *result
	}

	return m
}

func (g *Generator) transformations(props *apiextensions.JSONSchemaProps, input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
	result := props
	result = handleAdditionalProperties(result, input.Schema.Value.AdditionalProperties.Has)
	result = removeUnknownFormats(result)
	result = g.oneOfRefsTransform(result, input.Schema.Value.OneOf, input)

	return result
}

func (g *Generator) oneOfRefsTransform(props *apiextensions.JSONSchemaProps, oneOf openapi3.SchemaRefs, input converter.PropertyConvertInput) *apiextensions.JSONSchemaProps {
	if props.OneOf != nil && len(props.Properties) == 0 && props.AdditionalProperties == nil {
		result := props.DeepCopy()
		result.Type = "object"
		result.OneOf = nil
		if result.Properties == nil {
			result.Properties = map[string]apiextensions.JSONSchemaProps{}
		}

		options := make([]apiextensions.JSON, 0, len(oneOf))
		for _, v := range oneOf {
			if v.Ref == "" {
				// this transform does not apply here
				// return the original props
				return props
			}
			name := v.Ref
			name = name[strings.LastIndex(name, "/")+1:]
			name = strcase.LowerCamelCase(name)
			options = append(options, name)
			optionSchema := g.Convert(
				input.Child(v, input.ExtensionsSchemaRef, input.Depth+1, append(input.Path, name)),
			)
			if optionSchema == nil {
				continue
			}
			result.Properties[name] = *optionSchema
		}

		result.Properties["type"] = apiextensions.JSONSchemaProps{
			Type:        "string",
			Enum:        options,
			Description: "Type is the discriminator for the different possible values",
		}

		return result
	}

	return props
}

func handleAdditionalProperties(props *apiextensions.JSONSchemaProps, additionalPropertiesAllowed *bool) *apiextensions.JSONSchemaProps {
	if additionalPropertiesAllowed != nil && *additionalPropertiesAllowed {
		props.XPreserveUnknownFields = additionalPropertiesAllowed
	}
	return props
}

func removeUnknownFormats(props *apiextensions.JSONSchemaProps) *apiextensions.JSONSchemaProps {
	switch props.Format {
	case "int32", "int64", "float", "double", "byte", "date", "date-time", "password":
	// Valid formats https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#format
	case "":
		props.Format = ""
	default:
		props.Format = ""
	}
	return props
}
