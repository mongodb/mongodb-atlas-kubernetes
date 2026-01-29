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

import "github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generate"

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

	// ParsedConfig contains the parsed CRD configuration (lazily loaded)
	ParsedConfig *generate.ParsedConfig
}

// GetParsedConfig returns the parsed CRD config, loading it if necessary.
func (o *Options) GetParsedConfig() (*generate.ParsedConfig, error) {
	if o.ParsedConfig != nil {
		return o.ParsedConfig, nil
	}

	config, err := generate.ParseCRDConfig(o.InputPath, o.CRDKind)
	if err != nil {
		return nil, err
	}
	o.ParsedConfig = config
	return config, nil
}
