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
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type AtlasSdkVersionPlugin struct {
	NoOp
	crd         *apiextensions.CustomResourceDefinition
	definitions map[string]configv1alpha1.OpenAPIDefinition
}

func NewAtlasSdkVersionPlugin(crd *apiextensions.CustomResourceDefinition, definitions map[string]configv1alpha1.OpenAPIDefinition) *AtlasSdkVersionPlugin {
	return &AtlasSdkVersionPlugin{
		crd:         crd,
		definitions: definitions,
	}
}

func (p *AtlasSdkVersionPlugin) Name() string {
	return "atlas_sdk_version"
}

func (p *AtlasSdkVersionPlugin) ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error {
	pkg := p.definitions[mappingConfig.OpenAPIRef.Name].Package
	if pkg == "" {
		return nil
	}

	extensions := extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Extensions
	if extensions == nil {
		extensions = map[string]interface{}{}
	}
	extensions["x-atlas-sdk-version"] = pkg
	extensionsSchema.Properties["spec"].Value.Properties[mappingConfig.MajorVersion].Value.Extensions = extensions

	return nil
}
