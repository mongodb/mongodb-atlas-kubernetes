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

package plugins

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
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
		entryProps := g.ConvertProperty(entrySchema, extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Properties["entry"], &mappingConfig.EntryMapping, 0)

		entryProps.Description = fmt.Sprintf("The entry fields of the %v resource spec. These fields can be set for creating and updating %v.", s.crd.Spec.Names.Singular, s.crd.Spec.Names.Plural)
		s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion].Properties["entry"] = *entryProps
	}

	return nil
}
