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
	"net/url"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadOpenAPI(filePath string) (*openapi3.T, error) {
	loader := &openapi3.Loader{
		IsExternalRefsAllowed: true,
	}

	uri, err := url.Parse(filePath)
	if err == nil && uri.Scheme != "" && uri.Host != "" {
		return loader.LoadFromURI(uri)
	}

	return loader.LoadFromFile(filepath.Clean(filePath))
}
