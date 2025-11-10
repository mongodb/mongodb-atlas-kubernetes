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

package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	configv1alpha1 "tools/openapi2crd/pkg/apis/config/v1alpha1"
	"tools/openapi2crd/pkg/config"
	"tools/openapi2crd/pkg/exporter"
	"tools/openapi2crd/pkg/generator"
	"tools/openapi2crd/pkg/plugins"
)

const (
	outputOption = "output"
	configOption = "config"
	forceOption  = "force"
	crdsOption   = "crds"

	crdsDefaultValue = "all"
	readOnly         = os.O_RDONLY
)

func initConfig() {
	viper.AutomaticEnv()
}

type RunnerConfig struct {
	Input     string
	Output    string
	Overwrite bool
	Kinds     map[string]struct{}
}

func RunCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "openapi2crd SPEC_FILE",
		Short:         "Generate CustomResourceDefinition from OpenAPI 3.0 document",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := viper.GetString(configOption)
			outputPath := viper.GetString(outputOption)
			forceOverwrite := viper.GetBool(forceOption)
			crds := viper.GetString(crdsOption)

			c := &RunnerConfig{
				Input:     configPath,
				Output:    outputPath,
				Overwrite: forceOverwrite,
				Kinds:     make(map[string]struct{}),
			}

			if crds != crdsDefaultValue {
				kinds := strings.Split(crds, ",")
				c.Kinds = make(map[string]struct{}, len(kinds))
				for _, kind := range kinds {
					c.Kinds[strings.TrimSpace(kind)] = struct{}{}
				}
			}

			fs := afero.NewOsFs()

			return runOpenapi2crd(ctx, fs, c)
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(configOption, "c", "", "Path to the config file (required)")
	_ = cmd.MarkFlagRequired(configOption)
	cmd.Flags().BoolP(forceOption, "f", false, "Force overwrite the output file if it exists")
	cmd.Flags().String(crdsOption, crdsDefaultValue, "One or more Kind names to generate, separated by comma. Use 'all' to generate all CRDs.")
	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}

func runOpenapi2crd(ctx context.Context, fs afero.Fs, runnerConfig *RunnerConfig) error {
	file, err := fs.OpenFile(runnerConfig.Input, readOnly, 0o644)
	if err != nil {
		return fmt.Errorf("error opening the file %s: %w", runnerConfig.Input, err)
	}

	configData, err := afero.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading the file %s: %w", runnerConfig.Input, err)
	}

	cfg, err := config.Parse(configData)
	if err != nil {
		return fmt.Errorf("error parsing config: %w", err)
	}

	fsExporter, err := exporter.New(fs, runnerConfig.Output, runnerConfig.Overwrite)
	if err != nil {
		return fmt.Errorf("error creating the exporter: %w", err)
	}

	err = fsExporter.Start()
	if err != nil {
		return fmt.Errorf("error starting the exporter: %w", err)
	}

	definitionsMap := map[string]configv1alpha1.OpenAPIDefinition{}
	for _, def := range cfg.Spec.OpenAPIDefinitions {
		definitionsMap[def.Name] = def
	}

	catalog := plugins.NewCatalog()
	pluginSets, err := catalog.BuildSets(cfg.Spec.PluginSets)
	if err != nil {
		return fmt.Errorf("error creating plugin set: %w", err)
	}

	openapiLoader := config.NewKinOpeAPI(fs)
	atlasLoader := config.NewAtlas(openapiLoader)

	for _, crdConfig := range cfg.Spec.CRDConfig {
		_, shouldGen := runnerConfig.Kinds[crdConfig.GVK.Kind]
		if len(runnerConfig.Kinds) > 0 && !shouldGen {
			continue
		}

		pluginSet, err := plugins.GetPluginSet(pluginSets, crdConfig.PluginSet)
		if err != nil {
			return fmt.Errorf("error getting plugin set %q: %w", crdConfig.PluginSet, err)
		}

		g := generator.NewGenerator(definitionsMap, pluginSet, openapiLoader, atlasLoader)
		crd, err := g.Generate(ctx, &crdConfig)
		if err != nil {
			return err
		}

		err = fsExporter.Export(crd)
		if err != nil {
			return err
		}
	}

	return fsExporter.Close()
}
