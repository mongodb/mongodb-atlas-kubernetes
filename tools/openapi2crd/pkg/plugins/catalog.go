package plugins

import (
	"fmt"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
)

type PluginCatalog map[string]Plugin

func (p *PluginCatalog) Get(name string) (Plugin, error) {
	if plugin, ok := (*p)[name]; ok {
		return plugin, nil
	}

	return nil, fmt.Errorf("plugin %s not found", name)
}

func NewPluginCatalog(openAPIDefinitions map[string]configv1alpha1.OpenAPIDefinition) PluginCatalog {
	return PluginCatalog{
		"crd":                          NewCrdPlugin(),
		"majorVersion":                 NewMajorVersionPlugin(),
		"parameters":                   NewParametersPlugin(),
		"entry":                        NewEntryPlugin(),
		"status":                       NewStatusPlugin(),
		"sensitiveProperties":          NewSensitivePropertiesPlugin(),
		"skippedProperties":            NewSkippedPropertiesPlugin(),
		"readOnlyProperties":           NewReadOnlyPropertiesPlugin(),
		"readWriteOnlyProperties":      NewReadWriteOnlyPropertiesPlugin(),
		"references":                   NewReferencesPlugin(),
		"mutualExclusiveMajorVersions": NewMutualExclusiveMajorVersions(),
		"atlasSdkVersion":              NewAtlasSdkVersionPlugin(openAPIDefinitions),
	}
}
