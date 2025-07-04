package plugins

import "strings"

func jsonPath(path []string) string {
	result := strings.Join(path, ".")
	return strings.ReplaceAll(result, ".[*]", "[*]")
}
