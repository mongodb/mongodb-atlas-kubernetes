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
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/utils/ptr"
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

// schemaPropsToJSONProps converts openapi3.Schema to a JSONProps
func (g *Generator) schemaPropsToJSONProps(schemaRef *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps {
	var props *apiextensions.JSONSchemaProps

	if schemaRef == nil {
		return props
	}

	if len(path) == 0 {
		path = []string{"$"}
	}
	schemaProps := schemaRef.Value

	props = &apiextensions.JSONSchemaProps{
		// ID:               schemaProps.ID,
		// Schema:           apiextensions.JSONSchemaURL(string(schemaRef.Ref.)),
		// Ref:              ref,
		Description: schemaProps.Description,
		Type:        schemaProps.Type,
		//Format:      schemaProps.Format,
		Title: schemaProps.Title,
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
		MultipleOf:  schemaProps.MultipleOf,
		//Enum:        enumJSON(schemaProps.Enum),
		//MaxProperties:        castUInt64P(schemaProps.MaxProps),
		//MinProperties:        castUInt64(schemaProps.MinProps),
		Required:             schemaProps.Required,
		Items:                g.schemaToJSONSchemaPropsOrArray(schemaProps.Items, append(path, "[*]")),
		AllOf:                g.schemasToJSONSchemaPropsArray(schemaProps.AllOf, path),
		OneOf:                g.schemasToJSONSchemaPropsArray(schemaProps.OneOf, path),
		AnyOf:                g.schemasToJSONSchemaPropsArray(schemaProps.AnyOf, path),
		Not:                  g.schemaPropsToJSONProps(schemaProps.Not, path...),
		Properties:           g.schemasToJSONSchemaPropsMap(schemaProps.Properties, path),
		AdditionalProperties: g.schemaToJSONSchemaPropsOrBool(schemaProps.AdditionalProperties, path),
	}

	if props.Type == "" && props.Items == nil && len(props.Properties) == 0 {
		props.Type = "object"
		props.XPreserveUnknownFields = ptr.To(true)
	}

	// Apply custom transformations
	props = g.transformations(props, schemaRef, path)

	return props
}

func (g *Generator) transformations(props *apiextensions.JSONSchemaProps, schemaRef *openapi3.SchemaRef, path []string) *apiextensions.JSONSchemaProps {
	result := props
	result = handleAdditionalProperties(result, schemaRef.Value.AdditionalPropertiesAllowed)
	result = removeUnknownFormats(result)
	result = g.oneOfRefsTransform(result, schemaRef.Value.OneOf, path)
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
func (g *Generator) oneOfRefsTransform(props *apiextensions.JSONSchemaProps, oneOf openapi3.SchemaRefs, path []string) *apiextensions.JSONSchemaProps {
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
			result.Properties[name] = *g.schemaPropsToJSONProps(v, append(path, name)...)
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

func (g *Generator) schemasToJSONSchemaPropsArray(schemas openapi3.SchemaRefs, path []string) []apiextensions.JSONSchemaProps {
	var s []apiextensions.JSONSchemaProps
	for _, schema := range schemas {
		s = append(s, *g.schemaPropsToJSONProps(schema, path...))
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

func (g *Generator) schemaToJSONSchemaPropsOrArray(schema *openapi3.SchemaRef, path []string) *apiextensions.JSONSchemaPropsOrArray {
	if schema == nil {
		return nil
	}
	return &apiextensions.JSONSchemaPropsOrArray{
		Schema: g.schemaPropsToJSONProps(schema, path...),
	}
}

func (g *Generator) schemaToJSONSchemaPropsOrBool(schema *openapi3.SchemaRef, path []string) *apiextensions.JSONSchemaPropsOrBool {
	if schema == nil {
		return nil
	}

	return &apiextensions.JSONSchemaPropsOrBool{
		Schema: g.schemaPropsToJSONProps(schema, path...),
		Allows: true,
	}
}

func (g *Generator) schemasToJSONSchemaPropsMap(schemaMap openapi3.Schemas, path []string) map[string]apiextensions.JSONSchemaProps {
	m := make(map[string]apiextensions.JSONSchemaProps)
	for key, schema := range schemaMap {
		m[key] = *g.schemaPropsToJSONProps(schema, append(path, key)...)
	}
	return m
}

func castUInt64P(p *uint64) *int64 {
	if p == nil {
		return nil
	}
	val := int64(*p)
	return &val
}

func castUInt64(v uint64) *int64 {
	val := int64(v)
	if val == 0 {
		return nil
	}
	return &val
}
