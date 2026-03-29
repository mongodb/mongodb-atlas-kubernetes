package flatten

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// applyDiscriminatorTransformations finds all top-level schemas with a
// discriminator field (but no oneOf) and synthesises a oneOf from the
// discriminator mapping.
func applyDiscriminatorTransformations(root *yaml.Node) {
	schemas := asMapping(getPath(root, "components", "schemas"))
	if schemas == nil {
		return
	}
	for i := 0; i < len(schemas.Content)-1; i += 2 {
		name := schemas.Content[i].Value
		schema := asMapping(schemas.Content[i+1])
		if schema == nil {
			continue
		}
		if mappingHas(schema, "discriminator") {
			transformDiscriminatorToOneOf(name, schema)
		}
	}
}

func transformDiscriminatorToOneOf(name string, schema *yaml.Node) {
	if mappingHas(schema, "oneOf") {
		return // already has oneOf – nothing to do
	}
	disc := asMapping(mappingGet(schema, "discriminator"))
	if disc == nil {
		return
	}
	mapping := asMapping(mappingGet(disc, "mapping"))
	if mapping == nil {
		fmt.Fprintf(os.Stderr, "warning: skipping discriminator for %s: no mapping and no oneOf\n", name)
		return
	}

	// Collect unique refs in order.
	seen := map[string]bool{}
	var refs []string
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		ref := asString(mapping.Content[i+1])
		if !seen[ref] {
			seen[ref] = true
			refs = append(refs, ref)
		}
	}

	for _, ref := range refs {
		if schemaNameFromRef(ref) == name {
			fmt.Fprintf(os.Stderr, "error: %s.discriminator.mapping contains $ref to itself\n", name)
			return
		}
	}

	oneOf := newSequenceNode()
	for _, ref := range refs {
		oneOf.Content = append(oneOf.Content, newRefNode(ref))
	}
	mappingSet(schema, "oneOf", oneOf)
}
