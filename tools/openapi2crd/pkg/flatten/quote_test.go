package flatten

import "testing"

func TestNeedsQuotingYAML12(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"", true},

		{"true", true},
		{"True", true},
		{"TRUE", true},
		{"false", true},
		{"False", true},

		{"null", true},
		{"Null", true},
		{"NULL", true},
		{"~", true},

		{"yes", false},
		{"no", false},
		{"on", false},
		{"off", false},

		{"0", true},
		{"42", true},
		{"-1", true},
		{"+1", true},
		{"3.14", true},
		{".5", true},
		{"1e10", true},
		{"1.5e-3", true},
		{"0x1F", true},
		{"0o17", true},
		{".inf", true},
		{"+.inf", true},
		{"-.inf", true},
		{".nan", true},

		{"hello", false},
		{"hello world", false},
		{"snake_case", false},
		{"camelCase", false},
		{"v1", false},
		{"application/json", false},
		{"#/components/schemas/Foo", true},

		{"#comment", true},
		{"|literal", true},
		{">folded", true},
		{"!tag", true},
		{"%directive", true},
		{"@at", true},
		{"`backtick", true},
		{`"quoted`, true},
		{"'single", true},
		{"&anchor", true},
		{"*alias", true},
		{"?key", true},
		{"{flow", true},
		{"}flow", true},
		{"[seq", true},
		{"]seq", true},

		{":", true},
		{": ", true},
		{"key: value", true},
		{"key:", true},
		{"key:value", false},
		{"http://example.com", false},

		{"- item", true},
		{"-item", false},

		{"foo #bar", true},

		{"line1\nline2", true},
		{"line1\rline2", true},
	}

	for _, tt := range tests {
		got := needsQuotingYAML12(tt.value)
		if got != tt.want {
			t.Errorf("needsQuotingYAML12(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestIsYAML12Number(t *testing.T) {
	numbers := []string{
		"0", "1", "42", "-1", "+1",
		"3.14", ".5", "1.", "-3.14",
		"1e10", "1.5e-3", "1E+10",
		"0x1F", "0xdeadbeef", "0xDEAD",
		"0o17", "0o777",
		".inf", "+.inf", "-.inf", ".nan",
	}
	for _, v := range numbers {
		if !isYAML12Number(v) {
			t.Errorf("isYAML12Number(%q) = false, want true", v)
		}
	}

	nonNumbers := []string{
		"", "hello", "true", "false", "null",
		"1.2.3", "0b1010", "yes", "no",
	}
	for _, v := range nonNumbers {
		if isYAML12Number(v) {
			t.Errorf("isYAML12Number(%q) = true, want false", v)
		}
	}
}
