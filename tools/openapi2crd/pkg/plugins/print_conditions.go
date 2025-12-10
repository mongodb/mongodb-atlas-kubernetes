package plugins

import (
	"fmt"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type PrintConditions struct{}

func (p *PrintConditions) Name() string {
	return "print_conditions"
}

func (p *PrintConditions) Process(req *MappingProcessorRequest) error {
	index := findVersionIndex(req.CRD.Spec.Versions, req.CRD.APIVersion)
	if index == -1 {
		return fmt.Errorf("apiVersion %q not listed in spec", req.CRD.APIVersion)
	}
	req.CRD.Spec.Versions[index].AdditionalPrinterColumns = []apiextensions.CustomResourceColumnDefinition{
		{
			JSONPath: `.status.conditions[?(@.type=="Ready")].status`,
			Name:     "Ready",
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

func findVersionIndex(versions []apiextensions.CustomResourceDefinitionVersion, version string) int {
	if len(versions) == 1 && version == "" {
		return 0
	}
	for i, v := range versions {
		if v.Name == version {
			return i
		}
	}
	return -1
}
