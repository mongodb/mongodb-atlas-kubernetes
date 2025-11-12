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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

func TestBuildSets(t *testing.T) {
	tests := map[string]struct {
		setsDefinition []configv1alpha1.PluginSet
		expectedSet    []Set
		expectedErr    error
	}{
		"build sets": {
			setsDefinition: []configv1alpha1.PluginSet{
				{
					Name:    "set-1",
					Default: true,
					Plugins: []string{"base", "entry", "status"},
				},
				{
					Name:        "set-2",
					Default:     false,
					InheritFrom: "set-1",
					Plugins:     []string{"references", "read_only_properties", "read_write_properties"},
				},
			},
			expectedSet: []Set{
				{
					Name:      "set-1",
					Default:   true,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}},
					Property:  []PropertyPlugin{},
					Extension: []ExtensionPlugin{},
				},
				{
					Name:      "set-2",
					Default:   false,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
					Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
					Extension: []ExtensionPlugin{},
				},
			},
			expectedErr: nil,
		},
		"build sets in random order": {
			setsDefinition: []configv1alpha1.PluginSet{
				{
					Name:        "set-2",
					Default:     false,
					InheritFrom: "set-1",
					Plugins:     []string{"references", "read_only_properties", "read_write_properties"},
				},
				{
					Name:    "set-1",
					Default: true,
					Plugins: []string{"base", "entry", "status", "references_metadata"},
				},
			},
			expectedSet: []Set{
				{
					Name:      "set-1",
					Default:   true,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}},
					Property:  []PropertyPlugin{},
					Extension: []ExtensionPlugin{&ReferencesMetadata{}},
				},
				{
					Name:      "set-2",
					Default:   false,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
					Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
					Extension: []ExtensionPlugin{&ReferencesMetadata{}},
				},
			},
			expectedErr: nil,
		},
		"error on unknown plugin": {
			setsDefinition: []configv1alpha1.PluginSet{
				{
					Name:    "set-1",
					Default: true,
					Plugins: []string{"base", "entry", "unknown_plugin"},
				},
			},
			expectedSet: nil,
			expectedErr: errors.New("plugin unknown_plugin not found in catalog"),
		},
		"error on unknown inherited set": {
			setsDefinition: []configv1alpha1.PluginSet{
				{
					Name:        "set-2",
					Default:     false,
					InheritFrom: "set-1",
					Plugins:     []string{"references", "read_only_properties", "read_write_properties"},
				},
			},
			expectedSet: nil,
			expectedErr: errors.New("circular dependency detected for plugin set: set-2"),
		},
		"error on multiple default sets": {
			setsDefinition: []configv1alpha1.PluginSet{
				{
					Name:    "set-1",
					Default: true,
					Plugins: []string{"base", "entry", "status"},
				},
				{
					Name:    "set-2",
					Default: true,
					Plugins: []string{"references", "read_only_properties", "read_write_properties"},
				},
			},
			expectedSet: nil,
			expectedErr: errors.New("multiple default plugin sets defined"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewCatalog()
			set, err := c.BuildSets(tt.setsDefinition)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedSet, set)
		})
	}
}

func TestGetPluginSet(t *testing.T) {
	tests := map[string]struct {
		sets          []Set
		name          string
		expectedSet   *Set
		expectedError error
	}{
		"get existing set": {
			sets: []Set{
				{
					Name:      "set-1",
					Default:   true,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}},
					Property:  []PropertyPlugin{},
					Extension: []ExtensionPlugin{},
				},
				{
					Name:      "set-2",
					Default:   false,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
					Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
					Extension: []ExtensionPlugin{},
				},
			},
			name: "set-2",
			expectedSet: &Set{
				Name:      "set-2",
				Default:   false,
				CRD:       []CRDPlugin{&Base{}},
				Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
				Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
				Extension: []ExtensionPlugin{},
			},
			expectedError: nil,
		},
		"get default set": {
			sets: []Set{
				{
					Name:      "set-1",
					Default:   true,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}},
					Property:  []PropertyPlugin{},
					Extension: []ExtensionPlugin{},
				},
				{
					Name:      "set-2",
					Default:   false,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
					Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
					Extension: []ExtensionPlugin{},
				},
			},
			name: "",
			expectedSet: &Set{
				Name:      "set-1",
				Default:   true,
				CRD:       []CRDPlugin{&Base{}},
				Mapping:   []MappingPlugin{&Entry{}, &Status{}},
				Property:  []PropertyPlugin{},
				Extension: []ExtensionPlugin{},
			},
			expectedError: nil,
		},
		"error on unknown set": {
			sets: []Set{
				{
					Name:      "set-1",
					Default:   true,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}},
					Property:  []PropertyPlugin{},
					Extension: []ExtensionPlugin{},
				},
				{
					Name:      "set-2",
					Default:   false,
					CRD:       []CRDPlugin{&Base{}},
					Mapping:   []MappingPlugin{&Entry{}, &Status{}, &References{}},
					Property:  []PropertyPlugin{&ReadOnlyProperties{}, &ReadWriteProperties{}},
					Extension: []ExtensionPlugin{},
				},
			},
			name:          "set-3",
			expectedSet:   nil,
			expectedError: errors.New("pluginSet set-3 not found"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			set, err := GetPluginSet(tt.sets, tt.name)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedSet, set)
		})
	}
}
