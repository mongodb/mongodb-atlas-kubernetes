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

package plugins

import (
	"slices"

	"k8s.io/apimachinery/pkg/util/sets"
)

type ReadOnlyProperties struct{}

func (p *ReadOnlyProperties) Name() string {
	return "read_only_properties"
}

func (p *ReadOnlyProperties) Process(req *PropertyProcessorRequest) error {
	if req.PropertyConfig == nil || !req.PropertyConfig.Filters.ReadOnly {
		return nil
	}

	if req.OpenAPISchema.ReadOnly {
		return nil
	}

	required := sets.New(req.OpenAPISchema.Required...)
	for name, p := range req.OpenAPISchema.Properties {
		if !p.Value.ReadOnly {
			required.Delete(name)
		}
	}
	req.Property.Required = required.UnsortedList()
	slices.Sort(req.Property.Required)

	// ignore root
	if len(req.Path) == 1 {
		return nil
	}

	req.Property = nil

	return nil
}
