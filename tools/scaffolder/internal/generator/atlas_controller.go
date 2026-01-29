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

package generator

import (
	"fmt"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generate"
)

const AtlasControllersGeneratorName = "atlas-controllers"

func init() {
	Register(&AtlasControllersGenerator{})
}

// AtlasControllersGenerator generates controller files for CRDs.
type AtlasControllersGenerator struct{}

// Name returns the generator name.
func (g *AtlasControllersGenerator) Name() string {
	return AtlasControllersGeneratorName
}

// Description returns a human-readable description.
func (g *AtlasControllersGenerator) Description() string {
	return "Generates Atlas controller and handler files for CRDs"
}

// Generate runs the controller generation for a single CRD kind.
func (g *AtlasControllersGenerator) Generate(opts *Options) error {
	parsedConfig, err := opts.GetParsedConfig()
	if err != nil {
		return err
	}

	resourceName := parsedConfig.ResourceName

	// Parse reference fields for watch generation
	referenceFields, err := generate.ParseReferenceFields(opts.InputPath, opts.CRDKind)
	if err != nil {
		return fmt.Errorf("failed to parse reference fields: %w", err)
	}

	// Group references by target kind
	refsByKind := make(map[string][]generate.ReferenceField)
	for _, ref := range referenceFields {
		refsByKind[ref.ReferencedKind] = append(refsByKind[ref.ReferencedKind], ref)
	}

	// Generate the controller using the generate package
	if err := generate.GenerateController(
		opts.InputPath,
		opts.CRDKind,
		opts.ControllerOutDir,
		opts.TypesPath,
		opts.IndexerImportPath,
		parsedConfig,
		refsByKind,
		opts.Override,
	); err != nil {
		return fmt.Errorf("failed to generate controller: %w", err)
	}

	fmt.Printf("Successfully generated controller for resource %s with %d SDK versions\n",
		resourceName, len(parsedConfig.Mappings))

	return nil
}
