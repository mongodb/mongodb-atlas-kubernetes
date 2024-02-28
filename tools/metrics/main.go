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
