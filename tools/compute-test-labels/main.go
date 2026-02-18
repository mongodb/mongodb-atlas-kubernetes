// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type labelSet struct {
	prLabels     string
	intLabels    string
	e2eLabels    string
	e2e2Labels   string
	skipPrefixes string
}

func jsonDump(data any) string {
	r, _ := json.Marshal(data)
	return string(r)
}

func MatchWildcards(labels []string, testLabels []string, testType string) []string {
	matchedLabels := make(map[string]struct{})
	prefixMatch := regexp.MustCompile(fmt.Sprintf("^test/%s/(.+)$", testType))

	for _, label := range labels {
		if label == fmt.Sprintf("test/%s/*", testType) {
			for _, test := range testLabels {
				matchedLabels[test] = struct{}{}
			}
		} else {
			if prefixMatch.MatchString(label) {
				pattern := strings.Replace(prefixMatch.FindStringSubmatch(label)[1], "*", ".*", -1)
				for _, test := range testLabels {
					matched, _ := regexp.MatchString("^"+pattern+"$", test)
					if matched {
						matchedLabels[test] = struct{}{}
					}
				}
			}
		}
	}

	result := make([]string, 0, len(matchedLabels))
	for key := range matchedLabels {
		result = append(result, key)
	}

	return result
}

func FilterLabelsDoNotContain(labels []string, substr string) []string {
	filtered := make([]string, 0, len(labels))
	for _, label := range labels {
		if !strings.Contains(label, substr) {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func FilterLabelsContain(labels []string, substr string) []string {
	filtered := make([]string, 0, len(labels))
	for _, label := range labels {
		if strings.Contains(label, substr) {
			filtered = append(filtered, label)
		}
	}
	return filtered
}

func SkipLabelsByPrefix(labels []string, skipPrefixes []string) []string {
	if len(skipPrefixes) == 0 {
		return labels
	}
	filtered := make([]string, 0, len(labels))
	for _, label := range labels {
		if hasSkipPrefix(label, skipPrefixes) {
			continue
		}
		filtered = append(filtered, label)
	}
	return filtered
}

func hasSkipPrefix(label string, skipPrefixes []string) bool {
	for _, skipPrefix := range skipPrefixes {
		if strings.HasPrefix(label, skipPrefix) {
			return true
		}
	}
	return false
}

func computeTestLabels(out io.Writer, outputJSON bool, inputs *labelSet) error {
	var labels []string
	var intLabels []string
	var e2eLabels []string
	var e2e2Labels []string
	var skipPrefixes []string

	if err := json.Unmarshal([]byte(inputs.prLabels), &labels); err != nil {
		return fmt.Errorf("Error parsing PR labels: %w", err)
	}
	if len(inputs.intLabels) > 0 {
		if err := json.Unmarshal([]byte(inputs.intLabels), &intLabels); err != nil {
			return fmt.Errorf("Error parsing integration tests labels: %w", err)
		}
	}
	if len(inputs.e2eLabels) > 0 {
		if err := json.Unmarshal([]byte(inputs.e2eLabels), &e2eLabels); err != nil {
			return fmt.Errorf("Error parsing E2E tests labels: %w", err)
		}
	}
	if len(inputs.e2e2Labels) > 0 {
		if err := json.Unmarshal([]byte(inputs.e2e2Labels), &e2e2Labels); err != nil {
			return fmt.Errorf("Error parsing E2E2 tests labels: %w", err)
		}
	}
	if len(inputs.skipPrefixes) > 0 {
		if err := json.Unmarshal([]byte(inputs.skipPrefixes), &skipPrefixes); err != nil {
			return fmt.Errorf("Error parsing Skip prefixes tests labels: %w", err)
		}
	}

	matchedIntTests := MatchWildcards(labels, SkipLabelsByPrefix(intLabels, skipPrefixes), "int")
	matchedE2ETests := MatchWildcards(labels, SkipLabelsByPrefix(e2eLabels, skipPrefixes), "e2e")
	matchedE2E2Tests := MatchWildcards(labels, SkipLabelsByPrefix(e2e2Labels, skipPrefixes), "e2e2")
	// These have to be executed in their own environment )
	matchedE2EGovTests := FilterLabelsContain(matchedE2ETests, "atlas-gov")

	matchedE2ETests = FilterLabelsDoNotContain(matchedE2ETests, "atlas-gov")

	matchedIntTestsJSON, _ := json.Marshal(matchedIntTests)
	matchedE2ETestsJSON, _ := json.Marshal(matchedE2ETests)
	matchedE2E2TestsJSON, _ := json.Marshal(matchedE2E2Tests)
	matchedE2EGovTestsJSON, _ := json.Marshal(matchedE2EGovTests)

	if outputJSON {
		res := map[string]any{}
		res["int"] = matchedIntTests
		res["e2e"] = matchedE2ETests
		res["e2e2"] = matchedE2E2Tests
		res["e2e_gov"] = matchedE2EGovTests
		fmt.Fprintln(out, jsonDump(res))
		return nil
	}
	fmt.Fprintf(out, "Matched Integration Tests: %s\n", matchedIntTestsJSON)
	fmt.Fprintf(out, "Matched E2E Tests: %s\n", matchedE2ETestsJSON)
	fmt.Fprintf(out, "Matched E2E2 Tests: %s\n", matchedE2E2TestsJSON)
	fmt.Fprintf(out, "Matched E2E GOV Tests: %s\n", matchedE2EGovTestsJSON)
	return nil
}

func main() {
	useJSON := os.Getenv("USE_JSON") != ""
	inputs := labelSet{
		prLabels:     os.Getenv("PR_LABELS"),
		intLabels:    os.Getenv("INT_LABELS"),
		e2eLabels:    os.Getenv("E2E_LABELS"),
		e2e2Labels:   os.Getenv("E2E2_LABELS"),
		skipPrefixes: os.Getenv("SKIP_PREFIXES"),
	}
	if err := computeTestLabels(os.Stdout, useJSON, &inputs); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
