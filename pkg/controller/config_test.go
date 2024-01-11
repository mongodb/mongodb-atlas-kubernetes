package controller

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestParseDeletionProtection(t *testing.T) {
	t.Run("should use default values when flags or env vars were not set", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    true,
				SubObject: true,
			},
			config.DeletionProtection,
		)
	})

	t.Run("should do no action when flags where enabled", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app", "--object-deletion-protection", "--subobject-deletion-protection"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    true,
				SubObject: true,
			},
			config.DeletionProtection,
		)
	})

	t.Run("should do no action when flags where disabled", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app", "--object-deletion-protection=false", "--subobject-deletion-protection=false"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    false,
				SubObject: false,
			},
			config.DeletionProtection,
		)
	})

	//nolint:dupl
	t.Run("should use env vars when they are enabled and flags were not set", func(t *testing.T) {
		config := Config{}

		defer func(old []string) {
			os.Args = old
			assert.NoError(t, os.Unsetenv(objectDeletionProtectionEnvVar))
			assert.NoError(t, os.Unsetenv(subobjectDeletionProtectionEnvVar))
		}(os.Args)

		os.Args = []string{"app"}
		assert.NoError(t, os.Setenv(objectDeletionProtectionEnvVar, "true"))
		assert.NoError(t, os.Setenv(subobjectDeletionProtectionEnvVar, "true"))

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    true,
				SubObject: true,
			},
			config.DeletionProtection,
		)
	})

	//nolint:dupl
	t.Run("should use env vars when they are disabled and flags were not set", func(t *testing.T) {
		config := Config{}

		defer func(old []string) {
			os.Args = old
			assert.NoError(t, os.Unsetenv(objectDeletionProtectionEnvVar))
			assert.NoError(t, os.Unsetenv(subobjectDeletionProtectionEnvVar))
		}(os.Args)

		os.Args = []string{"app"}
		assert.NoError(t, os.Setenv(objectDeletionProtectionEnvVar, "false"))
		assert.NoError(t, os.Setenv(subobjectDeletionProtectionEnvVar, "false"))

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    false,
				SubObject: false,
			},
			config.DeletionProtection,
		)
	})

	t.Run("should use flags have precedence over env variables", func(t *testing.T) {
		config := Config{}

		defer func(old []string) {
			os.Args = old
			assert.NoError(t, os.Unsetenv(objectDeletionProtectionEnvVar))
			assert.NoError(t, os.Unsetenv(subobjectDeletionProtectionEnvVar))
		}(os.Args)

		os.Args = []string{"app", "--object-deletion-protection", "-subobject-deletion-protection"}
		assert.NoError(t, os.Setenv(objectDeletionProtectionEnvVar, "false"))
		assert.NoError(t, os.Setenv(subobjectDeletionProtectionEnvVar, "false"))

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.DeletionProtection.Object, objectDeletionProtectionFlag, ObjectDeletionProtectionDefault, "")
		flag.BoolVar(&config.DeletionProtection.SubObject, subobjectDeletionProtectionFlag, SubObjectDeletionProtectionDefault, "")
		flag.Parse()

		parseDeletionProtection(&config, flag.CommandLine)

		assert.Equal(
			t,
			DeletionProtection{
				Object:    true,
				SubObject: true,
			},
			config.DeletionProtection,
		)
	})
}

func TestAPISecretDefault(t *testing.T) {
	t.Run("should use default pod name and namespace", func(t *testing.T) {
		objectKey, err := APISecretDefault()
		assert.NoError(t, err)
		assert.Equal(
			t,
			client.ObjectKey{
				Namespace: "default",
				Name:      "mongodb-api-key",
			},
			objectKey,
		)
	})

	t.Run("should fail when operator pod name has wrong format", func(t *testing.T) {
		t.Setenv("OPERATOR_POD_NAME", "wrong_format")

		_, err := APISecretDefault()
		assert.ErrorContains(t, err, "the pod name must follow the format \"<deployment_name>-797f946f88-97f2q\" but got wrong_format")
	})

	t.Run("should use custom pod name and namespace", func(t *testing.T) {
		t.Setenv("OPERATOR_POD_NAME", "ako-797f946f88-97f2q")
		t.Setenv("OPERATOR_NAMESPACE", "atlas-operator")

		objectKey, err := APISecretDefault()
		assert.NoError(t, err)
		assert.Equal(
			t,
			client.ObjectKey{
				Namespace: "atlas-operator",
				Name:      "ako-api-key",
			},
			objectKey,
		)
	})
}

func TestApiSecret(t *testing.T) {
	t.Run("should use default namespace", func(t *testing.T) {
		cfg := Config{}
		err := apiSecret(&cfg)("custom-secret")
		assert.NoError(t, err)
		assert.Equal(
			t,
			client.ObjectKey{
				Namespace: "default",
				Name:      "custom-secret",
			},
			cfg.APISecret,
		)
	})

	t.Run("should fail when operator namespace doesn't exist", func(t *testing.T) {
		t.Setenv("OPERATOR_NAMESPACE", "atlas-operator")

		cfg := Config{}
		err := apiSecret(&cfg)("custom-secret")
		assert.NoError(t, err)
		assert.Equal(
			t,
			client.ObjectKey{
				Namespace: "atlas-operator",
				Name:      "custom-secret",
			},
			cfg.APISecret,
		)
	})
}
