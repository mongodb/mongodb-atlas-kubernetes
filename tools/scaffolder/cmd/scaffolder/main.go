package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/config"
	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/registry"

	// Import generators to register them via init()
	_ "github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/atlascontrollers"
	_ "github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/atlasexporters"
	_ "github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generators/indexers"
)

var (
	inputPath         string
	crdKind           string
	listCRDs          bool
	allCRDs           bool
	controllerOutDir  string
	indexerOutDir     string
	exporterOutDir    string
	typesPath         string
	indexerTypesPath  string
	indexerImportPath string
	override          bool
	generators        string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scaffolder",
		Short: "Generate Kubernetes controllers for MongoDB Atlas CRDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputPath == "" {
				return fmt.Errorf("--input is required")
			}

			if listCRDs {
				return config.PrintCRDs(inputPath)
			}

			// Validate that --all and --crd are mutually exclusive
			if allCRDs && crdKind != "" {
				return fmt.Errorf("--all and --crd flags are mutually exclusive, please specify only one")
			}

			if !allCRDs && crdKind == "" {
				return fmt.Errorf("either --crd or --all flag must be specified")
			}

			// Validate types path
			if err := validateGoImportPath(typesPath); err != nil {
				return fmt.Errorf("invalid --types-path: %w", err)
			}

			// Set defaults for new flags
			if indexerTypesPath == "" {
				indexerTypesPath = typesPath
			}
			if indexerImportPath == "" {
				indexerImportPath = "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers"
			}

			// Validate new flags
			if err := validateGoImportPath(indexerTypesPath); err != nil {
				return fmt.Errorf("invalid --indexer-types-path: %w", err)
			}
			if err := validateGoImportPath(indexerImportPath); err != nil {
				return fmt.Errorf("invalid --indexer-import-path: %w", err)
			}

			// Set default output directories
			if controllerOutDir == "" {
				controllerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "controller")
			}
			if indexerOutDir == "" {
				indexerOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "indexer")
			}
			if exporterOutDir == "" {
				exporterOutDir = filepath.Join("..", "mongodb-atlas-kubernetes", "internal", "exporter")
			}

			// Parse generator names
			selectedGenerators, err := parseGenerators(generators)
			if err != nil {
				return err
			}

			if allCRDs {
				return generateAllCRDs(inputPath, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, indexerTypesPath, indexerImportPath, override, selectedGenerators)
			}

			return runGenerators(inputPath, crdKind, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, indexerTypesPath, indexerImportPath, override, selectedGenerators)
		},
	}

	// Build default generators list
	defaultGenerators := strings.Join(registry.List(), ",")

	rootCmd.Flags().StringVar(&inputPath, "input", "", "Path to a CRD yaml file (required)")
	rootCmd.Flags().StringVar(&crdKind, "crd", "", "CRD kind to generate controller for. Can not be set together with --all")
	rootCmd.Flags().BoolVar(&listCRDs, "list", false, "List available CRDs in the input file")
	rootCmd.Flags().BoolVar(&allCRDs, "all", false, "Generate controllers for all CRDs in the input file. Can not be set together with --crd")
	rootCmd.Flags().StringVar(&controllerOutDir, "controller-out", "", "Output directory for controller files (default: ../mongodb-atlas-kubernetes/internal/controller)")
	rootCmd.Flags().StringVar(&indexerOutDir, "indexer-out", "", "Output directory for indexer files (default: ../mongodb-atlas-kubernetes/internal/indexer)")
	rootCmd.Flags().StringVar(&exporterOutDir, "exporter-out", "", "Output directory for exporter files (default: ../mongodb-atlas-kubernetes/internal/exporter)")
	rootCmd.Flags().StringVar(&typesPath, "types-path", "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", "Full import path to the API types package")
	rootCmd.Flags().StringVar(&indexerTypesPath, "indexer-types-path", "", "Full import path for type imports in indexers (defaults to --types-path value)")
	rootCmd.Flags().StringVar(&indexerImportPath, "indexer-import-path", "", "Full import path for indexer imports in controllers (defaults to github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/generated/indexers)")
	rootCmd.Flags().BoolVar(&override, "override", false, "Override existing versioned handler files (default: false)")
	rootCmd.Flags().StringVar(&generators, "generators", defaultGenerators, fmt.Sprintf("Comma-separated list of generators to run (available: %s)", defaultGenerators))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

