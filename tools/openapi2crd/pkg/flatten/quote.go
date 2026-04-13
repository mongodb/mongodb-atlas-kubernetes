package flatten

import (
	"regexp"
	"strings"
)

// needsQuotingYAML12 reports whether a plain scalar value must be quoted when
// targeting YAML 1.2 core schema.  A value that would be misinterpreted as
// null, bool, or a number must be quoted to preserve its string identity.
func needsQuotingYAML12(v string) bool {
	if v == "" {
		return true
	}
	switch strings.ToLower(v) {
	case "null", "~", "true", "false":
		return true
	}
	if isYAML12Number(v) {
		return true
	}
	return containsYAMLSpecialChars(v)
}

// containsYAMLSpecialChars reports whether a plain scalar would be syntactically
// invalid or ambiguous in YAML 1.2 block context.
func containsYAMLSpecialChars(v string) bool {
	if v == "" {
		return true
	}
	switch v[0] {
	case '#', '|', '>', '!', '%', '@', '`', '"', '\'':
		return true
	case '&', '*', '?':
		return true
	case '-':
		if len(v) > 1 && v[1] == ' ' {
			return true
		}
	case ':':
		if len(v) == 1 || v[1] == ' ' {
			return true
		}
	case '{', '}', '[', ']':
		return true
	}
	if strings.Contains(v, ": ") || strings.Contains(v, " #") {
		return true
	}
	if strings.HasSuffix(v, ":") {
		return true
	}
	if strings.ContainsAny(v, "\n\r") {
		return true
	}
	return false
}

// isYAML12Number reports whether v would be parsed as a numeric type under
// the YAML 1.2 core schema.
func isYAML12Number(v string) bool {
	if v == ".inf" || v == "+.inf" || v == "-.inf" || v == ".nan" {
		return true
	}
	return reYAML12Number.MatchString(v)
}

// reYAML12Number matches YAML 1.2 numeric literals.
var reYAML12Number = regexp.MustCompile(
	`^[-+]?(\.[0-9]+|[0-9]+(\.[0-9]*)?)([eE][-+]?[0-9]+)?$` +
		`|^[-+]?0o[0-7]+$` +
		`|^[-+]?0x[0-9a-fA-F]+$`,
)
