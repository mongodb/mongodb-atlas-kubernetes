package plugins

import (
	"errors"
	"fmt"

	configv1alpha1 "github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
)

type PluginSet struct {
	Name    string
	Default bool
	Plugins map[string]Plugin
}

type PluginSets []PluginSet

func (p *PluginSets) Get(name string) (*PluginSet, error) {
	for _, set := range *p {
		if set.Name == name {
			return &set, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("PluginSet %s not found", name))
}

func (p *PluginSets) Default() (*PluginSet, error) {
	for _, set := range *p {
		if set.Default {
			return &set, nil
		}
	}

	return nil, errors.New("no default plugin set defined")
}

func NewPluginSet(pluginSetDefinition []configv1alpha1.PluginSet, catalog PluginCatalog) ([]PluginSet, error) {
	orderedSet, err := orderPluginSet(pluginSetDefinition)
	if err != nil {
		return nil, err
	}

	sets := make([]PluginSet, 0, len(orderedSet))
	hasDefault := false

	for _, pluginSet := range orderedSet {
		if pluginSet.Default {
			if hasDefault {
				return nil, errors.New("multiple default plugin sets defined")
			}

			hasDefault = true
		}

		set := PluginSet{
			Name:    pluginSet.Name,
			Default: pluginSet.Default,
			Plugins: make(map[string]Plugin, len(pluginSet.Plugins)),
		}
		for _, plugin := range pluginSet.Plugins {
			p, err := catalog.Get(plugin)
			if err != nil {
				return nil, fmt.Errorf("failed to build plugin set %s: error getting plugin %s: %w", pluginSet.Name, plugin, err)
			}

			set.Plugins[p.Name()] = p
		}

		sets = append(sets, set)
	}

	if !hasDefault && len(sets) > 0 {
		sets[0].Default = true
	}

	return sets, nil
}

func orderPluginSet(sets []configv1alpha1.PluginSet) ([]configv1alpha1.PluginSet, error) {
	mapSets := make(map[string]configv1alpha1.PluginSet)
	orderedSets := make([]configv1alpha1.PluginSet, 0, len(sets))
	visitCount := make(map[string]int)

	for len(sets) > 0 {
		set := sets[0]
		sets = sets[1:]

		hasDependency := set.InheritFrom != ""
		_, dependencyIsMapped := mapSets[set.InheritFrom]

		if hasDependency && !dependencyIsMapped {
			visitCount[set.Name]++
			if visitCount[set.Name] > len(sets)+1 {
				return nil, fmt.Errorf("circular dependency detected for plugin set: %s", set.Name)
			}

			sets = append(sets, set)
			continue
		}

		orderedSets = append(orderedSets, set)
		mapSets[set.Name] = set
	}

	return orderedSets, nil
}
