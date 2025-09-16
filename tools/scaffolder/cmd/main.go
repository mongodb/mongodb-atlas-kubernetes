package main

import (
	"fmt"
	"os"
	"tools/scaffolder/internal/generate"

	"github.com/spf13/cobra"
)

var (
	configPath string
	crdKind    string
	listCRDs   bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "ako-controller-scaffolder",
		Short: "Generate Kubernetes controllers for MongoDB Atlas CRDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if configPath == "" {
				return fmt.Errorf("--config is required")
			}

			if listCRDs {
				return generate.PrintCRDs(configPath)
			}

			return generate.FromConfig(configPath, crdKind)
		},
	}

	rootCmd.Flags().StringVar(&configPath, "config", "", "Path to atlas2crd config file (required)")
	rootCmd.Flags().StringVar(&crdKind, "crd", "", "CRD kind to generate controller for")
	rootCmd.Flags().BoolVar(&listCRDs, "list", false, "List available CRDs from config file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
