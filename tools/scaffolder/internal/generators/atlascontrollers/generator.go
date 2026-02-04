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

// Package atlascontrollers provides the controller generator for Atlas CRDs.
package atlascontrollers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/indexers"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/registry"
)

// GeneratorName is the unique name for this generator.
const GeneratorName = "atlas-controllers"

func init() {
	registry.Register(&Generator{})
}

// Generator generates controller files for CRDs.
type Generator struct{}

// Name returns the generator name.
func (g *Generator) Name() string {
	return GeneratorName
}

// Description returns a human-readable description.
func (g *Generator) Description() string {
	return "Generates Atlas controller and handler files for CRDs"
}

// Generate runs the controller generation for a single CRD kind.
func (g *Generator) Generate(opts *registry.Options) error {
	parsedConfig, err := config.ParseCRDConfig(opts.InputPath, opts.CRDKind)
	if err != nil {
		return err
	}

	resourceName := parsedConfig.ResourceName

	// Parse reference fields for watch generation
	referenceFields, err := indexers.ParseReferenceFields(opts.InputPath, opts.CRDKind)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	// Group references by target kind
	refsByKind := make(map[string][]indexers.ReferenceField)
	for _, ref := range referenceFields {
		refsByKind[ref.ReferencedKind] = append(refsByKind[ref.ReferencedKind], ref)
	}

	// Set default directory if not provided
	controllerOutDir := opts.ControllerOutDir
	if controllerOutDir == "" {
		controllerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "controller")
	}

	controllerDir := filepath.Join(controllerOutDir, strings.ToLower(resourceName))

	if err := os.MkdirAll(controllerDir, 0755); err != nil {
		return fmt.Errorf("failed to create controller directory: %w", err)
	}

	if err := generateControllerFile(controllerDir, resourceName, opts.TypesPath, parsedConfig.Mappings); err != nil {
		return fmt.Errorf("failed to generate controller file: %w", err)
	}

	if err := generateMainHandlerFile(controllerDir, resourceName, opts.TypesPath, opts.IndexerImportPath, parsedConfig.Mappings, refsByKind, parsedConfig); err != nil {
		return fmt.Errorf("failed to generate main handler file: %w", err)
	}

	// Generate version-specific handlers
	for _, mapping := range parsedConfig.Mappings {
		if err := generateVersionHandlerFile(controllerDir, resourceName, opts.TypesPath, opts.IndexerImportPath, opts.InputPath, mapping, opts.Override); err != nil {
			return fmt.Errorf("failed to generate handler for version %s: %w", mapping.Version, err)
		}
	}

	fmt.Printf("Successfully generated controller for resource %s with %d SDK versions at %s\n",
		resourceName, len(parsedConfig.Mappings), controllerDir)

	return nil
}
