// Package flatten implements OpenAPI YAML flattening: it inlines $ref schemas,
// merges oneOf/anyOf/allOf composition, and removes unreferenced definitions.
package flatten

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Flatten takes raw OpenAPI YAML bytes and returns the flattened YAML bytes.
func Flatten(data []byte) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing YAML: %w", err)
	}
	if len(doc.Content) == 0 {
		return data, nil
	}

	root := doc.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected root to be a YAML mapping")
	}

	applyDiscriminatorTransformations(root)
	applyOneOfTransformations(root)
	applyAnyOfTransformations(root)
	applyAllOfTransformations(root)

	for removeUnusedSchemas(root) {
	}

	applyStructuralTransformations(root)

	return marshalYAML(root)
}

func marshalYAML(node *yaml.Node) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(node); err != nil {
		return nil, fmt.Errorf("marshaling YAML: %w", err)
	}
	if err := enc.Close(); err != nil {
		return nil, fmt.Errorf("closing YAML encoder: %w", err)
	}
	return buf.Bytes(), nil
}
