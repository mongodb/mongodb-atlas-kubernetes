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
	"fmt"
	"net/url"
	"path/filepath"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/goccy/go-yaml"
	"github.com/spf13/afero"
	"golang.org/x/sync/singleflight"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/flatten"
)

// Loader loads an OpenAPI spec based on the provided definition.
type Loader interface {
	Load(ctx context.Context, def v1alpha1.OpenAPIDefinition) (*openapi3.T, error)
}

type KinOpenAPI struct {
	fs    afero.Fs
	mu    sync.Mutex
	cache map[string]*openapi3.T
	group singleflight.Group
}

func NewKinOpenAPI(fs afero.Fs) *KinOpenAPI {
	return &KinOpenAPI{
		fs:    fs,
		cache: make(map[string]*openapi3.T),
	}
}

func (a *KinOpenAPI) load(_ context.Context, path string) (*openapi3.T, error) {
	// Fast path: return the cached spec without entering singleflight.
	// The mutex-guarded cache avoids the overhead of singleflight.Do on
	// every call after the spec has already been parsed.
	a.mu.Lock()
	if spec, ok := a.cache[path]; ok {
		a.mu.Unlock()
		return spec, nil
	}
	a.mu.Unlock()

	// Slow path: singleflight deduplicates concurrent parses of the same spec.
	// Returns interface{} because singleflight has no generic API.
	v, err, _ := a.group.Do(path, func() (interface{}, error) {
		loader := &openapi3.Loader{
			IsExternalRefsAllowed: true,
		}

		var (
			spec *openapi3.T
			err  error
		)

		if uri, ok := isURI(path); ok {
			spec, err = loader.LoadFromURI(uri)
		} else {
			var data []byte
			data, err = a.transform(path)
			if err != nil {
				return nil, fmt.Errorf("failed to transform the file %s: %w", path, err)
			}
			spec, err = loader.LoadFromData(data)
		}

		if err != nil {
			return nil, err
		}

		a.mu.Lock()
		a.cache[path] = spec
		a.mu.Unlock()

		return spec, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*openapi3.T), nil //nolint:forcetypeassert // singleflight returns interface{}; type is guaranteed by the closure above.
}

// LoadFlattened reads an OpenAPI spec from path, applies the same transform
// as Load (strip x-xgen-changelog), then flattens schema compositions
// (oneOf/anyOf/allOf/discriminator) before parsing with kin-openapi.
func (a *KinOpenAPI) loadFlattened(_ context.Context, path string) (*openapi3.T, error) {
	cacheKey := "flatten:" + path

	a.mu.Lock()
	if spec, ok := a.cache[cacheKey]; ok {
		a.mu.Unlock()
		return spec, nil
	}
	a.mu.Unlock()

	v, err, _ := a.group.Do(cacheKey, func() (interface{}, error) {
		data, err := a.transform(path)
		if err != nil {
			return nil, fmt.Errorf("failed to transform the file %s: %w", path, err)
		}

		data, err = flatten.Flatten(data)
		if err != nil {
			return nil, fmt.Errorf("failed to flatten the file %s: %w", path, err)
		}

		loader := &openapi3.Loader{
			IsExternalRefsAllowed: true,
		}
		spec, err := loader.LoadFromData(data)
		if err != nil {
			return nil, err
		}

		a.mu.Lock()
		a.cache[cacheKey] = spec
		a.mu.Unlock()

		return spec, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*openapi3.T), nil
}

func (a *KinOpenAPI) transform(path string) ([]byte, error) {
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
