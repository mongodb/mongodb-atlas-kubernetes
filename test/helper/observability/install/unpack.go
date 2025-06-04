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
			return os.MkdirAll(targetPath, 0755)
		}

		// Read the file content
		data, err := root.Assets.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read embedded file %s: %w", path, err)
		}

		// Write to the temp directory
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
