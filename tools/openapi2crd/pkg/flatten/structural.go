package flatten

import "gopkg.in/yaml.v3"

// applyStructuralTransformations cleans up the schema tree to be compatible
// with Kubernetes structural schemas. It runs after the existing composition
// transformations (oneOf/anyOf/allOf top-level handlers) to catch inline
// compositions that remain in property schemas.
//
// Specifically it:
//   - Flattens inline allOf by merging items' type/properties/required into the parent
//   - Flattens inline oneOf/anyOf by merging all variants' properties into the parent
//   - Removes additionalProperties when properties is non-empty (mutually exclusive in CRDs)
//   - Ensures every array schema has an items sub-schema
func applyStructuralTransformations(root *yaml.Node) {
	flattenInlineAllOf(root)
	flattenInlineCompositions(root, "oneOf")
	flattenInlineCompositions(root, "anyOf")
	fixAdditionalPropertiesConflict(root)
	ensureArrayHasItems(root)
}

// flattenInlineAllOf walks every schema that has allOf, merges the items'
// type/properties/required into the parent, then removes allOf.
func flattenInlineAllOf(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool { return mappingHas(m, "allOf") },
		func(_ string, m *yaml.Node) {
			mergeAllOfInto(m, root)
		},
	)
}

// mergeAllOfInto merges each allOf item into schema and removes the allOf key.
func mergeAllOfInto(schema *yaml.Node, root *yaml.Node) {
	allOf := asSequence(mappingGet(schema, "allOf"))
	if allOf == nil {
		return
	}

	for _, item := range allOf.Content {
		resolved := resolveOrInline(item, root)
		if resolved == nil {
			continue
		}
		if mappingGet(schema, "type") == nil {
			if t := mappingGet(resolved, "type"); t != nil {
				mappingSet(schema, "type", deepCopy(t))
			}
		}
		for _, meta := range []string{"description", "title"} {
			if mappingGet(schema, meta) == nil {
				if v := mappingGet(resolved, meta); v != nil {
					mappingSet(schema, meta, deepCopy(v))
				}
			}
		}
		mergePropertiesAndRequired(schema, resolved)
	}

	mappingDelete(schema, "allOf")
}

// flattenInlineCompositions walks every schema that has the given composition
// key (oneOf or anyOf), merges each variant's properties into the parent, then
// removes the composition key. Enum-only compositions are left untouched.
func flattenInlineCompositions(root *yaml.Node, key string) {
	walkMappings(root,
		func(m *yaml.Node) bool { return mappingHas(m, key) },
		func(_ string, m *yaml.Node) {
			seq := asSequence(mappingGet(m, key))
			if seq == nil {
				return
			}
			children := resolveSequenceRefs(seq, root, key)
			if len(children) > 0 && allHaveEnum(children) {
				mergeEnums(m, children, key)
				return
			}
			for _, item := range seq.Content {
				resolved := resolveOrInline(item, root)
				if resolved == nil {
					continue
				}
				mergePropertiesAndRequired(m, resolved)
			}
			mappingDelete(m, key)
		},
	)
}

// mergePropertiesAndRequired copies properties and required from src into dst,
// without overwriting existing keys.
func mergePropertiesAndRequired(dst, src *yaml.Node) {
	if srcProps := asMapping(mappingGet(src, "properties")); srcProps != nil {
		dstProps := asMapping(mappingGet(dst, "properties"))
		if dstProps == nil {
			dstProps = newMappingNode()
			mappingSet(dst, "properties", dstProps)
		}
		for i := 0; i < len(srcProps.Content)-1; i += 2 {
			k := srcProps.Content[i].Value
			if mappingGet(dstProps, k) == nil {
				mappingSet(dstProps, k, deepCopy(srcProps.Content[i+1]))
			}
		}
	}

	if srcReq := asSequence(mappingGet(src, "required")); srcReq != nil {
		dstReq := asSequence(mappingGet(dst, "required"))
		existing := map[string]bool{}
		if dstReq != nil {
			for _, v := range dstReq.Content {
				existing[asString(v)] = true
			}
		} else {
			dstReq = newSequenceNode()
		}
		added := false
		for _, v := range srcReq.Content {
			if s := asString(v); s != "" && !existing[s] {
				dstReq.Content = append(dstReq.Content, newStringNode(s))
				existing[s] = true
				added = true
			}
		}
		if added {
			mappingSet(dst, "required", dstReq)
		}
	}
}

// fixAdditionalPropertiesConflict removes additionalProperties from any schema
// that also has a non-empty properties map. In Kubernetes CRDs these two fields
// are mutually exclusive.
func fixAdditionalPropertiesConflict(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool {
			if !mappingHas(m, "additionalProperties") {
				return false
			}
			props := asMapping(mappingGet(m, "properties"))
			return props != nil && len(props.Content) > 0
		},
		func(_ string, m *yaml.Node) {
			mappingDelete(m, "additionalProperties")
		},
	)
}

// ensureArrayHasItems adds an empty items schema ({}) to any array schema that
// lacks one. Kubernetes requires every array to declare its items type.
func ensureArrayHasItems(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool {
			return asString(mappingGet(m, "type")) == "array" && !mappingHas(m, "items")
		},
		func(_ string, m *yaml.Node) {
			mappingSet(m, "items", newMappingNode())
		},
	)
}
