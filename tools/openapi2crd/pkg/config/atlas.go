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
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"golang.org/x/sync/singleflight"
)

type Atlas struct {
	fileLoader Loader
	mu         sync.Mutex
	pathCache  map[string]string
	group      singleflight.Group
}

func (a *Atlas) Load(ctx context.Context, pkg string) (*openapi3.T, error) {
	a.mu.Lock()
	cachedPath, ok := a.pathCache[pkg]
	a.mu.Unlock()

	var filename string
	if ok {
		filename = cachedPath
	} else {
		v, err, _ := a.group.Do(pkg, func() (interface{}, error) {
			path, err := getGoModulePath(ctx, pkg)
			if err != nil {
				return nil, fmt.Errorf("failed to load module path: %w", err)
			}

			resolved := filepath.Clean(filepath.Join(path, "..", "openapi", "atlas-api-transformed.yaml"))

			a.mu.Lock()
			a.pathCache[pkg] = resolved
			a.mu.Unlock()

			return resolved, nil
		})
		if err != nil {
			return nil, err
		}
		filename = v.(string)
	}

	return a.fileLoader.Load(ctx, filename)
}

func NewAtlas(loader Loader) *Atlas {
	return &Atlas{
		fileLoader: loader,
		pathCache:  make(map[string]string),
	}
}

func getGoModulePath(ctx context.Context, modulePath string) (string, error) {
	goCmd, err := exec.LookPath("go")
	if err != nil {
		return "", fmt.Errorf("go command not found in PATH: %w", err)
	}

	cmd := exec.CommandContext(ctx, goCmd, "list", "-f", "{{.Dir}}", modulePath)
	output, err := cmd.Output()
	if err != nil {
		// Check if the error is due to the module not being found
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "not a known module") || strings.Contains(stderr, "cannot find module") {
				return "", fmt.Errorf("module '%s' not found or not a dependency of the current project", modulePath)
			}
		}
		return "", fmt.Errorf("failed to run 'go list' for module '%s': %w, stderr: %s", modulePath, err, string(output))
	}

	return filepath.Clean(string(output)), nil
}
