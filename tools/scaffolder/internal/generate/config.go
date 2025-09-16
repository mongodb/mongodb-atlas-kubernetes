// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Config structures
type Config struct {
	metav1.TypeMeta `yaml:",inline"`
	Spec            Spec `yaml:"spec"`
}

type Spec struct {
	OpenAPI []OpenAPIConfig `yaml:"openapi,omitempty"`
	CRD     []CRDConfig     `yaml:"crd,omitempty"`
}

type OpenAPIConfig struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
}

type CRDConfig struct {
	GVK        metav1.GroupVersionKind `yaml:"gvk"`
	Categories []string                `yaml:"categories,omitempty"`
	ShortNames []string                `yaml:"shortNames,omitempty"`
	Mappings   []Mapping               `yaml:"mappings,omitempty"`
}

type Mapping struct {
	MajorVersion string      `yaml:"majorVersion"`
	OpenAPIRef   OpenAPIRef  `yaml:"openAPIRef"`
	Parameters   *Parameters `yaml:"parameters,omitempty"`
	Entry        *SchemaRef  `yaml:"entry,omitempty"`
	Status       *SchemaRef  `yaml:"status,omitempty"`
}

type OpenAPIRef struct {
	Name string `yaml:"name"`
}

type Parameters struct {
	Path       *PathInfo     `yaml:"path,omitempty"`
	Query      []QueryParam  `yaml:"query,omitempty"`
	Additional []interface{} `yaml:"additional,omitempty"`
}

type PathInfo struct {
	Template string `yaml:"template"`
}

type QueryParam struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type SchemaRef struct {
	Schema string `yaml:"$ref"`
}

// MappingWithConfig combines mapping with its OpenAPI config
type MappingWithConfig struct {
	Mapping       Mapping
	OpenAPIConfig OpenAPIConfig
}

// ParsedConfig contains all parsed configuration data
type ParsedConfig struct {
	Config       Config
	OpenAPIMap   map[string]OpenAPIConfig
	CRDMap       map[string]CRDConfig
	SelectedCRD  CRDConfig
	Mappings     []MappingWithConfig
	ResourceName string
}

// ParseAtlas2CRDConfig reads and parses the configuration file, validates CRD selection,
// and returns all necessary data for generating controllers and handlers
func ParseAtlas2CRDConfig(configPath, crdKind string) (*ParsedConfig, error) {
	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Create OpenAPI mapping
	openAPIMap := make(map[string]OpenAPIConfig)
	for _, openAPIConfig := range config.Spec.OpenAPI {
		openAPIMap[openAPIConfig.Name] = openAPIConfig
	}

	// Create CRD mapping and find selected CRD
	crdMap := make(map[string]CRDConfig)
	var selectedCRD CRDConfig
	var found bool

	for _, crd := range config.Spec.CRD {
		crdMap[crd.GVK.Kind] = crd
		if crd.GVK.Kind == crdKind {
			selectedCRD = crd
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("CRD kind '%s' not found in config", crdKind)
	}

	if len(selectedCRD.Mappings) == 0 {
		return nil, fmt.Errorf("no mappings found for CRD kind '%s'", crdKind)
	}

	// Build mappings with their OpenAPI configs
	var mappingsWithConfig []MappingWithConfig
	for _, mapping := range selectedCRD.Mappings {
		openAPIConfig, exists := openAPIMap[mapping.OpenAPIRef.Name]
		if !exists {
			return nil, fmt.Errorf("OpenAPI config '%s' not found for mapping", mapping.OpenAPIRef.Name)
		}

		mappingsWithConfig = append(mappingsWithConfig, MappingWithConfig{
			Mapping:       mapping,
			OpenAPIConfig: openAPIConfig,
		})
	}

	return &ParsedConfig{
		Config:       config,
		OpenAPIMap:   openAPIMap,
		CRDMap:       crdMap,
		SelectedCRD:  selectedCRD,
		Mappings:     mappingsWithConfig,
		ResourceName: crdKind,
	}, nil
}

// ListCRDs returns a list of all available CRDs from the config file
func ListCRDs(configPath string) ([]CRDConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return config.Spec.CRD, nil
}

