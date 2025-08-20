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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/atlas"
	"github.com/mongodb/atlas2crd/pkg/config"
	"github.com/mongodb/atlas2crd/pkg/plugins"
	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"sigs.k8s.io/yaml"
)

type Generator struct {
	definitions map[string]v1alpha1.OpenAPIDefinition
	pluginsSets plugins.PluginSets
}

func NewGenerator(openAPIDefinitions map[string]v1alpha1.OpenAPIDefinition, pluginsSets plugins.PluginSets) *Generator {
	return &Generator{
		definitions: openAPIDefinitions,
		pluginsSets: pluginsSets,
	}
}

func (g *Generator) Generate(ctx context.Context, crdConfig *v1alpha1.CRDConfig) (*apiextensions.CustomResourceDefinition, error) {
	crd := &apiextensions.CustomResourceDefinition{}

	extensionsSchema := openapi3.NewSchema()
	extensionsSchema.Properties = map[string]*openapi3.SchemaRef{
		"spec": {Value: &openapi3.Schema{Properties: map[string]*openapi3.SchemaRef{}}},
	}

	pluginSet, err := g.pluginsSets.Default()
	if err != nil {
		return nil, fmt.Errorf("error getting default plugin set: %w", err)
	}

	for _, p := range pluginSet.Plugins {
		if err = p.Process(processor.NewCRDInput(crd, crdConfig)); err != nil {
			return nil, fmt.Errorf("error processing CRD in plugin %q: %w", p.Name(), err)
		}
	}

	for _, mapping := range crdConfig.Mappings {
		def, ok := g.definitions[mapping.OpenAPIRef.Name]
		if !ok {
			return nil, fmt.Errorf("no OpenAPI definition named %q found", mapping.OpenAPIRef.Name)
		}

		var openApiSpec *openapi3.T

		path := def.Path
		if path == "" {
			var err error
			path, err = atlas.LoadOpenAPIPath(def.Package)
			if err != nil {
				return nil, fmt.Errorf("error loading OpenAPI package %q: %w", def.Package, err)
			}
		}

		openApiSpec, err := config.LoadOpenAPI(path)
		if err != nil {
			return nil, fmt.Errorf("error loading spec: %w", err)
		}
		for _, p := range pluginSet.Plugins {
			if err := p.Process(processor.NewMappingInput(crd, &mapping, openApiSpec, extensionsSchema)); err != nil {
				return nil, fmt.Errorf("error processing plugin %s: %w", p.Name(), err)
			}
		}
	}

	clearPropertiesWithoutExtensions(extensionsSchema)
	if len(extensionsSchema.Properties) > 0 {
		d, err := yaml.Marshal(extensionsSchema)
		if err != nil {
			return nil, fmt.Errorf("error marshaling extensions schema: %w", err)
		}
		if crd.Annotations == nil {
			crd.Annotations = make(map[string]string)
		}
		crd.Annotations["api-mappings"] = string(d)
	}

	if err := ValidateCRD(ctx, crd); err != nil {
		log.Printf("Error validating CRD: %v", err)
	}

	return crd, nil
}

func (g *Generator) majorVersions(config v1alpha1.CRDConfig) []string {
	var result []string
	for _, m := range config.Mappings {
		result = append(result, "- "+m.MajorVersion)
	}

	return result
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

	if schema.AdditionalProperties.Schema != nil && clearPropertiesWithoutExtensions(schema.AdditionalProperties.Schema.Value) {
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
