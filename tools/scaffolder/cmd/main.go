package main

import (
	"fmt"
	"os"
	"tools/scaffolder/internal/generate"

	"github.com/spf13/cobra"
)

var (
	inputPath         string
	crdKind           string
	listCRDs          bool
	controllerOutDir  string
	translationOutDir string
	indexerOutDir     string
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

			return generate.FromConfig(inputPath, crdKind, controllerOutDir, translationOutDir, indexerOutDir)
		},
	}

	rootCmd.Flags().StringVar(&inputPath, "input", "", "Path to openapi2crd result.yaml file (required)")
	rootCmd.Flags().StringVar(&crdKind, "crd", "", "CRD kind to generate controller for")
	rootCmd.Flags().BoolVar(&listCRDs, "list", false, "List available CRDs from result file")
	rootCmd.Flags().StringVar(&controllerOutDir, "controller-out", "", "Output directory for controller files (default: ../mongodb-atlas-kubernetes/internal/controller)")
	rootCmd.Flags().StringVar(&translationOutDir, "translation-out", "", "Output directory for translation files (default: ../mongodb-atlas-kubernetes/internal/translation)")
	rootCmd.Flags().StringVar(&indexerOutDir, "indexer-out", "", "Output directory for indexer files (default: ../mongodb-atlas-kubernetes/internal/indexer)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
