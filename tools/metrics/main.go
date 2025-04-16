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
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: [FORMAT=json|summary] [TILL=date] %s {regressions|flakiness}\n", os.Args[0])
		os.Exit(1)
	}
	query := strings.ToLower(os.Args[1])
	format := valueOrDefault(os.Getenv("FORMAT"), "json")
	end := mustParseRFC3999(os.Getenv("TILL"))
	if report, err := Report(NewDefaultQueryClient(), end, query, format); err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprint(os.Stdout, report)
	}
}

func valueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func mustParseRFC3999(date string) time.Time {
	if date == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		panic(err)
	}
	return t
}
