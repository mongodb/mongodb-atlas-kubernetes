package flatten

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func applyOneOfTransformations(root *yaml.Node) {
	walkMappings(root,
		func(m *yaml.Node) bool { return mappingHas(m, "oneOf") },
		func(path string, m *yaml.Node) {
			if isTopLevelSchema(path) {
				transformOneOf(m, root)
			}
		},
	)
}

func transformOneOf(schema *yaml.Node, root *yaml.Node) {
	oneOf := asSequence(mappingGet(schema, "oneOf"))
	if oneOf == nil {
		return
	}
	children := resolveSequenceRefs(oneOf, root, "oneOf")

	if allHaveEnum(children) {
		mergeEnums(schema, children, "oneOf")
	} else {
		transformOneOfProperties(schema, children, root)
	}
}

func transformOneOfProperties(schema *yaml.Node, children []*yaml.Node, root *yaml.Node) {
	discExt := buildDiscriminatorExtension(schema, root)

	for _, child := range children {
		props := asMapping(mappingGet(child, "properties"))
		if props == nil {
			props = syntheticPropsFromAllOf(child)
			if props == nil {
				typeStr := asString(mappingGet(child, "type"))
				fmt.Fprintf(os.Stderr, "warning: skipping non-object oneOf variant (type: %s)\n", typeStr)
				continue
			}
			mappingSet(child, "properties", props)
		}
		mergeChildPropertiesIntoParent(schema, props)
	}

	if discExt != nil {
		mappingSet(schema, "x-xgen-discriminator", discExt)
	}
	mappingDelete(schema, "discriminator")
	mappingDelete(schema, "oneOf")
}

// mergeChildPropertiesIntoParent merges child properties into the parent schema.
// Type-mismatched properties are deleted (except ignored ones); remaining
// properties use last-wins semantics, preserving the parent's description when
// the child's description differs.
func mergeChildPropertiesIntoParent(parent *yaml.Node, childProps *yaml.Node) {
	childCopy := deepCopy(childProps)

	parentProps := asMapping(mappingGet(parent, "properties"))
	if parentProps != nil {
		parentTypeOrRef := map[string]string{}
		for i := 0; i < len(parentProps.Content)-1; i += 2 {
			parentTypeOrRef[parentProps.Content[i].Value] = typeOrRef(asMapping(parentProps.Content[i+1]))
		}

		ignored := ignoredProperties()
		var keysToDelete []string
		for i := 0; i < len(childCopy.Content)-1; i += 2 {
			k := childCopy.Content[i].Value
			if pt, exists := parentTypeOrRef[k]; exists {
				ct := typeOrRef(asMapping(childCopy.Content[i+1]))
				if pt != ct && !ignored[k] {
					keysToDelete = append(keysToDelete, k)
				}
			}
		}
		for _, k := range keysToDelete {
			mappingDelete(childCopy, k)
			mappingDelete(childProps, k)
		}

		for i := 0; i < len(childCopy.Content)-1; i += 2 {
			k := childCopy.Content[i].Value
			parentProp := asMapping(mappingGet(parentProps, k))
			childProp := asMapping(childCopy.Content[i+1])
			if parentProp == nil || childProp == nil {
				continue
			}
			parentDesc := mappingGet(parentProp, "description")
			if parentDesc == nil {
				continue
			}
			if asString(mappingGet(childProp, "description")) != asString(parentDesc) {
				mappingSet(childProp, "description", parentDesc)
			}
		}
	}

	if parentProps == nil {
		parentProps = newMappingNode()
		mappingSet(parent, "properties", parentProps)
	}

	for i := 0; i < len(childCopy.Content)-1; i += 2 {
		mappingSet(parentProps, childCopy.Content[i].Value, childCopy.Content[i+1])
	}
}

// buildDiscriminatorExtension builds the x-xgen-discriminator extension node.
func buildDiscriminatorExtension(schema *yaml.Node, root *yaml.Node) *yaml.Node {
	disc := asMapping(mappingGet(schema, "discriminator"))
	if disc == nil {
		return nil
	}
	propNameNode := mappingGet(disc, "propertyName")
	mapping := asMapping(mappingGet(disc, "mapping"))
	if propNameNode == nil || mapping == nil {
		return nil
	}
	propertyName := asString(propNameNode)

	baseProps := map[string]bool{}
	if pp := asMapping(mappingGet(schema, "properties")); pp != nil {
		for i := 0; i < len(pp.Content)-1; i += 2 {
			baseProps[pp.Content[i].Value] = true
		}
	}

	mappingNode := newMappingNode()
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		discValue := mapping.Content[i].Value
		ref := asString(mapping.Content[i+1])

		_, child := resolveRef(root, ref)
		if child == nil {
			continue
		}

		allChildProps := collectAllProperties(child, root)
		allChildRequired := collectAllRequired(child, root)

		var typeSpecific []string
		for j := 0; j < len(allChildProps.Content)-1; j += 2 {
			k := allChildProps.Content[j].Value
			if !baseProps[k] {
				typeSpecific = append(typeSpecific, k)
			}
		}

		var tsRequired []string
		for _, p := range typeSpecific {
			if p != propertyName && allChildRequired[p] {
				tsRequired = append(tsRequired, p)
			}
		}

		entry := newMappingNode()
		propsSeq := newSequenceNode()
		for _, p := range typeSpecific {
			propsSeq.Content = append(propsSeq.Content, newStringNode(p))
		}
		mappingSet(entry, "properties", propsSeq)
		if len(tsRequired) > 0 {
			reqSeq := newSequenceNode()
			for _, r := range tsRequired {
				reqSeq.Content = append(reqSeq.Content, newStringNode(r))
			}
			mappingSet(entry, "required", reqSeq)
		}
		mappingSet(mappingNode, discValue, entry)
	}

	ext := newMappingNode()
	mappingSet(ext, "propertyName", newStringNode(propertyName))
	mappingSet(ext, "mapping", mappingNode)
	return ext
}

func collectAllProperties(child *yaml.Node, root *yaml.Node) *yaml.Node {
	result := newMappingNode()
	if props := asMapping(mappingGet(child, "properties")); props != nil {
		for i := 0; i < len(props.Content)-1; i += 2 {
			mappingSet(result, props.Content[i].Value, props.Content[i+1])
		}
	}
	if allOf := asSequence(mappingGet(child, "allOf")); allOf != nil {
		for _, item := range allOf.Content {
			resolved := resolveOrInline(item, root)
			if resolved == nil {
				continue
			}
			if props := asMapping(mappingGet(resolved, "properties")); props != nil {
				for i := 0; i < len(props.Content)-1; i += 2 {
					mappingSet(result, props.Content[i].Value, props.Content[i+1])
				}
			}
		}
	}
	return result
}

func collectAllRequired(child *yaml.Node, root *yaml.Node) map[string]bool {
	result := map[string]bool{}
	if req := asSequence(mappingGet(child, "required")); req != nil {
		for _, v := range req.Content {
			result[asString(v)] = true
		}
	}
	if allOf := asSequence(mappingGet(child, "allOf")); allOf != nil {
		for _, item := range allOf.Content {
			resolved := resolveOrInline(item, root)
			if resolved == nil {
				continue
			}
			if req := asSequence(mappingGet(resolved, "required")); req != nil {
				for _, v := range req.Content {
					result[asString(v)] = true
				}
			}
		}
	}
	return result
}
