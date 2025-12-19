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

package samples

import (
	"os"
	"path/filepath"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/test/helper/e2e2/yml"
)

// findRepoRoot finds the repository root by looking for go.mod file.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	// Fallback: try current directory
	return ".", nil
}

// LoadSampleObjects loads and parses YAML objects from config/samples.
// It finds the repository root and loads files from there.
func LoadSampleObjects(filename string) ([]client.Object, error) {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return nil, err
	}

	absPath := filepath.Join(repoRoot, "config", "samples", filename)
	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return yml.ParseObjects(f)
}

// MustLoadSampleObjects loads and parses YAML objects, panicking on error.
func MustLoadSampleObjects(filename string) []client.Object {
	objs, err := LoadSampleObjects(filename)
	if err != nil {
		panic(err)
	}
	return objs
}
