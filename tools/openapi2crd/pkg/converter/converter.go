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

package converter

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/processor"
	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/utils/ptr"
)

type PropertyConvertInput struct {
	Plugins          []processor.Processor
	Schema           *openapi3.SchemaRef
	ExtensionsSchema *openapi3.SchemaRef
	PropertyConfig   *configv1alpha1.PropertyMapping
	Depth            int
	Path             []string
}

func (i PropertyConvertInput) WithSchema(s *openapi3.SchemaRef) PropertyConvertInput {
	i.Schema = s

	return i
}

func (i PropertyConvertInput) WithExtensionsSchema(s *openapi3.SchemaRef) PropertyConvertInput {
	i.ExtensionsSchema = s

	return i
}

func (i PropertyConvertInput) WithDepthIncrement(inc int) PropertyConvertInput {
	i.Depth += inc

	return i
}

func (i PropertyConvertInput) WithPath(path []string) PropertyConvertInput {
	i.Path = path

	return i
}

type Converter interface {
	Convert(input PropertyConvertInput) *apiextensions.JSONSchemaProps
}

type PropertyConverter struct{}

func NewPropertyConverter() *PropertyConverter {
	return &PropertyConverter{}
}

func (c *PropertyConverter) Convert(input PropertyConvertInput) *apiextensions.JSONSchemaProps {
	depth := input.Depth
	if depth == 10 {
		return nil
	}

	schema := input.Schema
	if input.Schema == nil {
		return nil
	}

	path := input.Path
	if len(path) == 0 {
		path = []string{"$"}
	}

	propertySchema := schema.Value
	if propertySchema == nil {
		return nil
	}

	typ := ""
	if propertySchema.Type != nil && len(propertySchema.Type.Slice()) > 0 {
		typ = (*propertySchema.Type)[0]
	}
	example := apiextensions.JSON(propertySchema.Example)
	propertyConfig := input.PropertyConfig
	extensionsSchema := input.ExtensionsSchema

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
		Required:             propertySchema.Required,
		Items:                c.convertPropertyOrArray(input.WithSchema(propertySchema.Items).WithPath(append(path, "[*]"))),
		AllOf:                c.convertPropertySlice(propertySchema.AllOf, input),
		OneOf:                c.convertPropertySlice(propertySchema.OneOf, input),
		AnyOf:                c.convertPropertySlice(propertySchema.AnyOf, input),
		Properties:           c.convertPropertyMap(propertySchema.Properties, input),
		Not:                  c.Convert(input.WithSchema(propertySchema.Not).WithDepthIncrement(1)),
		AdditionalProperties: c.convertPropertyOrBool(input.WithSchema(propertySchema.AdditionalProperties.Schema)),
		Example:              &example,
	}

	for _, p := range input.Plugins {
		p.Process(processor.NewPropertyInput(propertyConfig, props, propertySchema, extensionsSchema, path...))
		if props == nil {
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
	props = c.transformations(
		props,
		PropertyConvertInput{
			Plugins:          input.Plugins,
			Schema:           schema,
			ExtensionsSchema: extensionsSchema,
			PropertyConfig:   propertyConfig,
			Depth:            depth,
			Path:             path,
		},
	)

	return props
}

func (c *PropertyConverter) convertPropertyOrBool(input PropertyConvertInput) *apiextensions.JSONSchemaPropsOrBool {
	if input.Schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema: c.Convert(input.WithDepthIncrement(1)),
		Allows: true,
	}
}

func (c *PropertyConverter) convertPropertyOrArray(input PropertyConvertInput) *apiextensions.JSONSchemaPropsOrArray {
	if input.Schema == nil {
		return nil
	}

	eSchema := input.ExtensionsSchema
	eSchema.Value.Items = openapi3.NewSchemaRef("", openapi3.NewSchema())

	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: c.Convert(input.WithSchema(eSchema).WithDepthIncrement(1)),
	}
}

func (c *PropertyConverter) convertPropertySlice(schemas openapi3.SchemaRefs, input PropertyConvertInput) []apiextensions.JSONSchemaProps {
	var s []apiextensions.JSONSchemaProps
	for _, schema := range schemas {
		input.Depth++
		result := c.Convert(input.WithSchema(schema).WithDepthIncrement(1))
		if result == nil {
			continue
		}
		s = append(s, *result)
	}
	return s
}

func (c *PropertyConverter) convertPropertyMap(schemaMap openapi3.Schemas, input PropertyConvertInput) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		childExtensionsSchema := openapi3.NewSchemaRef("", openapi3.NewSchema())
		result := c.Convert(
			input.WithSchema(schema).
				WithExtensionsSchema(childExtensionsSchema).
				WithPath(append(input.Path, key)).
				WithDepthIncrement(1))
		if result == nil {
			continue
		}

		propName := key
		if result.ID != "" { // workaround for the fact that CRD props do not let us specify its own property name
			propName = result.ID
			result.ID = ""
		}

		if input.ExtensionsSchema.Value.Properties == nil {
			input.ExtensionsSchema.Value.Properties = make(openapi3.Schemas)
		}
		input.ExtensionsSchema.Value.Properties[propName] = childExtensionsSchema

		m[propName] = *result
	}

	return m
}

func (c *PropertyConverter) transformations(props *apiextensions.JSONSchemaProps, input PropertyConvertInput) *apiextensions.JSONSchemaProps {
	result := props
	result = handleAdditionalProperties(result, input.Schema.Value.AdditionalProperties.Has)
	result = removeUnknownFormats(result)
	result = c.oneOfRefsTransform(result, input.Schema.Value.OneOf, input)

	return result
}

func (c *PropertyConverter) oneOfRefsTransform(props *apiextensions.JSONSchemaProps, oneOf openapi3.SchemaRefs, input PropertyConvertInput) *apiextensions.JSONSchemaProps {
	if props.OneOf != nil && len(props.Properties) == 0 && props.AdditionalProperties == nil {
		result := props.DeepCopy()
		result.Type = "object"
		result.OneOf = nil

		options := []apiextensions.JSON{}
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
			result := c.Convert(input.WithSchema(v).WithPath(append(input.Path, name)).WithDepthIncrement(1))
			if result == nil {
				continue
			}
			result.Properties[name] = *result
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
