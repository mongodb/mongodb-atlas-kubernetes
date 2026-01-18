package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/scaffolder/internal/generate"

	"github.com/spf13/cobra"
)

var (
	inputPath        string
	crdKind          string
	listCRDs         bool
	allCRDs          bool
	controllerOutDir string
	indexerOutDir    string
	exporterOutDir   string
	typesPath        string
	override         bool
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
				return generate.PrintCRDs(inputPath)
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

			if allCRDs {
				return generateAllCRDs(inputPath, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, override)
			}

			return generate.FromConfig(inputPath, crdKind, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, override)
		},
	}

	rootCmd.Flags().StringVar(&inputPath, "input", "", "Path to a CRD yaml file (required)")
	rootCmd.Flags().StringVar(&crdKind, "crd", "", "CRD kind to generate controller for. Can not be set together with --all")
	rootCmd.Flags().BoolVar(&listCRDs, "list", false, "List available CRDs in the input file")
	rootCmd.Flags().BoolVar(&allCRDs, "all", false, "Generate controllers for all CRDs in the input file. Can not be set together with --crd")
	rootCmd.Flags().StringVar(&controllerOutDir, "controller-out", "", "Output directory for controller files (default: ../mongodb-atlas-kubernetes/internal/controller)")
	rootCmd.Flags().StringVar(&indexerOutDir, "indexer-out", "", "Output directory for indexer files (default: ../mongodb-atlas-kubernetes/internal/indexer)")
	rootCmd.Flags().StringVar(&exporterOutDir, "exporter-out", "", "Output directory for indexer files (default: ../mongodb-atlas-kubernetes/internal/exporter)")
	rootCmd.Flags().StringVar(&typesPath, "types-path", "github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/nextapi/generated/v1", "Full import path to the API types package")
	rootCmd.Flags().BoolVar(&override, "override", false, "Override existing versioned handler files (default: false)")

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

func generateAllCRDs(inputPath, controllerOutDir, indexerOutDir, exporterOutDir, typesPath string, override bool) error {
	crds, err := generate.ListCRDs(inputPath)
	if err != nil {
		return fmt.Errorf("failed to list CRDs: %w", err)
	}

	fmt.Printf("Found %d CRD(s) to generate\n\n", len(crds))

	var results []CRDGenerationResult

	// Generate for each CRD
	for _, crd := range crds {
		fmt.Printf("Generating for CRD: %s...\n", crd.Kind)

		err = generate.FromConfig(inputPath, crd.Kind, controllerOutDir, indexerOutDir, exporterOutDir, typesPath, override)

		result := CRDGenerationResult{
			CRDKind: crd.Kind,
			Success: err == nil,
			Error:   err,
		}
		results = append(results, result)

		if err != nil {
			fmt.Printf("CRD: %s. Failed to generate because of an error: %v\n\n", crd.Kind, err)
		} else {
			fmt.Printf("CRD: %s. Generated controller and indexer. No errors\n\n", crd.Kind)
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
