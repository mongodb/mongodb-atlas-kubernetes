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

package refs

import (
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/pkg/crapi/objmap"
)

type OpenAPIMapping struct {
	Property string `json:"property"`
	Type     string `json:"type"`
}

func (oam OpenAPIMapping) targetPath() []string {
	return resolveXPath(oam.Property)
}

func resolveXPath(xpath string) []string {
	if strings.HasPrefix(xpath, "$.") {
		return objmap.AsPath(xpath[1:])
	}
	return objmap.AsPath(xpath)
}
