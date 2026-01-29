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

const AtlasExportersGeneratorName = "atlas-exporters"

func init() {
	Register(&AtlasExportersGenerator{})
}

// AtlasExportersGenerator generates exporter files for CRDs.
type AtlasExportersGenerator struct{}

// Name returns the generator name.
func (g *AtlasExportersGenerator) Name() string {
	return AtlasExportersGeneratorName
}

// Description returns a human-readable description.
func (g *AtlasExportersGenerator) Description() string {
	return "Generates Atlas exporter files for importing resources from Atlas API"
}

// Generate runs the exporter generation for a single CRD kind.
func (g *AtlasExportersGenerator) Generate(opts *Options) error {
	parsedConfig, err := opts.GetParsedConfig()
	if err != nil {
		return err
	}

	resourceName := parsedConfig.ResourceName

	// Generate exporter for each SDK version mapping
	for _, mapping := range parsedConfig.Mappings {
		if err := generate.GenerateResourceExporter(
			opts.InputPath,
			opts.CRDKind,
			resourceName,
			opts.TypesPath,
			opts.ExporterOutDir,
			mapping,
		); err != nil {
			return fmt.Errorf("failed to generate exporter for resource %s version %s: %w", resourceName, mapping.Version, err)
		}
	}

	return nil
}
