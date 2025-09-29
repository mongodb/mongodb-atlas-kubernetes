package plugins

import (
	"slices"

	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type SkippedProperties struct{}

func (p *SkippedProperties) Name() string {
	return "skipped_property"
}

func (p *SkippedProperties) Process(req *PropertyProcessorRequest) error {
	if isSkippedField(req.Path, req.PropertyConfig) {
		req.Property = nil

		return nil
	}

	requiredPaths := make(map[string]string)
	for _, r := range req.OpenAPISchema.Required {
		requiredPaths[jsonPath(append(req.Path, r))] = r
	}

	for _, s := range req.PropertyConfig.Filters.SkipProperties {
		if _, ok := requiredPaths[s]; ok {
			delete(requiredPaths, s)
		}
	}

	req.Property.Required = make([]string, 0, len(req.Property.Required))
	for _, r := range requiredPaths {
		req.Property.Required = append(req.Property.Required, r)
	}

	slices.Sort(req.Property.Required)

	return nil
}

func isSkippedField(path []string, mapping *configv1alpha1.PropertyMapping) bool {
	p := jsonPath(path)

	for _, skippedField := range mapping.Filters.SkipProperties {
		if skippedField == p {
			return true
		}
	}

	return false
}
