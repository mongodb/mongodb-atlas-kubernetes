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
//

package config

import (
	"fmt"

	"sigs.k8s.io/yaml"

	"github.com/mongodb/mongodb-atlas-kubernetes/tools/openapi2crd/pkg/apis/config/v1alpha1"
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
