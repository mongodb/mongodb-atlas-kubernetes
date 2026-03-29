package flatten

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// resolveSequenceRefs resolves the refs/inline items in a oneOf/anyOf sequence
// into concrete MappingNodes, logging warnings for unresolvable refs.
func resolveSequenceRefs(seq *yaml.Node, root *yaml.Node, kind string) []*yaml.Node {
	var result []*yaml.Node
	for _, item := range seq.Content {
		m := asMapping(item)
		if m == nil {
			continue
		}
		if ref := asString(mappingGet(m, "$ref")); ref != "" {
			_, child := resolveRef(root, ref)
			if child == nil {
				fmt.Fprintf(os.Stderr, "warning: could not resolve %s reference: %s\n", kind, ref)
				continue
			}
			result = append(result, child)
		} else if mappingHas(m, "type") || mappingHas(m, "properties") {
			result = append(result, m)
		}
	}
	return result
}

// resolveOrInline resolves a $ref or returns the node as-is if it is inline.
func resolveOrInline(item *yaml.Node, root *yaml.Node) *yaml.Node {
	m := asMapping(item)
	if m == nil {
		return nil
	}
	if ref := asString(mappingGet(m, "$ref")); ref != "" {
		_, resolved := resolveRef(root, ref)
		return resolved
	}
	return m
}

// allHaveEnum reports whether every mapping in the list has an "enum" key.
func allHaveEnum(children []*yaml.Node) bool {
	for _, c := range children {
		if !mappingHas(c, "enum") {
			return false
		}
	}
	return true
}

// mergeEnums merges enum values from children into schema, then deletes
// discriminator and the given composition key (oneOf or anyOf).
func mergeEnums(schema *yaml.Node, children []*yaml.Node, compositionKey string) {
	seen := map[string]bool{}
	var vals []*yaml.Node

	if existing := asSequence(mappingGet(schema, "enum")); existing != nil {
		for _, v := range existing.Content {
			if s := asString(v); s != "" {
				seen[s] = true
				vals = append(vals, v)
			}
		}
	}

	var lastType *yaml.Node
	for _, child := range children {
		childEnum := asSequence(mappingGet(child, "enum"))
		if childEnum == nil {
			continue
		}
		for _, v := range childEnum.Content {
			if s := asString(v); s != "" && !seen[s] {
				seen[s] = true
				vals = append(vals, v)
			}
		}
		if t := mappingGet(child, "type"); t != nil {
			lastType = t
		}
	}

	enumNode := newSequenceNode()
	enumNode.Content = vals
	mappingSet(schema, "enum", enumNode)
	if lastType != nil {
		mappingSet(schema, "type", lastType)
	}
	mappingDelete(schema, "discriminator")
	mappingDelete(schema, compositionKey)
}

// syntheticPropsFromAllOf builds a synthetic properties MappingNode from
// inline (non-$ref) allOf entries that have a properties key.
// Returns nil if there are no such entries.
func syntheticPropsFromAllOf(child *yaml.Node) *yaml.Node {
	allOf := asSequence(mappingGet(child, "allOf"))
	if allOf == nil {
		return nil
	}
	synth := newMappingNode()
	for _, item := range allOf.Content {
		m := asMapping(item)
		if m == nil || mappingHas(m, "$ref") {
			continue
		}
		if props := asMapping(mappingGet(m, "properties")); props != nil {
			for i := 0; i < len(props.Content)-1; i += 2 {
				mappingSet(synth, props.Content[i].Value, props.Content[i+1])
			}
		}
	}
	if len(synth.Content) == 0 {
		return nil
	}
	return synth
}

// typeOrRef returns the "type" or "$ref" value of a property schema node,
// used for mismatch detection.
func typeOrRef(node *yaml.Node) string {
	if node == nil {
		return ""
	}
	if t := asString(mappingGet(node, "type")); t != "" {
		return t
	}
	return asString(mappingGet(node, "$ref"))
}

// ignoredProperties returns the set of property names exempt from
// type-mismatch deletion (from the Atlas duplicate.ignore.json list).
func ignoredProperties() map[string]bool {
	return map[string]bool{
		"units":           true,
		"threshold":       true,
		"eventTypeName":   true,
		"currentValue":    true,
		"metricThreshold": true,
		"autoScaling":     true,
		"featureId":       true,
	}
}
