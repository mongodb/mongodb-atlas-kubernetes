package plugins

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type PrintConditions struct{}

func (p *PrintConditions) Name() string {
	return "print_conditions"
}

func (p *PrintConditions) Process(req *MappingProcessorRequest) error {
	req.CRD.Spec.AdditionalPrinterColumns = []apiextensions.CustomResourceColumnDefinition{
		{
			JSONPath: `.status.conditions[?(@.type=="Ready")].status`,
			Name:     "Ready",
			Type:     "string",
		},
		{
			JSONPath: `.status.conditions[?(@.type=="Ready")].reason`,
			Name:     "Reason",
			Type:     "string",
		},
		{
			JSONPath: `.status.conditions[?(@.type=="State")].reason`,
			Name:     "State",
			Type:     "string",
		},
	}
	return nil
}
