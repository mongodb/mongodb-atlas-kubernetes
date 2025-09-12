/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-yaml"
)

func LoadOpenAPI(filePath string) (*openapi3.T, error) {
	loader := &openapi3.Loader{
		IsExternalRefsAllowed: true,
	}

	uri, err := url.Parse(filePath)
	if err == nil && uri.Scheme != "" && uri.Host != "" {
		return loader.LoadFromURI(uri)
	}

	filePath = filepath.Clean(filePath)
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	result := make(map[string]interface{})
	err = yaml.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}
	removeXGenChangelog(result)
	b, err = yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return loader.LoadFromData(b)
}

func removeXGenChangelog(m map[string]interface{}) {
	for key, val := range m {
		if key == "x-xgen-changelog" {
			delete(m, key)
			continue
		}
		switch v := val.(type) {
		case map[string]interface{}:
			removeXGenChangelog(v)
		}
	}
}