type CRDGenerationResult struct {
	CRDKind string
	Success bool
	Error   error
}

func validateGoImportPath(path string) error {
	if path == "" {
		return fmt.Errorf("import path cannot be empty")
	}
	// Basic validation: should not start with "/", should contain at least one "/"
	if path[0] == '/' {
		return fmt.Errorf("import path should not start with '/'")
	}
	if !strings.Contains(path, "/") {
		return fmt.Errorf("import path should contain at least one '/'")
	}
	// Should not contain spaces or other invalid characters
	if strings.ContainsAny(path, " \t\n\r") {
		return fmt.Errorf("import path should not contain whitespace")
	}
	return nil
}

// parseGenerators parses the comma-separated generator names and returns the corresponding generators.
func parseGenerators(generatorsFlag string) ([]registry.Generator, error) {
	if generatorsFlag == "" {
		return registry.All(), nil
	}

	names := strings.Split(generatorsFlag, ",")
	for i := range names {
		names[i] = strings.TrimSpace(names[i])
	}

	return registry.GetByNames(names)
}

// runGenerators runs the specified generators for a single CRD kind.
func runGenerators(inputPath, crdKind, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, indexerTypesPath, indexerImportPath string, override bool, gens []registry.Generator) error {
	opts := &registry.Options{
		InputPath:         inputPath,
		CRDKind:           crdKind,
		ControllerOutDir:  controllerOutDir,
		IndexerOutDir:     indexerOutDir,
		ExporterOutDir:    exporterOutDir,
		TypesPath:         typesPath,
		IndexerTypesPath:  indexerTypesPath,
		IndexerImportPath: indexerImportPath,
		Override:          override,
	}

	for _, gen := range gens {
		fmt.Printf("Running generator: %s\n", gen.Name())
		if err := gen.Generate(opts); err != nil {
			return fmt.Errorf("generator %s failed: %w", gen.Name(), err)
		}
	}

	return nil
}

func generateAllCRDs(inputPath, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, indexerTypesPath, indexerImportPath string, override bool, gens []registry.Generator) error {
	crds, err := config.ListCRDs(inputPath)
	if err != nil {
		return fmt.Errorf("failed to list CRDs: %w", err)
	}

	fmt.Printf("Found %d CRD(s) to generate\n", len(crds))
	fmt.Printf("Running generators: %s\n\n", generatorNames(gens))

	var results []CRDGenerationResult

	// Generate for each CRD
	for _, crd := range crds {
		fmt.Printf("Generating for CRD: %s...\n", crd.Kind)

		err = runGenerators(inputPath, crd.Kind, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, indexerTypesPath, indexerImportPath, override, gens)

		result := CRDGenerationResult{
			CRDKind: crd.Kind,
			Success: err == nil,
			Error:   err,
		}
		results = append(results, result)

		if err != nil {
			fmt.Printf("CRD: %s. Failed to generate because of an error: %v\n\n", crd.Kind, err)
		} else {
			fmt.Printf("CRD: %s. Generated successfully. No errors\n\n", crd.Kind)
		}
	}

	// Print summary
	fmt.Println("=== Generation Summary ===")
	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("OK  %s\n", result.CRDKind)
		} else {
			failureCount++
			fmt.Printf("ERR %s: %v\n", result.CRDKind, result.Error)
		}
	}

	fmt.Printf("\nTotal: %d, Success: %d, Failed: %d\n", len(results), successCount, failureCount)

	if failureCount > 0 {
		os.Exit(1)
	}

	return nil
}

// generatorNames returns a comma-separated list of generator names.
func generatorNames(gens []registry.Generator) string {
	names := make([]string, len(gens))
	for i, gen := range gens {
		names[i] = gen.Name()
	}
	return strings.Join(names, ", ")
}
