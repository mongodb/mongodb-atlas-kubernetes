package crd2go

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
)

func ParseCRD(r io.Reader) (*apiextensions.CustomResourceDefinition, error) {
	crd := apiextensions.CustomResourceDefinition{}

	yml, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read all YAML input: %w", err)
	}
	err = yaml.Unmarshal(yml, &crd)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CRD YAML: %w", err)
	}
	return &crd, nil
}
