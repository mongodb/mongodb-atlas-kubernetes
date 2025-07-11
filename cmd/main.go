package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/josvazg/crd2go/internal/crd2go"
)

func main() {
	var input, output, reserved, skipList string
	flag.StringVar(&input, "input", "crds.yaml", "input YAML to process")
	flag.StringVar(&output, "output", ".", "output directory to produce source code to")
	flag.StringVar(&reserved, "reserved", "", "comma separated list of types names to avoid")
	flag.StringVar(&skipList, "skipList", "", "comma separated list of CRD to skip code generation for")
	flag.Parse()
	reservations := commaSeparatedListFrom(reserved)
	skips := commaSeparatedListFrom(skipList)
	err := generate(output, input, reservations, skips)
	if err != nil {
		log.Fatalf("Failed to generate go structs: %v", err)
	}
	log.Printf("Code generated at %s", output)
}

func generate(output, input string, reservations, skips []string) error {
	in, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", input, err)
	}
	preloaded := crd2go.KnownTypes()
	for _, name := range reservations {
		preloaded = append(preloaded, crd2go.NewOpaqueType(name))
	}
	cfg := crd2go.GenerateConfig{
		Version: crd2go.FirstVersion,
		Skips: skips,
		PreloadedTypes: preloaded,
	}
	return crd2go.GenerateStream(crd2go.CodeFileForCRDAtPath(output), in, &cfg)
}

func commaSeparatedListFrom(arg string) []string {
	if arg != "" {
		return strings.Split(arg, ",")
	}
	return []string{}
}
