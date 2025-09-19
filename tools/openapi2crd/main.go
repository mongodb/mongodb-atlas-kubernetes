// Copyright 2025 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"tools/openapi2crd/pkg/config"
	"tools/openapi2crd/pkg/exporter"
	"tools/openapi2crd/pkg/generator"
)

const (
	outputOption = "output"
	configOption = "config"
)

// RootCmd defines the root cli command
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "openapi2crd SPEC_FILE",
		Short:         "Generate CustomResourceDefinition from OpenAPI 3.0 document",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			outputOptionValue := viper.GetString(outputOption)
			exporter, err := exporter.New(outputOptionValue)
			if err != nil {
				return err
			}

			configPath := viper.GetString(configOption)
			raw, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("error reading config: %w", err)
			}

			cfg, err := config.Parse(raw)
			if err != nil {
				return fmt.Errorf("error parsing config: %w", err)
			}

			ctx := context.Background()
			for _, crdConfig := range cfg.Spec.CRDConfig {
				g := generator.NewGenerator(crdConfig, cfg.Spec.OpenAPIDefinitions)
				crd, err := g.Generate(ctx)
				if err != nil {
					return err
				}

				err = exporter.Export(crd)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(configOption, "c", "", "Path to the config file (required)")
	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	// Run the cli
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
