package flatten

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func applyAllOfTransformations(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool { return mappingHas(m, "oneOf") },
		func(path string, m *yaml.Node) {
			if isTopLevelSchema(path) {
				transformAllOf(m, lastPathComponent(path), root)
			}
		},
	)
}

func transformAllOf(schema *yaml.Node, parentName string, root *yaml.Node) {
	oneOf := asSequence(mappingGet(schema, "oneOf"))
	if oneOf == nil {
		return
	}

	// Snapshot of parent without oneOf/discriminator for injection into children.
	expandedParent := copyExcluding(schema, "oneOf", "discriminator")

	for _, item := range oneOf.Content {
		refNode := asMapping(item)
		if refNode == nil {
			continue
		}
		ref := asString(mappingGet(refNode, "$ref"))
		if ref == "" {
			continue
		}
		_, child := resolveRef(root, ref)
		if child == nil {
			fmt.Fprintf(os.Stderr, "error: missing object reference %s for %s\n", ref, parentName)
			continue
		}
		if removeParentFromAllOf(child, parentName) {
			allOf := asSequence(mappingGet(child, "allOf"))
			if allOf == nil {
				allOf = newSequenceNode()
				mappingSet(child, "allOf", allOf)
			}
			allOf.Content = append(allOf.Content, expandedParent)
			flattenAllOf(child, root)
		}
	}

	mappingDelete(schema, "properties")
	mappingDelete(schema, "required")
}

// copyExcluding returns a deep copy of schema omitting the listed keys.
func copyExcluding(schema *yaml.Node, exclude ...string) *yaml.Node {
	excl := map[string]bool{}
	for _, e := range exclude {
		excl[e] = true
	}
	dst := newMappingNode()
	for i := 0; i < len(schema.Content)-1; i += 2 {
		key := schema.Content[i]
		val := schema.Content[i+1]
		if !excl[key.Value] {
			dst.Content = append(dst.Content, deepCopy(key), deepCopy(val))
		}
	}
	return dst
}

// removeParentFromAllOf removes allOf entries referencing parentName.
// Returns true if anything was removed.
func removeParentFromAllOf(child *yaml.Node, parentName string) bool {
	allOf := asSequence(mappingGet(child, "allOf"))
	if allOf == nil {
		return false
	}
	before := len(allOf.Content)
	var kept []*yaml.Node
	for _, item := range allOf.Content {
		m := asMapping(item)
		if m != nil {
			if ref := asString(mappingGet(m, "$ref")); schemaNameFromRef(ref) == parentName {
				continue
			}
		}
		kept = append(kept, item)
	}
	allOf.Content = kept
	return len(allOf.Content) != before
}

// flattenAllOf merges all allOf entries into the object's properties/required.
func flattenAllOf(obj *yaml.Node, root *yaml.Node) {
	allOf := asSequence(mappingGet(obj, "allOf"))
	if allOf == nil {
		return
	}

	props := asMapping(mappingGet(obj, "properties"))
	if props == nil {
		props = newMappingNode()
		mappingSet(obj, "properties", props)
	}

	reqSet := map[string]bool{}
	var reqOrder []string
	if existing := asSequence(mappingGet(obj, "required")); existing != nil {
		for _, v := range existing.Content {
			if s := asString(v); s != "" && !reqSet[s] {
				reqSet[s] = true
				reqOrder = append(reqOrder, s)
			}
		}
	}

	for _, item := range allOf.Content {
		resolved := resolveOrInline(item, root)
		if resolved == nil {
			continue
		}
		if rProps := asMapping(mappingGet(resolved, "properties")); rProps != nil {
			for i := 0; i < len(rProps.Content)-1; i += 2 {
				mappingSet(props, rProps.Content[i].Value, deepCopy(rProps.Content[i+1]))
			}
		}
		if rReq := asSequence(mappingGet(resolved, "required")); rReq != nil {
			for _, v := range rReq.Content {
				if s := asString(v); s != "" && !reqSet[s] {
					reqSet[s] = true
					reqOrder = append(reqOrder, s)
				}
			}
		}
	}

	if len(reqOrder) > 0 {
		reqNode := newSequenceNode()
		for _, r := range reqOrder {
			reqNode.Content = append(reqNode.Content, newStringNode(r))
		}
		mappingSet(obj, "required", reqNode)
	}

	mappingDelete(obj, "allOf")
}
