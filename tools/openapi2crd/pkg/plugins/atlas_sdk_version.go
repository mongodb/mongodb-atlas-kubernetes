package plugins

import (
	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/processor"
)

type AtlasSdkVersionPlugin struct {
	definitions map[string]configv1alpha1.OpenAPIDefinition
}

func (p *AtlasSdkVersionPlugin) Name() string {
	return "atlas_sdk_version"
}

func (p *AtlasSdkVersionPlugin) Process(input processor.Input) error {
	i, ok := input.(*processor.MappingInput)
	if !ok {
		return nil // No operation to perform
	}

	mappingConfig := i.MappingConfig
	extensionsSchema := i.ExtensionsSchema

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

func NewAtlasSdkVersionPlugin(definitions map[string]configv1alpha1.OpenAPIDefinition) *AtlasSdkVersionPlugin {
	return &AtlasSdkVersionPlugin{
		definitions: definitions,
	}
}
