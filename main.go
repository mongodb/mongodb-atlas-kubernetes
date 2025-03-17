package main

import (
	"context"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"os"
	"strings"

	"fybrik.io/openapi2crd/pkg/config"
	"fybrik.io/openapi2crd/pkg/exporter"
	"fybrik.io/openapi2crd/pkg/generator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	outputOption       = "output"
	resourcesOption    = "input"
	gvkOption          = "gvk"
	pathOption         = "path"
	verbOption         = "verb"
	categoryOption     = "category"
	majorVersionOption = "major-version"
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
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			specOptionValue := args[0]

			openapiSpec, err := config.LoadOpenAPI(specOptionValue)
			if err != nil {
				return err
			}

			// Fix known types (ref: https://github.com/kubernetes/kubernetes/issues/62329)
			openapiSpec.Components.Schemas["k8s.io/apimachinery/pkg/util/intstr.IntOrString"] = openapi3.NewSchemaRef("", &openapi3.Schema{
				AnyOf: openapi3.SchemaRefs{
					{
						Value: openapi3.NewStringSchema(),
					},
					{
						Value: openapi3.NewIntegerSchema(),
					},
				},
			})

			gvks := viper.GetStringSlice(gvkOption)
			paths := viper.GetStringSlice(pathOption)
			verbs := viper.GetStringSlice(verbOption)
			categories := viper.GetStringSlice(categoryOption)
			majorVersion := viper.GetString(majorVersionOption)

			if len(gvks) != len(paths) {
				return fmt.Errorf("number of gvks does not match the number of paths")
			}

			if len(verbs) != len(paths) {
				return fmt.Errorf("number of gvks does not match the number of paths")
			}

			outputOptionValue := viper.GetString(outputOption)
			exporter, err := exporter.New(outputOptionValue)
			if err != nil {
				return err
			}

			ctx := context.Background()
			for i, gvk := range gvks {
				g := generator.NewGenerator(majorVersion, paths[i], verbs[i], gvk, categories, openapiSpec)
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
	cmd.Flags().StringP(resourcesOption, "i", "", "Path to a directory with CustomResourceDefinition YAML files (required unless -g is used)")
	cmd.Flags().StringP(majorVersionOption, "m", "", "The CRD major version")
	cmd.Flags().StringSliceP(gvkOption, "g", []string{}, "The group/version/kind CRDs to create")
	cmd.Flags().StringSliceP(pathOption, "s", []string{}, "The OpenAPI paths to convert.")
	cmd.Flags().StringSliceP(verbOption, "v", []string{}, "The OpenAPI verbs to convert.")
	cmd.Flags().StringSliceP(categoryOption, "c", []string{}, "The CRD categories to apply.")

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
