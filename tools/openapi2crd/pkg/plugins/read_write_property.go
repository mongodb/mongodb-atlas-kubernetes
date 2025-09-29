package plugins

import (
	"slices"

	"k8s.io/apimachinery/pkg/util/sets"
)

type ReadWriteProperty struct{}

func (p *ReadWriteProperty) Name() string {
	return "read_write_property"
}

func (p *ReadWriteProperty) Process(req *PropertyProcessorRequest) error {
	if req.PropertyConfig == nil || !req.PropertyConfig.Filters.ReadWriteOnly {
		return nil
	}

	if req.OpenAPISchema.ReadOnly {
		req.Property = nil

		return nil
	}

	required := sets.New(req.OpenAPISchema.Required...)
	for name, prop := range req.OpenAPISchema.Properties {
		if prop.Value.ReadOnly {
			required.Delete(name)
		}
	}
	req.Property.Required = required.UnsortedList()
	slices.Sort(req.Property.Required)

	return nil
}
