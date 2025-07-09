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

func (s *EntryPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	var entrySchema *openapi3.SchemaRef
	switch {
	case mappingConfig.EntryMapping.Schema != "":
		var ok bool
		entrySchema, ok = openApiSpec.Components.Schemas[mappingConfig.EntryMapping.Schema]
		if !ok {
			return fmt.Errorf("entry schema %q not found in openapi spec", mappingConfig.EntryMapping.Schema)
		}
	case mappingConfig.EntryMapping.Path.Name != "":
		entrySchema = openApiSpec.Paths[mappingConfig.EntryMapping.Path.Name].Operations()[strings.ToUpper(mappingConfig.EntryMapping.Path.Verb)].RequestBody.Value.Content[mappingConfig.EntryMapping.Path.RequestBody.MimeType].Schema
	default:
		return errors.New("entry schema not found in spec")
	}

	extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Properties: map[string]*openapi3.SchemaRef{
				"entry": {Value: &openapi3.Schema{}},
			},
		},
	}

	entryProps := g.ConvertProperty(entrySchema, extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Properties["entry"], &mappingConfig.EntryMapping)
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
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion].Properties["entry"] = *entryProps
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
