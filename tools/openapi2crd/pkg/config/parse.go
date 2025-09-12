package config

import (
	"fmt"
	"github.com/mongodb/atlas2crd/pkg/apis/config/v1alpha1"
	"sigs.k8s.io/yaml"
)

func Parse(raw []byte) (*v1alpha1.Config, error) {
	cfg := v1alpha1.Config{}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config type: %w", err)
	}

	if cfg.Kind != "Config" || cfg.APIVersion != "atlas2crd.mongodb.com/v1alpha1" {
		return nil, fmt.Errorf("invalid config type: %s", cfg.Kind)
	}

	return &cfg, nil
}
