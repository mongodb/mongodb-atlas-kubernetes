package flatten

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// mappingGet returns the value node for the given key in a MappingNode, or nil.
func mappingGet(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

// mappingSet sets or adds a key-value pair in a MappingNode.
func mappingSet(node *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			node.Content[i+1] = value
			return
		}
	}
	node.Content = append(node.Content, newStringNode(key), value)
}

// mappingDelete removes a key from a MappingNode. Returns true if the key was found.
func mappingDelete(node *yaml.Node, key string) bool {
	for i := 0; i < len(node.Content)-1; i += 2 {
		if node.Content[i].Value == key {
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return true
		}
	}
	return false
}

// mappingHas reports whether a MappingNode contains the given key.
func mappingHas(node *yaml.Node, key string) bool {
	return mappingGet(node, key) != nil
}

// asMapping returns n if it is a MappingNode (resolving aliases), or nil.
func asMapping(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.AliasNode {
		return asMapping(n.Alias)
	}
	if n.Kind == yaml.MappingNode {
		return n
	}
	return nil
}

// asSequence returns n if it is a SequenceNode (resolving aliases), or nil.
func asSequence(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.AliasNode {
		return asSequence(n.Alias)
	}
	if n.Kind == yaml.SequenceNode {
		return n
	}
	return nil
}

// asString returns the string value of a ScalarNode (resolving aliases), or "".
func asString(n *yaml.Node) string {
	if n == nil {
		return ""
	}
	if n.Kind == yaml.AliasNode {
		return asString(n.Alias)
	}
	if n.Kind == yaml.ScalarNode {
		return n.Value
	}
	return ""
}

// getPath navigates a chain of mapping keys from root, returning the final node.
func getPath(root *yaml.Node, keys ...string) *yaml.Node {
	cur := root
	for _, k := range keys {
		cur = mappingGet(cur, k)
		if cur == nil {
			return nil
		}
	}
	return cur
}

// resolveRef resolves "#/components/schemas/<Name>" and returns (name, node).
func resolveRef(root *yaml.Node, ref string) (string, *yaml.Node) {
	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return "", nil
	}
	name := strings.TrimPrefix(ref, prefix)
	schemas := getPath(root, "components", "schemas")
	if schemas == nil {
		return name, nil
	}
	return name, asMapping(mappingGet(schemas, name))
}

// schemaNameFromRef extracts the schema name from "#/components/schemas/<Name>".
func schemaNameFromRef(ref string) string {
	return strings.TrimPrefix(ref, "#/components/schemas/")
}

// isTopLevelSchema reports whether the dot-separated path is ".components.schemas.<name>".
func isTopLevelSchema(path string) bool {
	parts := strings.Split(path, ".")
	// path starts with "." → ["", "components", "schemas", "<name>"]
	return len(parts) == 4 && parts[1] == "components" && parts[2] == "schemas"
}

// lastPathComponent returns the last "."-separated component of a path.
func lastPathComponent(path string) string {
	parts := strings.Split(path, ".")
	return parts[len(parts)-1]
}

// walkMappings traverses the YAML tree depth-first and calls fn for every
// MappingNode that passes the filter, along with its dot-path.
func walkMappings(root *yaml.Node, filter func(*yaml.Node) bool, fn func(path string, n *yaml.Node)) {
	type item struct {
		path string
		node *yaml.Node
	}
	stack := []item{{path: "", node: root}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		n := cur.node
		if n == nil {
			continue
		}
		if n.Kind == yaml.AliasNode {
			n = n.Alias
		}

		switch n.Kind {
		case yaml.MappingNode:
			if filter(n) {
				fn(cur.path, n)
			}
			for i := 0; i < len(n.Content)-1; i += 2 {
				key := n.Content[i].Value
				val := n.Content[i+1]
				if val.Kind == yaml.MappingNode || val.Kind == yaml.SequenceNode || val.Kind == yaml.AliasNode {
					stack = append(stack, item{path: cur.path + "." + key, node: val})
				}
			}
		case yaml.SequenceNode:
			for i, child := range n.Content {
				if child.Kind == yaml.MappingNode || child.Kind == yaml.SequenceNode || child.Kind == yaml.AliasNode {
					stack = append(stack, item{
						path: fmt.Sprintf("%s.%d", cur.path, i),
						node: child,
					})
				}
			}
		}
	}
}

// collectRefs returns every "$ref" value found anywhere in the tree.
func collectRefs(node *yaml.Node) map[string]bool {
	refs := map[string]bool{}
	collectRefsInto(node, refs)
	return refs
}

func collectRefsInto(node *yaml.Node, refs map[string]bool) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content)-1; i += 2 {
			key := node.Content[i].Value
			val := node.Content[i+1]
			if key == "$ref" {
				refs[val.Value] = true
			} else {
				collectRefsInto(val, refs)
			}
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			collectRefsInto(child, refs)
		}
	case yaml.DocumentNode:
		for _, child := range node.Content {
			collectRefsInto(child, refs)
		}
	case yaml.AliasNode:
		collectRefsInto(node.Alias, refs)
	}
}

// deepCopy returns a deep copy of a yaml.Node, inlining any aliases.
func deepCopy(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	if n.Kind == yaml.AliasNode {
		return deepCopy(n.Alias)
	}
	dst := *n
	if len(n.Content) > 0 {
		dst.Content = make([]*yaml.Node, len(n.Content))
		for i, child := range n.Content {
			dst.Content[i] = deepCopy(child)
		}
	}
	return &dst
}

// normalizeScalarStyles walks the entire tree and sets the correct quoting style
// on every string scalar, based on YAML 1.2 core-schema rules.
func normalizeScalarStyles(node *yaml.Node) {
	if node == nil {
		return
	}
	switch node.Kind {
	case yaml.ScalarNode:
		if node.Tag == "!!str" {
			v := node.Value
			if strings.Contains(v, "\n") {
				node.Style = yaml.LiteralStyle
			} else if needsQuotingYAML12(v) {
				if strings.Contains(v, `"`) && !strings.Contains(v, "'") {
					node.Style = yaml.SingleQuotedStyle
				} else {
					node.Style = yaml.DoubleQuotedStyle
				}
			} else {
				node.Style = 0
			}
		} else {
			node.Style = 0
		}
	case yaml.MappingNode, yaml.SequenceNode, yaml.DocumentNode:
		for _, child := range node.Content {
			normalizeScalarStyles(child)
		}
	case yaml.AliasNode:
	}
}

// ---- Node constructors ----

func newStringNode(value string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: value,
	}
}

func newMappingNode() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
}

func newSequenceNode() *yaml.Node {
	return &yaml.Node{
		Kind: yaml.SequenceNode,
		Tag:  "!!seq",
	}
}

func newRefNode(ref string) *yaml.Node {
	m := newMappingNode()
	mappingSet(m, "$ref", newStringNode(ref))
	return m
}
