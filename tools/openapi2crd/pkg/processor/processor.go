package processor

import (
	"github.com/getkin/kin-openapi/openapi3"
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/converter"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type Input interface {
	processorInput()
}

type CRDInput struct {
	CRD       *apiextensions.CustomResourceDefinition
	CRDConfig *configv1alpha1.CRDConfig
}

func (i *CRDInput) processorInput() {}

func NewCRDInput(crd *apiextensions.CustomResourceDefinition, crdConfig *configv1alpha1.CRDConfig) *CRDInput {
	return &CRDInput{
		CRD:       crd,
		CRDConfig: crdConfig,
	}
}

type MappingInput struct {
	CRD              *apiextensions.CustomResourceDefinition
	MappingConfig    *configv1alpha1.CRDMapping
	OpenAPISpec      *openapi3.T
	ExtensionsSchema *openapi3.Schema
	Converter        converter.Converter
}

func (i *MappingInput) processorInput() {}

func NewMappingInput(crd *apiextensions.CustomResourceDefinition, mappingConfig *configv1alpha1.CRDMapping, openApiSpec *openapi3.T, extensionsSchema *openapi3.Schema) *MappingInput {
	return &MappingInput{
		CRD:              crd,
		MappingConfig:    mappingConfig,
		OpenAPISpec:      openApiSpec,
		ExtensionsSchema: extensionsSchema,
	}
}

type PropertyInput struct {
	PropertyConfig   *configv1alpha1.PropertyMapping
	KubeSchema       *apiextensions.JSONSchemaProps
	OpenAPISchema    *openapi3.Schema
	ExtensionsSchema *openapi3.SchemaRef
	Path             []string
}

func (i PropertyInput) processorInput() {}

func NewPropertyInput(propertyConfig *configv1alpha1.PropertyMapping, kubeSchema *apiextensions.JSONSchemaProps, openApiSchema *openapi3.Schema, extensionsSchema *openapi3.SchemaRef, path ...string) *PropertyInput {
	return &PropertyInput{
		PropertyConfig:   propertyConfig,
		KubeSchema:       kubeSchema,
		OpenAPISchema:    openApiSchema,
		ExtensionsSchema: extensionsSchema,
		Path:             path,
	}
}

type Processor interface {
	Process(input Input) error
}
