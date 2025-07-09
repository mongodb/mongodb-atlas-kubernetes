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
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mongodb/atlas2crd/pkg/atlas"
	"log"
	"sigs.k8s.io/yaml"
	"strings"

	"github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/config"
	"github.com/mongodb/atlas2crd/pkg/plugins"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/utils/ptr"
)

type Generator struct {
	config      v1alpha1.CRDConfig
	definitions map[string]v1alpha1.OpenAPIDefinition
	plugins     []plugins.Plugin
}

func NewGenerator(crdConfig v1alpha1.CRDConfig, definitions []v1alpha1.OpenAPIDefinition) *Generator {
	definitionsMap := map[string]v1alpha1.OpenAPIDefinition{}
	for _, def := range definitions {
		definitionsMap[def.Name] = def
	}
	return &Generator{
		config:      crdConfig,
		definitions: definitionsMap,
	}
}

func (g *Generator) majorVersions() []string {
	var result []string
	for _, m := range g.config.Mappings {
		result = append(result, "- "+m.MajorVersion)
	}
	return result
}

func (g *Generator) Generate(ctx context.Context) (*apiextensions.CustomResourceDefinition, error) {
	pluralGvk, singularGvk := guessKindToResource(g.config.GVK)

	crd := &apiextensions.CustomResourceDefinition{
		ObjectMeta: v1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", pluralGvk.Resource, pluralGvk.Group),
		},
		Spec: apiextensions.CustomResourceDefinitionSpec{
			Group: pluralGvk.Group,
			Scope: apiextensions.NamespaceScoped,
			Names: apiextensions.CustomResourceDefinitionNames{
				Kind:     g.config.GVK.Kind,
				ListKind: fmt.Sprintf("%sList", g.config.GVK.Kind),
				Plural:   pluralGvk.Resource,
				Singular: singularGvk.Resource,
			},
			Versions: []apiextensions.CustomResourceDefinitionVersion{
				{
					Name:    g.config.GVK.Version,
					Served:  true,
					Storage: true,
				},
			},
			PreserveUnknownFields: ptr.To(false),
			Validation: &apiextensions.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensions.JSONSchemaProps{
					Type:        "object",
					Description: fmt.Sprintf("A %v, managed by the MongoDB Kubernetes Atlas Operator.", singularGvk.Resource),
					Properties: map[string]apiextensions.JSONSchemaProps{
						"spec": {
							Type: "object",
							Description: fmt.Sprintf(`Specification of the %v supporting the following versions:

%v

At most one versioned spec can be specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource, strings.Join(g.majorVersions(), "\n")),
							Properties: map[string]apiextensions.JSONSchemaProps{},
						},
						"status": {
							Type:        "object",
							Description: fmt.Sprintf(`Most recently observed read-only status of the %v for the specified resource version. This data may not be up to date and is populated by the system. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status`, singularGvk.Resource),
							Properties:  map[string]apiextensions.JSONSchemaProps{},
						},
					},
				},
			},
		},
	}

	g.plugins = []plugins.Plugin{
		plugins.NewMajorVersionPlugin(crd),
		plugins.NewParametersPlugin(crd),
		plugins.NewEntryPlugin(crd),
		plugins.NewStatusPlugin(crd),
		plugins.NewSensitivePropertiesPlugin(),
		plugins.NewSkippedPropertiesPlugin(),
		plugins.NewReadOnlyPropertiesPlugin(),
		plugins.NewReadWriteOnlyPropertiesPlugin(),
		plugins.NewReferencesPlugin(crd),
		plugins.NewMutualExclusiveMajorVersions(crd),
	}

	extensionsSchema := openapi3.NewSchema()
	extensionsSchema.Properties = map[string]*openapi3.SchemaRef{
		"spec": {Value: &openapi3.Schema{Properties: map[string]*openapi3.SchemaRef{}}},
	}

	for _, mapping := range g.config.Mappings {
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
		for _, p := range g.plugins {
			if err := p.ProcessMapping(g, &mapping, openApiSpec, extensionsSchema); err != nil {
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

	crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties["conditions"] = apiextensions.JSONSchemaProps{
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
	}

	crd.Status.StoredVersions = []string{}

	// enable status subresource
	crd.Spec.Subresources = &apiextensions.CustomResourceSubresources{
		Status: &apiextensions.CustomResourceSubresourceStatus{},
	}

	crd.Spec.Names.Categories = g.config.Categories
	crd.Spec.Names.ShortNames = g.config.ShortNames

	for _, version := range crd.Spec.Versions {
		if version.Storage {
			crd.Status.StoredVersions = append(crd.Status.StoredVersions, version.Name)
		}
	}

	for _, p := range g.plugins {
		if err := p.ProcessCRD(g, &g.config); err != nil {
			return nil, fmt.Errorf("error processing CRD in plugin %q: %w", p.Name(), err)
		}
	}

	if err := ValidateCRD(ctx, crd); err != nil {
		log.Printf("Error validating CRD: %v", err)
	}

	return crd, nil
}

func guessKindToResource(gvk v1.GroupVersionKind) ( /*plural*/ runtimeschema.GroupVersionResource /*singular*/, runtimeschema.GroupVersionResource) {
	runtimeGVK := runtimeschema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
	kindName := runtimeGVK.Kind
	if len(kindName) == 0 {
		return runtimeschema.GroupVersionResource{}, runtimeschema.GroupVersionResource{}
	}
	singularName := strings.ToLower(kindName)
	singular := runtimeGVK.GroupVersion().WithResource(singularName)

	switch string(singularName[len(singularName)-1]) {
	case "s":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "x":
		return runtimeGVK.GroupVersion().WithResource(singularName + "es"), singular
	case "y":
		return runtimeGVK.GroupVersion().WithResource(strings.TrimSuffix(singularName, "y") + "ies"), singular
	}

	return runtimeGVK.GroupVersion().WithResource(singularName + "s"), singular
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
