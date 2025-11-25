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

package unstructured

import "strings"

// AsPath translates the given simplified xpath expression into a sequence of
// path entries. This ia very shallow xpath formatter, not full xpath compliant
func AsPath(xpath string) []string {
	if strings.HasPrefix(xpath, ".") {
		return AsPath(xpath[1:])
	}
	return strings.Split(xpath, ".")
}

// Base returns the base of the given path, namely the last name in the array
func Base(path []string) string {
	if len(path) == 0 {
		return ""
	}
	lastIndex := len(path) - 1
	return path[lastIndex]
}

// Dir returns the parent path of path
func Dir(path []string) []string {
	if len(path) == 0 {
		return path
	}
	return path[:len(path)-1]
}
