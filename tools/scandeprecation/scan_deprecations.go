package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"slices"
	"strings"
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

	// Take file from GH CLI tool (which has already aggregated logs)
	file, err := os.Open("/Users/roo.thorp/test/log.out") // TODO: This is a local test file currently :)
	if err != nil {
		log.Fatal("failed to open log file", zap.Error(err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var responses []DeprecationResponse
	var out strings.Builder
	// Markdown table header for GH Comment
	out.WriteString("|Type|Java Method|Date|\n|------|------|------|\n")

	for scanner.Scan() {
		example := scanner.Text()
		// Non-JSON logs mean we split by tabs
		split := strings.Split(example, "\t")

		// Last element is the JSON "body" of the log line
		example = split[len(split)-1]

		res := DeprecationResponse{}
		err = json.Unmarshal([]byte(example), &res)
		if err != nil {
			log.Warn("failed to unmarshal JSON", zap.Error(err))
			continue
		}
		responses = append(responses, res)
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
