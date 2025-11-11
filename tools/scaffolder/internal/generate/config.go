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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type OpenAPIConfig struct {
	Name    string `yaml:"name"`
	Package string `yaml:"package"`
}

type MappingWithConfig struct {
	Version       string
	OpenAPIConfig OpenAPIConfig
}

type ParsedConfig struct {
	SelectedCRD    CRDInfo
	Mappings       []MappingWithConfig
	ResourceName   string
	APIVersion     string // API version package (e.g., "v1", "v3")
	StorageVersion string
	PluralName     string
	CRDGroup       string
}

type CRDDocument struct {
	APIVersion string                          `yaml:"apiVersion"`
	Kind       string                          `yaml:"kind"`
	Metadata   CRDMetadata                     `yaml:"metadata"`
	Spec       v1.CustomResourceDefinitionSpec `yaml:"spec"`
}

type CRDMetadata struct {
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations"`
}

type APIMapping struct {
	Properties map[string]SpecProperty `yaml:"properties"`
}

type SpecProperty struct {
	Properties map[string]VersionProperty `yaml:"properties"`
}

type VersionProperty struct {
	AtlasSDKVersion string `yaml:"x-atlas-sdk-version"`
}

type CRDInfo struct {
	Kind       string
	Group      string
	Version    string
	Plural     string
	ShortNames []string
	Categories []string
	Versions   []CRDVersionInfo
}

type CRDVersionInfo struct {
	Version         string
	AtlasSDKVersion string
	SDKPackage      string
}

func ParseCRDResult(resultPath, crdKind string) (*CRDInfo, error) {
	file, err := os.Open(resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open result file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for {
		crdInfo, err := parseNextCRD(scanner)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if crdInfo != nil && crdInfo.Kind == crdKind {
			return crdInfo, nil
		}
	}

	return nil, fmt.Errorf("CRD kind '%s' not found in result file", crdKind)
}

func parseNextCRD(scanner *bufio.Scanner) (*CRDInfo, error) {
	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if len(strings.TrimSpace(buffer.String())) > 0 {
				crdInfo, err := decodeCRDDocument(buffer.Bytes())
				if err != nil {
					buffer.Reset()
					continue
				}
				return crdInfo, nil
			}
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		buffer.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	if buffer.Len() > 0 {
		crdInfo, err := decodeCRDDocument(buffer.Bytes())
		if err != nil {
			return nil, err
		}
		return crdInfo, nil
	}

	return nil, io.EOF
}

func decodeCRDDocument(content []byte) (*CRDInfo, error) {
	var crd CRDDocument
	if err := yaml.Unmarshal(content, &crd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if crd.Kind != "CustomResourceDefinition" {
		return nil, fmt.Errorf("not a CustomResourceDefinition")
	}

	crdInfo := &CRDInfo{
		Kind:       crd.Spec.Names.Kind,
		Group:      crd.Spec.Group,
		Plural:     crd.Spec.Names.Plural,
		ShortNames: crd.Spec.Names.ShortNames,
		Categories: crd.Spec.Names.Categories,
	}

	if apiMappingsStr, exists := crd.Metadata.Annotations["api-mappings"]; exists {
		var apiMapping APIMapping
		if err := yaml.Unmarshal([]byte(apiMappingsStr), &apiMapping); err != nil {
			return nil, fmt.Errorf("failed to parse api-mappings annotation: %w", err)
		}

		if specProp, exists := apiMapping.Properties["spec"]; exists {
			for versionName, versionProp := range specProp.Properties {
				if versionProp.AtlasSDKVersion != "" {
					crdInfo.Versions = append(crdInfo.Versions, CRDVersionInfo{
						Version:         versionName,
						AtlasSDKVersion: versionProp.AtlasSDKVersion,
						SDKPackage:      versionProp.AtlasSDKVersion,
					})
				}
			}
		}
	}

	if len(crd.Spec.Versions) > 0 {
		crdInfo.Version = crd.Spec.Versions[0].Name
	}

	return crdInfo, nil
}

func ParseCRDConfig(resultPath, crdKind string) (*ParsedConfig, error) {
	crdInfo, err := ParseCRDResult(resultPath, crdKind)
	if err != nil {
		return nil, err
	}

	var mappings []MappingWithConfig
	for _, version := range crdInfo.Versions {
		mappings = append(mappings, MappingWithConfig{
			Version: version.Version,
			OpenAPIConfig: OpenAPIConfig{
				Name:    version.Version,
				Package: version.SDKPackage,
			},
		})
	}

	const apiVersion = "v3"
	// TODO: consider parsing target package from the api version
	// if crdInfo.Version != "" && crdInfo.Version != "v1" {
	// 	apiVersion = "v3"
	// }

	// Extract storage version - default to "v1" if not found
	storageVersion := "v1"
	if crdInfo.Version != "" {
		storageVersion = crdInfo.Version
	}

	return &ParsedConfig{
		SelectedCRD:    *crdInfo,
		Mappings:       mappings,
		ResourceName:   crdInfo.Kind,
		APIVersion:     apiVersion,
		StorageVersion: storageVersion,
		PluralName:     crdInfo.Plural,
		CRDGroup:       crdInfo.Group,
	}, nil
}

func ListCRDs(resultPath string) ([]CRDInfo, error) {
	file, err := os.Open(resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open result file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close result file: %v\n", err)
		}
	}()

	var crds []CRDInfo
	scanner := bufio.NewScanner(file)

	for {
		crdInfo, err := parseNextCRD(scanner)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if crdInfo != nil {
			crds = append(crds, *crdInfo)
		}
	}

	if len(crds) == 0 {
		return nil, fmt.Errorf("no CustomResourceDefinition documents found in '%s' file", resultPath)
	}

	return crds, nil
}
