package plugins

import (
	"errors"
	"testing"
	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"

	"github.com/stretchr/testify/assert"
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
					Plugins:     []string{"reference", "read_only_property", "read_write_property"},
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
					Plugins:     []string{"reference", "read_only_property", "read_write_property"},
				},
				{
					Name:    "set-1",
					Default: true,
					Plugins: []string{"base", "entry", "status", "reference_metadata"},
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
					Plugins:     []string{"reference", "read_only_property", "read_write_property"},
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
					Plugins: []string{"reference", "read_only_property", "read_write_property"},
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
