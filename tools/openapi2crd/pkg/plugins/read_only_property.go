package plugins

import (
	"slices"

	"k8s.io/apimachinery/pkg/util/sets"
)

type ReadOnlyProperty struct{}

func (p *ReadOnlyProperty) Name() string {
	return "read_only_property"
}

func (p *ReadOnlyProperty) Process(req *PropertyProcessorRequest) error {
	if req.PropertyConfig == nil || !req.PropertyConfig.Filters.ReadOnly {
		return nil
	}

	if req.OpenAPISchema.ReadOnly {
		return nil
	}

	required := sets.New(req.OpenAPISchema.Required...)
	for name, p := range req.OpenAPISchema.Properties {
		if !p.Value.ReadOnly {
			required.Delete(name)
		}
	}
	req.Property.Required = required.UnsortedList()
	slices.Sort(req.Property.Required)

	// ignore root
	if len(req.Path) == 1 {
		return nil
	}

	req.Property = nil

	return nil
}
