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

const IndexersGeneratorName = "indexers"

func init() {
	Register(&IndexersGenerator{})
}

// IndexersGenerator generates indexer files for CRDs.
type IndexersGenerator struct{}

// Name returns the generator name.
func (g *IndexersGenerator) Name() string {
	return IndexersGeneratorName
}

// Description returns a human-readable description.
func (g *IndexersGenerator) Description() string {
	return "Generates indexer files for CRD reference field lookups"
}

// Generate runs the indexer generation for a single CRD kind.
func (g *IndexersGenerator) Generate(opts *Options) error {
	if err := generate.GenerateIndexers(opts.InputPath, opts.CRDKind, opts.IndexerOutDir, opts.IndexerTypesPath); err != nil {
		return fmt.Errorf("failed to generate indexers: %w", err)
	}
	return nil
}
