package kube

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation"
)

func TestNormalizeIdentifier(t *testing.T) {
	t.Run("Valid Identifier", func(t *testing.T) {
		successCases := []string{
			"a", "ab", "abc", "a1", "a-1", "a--1--2--b",
			"a.a", "ab.a", "abc.a", "a1.a", "a-1.a", "a--1--2--b.a",
			"a.1", "ab.1", "abc.1", "a1.1", "a-1.1", "a--1--2--b.1",
			"0.a", "01.a", "012.a", "1a.a", "1-a.a", "1--a--b--2",
			"a.b.c.d.e", "aa.bb.cc.dd.ee", "1.2.3.4.5", "11.22.33.44.55", strings.Repeat("a", 253),
		}
		for i := range successCases {
			var normalized string
			if normalized = NormalizeIdentifier(successCases[i]); normalized != successCases[i] {
				t.Errorf("case[%d]: %q: the output is %q", i, successCases[i], normalized)
			}
			assert.Len(t, validation.IsDNS1123Subdomain(normalized), 0)
		}
	})
	t.Run("Invalid Identifier", func(t *testing.T) {
		successCases := [][]string{
			{"a_b", "a-b"},
			{"a__b", "a-b"},
			{"a.b#c", "a.b-c"},
			{"ab$", "ab"},
			{"ab*c$d", "ab-c-d"},
			{"a/b/c/", "a-b-c"},
			{"*A*", "a"},
			{"aa:bb:1", "aa-bb-1"},
			{strings.Repeat("a", 254), strings.Repeat("a", 253)},
		}
		for i := range successCases {
			var normalized string
			if normalized = NormalizeIdentifier(successCases[i][0]); normalized != successCases[i][1] {
				t.Errorf("case[%d]: %q: the output is %q (expected %q)", i, successCases[i][0], normalized, successCases[i][1])
			}
			assert.Len(t, validation.IsDNS1123Subdomain(normalized), 0, "case[%d]: %q: the output is %q", i, successCases[i][0], normalized)
		}
	})
}
func TestNormalizeLabelValue(t *testing.T) {
	t.Run("Valid Label Value", func(t *testing.T) {
		successCases := []string{
			"a", "ab", "abc", "a1", "a-1", "a--1--2--b",
			"a_b", "a.b",
			strings.Repeat("a", 63),
		}
		for i := range successCases {
			var normalized string
			if normalized = NormalizeLabelValue(successCases[i]); normalized != successCases[i] {
				t.Errorf("case[%d]: %q: the output is %q", i, successCases[i], normalized)
			}
			assert.Len(t, validation.IsValidLabelValue(normalized), 0)
		}
	})
	t.Run("Invalid Label Value", func(t *testing.T) {
		successCases := [][]string{
			{"a.#b", "a.-b"},
			{"ab#", "ab"},
			{"ab*c$d", "ab-c-d"},
			{"/a/b/c/", "a-b-c"},
			{"A*", "a"},
			{"aa:bb:1", "aa-bb-1"},
			{strings.Repeat("a", 63) + "bb", strings.Repeat("a", 63)},
			{strings.Repeat("a", 62) + "%", strings.Repeat("a", 62)},
		}
		for i := range successCases {
			var normalized string
			if normalized = NormalizeLabelValue(successCases[i][0]); normalized != successCases[i][1] {
				t.Errorf("case[%d]: %q: the output is %q (expected %q)", i, successCases[i][0], normalized, successCases[i][1])
			}
			assert.Len(t, validation.IsValidLabelValue(normalized), 0)
		}
	})
}
