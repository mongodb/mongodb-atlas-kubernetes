package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/josvazg/crd2go/internal/crd2go"
)

func main() {
	var input, output string
	flag.StringVar(&input, "input", "crds.yaml", "input YAML to process")
	flag.StringVar(&output, "output", "structs.go", "output Go to produce")
	flag.Parse()
	err := generate(output, input)
	if err != nil {
		log.Fatalf("Failed to generate go sturcts: %v", err)
	}
	log.Printf("Semantics applied to %s", output)
}

func generate(output, input string) error {
	o, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %w", output, err)
	}
	defer o.Close()
	i, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", input, err)
	}
	return crd2go.Generate(o, i)
}
