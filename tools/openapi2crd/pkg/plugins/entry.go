package plugins

import (
	"fmt"
	"strings"

	"tools/openapi2crd/pkg/converter"

	"github.com/getkin/kin-openapi/openapi3"
)

// Entry is a plugin that processes the entry mapping configuration and adds the entry schema to the CRD's spec validation schema.
// It requires the base and major_version plugin to be run first.
type Entry struct{}

func (p *Entry) Name() string {
	return "entry"
}

func (p *Entry) Process(req *MappingProcessorRequest) error {
	var entrySchema *openapi3.SchemaRef
	switch {
	case req.MappingConfig.EntryMapping.Schema != "":
		var ok bool
		entrySchema, ok = req.OpenAPISpec.Components.Schemas[req.MappingConfig.EntryMapping.Schema]
		if !ok {
			return fmt.Errorf("entry schema %q not found in openapi spec", req.MappingConfig.EntryMapping.Schema)
		}
	case req.MappingConfig.EntryMapping.Path.Name != "":
		entrySchema = req.OpenAPISpec.Paths.
			Find(req.MappingConfig.EntryMapping.Path.Name).
			Operations()[strings.ToUpper(req.MappingConfig.EntryMapping.Path.Verb)].
			RequestBody.Value.Content[req.MappingConfig.EntryMapping.Path.RequestBody.MimeType].Schema
	}

	req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"entry": {Value: &openapi3.Schema{}},
			},
		},
	}

	if entrySchema != nil {
		entryProps := req.Converter.Convert(
			converter.PropertyConvertInput{
				Schema:              entrySchema,
				ExtensionsSchemaRef: req.ExtensionsSchema.Properties["spec"].Value.Properties[req.MappingConfig.MajorVersion].Value.Properties["entry"],
				PropertyConfig:      &req.MappingConfig.EntryMapping,
				Depth:               0,
				Path:                nil,
			},
		)

		entryProps.Description = fmt.Sprintf(
			"The entry fields of the %v resource spec. These fields can be set for creating and updating %v.",
			req.CRD.Spec.Names.Singular,
			req.CRD.Spec.Names.Plural,
		)
		req.CRD.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[req.MappingConfig.MajorVersion].Properties["entry"] = *entryProps
	}

	return nil
}
