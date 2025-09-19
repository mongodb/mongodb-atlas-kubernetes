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

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type StatusPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &StatusPlugin{}

func NewStatusPlugin(crd *apiextensions.CustomResourceDefinition) *StatusPlugin {
	return &StatusPlugin{
		crd: crd,
	}
}

func (s *StatusPlugin) Name() string {
	return "status"
}

func (s *StatusPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	if mappingConfig.StatusMapping.Schema == "" {
		return nil
	}

	statusSchema, ok := openApiSpec.Components.Schemas[mappingConfig.StatusMapping.Schema]
	if !ok {
		return fmt.Errorf("status schema %q not found in openapi spec", mappingConfig.StatusMapping.Schema)
	}

	statusProps := g.ConvertProperty(statusSchema, openapi3.NewSchemaRef("", openapi3.NewSchema()), &mappingConfig.StatusMapping, 0)
	statusProps.Description = fmt.Sprintf("The last observed Atlas state of the %v resource for version %v.", s.crd.Spec.Names.Singular, mappingConfig.MajorVersion)
	if statusProps != nil {
		s.crd.Spec.Validation.OpenAPIV3Schema.Properties["status"].Properties[mappingConfig.MajorVersion] = *statusProps
	}

	return nil
}
