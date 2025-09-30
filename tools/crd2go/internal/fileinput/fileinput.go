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

package fileinput

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MustBeSafe(filename string) string {
	cleanPath, err := Safe(filename)
	if err != nil {
		panic(err)
	}
	return cleanPath
}

func Safe(filename string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory to sanitize path %q: %w", filename, err)
	}
	return SafeAt(cwd, filename)
}

func MustBeSafeAt(allowedBase, filename string) string {
	cleanPath, err := SafeAt(allowedBase, filename)
	if err != nil {
		panic(err)
	}
	return cleanPath
}

func SafeAt(allowedBase, filename string) (string, error) {
	absAllowedBase, err := filepath.Abs(allowedBase)
	if err != nil {
		return "", fmt.Errorf("absolute base path computation failed: %w", err)
	}
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("absolute filename path computation failed: %w", err)
	}
	cleanPath := filepath.Clean(absPath)
	if !strings.HasPrefix(cleanPath, absAllowedBase) {
		return "", fmt.Errorf("Unsafe input path %q not in %q (clean path %q)", filename, allowedBase, cleanPath)
	}
	return cleanPath, nil
}
