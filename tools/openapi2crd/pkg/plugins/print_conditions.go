package plugins

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

type PrintConditions struct{}

func (p *PrintConditions) Name() string {
	return "print_conditions"
}

func (p *PrintConditions) Process(req *MappingProcessorRequest) error {
	// Given that new SDK versions are supported by different spec subfields, it
	// is unlikely there will be more than one supported Kubernetes version.
	//
	// The CRD validation fails if all versions use the same print columns, in
	// such case they have to be set at the top level to avoid errors.
	// Note that the produced YAML still places the print columns at the version 
	// level anyways.
	//
	// Check out the upstream code for more details:
	// https://github.com/kubernetes/apiextensions-apiserver/blob/a780e0393511161d7ef1e6466035181a4f84f347/pkg/apis/apiextensions/validation/validation.go#L438C5-L438C34
	// https://github.com/kubernetes/apiextensions-apiserver/blob/a780e0393511161d7ef1e6466035181a4f84f347/pkg/apis/apiextensions/validation/validation.go#L762
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
