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

package config

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

// CompositeLoader implements Loader by routing to the appropriate internal
// loader based on the OpenAPIDefinition fields.
type CompositeLoader struct {
	file  *KinOpenAPI
	atlas *Atlas
}

func NewLoader(file *KinOpenAPI, atlas *Atlas) *CompositeLoader {
	return &CompositeLoader{file: file, atlas: atlas}
}

func (c *CompositeLoader) Load(ctx context.Context, def v1alpha1.OpenAPIDefinition) (*openapi3.T, error) {
	if def.Package != "" {
		return c.loadFromPackage(ctx, def)
	}

	return c.loadFromPath(ctx, def)
}

func (c *CompositeLoader) loadFromPackage(ctx context.Context, def v1alpha1.OpenAPIDefinition) (*openapi3.T, error) {
	if def.Flatten {
		spec, err := c.atlas.loadFlattenedFromPackage(ctx, def.Package, def.Path)
		if err != nil {
			return nil, fmt.Errorf("error loading and flattening Atlas OpenAPI package %q: %w", def.Package, err)
		}
		return spec, nil
	}

	spec, err := c.atlas.loadFromPackage(ctx, def.Package, def.Path)
	if err != nil {
		return nil, fmt.Errorf("error loading Atlas OpenAPI package %q: %w", def.Package, err)
	}
	return spec, nil
}

func (c *CompositeLoader) loadFromPath(ctx context.Context, def v1alpha1.OpenAPIDefinition) (*openapi3.T, error) {
	if def.Flatten {
		spec, err := c.file.loadFlattened(ctx, def.Path)
		if err != nil {
			return nil, fmt.Errorf("error loading and flattening spec %q: %w", def.Path, err)
		}
		return spec, nil
	}

	spec, err := c.file.load(ctx, def.Path)
	if err != nil {
		return nil, fmt.Errorf("error loading spec: %w", err)
	}
	return spec, nil
}
