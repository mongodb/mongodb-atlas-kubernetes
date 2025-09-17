package main

import (
	"fmt"
	"os"
	"tools/scaffolder/internal/generate"

	"github.com/spf13/cobra"
)

var (
	inputPath string
	crdKind   string
	listCRDs  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ako-controller-scaffolder",
		Short: "Generate Kubernetes controllers for MongoDB Atlas CRDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if inputPath == "" {
				return fmt.Errorf("--input is required")
			}

			if listCRDs {
				return generate.PrintCRDs(inputPath)
			}

			return generate.FromConfig(inputPath, crdKind)
		},
	}

	rootCmd.Flags().StringVar(&inputPath, "input", "", "Path to openapi2crd result.yaml file (required)")
	rootCmd.Flags().StringVar(&crdKind, "crd", "", "CRD kind to generate controller for")
	rootCmd.Flags().BoolVar(&listCRDs, "list", false, "List available CRDs from result file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
