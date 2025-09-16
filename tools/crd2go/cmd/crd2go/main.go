package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/internal/checkerr"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/pkg/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/crd2go/pkg/crd2go"
)

func main() {
	var input, output, config string
	flag.StringVar(&input, "input", "", "input YAML to process")
	flag.StringVar(&output, "output", "", "output directory to produce source code to")
	flag.StringVar(&config, "config", "crd2go.yaml", "YAML file with the CRD2Go config")
	flag.Parse()

	cfg, err := generate(input, output, config)
	if err != nil {
		log.Fatalf("Failed to generate go structs: %v", err)
	}
	log.Printf("Code generated at %s", cfg.Output)
}

func generate(input, output, config string) (*config.Config, error) {
	f, err := os.Open(config)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file: %w", err)
	}
	defer checkerr.CheckErr("closing config file", f.Close)
	cfg, err := crd2go.LoadConfig(f)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	if input != "" {
		cfg.Input = input
	}
	if output != "" {
		cfg.Output = output
	}
	if err := crd2go.GenerateToDir(cfg); err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}
	return cfg, nil
}
