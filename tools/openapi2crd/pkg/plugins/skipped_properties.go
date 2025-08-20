package plugins

import (
	"slices"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/atlas2crd/pkg/processor"
)

type SkippedProperties struct{}

func (p *SkippedProperties) Name() string {
	return "skipped_properties"
}

func (p *SkippedProperties) Process(input processor.Input) error {
	i, ok := input.(processor.PropertyInput)
	if !ok {
		return nil
	}
	propertyConfig := i.PropertyConfig
	props := i.KubeSchema
	propertySchema := i.OpenAPISchema
	path := i.Path

	if isSkippedField(path, propertyConfig) {
		props = nil
		return nil
	}

	if propertyConfig == nil || len(propertyConfig.Filters.SkipProperties) == 0 {
		return nil
	}

	requiredPaths := make(map[string]string)
	for _, r := range propertySchema.Required {
		requiredPaths[jsonPath(append(path, r))] = r
	}

	for _, s := range propertyConfig.Filters.SkipProperties {
		if _, ok := requiredPaths[s]; ok {
			delete(requiredPaths, s)
		}
	}

	props.Required = make([]string, 0, len(props.Required))
	for _, r := range requiredPaths {
		props.Required = append(props.Required, r)
	}

	slices.Sort(props.Required)

	return nil
}

func NewSkippedPropertiesPlugin() *SkippedProperties {
	return &SkippedProperties{}
}

func isSkippedField(path []string, mapping *configv1alpha1.PropertyMapping) bool {
	if mapping == nil {
		return false
	}

	p := jsonPath(path)

	for _, skippedField := range mapping.Filters.SkipProperties {
		if skippedField == p {
			return true
		}
	}

	return false
}
