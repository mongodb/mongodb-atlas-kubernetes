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

package install

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	root "github.com/mongodb/mongodb-atlas-kubernetes/v2"
)

// Unpack unpacks assets to a temporary directory and returns the path to that directory or an error if it fails.
// use `defer os.RemoveAll(...)` to clean up the temporary directory after use.
func Unpack() (string, error) {
	tempDir, err := os.MkdirTemp("", "unpacked-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	err = fs.WalkDir(root.Assets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(tempDir, path)

		if d.IsDir() {
			//nolint:gosec
			return os.MkdirAll(targetPath, 0755)
		}

		// Read the file content
		data, err := root.Assets.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write to the temp directory
		//nolint:gosec
		err = os.WriteFile(targetPath, data, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		return nil
	})

	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to unpack files: %w", err)
	}

	return tempDir, nil
}
