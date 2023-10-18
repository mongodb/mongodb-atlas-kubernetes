package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_configureDeletionProtection(t *testing.T) {
	t.Run("should do no action when config is nil", func(t *testing.T) {
		var config *Config
		configureDeletionProtection(config)

		assert.Nil(t, config)
	})

	t.Run("should use default values when flags or env vars were not set", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: true,
			},
			config,
		)
	})

	t.Run("should do no action when flags where enabled", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app", "--object-deletion-protection", "-subobject-deletion-protection"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: true,
			},
			config,
		)
	})

	t.Run("should do no action when flags where disabled", func(t *testing.T) {
		config := Config{}

		defer func(old []string) { os.Args = old }(os.Args)

		os.Args = []string{"app", "--object-deletion-protection=false", "-subobject-deletion-protection=false"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    false,
				SubObjectDeletionProtection: false,
			},
			config,
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
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: true,
			},
			config,
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
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    false,
				SubObjectDeletionProtection: false,
			},
			config,
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
		flag.BoolVar(&config.ObjectDeletionProtection, objectDeletionProtectionFlag, objectDeletionProtectionDefault, "")
		flag.BoolVar(&config.SubObjectDeletionProtection, subobjectDeletionProtectionFlag, subobjectDeletionProtectionDefault, "")
		flag.Parse()

		configureDeletionProtection(&config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: true,
			},
			config,
		)
	})
}
