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

// Package registry provides the generator registry and shared types.
package registry

import (
	"fmt"
	"sort"
	"sync"
)

// Generator is the interface that all generators must implement.
// Each generator is responsible for generating a specific type of output
// (controllers, exporters, indexers, etc.).
type Generator interface {
	// Name returns the unique name of this generator (e.g., "atlas-controllers", "indexers")
	Name() string

	// Description returns a human-readable description of what this generator does
	Description() string

	// Generate runs the generator for a single CRD kind
	Generate(opts *Options) error
}

// Options contains all the configuration needed by generators.
// This is passed to each generator's Generate method.
type Options struct {
	// InputPath is the path to the CRD YAML file
	InputPath string

	// CRDKind is the specific CRD kind to generate for
	CRDKind string

	// ControllerOutDir is the output directory for controller files
	ControllerOutDir string

	// IndexerOutDir is the output directory for indexer files
	IndexerOutDir string

	// ExporterOutDir is the output directory for exporter files
	ExporterOutDir string

	// TypesPath is the full import path to the API types package
	TypesPath string

	// IndexerTypesPath is the full import path for type imports in indexers
	IndexerTypesPath string

	// IndexerImportPath is the full import path for indexer imports in controllers
	IndexerImportPath string

	// Override determines whether to override existing versioned handler files
	Override bool
}

// registry holds all registered generators
var (
	registryMu sync.RWMutex
	registry   = make(map[string]Generator)
)

// Register adds a generator to the registry.
// This is typically called from init() functions in generator packages.
// Panics if a generator with the same name is already registered.
func Register(g Generator) {
	registryMu.Lock()
	defer registryMu.Unlock()

	name := g.Name()
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("generator %q already registered", name))
	}
	registry[name] = g
}

// List returns all registered generator names in sorted order.
func List() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// All returns all registered generators.
func All() []Generator {
	registryMu.RLock()
	defer registryMu.RUnlock()

	generators := make([]Generator, 0, len(registry))
	for _, g := range registry {
		generators = append(generators, g)
	}
	return generators
}

// GetByNames returns generators matching the given names.
// Returns an error if any name is not found.
func GetByNames(names []string) ([]Generator, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()

	generators := make([]Generator, 0, len(names))
	for _, name := range names {
		g, exists := registry[name]
		if !exists {
			return nil, fmt.Errorf("unknown generator: %q (available: %v)", name, List())
		}
		generators = append(generators, g)
	}
	return generators, nil
}
