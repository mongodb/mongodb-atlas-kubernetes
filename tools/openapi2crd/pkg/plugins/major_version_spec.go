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

type MajorVersionSpecPlugin struct {
	NoOp
	crd *apiextensions.CustomResourceDefinition
}

var _ Plugin = &MajorVersionSpecPlugin{}

func NewMajorVersionPlugin(crd *apiextensions.CustomResourceDefinition) *MajorVersionSpecPlugin {
	return &MajorVersionSpecPlugin{
		crd: crd,
	}
}

func (s *MajorVersionSpecPlugin) Name() string {
	return "major_version"
}

func (s *MajorVersionSpecPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	s.crd.Spec.Validation.OpenAPIV3Schema.Properties["spec"].Properties[mappingConfig.MajorVersion] = apiextensions.JSONSchemaProps{
		Type:        "object",
		Description: fmt.Sprintf("The spec of the %v resource for version %v.", s.crd.Spec.Names.Singular, mappingConfig.MajorVersion),
		Properties:  map[string]apiextensions.JSONSchemaProps{},
	}
	return nil
}
