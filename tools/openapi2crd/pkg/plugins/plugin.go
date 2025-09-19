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

type Plugin interface {
	Name() string
	ProcessCRD(g Generator, crdConfig *configv1alpha1.CRDConfig) error
	ProcessMapping(g Generator, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) error
	ProcessProperty(g Generator, propertyConfig *configv1alpha1.PropertyMapping, props *apiextensions.JSONSchemaProps, propertySchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *apiextensions.JSONSchemaProps
}

type Generator interface {
	ConvertProperty(schema, extensionsSchema *openapi3.SchemaRef, propertyConfig *configv1alpha1.PropertyMapping, depth int, path ...string) *apiextensions.JSONSchemaProps
}
