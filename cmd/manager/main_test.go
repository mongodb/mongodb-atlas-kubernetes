package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mongodb/mongodb-atlas-kubernetes/pkg/version"
)

const (
	nonReleaseVersion = "1.8.0-30-g81233c6-dirty"
	releaseVersion    = "1.9.0-certified"
)

func TestDeletionProtectionDisabledByDefault(t *testing.T) {
	os.Unsetenv(objectDeletionProtectionEnvVar)
	os.Unsetenv(subobjectDeletionProtectionEnvVar)

	cfg := Config{
		ObjectDeletionProtection:    objectDeletionProtectionDefault,
		SubObjectDeletionProtection: subobjectDeletionProtectionDefault,
	}
	enableDeletionProtectionFromEnvVars(&cfg, version.DefaultVersion)

	assert.Equal(t, false, cfg.ObjectDeletionProtection)
	assert.Equal(t, false, cfg.SubObjectDeletionProtection)
}

func TestDeletionProtectionIgnoredOnReleases(t *testing.T) {
	version.Version = releaseVersion
	os.Setenv(objectDeletionProtectionEnvVar, "On")
	os.Setenv(subobjectDeletionProtectionEnvVar, "On")

	cfg := Config{
		ObjectDeletionProtection:    objectDeletionProtectionDefault,
		SubObjectDeletionProtection: subobjectDeletionProtectionDefault,
	}
	enableDeletionProtectionFromEnvVars(&cfg, releaseVersion)

	assert.Equal(t, false, cfg.ObjectDeletionProtection)
	assert.Equal(t, false, cfg.SubObjectDeletionProtection)
}

func TestDeletionProtectionEnabledAsEnvVars(t *testing.T) {
	testCases := []struct {
		title            string
		objDelProtect    bool
		subObjDelProtect bool
	}{
		{
			"both env vars set on non release version enables both protections",
			true,
			true,
		},
		{
			"obj env var set on non release version enables obj protection only",
			true,
			false,
		},
		{
			"subobj env var set on non release version enables subobj protection only",
			false,
			true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			os.Unsetenv(objectDeletionProtectionEnvVar)
			os.Unsetenv(subobjectDeletionProtectionEnvVar)
			if tc.objDelProtect {
				os.Setenv(objectDeletionProtectionEnvVar, "On")
			}
			if tc.subObjDelProtect {
				os.Setenv(subobjectDeletionProtectionEnvVar, "On")
			}

			cfg := Config{
				ObjectDeletionProtection:    objectDeletionProtectionDefault,
				SubObjectDeletionProtection: subobjectDeletionProtectionDefault,
			}
			enableDeletionProtectionFromEnvVars(&cfg, nonReleaseVersion)

			assert.Equal(t, tc.objDelProtect, cfg.ObjectDeletionProtection)
			assert.Equal(t, tc.subObjDelProtect, cfg.SubObjectDeletionProtection)
		})
	}
}
