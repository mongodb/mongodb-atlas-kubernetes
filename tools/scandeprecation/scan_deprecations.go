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
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"go.uber.org/zap"
)

type DeprecationResponse struct {
	Type       string `json:"type"`
	Date       string `json:"date"`
	JavaMethod string `json:"javaMethod"`
}

func main() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log := l.Sugar()

	// Take input from stdin (pipe from GitHub CLI)
	scanner := bufio.NewScanner(os.Stdin)
	var responses []DeprecationResponse
	var out strings.Builder
	// Markdown table header for GH Comment
	out.WriteString("|Type|Java Method|Date|\n|------|------|------|\n")

	for scanner.Scan() {
		line, err := parseLogLine(scanner.Text())
		if err != nil {
			log.Warn(err.Error())
			continue
		}
		responses = append(responses, line)
	}

	// Quit out if there is no deprecations logged
	if len(responses) == 0 {
		log.Info("no responses found")
		os.Exit(0)
	}

	// Sort & Compact to remove duplicates
	slices.SortFunc(responses, func(a, b DeprecationResponse) int {
		if a.JavaMethod != b.JavaMethod {
			return cmp.Compare(a.JavaMethod, b.JavaMethod)
		}
		if a.Type != b.Type {
			return cmp.Compare(a.Type, b.Type)
		}
		return cmp.Compare(a.Date, b.Date)
	})
	responses = slices.Compact(responses)

	// Build the Markdown table
	for _, res := range responses {
		out.WriteString(res.Type + "|" + res.JavaMethod + "|" + res.Date + "|\n")
	}

	// Print to stdout
	fmt.Println(out.String())
}

func parseLogLine(line string) (DeprecationResponse, error) {
	// Non-JSON logs mean we split by tabs
	split := strings.Split(line, "\t")
	fmt.Println(split)

	// Last element is the JSON "body" of the log line
	line = split[len(split)-1]
	line = strings.Replace(line, "***", "{", 1)
	line = strings.Replace(line, "***", "}", 1)

	res := DeprecationResponse{}
	err := json.Unmarshal([]byte(line), &res)

	return res, err
}
