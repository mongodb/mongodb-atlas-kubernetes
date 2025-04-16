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

package version

import (
	"regexp"
	"strings"
)

const DefaultVersion = "unknown"

// Version set by the linker during link time.
var Version = DefaultVersion

func IsRelease(v string) bool {
	return v != DefaultVersion &&
		regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+[-certified]*$`).Match([]byte(strings.TrimSpace(v)))
}
