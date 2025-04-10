package crd2go

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	// ErrEmptyDoc fails when an non empty YAML object was expected
	ErrEmptyDoc = errors.New("empty document")

	// ErrNotAnObject is emited when the YAMl was suposed to contain an object
	ErrNotAnObject = errors.New("not a YAML object")
)

type CRD struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Metadata   Metadata `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type Metadata struct {
	Name string `yaml:"name"`
}

type Spec struct {
	Group    string    `yaml:"group"`
	Scope    string    `yaml:"scope"`
	Names    SpecNames `yaml:"names"`
	Versions []Version `yaml:"versions"`
}

type SpecNames struct {
	Kind     string `yaml:"kind"`
	ListKind string `yaml:"listKind"`
}

type Version struct {
	Name         string         `yaml:"name"`
	Schema       VersionSchema  `yaml:"schema"`
	Served       bool           `yaml:"served"`
	Storage      bool           `yaml:"storage"`
	Subresources map[string]any `yaml:"subresources"`
}

type VersionSchema struct {
	OpenAPIV3Schema OpenAPISchema `yaml:"openAPIV3Schema"`
}

type OpenAPISchema struct {
	Type        string                   `yaml:"type"`
	Properties  map[string]OpenAPISchema `yaml:"properties"`
	Items       *OpenAPISchema           `yaml:"items,omitempty"`
	Required    []string                 `yaml:"required,omitempty"`
	Title       *string                  `yaml:"title,omitempty"`
	Description *string                  `yaml:"description,omitempty"`
	MinLength   *int                     `yaml:"minLength,omitempty"`
	MaxLength   *int                     `yaml:"maxLength,omitempty"`
	Format      *string                  `yaml:"format,omitempty"`
}

func ParseCRD(r io.Reader) (*CRD, error) {
	crd := CRD{}
	err := yaml.NewDecoder(r).Decode(&crd)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CRD YAML: %w", err)
	}
	return &crd, nil
}
