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

package operator

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func repositoryDir() string {
	// Caller(0) returns the path to the calling test file rather than the path to this framework file. That
	// precludes assuming how many directories are between the file and the repo root. It's therefore necessary
	// to search in the hierarchy for an indication of a path that looks like the repo root.
	//nolint:dogsled
	_, sourceFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(sourceFile)
	for {
		// go.mod should always exist in the repo root
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		} else if errors.Is(err, os.ErrNotExist) {
			currentDir, err = filepath.Abs(filepath.Join(currentDir, ".."))
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return currentDir
}
