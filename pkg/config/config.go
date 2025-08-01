package config

import (
	"io"
)

// CodeWriterFunc is a function type that takes a CRD and returns a writer for the generated code
type CodeWriterFunc func(filename string, overwrite bool) (io.WriteCloser, error)

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
}
