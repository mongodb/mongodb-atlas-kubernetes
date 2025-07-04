package plugins

import (
	"errors"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"sigs.k8s.io/yaml"
	"strings"
)

type EntryPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &EntryPlugin{}

func NewEntryPlugin(crd *apiextensions.CustomResourceDefinition) *EntryPlugin {
	return &EntryPlugin{
		crd: crd,
	}
}

func (s *EntryPlugin) Name() string {
	return "entry"
}

func (s *EntryPlugin) ProcessMapping(g Generator, mapping *configv1alpha1.CRDMapping, openApiSpec *openapi3.T) error {
	var entrySchemaRef *openapi3.SchemaRef
	switch {
	case mapping.EntryMapping.Schema != "":
		var ok bool
		entrySchemaRef, ok = openApiSpec.Components.Schemas[mapping.EntryMapping.Schema]
		if !ok {
			return fmt.Errorf("entry schema %q not found in openapi spec", mapping.EntryMapping.Schema)
		}
	case mapping.EntryMapping.Path.Name != "":
		entrySchemaRef = openApiSpec.Paths[mapping.EntryMapping.Path.Name].Operations()[strings.ToUpper(mapping.EntryMapping.Path.Verb)].RequestBody.Value.Content[mapping.EntryMapping.Path.RequestBody.MimeType].Schema
	default:
		return errors.New("entry schema not found in spec")
	}

	extensionsSchema := openapi3.NewSchema()
	extensionsSchema.Properties = map[string]*openapi3.SchemaRef{
		"spec": {Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				mapping.MajorVersion: {Value: &openapi3.Schema{
					Properties: map[string]*openapi3.SchemaRef{
						"entry": {Value: &openapi3.Schema{}},
					},
				}},
			},
		}},
	}
	entryProps := g.ConvertProperty(entrySchemaRef, &mapping.EntryMapping, extensionsSchema.Properties["spec"].Value.Properties[mapping.MajorVersion].Value.Properties["entry"])
	clearPropertiesWithoutExtensions(extensionsSchema)
	if len(extensionsSchema.Properties) > 0 {
		d, err := yaml.Marshal(extensionsSchema)
		if err != nil {
			return fmt.Errorf("error marshaling extensions schema: %w", err)
		}
		if s.crd.Annotations == nil {
			s.crd.Annotations = make(map[string]string)
		}
		s.crd.Annotations["api-mappings"] = string(d)
	}

	entryProps.Description = fmt.Sprintf("The entry fields of the %v resource spec. These fields can be set for creating and updating %v.", s.crd.Spec.Names.Singular, s.crd.Spec.Names.Plural)
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mapping.MajorVersion].Properties["entry"] = *entryProps
	return nil
}

func clearPropertiesWithoutExtensions(schema *openapi3.Schema) bool {
	if schema == nil {
		return false
	}
	hasExtensions := len(schema.Extensions) > 0

	var toDelete []string
	for k, prop := range schema.Properties {
		if !clearPropertiesWithoutExtensions(prop.Value) {
			toDelete = append(toDelete, k)
		} else {
			hasExtensions = true
		}
	}
	for _, k := range toDelete {
		delete(schema.Properties, k)
	}

	if schema.AdditionalProperties != nil && clearPropertiesWithoutExtensions(schema.AdditionalProperties.Value) {
		hasExtensions = true
	}
	if schema.Items != nil && clearPropertiesWithoutExtensions(schema.Items.Value) {
		hasExtensions = true
	}
	for _, ref := range schema.AllOf {
		if clearPropertiesWithoutExtensions(ref.Value) {
			hasExtensions = true
		}
	}
	for _, ref := range schema.AnyOf {
		if clearPropertiesWithoutExtensions(ref.Value) {
			hasExtensions = true
		}
	}
	for _, ref := range schema.OneOf {
		if clearPropertiesWithoutExtensions(ref.Value) {
			hasExtensions = true
		}
	}
	return hasExtensions
}
