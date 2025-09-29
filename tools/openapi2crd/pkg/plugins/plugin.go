package plugins

import (
	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"

	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
	"tools/openapi2crd/pkg/converter"
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
var _ MappingPlugin = &Parameter{}
var _ MappingPlugin = &Reference{}
var _ MappingPlugin = &MajorVersion{}
var _ PropertyPlugin = &ReadOnlyProperty{}
var _ PropertyPlugin = &ReadWriteProperty{}
var _ PropertyPlugin = &SensitiveProperty{}
var _ PropertyPlugin = &SkippedProperty{}
var _ ExtensionPlugin = &AtlasSdkVersionPlugin{}
var _ ExtensionPlugin = &ReferenceMetadata{}
