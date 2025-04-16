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
	"os"
	"regexp"
	"strings"
)

func jsonDump(data interface{}) string {
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

func main() {
	envPRLabels := os.Getenv("PR_LABELS")
	envIntLabels := os.Getenv("INT_LABELS")
	envE2ELabels := os.Getenv("E2E_LABELS")
	envUseJSON := os.Getenv("USE_JSON")

	var labels []string
	var intLabels []string
	var e2eLabels []string

	if err := json.Unmarshal([]byte(envPRLabels), &labels); err != nil {
		fmt.Printf("Error parsing PR labels: %v\n", err)
		return
	}
	if err := json.Unmarshal([]byte(envIntLabels), &intLabels); err != nil {
		fmt.Printf("Error parsing integration tests labels: %v\n", err)
		return
	}
	if err := json.Unmarshal([]byte(envE2ELabels), &e2eLabels); err != nil {
		fmt.Printf("Error parsing E2E tests labels: %v\n", err)
		return
	}

	matchedIntTests := MatchWildcards(labels, intLabels, "int")
	matchedE2ETests := MatchWildcards(labels, e2eLabels, "e2e")

	matchedIntTestsJSON, _ := json.Marshal(matchedIntTests)
	matchedE2ETestsJSON, _ := json.Marshal(matchedE2ETests)

	if envUseJSON != "" {
		res := map[string]any{}
		res["int"] = matchedIntTests
		res["e2e"] = matchedE2ETests
		fmt.Println(jsonDump(res))
		return
	}
	fmt.Printf("Matched Integration Tests: %s\n", matchedIntTestsJSON)
	fmt.Printf("Matched E2E Tests: %s\n", matchedE2ETestsJSON)
}
