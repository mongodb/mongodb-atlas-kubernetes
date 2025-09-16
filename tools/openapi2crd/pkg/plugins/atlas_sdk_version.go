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
