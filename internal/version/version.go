package version

import (
	"regexp"
	"strings"
)

const DefaultVersion = "unknown"

// Version set by the linker during link time.
var Version = DefaultVersion

func IsRelease(v string) bool {
	return v != DefaultVersion &&
		regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+[-certified]*$`).Match([]byte(strings.TrimSpace(v)))
}
