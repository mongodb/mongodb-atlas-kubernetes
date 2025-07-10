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
	var input, output, reserved string
	flag.StringVar(&input, "input", "crds.yaml", "input YAML to process")
	flag.StringVar(&output, "output", ".", "output directory to produce source code to")
	flag.StringVar(&reserved, "reserved", "", "comma separated list of types names to avoid")
	flag.Parse()
	reservations := []string{}
	if reserved != "" {
		reservations = strings.Split(reserved, ",")
	}
	err := generate(output, input, reservations)
	if err != nil {
		log.Fatalf("Failed to generate go structs: %v", err)
	}
	log.Printf("Code generated at %s", output)
}

func generate(output, input string, reservations []string) error {
	i, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %w", input, err)
	}
	preloaded := crd2go.KnownTypes()
	for _, name := range reservations {
		preloaded = append(preloaded, crd2go.NewStruct(name, nil))
	}
	return crd2go.GenerateStream(crd2go.CodeFileForCRDAtPath(output), i, crd2go.FirstVersion, preloaded...)
}
