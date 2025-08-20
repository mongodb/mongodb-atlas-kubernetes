package plugins

import (
	"slices"

	"github.com/mongodb/atlas2crd/pkg/processor"
	"k8s.io/apimachinery/pkg/util/sets"
)

type ReadWriteOnlyProperties struct{}

func (p *ReadWriteOnlyProperties) Name() string {
	return "read_write_only_properties"
}

func (p *ReadWriteOnlyProperties) Process(input processor.Input) error {
	i, ok := input.(processor.PropertyInput)
	if !ok {
		return nil
	}
	propertyConfig := i.PropertyConfig
	props := i.KubeSchema
	propertySchema := i.OpenAPISchema
	path := i.Path

	if propertyConfig == nil || !propertyConfig.Filters.ReadWriteOnly {
		return nil
	}

	if propertySchema.ReadOnly {
		props = nil
		return nil
	}

	required := sets.New(propertySchema.Required...)
	for name, p := range propertySchema.Properties {
		if p.Value.ReadOnly {
			required.Delete(name)
		}
	}
	props.Required = required.UnsortedList()
	slices.Sort(props.Required)

	// ignore root
	if len(path) == 1 {
		return nil
	}

	return nil
}

func NewReadWriteOnlyPropertiesPlugin() *ReadWriteOnlyProperties {
	return &ReadWriteOnlyProperties{}
}
