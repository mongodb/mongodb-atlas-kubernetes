package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: [FORMAT=json|summary] %s {regressions|flakiness}\n", os.Args[0])
		os.Exit(1)
	}
	query := strings.ToLower(os.Args[1])
	if report, err := report(query); err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprint(os.Stdout, report)
	}
}
