package plugins

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/processor"
)

type EntryPlugin struct{}

func (p *EntryPlugin) Name() string {
	return "entry"
}

func (p *EntryPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	openApiSpec := i.OpenAPISpec
	extensionsSchema := i.ExtensionsSchema
	crd := i.CRD

	var entrySchema *openapi3.SchemaRef
	switch {
	case mappingConfig.EntryMapping.Schema != "":
		var ok bool
		entrySchema, ok = openApiSpec.Components.Schemas[mappingConfig.EntryMapping.Schema]
		if !ok {
			return fmt.Errorf("entry schema %q not found in openapi spec", mappingConfig.EntryMapping.Schema)
		}
	case mappingConfig.EntryMapping.Path.Name != "":
		entrySchema = openApiSpec.Paths.Find(mappingConfig.EntryMapping.Path.Name).Operations()[strings.ToUpper(mappingConfig.EntryMapping.Path.Verb)].RequestBody.Value.Content[mappingConfig.EntryMapping.Path.RequestBody.MimeType].Schema
	}

	extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"entry": {Value: &openapi3.Schema{}},
			},
		},
	}

	if entrySchema != nil {
		entryProps := i.Converter.Convert(entrySchema, extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Properties["entry"], &mappingConfig.EntryMapping, 0)

		entryProps.Description = fmt.Sprintf("The entry fields of the %v resource spec. These fields can be set for creating and updating %v.", crd.Spec.Names.Singular, crd.Spec.Names.Plural)
		crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion].Properties["entry"] = *entryProps
	}

	return nil
}

func NewEntryPlugin() *EntryPlugin {
	return &EntryPlugin{}
}
