package featureflags

import (
	"strings"
)

const (
	featurePrefix    = "FEATURE_"
	featureSeparator = "="
)

type FeatureFlags struct {
	features map[string]string
}

type EnvLister func() []string

// NewFeatureFlags creates a new instance of FeatureFlags and reads feature flags from the ENV
func NewFeatureFlags(envVarsLister EnvLister) *FeatureFlags {
	instance := &FeatureFlags{}
	envs := envVarsLister()
	result := map[string]string{}
	for _, e := range envs {
		if strings.HasPrefix(e, featurePrefix) {
			keyVal := strings.SplitN(e, featureSeparator, 2)
			if len(keyVal) == 2 {
				result[keyVal[0]] = keyVal[1]
			}
			result[e] = keyVal[0]
		}
	}
	instance.features = result
	return instance
}

func (f *FeatureFlags) IsFeaturePresent(featureName string) bool {
	_, ok := f.features[featureName]
	return ok
}

func (f *FeatureFlags) GetFeatureValue(featureName string) string {
	v, ok := f.features[featureName]
	if !ok {
		return ""
	}
	return v
}
