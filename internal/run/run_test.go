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

package run

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/mongodb/mongodb-atlas-kubernetes/v2/internal/featureflags"
)

func Test_configureDeletionProtection(t *testing.T) {
	t.Run("should do no action when config is nil", func(t *testing.T) {
		var config *Config
		configureDeletionProtection(flag.CommandLine, config)

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

		configureDeletionProtection(flag.CommandLine, &config)

		assert.Equal(
			t,
			Config{
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: false,
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

		configureDeletionProtection(flag.CommandLine, &config)

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

		configureDeletionProtection(flag.CommandLine, &config)

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

		configureDeletionProtection(flag.CommandLine, &config)

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

		configureDeletionProtection(flag.CommandLine, &config)

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

		configureDeletionProtection(flag.CommandLine, &config)

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

func TestInitCustomZapLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		wantLevel zapcore.Level
		wantErr   bool
	}{
		{
			name:      "valid string level info with json encoding",
			level:     "info",
			wantLevel: zapcore.InfoLevel,
			wantErr:   false,
		},
		{
			name:      "valid string level debug with console encoding",
			level:     "debug",
			wantLevel: zapcore.DebugLevel,
			wantErr:   false,
		},
		{
			name:      "valid numeric level",
			level:     "-1",
			wantLevel: zapcore.Level(-1),
			wantErr:   false,
		},
		{
			name:    "valid numeric level",
			level:   "-255",
			wantErr: true,
		},
		{
			name:    "invalid level",
			level:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := initCustomZapLogger(tt.level, "json")

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, logger)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)

			if logger == nil {
				return
			}
			// Verify logger configuration
			loggerImpl := logger.Core()
			assert.True(t, loggerImpl.Enabled(tt.wantLevel))
		})
	}
}

func TestParseConfiguration(t *testing.T) {
	os.Setenv("OPERATOR_NAMESPACE", "atlas-operator")
	os.Setenv("OPERATOR_POD_NAME", "podname-797f946f88-97f2q")
	for _, tc := range []struct {
		name    string
		args    []string
		want    Config
		wantErr string
	}{
		{
			name: "empty args default config",
			args: []string{},
			want: Config{
				AtlasDomain:          "https://cloud.mongodb.com/",
				EnableLeaderElection: false,
				MetricsAddr:          ":8080",
				WatchedNamespaces:    nil,
				ProbeAddr:            ":8081",
				GlobalAPISecret: client.ObjectKey{
					Namespace: "atlas-operator",
					Name:      "podname-api-key",
				},
				LogLevel:                    "info",
				LogEncoder:                  "json",
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: false,
				IndependentSyncPeriod:       15,
				FeatureFlags:                featureflags.NewFeatureFlags(os.Environ),
				DryRun:                      false,
			},
		},
		{
			name: "typical test args",
			args: []string{
				"--log-level=-9",
				"--global-api-secret-name=mongodb-atlas-operator-api-key",
				"--log-encoder=json",
				`--atlas-domain=https://cloud-qa.mongodb.com`,
			},
			want: Config{
				AtlasDomain:          "https://cloud-qa.mongodb.com",
				EnableLeaderElection: false,
				MetricsAddr:          ":8080",
				WatchedNamespaces:    nil,
				ProbeAddr:            ":8081",
				GlobalAPISecret: client.ObjectKey{
					Namespace: "atlas-operator",
					Name:      "mongodb-atlas-operator-api-key",
				},
				LogLevel:                    "-9",
				LogEncoder:                  "json",
				ObjectDeletionProtection:    true,
				SubObjectDeletionProtection: false,
				IndependentSyncPeriod:       15,
				FeatureFlags:                featureflags.NewFeatureFlags(os.Environ),
				DryRun:                      false,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			got, err := parseConfiguration(fs, tc.args)
			if tc.wantErr == "" {
				assert.Equal(t, tc.want, got)
			} else {
				assert.ErrorContains(t, err, tc.wantErr)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
