package config

import (
	"io"
)

// CodeWriterFunc is a function type that takes a CRD and returns a writer for the generated code
type CodeWriterFunc func(filename string, overwrite bool) (io.WriteCloser, error)

// GenDeepCopy controls how deep copy generation is handled
type GenDeepCopy string

const (
	// GenDeepCopyAuto will run controller-gen if present in the path
	GenDeepCopyAuto = "auto"

	// GenDeepCopyOff will not try to controller-gen
	GenDeepCopyOff = "off"

	// GenDeepCopyForced will run controller-gen and fail if it fails in any way
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
