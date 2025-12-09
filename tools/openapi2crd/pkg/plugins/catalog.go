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
	"fmt"

	configv1alpha1 "github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
)

type Set struct {
	Name      string
	Default   bool
	CRD       []CRDPlugin
	Mapping   []MappingPlugin
	Property  []PropertyPlugin
	Extension []ExtensionPlugin
}

type Catalog struct {
	crd       map[string]CRDPlugin
	mapping   map[string]MappingPlugin
	property  map[string]PropertyPlugin
	extension map[string]ExtensionPlugin
}

func (c *Catalog) IsCRD(name string) bool {
	_, ok := c.crd[name]

	return ok
}

func (c *Catalog) IsMapping(name string) bool {
	_, ok := c.mapping[name]

	return ok
}

func (c *Catalog) IsProperty(name string) bool {
	_, ok := c.property[name]

	return ok
}

func (c *Catalog) IsMappingExtension(name string) bool {
	_, ok := c.extension[name]

	return ok
}

func (c *Catalog) BuildSets(setsDefinition []configv1alpha1.PluginSet) ([]Set, error) {
	orderedSet, err := orderPluginSet(setsDefinition)
	if err != nil {
		return nil, err
	}

	sets := make([]Set, 0, len(orderedSet))
	hasDefault := false

	for _, pluginSet := range orderedSet {
		if pluginSet.Default {
			if hasDefault {
				return nil, errors.New("multiple default plugin sets defined")
			}

			hasDefault = true
		}

		var set Set

		if pluginSet.InheritFrom != "" {
			var parent *Set
			for _, s := range sets {
				if s.Name == pluginSet.InheritFrom {
					parent = &s
					break
				}
			}

			if parent == nil {
				return nil, fmt.Errorf("parent plugin set %s not found for plugin set %s", pluginSet.InheritFrom, pluginSet.Name)
			}

			set = *parent
			set.Name = pluginSet.Name
			set.Default = pluginSet.Default
		} else {
			set = Set{
				Name:      pluginSet.Name,
				Default:   pluginSet.Default,
				CRD:       make([]CRDPlugin, 0, len(c.crd)),
				Mapping:   make([]MappingPlugin, 0, len(c.mapping)),
				Property:  make([]PropertyPlugin, 0, len(c.property)),
				Extension: make([]ExtensionPlugin, 0, len(c.extension)),
			}
		}
		for _, plugin := range pluginSet.Plugins {
			switch {
			case c.IsCRD(plugin):
				set.CRD = append(set.CRD, c.crd[plugin])
			case c.IsMapping(plugin):
				set.Mapping = append(set.Mapping, c.mapping[plugin])
			case c.IsProperty(plugin):
				set.Property = append(set.Property, c.property[plugin])
			case c.IsMappingExtension(plugin):
				set.Extension = append(set.Extension, c.extension[plugin])
			default:
				return nil, fmt.Errorf("plugin %s not found in catalog", plugin)
			}
		}

		sets = append(sets, set)
	}

	return sets, nil
}

func NewCatalog() *Catalog {
	return &Catalog{
		crd: map[string]CRDPlugin{
			"base":                            &Base{},
			"mutual_exclusive_major_versions": &MutualExclusiveMajorVersions{},
		},
		mapping: map[string]MappingPlugin{
			"major_version":          &MajorVersion{},
			"parameters":             &Parameters{},
			"entry":                  &Entry{},
			"status":                 &Status{},
			"references":             &References{},
			"connection_secret":      &ConnectionSecret{},
			"mutual_exclusive_group": &MutualExclusiveGroup{},
		},
		property: map[string]PropertyPlugin{
			"sensitive_properties":  &SensitiveProperties{},
			"skipped_properties":    &SkippedProperties{},
			"read_only_properties":  &ReadOnlyProperties{},
			"read_write_properties": &ReadWriteProperties{},
		},
		extension: map[string]ExtensionPlugin{
			"atlas_sdk_version":    &AtlasSdkVersionPlugin{},
			"reference_extensions": &ReferenceExtensions{},
		},
	}
}

func GetPluginSet(sets []Set, name string) (*Set, error) {
	for _, set := range sets {
		if name == "" && set.Default {
			return &set, nil
		}

		if set.Name == name {
			return &set, nil
		}
	}

	return nil, fmt.Errorf("pluginSet %s not found", name)
}

func orderPluginSet(setsDefinition []configv1alpha1.PluginSet) ([]configv1alpha1.PluginSet, error) {
	mapSets := make(map[string]configv1alpha1.PluginSet)
	orderedSets := make([]configv1alpha1.PluginSet, 0, len(setsDefinition))
	visitCount := make(map[string]int)

	for len(setsDefinition) > 0 {
		set := setsDefinition[0]
		setsDefinition = setsDefinition[1:]

		hasDependency := set.InheritFrom != ""
		_, dependencyIsMapped := mapSets[set.InheritFrom]

		if hasDependency && !dependencyIsMapped {
			visitCount[set.Name]++
			if visitCount[set.Name] > len(setsDefinition)+1 {
				return nil, fmt.Errorf("circular dependency detected for plugin set: %s", set.Name)
			}

			setsDefinition = append(setsDefinition, set)
			continue
		}

		orderedSets = append(orderedSets, set)
		mapSets[set.Name] = set
	}

	return orderedSets, nil
}
