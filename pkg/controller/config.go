package controller

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DomainDefault                      = "https://cloud.mongodb.com/"
	ObjectDeletionProtectionDefault    = true
	SubObjectDeletionProtectionDefault = true
	operatorPodNameDefault             = "mongodb-atlas-operator"
	operatorNamespaceDefault           = "default"

	objectDeletionProtectionFlag      = "object-deletion-protection"
	subobjectDeletionProtectionFlag   = "subobject-deletion-protection"
	objectDeletionProtectionEnvVar    = "OBJECT_DELETION_PROTECTION"
	subobjectDeletionProtectionEnvVar = "SUBOBJECT_DELETION_PROTECTION"
)

type Config struct {
	Domain             string
	APISecret          client.ObjectKey
	DeletionProtection DeletionProtection
	FeatureFlags       *featureflags.FeatureFlags
}

type DeletionProtection struct {
	Object    bool
	SubObject bool
}

func (c *Config) ParseFlagsFromEnv(fs *flag.FlagSet) {
	if c == nil {
		return
	}

	parseDeletionProtection(c, fs)
	c.FeatureFlags = featureflags.NewFeatureFlags(os.Environ)
}

func DefaultConfig() (Config, error) {
	apiSecretDefault, err := APISecretDefault()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Domain:    DomainDefault,
		APISecret: apiSecretDefault,
		DeletionProtection: DeletionProtection{
			Object:    ObjectDeletionProtectionDefault,
			SubObject: SubObjectDeletionProtectionDefault,
		},
	}, nil
}

func RegisterFlags(conf *Config, fs *flag.FlagSet) {
	fs.StringVar(&conf.Domain, "atlas-domain", DomainDefault, "the Atlas URL domain name (with slash in the end).")
	fs.Func("global-api-secret-name", "The name of the Secret that contains Atlas API keys. "+
		"It is used by the Operator if AtlasProject configuration doesn't contain API key reference. Defaults to <deployment_name>-api-key.", apiSecret(conf))
	fs.BoolVar(&conf.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "Defines if the operator deletes Atlas resource "+
		"when a Custom Resource is deleted")
	fs.BoolVar(&conf.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "Defines if the operator overwrites "+
		"(and consequently delete) subresources that were not previously created by the operator")
}

func APISecretDefault() (client.ObjectKey, error) {
	parts := strings.Split(getOperatorPodName(), "-")
	if len(parts) <= 2 {
		return client.ObjectKey{}, fmt.Errorf("the pod name must follow the format \"<deployment_name>-797f946f88-97f2q\" but got %s", getOperatorPodName())
	}
	deploymentName := strings.Join(parts[0:len(parts)-2], "-")

	return client.ObjectKey{Namespace: GetOperatorNamespace(), Name: fmt.Sprintf("%s-api-key", deploymentName)}, nil
}

func GetOperatorNamespace() string {
	operatorNamespace := operatorNamespaceDefault
	if customNamespace, found := os.LookupEnv("OPERATOR_NAMESPACE"); found {
		operatorNamespace = customNamespace
	}

	return operatorNamespace
}

func apiSecret(conf *Config) func(string) error {
	return func(secretName string) error {
		conf.APISecret = client.ObjectKey{Namespace: GetOperatorNamespace(), Name: secretName}

		return nil
	}
}

func getOperatorPodName() string {
	operatorPodName := operatorPodNameDefault
	if customPodName, found := os.LookupEnv("OPERATOR_POD_NAME"); found {
		operatorPodName = customPodName
	}

	return operatorPodName
}

func parseDeletionProtection(c *Config, fs *flag.FlagSet) {
	objectDeletionSet := false
	subObjectDeletionSet := false

	fs.Visit(func(f *flag.Flag) {
		if f.Name == objectDeletionProtectionFlag {
			objectDeletionSet = true
		}

		if f.Name == subobjectDeletionProtectionFlag {
			subObjectDeletionSet = true
		}
	})

	if !objectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(objectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			c.DeletionProtection.Object = true
		case "false":
			c.DeletionProtection.Object = false
		default:
			c.DeletionProtection.Object = ObjectDeletionProtectionDefault
		}
	}

	if !subObjectDeletionSet {
		objDeletion := strings.ToLower(os.Getenv(subobjectDeletionProtectionEnvVar))
		switch objDeletion {
		case "true":
			c.DeletionProtection.SubObject = true
		case "false":
			c.DeletionProtection.SubObject = false
		default:
			c.DeletionProtection.SubObject = SubObjectDeletionProtectionDefault
		}
	}
}
