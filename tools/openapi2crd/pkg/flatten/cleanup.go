package flatten

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// removeUnusedSchemas removes schema definitions that are not referenced
// anywhere in the document. Returns true if anything was removed.
func removeUnusedSchemas(root *yaml.Node) bool {
	refs := collectRefs(root)

	used := map[string]bool{}
	for ref := range refs {
		if name, ok := strings.CutPrefix(ref, "#/components/schemas/"); ok {
			used[name] = true
		}
	}

	schemas := asMapping(getPath(root, "components", "schemas"))
	if schemas == nil {
		return false
	}

	before := len(schemas.Content)
	var kept []*yaml.Node
	for i := 0; i < len(schemas.Content)-1; i += 2 {
		key := schemas.Content[i]
		val := schemas.Content[i+1]
		if used[key.Value] {
			kept = append(kept, key, val)
		}
	}
	schemas.Content = kept
	return len(schemas.Content) != before
}
