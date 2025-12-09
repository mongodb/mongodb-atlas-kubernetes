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
	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/converter"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type CRDProcessorRequest struct {
	CRD       *apiextensions.CustomResourceDefinition
	CRDConfig *configv1alpha1.CRDConfig
}

type MappingProcessorRequest struct {
	CRD              *apiextensions.CustomResourceDefinition
	MappingConfig    *configv1alpha1.CRDMapping
	OpenAPISpec      *openapi3.T
	ExtensionsSchema *openapi3.Schema
	Converter        converter.Converter
}

type PropertyProcessorRequest struct {
	Property         *apiextensions.JSONSchemaProps
	PropertyConfig   *configv1alpha1.PropertyMapping
	OpenAPISchema    *openapi3.Schema
	ExtensionsSchema *openapi3.SchemaRef
	Path             []string
}

type ExtensionProcessorRequest struct {
	ExtensionsSchema *openapi3.Schema
	ApiDefinitions   map[string]configv1alpha1.OpenAPIDefinition
	MappingConfig    *configv1alpha1.CRDMapping
}

type Plugin[R any] interface {
	Name() string
	Process(request R) error
}

type CRDPlugin = Plugin[*CRDProcessorRequest]
type MappingPlugin = Plugin[*MappingProcessorRequest]
type PropertyPlugin = Plugin[*PropertyProcessorRequest]
type ExtensionPlugin = Plugin[*ExtensionProcessorRequest]

var _ CRDPlugin = &Base{}
var _ CRDPlugin = &MutualExclusiveMajorVersions{}
var _ MappingPlugin = &Entry{}
var _ MappingPlugin = &Status{}
var _ MappingPlugin = &Parameters{}
var _ MappingPlugin = &References{}
var _ MappingPlugin = &MajorVersion{}
var _ PropertyPlugin = &ReadOnlyProperties{}
var _ PropertyPlugin = &ReadWriteProperties{}
var _ PropertyPlugin = &SensitiveProperties{}
var _ PropertyPlugin = &SkippedProperties{}
var _ ExtensionPlugin = &AtlasSdkVersionPlugin{}
var _ ExtensionPlugin = &ReferenceExtensions{}
