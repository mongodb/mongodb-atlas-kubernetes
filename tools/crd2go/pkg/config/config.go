// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"io"
)

// CodeWriterFunc is a function type that takes a CRD and returns a writer for the generated code
type CodeWriterFunc func(filename string, overwrite bool) (io.WriteCloser, error)

// GenDeepCopy controls how deep copy generation is handled
type GenDeepCopy string

const (
	// GenDeepCopyAuto runs controller-gen when present in $PATH
	GenDeepCopyAuto = "auto"

	// GenDeepCopyOff will skip controller-gen
	GenDeepCopyOff = "off"

	// GenDeepCopyForced always runs controller-gen after CRD code generation
	GenDeepCopyForced = "forced"
)

// Config holds all CLI configurable parameters
type Config struct {
	CoreConfig `yaml:",inline"`

	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

// ImportedTypeConfig holds one imported type information
type ImportedTypeConfig struct {
	ImportInfo `yaml:",inline"`
	Name       string `yaml:"name"`
}

// ImportInfo holds the import path and alias for existing types
type ImportInfo struct {
	Alias string
	Path  string
}

// CoreConfig holds the subset of the config witout the input and output fields
type CoreConfig struct {
	Version  string               `yaml:"version"`
	Reserved []string             `yaml:"reserved"`
	SkipList []string             `yaml:"skipList"`
	Renames  map[string]string    `yaml:"renames"`
	Imports  []ImportedTypeConfig `yaml:"imports"`
	DeepCopy DeepCopy             `yaml:"deepCopy"`
}

type DeepCopy struct {
	Generate          GenDeepCopy `yaml:"generate"`
	ControllerGenPath string      `yaml:"controllerGenPath"`
}
