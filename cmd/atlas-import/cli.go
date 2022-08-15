package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mongodb/mongodb-atlas-kubernetes/cmd/atlas-import/importer"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// TODO add a debug flag and instantiate logger based on that

func generateBaseConfig(cmd *cobra.Command) importer.AtlasImportConfig {
	baseConfig := importer.AtlasImportConfig{
		AtlasDomain:     "https://cloud-qa.mongodb.com",
		ImportNamespace: "default",
	}
	orgID, publicKey, privateKey := getConnectionStrings(cmd)
	baseConfig.PublicKey = publicKey
	baseConfig.PrivateKey = privateKey
	baseConfig.OrgID = orgID
	domain, err := cmd.Flags().GetString(domainFlag)
	if err == nil && domain != "" {
		baseConfig.AtlasDomain = domain
	}
	namespace, err := cmd.Flags().GetString(namespaceFlag)
	if err == nil && domain != "" {
		baseConfig.ImportNamespace = namespace
	}
	return baseConfig
}

func getConnectionStrings(cmd *cobra.Command) (string, string, string) {
	orgID, err := cmd.Flags().GetString(orgFlag)
	if err != nil {
		importer.Log.Error(err.Error())
		importer.Log.Fatal("Please provide following arg : " + orgFlag)
	}
	publicKey, err := cmd.Flags().GetString(publicKeyFlag)
	if err != nil {
		importer.Log.Error(err.Error())
		importer.Log.Fatal("Please provide following arg : " + publicKeyFlag)
	}
	privateKey, err := cmd.Flags().GetString(privateKeyFlag)
	if err != nil {
		importer.Log.Error(err.Error())
		importer.Log.Fatal("Please provide following arg : " + privateKeyFlag)
	}
	return orgID, publicKey, privateKey
}

func parseYAMLFile(filename string) (*importer.AtlasImportConfig, error) {
	absPath := filename
	if !filepath.IsAbs(filename) {
		abs, err := filepath.Abs(filename)
		if err != nil {
			importer.Log.Error("Error with file path, try specifying the absolute path")
			return nil, err
		}
		absPath = abs
	}
	yamlFile, err := os.ReadFile(filepath.Clean(absPath))
	if err != nil {
		importer.Log.Error("Error reading file " + filename + " : " + err.Error())
		return nil, err
	}
	var config importer.AtlasImportConfig
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		importer.Log.Error("Error parsing YAML from file : " + err.Error())
		return nil, err
	}
	return &config, nil
}

func main() {
	importer.Log, _ = zap.NewDevelopment()
	Execute()
}

func run(config importer.AtlasImportConfig) {
	err := importer.RunImports(config)
	if err != nil {
		importer.Log.Fatal(err.Error())
	}
}

var rootCmd = &cobra.Command{
	Use:   "atlas-import",
	Short: "CLI tool to import your atlas resources into Kubernetes",
	Long: `CLI tool to import your atlas resources into Kubernetes. This tool allows you
to either use a configuration file, or CLI arguments.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You called atlas-import")
	},
}

var fromConfigCmd = &cobra.Command{
	Use:     "from-config",
	Aliases: []string{"conf"},
	Short:   "Import resources using a config file",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]
		importer.Log.Info("Running imports with config file : " + configFile)
		config, err := parseYAMLFile(configFile)
		if err != nil {
			importer.Log.Error(err.Error())
			importer.Log.Fatal("Couldn't read provided configuration file")
		}
		run(*config)
	},
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Import an Atlas Project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		importProjectID := args[0]
		importAllDeployments, _ := cmd.Flags().GetBool(allFlag)
		deploymentsList, _ := cmd.Flags().GetStringSlice(deploymentsFlag)
		importer.Log.Info("Importing project with ID : " + importProjectID)
		config := generateBaseConfig(cmd)
		config.ImportAll = false
		config.ImportedProjects = []importer.ImportedProject{
			{
				ID:          importProjectID,
				ImportAll:   importAllDeployments,
				Deployments: deploymentsList,
			},
		}
		run(config)
	},
}

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Import all your Atlas resources in Kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		importer.Log.Info("Importing all resources")
		config := generateBaseConfig(cmd)
		config.ImportAll = true
		run(config)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const orgFlag = "org"
const publicKeyFlag = "publickey"
const privateKeyFlag = "privatekey"

const namespaceFlag = "import-namespace"
const domainFlag = "domain"

const allFlag = "all"
const deploymentsFlag = "deployments"

func init() {
	rootCmd.AddCommand(fromConfigCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(allCmd)

	//TODO replace authentication method, can store a secret in cluster like for operator
	rootCmd.PersistentFlags().String(orgFlag, "", "Your Atlas organization ID")
	rootCmd.PersistentFlags().String(publicKeyFlag, "", "Your Atlas organization public key")
	rootCmd.PersistentFlags().String(privateKeyFlag, "", "Your Atlas organization private key")

	rootCmd.PersistentFlags().String(namespaceFlag, "", "Kubernetes namespace in which to instantiate resources")
	rootCmd.PersistentFlags().String(domainFlag, "", "Atlas domain name")

	//Will be true if --all is added, and false otherwise
	projectCmd.PersistentFlags().Bool(allFlag, false, "Import all Deployments for given project")
	//Deployments should be specified as : --deployments="dep1,dep2,dep3"
	projectCmd.PersistentFlags().StringSlice(deploymentsFlag, []string{}, "List of Deployments ID to import for given project")
}
