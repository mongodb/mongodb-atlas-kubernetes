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
	"k8s.io/utils/ptr"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
)

func FilterSchemaProps(key string, relaxed bool, schema *openapi3.SchemaRef, predicate func(string, *openapi3.SchemaRef) bool) *openapi3.SchemaRef {
	valueCopy := *schema.Value
	schemaCopy := &openapi3.SchemaRef{
		Ref:   schema.Ref,
		Value: &valueCopy,
	}
	schemaValue := schemaCopy.Value
	isFiltered := predicate(key, schema)

	hasFilteredProps := false
	filteredProps := make(openapi3.Schemas)
	for key, schema := range schemaValue.Properties {
		filtered := FilterSchemaProps(key, relaxed, schema, predicate)
		if filtered != nil {
			filteredProps[key] = filtered
			hasFilteredProps = true
		}
	}
	schemaValue.Properties = filteredProps
	var required []string
	for _, r := range schemaValue.Required {
		if _, ok := filteredProps[r]; ok {
			required = append(required, r)
		}
	}
	schemaValue.Required = required

	hasFilteredItems := false
	if schemaValue.Items != nil {
		filteredItems := FilterSchemaProps(key+".items", relaxed, schemaValue.Items, predicate)
		if !isFiltered || filteredItems != nil {
			schemaValue.Items = filteredItems
		}
		if filteredItems != nil {
			hasFilteredItems = true
		}
	}

	isRelaxed := relaxed && (hasFilteredProps || hasFilteredItems)

	if isFiltered || isRelaxed {
		return schemaCopy
	}

	return nil
}

func jsonPath(path []string) string {
	result := strings.Join(path, ".")
	return strings.ReplaceAll(result, ".[*]", "[*]")
}

func isSkippedField(path []string, mapping *configv1alpha1.FieldMapping) bool {
	if mapping == nil {
		return false
	}

	p := jsonPath(path)

	for _, skippedField := range mapping.Filters.SkipProperties {
		if skippedField == p {
			return true
		}
	}

	return false
}

// SchemaPropsToJSONProps converts openapi3.Schema to a JSONProps
func (g *Generator) ConvertProperty(schema, extensionsSchema *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, path ...string) *apiextensions.JSONSchemaProps {
	if schema == nil {
		return nil
	}

	if len(path) == 0 {
		path = []string{"$"}
	}

	if isSkippedField(path, mapping) {
		return nil
	}

	var skipProperties []string
	if mapping != nil {
		skipProperties = mapping.Filters.SkipProperties
	}

	propertySchema := schema.Value
	example := apiextensions.JSON(propertySchema.Example)
	props := &apiextensions.JSONSchemaProps{
		//ID:               schemaProps.ID,
		//Schema:           apiextensions.JSONSchemaURL(string(schemaRef.Ref.)),
		//Ref:              ref,
		Description: propertySchema.Description,
		Type:        propertySchema.Type,
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
		Required:             filterSlice(propertySchema.Required, skipProperties),
		Items:                g.convertPropertyOrArray(propertySchema.Items, extensionsSchema, mapping, append(path, "[*]")),
		AllOf:                g.convertPropertySlice(propertySchema.AllOf, mapping, extensionsSchema, path),
		OneOf:                g.convertPropertySlice(propertySchema.OneOf, mapping, extensionsSchema, path),
		AnyOf:                g.convertPropertySlice(propertySchema.AnyOf, mapping, extensionsSchema, path),
		Not:                  g.ConvertProperty(propertySchema.Not, extensionsSchema, mapping, path...),
		Properties:           g.ConvertPropertyMap(propertySchema.Properties, extensionsSchema, mapping, path...),
		AdditionalProperties: g.convertPropertyOrBool(propertySchema.AdditionalProperties, extensionsSchema, mapping, path),
		Example:              &example,
	}

	for _, p := range g.plugins {
		p.ProcessProperty(g, mapping, props, propertySchema, extensionsSchema, path...)
	}

	if props.Type == "" && props.Items == nil && len(props.Properties) == 0 {
		props.Type = "object"
		props.XPreserveUnknownFields = ptr.To(true)
	}

	// Apply custom transformations
	props = g.transformations(props, schema, mapping, extensionsSchema, path)

	return props
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func filterSlice(source, by []string) []string {
	filtered := []string{}
	for _, s := range source {
		if !contains(by, "$."+s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (g *Generator) transformations(props *apiextensions.JSONSchemaProps, schemaRef *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, extensionsSchema *openapi3.SchemaRef, path []string) *apiextensions.JSONSchemaProps {
	result := props
	result = handleAdditionalProperties(result, schemaRef.Value.AdditionalPropertiesAllowed)
	result = removeUnknownFormats(result)
	result = g.oneOfRefsTransform(result, schemaRef.Value.OneOf, mapping, extensionsSchema, path)
	return result
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

// oneOfRefsTransform transforms oneOf with a list of $ref to a list of nullable properties
func (g *Generator) oneOfRefsTransform(props *apiextensions.JSONSchemaProps, oneOf openapi3.SchemaRefs, mapping *configv1alpha1.FieldMapping, extensionsSchema *openapi3.SchemaRef, path []string) *apiextensions.JSONSchemaProps {
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
			result.Properties[name] = *g.ConvertProperty(v, extensionsSchema, mapping, append(path, name)...)
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

func (g *Generator) convertPropertySlice(schemas openapi3.SchemaRefs, mapping *configv1alpha1.FieldMapping, extensionsSchema *openapi3.SchemaRef, path []string) []apiextensions.JSONSchemaProps {
	var s []apiextensions.JSONSchemaProps
	for _, schema := range schemas {
		s = append(s, *g.ConvertProperty(schema, extensionsSchema, mapping, path...))
	}
	return s
}

// enumJSON converts []interface{} to []JSON
func enumJSON(enum []interface{}) []apiextensions.JSON {
	var s []apiextensions.JSON
	for _, elt := range enum {
		s = append(s, apiextensions.JSON(elt))
	}
	return s
}

func (g *Generator) convertPropertyOrArray(schema, extensionsSchema *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, path []string) *apiextensions.JSONSchemaPropsOrArray {
	if schema == nil {
		return nil
	}
	extensionsSchema.Value.Items = openapi3.NewSchemaRef("", openapi3.NewSchema())
	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: g.ConvertProperty(schema, extensionsSchema.Value.Items, mapping, path...),
	}
}

func (g *Generator) convertPropertyOrBool(schema, extensionsSchema *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, path []string) *apiextensions.JSONSchemaPropsOrBool {
	if schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema: g.ConvertProperty(schema, extensionsSchema, mapping, path...),
		Allows: true,
	}
}

func (g *Generator) ConvertPropertyMap(schemaMap openapi3.Schemas, extensionsSchema *openapi3.SchemaRef, mapping *configv1alpha1.FieldMapping, path ...string) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		childExtensionsSchema := openapi3.NewSchemaRef("", openapi3.NewSchema())
		result := g.ConvertProperty(schema, childExtensionsSchema, mapping, append(path, key)...)
		if result == nil {
			continue
		}

		propName := key
		if result.ID != "" { // workaround for the fact that CRD props do not let us specify the name
			propName = result.ID
			result.ID = ""
		}

		if extensionsSchema.Value.Properties == nil {
			extensionsSchema.Value.Properties = make(openapi3.Schemas)
		}
		extensionsSchema.Value.Properties[propName] = childExtensionsSchema

		m[propName] = *result
	}
	return m
}
