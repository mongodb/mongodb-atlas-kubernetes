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
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"
)

type Loader interface {
	Load(ctx context.Context, path string) (*openapi3.T, error)
}

type KinOpeAPI struct {
	fs afero.Fs
}

func NewKinOpeAPI(fs afero.Fs) *KinOpeAPI {
	return &KinOpeAPI{
		fs: fs,
	}
}

func (a *KinOpeAPI) Load(_ context.Context, path string) (*openapi3.T, error) {
	loader := &openapi3.Loader{
		IsExternalRefsAllowed: true,
	}

	if uri, ok := isURI(path); ok {
		return loader.LoadFromURI(uri)
	}

	data, err := a.transform(path)
	if err != nil {
		return nil, fmt.Errorf("failed to transform the file %s: %w", path, err)
	}

	return loader.LoadFromData(data)
}

func (a *KinOpeAPI) transform(path string) ([]byte, error) {
	filePath := filepath.Clean(path)

	data, err := afero.ReadFile(a.fs, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read the file %s: %w", filePath, err)
	}

	result := make(map[string]interface{})
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}

	removeXGenChangelog(result)

	data, err = yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return data, nil
}

func isURI(path string) (*url.URL, bool) {
	uri, err := url.Parse(path)

	if err == nil && uri.Scheme != "" && uri.Host != "" {
		return uri, true
	}

	return nil, false
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
