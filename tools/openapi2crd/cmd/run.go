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

	readOnly = os.O_RDONLY
)

func initConfig() {
	viper.AutomaticEnv()
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
			fs := afero.NewOsFs()

			outputOptionValue := viper.GetString(outputOption)
			forceOptionValue := viper.GetBool(forceOption)
			fsExporter, err := exporter.New(fs, outputOptionValue, forceOptionValue)
			if err != nil {
				return err
			}

			err = fsExporter.Start()
			if err != nil {
				return fmt.Errorf("error starting the exporter: %w", err)
			}

			configPath := viper.GetString(configOption)
			file, err := fs.OpenFile(configPath, readOnly, 0o644)
			if err != nil {
				return fmt.Errorf("error opening the file %s: %w", configPath, err)
			}

			configData, err := afero.ReadAll(file)
			if err != nil {
				return fmt.Errorf("error reading the file %s: %w", configPath, err)
			}

			cfg, err := config.Parse(configData)
			if err != nil {
				return fmt.Errorf("error parsing config: %w", err)
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

			for _, crdConfig := range cfg.Spec.CRDConfig {
				pluginSet, err := plugins.GetPluginSet(pluginSets, crdConfig.PluginSet)
				if err != nil {
					return fmt.Errorf("error getting plugin set %q: %w", crdConfig.PluginSet, err)
				}

				g := generator.NewGenerator(definitionsMap, pluginSet, openapiLoader)
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
		},
	}

	cmd.Flags().StringP(outputOption, "o", "", "Path to output file (required)")
	_ = cmd.MarkFlagRequired(outputOption)
	cmd.Flags().StringP(configOption, "c", "", "Path to the config file (required)")
	cmd.Flags().BoolP(forceOption, "f", false, "Force overwrite the output file if it exists")
	cobra.OnInitialize(initConfig)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}
