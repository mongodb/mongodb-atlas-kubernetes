package flatten

import (
	"log/slog"

	"gopkg.in/yaml.v3"
)

func applyAnyOfTransformations(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool { return mappingHas(m, "anyOf") },
		func(path string, m *yaml.Node) {
			if isTopLevelSchema(path) {
				transformAnyOf(path, m, root)
			}
		},
	)
}

func transformAnyOf(path string, schema *yaml.Node, root *yaml.Node) {
	anyOf := asSequence(mappingGet(schema, "anyOf"))
	if anyOf == nil {
		return
	}
	children := resolveSequenceRefs(anyOf, root, "anyOf")

	if allHaveEnum(children) {
		mergeEnums(schema, children, "anyOf")
	} else {
		transformAnyOfProperties(path, schema, children)
	}
}

func transformAnyOfProperties(path string, schema *yaml.Node, children []*yaml.Node) {
	for _, child := range children {
		props := asMapping(mappingGet(child, "properties"))
		if props == nil {
			props = syntheticPropsFromAllOf(child)
			if props == nil {
				typeStr := asString(mappingGet(child, "type"))
				slog.Warn("skipping anyOf variant without mergeable properties", "phase", "flatten", "path", path, "type", typeStr)
				continue
			}
			mappingSet(child, "properties", props)
		}

		childCopy := deepCopy(props)
		parentProps := asMapping(mappingGet(schema, "properties"))
		if parentProps == nil {
			parentProps = newMappingNode()
			mappingSet(schema, "properties", parentProps)
		}
		for i := 0; i < len(childCopy.Content)-1; i += 2 {
			mappingSet(parentProps, childCopy.Content[i].Value, childCopy.Content[i+1])
		}
	}

	mappingDelete(schema, "discriminator")
	mappingDelete(schema, "anyOf")
}
